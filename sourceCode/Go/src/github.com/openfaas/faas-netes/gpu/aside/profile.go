//@file: profile.go
//@author: Yanan Yang
//@date: 2020/5/6
//@note: for function resource usage profiling, run as an aside thread when funcProfileCache misss
package aside

import (
	"fmt"
	"github.com/openfaas/faas-netes/gpu/metrics"
	"github.com/openfaas/faas-netes/gpu/repository"
	gTypes "github.com/openfaas/faas-netes/gpu/types"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ProfileFunc(funcName string, profTimeWindow time.Duration) {
	go func(){
		fmt.Printf("start to profile function %s, timeWindow = %d s \n", funcName, profTimeWindow/1000000000)
		// waiting prometheus to collect function cpu usage
		time.Sleep(profTimeWindow)
		MaxFuncCpuUsage := float64(-1)
		prometheusQuery := metrics.NewPrometheusQuery("192.168.0.111", 31113, &http.Client{})
		expr := url.QueryEscape(`sum(rate(container_cpu_usage_seconds_total{namespace=~"//FuncCpuCoreBind"}[1m])) by (pod)`)
		for { // while
			results, fetchErr := prometheusQuery.Fetch(expr)
			if fetchErr != nil {
				log.Printf("Error querying Prometheus API: %s\n", fetchErr.Error())
				time.Sleep(time.Second*3)
				continue
			}
			for _, v := range results.Data.Result {
				podNameSplit := strings.Split(v.Metric.PodName, "-")
				if podNameSplit == nil || len(podNameSplit)== 0 {
					log.Printf("profile: pod profile is nil, sleep 2s and continue \n")
					continue
				}
				if podNameSplit[0] == funcName {
					funcCpuUsage, strconvErr := strconv.ParseFloat(v.Value[1].(string),32)
					if strconvErr != nil {
						log.Printf("profile: unable to convert value for metric: %s", strconvErr)
						continue
					}
					if MaxFuncCpuUsage < funcCpuUsage {
						MaxFuncCpuUsage = funcCpuUsage
					}
				}
			}
			if MaxFuncCpuUsage > 0 {
				repository.UpdateFuncProfileCache(&gTypes.FuncProfile{
					FunctionName: funcName,
					MaxCpuCoreUsage: MaxFuncCpuUsage,
				})
				break
			} else {
				time.Sleep(time.Second*3)
				continue
			}
		}
		log.Printf("end to profile function %s, podCpuAvgUsage[1m]= %f \n", funcName, MaxFuncCpuUsage )
	}()
}