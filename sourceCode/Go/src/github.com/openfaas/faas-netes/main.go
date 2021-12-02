// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

//package main
//
//import (
//	"flag"
//	cpuRepository "github.com/openfaas/faas-netes/cpu/repository"
//	"github.com/openfaas/faas-netes/gpu/aside"
//	 "github.com/openfaas/faas-netes/gpu/repository"
//	"log"
//	"os"
//	"time"
//
//	"github.com/openfaas/faas-provider/proxy"
//	"k8s.io/client-go/kubernetes"
//
//	"github.com/openfaas-incubator/openfaas-operator/pkg/signals"
//	"github.com/openfaas/faas-netes/handlers"
//	"github.com/openfaas/faas-netes/k8s"
//	"github.com/openfaas/faas-netes/types"
//	"github.com/openfaas/faas-netes/version"
//	bootstrap "github.com/openfaas/faas-provider"
//	"github.com/openfaas/faas-provider/logs"
//	bootTypes "github.com/openfaas/faas-provider/types"
//	kubeinformers "k8s.io/client-go/informers"
//
//	"k8s.io/client-go/tools/clientcmd"
//)
//
//func main() {
//
//	repository.Init()
//
//	var kubeconfig string
//	var masterURL string
//
//	flag.StringVar(&kubeconfig, "kubeconfig", "",
//		"Path to a kubeconfig. Only required if out-of-cluster.")
//	flag.StringVar(&masterURL, "master", "",
//		"The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
//	flag.Parse()
//
//	clientCmdConfig, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
//	if err != nil {
//		log.Fatalf("Error building kubeconfig: %s", err.Error())
//	}
//
//	clientset, err := kubernetes.NewForConfig(clientCmdConfig)
//	if err != nil {
//		log.Fatalf("Error building Kubernetes clientset: %s", err.Error())
//
//	}
//
//	functionNamespace := "default"
//
//	if namespace, exists := os.LookupEnv("function_namespace"); exists {
//		functionNamespace = namespace
//	}
//
//	readConfig := types.ReadConfig{}
//	osEnv := bootTypes.OsEnv{}
//	cfg, err := readConfig.Read(osEnv)
//	if err != nil {
//		log.Fatalf("Error reading config: %s", err.Error())
//	}
//
//	log.Printf("HTTP Read Timeout: %s\n", cfg.FaaSConfig.GetReadTimeout())
//	log.Printf("HTTP Write Timeout: %s\n", cfg.FaaSConfig.WriteTimeout)
//	log.Printf("HTTPProbe: %v\n", cfg.HTTPProbe)
//	log.Printf("SetNonRootUser: %v\n", cfg.SetNonRootUser)
//
//	deployConfig := k8s.DeploymentConfig {
//		RuntimeHTTPPort: 8080,
//		HTTPProbe:       cfg.HTTPProbe,
//		SetNonRootUser:  cfg.SetNonRootUser,
//		ReadinessProbe: &k8s.ProbeConfig{
//			InitialDelaySeconds: int32(cfg.ReadinessProbeInitialDelaySeconds),
//			TimeoutSeconds:      int32(cfg.ReadinessProbeTimeoutSeconds),
//			PeriodSeconds:       int32(cfg.ReadinessProbePeriodSeconds),
//		},
//		LivenessProbe: &k8s.ProbeConfig{
//			InitialDelaySeconds: int32(cfg.LivenessProbeInitialDelaySeconds),
//			TimeoutSeconds:      int32(cfg.LivenessProbeTimeoutSeconds),
//			PeriodSeconds:       int32(cfg.LivenessProbePeriodSeconds),
//		},
//		ImagePullPolicy: cfg.ImagePullPolicy,
//	}
//
//	factory := k8s.NewFunctionFactory(clientset, deployConfig)
//
//	defaultResync := time.Second * 5
//	kubeInformerOpt := kubeinformers.WithNamespace(functionNamespace)
//	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(clientset, defaultResync, kubeInformerOpt)
//
//	// set up signals so we handle the first shutdown signal gracefully
//	stopCh := signals.SetupSignalHandler()
//
//	endpointsInformer := kubeInformerFactory.Core().V1().Endpoints()
//	go kubeInformerFactory.Start(stopCh)
//	lister := endpointsInformer.Lister()
//
//	functionLookup := k8s.NewFunctionLookup(functionNamespace, lister)
//
//	bootstrapHandlers := bootTypes.FaaSHandlers {
//		FunctionProxy:        proxy.NewHandlerFunc(cfg.FaaSConfig, functionLookup),
//		DeleteHandler:        handlers.MakeDeleteHandler(functionNamespace, clientset),
//		DeployHandler:        handlers.MakeDeployHandler(functionNamespace, factory, clientset),
//		FunctionReader:       handlers.MakeFunctionReader(functionNamespace, clientset),
//		ReplicaReader:        handlers.MakeReplicaReader(functionNamespace, clientset),
//		ReplicaUpdater:       handlers.MakeReplicaUpdater(functionNamespace, clientset),
//		UpdateHandler:        handlers.MakeUpdateHandler(functionNamespace, factory),
//		HealthHandler:        handlers.MakeHealthHandler(),
//		InfoHandler:          handlers.MakeInfoHandler(version.BuildVersion(), version.GitCommit),
//		SecretHandler:        handlers.MakeSecretHandler(functionNamespace, clientset),
//		LogHandler:           logs.NewLogHandlerFunc(k8s.NewLogRequestor(clientset, functionNamespace), cfg.FaaSConfig.WriteTimeout),
//		ListNamespaceHandler: handlers.MakeNamespacesLister(functionNamespace, clientset),
//	}
//	go aside.RpsDispatcherMonitor(functionNamespace, cfg.LoadGenHost, cfg.LoadGenPort)
//	cpuRepository.InitializeCluster(clientset)
//	bootstrap.Serve(&bootstrapHandlers, &cfg.FaaSConfig)
//}


package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/openfaas/faas-netes/gpu/repository"
	ptypes "github.com/openfaas/faas-provider/types"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"sync"
	"time"

	//"github.com/openfaas/faas-netes/gpu/aside"
	scheduler "github.com/openfaas/faas-netes/gpu/controller"
	//"github.com/openfaas/faas-netes/gpu/repository"
	"github.com/openfaas/faas-netes/gpu/tools"
	gpuTypes "github.com/openfaas/faas-netes/gpu/types"
	"github.com/openfaas/faas-netes/k8s"
	types "github.com/openfaas/faas-netes/types"
	bootTypes "github.com/openfaas/faas-provider/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	//"time"

	//"time"
)

func main() {
	time.Local = time.FixedZone("CST", 3600*8)
	log.SetFlags(log.Lmicroseconds)
	repository.Init()
	scheduler.InitProfiler()

	var kubeconfig string
	var masterURL string

	flag.StringVar(&kubeconfig, "kubeconfig", "C:/Go/workspace/k8s/config",
		"Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "",
		"The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()

	clientCmdConfig, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(clientCmdConfig)
	if err != nil {
		panic(err.Error())
	}
	readConfig := types.ReadConfig{}
	osEnv := bootTypes.OsEnv{}
	cfg, err := readConfig.Read(osEnv)
	if err != nil {
		log.Fatalf("Error reading config: %s", err.Error())
	}

	cfg.HTTPProbe = true

	deployConfig := k8s.DeploymentConfig {
		RuntimeHTTPPort: 8080,
		HTTPProbe:       cfg.HTTPProbe,
		SetNonRootUser:  cfg.SetNonRootUser,
		ReadinessProbe: &k8s.ProbeConfig {
			InitialDelaySeconds: int32(cfg.ReadinessProbeInitialDelaySeconds),
			TimeoutSeconds:      int32(cfg.ReadinessProbeTimeoutSeconds),
			PeriodSeconds:       int32(cfg.ReadinessProbePeriodSeconds),
		},
		LivenessProbe: &k8s.ProbeConfig {
			InitialDelaySeconds: int32(cfg.LivenessProbeInitialDelaySeconds),
			TimeoutSeconds:      int32(cfg.LivenessProbeTimeoutSeconds),
			PeriodSeconds:       int32(cfg.LivenessProbePeriodSeconds),
		},
		ImagePullPolicy: cfg.ImagePullPolicy,
	}

	factory := k8s.NewFunctionFactory(clientset, deployConfig)
	log.Println("factory:", factory)

	//testEstimator(9999)
	//scheduler.InferResourceConfigsWithBatch("resnet-50",350.0, 1,200)

	//testScaleIn()
	//testScheRSWA()
	//testScheDRP2()
	//testEstimator(850)
	//repository.RegisterFuncDeploy("resnet50")
	//repository.AddVirtualFuncPodConfig("resnet50")
	//repository.AddFuncPodConfig("resnet50", &gpuTypes.FuncPodConfig {
	//	FuncPodType:          "p",
	//	FuncPodName:          "p",
	//	BatchSize:            0,
	//	CpuThreads:           0,
	//	GpuCorePercent:       0,
	//	GpuMemoryRate:        0,
	//	ExecutionTime:        0,
	//	ReqPerSecondLottery:  0,
	//	ReqPerSecondMax:      20,
	//	ReqPerSecondMin:      1,
	//	InactiveCounter:      0,
	//	FuncPodIp:            "",
	//	NodeGpuCpuAllocation: &gpuTypes.NodeGpuCpuAllocation{
	//		NodeTh:        0,
	//		CudaDeviceTh:  0,
	//		SocketTh:      0,
	//		CpuCoreThList: []int{0},
	//	},
	//})
	//for i:=0;i<10;i++{
	//	repository.AddFuncPodConfig("resnet50", &gpuTypes.FuncPodConfig{
	//		FuncPodType:          "i",
	//		FuncPodName:          "pod"+strconv.Itoa(i),
	//		BatchSize:            0,
	//		CpuThreads:           0,
	//		GpuCorePercent:       0,
	//		GpuMemoryRate:        0,
	//		ExecutionTime:        0,
	//		ReqPerSecondLottery:  0,
	//		ReqPerSecondMax:      int32(rand.Intn(20)+10),
	//		ReqPerSecondMin:      int32(rand.Intn(8)),
	//		InactiveCounter:      0,
	//		FuncPodIp:            "",
	//		NodeGpuCpuAllocation: &gpuTypes.NodeGpuCpuAllocation{
	//			NodeTh:        0,
	//			CudaDeviceTh:  0,
	//			SocketTh:      0,
	//			CpuCoreThList: []int{0},
	//		},
	//	})
	//}
	//
	//go aside.RpsDispatcherMonitor("","192.168.1.109",8080)
	//
	//time.Sleep(time.Second*40)


	//repository.AddFuncPodConfig("test", &gpuTypes.FuncPodConfig{
	//	FuncPodType:          "i",
	//	FuncPodName:          "i",
	//	BatchSize:            0,
	//	CpuThreads:           0,
	//	GpuCorePercent:       0,
	//	GpuMemoryRate:        0,
	//	ExecutionTime:        0,
	//	ReqPerSecondLottery:  0,
	//	ReqPerSecondMax:      10,
	//	ReqPerSecondMin:      0,
	//	InactiveCounter:      0,
	//	FuncPodIp:            "",
	//	NodeGpuCpuAllocation: &gpuTypes.NodeGpuCpuAllocation{
	//		NodeTh:        0,
	//		CudaDeviceTh:  0,
	//		SocketTh:      0,
	//		CpuCoreThList: []int{0},
	//	},
	//})
	//time.Sleep(time.Second*4)
	//
	//for  j:=0; j<10; j++ {
	//	log.Println()
	//	len := 1000
	//	w := 0
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	for i := 0; i < len; i++ {
	//		//time.Sleep(time.Millisecond * 10)
	//		counter := int32(0)
	//		funcc := repository.GetFunc("test")
	//
	//		var keys []string
	//		for _, v := range funcc.FuncPodConfigMap {
	//			if v.FuncPodType != "p" {
	//				keys = append(keys, v.FuncPodName)
	//			}
	//
	//		}
	//		sort.Strings(keys)
	//
	//
	//		winner := r.Intn(int(funcc.FuncPodTotalLottery))
	//		//winner := w%(int(6)+int(vpodLottery))
	//		w++
	//
	//		for _, v := range keys {
	//			counter = counter + funcc.FuncPodConfigMap[v].ReqPerSecondLottery
	//			if counter > int32(winner) {
	//				//log.Println("w==",w,"--->",funcc.FuncPodConfigMap[v].FuncPodType,"counter=",counter, "winner=",winner)
	//				funcc.FuncPodConfigMap[v].ReqPerSecondMin++
	//				break
	//			} else {
	//			}
	//		}
	//		//counter = 0
	//		//for _, item := range List {
	//		//	if item.PodName != "p" {
	//		//		counter = counter + item.Lottary
	//		//		if counter > int32(winner) {
	//		//			item.MinReq++
	//		//			break
	//		//		}
	//		//	}
	//		//}
	//	}
	//	for _, v := range repository.GetFunc("test").FuncPodConfigMap {
	//		log.Println(v.ReqPerSecondLottery, " -->", float32(v.ReqPerSecondMin), " w=", w)
	//		v.ReqPerSecondMin=0
	//	}
	//	//for _, item := range List {
	//	//	fmt.Println(item.Lottary," -->",float32(item.MinReq))
	//	//}
	//
	//}
	//time.Sleep(time.Second*40)
	//prometheusQuery := metrics.NewLoadGenQuery("192.168.1.109", 8080, &http.Client{})
	//
	//results, fetchErr := prometheusQuery.Fetch()
	//if fetchErr != nil {
	//	log.Printf("Error querying Prometheus API: %s \n", fetchErr.Error())
	//}
	//log.Printf("%+v\n",results)
	//for _ , item := range *results {
	//	log.Printf("%+v\n",item)
	//}
	//resolve()

	//var jsonBlob = []byte(`[
	//	{"Name": "Platypus", "Order": 0},
	//	{"Name": "Quoll",    "Order": 1}
	//]`)
	//log.Printf("%s\n",jsonBlob)
	//type Animal struct {
	//	Name  string
	//	Order string
	//}
	//var animals []Animal
	//
	//err = json.Unmarshal(jsonBlob, &animals)
	//if err != nil {
	//	fmt.Println("error:", err)
	//}
	//fmt.Printf("%+v", animals)
	/*for k := 0; k < 2 && neededCores > 0; k++ {
		if sort.SearchInts(cpuCoreThList, k) == k {

		}
		cpuCoreThList = append(cpuCoreThList, k)
		neededCores--
	}*/


	//aside.RpsDispatcherMonitor()
	//testScheduler()
	/* for i:=0;i<300;i++{
		    a:= rand.Intn(700)+50
		    a++
	    	if i>200 {
			    fmt.Print(rand.Intn(700)+50,",")
		    }
		    //265,237,260,540,348,441,734,217,421,276,131,216,343,169,102,735,237,478,634,74,662,666,290,201,390,594,255,151,52,417,128,72,92,393,716,685,454,107,121,580,350,670,220,397,415,179,598,145,350,579,481,736,549,97,738,53,68,69,407,225,603,67,320,336,270,572,411,570,704,610,531,66,389,683,454,159,269,290,690,534,185,271,624,620,88,494,124,374,556,102,472,167,515,229,79,71,561,303,309,
		    //481,736,549,97,738,53,68,69,407,225,603,67,320,336,270,572,411,570,704,610,531,66,389,683,454,159,269,290,690,534,185,271,624,620,88,494,124,374,556,102,472,167,515,229,79,71,561,303,309,558,278,124,125,356,636,464,709,73,337,424,281,733,366,327,386,355,739,637,559,109,52,495,303,370,555,259,706,530,727,428,643,58,315,234,73,623,303,477,98,81,356,581,227,604,280,137,654,343,106,
	    }*/
	// MIN 为用户自定义的比较精度
	//supportBatchGroup :=[]int32{16,8,4,2,1}
	//repository.RegisterFuncDeploy("resnet50")
	//repository.UpdateFuncProfileCache(&gpuTypes.FuncProfile {
	//	FunctionName:    "resnet50",
	//	MaxCpuCoreUsage: 0.9,
	//	MinCpuCoreUsage: 0.1,
	//	AvgCpuCoreUsage: 0.44,
	//})
	//CreatePreWarmPod("resnet50", "", 400, 1, nil)
	//ScaleUpFuncCapacity("resnet50", "", 400, 800, supportBatchGroup, nil)


	//prometheusQuery := metrics.NewPrometheusQuery(prometheusHost, prometheusPort, &http.Client{})
	//expr := url.QueryEscape(`sum(rate(gateway_function_invocation_total{function_name=~".*", code=~".*", kubernetes_namespace="`+funcNamespace+`"}[10s])) by (function_name)`)
	//for i:=-50;i<150;i++ {
	//	fmt.Println(int32(math.Sin(float64(i)*math.Pi/100)*400+400))
	//}


	//testDeploy(clientset, factory)
	//
	////testReader("openfaas-fn","sleep", clientset)
	//testReaderList("test",clientset)
	//time.Sleep(time.Second*120)
	//testDelete("test","sleep",clientset)

	//time.Sleep(time.Second*4)

	//cpuRepository.InitializeCluster(clientset)

	/*fn := []string{"sleep","sleep2"}
	for i:=0; i< len(fn); i++ {
		aside.ProfileFunc(fn[i], time.Second*60)
	}
	time.Sleep(time.Second*100)*/
	//testReader("openfaas-fn","sleep", clientset)
	//testReaderList("openfaas-fn", clientset)
	/*alert.AlertMonitor(clientset, time.Second*10)

	testDeploy(clientset, factory)
	repository.PrintClusterCapConfig()
	repository.PrintFuncDeployStatusMap()
	time.Sleep(time.Second*60)
	repository.UpdateFuncExpectedReplicas("sleep",5)
	repository.PrintClusterCapConfig()
	repository.PrintFuncDeployStatusMap()

	testReplicas("test","sleep", clientset)
	repository.PrintClusterCapConfig()
	repository.PrintFuncDeployStatusMap()
	time.Sleep(time.Second*60)
	testReader("test","sleep", clientset)

	repository.UpdateFuncExpectedReplicas("sleep",2)
	repository.PrintClusterCapConfig()
	repository.PrintFuncDeployStatusMap()
	testReplicas("test","sleep", clientset)
	repository.PrintClusterCapConfig()
	repository.PrintFuncDeployStatusMap()
	time.Sleep(time.Second*600)
	//testReader("test","sleep", clientset)
	testDelete("test","sleep",clientset)
	repository.PrintClusterCapConfig()
	repository.PrintFuncDeployStatusMap()*/
	//Post("http://192.168.1.120:31212/function/resnet50", "{\"instances\": [1.0, 2.0, 5.0]}", "application/json")



}
var readLock sync.RWMutex
func updateFunctionsLottery(funcDeployStatus *gpuTypes.FuncDeployStatus){
	readLock.RLock()
	defer readLock.RUnlock()
	for i:=0;i<50;i++ {
		count := 0
		for _, v := range funcDeployStatus.FuncPodConfigMap {
			log.Println(v.FuncPodType)
			count ++
		}
		if count == 0 {
			log.Println("return")
			return
		}
		time.Sleep(time.Millisecond*100)
	}
}
func Get(url string) (response string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	defer resp.Body.Close()
	if error != nil {
		panic(error)
	}

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	response = result.String()
	return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func Post(url string, data interface{}, contentType string) (content string) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n",req.Body)
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}

func packing(){


	maxResourceQuotaNagDiffIndex := -1
	minResourceQuotaPosDiffIndex := -1
	maxResourceQuotaNagDiff := float64(-999)
	minResourceQuotaPosDiff := float64(999)
	tempGpuCoreQuota := float64(0)
	tempCpuQuota := float64(0)
	tempDiffQuota := float64(0)
	pickConfigIndex := -1

	type Config struct {
		CpuThreads int32
		GpuCorePercent int32
	}
	resourcesConfigs := []Config{}
	for i:=2; i<=16; i=i*2{
		for j:=10; j<=100; j=j+10 {
			resourcesConfigs = append(resourcesConfigs, Config{
				CpuThreads:     int32(i),
				GpuCorePercent: int32(j),
			})
		}
	}


	// cpu is dominantly remained resource
	for k := 0; k < len(resourcesConfigs); k++ {
		tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / 20
		tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / 100
		tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
		log.Printf("k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
		if gpuTypes.Greater(tempDiffQuota,0) {
			if gpuTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
				minResourceQuotaPosDiff = tempDiffQuota
				minResourceQuotaPosDiffIndex = k
			}
		} else if gpuTypes.Less(tempDiffQuota,0) {
			if gpuTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
				maxResourceQuotaNagDiff = tempDiffQuota
				maxResourceQuotaNagDiffIndex = k
			}
		}
	}
	if minResourceQuotaPosDiffIndex == -1 {
		pickConfigIndex = maxResourceQuotaNagDiffIndex
	} else {
		pickConfigIndex = minResourceQuotaPosDiffIndex
	}
	log.Printf("maxResourceQuotaNagDiffIndex=%d, maxResourceQuotaNagDiff=%f\n", maxResourceQuotaNagDiffIndex, maxResourceQuotaNagDiff)
	log.Printf("minResourceQuotaPosDiffIndex=%d, minResourceQuotaPosDiff=%f\n", minResourceQuotaPosDiffIndex, minResourceQuotaPosDiff)
	log.Printf("pickConfigIndex=%d\n", pickConfigIndex)


	maxResourceQuotaNagDiffIndex = -1
	minResourceQuotaPosDiffIndex = -1
	maxResourceQuotaNagDiff = -999
	minResourceQuotaPosDiff = 999
	tempGpuCoreQuota = 0
	tempCpuQuota = 0
	tempDiffQuota = 0
	pickConfigIndex = -1
	// cpu is dominantly remained resource
	for k := 0; k < len(resourcesConfigs); k++ {
		tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20)
		tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / 100
		tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
		log.Printf("k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
		if gpuTypes.Greater(tempDiffQuota,0) {
			if gpuTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
				minResourceQuotaPosDiff = tempDiffQuota
				minResourceQuotaPosDiffIndex = k
			}
		} else if gpuTypes.Less(tempDiffQuota,0) {
			if gpuTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
				maxResourceQuotaNagDiff = tempDiffQuota
				maxResourceQuotaNagDiffIndex = k
			}
		}
	}
	if minResourceQuotaPosDiffIndex == -1 {
		pickConfigIndex = maxResourceQuotaNagDiffIndex
	} else {
		pickConfigIndex = minResourceQuotaPosDiffIndex
	}
	log.Printf("maxResourceQuotaNagDiffIndex=%d, maxResourceQuotaNagDiff=%f\n", maxResourceQuotaNagDiffIndex, maxResourceQuotaNagDiff)
	log.Printf("minResourceQuotaPosDiffIndex=%d, minResourceQuotaPosDiff=%f\n", minResourceQuotaPosDiffIndex, minResourceQuotaPosDiff)
	log.Printf("pickConfigIndex=%d\n", pickConfigIndex)

}

func resolve (){
	type Config struct {
		Lottary int32
		MinReq int32
		MaxReq int32
		PodName string
	}
	type Func struct {
		MinReqCap int32
		MaxReqCap int32
		LottarySum int32
		ConfigMap map[string]*Config
	}
	funcConfig := map[string]*Config{}
	List :=[]*Config{}
	List = append(List,&Config{
		Lottary: 1,
		MinReq:  0,
		MaxReq:  0,
		PodName: "pod1",
	} )
	List = append(List,&Config{
		Lottary: 2,
		MinReq:  0,
		MaxReq:  0,
		PodName: "pod2",
	} )
	List = append(List,&Config{
		Lottary: 3,
		MinReq:  0,
		MaxReq:  0,
		PodName: "pod3",
	} )
	List = append(List,&Config{
		Lottary: 4,
		MinReq:  0,
		MaxReq:  0,
		PodName: "podv",
	} )

	funcConfig["pod1"]= &Config{
		Lottary: 1,
		MinReq:  0,
		MaxReq:  0,
		PodName: "pod1",
	}
	funcConfig["pod2"]= &Config{
		Lottary: 2,
		MinReq:  0,
		MaxReq:  0,
		PodName: "pod2",
	}
	funcConfig["pod3"]= &Config{
		Lottary: 3,
		MinReq:  0,
		MaxReq:  0,
		PodName: "pod3",
	}

	funcConfig["v"]= &Config{
		Lottary: 4,
		MinReq:  0,
		MaxReq:  0,
		PodName: "v",
	}
	start := time.Now()
	len := 100
	w:=0
	for i:=0; i<len; i++ {
		counter := int32(0)
		vpodLottery := int32(0)
		vpodConfig, exist := funcConfig["v"]
		if exist {
			vpodLottery = vpodConfig.Lottary
		}
		winner := rand.Intn(int(6)+int(vpodLottery))
		//winner := w%(int(6)+int(vpodLottery))
		w++
		for _, v := range funcConfig {
			v.PodName=v.PodName
		}
		for _, v := range funcConfig {
			if v.PodName != "p" {
				counter = counter + v.Lottary
				if counter > int32(winner) {
					v.MinReq++
					break
				}
			}
		}
		counter = 0
		for _, item := range List {
			if item.PodName != "p" {
				counter = counter + item.Lottary
				if counter > int32(winner) {
					item.MinReq++
					break
				}
			}
		}
	}
	fmt.Println(time.Since(start).Seconds())
	for _, v := range funcConfig {
		fmt.Println(v.Lottary," -->",float32(v.MinReq))
	}
	for _, item := range List {
		fmt.Println(item.Lottary," -->",float32(item.MinReq))
	}
}



/**
* init box
 */
type BoxConfig struct {
	CpuThreadsCap int32
	GpuCorePercentCap int32 // GPU cores capacity is 100
	GpuMemoryRateCap float64  //GPU memory capacity is 1.0
}
func InitBox(boxNum int) []*BoxConfig {
	var resourcesConfigs []*BoxConfig
	for i:=0; i<boxNum; i++{
		resourcesConfigs = append(resourcesConfigs, &BoxConfig{
			CpuThreadsCap:     20, //20
			GpuCorePercentCap: 100, //100
			GpuMemoryRateCap:  1.0, //1.0
		})
	}
	return resourcesConfigs
}
func testScheDRP(){

	/*latencySLO_100s := []int{231,437,597,309,231,68,375,90,306,450,544,
	161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,
	737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,
	340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,
	738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,
	526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,
	383,693,541,552,728,186,196,657,690,453,602,693,455,248}*/
	latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468,351,324,439,670,493,533,626,508,123,234,640,182,705,614,205,453,580,631,93,98,687,590,271,384,110,264,574,171,522,173,555,468,325,570,489,733,725,352,174,63,138,431,123,599,80,248,474,516,546,429,395,248,250,595,639,539,531,352,533,178,348,104,285,106,345,698,660,203,598,135,618,254,479,705,571,626,334,92,608,274,729,711,605,487,547,709,242,599,673,730,304,74,235,482,234,365,152,251,253,269,600,173,59,445,277,638,462,619,576,79,56,350,521,675,92,495,221,256,430,63,365,509,399,427,590,298,488,581,279,219,272,649,149,385,573,212,520,460,700,744,380,575,302,515,288,454,55,154,368,735,206,678,689,494,446,110,189,609,573,701,130,452,295,335,209,264,493,624,247,189,299,645,490,425,341,221,551,156,56,352,363,422,490,305,338,690,627,508,92,415,141,324,664,159,514,418,96,497,92,433,128,283,572,230,709,265,87,102,191,254,178,457,351,100,356,398,255,552,657,556,221,598,648,364,56,624,429,284,619,350,328,125,203,312,76,68,229,615,557,715,169,92,92,468,435,729,658,363,179,548,327,50,601,662,366,466,583,710,490,306,565,166,521,571,421,565,699,603,290,260,76,507,628,307,223,415,709,88,427,648,174,468,611,581,640,554,132,313,536,129,419,263,153,742,296,167,73,344,631,606,174,532,704,196,100,268,82,694,559,115,342,218,540,548,687,260,732,726,724,245,608,475,360,680,416,139,388,617,726,698,379,267,87,586,356,358,370,616,227,288,348,527,122,359,71,172,689,176,92,723,704,696,422,82,83,496,422,452,719,555,203,652,268,232,155,586,736,324,188,358,472,349,592,538,172,522,253,551,345,335,432,595,612,331,347,452,700,379,664,649,73,586,56,512,134,108,129,106,304,554,81,167,416,254,560,379,326,554,102,371,103,726,402,551,564,677,567,445,622,391,196,431,581,674,683,151,460,424,480,602,601,610,320,530,219,268,496,360,637,191,682,688,706,448,293,631,658,376,550,387,378,294,535,401,225,391,384,192,677,244,635,291,417,424,301,442,727,170,634,684,728,712,703,284,336,364,291,80,512,312,247,693,662,479,535,612,588,319,284,641,717,263,589,355,244,736,533,252,589,312,616,723,639,577,169,397,321,175,667,381,716,687,701,346,369,289,641,536,709,715,394,83,358,728,687,193,640,711,239,321,480,524,630,602,228,360,563,628,576,320,108,208,336,453,501,350,398,742,162,413,59,237,713,332,268,535,601,224,95,604,92,600,447,638,305,659,222,355,166,315,490,223,630,306,125,371,316,147,269,253,229,368,681,338,394,81,360,715,660,261,190,332,323,320,62,466,363,270,417,121,429,626,462,697,725,63,304,109,84,211,689,670,223,225,285,196,695,113,304,645,192,247,578,687,730,466,283,720,336,95,86,80,221,727,74,192,356,113,305,463,144,628,493,599,651,717,703,261,356,327,400,527,686,308,361,408,162,693,650,201,377,553,590,707,341,539,641,584,314,517,163,600,654,692,472,417,340,581,584,56,509,623,248,540,220,261,75,696,90,423,335,420,351,541,281,657,654,238,176,646,309,117,218,155,263,676,201,70,544,506,55,317,80,546,709,533,629,300,232,581,453,593,132,549,653,216,566,693,302,733,569,518,493,526,123,98,593,675,341,127,505,457,147,650,423,105,278,118,190,696,697,202,716,405,97,280,616,143,541,680,215,132,153,308,123,413,459,476,438,230,538,201,678,511,749,288,716,571,293,214,550,259,705,518,93,707,100,134,661,381,243,297,656,272,541,553,545,230,467,555,120,703,142,526,424,673,438,306,745,385,261,688,688,145,261,478,657,630,209,151,92,468,62,416,740,72,138,439,383,651,691,254,479,472,135,468,238,735,187,375,70,192,333,121,88,589,424,396,413,511,717,655,385,290,133,262,136,672,234,363,79,429,130,176,293,64,576,578,582,135,244,240,89,379,479,154,207,52,673,335,518,410,576,100,161,600,275,562,66,557,645,347,602,275,663,516,67,399,679,538,712,734,213,109,526,102,736,350,675,310,260,356,134,329,573,303,665,243,186,646,211,293,690,535,656,661,504,372,133,473,75,469,443,167,420,567,655,670,57,678,420,620,125,635,619,131,382,71,346,564,342,57,330,601,749,442,513,327,359,448,728,123,321,160,122,321,346,586,559,325,421,434,482,447,391,509,615,137,276,744,710,667,431,411,270,527,132,543,378,492,733,529,284,404,374,747,600,129,205,573,607,386,244,101,726,747,482,307,215,412,180,152,609,82,262,351,424,628,82,83,71,157,75,667,50,589,746,396,662,448,68,475,82,626,659,584,313,674,425,182,256,60,601,284,511,718,653,681,378,672,422,537,95,611,586,658,156,421,404,586,638,288,174,578,439,583,618,579,247,213,57,170,247,98,124,726,461,555,636,391,692,483,699,503,170,450,560,610,573,569,359,695,714,83,189,449,308,160,285,130,196,730,638,626,456,319,186,153,442,280,562,459,723,494,317,250,605,207,678,527,144,184,50,60,594,344,255,477,296,596,103,171,506,669,390,319,599,723,690,457,583,60,445,57,119,616,102,689,313,646,456,425,88,167,403,140,464,170,240,238,144,285,495,496,673,100,129,706,367,616,440,386,448,336,492,481,526,267,428,68,59,202,83,678,605,354,532,205,73,596,684,531,677,522,259,52,199,96,243,229,642,146,189,389,141,262,239,100,443,519,69,185,144,127,131,160,97,558,203,708,398,582,162,288,121,737,496,344,640,494,264,492,316,414,246,59,206,310,503,670,716,736,467,179,497,189,313,638,387,359,687,329,712,116,330,738,141,621,456,68,102,141,303,146,625,337,707,702,702,441,86,496,363,499,373,508,568,140,50,386,249,638,231,622,135,595,240,173,224,622,419,405,294,249,103,115,340,336,412,347,187,345,507,672,115,360,522,218,103,481,80,638,108,408,232,395,192,517,285,135,478,338,607,301,239,429,595,289,501,391,372,460,53,455,383,427,714,411,503,361,356,473,413,52,96,691,327,732,585,710,456,461,447,373,552,292,601,269,54,526,266,621,252,226,241,505,714,355,684,412,419,359,493,606,457,327,107,397,622,266,615,594,670,681,280,418,658,270,612,561,676,174,637,195,421,662,391,88,744,314,151,240,552,599,504,376,626,613,272,691,652,365,533,272,469,592,544,331,554,740,237,676,527,515,305,247,421,598,559,477,537,383,261,384,126,280,397,66,304,559,428,604,335,629,372,586,662,665,623,109,491,619,608,585,701,599,439,724,412,167,445,157,261,70,605,161,136,288,432,279,374,638,395,621,127,174,626,336,453,430,582,324,623,379,187,68,413,600,672,75,609,315,711,623,395,179,248,168,334,501,569,553,260,433,142,397,202,360,196,521,370,432,722,703,488,545,591,159,495,398,217,362,64,54,675,148,712,492,433,81,558,362,679,708,464,386,721,102,208,400,119,373,410,143,569,523,580,744,531,96,726,619,620,734,111,528,437,428,603,521,725,376,466,334,683,56,99,233,341,63,446,637,71,635,182,317,400,639,744,407,68,597,392,645,726,486,403,193,465,726,330,506,64,561,132,387,196,375,361,129,95,126,388,223,539,360,279,241,216,208,462,480,511,240,730,579,415,378,93,186,216,97,312,130,492,631,406,725,399,626,187,55,518,262,419,317,77,259,382,199,384,593,398,682,435,169,738,115,489,658,631,384,234,670,526,299,56,104,357,156,131,448,622,234,697,444,183,310,353,516,68,433,299,720,69,749,214,511,497,276,661,498,119,515,170,419,65,412,326,494,236,674,593,747,689,350,487,67,345,385,631,151,176,247,494,649,465,273,261,285,643,163,280,717,160,179,311,337,295,87,409,343,462,347,77,689,362,96,733,440,734,256,614,103,736,508,68,658,435,745,360,224,548,224,612,336,481,661,337,473,495,675,297,634,415,339,276,593,443,603,170,326,166,323,435,202,264,642,414,693,67,258,551,649,661,730,299,684,186,53,647,639,252,298,102,553,721,497,149,625,88,502,702,152,517,88,53,610,81,92,416,643,84,707,561,173,188,509,635,313,442,300,65,460,575,255,613,589,658,675,274,613,206,225,580,74,693,612,513,549,666,126,508,624,671,357,230,197,498,366,598,537,73,54,197,677,556,149,455,117,221,644,698,331,485,729,407,663,66,518,411,745,714,421,529,471,579,336,479,449,471,388,569,53,230,340,66,727,107,745,691,539,667,281,722,511,577,381,619,748,251,331,732,407,373,638,154,313,439,103,225,452,692,629,341,332,712,556,612,534,334,600,509,157,214,604,565,617,210,248,474,672,456,226,303,622,208,406,739,613,356,547,536,676,140,373,576,320,154,338,167,686,346,732,247,53,483,527,624,577,514,696,218,391,481,655,628,728,283,462,188,280,188,87,194,597,676,205,221,690,395,537,293,748,608,110,62,315,410,360,282,120,714,700,510,153,233,178,723,232,268,678,456,470,563,324,311,155,113,401,393,557,504,528,294,608,367,290,698,571,659,189,523,237,239,223,165,514,104,264,335,57,316,290,408,253,77,602,265,53,511,151,101,541,174,128,203,375,688,475,277,643,224,566,427,257,611,340,736,736,563,660,741,559,706,278,510,543,422,732,485,100,474,478,477,516,422,365,195,264,630,89,204,654,333,618,389,185,62,534,380,713,438,471,738,620,293,408,681,639,219,248,693,89,616,129,250,439,339,654,486,158,131,317,244,245,375,490,312,593,580,742,713,202,300,219,678,670,511,681,254,194,65,63,54,596,679,105,342,422,667,692,68,526,488,397,605,489,83,661,295,346,238,443,441,727,556,451,115,426,120,79,670,673,686,449,327,150,380,600,365,133,678,90,709,737,330,349,265,571,598,372,109,274,636,283,642,189,371,185,594,745,254,684,642,656,666,232,451,208,665,186,685,280,267,171,202,583,530,222,217,121,403,654,324,147,207,186,602,366,645,95,649,396,330,664,343,660,636,245,718,605,559,331,325,414,527,99,155,190,261,177,496,545,256,537,116,483,55,596,569,591,409,574,747,78,197,242,106,562,400,95,632,422,338,606,61,392,639,418,98,651,481,274,732,63,576,644,174,183,340,61,96,526,117,260,446,312,444,617,413,169,722,124,103,423,616,303,124,386,432,376,85,350,572,593,186,747,330,72,600,205,354,249,320,548,329,208,229,481,725,584,526,144,319,587,254,320,371,294,307,147,165,277,248,286,749,517,571,222,141,450,671,97,299,610,320,614,723,286,418,334,174,101,150,741,649,96,245,340,359,76,273,111,76,177,195,484,545,540,445,347,228,597,92,539,642,188,557,367,165,177,259,700,51,496,260,413,605,328,542,224,573,282,676,301,206,197,260,333,370,564,484,226,94,164,638,244,349,84,359,262,460,122,704,100,568,155,443,640,599,747,236,586,257,123,706,134,430,571,314,641,632,683,618,579,598,210,464,635,325,211,624,461,695,442,742,78,113,240,223,338,330,231,672,709,742,636,223,665,346,671,265,552,296,112,204,113,551,67,514,339,339,323,251,701,748,629,319,318,255,461,96,199,407,285,353,746,638,726,341,234,109,662,72,707,543,423,454,261,343,60,706,698,236,581,329,678,695,509,608,479,231,366,441,186,50,583,392,59,362,346,178,684,550,272,511,483,81,656,501,422,417,216,115,92,621,62,528,619,568,193,442,409,308,470,290,491,712,514,722,265,506,318,483,200,115,235,196,742,447,642,72,146,99,597,652,567,369,84,644,555,292,310,355,318,215,710,124,325,237,475,620,678,512,678,342,118,488,725,206,240,694,189,150,186,455,551,110,371,423,132,549,643,275,291,459,618,367,455,237,537,649,402,707,734,449,642,570,227,392,419,297,735,575,123,622,60,333,575,435,435,514,627,692,327,478,362,184,613,724,356,196,472,200,596,601,706,92,693,132,134,445,470,654,350,578,430,542,587,230,283,536,556,666,560,420,642,212,197,159,368,714,131,358,157,360,529,367,515,673,228,551,731,577,58,337,584,291,175,212,179,661,238,183,249,522,740,239,546,707,419,442,557,299,564,572,240,78,542,249,482,167,526,52,365,517,584,551,124,345,460,436,700,79,745,74,655,701,158,241,473,287,427,743,727,636,501,409,219,126,707,613,595,546,632,695,320,96,67,406,210,118,110,201,281,111,105,582,542,591,329,156,147,636,175,339,143,736,464,347,542,112,72,161,337,548,704,665,410,171,654,649,189,78,598,340,441,660,109,260,278,445,255,638,499,120,546,268,168,418,674,715,570,56,271,526,730,228,556,445,746,193,173,727,223,734,476,693,644,396,287,64,180,553,511,87,379,397,166,733,629,109,741,704,110,341,220,452,613,148,644,449,422,605,549,450,567,598,661,657,483,113,303,361,78,524,484,671,702,643,637,206,505,432,50,702,652,384,337,223,398,227,434,626,439,463,476,192,342,143,501,620,216,431,226,535,298,198,647,571,637,437,241,682,691,54,238,394,595,134,386,248,301,141,588,416,374,106,513,114,158,457,313,353,477,495,202,203,205,206,363,644,563,292,77,202,81,184,354,524,111,425,435,554,309,100,375,657,536,170,410,586,159,136,535,455,601,709,281,729,616,528,395,245,676,420,245,343,461,88,392,167,112,225,308,529,645,427,745,102,266,674,152,154,389,500,694,161,203,87,273,483,293,73,127,296,320,644,179,686,704,722,556,124,291,278,724,371,597,339,80,410,397,73,649,710,364,541,370}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468,351,324,439,670,493,533,626,508,123,234,640,182,705,614,205,453,580,631,93,98,687,590,271,384,110,264,574,171,522,173,555,468,325,570,489,733,725,352,174,63,138,431,123,599,80,248,474,516,546,429,395,248,250,595,639,539,531,352,533,178,348,104,285,106,345,698,660,203,598,135,618,254,479,705,571,626,334,92,608,274,729,711,605,487,547,709,242,599,673,730,304,74,235,482,234,365,152,251,253,269,600,173,59,445,277,638,462,619,576,79,56,350,521,675,92,495,221,256,430,63,365,509,399,427,590,298,488,581,279,219,272,649,149,385,573,212,520,460,700,744,380,575,302,515,288,454,55,154,368,735,206,678,689,494,446,110,189,609,573,701,130,452,295,335,209,264,493,624,247,189,299,645,490,425,341,221,551,156,56,352,363,422,490,305,338,690,627,508,92,415,141,324,664,159,514,418,96,497,92,433,128,283,572,230,709,265,87,102,191,254,178,457,351,100,356,398,255,552,657,556,221,598,648,364,56,624,429,284,619,350,328,125,203,312,76,68,229,615,557,715,169,92,92,468,435,729,658,363,179,548,327,50,601,662,366,466,583,710,490,306,565,166,521,571,421,565,699,603,290,260,76,507,628,307,223,415,709,88,427,648,174,468,611,581,640,554,132,313,536,129,419,263,153,742,296,167,73,344,631,606,174,532,704,196,100,268,82,694,559,115,342,218,540,548,687,260,732,726,724,245,608,475,360,680,416,139,388,617,726,698,379,267,87,586,356,358,370,616,227,288,348,527,122,359,71,172,689,176,92,723,704,696,422,82,83,496,422,452,719,555,203,652,268,232,155,586,736,324,188,358,472,349,592,538,172,522,253,551,345,335,432,595,612,331,347,452,700,379,664,649,73,586,56,512,134,108,129,106,304,554,81,167,416,254,560,379,326,554,102,371,103,726,402,551,564,677,567,445,622,391,196,431,581,674,683,151,460,424,480,602,601,610,320,530,219,268,496,360,637,191,682,688,706,448,293,631,658,376,550,387,378,294,535,401,225,391,384,192,677,244,635,291,417,424,301,442,727,170,634,684,728,712,703,284,336,364,291,80,512,312,247,693,662,479,535,612,588,319,284,641,717,263,589,355,244,736,533,252,589,312,616,723,639,577,169,397,321,175,667,381,716,687,701,346,369,289,641,536,709,715,394,83,358,728,687,193,640,711,239,321,480,524,630,602,228,360,563,628,576,320,108,208,336,453,501,350,398,742,162,413,59,237,713,332,268,535,601,224,95,604,92,600,447,638,305,659,222,355,166,315,490,223,630,306,125,371,316,147,269,253,229,368,681,338,394,81,360,715,660,261,190,332,323,320,62,466,363,270,417,121,429,626,462,697,725,63,304,109,84,211,689,670,223,225,285,196,695,113,304,645,192,247,578,687,730,466,283,720,336,95,86,80,221,727,74,192,356,113,305,463,144,628,493,599,651,717,703,261,356,327,400,527,686,308,361,408,162,693,650,201,377,553,590,707,341,539,641,584,314,517,163,600,654,692,472,417,340,581,584,56,509,623,248,540,220,261,75,696,90,423,335,420,351,541,281,657,654,238,176,646,309,117,218,155,263,676,201,70,544,506,55,317,80,546,709,533,629,300,232,581,453,593,132,549,653,216,566,693,302,733,569,518,493,526,123,98,593,675,341,127,505,457,147,650,423,105,278,118,190,696,697,202,716,405,97,280,616,143,541,680,215,132,153,308,123,413,459,476,438,230,538,201,678,511,749,288,716,571,293,214,550,259,705,518,93,707,100,134,661,381,243,297,656,272,541,553,545,230,467,555,120,703,142,526,424,673,438,306,745,385,261,688,688,145,261,478,657,630,209,151,92,468,62,416,740,72,138,439,383,651,691,254,479,472,135,468,238,735,187,375,70,192,333,121,88,589,424,396,413,511,717,655,385,290,133,262,136,672,234,363,79,429,130,176,293,64,576,578,582,135,244,240,89,379,479,154,207,52,673,335,518,410,576,100,161,600,275,562,66,557,645,347,602,275,663,516,67,399,679,538,712,734,213,109,526,102,736,350,675,310,260,356,134,329,573,303,665,243,186,646,211,293,690,535,656,661,504,372,133,473,75,469,443,167,420,567,655,670,57,678,420,620,125,635,619,131,382,71,346,564,342,57,330,601,749,442,513,327,359,448,728,123,321,160,122,321,346,586,559,325,421,434,482,447,391,509,615,137,276,744,710,667,431,411,270,527,132,543,378,492,733,529,284,404,374,747,600,129,205,573,607,386,244,101,726,747,482,307,215,412,180,152,609,82,262,351,424,628,82,83,71,157,75,667,50,589,746,396,662,448,68,475,82,626,659,584,313,674,425,182,256,60,601,284,511,718,653,681,378,672,422,537,95,611,586,658,156,421,404,586,638,288,174,578,439,583,618,579,247,213,57,170,247,98,124,726,461,555,636,391,692,483,699,503,170,450,560,610,573,569,359,695,714,83,189,449,308,160,285,130,196,730,638,626,456,319,186,153,442,280,562,459,723,494,317,250,605,207,678,527,144,184,50,60,594,344,255,477,296,596,103,171,506,669,390,319,599,723,690,457,583,60,445,57,119,616,102,689,313,646,456,425,88,167,403,140,464,170,240,238,144,285,495,496,673,100,129,706,367,616,440,386,448,336,492,481,526,267,428,68,59,202,83,678,605,354,532,205,73,596,684,531,677,522,259,52,199,96,243,229,642,146,189,389,141,262,239,100,443,519,69,185,144,127,131,160,97,558,203,708,398,582,162,288,121,737,496,344,640,494,264,492,316,414,246,59,206,310,503,670,716,736,467,179,497,189,313,638,387,359,687,329,712,116,330,738,141,621,456,68,102,141,303,146,625,337,707,702,702,441,86,496,363,499,373,508,568,140,50,386,249,638,231,622,135,595,240,173,224,622,419,405,294,249,103,115,340,336,412,347,187,345,507,672,115,360,522,218,103,481,80,638,108,408,232,395,192,517,285,135,478,338,607,301,239,429,595,289,501,391,372,460,53,455,383,427,714,411,503,361,356,473,413,52,96,691,327,732,585,710,456,461,447,373,552,292,601,269,54,526,266,621,252,226,241,505,714,355,684,412,419,359,493,606,457,327,107,397,622,266,615,594,670,681,280,418,658,270,612,561,676,174,637,195,421,662,391,88,744,314,151,240,552,599,504,376,626,613,272,691,652,365,533,272,469,592,544,331,554,740,237,676,527,515,305,247,421,598,559,477,537,383,261,384,126,280,397,66,304,559,428,604,335,629,372,586,662,665,623,109,491,619,608,585,701,599,439,724,412,167,445,157,261,70,605,161,136,288,432,279,374,638,395,621,127,174,626,336,453,430,582,324,623,379,187,68,413,600,672,75,609,315,711,623,395,179,248,168,334,501,569,553,260,433,142,397,202,360,196,521,370,432,722,703,488,545,591,159,495,398,217,362,64,54,675,148,712,492,433,81,558,362,679,708,464,386,721,102,208,400,119,373,410,143,569,523,580,744,531,96,726,619,620,734,111,528,437,428,603,521,725,376,466,334,683,56,99,233,341,63,446,637,71,635,182,317,400,639,744,407,68,597,392,645,726,486,403,193,465,726,330,506,64,561,132,387,196,375,361,129,95,126,388,223,539,360,279,241,216,208,462,480,511,240,730,579,415,378,93,186,216,97,312,130,492,631,406,725,399,626,187,55,518,262,419,317,77,259,382,199,384,593,398,682,435,169,738,115,489,658,631,384,234,670,526,299,56,104,357,156,131,448,622,234,697,444,183,310,353,516,68,433,299,720,69,749,214,511,497,276,661,498,119,515,170,419,65,412,326,494,236,674,593,747,689,350,487,67,345,385,631,151,176,247,494,649,465,273,261,285,643,163,280,717,160,179,311,337,295,87,409,343,462,347,77,689,362,96,733,440,734,256,614,103,736,508,68,658,435,745,360,224,548,224,612,336,481,661,337,473,495,675,297,634,415,339,276,593,443,603,170,326,166,323,435,202,264,642,414,693,67,258,551,649,661,730,299,684,186,53,647,639,252,298,102,553,721,497,149,625,88,502,702,152,517,88,53,610,81,92,416,643,84,707,561,173,188,509,635,313,442,300,65,460,575,255,613,589,658,675,274,613,206,225,580,74,693,612,513,549,666,126,508,624,671,357,230,197,498,366,598,537,73,54,197,677,556,149,455,117,221,644,698,331,485,729,407,663,66,518,411,745,714,421,529,471,579,336,479,449,471,388,569,53,230,340,66,727,107,745,691,539,667,281,722,511,577,381,619,748,251,331,732,407,373,638,154,313,439,103,225,452,692,629,341,332,712,556,612,534,334,600,509,157,214,604,565,617,210,248,474,672,456,226,303,622,208,406,739,613,356,547,536,676,140,373,576,320,154,338,167,686,346,732,247,53,483,527,624,577,514,696,218,391,481,655,628,728,283,462,188,280,188,87,194,597,676,205,221,690,395,537,293,748,608,110,62,315,410,360,282,120,714,700,510,153,233,178,723,232,268,678,456,470,563,324,311,155,113,401,393,557,504,528,294,608,367,290,698,571,659,189,523,237,239,223,165,514,104,264,335,57,316,290,408,253,77,602,265,53,511,151,101,541,174,128,203,375,688,475,277,643,224,566,427,257,611,340,736,736,563,660,741,559,706,278,510,543,422,732,485,100,474,478,477,516,422,365,195,264,630,89,204,654,333,618,389,185,62,534,380,713,438,471,738,620,293,408,681,639,219,248,693,89,616,129,250,439,339,654,486,158,131,317,244,245,375,490,312,593,580,742,713,202,300,219,678,670,511,681,254,194,65,63,54,596,679,105,342,422,667,692,68,526,488,397,605,489,83,661,295,346,238,443,441,727,556,451,115,426,120,79,670,673,686,449,327,150,380,600,365,133,678,90,709,737,330,349,265,571,598,372,109,274,636,283,642,189,371,185,594,745,254,684,642,656,666,232,451,208,665,186,685,280,267,171,202,583,530,222,217,121,403,654,324,147,207,186,602,366,645,95,649,396,330,664,343,660,636,245,718,605,559,331,325,414,527,99,155,190,261,177,496,545,256,537,116,483,55,596,569,591,409,574,747,78,197,242,106,562,400,95,632,422,338,606,61,392,639,418,98,651,481,274,732,63,576,644,174,183,340,61,96,526,117,260,446,312,444,617,413,169,722,124,103,423,616,303,124,386,432,376,85,350,572,593,186,747,330,72,600,205,354,249,320,548,329,208,229,481,725,584,526,144,319,587,254,320,371,294,307,147,165,277,248,286,749,517,571,222,141,450,671,97,299,610,320,614,723,286,418,334,174,101,150,741,649,96,245,340,359,76,273,111,76,177,195,484,545,540,445,347,228,597,92,539,642,188,557,367,165,177,259,700,51,496,260,413,605,328,542,224,573,282,676,301,206,197,260,333,370,564,484,226,94,164,638,244,349,84,359,262,460,122,704,100,568,155,443,640,599,747,236,586,257,123,706,134,430,571,314,641,632,683,618,579,598,210,464,635,325,211,624,461,695,442,742,78,113,240,223,338,330,231,672,709,742,636,223,665,346,671,265,552,296,112,204,113,551,67,514,339,339,323,251,701,748,629,319,318,255,461,96,199,407,285,353,746,638,726,341,234,109,662,72,707,543,423,454,261,343,60,706,698,236,581,329,678,695,509,608,479,231,366,441,186,50,583,392,59,362,346,178,684,550,272,511,483,81,656,501,422,417,216,115,92,621,62,528,619,568,193,442,409,308,470,290,491,712,514,722,265,506,318,483,200,115,235,196,742,447,642,72,146,99,597,652,567,369,84,644,555,292,310,355,318,215,710,124,325,237,475,620,678,512,678,342,118,488,725,206,240,694,189,150,186,455,551,110,371,423,132,549,643,275,291,459,618,367,455,237,537,649,402,707,734,449,642,570,227,392,419,297,735,575,123,622,60,333,575,435,435,514,627,692,327,478,362,184,613,724,356,196,472,200,596,601,706,92,693,132,134,445,470,654,350,578,430,542,587,230,283,536,556,666,560,420,642,212,197,159,368,714,131,358,157,360,529,367,515,673,228,551,731,577,58,337,584,291,175,212,179,661,238,183,249,522,740,239,546,707,419,442,557,299,564,572,240,78,542,249,482,167,526,52,365,517,584,551,124,345,460,436,700,79,745,74,655,701,158,241,473,287,427,743,727,636,501,409,219,126,707,613,595,546,632,695,320,96,67,406,210,118,110,201,281,111,105,582,542,591,329,156,147,636,175,339,143,736,464,347,542,112,72,161,337,548,704,665,410,171,654,649,189,78,598,340,441,660,109,260,278,445,255,638,499,120,546,268,168,418,674,715,570,56,271,526,730,228,556,445,746,193,173,727,223,734,476,693,644,396,287,64,180,553,511,87,379,397,166,733,629,109,741,704,110,341,220,452,613,148,644,449,422,605,549,450,567,598,661,657,483,113,303,361,78,524,484,671,702,643,637,206,505,432,50,702,652,384,337,223,398,227,434,626,439,463,476,192,342,143,501,620,216,431,226,535,298,198,647,571,637,437,241,682,691,54,238,394,595,134,386,248,301,141,588,416,374,106,513,114,158,457,313,353,477,495,202,203,205,206,363,644,563,292,77,202,81,184,354,524,111,425,435,554,309,100,375,657,536,170,410,586,159,136,535,455,601,709,281,729,616,528,395,245,676,420,245,343,461,88,392,167,112,225,308,529,645,427,745,102,266,674,152,154,389,500,694,161,203,87,273,483,293,73,127,296,320,644,179,686,704,722,556,124,291,278,724,371,597,339,80,410,397,73,649,710,364,541,370,129,453,334,167,202,365,630,155,144,406,392,51,269,377,467,511,153,50,279,298,641,242,615,700,565,537,669,352,66,373,229,57,445,196,154,678,110,589,650,575,567,62,102,163,683,712,531,403,186,620,622,582,220,405,316,70,323,235,420,335,349,667,230,573,737,668,587,520,712,116,303,429,660,443,78,462,441,384,634,80,78,503,180,700,435,653,309,144,695,146,292,420,275,679,303,513,404,297,273,450,168,91,444,580,573,394,334,616,278,83,335,289,97,657,317,455,464,673,526,652,419,471,224,556,371,237,412,184,574,521,118,449,379,247,411,706,456,597,583,156,135,451,214,536,243,735,133,404,441,444,718,631,97,88,274,736,385,532,736,229,335,444,594,188,204,568,336,736,546,658,635,411,234,280,271,683,588,693,354,518,283,666,133,651,431,383,539,600,267,146,220,436,643,230,613,485,208,228,544,600,269,334,230,731,279,355,716,475,323,355,254,667,418,580,661,348,573,739,345,667,383,320,686,660,554,450,637,611,146,445,674,370,475,682,578,261,716,280,148,99,569,464,749,638,229,526,389,155,211,692,741,701,479,666,384,515,394,456,525,217,530,287,560,710,673,572,696,335,268,675,445,301,629,440,385,104,633,277,544,583,289,124,567,679,737,511,502,629,92,216,239,564,649,566,233,105,496,523,359,513,312,468,187,204,399,625,140,353,563,380,366,653,343,335,749,55,389,278,574,217,353,368,616,98,639,171,513,507,646,712,64,413,312,387,418,742,743,266,625,451,518,508,150,401,700,707,612,586,623,510,597,85,464,139,710,164,644,276,232,568,242,733,271,733,665,414,393,600,457,466,320,385,189,119,197,511,241,694,451,233,55,458,632,625,680,160,168,679,135,159,179,414,554,641,420,595,675,740,313,529,157,580,194,746,430,709,54,258,155,599,720,712,329,119,329,310,431,589,272,148,686,396,357,55,261,597,205,410,616,444,523,690,290,79,614,448,51,454,462,735,560,553,458,488,380,82,191,330,416,651,247,665,546,57,552,55,603,324,270,147,355,650,298,324,339,277,669,239,133,386,64,572,289,683,60,373,476,525,107,587,94,156,531,192,367,718,478,483,388,413,230,430,61,540,508,483,657,638,239,288,163,237,284,340,146,531,446,95,522,575,543,604,526,315,69,203,246,538,139,284,668,605,719,193,732,66,423,311,577,221,294,481,477,732,65,619,649,245,385,684,596,506,613,150,153,741,124,516,671,563,623,739,509,111,320,312,543,76,305,619,691,420,474,170,459,478,561,438,473,573,77,615,305,215,565,451,358,437,301,574,597,586,457,64,223,238,617,137,710,514,237,100,611,734,494,389,646,229,354,325,306,670,612,725,728,256,708,650,597,480,384,230,67,398,66,245,209,342,160,132,214,340,424,153,238,693,539,54,356,662,227,748,355,314,488,440,356,145,421,241,746,561,546,532,225,90,645,715,524,446,259,384,348,537,263,69,648,705,693,339,334,402,332,604,51,224,227,172,512,308,715,453,448,694,503,156,112,135,541,116,652,336,312,564,347,214,301,357,215,713,692,110,51,245,491,455,152,655,267,431,340,321,360,177,389,514,631,185,489,446,437,438,532,145,420,138,522,666,321,391,390,135,509,662,589,322,112,221,628,230,131,743,486,170,237,274,247,560,676,451,436,647,495,623,420,274,241,298,226,93,312,300,572,177,662,375,739,336,600,678,673,612,552,566,212,199,719,246,177,502,228,393,107,602,223,232,249,732,344,630,325,743,508,534,547,130,235,302,271,329,439,673,69,609,455,388,442,354,452,316,491,592,109,452,111,423,735,436,204,129,227,67,290,563,232,338,608,509,658,644,347,583,161,598,728,105,728,105,582,480,138,109,472,731,434,571,566,383,250,439,371,453,236,195,146,403,264,384,317,357,477,633,502,237,691,70,707,366,189,346,54,496,431,733,507,539,201,274,627,609,189,277,148,253,408,283,102,575,507,740,340,232,534,267,457,713,232,305,471,726,424,148,745,118,220,466,319,579,507,438,570,539,169,126,673,117,229,51,310,645,300,223,117,289,317,360,162,693,145,238,453,675,452,50,65,173,157,567,539,666,54,583,613,98,312,229,567,66,740,641,509,333,300,161,670,324,548,271,615,213,638,643,226,515,325,337,324,347,401,715,598,176,127,738,337,385,63,135,591,525,638,183,95,683,171,67,745,111,302,71,255,89,360,80,534,233,481,577,699,556,693,404,653,509,65,746,223,171,139,616,113,298,75,209,733,202,650,496,105,392,556,508,411,592,155,53,473,542,237,521,367,657,353,145,292,247,427,599,173,347,494,264,525,129,273,570,413,685,139,412,427,194,353,229,250,671,739,186,736,109,622,675,663,618,476,327,588,744,198,184,461,74,486,445,400,290,158,61,440,142,629,742,696,563,233,345,89,278,245,426,369,215,604,105,542,347,552,429,444,541,103,428,613,76,468,500,414,361,485,192,595,341,214,198,149,662,270,547,210,252,674,347,372,67,499,708,721,91,258,740,102,489,663,125,421,592,624,61,625,128,662,323,362,580,58,410,356,51,423,80,409,95,478,88,472,635,188,254,739,665,127,62,731,566,151,608,233,483,703,691,366,193,289,218,326,300,140,313,137,196,700,368,318,534,227,593,101,701,682,485,216,125,499,608,62,520,265,317,348,570,160,471,578,320,63,564,515,404,223,171,703,283,306,103,214,217,373,569,149,581,194,170,299,620,233,147,743,304,656,717,72,632,560,98,236,336,749,641,628,698,450,742,110,388,180,586,650,666,536,521,551,330,367,245,431,619,327,320,579,542,594,629,263,302,55,287,476,743,720,231,517,213,727,592,512,596,250,350,399,506,177,141,243,442,420,690,467,151,605,348,440,558,498,619,746,354,111,210,491,736,356,687,138,71,417,519,260,386,307,580,246,614,204,748,169,743,663,563,142,723,265,479,329,309,481,93,78,90,512,521,434,558,350,546,744,386,412,116,412,347,108,737,457,524,246,365,625,634,283,595,181,79,367,127,704,623,264,620,394,687,152,615,651,713,463,232,148,256,50,527,556,509,318,92,166,464,704,392,514,216,177,100,341,565,690,494,650,327,458,479,551,536,434,326,481,119,648,74,154,631,576,207,273,439,193,185,405,263,667,643,73,331,466,354,628,217,147,652,252,415,606,485,453,394,482,633,320,101,653,580,220,222,405,685,521,92,461,448,593,345,484,134,199,497,249,121,649,415,360,563,447,521,169,668,702,519,223,548,297,288,147,78,290,431,612,66,368,490,403,718,657,441,590,118,88,169,68,220,311,196,539,133,660,155,660,183,662,354,596,624,72,659,114,246,696,57,610,676,482,527,674,364,680,254,339,282,75,63,586,744,709,250,567,171,571,720,228,700,376,70,194,612,740,731,483,628,293,681,92,546,559,147,304,707,725,358,643,610,213,178,744,188,354,76,182,517,508,516,582,265,550,661,685,146,310,234,486,332,425,170,262,525,74,391,232,516,235,635,248,287,127,66,476,83,284,490,189,265,97,338,375,746,216,412,152,746,524,58,329,708,649,263,741,260,540,581,482,686,133,540,120,296,705,165,386,299,352,140,741,316,728,427,496,282,310,546,390,461,65,103,93,419,747,234,525,422,621,511,534,78,615,652,690,575,560,139,485,233,140,452,352,445,176,234,510,442,541,611,130,480,623,190,189,540,187,583,380,367,90,254,659,111,553,541,672,573,537,354,745,713,95,581,208,414,354,421,573,605,209,293,110,318,335,182,495,308,599,660,361,93,618,436,710,98,713,649,71,254,250,175,157,446,330,174,602,252,629,575,622,477,474,311,247,148,445,635,91,246,565,265,563,459,628,119,471,329,170,79,117,697,718,198,247,270,272,651,573,664,505,281,565,227,712,531,453,618,348,700,595,381,50,471,556,146,134,486,186,320,380,391,545,319,510,363,164,104,314,73,70,744,493,583,538,158,249,86,348,656,426,466,231,162,295,100,514,655,317,674,339,497,485,706,65,245,563,673,618,102,156,742,691,104,317,75,415,603,679,136,231,471,180,226,693,661,542,680,198,136,202,64,722,746,584,326,632,236,678,454,355,736,203,608,267,255,80,681,358,254,550,746,684,586,573,126,109,354,732,136,482,250,56,228,326,276,662,155,336,654,478,293,405,710,578,378,686,723,205,537,612,696,391,290,77,547,441,128,347,153,463,330,682,504,314,228,179,145,739,438,90,735,306,588,195,663,521,55,532,216,382,522,533,683,382,676,112,432,714,319,76,574,407,507,374,562,82,586,65,707,459,696,61,250,626,536,155,281,665,403,312,659,79,521,449,140,151,315,313,357,382,290,150,191,350,457,648,571,163,413,561,280,379,557,61,686,202,721,179,79,337,216,737,316,119,133,459,355,735,268,667,125,218,116,509,255,652,483,103,560,712,198,70,739,344,581,105,164,467,724,685,382,214,657,161,616,737,488,548,272,116,707,507,548,66,616,374,351,152,251,108,96,354,101,416,373,592,630,576,476,352,216,457,473,549,300,526,211,512,143,221,741,116,113,426,115,523,413,346,486,429,670,161,694,239,357,282,381,605,683,170,497,555,77,538,282,679,58,65,211,578,437,593,443,267,110,657,464,484,113,154,311,635,718,211,310,696,670,195,474,294,672,356,269,743,199,549,406,108,484,63,261,242,113,350,634,502,431,200,655,146,194,103,453,196,63,687,152,185,157,447,289,611,292,430,161,599,562,121,69,707,549,258,456,220,386,179,653,448,259,717,326,87,211,172,648,93,72,175,647,120,460,494,740,457,284,302,240,128,588,266,104,180,595,249,592,422,146,595,680,208,126,535,502,220,194,171,570,83,185,531,556,709,565,441,718,568,732,214,190,601,440,353,462,466,546,175,421,303,559,582,383,522,561,144,535,259,357,495,54,647,656,65,255,642,663,610,684,630,373,561,420,624,718,539,605,281,493,175,78,401,609,645,122,229,172,421,76,191,246,391,329,263,113,740,61,77,55,313,324,257,676,94,630,215,653,452,600,110,407,178,230,197,472,743,714,52,424,311,679,304,607,666,123,368,544,275,201,272,99,701,351,215,191,531,447,304,213,732,419,334,311,615,413,620,330,656,728,69,653,199,162,61,97,474,306,384,161,495,394,50,625,191,337,93,290,317,627,392,630,261,626,611,433,173,741,519,112,84,158,94,102,234,717,535,597,274,505,650,645,467,95,280,451,169,166,624,486,58,625,613,133,333,376,212,458,85,119,694,664,300,158,255,149,540,106,372,589,682,550,279,66,581,300,697,530,535,681,428,297,153,113,420,413,249,367,477,216,210,430,686,673,265,647,325,361,350,259,577,233,675,519,272,471,607,678,380,509,286,66,61,464,455,683,249,107,416,430,93,735,127,442,125,638,465,575,127,91,158,432,569,223,499,211,734,480,682,402,59,395,295,582,256,90,313,615,207,677,181,324,581,453,391,714,487,317,222,242,166,654,569,190,419,661,197,335,549,350,746,181,600,580,237,647,640,108,617,572,231,67,317,476,383,583,630,221,742,657,62,252,376,658,323,404,499,374,274,253,185,365,399,440,408,531,156,86,663,422,158,129,67,337,160,314,335,197,431,373,534,712,688,645,403,472,511,210,199,393,629,109,510,434,371,363,66,364,582,452,112,64,625,340,465,410,290,582,628,630,453,456,673,124,189,254,616,478,498,163,687,405,355,332,278,384,285,181,228,432,561,438,687,80,590,352,134,190,279,62,624,523,427,564,629,459,230,176,555,163,366,485,694,152,288,316,359,116,396,287,342,692,101,379,690,150,369,651,144,622,382,438,413,62,511,607,686,230,126,387,123,381,436,75,644,237,434,588,577,423,657,542,170,513,439,142,739,82,728,295,273,231,465,633,100,207,228,495,393,168,674,602,196,105,515,78,628,716,301,80,440,740,160,551,189,370,304,272,732,597,261,156,286,518,465,56,527,439,84,587,517,375,202,655,127,668,331,341,746,354,368,301,551,254,447,225,515,271,514,162,649,585,462,304,734,368,731,710,645,200,692,327,465,512,432,747,415,559,147,562,60,412,380,327,490,670,556,527,121,185,557,235,484,699,218,418,702,490,133,154,287,58,208,368,670,181,635,636,157,608,536,107,186,505,118,361,498,660,633,153,83,638,648,469,487,629,148,401,243,168,428,556,250,537,208,72,537,253,613,467,163,448,174,480,273,189,310,345,450,117,704,391,340,225,132,659,80,719,64,362,52,502,106,498,605,635,258,76,324,50,698,506,416,207,701,286,647,363,566,188,489,644,287,346,370,534,210,70,578,579,301,524,283,477,412,738,228,404,412,636,740,124,656,193,712,358,259,183,700,627,470,644,277,615,535,606,570,192,186,590,103,51,743,531,541,396,154,349,125,640,296,543,162,335,656,256,53,404,231,623,660,61,147,700,233,90,213,186,657,637,683,291,528,605,230,691,239,626,58,171,301,595,238,528,122,77,491,459,746,341,346,739,343,150,640,99,744,545,552,380,370,534,607,717,142,463,134,309,659,140,745,163,463,563,732,81,613,224,479,501,325,217,717,747,155,511,425,366,652,245,637,467,441,636,491,77,391,135,696,194,424,482,132,683,162,739,323,392,140,65,531,719,262,440,642,85,143,118,168,261,179,629,451,91,259,251,595,262,447,670,74,670,225,259,730,259,645,103,503,644,450,489,234,284,547,237,208,198,601,79,412,355,727,296,686,539,128,206,198,98,502,714,693,736,164,553,185,161,498,728,352,709,193,499,616,177,510,748,464,668,727,693,481,565,287,335,189,218,569,82,237,120,660,579,466,484,559,405,616,294,749,92,233,594,659,664,569,401,564,599,159,434,273,352,240,606,250,53,219,108,59,512,648,740,716,436,182,183,303,214,90,351,174,521,115,431,254,382,736,537,678,151,234,203,148,508,132,547,354,578,94,71,569,147,313,570,565,272,437,276,471,263,344,55,92,659,390,471,721,661,240,539,313,716,363,345,174,420,607,614,185,599,405,465,157,694,189,719,270,629,694,647,436,201,593,725,732,89,497,137,412,96,142,702,172,428,718,168,498,696,361,287,506,536,726,416,362,85,617,429,489,738,95,582,180,724,162,514,724,574,512,584,90,551,113,564,283,458,198,79,563,349,558,355,381,261,253,324,724,117,705,571,601,97,271,111,71,487,351,89,378,566,254,573,99,518,425,565,561,78,313,590,533,109,472,112,356,70,531,444,307,530,158,673,203,240,137,165,419,558,702,370,571,509,374,247,361,183,255,741,469,364,699,653,636,477,442,421,250,55,604,642,264,157,711,217,284,211,343,444,666,469,238,190,197,214,286,453,317,433,112,427,717,235,305,141,492,635,256,687,365,196,204,118,230,657,63,515,713,375,302,600,438,708,105,529,168,382,64,57,132,404,232,603,81,210,391,476,146,189,134,508,215,707,721,305,182,618,150,455,350,184,420,264,327,357,608,629,66,409,372,139,392,96,550,207,599,503,493,319,638,371,677,690,262,57,76,301,472,304,254,625,568,77,655,427,321,679,315,485,334,664,600,733,426,433,717,578,420,142,154,661,416,183,558,196,274,199,300,681,646,103,570,271,509,480,328,366,98,383,320,463,652,394,626,680,466,257,603,555,329,222,166,591,260,252,491,485,619,422,538,269,73,573,116,122,531,424,203,252,747,160,136,539,641,321,575,169,642,421,381,580,189,698,172,746,460,52,619,390,379,509,582,485,727,226,544,318,566,277,392,689,577,120,648,634,397,628,55,365,61,82,324,709,225,82,671,181,123,169,437,687,281,694,155,314,613,733,211,204,572,194,172,502,467,542,383,723,110,544,700,55,623,748,239,155,264,500,382,431,265,309,59,149,318,373,140,695,183,219,303,348,454,554,623,60,195,631,742,307,387,384,239,620,116,146,609,439,536,252,208,449,357,494,237,140,537,84,631,677,492,59,612,605,665,663,379,680,573,727,569,369,676,91,690,446,133,712,79,275,358,57,562,396,668,557,371,234,408,569,423,274,410,677,616,81,524,258,130,718,138,279,551,729,562,668,601,574,196,518,153,99,643,508,333,272,445,689,114,586,327,702,350,602,296,693,145,455,172,694,642,608,274,226,482,426,121,565,716,220,202,338,432,166,624,434,147,632,522,204,646,705,223,626,284,572,218,614,656,627,408,137,728,664,468,434,428,709,391,458,196,361,748,618,454,537,121,70,164,704,120,440,239,521,619,121,605,87,443,516,248,570,242,461,355,571,208,452,706,67,408,137,104,742,54,212,512,745,224,168,577,454,162,282,110,593,367,685,342,395,418,205,393,348,255,577,377,746,259,206,210,114,451,720,231,705,396,369,168,259,285,331,725,669,512,435,137,84,483,53,495,314,462,201,203,696,740,170,131,357,485,118,144,572,726,180,233,517,360,157,685,334,734,541,312,621,600,231,80,139,662,86,211,465,243,491,579,130,295,310,333,665,557,530,73,734,673,81,207,181,649,247,733,536,209,603,735,589,359,704,388,374,519,590,111,213,140,72,668,613,486,57,313,381,187,111,533,566,109,742,60,681,597,168,726,575,634,228,656,315,191,448,242,175,541,609,485,247,227,242,635,233,618,446,173,231,356,73,330,192,694,185,218,735,559,68,669,488,396,625,161,716,598,530,650,203,735,569,561,116,134,158,686,539,261,389,202,401,143,634,129,261,203,563,558,618,437,358,523,112,640,427,86,537,494,633,315,250,725,427,85,348,279,215,104,547,491,176,284,536,743,729,684,64,337,237,86,367,378,400,468,139,87,201,716,125,255,553,193,593,666,363,165,328,381,148,324,55,672,491,437,441,463,679,319,255,332,393,211,399,134,396,674,249,324,137,183,395,193,136,393,219,227,435,540,678,449,410,88,330,224,282,697,300,383,579,169,547,507,370,121,719,78,385,186,216,440,360,641,195,179,219,353,435,264,154,103,279,59,369,591,79,689,563,99,225,302,83,300,739,580,512,136,211,674,316,706,482,169,308,711,342,711,227,706,166,244,514,400,257,340,579,230,619,708,339,195,462,732,687,720,123,411,487,357,256,136,575,474,719,231,418,306,195,668,427,322,555,646,658,669,336,509,443,515,701,216,551,695,715,660,410,589,188,404,676,206,122,398,566,183,144,97,559,189,539,501,727,748,216,662,214,306,623,664,59,474,549,720,56,132,526,108,524,421,302,670,643,67,511,398,568,499,302,135,402,279,540,491,518,172,599,668,236,130,352,204,113,366,374,445,422,635,459,50,742,88,439,690,377,591,80,364,146,547,625,741,571,88,662,145,331,177,68,447,84,643,300,79,300,366,438,178,422,358,62,515,313,252,393,292,666,124,321,311,309,593,412,555,304,675,355,496,516,366,494,736,162,462,446,295,245,432,394,348,167,170,104,354,56,532,705,511,670,187,312,249,637,350,217,665,255,165,285,532,528,667,749,360,633,398,610,342,229,101,97,346,369,695,407,572,78,664,617,186,346,57,693,267,131,526,114,166,473,316,622,455,280,534,342,424,194,345,289,736,136,414,544,456,419,274,360,579,557,331,694,296,410,466,102,147,616,157,114,252,104,670,170,616,492,706,321,263,435,549,681,678,324,137,350,346,474,130,321,745,668,245,364,406,298,659,556,571,103,726,557,737,687,743,308,146,244,530,490,211,155,284,726,563,632,173,131,699,558,414,667,504,128,492,724,677,127,613,468,656,80,653,51,283,191,470,457,552,199,445,474,344,150,241,619,425,552,513,174,264,179,166,625,59,747,549,177,347,146,682,69,183,52,496,610,493,434,79,683,707,290,66,334,491,283,511,261,405,227,648,197,562,570,318,86,473,666,655,248,83,278,519,182,549,339,599,677,379,583,262,245,414,389,518,557,279,51,540,682,254,283,709,436,633,457,203,378,724,80,299,236,61,458,168,491,570,55,278,250,617,85,688,499,460,707,590,195,253,180,315,453,665,348,644,613,108,130,638,433,378,76,111,462,82,205,147,271,465,214,136,230,269,460,286,170,585,406,362,670,676,474,188,424,572,337,426,561,556,118,606,411,321,64,191,121,562,321,502,86,371,115,444,144,627,383,156,592,290,137,431,640,294,374,645,287,159,473,629,420,393,552,260,310,601,541,447,273,530,612,512,674,92,454,226,604,696,422,646,479,315,175,553,400,308,411,140,98,352,421,73,79,707,639,647,384,99,338,150,289,665,425,85,85,328,55,676,563,438,479,120,273,559,541,289,422,202,197,227,358,217,62,535,353,89,684,238,529,137,266,742,440,263,283,277,465,514,700,99,716,441,590,703,696,289,108,666,158,129,731,509,75,88,649,109,593,717,485,366,175,643,103,341,275,388,430,464,317,186,534,437,84,433,377,623,280,254,716,316,744,642,150,505,273,50,364,621,607,199,638,147,213,73,130,364,579,223,98,641,749,238,281,348,130,566,270,168,236,616,745,623,551,156,702,717,279,575,315,443,217,631,477,548,324,606,405,680,173,636,588,276,341,650,388,377,208,246,225,246,106,355,723,247,212,553,469,332,234,653,294,555,248,409,700,561,746,426,284,260,662,241,422,641,663,271,415,601,182,266,337,73,531,202,500,78,182,648,315,507,453,486,233,690,531,602,155,321,287,163,75,247,303,297,518,505,503,520,521,212,418,383,747,731,154,692,449,83,741,156,456,723,531,488,687,504,573,517,387,547,73,307,644,481,673,708,368,645,358,698,702,288,434,127,470,157,158,229,217,438,433,386,394,460,209,322,626,98,525,173,268,477,402,638,98,282,466,320,597,125,339,685,653,741,689,436,275,326,430,314,717,337,612,173,314,285,530,265,531,102,273,182,536,289,94,454,474,408,578,719,210,140,596,592,146,127,651,76,356,556,643,211,74,711,568,585,559,523,263,552,179,494,505,537,117,353,685,656,571,53,212,670,697,235,607,94,347,166,311,99,637,605,385,70,599,375,555,113,188,362,393,592,425,135,185,508,513,649,71,351,602,483,672,444,703,657,327,478,727,324,61,544,345,349,396,408,553,723,526,685,610,75,455,632,362,362,298,90,382,340,623,466,510,172,551,81,172,420,75,309,631,316,677,511,352,336,351,626,453,281,274,426,127,609,652,143,565,664,743,490,515,252,256,104,307,150,545,460,485,641,465,237,551,155,298,632,591,85,671,286,442,569,701,348,154,86,296,293,170,563,566,155,534,526,496,448,259,699,111,562,418,447,306,743,488,313,472,516,266,502,685}

	//fmt.Println(len(latencySLO_500s))
	//fmt.Println(len(latencySLO_1000s))
	//fmt.Println(len(latencySLO_2000s))
	//fmt.Println(len(latencySLO_5000s))
	fmt.Println("scale=", len(latencySLO_100s))


	clusterCapConfig := InitBox(5000)
	cpuOverSell:= int32(0) //threads
	gpuOverSell:= int32(0) //SM percentage

	start := time.Now()

	cpuTotalConsum := int32(0)
	gpuTotalConsum := int32(0)
	memoryTotalConsum := float64(0)

	for l:=0; l<len(latencySLO_100s); l++ {

		resourcesConfigs :=  testEstimator(float64(latencySLO_100s[l]))
		if len(resourcesConfigs) == 0 {
			continue
		}

		cpuConsumedRate := float64(0)
		gpuMemConsumedRate := float64(0)
		gpuCoreConsumedRate := float64(0)

		maxResourceQuotaNagDiffIndex := -1
		minResourceQuotaPosDiffIndex := -1
		pickConfigIndex := -1
		maxResourceQuotaNagDiff := float64(-999)
		minResourceQuotaPosDiff := float64(999)
		tempGpuCoreQuota := float64(0)
		tempCpuQuota := float64(0)
		tempDiffQuota := float64(0)
		for i := 0; i < len(clusterCapConfig); i++ { // per node
			/** CPU GPU consumed rate **/

			cpuConsumedRate = 1.0 - float64(clusterCapConfig[i].CpuThreadsCap) / float64(20 + cpuOverSell) // cpu usage rate in node i socket j
			gpuMemConsumedRate = 1.0 - clusterCapConfig[i].GpuMemoryRateCap
			gpuCoreConsumedRate = 1.0 - float64(clusterCapConfig[i].GpuCorePercentCap) / float64(100 + gpuOverSell)
			//log.Printf("scheduler: warm node=%dth, socket=%dth, GPU=%dth, cpuConsumedRate=%f, gpuMemConsumedRate=%f, gpuCoreConsumedRate=%f",
			//	i, j, j+1, cpuConsumedRate, gpuMemConsumedRate, gpuCoreConsumedRate)
			/**
			 * allocate resource
			 */
			maxResourceQuotaNagDiffIndex = -1
			minResourceQuotaPosDiffIndex = -1
			pickConfigIndex = -1
			maxResourceQuotaNagDiff = float64(-999)
			minResourceQuotaPosDiff = float64(999)
			if LessEqual(cpuConsumedRate, gpuMemConsumedRate) && LessEqual(cpuConsumedRate, gpuCoreConsumedRate) { // cpu is dominantly remained resource
				for k := 0; k < len(resourcesConfigs); k++ {
					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20 + cpuOverSell)
					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
					tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
					if Greater(tempDiffQuota,0) {
						if Less(tempDiffQuota, minResourceQuotaPosDiff) {
							minResourceQuotaPosDiff = tempDiffQuota
							minResourceQuotaPosDiffIndex = k
						} else if Equal(tempDiffQuota, minResourceQuotaPosDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							minResourceQuotaPosThroughIntensity := float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, minResourceQuotaPosThroughIntensity) {
								minResourceQuotaPosDiffIndex = k
							}
						}
					} else {
						if Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
							maxResourceQuotaNagDiff = tempDiffQuota
							maxResourceQuotaNagDiffIndex = k
						} else if Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							maxResourceQuotaPosThroughIntensity := float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, maxResourceQuotaPosThroughIntensity) {
								maxResourceQuotaNagDiffIndex = k
							}
						}
					}

				}
				//log.Printf("scheduler: warm CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
			} else if LessEqual(gpuMemConsumedRate, cpuConsumedRate) && LessEqual(gpuMemConsumedRate, gpuCoreConsumedRate) { // GPU mem is dominantly remained resource
				if LessEqual(cpuConsumedRate, gpuCoreConsumedRate) {
					for k := 0; k < len(resourcesConfigs); k++ {
						tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20 + cpuOverSell)
						tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
						tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
						//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
						if Greater(tempDiffQuota,0) {
							if Less(tempDiffQuota, minResourceQuotaPosDiff) {
								minResourceQuotaPosDiff = tempDiffQuota
								minResourceQuotaPosDiffIndex = k
							} else if Equal(tempDiffQuota, minResourceQuotaPosDiff) {
								tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
								minResourceQuotaPosThroughIntensity := float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
									(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
										float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
								if Greater(tempThroughIntensity, minResourceQuotaPosThroughIntensity) {
									minResourceQuotaPosDiffIndex = k
								}
							}
						} else {
							if Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
								maxResourceQuotaNagDiff = tempDiffQuota
								maxResourceQuotaNagDiffIndex = k
							} else if Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
								tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
								maxResourceQuotaPosThroughIntensity := float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
									(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
										float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
								if Greater(tempThroughIntensity, maxResourceQuotaPosThroughIntensity) {
									maxResourceQuotaNagDiffIndex = k
								}
							}
						}
					}
					//log.Printf("scheduler: warm GPU memory and CPU are in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
					//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
				} else {
					for k := 0; k < len(resourcesConfigs); k++ {
						tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20 + cpuOverSell)
						tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
						tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
						//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
						if Greater(tempDiffQuota,0) {
							if Less(tempDiffQuota, minResourceQuotaPosDiff) {
								minResourceQuotaPosDiff = tempDiffQuota
								minResourceQuotaPosDiffIndex = k
							} else if Equal(tempDiffQuota, minResourceQuotaPosDiff) {
								tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
								minResourceQuotaPosThroughIntensity := float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
									(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
										float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
								if Greater(tempThroughIntensity, minResourceQuotaPosThroughIntensity) {
									minResourceQuotaPosDiffIndex = k
								}
							}
						} else {
							if Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
								maxResourceQuotaNagDiff = tempDiffQuota
								maxResourceQuotaNagDiffIndex = k
							} else if Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
								tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
								maxResourceQuotaPosThroughIntensity := float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
									(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
										float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
								if Greater(tempThroughIntensity, maxResourceQuotaPosThroughIntensity) {
									maxResourceQuotaNagDiffIndex = k
								}
							}
						}
					}
					//log.Printf("scheduler: warm GPU memory and GPU are in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
					//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
				}
			} else if LessEqual(gpuCoreConsumedRate, cpuConsumedRate) && LessEqual(gpuCoreConsumedRate, gpuMemConsumedRate) { // GPU core is dominantly remained resource
				for k := 0; k < len(resourcesConfigs); k++ {
					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20 + cpuOverSell)
					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
					tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
					if Greater(tempDiffQuota,0) {
						if Less(tempDiffQuota, minResourceQuotaPosDiff) {
							minResourceQuotaPosDiff = tempDiffQuota
							minResourceQuotaPosDiffIndex = k
						} else if Equal(tempDiffQuota, minResourceQuotaPosDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							minResourceQuotaPosThroughIntensity := float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, minResourceQuotaPosThroughIntensity) {
								minResourceQuotaPosDiffIndex = k
							}
						}
					} else {
						if Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
							maxResourceQuotaNagDiff = tempDiffQuota
							maxResourceQuotaNagDiffIndex = k
						} else if Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							maxResourceQuotaPosThroughIntensity := float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, maxResourceQuotaPosThroughIntensity) {
								maxResourceQuotaNagDiffIndex = k
							}
						}
					}
				}
				//log.Printf("scheduler: warm GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
			} else {
				fmt.Printf("error: in node %d\n",i)
				return
			}
			if minResourceQuotaPosDiffIndex == -1 {
				pickConfigIndex = maxResourceQuotaNagDiffIndex
				//log.Printf("scheduler: warm choosed %dth resourceConfigs with maxResourceQuotaNagDiff=%f\n",
				//	pickConfigIndex, maxResourceQuotaNagDiff)
			} else {
				pickConfigIndex = minResourceQuotaPosDiffIndex
				//log.Printf("scheduler: warm choosed %dth resourceConfigs with minResourceQuotaPosDiff=%f\n",
				//	pickConfigIndex, minResourceQuotaPosDiff)
			}
			// update GPU memory allocation

			/**
			 * find a node to place function pod
			 */
			if clusterCapConfig[i].CpuThreadsCap + cpuOverSell >= resourcesConfigs[pickConfigIndex].CpuThreads &&
				clusterCapConfig[i].GpuCorePercentCap + gpuOverSell >= resourcesConfigs[pickConfigIndex].GpuCorePercent &&
				GreaterEqual(clusterCapConfig[i].GpuMemoryRateCap, resourcesConfigs[pickConfigIndex].GpuMemoryRate) {

				clusterCapConfig[i].CpuThreadsCap -= resourcesConfigs[pickConfigIndex].CpuThreads
				clusterCapConfig[i].GpuCorePercentCap -= resourcesConfigs[pickConfigIndex].GpuCorePercent
				clusterCapConfig[i].GpuMemoryRateCap -= resourcesConfigs[pickConfigIndex].GpuMemoryRate

				cpuTotalConsum+=resourcesConfigs[pickConfigIndex].CpuThreads
				gpuTotalConsum+=resourcesConfigs[pickConfigIndex].GpuCorePercent
				memoryTotalConsum+=resourcesConfigs[pickConfigIndex].GpuMemoryRate
				//	fmt.Printf("place %dth Pod %+v to %dth node\n",l,resourcesConfigs[pickConfigIndex], i)
				break

			} // check the next <CPU socket and GPU> to place function pod
		} // per socket
	}
	fmt.Println("Solve Time: ", time.Since(start))

	boxNum :=0

	for j:=0; j< len(clusterCapConfig); j++ {
		if Equal(clusterCapConfig[j].GpuMemoryRateCap,1.0) &&
			(clusterCapConfig[j].GpuCorePercentCap == 100 ) &&
			clusterCapConfig[j].CpuThreadsCap == 20 {
		} else {
			boxNum++
			fmt.Printf("%f\t%f\t%f \n",
				float64(clusterCapConfig[j].CpuThreadsCap) / float64(20 + cpuOverSell),
				float64(clusterCapConfig[j].GpuCorePercentCap) / float64(100 +gpuOverSell),
				clusterCapConfig[j].GpuMemoryRateCap)
		}

	}
	fmt.Println("Total Box:", boxNum)

	fmt.Println("Optimized Box:")
	fmt.Println(cpuTotalConsum/(20 + cpuOverSell)+1)
	fmt.Println(gpuTotalConsum/(100 +gpuOverSell)+1)
	fmt.Println(memoryTotalConsum)

}
func testScheDRP2(){

	/*latencySLO_100s := []int{231,437,597,309,231,68,375,90,306,450,544,
	161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,
	737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,
	340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,
	738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,
	526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,
	383,693,541,552,728,186,196,657,690,453,602,693,455,248}*/
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468,351,324,439,670,493,533,626,508,123,234,640,182,705,614,205,453,580,631,93,98,687,590,271,384,110,264,574,171,522,173,555,468,325,570,489,733,725,352,174,63,138,431,123,599,80,248,474,516,546,429,395,248,250,595,639,539,531,352,533,178,348,104,285,106,345,698,660,203,598,135,618,254,479,705,571,626,334,92,608,274,729,711,605,487,547,709,242,599,673,730,304,74,235,482,234,365,152,251,253,269,600,173,59,445,277,638,462,619,576,79,56,350,521,675,92,495,221,256,430,63,365,509,399,427,590,298,488,581,279,219,272,649,149,385,573,212,520,460,700,744,380,575,302,515,288,454,55,154,368,735,206,678,689,494,446,110,189,609,573,701,130,452,295,335,209,264,493,624,247,189,299,645,490,425,341,221,551,156,56,352,363,422,490,305,338,690,627,508,92,415,141,324,664,159,514,418,96,497,92,433,128,283,572,230,709,265,87,102,191,254,178,457,351,100,356,398,255,552,657,556,221,598,648,364,56,624,429,284,619,350,328,125,203,312,76,68,229,615,557,715,169,92,92,468,435,729,658,363,179,548,327,50,601,662,366,466,583,710,490,306,565,166,521,571,421,565,699,603,290,260,76,507,628,307,223,415,709,88,427,648,174,468,611,581,640,554,132,313,536,129,419,263,153,742,296,167,73,344,631,606,174,532,704,196,100,268,82,694,559,115,342,218,540,548,687,260,732,726,724,245,608,475,360,680,416,139,388,617,726,698,379,267,87,586,356,358,370,616,227,288,348,527,122,359,71,172,689,176,92,723,704,696,422,82,83,496,422,452,719,555,203,652,268,232,155,586,736,324,188,358,472,349,592,538,172,522,253,551,345,335,432,595,612,331,347,452,700,379,664,649,73,586,56,512,134,108,129,106,304,554,81,167,416,254,560,379,326,554,102,371,103,726,402,551,564,677,567,445,622,391,196,431,581,674,683,151,460,424,480,602,601,610,320,530,219,268,496,360,637,191,682,688,706,448,293,631,658,376,550,387,378,294,535,401,225,391,384,192,677,244,635,291,417,424,301,442,727,170,634,684,728,712,703,284,336,364,291,80,512,312,247,693,662,479,535,612,588,319,284,641,717,263,589,355,244,736,533,252,589,312,616,723,639,577,169,397,321,175,667,381,716,687,701,346,369,289,641,536,709,715,394,83,358,728,687,193,640,711,239,321,480,524,630,602,228,360,563,628,576,320,108,208,336,453,501,350,398,742,162,413,59,237,713,332,268,535,601,224,95,604,92,600,447,638,305,659,222,355,166,315,490,223,630,306,125,371,316,147,269,253,229,368,681,338,394,81,360,715,660,261,190,332,323,320,62,466,363,270,417,121,429,626,462,697,725,63,304,109,84,211,689,670,223,225,285,196,695,113,304,645,192,247,578,687,730,466,283,720,336,95,86,80,221,727,74,192,356,113,305,463,144,628,493,599,651,717,703,261,356,327,400,527,686,308,361,408,162,693,650,201,377,553,590,707,341,539,641,584,314,517,163,600,654,692,472,417,340,581,584,56,509,623,248,540,220,261,75,696,90,423,335,420,351,541,281,657,654,238,176,646,309,117,218,155,263,676,201,70,544,506,55,317,80,546,709,533,629,300,232,581,453,593,132,549,653,216,566,693,302,733,569,518,493,526,123,98,593,675,341,127,505,457,147,650,423,105,278,118,190,696,697,202,716,405,97,280,616,143,541,680,215,132,153,308,123,413,459,476,438,230,538,201,678,511,749,288,716,571,293,214,550,259,705,518,93,707,100,134,661,381,243,297,656,272,541,553,545,230,467,555,120,703,142,526,424,673,438,306,745,385,261,688,688,145,261,478,657,630,209,151,92,468,62,416,740,72,138,439,383,651,691,254,479,472,135,468,238,735,187,375,70,192,333,121,88,589,424,396,413,511,717,655,385,290,133,262,136,672,234,363,79,429,130,176,293,64,576,578,582,135,244,240,89,379,479,154,207,52,673,335,518,410,576,100,161,600,275,562,66,557,645,347,602,275,663,516,67,399,679,538,712,734,213,109,526,102,736,350,675,310,260,356,134,329,573,303,665,243,186,646,211,293,690,535,656,661,504,372,133,473,75,469,443,167,420,567,655,670,57,678,420,620,125,635,619,131,382,71,346,564,342,57,330,601,749,442,513,327,359,448,728,123,321,160,122,321,346,586,559,325,421,434,482,447,391,509,615,137,276,744,710,667,431,411,270,527,132,543,378,492,733,529,284,404,374,747,600,129,205,573,607,386,244,101,726,747,482,307,215,412,180,152,609,82,262,351,424,628,82,83,71,157,75,667,50,589,746,396,662,448,68,475,82,626,659,584,313,674,425,182,256,60,601,284,511,718,653,681,378,672,422,537,95,611,586,658,156,421,404,586,638,288,174,578,439,583,618,579,247,213,57,170,247,98,124,726,461,555,636,391,692,483,699,503,170,450,560,610,573,569,359,695,714,83,189,449,308,160,285,130,196,730,638,626,456,319,186,153,442,280,562,459,723,494,317,250,605,207,678,527,144,184,50,60,594,344,255,477,296,596,103,171,506,669,390,319,599,723,690,457,583,60,445,57,119,616,102,689,313,646,456,425,88,167,403,140,464,170,240,238,144,285,495,496,673,100,129,706,367,616,440,386,448,336,492,481,526,267,428,68,59,202,83,678,605,354,532,205,73,596,684,531,677,522,259,52,199,96,243,229,642,146,189,389,141,262,239,100,443,519,69,185,144,127,131,160,97,558,203,708,398,582,162,288,121,737,496,344,640,494,264,492,316,414,246,59,206,310,503,670,716,736,467,179,497,189,313,638,387,359,687,329,712,116,330,738,141,621,456,68,102,141,303,146,625,337,707,702,702,441,86,496,363,499,373,508,568,140,50,386,249,638,231,622,135,595,240,173,224,622,419,405,294,249,103,115,340,336,412,347,187,345,507,672,115,360,522,218,103,481,80,638,108,408,232,395,192,517,285,135,478,338,607,301,239,429,595,289,501,391,372,460,53,455,383,427,714,411,503,361,356,473,413,52,96,691,327,732,585,710,456,461,447,373,552,292,601,269,54,526,266,621,252,226,241,505,714,355,684,412,419,359,493,606,457,327,107,397,622,266,615,594,670,681,280,418,658,270,612,561,676,174,637,195,421,662,391,88,744,314,151,240,552,599,504,376,626,613,272,691,652,365,533,272,469,592,544,331,554,740,237,676,527,515,305,247,421,598,559,477,537,383,261,384,126,280,397,66,304,559,428,604,335,629,372,586,662,665,623,109,491,619,608,585,701,599,439,724,412,167,445,157,261,70,605,161,136,288,432,279,374,638,395,621,127,174,626,336,453,430,582,324,623,379,187,68,413,600,672,75,609,315,711,623,395,179,248,168,334,501,569,553,260,433,142,397,202,360,196,521,370,432,722,703,488,545,591,159,495,398,217,362,64,54,675,148,712,492,433,81,558,362,679,708,464,386,721,102,208,400,119,373,410,143,569,523,580,744,531,96,726,619,620,734,111,528,437,428,603,521,725,376,466,334,683,56,99,233,341,63,446,637,71,635,182,317,400,639,744,407,68,597,392,645,726,486,403,193,465,726,330,506,64,561,132,387,196,375,361,129,95,126,388,223,539,360,279,241,216,208,462,480,511,240,730,579,415,378,93,186,216,97,312,130,492,631,406,725,399,626,187,55,518,262,419,317,77,259,382,199,384,593,398,682,435,169,738,115,489,658,631,384,234,670,526,299,56,104,357,156,131,448,622,234,697,444,183,310,353,516,68,433,299,720,69,749,214,511,497,276,661,498,119,515,170,419,65,412,326,494,236,674,593,747,689,350,487,67,345,385,631,151,176,247,494,649,465,273,261,285,643,163,280,717,160,179,311,337,295,87,409,343,462,347,77,689,362,96,733,440,734,256,614,103,736,508,68,658,435,745,360,224,548,224,612,336,481,661,337,473,495,675,297,634,415,339,276,593,443,603,170,326,166,323,435,202,264,642,414,693,67,258,551,649,661,730,299,684,186,53,647,639,252,298,102,553,721,497,149,625,88,502,702,152,517,88,53,610,81,92,416,643,84,707,561,173,188,509,635,313,442,300,65,460,575,255,613,589,658,675,274,613,206,225,580,74,693,612,513,549,666,126,508,624,671,357,230,197,498,366,598,537,73,54,197,677,556,149,455,117,221,644,698,331,485,729,407,663,66,518,411,745,714,421,529,471,579,336,479,449,471,388,569,53,230,340,66,727,107,745,691,539,667,281,722,511,577,381,619,748,251,331,732,407,373,638,154,313,439,103,225,452,692,629,341,332,712,556,612,534,334,600,509,157,214,604,565,617,210,248,474,672,456,226,303,622,208,406,739,613,356,547,536,676,140,373,576,320,154,338,167,686,346,732,247,53,483,527,624,577,514,696,218,391,481,655,628,728,283,462,188,280,188,87,194,597,676,205,221,690,395,537,293,748,608,110,62,315,410,360,282,120,714,700,510,153,233,178,723,232,268,678,456,470,563,324,311,155,113,401,393,557,504,528,294,608,367,290,698,571,659,189,523,237,239,223,165,514,104,264,335,57,316,290,408,253,77,602,265,53,511,151,101,541,174,128,203,375,688,475,277,643,224,566,427,257,611,340,736,736,563,660,741,559,706,278,510,543,422,732,485,100,474,478,477,516,422,365,195,264,630,89,204,654,333,618,389,185,62,534,380,713,438,471,738,620,293,408,681,639,219,248,693,89,616,129,250,439,339,654,486,158,131,317,244,245,375,490,312,593,580,742,713,202,300,219,678,670,511,681,254,194,65,63,54,596,679,105,342,422,667,692,68,526,488,397,605,489,83,661,295,346,238,443,441,727,556,451,115,426,120,79,670,673,686,449,327,150,380,600,365,133,678,90,709,737,330,349,265,571,598,372,109,274,636,283,642,189,371,185,594,745,254,684,642,656,666,232,451,208,665,186,685,280,267,171,202,583,530,222,217,121,403,654,324,147,207,186,602,366,645,95,649,396,330,664,343,660,636,245,718,605,559,331,325,414,527,99,155,190,261,177,496,545,256,537,116,483,55,596,569,591,409,574,747,78,197,242,106,562,400,95,632,422,338,606,61,392,639,418,98,651,481,274,732,63,576,644,174,183,340,61,96,526,117,260,446,312,444,617,413,169,722,124,103,423,616,303,124,386,432,376,85,350,572,593,186,747,330,72,600,205,354,249,320,548,329,208,229,481,725,584,526,144,319,587,254,320,371,294,307,147,165,277,248,286,749,517,571,222,141,450,671,97,299,610,320,614,723,286,418,334,174,101,150,741,649,96,245,340,359,76,273,111,76,177,195,484,545,540,445,347,228,597,92,539,642,188,557,367,165,177,259,700,51,496,260,413,605,328,542,224,573,282,676,301,206,197,260,333,370,564,484,226,94,164,638,244,349,84,359,262,460,122,704,100,568,155,443,640,599,747,236,586,257,123,706,134,430,571,314,641,632,683,618,579,598,210,464,635,325,211,624,461,695,442,742,78,113,240,223,338,330,231,672,709,742,636,223,665,346,671,265,552,296,112,204,113,551,67,514,339,339,323,251,701,748,629,319,318,255,461,96,199,407,285,353,746,638,726,341,234,109,662,72,707,543,423,454,261,343,60,706,698,236,581,329,678,695,509,608,479,231,366,441,186,50,583,392,59,362,346,178,684,550,272,511,483,81,656,501,422,417,216,115,92,621,62,528,619,568,193,442,409,308,470,290,491,712,514,722,265,506,318,483,200,115,235,196,742,447,642,72,146,99,597,652,567,369,84,644,555,292,310,355,318,215,710,124,325,237,475,620,678,512,678,342,118,488,725,206,240,694,189,150,186,455,551,110,371,423,132,549,643,275,291,459,618,367,455,237,537,649,402,707,734,449,642,570,227,392,419,297,735,575,123,622,60,333,575,435,435,514,627,692,327,478,362,184,613,724,356,196,472,200,596,601,706,92,693,132,134,445,470,654,350,578,430,542,587,230,283,536,556,666,560,420,642,212,197,159,368,714,131,358,157,360,529,367,515,673,228,551,731,577,58,337,584,291,175,212,179,661,238,183,249,522,740,239,546,707,419,442,557,299,564,572,240,78,542,249,482,167,526,52,365,517,584,551,124,345,460,436,700,79,745,74,655,701,158,241,473,287,427,743,727,636,501,409,219,126,707,613,595,546,632,695,320,96,67,406,210,118,110,201,281,111,105,582,542,591,329,156,147,636,175,339,143,736,464,347,542,112,72,161,337,548,704,665,410,171,654,649,189,78,598,340,441,660,109,260,278,445,255,638,499,120,546,268,168,418,674,715,570,56,271,526,730,228,556,445,746,193,173,727,223,734,476,693,644,396,287,64,180,553,511,87,379,397,166,733,629,109,741,704,110,341,220,452,613,148,644,449,422,605,549,450,567,598,661,657,483,113,303,361,78,524,484,671,702,643,637,206,505,432,50,702,652,384,337,223,398,227,434,626,439,463,476,192,342,143,501,620,216,431,226,535,298,198,647,571,637,437,241,682,691,54,238,394,595,134,386,248,301,141,588,416,374,106,513,114,158,457,313,353,477,495,202,203,205,206,363,644,563,292,77,202,81,184,354,524,111,425,435,554,309,100,375,657,536,170,410,586,159,136,535,455,601,709,281,729,616,528,395,245,676,420,245,343,461,88,392,167,112,225,308,529,645,427,745,102,266,674,152,154,389,500,694,161,203,87,273,483,293,73,127,296,320,644,179,686,704,722,556,124,291,278,724,371,597,339,80,410,397,73,649,710,364,541,370}
	latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468,351,324,439,670,493,533,626,508,123,234,640,182,705,614,205,453,580,631,93,98,687,590,271,384,110,264,574,171,522,173,555,468,325,570,489,733,725,352,174,63,138,431,123,599,80,248,474,516,546,429,395,248,250,595,639,539,531,352,533,178,348,104,285,106,345,698,660,203,598,135,618,254,479,705,571,626,334,92,608,274,729,711,605,487,547,709,242,599,673,730,304,74,235,482,234,365,152,251,253,269,600,173,59,445,277,638,462,619,576,79,56,350,521,675,92,495,221,256,430,63,365,509,399,427,590,298,488,581,279,219,272,649,149,385,573,212,520,460,700,744,380,575,302,515,288,454,55,154,368,735,206,678,689,494,446,110,189,609,573,701,130,452,295,335,209,264,493,624,247,189,299,645,490,425,341,221,551,156,56,352,363,422,490,305,338,690,627,508,92,415,141,324,664,159,514,418,96,497,92,433,128,283,572,230,709,265,87,102,191,254,178,457,351,100,356,398,255,552,657,556,221,598,648,364,56,624,429,284,619,350,328,125,203,312,76,68,229,615,557,715,169,92,92,468,435,729,658,363,179,548,327,50,601,662,366,466,583,710,490,306,565,166,521,571,421,565,699,603,290,260,76,507,628,307,223,415,709,88,427,648,174,468,611,581,640,554,132,313,536,129,419,263,153,742,296,167,73,344,631,606,174,532,704,196,100,268,82,694,559,115,342,218,540,548,687,260,732,726,724,245,608,475,360,680,416,139,388,617,726,698,379,267,87,586,356,358,370,616,227,288,348,527,122,359,71,172,689,176,92,723,704,696,422,82,83,496,422,452,719,555,203,652,268,232,155,586,736,324,188,358,472,349,592,538,172,522,253,551,345,335,432,595,612,331,347,452,700,379,664,649,73,586,56,512,134,108,129,106,304,554,81,167,416,254,560,379,326,554,102,371,103,726,402,551,564,677,567,445,622,391,196,431,581,674,683,151,460,424,480,602,601,610,320,530,219,268,496,360,637,191,682,688,706,448,293,631,658,376,550,387,378,294,535,401,225,391,384,192,677,244,635,291,417,424,301,442,727,170,634,684,728,712,703,284,336,364,291,80,512,312,247,693,662,479,535,612,588,319,284,641,717,263,589,355,244,736,533,252,589,312,616,723,639,577,169,397,321,175,667,381,716,687,701,346,369,289,641,536,709,715,394,83,358,728,687,193,640,711,239,321,480,524,630,602,228,360,563,628,576,320,108,208,336,453,501,350,398,742,162,413,59,237,713,332,268,535,601,224,95,604,92,600,447,638,305,659,222,355,166,315,490,223,630,306,125,371,316,147,269,253,229,368,681,338,394,81,360,715,660,261,190,332,323,320,62,466,363,270,417,121,429,626,462,697,725,63,304,109,84,211,689,670,223,225,285,196,695,113,304,645,192,247,578,687,730,466,283,720,336,95,86,80,221,727,74,192,356,113,305,463,144,628,493,599,651,717,703,261,356,327,400,527,686,308,361,408,162,693,650,201,377,553,590,707,341,539,641,584,314,517,163,600,654,692,472,417,340,581,584,56,509,623,248,540,220,261,75,696,90,423,335,420,351,541,281,657,654,238,176,646,309,117,218,155,263,676,201,70,544,506,55,317,80,546,709,533,629,300,232,581,453,593,132,549,653,216,566,693,302,733,569,518,493,526,123,98,593,675,341,127,505,457,147,650,423,105,278,118,190,696,697,202,716,405,97,280,616,143,541,680,215,132,153,308,123,413,459,476,438,230,538,201,678,511,749,288,716,571,293,214,550,259,705,518,93,707,100,134,661,381,243,297,656,272,541,553,545,230,467,555,120,703,142,526,424,673,438,306,745,385,261,688,688,145,261,478,657,630,209,151,92,468,62,416,740,72,138,439,383,651,691,254,479,472,135,468,238,735,187,375,70,192,333,121,88,589,424,396,413,511,717,655,385,290,133,262,136,672,234,363,79,429,130,176,293,64,576,578,582,135,244,240,89,379,479,154,207,52,673,335,518,410,576,100,161,600,275,562,66,557,645,347,602,275,663,516,67,399,679,538,712,734,213,109,526,102,736,350,675,310,260,356,134,329,573,303,665,243,186,646,211,293,690,535,656,661,504,372,133,473,75,469,443,167,420,567,655,670,57,678,420,620,125,635,619,131,382,71,346,564,342,57,330,601,749,442,513,327,359,448,728,123,321,160,122,321,346,586,559,325,421,434,482,447,391,509,615,137,276,744,710,667,431,411,270,527,132,543,378,492,733,529,284,404,374,747,600,129,205,573,607,386,244,101,726,747,482,307,215,412,180,152,609,82,262,351,424,628,82,83,71,157,75,667,50,589,746,396,662,448,68,475,82,626,659,584,313,674,425,182,256,60,601,284,511,718,653,681,378,672,422,537,95,611,586,658,156,421,404,586,638,288,174,578,439,583,618,579,247,213,57,170,247,98,124,726,461,555,636,391,692,483,699,503,170,450,560,610,573,569,359,695,714,83,189,449,308,160,285,130,196,730,638,626,456,319,186,153,442,280,562,459,723,494,317,250,605,207,678,527,144,184,50,60,594,344,255,477,296,596,103,171,506,669,390,319,599,723,690,457,583,60,445,57,119,616,102,689,313,646,456,425,88,167,403,140,464,170,240,238,144,285,495,496,673,100,129,706,367,616,440,386,448,336,492,481,526,267,428,68,59,202,83,678,605,354,532,205,73,596,684,531,677,522,259,52,199,96,243,229,642,146,189,389,141,262,239,100,443,519,69,185,144,127,131,160,97,558,203,708,398,582,162,288,121,737,496,344,640,494,264,492,316,414,246,59,206,310,503,670,716,736,467,179,497,189,313,638,387,359,687,329,712,116,330,738,141,621,456,68,102,141,303,146,625,337,707,702,702,441,86,496,363,499,373,508,568,140,50,386,249,638,231,622,135,595,240,173,224,622,419,405,294,249,103,115,340,336,412,347,187,345,507,672,115,360,522,218,103,481,80,638,108,408,232,395,192,517,285,135,478,338,607,301,239,429,595,289,501,391,372,460,53,455,383,427,714,411,503,361,356,473,413,52,96,691,327,732,585,710,456,461,447,373,552,292,601,269,54,526,266,621,252,226,241,505,714,355,684,412,419,359,493,606,457,327,107,397,622,266,615,594,670,681,280,418,658,270,612,561,676,174,637,195,421,662,391,88,744,314,151,240,552,599,504,376,626,613,272,691,652,365,533,272,469,592,544,331,554,740,237,676,527,515,305,247,421,598,559,477,537,383,261,384,126,280,397,66,304,559,428,604,335,629,372,586,662,665,623,109,491,619,608,585,701,599,439,724,412,167,445,157,261,70,605,161,136,288,432,279,374,638,395,621,127,174,626,336,453,430,582,324,623,379,187,68,413,600,672,75,609,315,711,623,395,179,248,168,334,501,569,553,260,433,142,397,202,360,196,521,370,432,722,703,488,545,591,159,495,398,217,362,64,54,675,148,712,492,433,81,558,362,679,708,464,386,721,102,208,400,119,373,410,143,569,523,580,744,531,96,726,619,620,734,111,528,437,428,603,521,725,376,466,334,683,56,99,233,341,63,446,637,71,635,182,317,400,639,744,407,68,597,392,645,726,486,403,193,465,726,330,506,64,561,132,387,196,375,361,129,95,126,388,223,539,360,279,241,216,208,462,480,511,240,730,579,415,378,93,186,216,97,312,130,492,631,406,725,399,626,187,55,518,262,419,317,77,259,382,199,384,593,398,682,435,169,738,115,489,658,631,384,234,670,526,299,56,104,357,156,131,448,622,234,697,444,183,310,353,516,68,433,299,720,69,749,214,511,497,276,661,498,119,515,170,419,65,412,326,494,236,674,593,747,689,350,487,67,345,385,631,151,176,247,494,649,465,273,261,285,643,163,280,717,160,179,311,337,295,87,409,343,462,347,77,689,362,96,733,440,734,256,614,103,736,508,68,658,435,745,360,224,548,224,612,336,481,661,337,473,495,675,297,634,415,339,276,593,443,603,170,326,166,323,435,202,264,642,414,693,67,258,551,649,661,730,299,684,186,53,647,639,252,298,102,553,721,497,149,625,88,502,702,152,517,88,53,610,81,92,416,643,84,707,561,173,188,509,635,313,442,300,65,460,575,255,613,589,658,675,274,613,206,225,580,74,693,612,513,549,666,126,508,624,671,357,230,197,498,366,598,537,73,54,197,677,556,149,455,117,221,644,698,331,485,729,407,663,66,518,411,745,714,421,529,471,579,336,479,449,471,388,569,53,230,340,66,727,107,745,691,539,667,281,722,511,577,381,619,748,251,331,732,407,373,638,154,313,439,103,225,452,692,629,341,332,712,556,612,534,334,600,509,157,214,604,565,617,210,248,474,672,456,226,303,622,208,406,739,613,356,547,536,676,140,373,576,320,154,338,167,686,346,732,247,53,483,527,624,577,514,696,218,391,481,655,628,728,283,462,188,280,188,87,194,597,676,205,221,690,395,537,293,748,608,110,62,315,410,360,282,120,714,700,510,153,233,178,723,232,268,678,456,470,563,324,311,155,113,401,393,557,504,528,294,608,367,290,698,571,659,189,523,237,239,223,165,514,104,264,335,57,316,290,408,253,77,602,265,53,511,151,101,541,174,128,203,375,688,475,277,643,224,566,427,257,611,340,736,736,563,660,741,559,706,278,510,543,422,732,485,100,474,478,477,516,422,365,195,264,630,89,204,654,333,618,389,185,62,534,380,713,438,471,738,620,293,408,681,639,219,248,693,89,616,129,250,439,339,654,486,158,131,317,244,245,375,490,312,593,580,742,713,202,300,219,678,670,511,681,254,194,65,63,54,596,679,105,342,422,667,692,68,526,488,397,605,489,83,661,295,346,238,443,441,727,556,451,115,426,120,79,670,673,686,449,327,150,380,600,365,133,678,90,709,737,330,349,265,571,598,372,109,274,636,283,642,189,371,185,594,745,254,684,642,656,666,232,451,208,665,186,685,280,267,171,202,583,530,222,217,121,403,654,324,147,207,186,602,366,645,95,649,396,330,664,343,660,636,245,718,605,559,331,325,414,527,99,155,190,261,177,496,545,256,537,116,483,55,596,569,591,409,574,747,78,197,242,106,562,400,95,632,422,338,606,61,392,639,418,98,651,481,274,732,63,576,644,174,183,340,61,96,526,117,260,446,312,444,617,413,169,722,124,103,423,616,303,124,386,432,376,85,350,572,593,186,747,330,72,600,205,354,249,320,548,329,208,229,481,725,584,526,144,319,587,254,320,371,294,307,147,165,277,248,286,749,517,571,222,141,450,671,97,299,610,320,614,723,286,418,334,174,101,150,741,649,96,245,340,359,76,273,111,76,177,195,484,545,540,445,347,228,597,92,539,642,188,557,367,165,177,259,700,51,496,260,413,605,328,542,224,573,282,676,301,206,197,260,333,370,564,484,226,94,164,638,244,349,84,359,262,460,122,704,100,568,155,443,640,599,747,236,586,257,123,706,134,430,571,314,641,632,683,618,579,598,210,464,635,325,211,624,461,695,442,742,78,113,240,223,338,330,231,672,709,742,636,223,665,346,671,265,552,296,112,204,113,551,67,514,339,339,323,251,701,748,629,319,318,255,461,96,199,407,285,353,746,638,726,341,234,109,662,72,707,543,423,454,261,343,60,706,698,236,581,329,678,695,509,608,479,231,366,441,186,50,583,392,59,362,346,178,684,550,272,511,483,81,656,501,422,417,216,115,92,621,62,528,619,568,193,442,409,308,470,290,491,712,514,722,265,506,318,483,200,115,235,196,742,447,642,72,146,99,597,652,567,369,84,644,555,292,310,355,318,215,710,124,325,237,475,620,678,512,678,342,118,488,725,206,240,694,189,150,186,455,551,110,371,423,132,549,643,275,291,459,618,367,455,237,537,649,402,707,734,449,642,570,227,392,419,297,735,575,123,622,60,333,575,435,435,514,627,692,327,478,362,184,613,724,356,196,472,200,596,601,706,92,693,132,134,445,470,654,350,578,430,542,587,230,283,536,556,666,560,420,642,212,197,159,368,714,131,358,157,360,529,367,515,673,228,551,731,577,58,337,584,291,175,212,179,661,238,183,249,522,740,239,546,707,419,442,557,299,564,572,240,78,542,249,482,167,526,52,365,517,584,551,124,345,460,436,700,79,745,74,655,701,158,241,473,287,427,743,727,636,501,409,219,126,707,613,595,546,632,695,320,96,67,406,210,118,110,201,281,111,105,582,542,591,329,156,147,636,175,339,143,736,464,347,542,112,72,161,337,548,704,665,410,171,654,649,189,78,598,340,441,660,109,260,278,445,255,638,499,120,546,268,168,418,674,715,570,56,271,526,730,228,556,445,746,193,173,727,223,734,476,693,644,396,287,64,180,553,511,87,379,397,166,733,629,109,741,704,110,341,220,452,613,148,644,449,422,605,549,450,567,598,661,657,483,113,303,361,78,524,484,671,702,643,637,206,505,432,50,702,652,384,337,223,398,227,434,626,439,463,476,192,342,143,501,620,216,431,226,535,298,198,647,571,637,437,241,682,691,54,238,394,595,134,386,248,301,141,588,416,374,106,513,114,158,457,313,353,477,495,202,203,205,206,363,644,563,292,77,202,81,184,354,524,111,425,435,554,309,100,375,657,536,170,410,586,159,136,535,455,601,709,281,729,616,528,395,245,676,420,245,343,461,88,392,167,112,225,308,529,645,427,745,102,266,674,152,154,389,500,694,161,203,87,273,483,293,73,127,296,320,644,179,686,704,722,556,124,291,278,724,371,597,339,80,410,397,73,649,710,364,541,370,129,453,334,167,202,365,630,155,144,406,392,51,269,377,467,511,153,50,279,298,641,242,615,700,565,537,669,352,66,373,229,57,445,196,154,678,110,589,650,575,567,62,102,163,683,712,531,403,186,620,622,582,220,405,316,70,323,235,420,335,349,667,230,573,737,668,587,520,712,116,303,429,660,443,78,462,441,384,634,80,78,503,180,700,435,653,309,144,695,146,292,420,275,679,303,513,404,297,273,450,168,91,444,580,573,394,334,616,278,83,335,289,97,657,317,455,464,673,526,652,419,471,224,556,371,237,412,184,574,521,118,449,379,247,411,706,456,597,583,156,135,451,214,536,243,735,133,404,441,444,718,631,97,88,274,736,385,532,736,229,335,444,594,188,204,568,336,736,546,658,635,411,234,280,271,683,588,693,354,518,283,666,133,651,431,383,539,600,267,146,220,436,643,230,613,485,208,228,544,600,269,334,230,731,279,355,716,475,323,355,254,667,418,580,661,348,573,739,345,667,383,320,686,660,554,450,637,611,146,445,674,370,475,682,578,261,716,280,148,99,569,464,749,638,229,526,389,155,211,692,741,701,479,666,384,515,394,456,525,217,530,287,560,710,673,572,696,335,268,675,445,301,629,440,385,104,633,277,544,583,289,124,567,679,737,511,502,629,92,216,239,564,649,566,233,105,496,523,359,513,312,468,187,204,399,625,140,353,563,380,366,653,343,335,749,55,389,278,574,217,353,368,616,98,639,171,513,507,646,712,64,413,312,387,418,742,743,266,625,451,518,508,150,401,700,707,612,586,623,510,597,85,464,139,710,164,644,276,232,568,242,733,271,733,665,414,393,600,457,466,320,385,189,119,197,511,241,694,451,233,55,458,632,625,680,160,168,679,135,159,179,414,554,641,420,595,675,740,313,529,157,580,194,746,430,709,54,258,155,599,720,712,329,119,329,310,431,589,272,148,686,396,357,55,261,597,205,410,616,444,523,690,290,79,614,448,51,454,462,735,560,553,458,488,380,82,191,330,416,651,247,665,546,57,552,55,603,324,270,147,355,650,298,324,339,277,669,239,133,386,64,572,289,683,60,373,476,525,107,587,94,156,531,192,367,718,478,483,388,413,230,430,61,540,508,483,657,638,239,288,163,237,284,340,146,531,446,95,522,575,543,604,526,315,69,203,246,538,139,284,668,605,719,193,732,66,423,311,577,221,294,481,477,732,65,619,649,245,385,684,596,506,613,150,153,741,124,516,671,563,623,739,509,111,320,312,543,76,305,619,691,420,474,170,459,478,561,438,473,573,77,615,305,215,565,451,358,437,301,574,597,586,457,64,223,238,617,137,710,514,237,100,611,734,494,389,646,229,354,325,306,670,612,725,728,256,708,650,597,480,384,230,67,398,66,245,209,342,160,132,214,340,424,153,238,693,539,54,356,662,227,748,355,314,488,440,356,145,421,241,746,561,546,532,225,90,645,715,524,446,259,384,348,537,263,69,648,705,693,339,334,402,332,604,51,224,227,172,512,308,715,453,448,694,503,156,112,135,541,116,652,336,312,564,347,214,301,357,215,713,692,110,51,245,491,455,152,655,267,431,340,321,360,177,389,514,631,185,489,446,437,438,532,145,420,138,522,666,321,391,390,135,509,662,589,322,112,221,628,230,131,743,486,170,237,274,247,560,676,451,436,647,495,623,420,274,241,298,226,93,312,300,572,177,662,375,739,336,600,678,673,612,552,566,212,199,719,246,177,502,228,393,107,602,223,232,249,732,344,630,325,743,508,534,547,130,235,302,271,329,439,673,69,609,455,388,442,354,452,316,491,592,109,452,111,423,735,436,204,129,227,67,290,563,232,338,608,509,658,644,347,583,161,598,728,105,728,105,582,480,138,109,472,731,434,571,566,383,250,439,371,453,236,195,146,403,264,384,317,357,477,633,502,237,691,70,707,366,189,346,54,496,431,733,507,539,201,274,627,609,189,277,148,253,408,283,102,575,507,740,340,232,534,267,457,713,232,305,471,726,424,148,745,118,220,466,319,579,507,438,570,539,169,126,673,117,229,51,310,645,300,223,117,289,317,360,162,693,145,238,453,675,452,50,65,173,157,567,539,666,54,583,613,98,312,229,567,66,740,641,509,333,300,161,670,324,548,271,615,213,638,643,226,515,325,337,324,347,401,715,598,176,127,738,337,385,63,135,591,525,638,183,95,683,171,67,745,111,302,71,255,89,360,80,534,233,481,577,699,556,693,404,653,509,65,746,223,171,139,616,113,298,75,209,733,202,650,496,105,392,556,508,411,592,155,53,473,542,237,521,367,657,353,145,292,247,427,599,173,347,494,264,525,129,273,570,413,685,139,412,427,194,353,229,250,671,739,186,736,109,622,675,663,618,476,327,588,744,198,184,461,74,486,445,400,290,158,61,440,142,629,742,696,563,233,345,89,278,245,426,369,215,604,105,542,347,552,429,444,541,103,428,613,76,468,500,414,361,485,192,595,341,214,198,149,662,270,547,210,252,674,347,372,67,499,708,721,91,258,740,102,489,663,125,421,592,624,61,625,128,662,323,362,580,58,410,356,51,423,80,409,95,478,88,472,635,188,254,739,665,127,62,731,566,151,608,233,483,703,691,366,193,289,218,326,300,140,313,137,196,700,368,318,534,227,593,101,701,682,485,216,125,499,608,62,520,265,317,348,570,160,471,578,320,63,564,515,404,223,171,703,283,306,103,214,217,373,569,149,581,194,170,299,620,233,147,743,304,656,717,72,632,560,98,236,336,749,641,628,698,450,742,110,388,180,586,650,666,536,521,551,330,367,245,431,619,327,320,579,542,594,629,263,302,55,287,476,743,720,231,517,213,727,592,512,596,250,350,399,506,177,141,243,442,420,690,467,151,605,348,440,558,498,619,746,354,111,210,491,736,356,687,138,71,417,519,260,386,307,580,246,614,204,748,169,743,663,563,142,723,265,479,329,309,481,93,78,90,512,521,434,558,350,546,744,386,412,116,412,347,108,737,457,524,246,365,625,634,283,595,181,79,367,127,704,623,264,620,394,687,152,615,651,713,463,232,148,256,50,527,556,509,318,92,166,464,704,392,514,216,177,100,341,565,690,494,650,327,458,479,551,536,434,326,481,119,648,74,154,631,576,207,273,439,193,185,405,263,667,643,73,331,466,354,628,217,147,652,252,415,606,485,453,394,482,633,320,101,653,580,220,222,405,685,521,92,461,448,593,345,484,134,199,497,249,121,649,415,360,563,447,521,169,668,702,519,223,548,297,288,147,78,290,431,612,66,368,490,403,718,657,441,590,118,88,169,68,220,311,196,539,133,660,155,660,183,662,354,596,624,72,659,114,246,696,57,610,676,482,527,674,364,680,254,339,282,75,63,586,744,709,250,567,171,571,720,228,700,376,70,194,612,740,731,483,628,293,681,92,546,559,147,304,707,725,358,643,610,213,178,744,188,354,76,182,517,508,516,582,265,550,661,685,146,310,234,486,332,425,170,262,525,74,391,232,516,235,635,248,287,127,66,476,83,284,490,189,265,97,338,375,746,216,412,152,746,524,58,329,708,649,263,741,260,540,581,482,686,133,540,120,296,705,165,386,299,352,140,741,316,728,427,496,282,310,546,390,461,65,103,93,419,747,234,525,422,621,511,534,78,615,652,690,575,560,139,485,233,140,452,352,445,176,234,510,442,541,611,130,480,623,190,189,540,187,583,380,367,90,254,659,111,553,541,672,573,537,354,745,713,95,581,208,414,354,421,573,605,209,293,110,318,335,182,495,308,599,660,361,93,618,436,710,98,713,649,71,254,250,175,157,446,330,174,602,252,629,575,622,477,474,311,247,148,445,635,91,246,565,265,563,459,628,119,471,329,170,79,117,697,718,198,247,270,272,651,573,664,505,281,565,227,712,531,453,618,348,700,595,381,50,471,556,146,134,486,186,320,380,391,545,319,510,363,164,104,314,73,70,744,493,583,538,158,249,86,348,656,426,466,231,162,295,100,514,655,317,674,339,497,485,706,65,245,563,673,618,102,156,742,691,104,317,75,415,603,679,136,231,471,180,226,693,661,542,680,198,136,202,64,722,746,584,326,632,236,678,454,355,736,203,608,267,255,80,681,358,254,550,746,684,586,573,126,109,354,732,136,482,250,56,228,326,276,662,155,336,654,478,293,405,710,578,378,686,723,205,537,612,696,391,290,77,547,441,128,347,153,463,330,682,504,314,228,179,145,739,438,90,735,306,588,195,663,521,55,532,216,382,522,533,683,382,676,112,432,714,319,76,574,407,507,374,562,82,586,65,707,459,696,61,250,626,536,155,281,665,403,312,659,79,521,449,140,151,315,313,357,382,290,150,191,350,457,648,571,163,413,561,280,379,557,61,686,202,721,179,79,337,216,737,316,119,133,459,355,735,268,667,125,218,116,509,255,652,483,103,560,712,198,70,739,344,581,105,164,467,724,685,382,214,657,161,616,737,488,548,272,116,707,507,548,66,616,374,351,152,251,108,96,354,101,416,373,592,630,576,476,352,216,457,473,549,300,526,211,512,143,221,741,116,113,426,115,523,413,346,486,429,670,161,694,239,357,282,381,605,683,170,497,555,77,538,282,679,58,65,211,578,437,593,443,267,110,657,464,484,113,154,311,635,718,211,310,696,670,195,474,294,672,356,269,743,199,549,406,108,484,63,261,242,113,350,634,502,431,200,655,146,194,103,453,196,63,687,152,185,157,447,289,611,292,430,161,599,562,121,69,707,549,258,456,220,386,179,653,448,259,717,326,87,211,172,648,93,72,175,647,120,460,494,740,457,284,302,240,128,588,266,104,180,595,249,592,422,146,595,680,208,126,535,502,220,194,171,570,83,185,531,556,709,565,441,718,568,732,214,190,601,440,353,462,466,546,175,421,303,559,582,383,522,561,144,535,259,357,495,54,647,656,65,255,642,663,610,684,630,373,561,420,624,718,539,605,281,493,175,78,401,609,645,122,229,172,421,76,191,246,391,329,263,113,740,61,77,55,313,324,257,676,94,630,215,653,452,600,110,407,178,230,197,472,743,714,52,424,311,679,304,607,666,123,368,544,275,201,272,99,701,351,215,191,531,447,304,213,732,419,334,311,615,413,620,330,656,728,69,653,199,162,61,97,474,306,384,161,495,394,50,625,191,337,93,290,317,627,392,630,261,626,611,433,173,741,519,112,84,158,94,102,234,717,535,597,274,505,650,645,467,95,280,451,169,166,624,486,58,625,613,133,333,376,212,458,85,119,694,664,300,158,255,149,540,106,372,589,682,550,279,66,581,300,697,530,535,681,428,297,153,113,420,413,249,367,477,216,210,430,686,673,265,647,325,361,350,259,577,233,675,519,272,471,607,678,380,509,286,66,61,464,455,683,249,107,416,430,93,735,127,442,125,638,465,575,127,91,158,432,569,223,499,211,734,480,682,402,59,395,295,582,256,90,313,615,207,677,181,324,581,453,391,714,487,317,222,242,166,654,569,190,419,661,197,335,549,350,746,181,600,580,237,647,640,108,617,572,231,67,317,476,383,583,630,221,742,657,62,252,376,658,323,404,499,374,274,253,185,365,399,440,408,531,156,86,663,422,158,129,67,337,160,314,335,197,431,373,534,712,688,645,403,472,511,210,199,393,629,109,510,434,371,363,66,364,582,452,112,64,625,340,465,410,290,582,628,630,453,456,673,124,189,254,616,478,498,163,687,405,355,332,278,384,285,181,228,432,561,438,687,80,590,352,134,190,279,62,624,523,427,564,629,459,230,176,555,163,366,485,694,152,288,316,359,116,396,287,342,692,101,379,690,150,369,651,144,622,382,438,413,62,511,607,686,230,126,387,123,381,436,75,644,237,434,588,577,423,657,542,170,513,439,142,739,82,728,295,273,231,465,633,100,207,228,495,393,168,674,602,196,105,515,78,628,716,301,80,440,740,160,551,189,370,304,272,732,597,261,156,286,518,465,56,527,439,84,587,517,375,202,655,127,668,331,341,746,354,368,301,551,254,447,225,515,271,514,162,649,585,462,304,734,368,731,710,645,200,692,327,465,512,432,747,415,559,147,562,60,412,380,327,490,670,556,527,121,185,557,235,484,699,218,418,702,490,133,154,287,58,208,368,670,181,635,636,157,608,536,107,186,505,118,361,498,660,633,153,83,638,648,469,487,629,148,401,243,168,428,556,250,537,208,72,537,253,613,467,163,448,174,480,273,189,310,345,450,117,704,391,340,225,132,659,80,719,64,362,52,502,106,498,605,635,258,76,324,50,698,506,416,207,701,286,647,363,566,188,489,644,287,346,370,534,210,70,578,579,301,524,283,477,412,738,228,404,412,636,740,124,656,193,712,358,259,183,700,627,470,644,277,615,535,606,570,192,186,590,103,51,743,531,541,396,154,349,125,640,296,543,162,335,656,256,53,404,231,623,660,61,147,700,233,90,213,186,657,637,683,291,528,605,230,691,239,626,58,171,301,595,238,528,122,77,491,459,746,341,346,739,343,150,640,99,744,545,552,380,370,534,607,717,142,463,134,309,659,140,745,163,463,563,732,81,613,224,479,501,325,217,717,747,155,511,425,366,652,245,637,467,441,636,491,77,391,135,696,194,424,482,132,683,162,739,323,392,140,65,531,719,262,440,642,85,143,118,168,261,179,629,451,91,259,251,595,262,447,670,74,670,225,259,730,259,645,103,503,644,450,489,234,284,547,237,208,198,601,79,412,355,727,296,686,539,128,206,198,98,502,714,693,736,164,553,185,161,498,728,352,709,193,499,616,177,510,748,464,668,727,693,481,565,287,335,189,218,569,82,237,120,660,579,466,484,559,405,616,294,749,92,233,594,659,664,569,401,564,599,159,434,273,352,240,606,250,53,219,108,59,512,648,740,716,436,182,183,303,214,90,351,174,521,115,431,254,382,736,537,678,151,234,203,148,508,132,547,354,578,94,71,569,147,313,570,565,272,437,276,471,263,344,55,92,659,390,471,721,661,240,539,313,716,363,345,174,420,607,614,185,599,405,465,157,694,189,719,270,629,694,647,436,201,593,725,732,89,497,137,412,96,142,702,172,428,718,168,498,696,361,287,506,536,726,416,362,85,617,429,489,738,95,582,180,724,162,514,724,574,512,584,90,551,113,564,283,458,198,79,563,349,558,355,381,261,253,324,724,117,705,571,601,97,271,111,71,487,351,89,378,566,254,573,99,518,425,565,561,78,313,590,533,109,472,112,356,70,531,444,307,530,158,673,203,240,137,165,419,558,702,370,571,509,374,247,361,183,255,741,469,364,699,653,636,477,442,421,250,55,604,642,264,157,711,217,284,211,343,444,666,469,238,190,197,214,286,453,317,433,112,427,717,235,305,141,492,635,256,687,365,196,204,118,230,657,63,515,713,375,302,600,438,708,105,529,168,382,64,57,132,404,232,603,81,210,391,476,146,189,134,508,215,707,721,305,182,618,150,455,350,184,420,264,327,357,608,629,66,409,372,139,392,96,550,207,599,503,493,319,638,371,677,690,262,57,76,301,472,304,254,625,568,77,655,427,321,679,315,485,334,664,600,733,426,433,717,578,420,142,154,661,416,183,558,196,274,199,300,681,646,103,570,271,509,480,328,366,98,383,320,463,652,394,626,680,466,257,603,555,329,222,166,591,260,252,491,485,619,422,538,269,73,573,116,122,531,424,203,252,747,160,136,539,641,321,575,169,642,421,381,580,189,698,172,746,460,52,619,390,379,509,582,485,727,226,544,318,566,277,392,689,577,120,648,634,397,628,55,365,61,82,324,709,225,82,671,181,123,169,437,687,281,694,155,314,613,733,211,204,572,194,172,502,467,542,383,723,110,544,700,55,623,748,239,155,264,500,382,431,265,309,59,149,318,373,140,695,183,219,303,348,454,554,623,60,195,631,742,307,387,384,239,620,116,146,609,439,536,252,208,449,357,494,237,140,537,84,631,677,492,59,612,605,665,663,379,680,573,727,569,369,676,91,690,446,133,712,79,275,358,57,562,396,668,557,371,234,408,569,423,274,410,677,616,81,524,258,130,718,138,279,551,729,562,668,601,574,196,518,153,99,643,508,333,272,445,689,114,586,327,702,350,602,296,693,145,455,172,694,642,608,274,226,482,426,121,565,716,220,202,338,432,166,624,434,147,632,522,204,646,705,223,626,284,572,218,614,656,627,408,137,728,664,468,434,428,709,391,458,196,361,748,618,454,537,121,70,164,704,120,440,239,521,619,121,605,87,443,516,248,570,242,461,355,571,208,452,706,67,408,137,104,742,54,212,512,745,224,168,577,454,162,282,110,593,367,685,342,395,418,205,393,348,255,577,377,746,259,206,210,114,451,720,231,705,396,369,168,259,285,331,725,669,512,435,137,84,483,53,495,314,462,201,203,696,740,170,131,357,485,118,144,572,726,180,233,517,360,157,685,334,734,541,312,621,600,231,80,139,662,86,211,465,243,491,579,130,295,310,333,665,557,530,73,734,673,81,207,181,649,247,733,536,209,603,735,589,359,704,388,374,519,590,111,213,140,72,668,613,486,57,313,381,187,111,533,566,109,742,60,681,597,168,726,575,634,228,656,315,191,448,242,175,541,609,485,247,227,242,635,233,618,446,173,231,356,73,330,192,694,185,218,735,559,68,669,488,396,625,161,716,598,530,650,203,735,569,561,116,134,158,686,539,261,389,202,401,143,634,129,261,203,563,558,618,437,358,523,112,640,427,86,537,494,633,315,250,725,427,85,348,279,215,104,547,491,176,284,536,743,729,684,64,337,237,86,367,378,400,468,139,87,201,716,125,255,553,193,593,666,363,165,328,381,148,324,55,672,491,437,441,463,679,319,255,332,393,211,399,134,396,674,249,324,137,183,395,193,136,393,219,227,435,540,678,449,410,88,330,224,282,697,300,383,579,169,547,507,370,121,719,78,385,186,216,440,360,641,195,179,219,353,435,264,154,103,279,59,369,591,79,689,563,99,225,302,83,300,739,580,512,136,211,674,316,706,482,169,308,711,342,711,227,706,166,244,514,400,257,340,579,230,619,708,339,195,462,732,687,720,123,411,487,357,256,136,575,474,719,231,418,306,195,668,427,322,555,646,658,669,336,509,443,515,701,216,551,695,715,660,410,589,188,404,676,206,122,398,566,183,144,97,559,189,539,501,727,748,216,662,214,306,623,664,59,474,549,720,56,132,526,108,524,421,302,670,643,67,511,398,568,499,302,135,402,279,540,491,518,172,599,668,236,130,352,204,113,366,374,445,422,635,459,50,742,88,439,690,377,591,80,364,146,547,625,741,571,88,662,145,331,177,68,447,84,643,300,79,300,366,438,178,422,358,62,515,313,252,393,292,666,124,321,311,309,593,412,555,304,675,355,496,516,366,494,736,162,462,446,295,245,432,394,348,167,170,104,354,56,532,705,511,670,187,312,249,637,350,217,665,255,165,285,532,528,667,749,360,633,398,610,342,229,101,97,346,369,695,407,572,78,664,617,186,346,57,693,267,131,526,114,166,473,316,622,455,280,534,342,424,194,345,289,736,136,414,544,456,419,274,360,579,557,331,694,296,410,466,102,147,616,157,114,252,104,670,170,616,492,706,321,263,435,549,681,678,324,137,350,346,474,130,321,745,668,245,364,406,298,659,556,571,103,726,557,737,687,743,308,146,244,530,490,211,155,284,726,563,632,173,131,699,558,414,667,504,128,492,724,677,127,613,468,656,80,653,51,283,191,470,457,552,199,445,474,344,150,241,619,425,552,513,174,264,179,166,625,59,747,549,177,347,146,682,69,183,52,496,610,493,434,79,683,707,290,66,334,491,283,511,261,405,227,648,197,562,570,318,86,473,666,655,248,83,278,519,182,549,339,599,677,379,583,262,245,414,389,518,557,279,51,540,682,254,283,709,436,633,457,203,378,724,80,299,236,61,458,168,491,570,55,278,250,617,85,688,499,460,707,590,195,253,180,315,453,665,348,644,613,108,130,638,433,378,76,111,462,82,205,147,271,465,214,136,230,269,460,286,170,585,406,362,670,676,474,188,424,572,337,426,561,556,118,606,411,321,64,191,121,562,321,502,86,371,115,444,144,627,383,156,592,290,137,431,640,294,374,645,287,159,473,629,420,393,552,260,310,601,541,447,273,530,612,512,674,92,454,226,604,696,422,646,479,315,175,553,400,308,411,140,98,352,421,73,79,707,639,647,384,99,338,150,289,665,425,85,85,328,55,676,563,438,479,120,273,559,541,289,422,202,197,227,358,217,62,535,353,89,684,238,529,137,266,742,440,263,283,277,465,514,700,99,716,441,590,703,696,289,108,666,158,129,731,509,75,88,649,109,593,717,485,366,175,643,103,341,275,388,430,464,317,186,534,437,84,433,377,623,280,254,716,316,744,642,150,505,273,50,364,621,607,199,638,147,213,73,130,364,579,223,98,641,749,238,281,348,130,566,270,168,236,616,745,623,551,156,702,717,279,575,315,443,217,631,477,548,324,606,405,680,173,636,588,276,341,650,388,377,208,246,225,246,106,355,723,247,212,553,469,332,234,653,294,555,248,409,700,561,746,426,284,260,662,241,422,641,663,271,415,601,182,266,337,73,531,202,500,78,182,648,315,507,453,486,233,690,531,602,155,321,287,163,75,247,303,297,518,505,503,520,521,212,418,383,747,731,154,692,449,83,741,156,456,723,531,488,687,504,573,517,387,547,73,307,644,481,673,708,368,645,358,698,702,288,434,127,470,157,158,229,217,438,433,386,394,460,209,322,626,98,525,173,268,477,402,638,98,282,466,320,597,125,339,685,653,741,689,436,275,326,430,314,717,337,612,173,314,285,530,265,531,102,273,182,536,289,94,454,474,408,578,719,210,140,596,592,146,127,651,76,356,556,643,211,74,711,568,585,559,523,263,552,179,494,505,537,117,353,685,656,571,53,212,670,697,235,607,94,347,166,311,99,637,605,385,70,599,375,555,113,188,362,393,592,425,135,185,508,513,649,71,351,602,483,672,444,703,657,327,478,727,324,61,544,345,349,396,408,553,723,526,685,610,75,455,632,362,362,298,90,382,340,623,466,510,172,551,81,172,420,75,309,631,316,677,511,352,336,351,626,453,281,274,426,127,609,652,143,565,664,743,490,515,252,256,104,307,150,545,460,485,641,465,237,551,155,298,632,591,85,671,286,442,569,701,348,154,86,296,293,170,563,566,155,534,526,496,448,259,699,111,562,418,447,306,743,488,313,472,516,266,502,685}

	//fmt.Println(len(latencySLO_500s))
	//fmt.Println(len(latencySLO_1000s))
	//fmt.Println(len(latencySLO_2000s))
	//fmt.Println(len(latencySLO_5000s))
	fmt.Println("scale=", len(latencySLO_100s))


	clusterCapConfig := InitBox(5000)
	cpuOverSell:= int32(0) //threads
	gpuOverSell:= int32(0) //SM percentage

	start := time.Now()

	cpuTotalConsum := int32(0)
	gpuTotalConsum := int32(0)
	memoryTotalConsum := float64(0)

	for l:=0; l<len(latencySLO_100s); l++ {

		resourcesConfigs :=  testEstimator(float64(latencySLO_100s[l]))
		if len(resourcesConfigs) == 0 {
			continue
		}

		cpuConsumedRate := float64(0)
		//gpuMemConsumedRate := float64(0)
		gpuCoreConsumedRate := float64(0)

		maxResourceQuotaNagDiffIndex := -1
		minResourceQuotaPosDiffIndex := -1
		pickConfigIndex := -1
		maxResourceQuotaNagDiff := float64(-999)
		minResourceQuotaPosDiff := float64(999)
		tempGpuCoreQuota := float64(0)
		tempCpuQuota := float64(0)
		tempDiffQuota := float64(0)
		for i := 0; i < len(clusterCapConfig); i++ { // per node
			/** CPU GPU consumed rate **/

			cpuConsumedRate = 1.0 - float64(clusterCapConfig[i].CpuThreadsCap) / float64(20 + cpuOverSell) // cpu usage rate in node i socket j
			//gpuMemConsumedRate = 1.0 - clusterCapConfig[i].GpuMemoryRateCap
			gpuCoreConsumedRate = 1.0 - float64(clusterCapConfig[i].GpuCorePercentCap) / float64(100 + gpuOverSell)
			//log.Printf("scheduler: warm node=%dth, socket=%dth, GPU=%dth, cpuConsumedRate=%f, gpuMemConsumedRate=%f, gpuCoreConsumedRate=%f",
			//	i, j, j+1, cpuConsumedRate, gpuMemConsumedRate, gpuCoreConsumedRate)
			/**
			 * allocate resource
			 */
			maxResourceQuotaNagDiffIndex = -1
			minResourceQuotaPosDiffIndex = -1
			pickConfigIndex = -1
			maxResourceQuotaNagDiff = float64(-999)
			minResourceQuotaPosDiff = float64(999)

			if LessEqual(cpuConsumedRate, gpuCoreConsumedRate) { // cpu is dominantly remained resource
				for k := 0; k < len(resourcesConfigs); k++ {
					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20 + cpuOverSell)
					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
					tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
					if Greater(tempDiffQuota,0) {
						if Less(tempDiffQuota, minResourceQuotaPosDiff) {
							minResourceQuotaPosDiff = tempDiffQuota
							minResourceQuotaPosDiffIndex = k
						} else if Equal(tempDiffQuota, minResourceQuotaPosDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							minResourceQuotaPosThroughIntensity := float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, minResourceQuotaPosThroughIntensity) {
								minResourceQuotaPosDiffIndex = k
							}
						}
					} else {
						if Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
							maxResourceQuotaNagDiff = tempDiffQuota
							maxResourceQuotaNagDiffIndex = k
						} else if Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							maxResourceQuotaPosThroughIntensity := float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, maxResourceQuotaPosThroughIntensity) {
								maxResourceQuotaNagDiffIndex = k
							}
						}
					}

				}
				//log.Printf("scheduler: warm CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
			} else if LessEqual(gpuCoreConsumedRate, cpuConsumedRate) { // GPU core is dominantly remained resource
				for k := 0; k < len(resourcesConfigs); k++ {
					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20 + cpuOverSell)
					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
					tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
					if Greater(tempDiffQuota,0) {
						if Less(tempDiffQuota, minResourceQuotaPosDiff) {
							minResourceQuotaPosDiff = tempDiffQuota
							minResourceQuotaPosDiffIndex = k
						} else if Equal(tempDiffQuota, minResourceQuotaPosDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							minResourceQuotaPosThroughIntensity := float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, minResourceQuotaPosThroughIntensity) {
								minResourceQuotaPosDiffIndex = k
							}
						}
					} else {
						if Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
							maxResourceQuotaNagDiff = tempDiffQuota
							maxResourceQuotaNagDiffIndex = k
						} else if Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
							tempThroughIntensity := float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							maxResourceQuotaPosThroughIntensity := float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(20 + cpuOverSell) +
									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if Greater(tempThroughIntensity, maxResourceQuotaPosThroughIntensity) {
								maxResourceQuotaNagDiffIndex = k
							}
						}
					}
				}
				//log.Printf("scheduler: warm GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
			} else {
				fmt.Printf("error: in node %d\n",i)
				return
			}
			if minResourceQuotaPosDiffIndex == -1 {
				pickConfigIndex = maxResourceQuotaNagDiffIndex
				//log.Printf("scheduler: warm choosed %dth resourceConfigs with maxResourceQuotaNagDiff=%f\n",
				//	pickConfigIndex, maxResourceQuotaNagDiff)
			} else {
				pickConfigIndex = minResourceQuotaPosDiffIndex
				//log.Printf("scheduler: warm choosed %dth resourceConfigs with minResourceQuotaPosDiff=%f\n",
				//	pickConfigIndex, minResourceQuotaPosDiff)
			}
			// update GPU memory allocation

			/**
			 * find a node to place function pod
			 */
			if clusterCapConfig[i].CpuThreadsCap + cpuOverSell >= resourcesConfigs[pickConfigIndex].CpuThreads &&
				clusterCapConfig[i].GpuCorePercentCap + gpuOverSell >= resourcesConfigs[pickConfigIndex].GpuCorePercent &&
				GreaterEqual(clusterCapConfig[i].GpuMemoryRateCap, resourcesConfigs[pickConfigIndex].GpuMemoryRate) {

				clusterCapConfig[i].CpuThreadsCap -= resourcesConfigs[pickConfigIndex].CpuThreads
				clusterCapConfig[i].GpuCorePercentCap -= resourcesConfigs[pickConfigIndex].GpuCorePercent
				clusterCapConfig[i].GpuMemoryRateCap -= resourcesConfigs[pickConfigIndex].GpuMemoryRate

				cpuTotalConsum+=resourcesConfigs[pickConfigIndex].CpuThreads
				gpuTotalConsum+=resourcesConfigs[pickConfigIndex].GpuCorePercent
				memoryTotalConsum+=resourcesConfigs[pickConfigIndex].GpuMemoryRate
				//	fmt.Printf("place %dth Pod %+v to %dth node\n",l,resourcesConfigs[pickConfigIndex], i)
				break

			} // check the next <CPU socket and GPU> to place function pod
		} // per socket
	}
	fmt.Println("Solve Time: ", time.Since(start))

	boxNum :=0

	for j:=0; j< len(clusterCapConfig); j++ {
		if Equal(clusterCapConfig[j].GpuMemoryRateCap,1.0) &&
			(clusterCapConfig[j].GpuCorePercentCap == 100 ) &&
			clusterCapConfig[j].CpuThreadsCap == 20 {
		} else {
			boxNum++
			fmt.Printf("%f\t%f\t%f \n",
				float64(clusterCapConfig[j].CpuThreadsCap) / float64(20 + cpuOverSell),
				float64(clusterCapConfig[j].GpuCorePercentCap) / float64(100 +gpuOverSell),
				clusterCapConfig[j].GpuMemoryRateCap)
		}

	}
	fmt.Println("Total Box:", boxNum)

	fmt.Println("Optimized Box:")
	fmt.Println(cpuTotalConsum/(20 + cpuOverSell)+1)
	fmt.Println(gpuTotalConsum/(100 +gpuOverSell)+1)
	fmt.Println(memoryTotalConsum)

}
func testScheRSWA(){

	latencySLO_100s := []int{231,437,597,309,231,68,375,90,306,450,544,
		161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,
		737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,
		340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,
		738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,
		526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,
		383,693,541,552,728,186,196,657,690,453,602,693,455,248}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468,351,324,439,670,493,533,626,508,123,234,640,182,705,614,205,453,580,631,93,98,687,590,271,384,110,264,574,171,522,173,555,468,325,570,489,733,725,352,174,63,138,431,123,599,80,248,474,516,546,429,395,248,250,595,639,539,531,352,533,178,348,104,285,106,345,698,660,203,598,135,618,254,479,705,571,626,334,92,608,274,729,711,605,487,547,709,242,599,673,730,304,74,235,482,234,365,152,251,253,269,600,173,59,445,277,638,462,619,576,79,56,350,521,675,92,495,221,256,430,63,365,509,399,427,590,298,488,581,279,219,272,649,149,385,573,212,520,460,700,744,380,575,302,515,288,454,55,154,368,735,206,678,689,494,446,110,189,609,573,701,130,452,295,335,209,264,493,624,247,189,299,645,490,425,341,221,551,156,56,352,363,422,490,305,338,690,627,508,92,415,141,324,664,159,514,418,96,497,92,433,128,283,572,230,709,265,87,102,191,254,178,457,351,100,356,398,255,552,657,556,221,598,648,364,56,624,429,284,619,350,328,125,203,312,76,68,229,615,557,715,169,92,92,468,435,729,658,363,179,548,327,50,601,662,366,466,583,710,490,306,565,166,521,571,421,565,699,603,290,260,76,507,628,307,223,415,709,88,427,648,174,468,611,581,640,554,132,313,536,129,419,263,153,742,296,167,73,344,631,606,174,532,704,196,100,268,82,694,559,115,342,218,540,548,687,260,732,726,724,245,608,475,360,680,416,139,388,617,726,698,379,267,87,586,356,358,370,616,227,288,348,527,122,359,71,172,689,176,92,723,704,696,422,82,83,496,422,452,719,555,203,652,268,232,155,586,736,324,188,358,472,349,592,538,172,522,253,551,345,335,432,595,612,331,347,452,700,379,664,649,73,586,56,512,134,108,129,106,304,554,81,167,416,254,560,379,326,554,102,371,103,726,402,551,564,677,567,445,622,391,196,431,581,674,683,151,460,424,480,602,601,610,320,530,219,268,496,360,637,191,682,688,706,448,293,631,658,376,550,387,378,294,535,401,225,391,384,192,677,244,635,291,417,424,301,442,727,170,634,684,728,712,703,284,336,364,291,80,512,312,247,693,662,479,535,612,588,319,284,641,717,263,589,355,244,736,533,252,589,312,616,723,639,577,169,397,321,175,667,381,716,687,701,346,369,289,641,536,709,715,394,83,358,728,687,193,640,711,239,321,480,524,630,602,228,360,563,628,576,320,108,208,336,453,501,350,398,742,162,413,59,237,713,332,268,535,601,224,95,604,92,600,447,638,305,659,222,355,166,315,490,223,630,306,125,371,316,147,269,253,229,368,681,338,394,81,360,715,660,261,190,332,323,320,62,466,363,270,417,121,429,626,462,697,725,63,304,109,84,211,689,670,223,225,285,196,695,113,304,645,192,247,578,687,730,466,283,720,336,95,86,80,221,727,74,192,356,113,305,463,144,628,493,599,651,717,703,261,356,327,400,527,686,308,361,408,162,693,650,201,377,553,590,707,341,539,641,584,314,517,163,600,654,692,472,417,340,581,584,56,509,623,248,540,220,261,75,696,90,423,335,420,351,541,281,657,654,238,176,646,309,117,218,155,263,676,201,70,544,506,55,317,80,546,709,533,629,300,232,581,453,593,132,549,653,216,566,693,302,733,569,518,493,526,123,98,593,675,341,127,505,457,147,650,423,105,278,118,190,696,697,202,716,405,97,280,616,143,541,680,215,132,153,308,123,413,459,476,438,230,538,201,678,511,749,288,716,571,293,214,550,259,705,518,93,707,100,134,661,381,243,297,656,272,541,553,545,230,467,555,120,703,142,526,424,673,438,306,745,385,261,688,688,145,261,478,657,630,209,151,92,468,62,416,740,72,138,439,383,651,691,254,479,472,135,468,238,735,187,375,70,192,333,121,88,589,424,396,413,511,717,655,385,290,133,262,136,672,234,363,79,429,130,176,293,64,576,578,582,135,244,240,89,379,479,154,207,52,673,335,518,410,576,100,161,600,275,562,66,557,645,347,602,275,663,516,67,399,679,538,712,734,213,109,526,102,736,350,675,310,260,356,134,329,573,303,665,243,186,646,211,293,690,535,656,661,504,372,133,473,75,469,443,167,420,567,655,670,57,678,420,620,125,635,619,131,382,71,346,564,342,57,330,601,749,442,513,327,359,448,728,123,321,160,122,321,346,586,559,325,421,434,482,447,391,509,615,137,276,744,710,667,431,411,270,527,132,543,378,492,733,529,284,404,374,747,600,129,205,573,607,386,244,101,726,747,482,307,215,412,180,152,609,82,262,351,424,628,82,83,71,157,75,667,50,589,746,396,662,448,68,475,82,626,659,584,313,674,425,182,256,60,601,284,511,718,653,681,378,672,422,537,95,611,586,658,156,421,404,586,638,288,174,578,439,583,618,579,247,213,57,170,247,98,124,726,461,555,636,391,692,483,699,503,170,450,560,610,573,569,359,695,714,83,189,449,308,160,285,130,196,730,638,626,456,319,186,153,442,280,562,459,723,494,317,250,605,207,678,527,144,184,50,60,594,344,255,477,296,596,103,171,506,669,390,319,599,723,690,457,583,60,445,57,119,616,102,689,313,646,456,425,88,167,403,140,464,170,240,238,144,285,495,496,673,100,129,706,367,616,440,386,448,336,492,481,526,267,428,68,59,202,83,678,605,354,532,205,73,596,684,531,677,522,259,52,199,96,243,229,642,146,189,389,141,262,239,100,443,519,69,185,144,127,131,160,97,558,203,708,398,582,162,288,121,737,496,344,640,494,264,492,316,414,246,59,206,310,503,670,716,736,467,179,497,189,313,638,387,359,687,329,712,116,330,738,141,621,456,68,102,141,303,146,625,337,707,702,702,441,86,496,363,499,373,508,568,140,50,386,249,638,231,622,135,595,240,173,224,622,419,405,294,249,103,115,340,336,412,347,187,345,507,672,115,360,522,218,103,481,80,638,108,408,232,395,192,517,285,135,478,338,607,301,239,429,595,289,501,391,372,460,53,455,383,427,714,411,503,361,356,473,413,52,96,691,327,732,585,710,456,461,447,373,552,292,601,269,54,526,266,621,252,226,241,505,714,355,684,412,419,359,493,606,457,327,107,397,622,266,615,594,670,681,280,418,658,270,612,561,676,174,637,195,421,662,391,88,744,314,151,240,552,599,504,376,626,613,272,691,652,365,533,272,469,592,544,331,554,740,237,676,527,515,305,247,421,598,559,477,537,383,261,384,126,280,397,66,304,559,428,604,335,629,372,586,662,665,623,109,491,619,608,585,701,599,439,724,412,167,445,157,261,70,605,161,136,288,432,279,374,638,395,621,127,174,626,336,453,430,582,324,623,379,187,68,413,600,672,75,609,315,711,623,395,179,248,168,334,501,569,553,260,433,142,397,202,360,196,521,370,432,722,703,488,545,591,159,495,398,217,362,64,54,675,148,712,492,433,81,558,362,679,708,464,386,721,102,208,400,119,373,410,143,569,523,580,744,531,96,726,619,620,734,111,528,437,428,603,521,725,376,466,334,683,56,99,233,341,63,446,637,71,635,182,317,400,639,744,407,68,597,392,645,726,486,403,193,465,726,330,506,64,561,132,387,196,375,361,129,95,126,388,223,539,360,279,241,216,208,462,480,511,240,730,579,415,378,93,186,216,97,312,130,492,631,406,725,399,626,187,55,518,262,419,317,77,259,382,199,384,593,398,682,435,169,738,115,489,658,631,384,234,670,526,299,56,104,357,156,131,448,622,234,697,444,183,310,353,516,68,433,299,720,69,749,214,511,497,276,661,498,119,515,170,419,65,412,326,494,236,674,593,747,689,350,487,67,345,385,631,151,176,247,494,649,465,273,261,285,643,163,280,717,160,179,311,337,295,87,409,343,462,347,77,689,362,96,733,440,734,256,614,103,736,508,68,658,435,745,360,224,548,224,612,336,481,661,337,473,495,675,297,634,415,339,276,593,443,603,170,326,166,323,435,202,264,642,414,693,67,258,551,649,661,730,299,684,186,53,647,639,252,298,102,553,721,497,149,625,88,502,702,152,517,88,53,610,81,92,416,643,84,707,561,173,188,509,635,313,442,300,65,460,575,255,613,589,658,675,274,613,206,225,580,74,693,612,513,549,666,126,508,624,671,357,230,197,498,366,598,537,73,54,197,677,556,149,455,117,221,644,698,331,485,729,407,663,66,518,411,745,714,421,529,471,579,336,479,449,471,388,569,53,230,340,66,727,107,745,691,539,667,281,722,511,577,381,619,748,251,331,732,407,373,638,154,313,439,103,225,452,692,629,341,332,712,556,612,534,334,600,509,157,214,604,565,617,210,248,474,672,456,226,303,622,208,406,739,613,356,547,536,676,140,373,576,320,154,338,167,686,346,732,247,53,483,527,624,577,514,696,218,391,481,655,628,728,283,462,188,280,188,87,194,597,676,205,221,690,395,537,293,748,608,110,62,315,410,360,282,120,714,700,510,153,233,178,723,232,268,678,456,470,563,324,311,155,113,401,393,557,504,528,294,608,367,290,698,571,659,189,523,237,239,223,165,514,104,264,335,57,316,290,408,253,77,602,265,53,511,151,101,541,174,128,203,375,688,475,277,643,224,566,427,257,611,340,736,736,563,660,741,559,706,278,510,543,422,732,485,100,474,478,477,516,422,365,195,264,630,89,204,654,333,618,389,185,62,534,380,713,438,471,738,620,293,408,681,639,219,248,693,89,616,129,250,439,339,654,486,158,131,317,244,245,375,490,312,593,580,742,713,202,300,219,678,670,511,681,254,194,65,63,54,596,679,105,342,422,667,692,68,526,488,397,605,489,83,661,295,346,238,443,441,727,556,451,115,426,120,79,670,673,686,449,327,150,380,600,365,133,678,90,709,737,330,349,265,571,598,372,109,274,636,283,642,189,371,185,594,745,254,684,642,656,666,232,451,208,665,186,685,280,267,171,202,583,530,222,217,121,403,654,324,147,207,186,602,366,645,95,649,396,330,664,343,660,636,245,718,605,559,331,325,414,527,99,155,190,261,177,496,545,256,537,116,483,55,596,569,591,409,574,747,78,197,242,106,562,400,95,632,422,338,606,61,392,639,418,98,651,481,274,732,63,576,644,174,183,340,61,96,526,117,260,446,312,444,617,413,169,722,124,103,423,616,303,124,386,432,376,85,350,572,593,186,747,330,72,600,205,354,249,320,548,329,208,229,481,725,584,526,144,319,587,254,320,371,294,307,147,165,277,248,286,749,517,571,222,141,450,671,97,299,610,320,614,723,286,418,334,174,101,150,741,649,96,245,340,359,76,273,111,76,177,195,484,545,540,445,347,228,597,92,539,642,188,557,367,165,177,259,700,51,496,260,413,605,328,542,224,573,282,676,301,206,197,260,333,370,564,484,226,94,164,638,244,349,84,359,262,460,122,704,100,568,155,443,640,599,747,236,586,257,123,706,134,430,571,314,641,632,683,618,579,598,210,464,635,325,211,624,461,695,442,742,78,113,240,223,338,330,231,672,709,742,636,223,665,346,671,265,552,296,112,204,113,551,67,514,339,339,323,251,701,748,629,319,318,255,461,96,199,407,285,353,746,638,726,341,234,109,662,72,707,543,423,454,261,343,60,706,698,236,581,329,678,695,509,608,479,231,366,441,186,50,583,392,59,362,346,178,684,550,272,511,483,81,656,501,422,417,216,115,92,621,62,528,619,568,193,442,409,308,470,290,491,712,514,722,265,506,318,483,200,115,235,196,742,447,642,72,146,99,597,652,567,369,84,644,555,292,310,355,318,215,710,124,325,237,475,620,678,512,678,342,118,488,725,206,240,694,189,150,186,455,551,110,371,423,132,549,643,275,291,459,618,367,455,237,537,649,402,707,734,449,642,570,227,392,419,297,735,575,123,622,60,333,575,435,435,514,627,692,327,478,362,184,613,724,356,196,472,200,596,601,706,92,693,132,134,445,470,654,350,578,430,542,587,230,283,536,556,666,560,420,642,212,197,159,368,714,131,358,157,360,529,367,515,673,228,551,731,577,58,337,584,291,175,212,179,661,238,183,249,522,740,239,546,707,419,442,557,299,564,572,240,78,542,249,482,167,526,52,365,517,584,551,124,345,460,436,700,79,745,74,655,701,158,241,473,287,427,743,727,636,501,409,219,126,707,613,595,546,632,695,320,96,67,406,210,118,110,201,281,111,105,582,542,591,329,156,147,636,175,339,143,736,464,347,542,112,72,161,337,548,704,665,410,171,654,649,189,78,598,340,441,660,109,260,278,445,255,638,499,120,546,268,168,418,674,715,570,56,271,526,730,228,556,445,746,193,173,727,223,734,476,693,644,396,287,64,180,553,511,87,379,397,166,733,629,109,741,704,110,341,220,452,613,148,644,449,422,605,549,450,567,598,661,657,483,113,303,361,78,524,484,671,702,643,637,206,505,432,50,702,652,384,337,223,398,227,434,626,439,463,476,192,342,143,501,620,216,431,226,535,298,198,647,571,637,437,241,682,691,54,238,394,595,134,386,248,301,141,588,416,374,106,513,114,158,457,313,353,477,495,202,203,205,206,363,644,563,292,77,202,81,184,354,524,111,425,435,554,309,100,375,657,536,170,410,586,159,136,535,455,601,709,281,729,616,528,395,245,676,420,245,343,461,88,392,167,112,225,308,529,645,427,745,102,266,674,152,154,389,500,694,161,203,87,273,483,293,73,127,296,320,644,179,686,704,722,556,124,291,278,724,371,597,339,80,410,397,73,649,710,364,541,370}
	//latencySLO_100s :=[]int{231,437,597,309,231,68,375,90,306,450,544,161,712,739,478,124,561,595,687,356,445,416,578,308,597,697,737,238,740,465,491,158,637,481,179,406,687,281,335,376,63,340,644,113,483,597,128,474,609,103,707,71,139,449,450,355,738,88,153,105,601,260,555,606,216,278,211,352,433,696,513,526,352,668,697,544,127,113,546,370,673,303,687,683,391,109,383,693,541,552,728,186,196,657,690,453,602,693,455,248,175,501,265,207,237,60,260,135,540,482,348,603,441,632,734,447,217,87,421,144,276,252,131,329,216,420,343,236,169,131,102,425,735,160,237,399,478,468,634,453,74,197,662,182,666,89,290,736,201,326,390,601,594,314,255,333,151,540,52,708,417,281,128,204,72,273,92,58,393,718,716,60,685,190,454,612,107,65,121,489,580,663,350,309,670,533,220,334,397,260,415,612,179,570,598,606,145,216,350,206,579,442,481,727,736,570,549,512,97,342,738,361,53,138,68,706,69,157,407,602,225,531,603,545,67,643,320,346,336,82,270,110,572,179,411,310,570,629,704,714,610,401,531,507,66,150,389,87,683,711,454,135,159,465,269,664,290,312,690,50,534,323,185,493,271,740,624,419,620,86,88,522,494,245,124,687,374,143,556,598,102,270,472,321,167,601,515,632,229,142,79,581,71,105,561,693,303,516,309,417,558,358,278,203,124,634,125,624,356,187,636,717,464,508,709,696,73,574,337,407,424,689,281,296,733,178,366,676,327,362,386,301,355,712,739,148,637,480,559,78,109,211,52,368,495,121,303,590,370,431,555,626,259,516,706,345,530,384,727,517,428,125,643,230,58,181,315,608,234,392,73,212,623,322,303,308,477,646,98,163,81,724,356,190,581,382,227,52,604,233,280,536,137,312,654,71,343,248,106,615,572,95,144,156,149,574,570,194,509,57,174,207,303,424,325,590,539,364,95,489,631,481,101,335,179,653,245,242,322,525,227,512,723,167,209,457,130,55,56,224,694,273,244,572,483,624,460,680,645,64,473,600,59,727,595,221,658,648,184,729,504,571,537,476,360,671,317,551,299,138,70,724,251,490,123,646,416,625,160,522,163,264,309,75,696,71,494,113,686,537,538,473,375,76,469,450,264,580,626,350,401,84,569,716,363,295,391,592,289,444,71,724,351,300,335,538,514,510,366,269,132,666,363,466,372,151,505,627,713,367,408,382,229,125,83,290,77,284,502,478,187,497,519,662,696,659,354,506,564,339,573,92,614,589,580,269,203,210,444,535,371,463,651,378,707,334,612,233,671,300,265,490,442,292,241,748,51,357,162,645,466,602,633,105,319,490,244,648,640,467,332,564,420,632,218,226,576,646,375,318,112,693,701,471,563,159,743,220,687,636,191,740,248,435,251,308,68,442,385,265,731,98,506,367,624,349,123,401,711,510,642,574,464,568,313,85,98,500,290,202,340,516,749,710,444,360,101,316,51,372,710,684,163,312,448,367,243,690,403,468,451,456,538,421,727,392,367,657,434,300,257,272,295,185,398,680,254,612,502,325,264,457,703,251,423,251,745,400,295,477,449,588,549,265,425,593,290,310,89,543,301,262,228,231,456,429,365,335,361,498,364,243,386,205,510,186,456,233,478,151,562,629,613,93,637,582,273,273,575,740,227,351,551,97,58,657,709,564,346,302,114,638,487,650,442,288,345,147,94,328,58,77,713,404,402,461,351,488,622,338,530,191,121,91,264,377,186,288,134,494,187,648,573,444,245,212,685,329,457,478,732,145,306,128,739,422,191,165,452,114,70,439,53,454,250,710,129,737,165,635,448,713,375,503,701,530,282,612,682,508,414,488,299,138,259,566,406,137,422,358,594,583,296,593,351,423,594,55,50,56,705,613,428,641,211,390,562,228,394,586,362,333,561,85,521,549,408,254,551,533,266,331,746,106,639,566,67,724,727,353,516,311,534,573,279,361,633,154,500,712,541,597,460,687,514,723,507,53,631,488,567,214,726,516,285,55,435,455,529,491,713,654,269,497,133,686,460,653,180,297,457,522,454,348,341,175,387,122,393,479,654,364,239,568,542,508,300,330,328,667,75,163,319,185,660,602,121,122,459,626,237,538,596,68,383,268,733,516,643,124,71,423,398,296,154,104,213,610,141,409,382,247,549,669,449,475,633,228,280,407,167,252,72,381,222,678,732,255,668,335,643,98,675,566,121,565,726,631,744,391,59,614,208,92,140,587,471,694,161,63,374,461,483,192,632,468,276,320,520,512,535,479,213,672,701,607,114,732,702,104,123,597,279,635,234,399,109,567,446,50,522,522,320,548,543,59,295,284,192,555,225,204,241,177,536,545,566,593,292,160,405,293,353,267,654,579,70,190,587,332,303,190,353,595,76,353,142,615,272,352,142,449,438,131,337,188,434,156,684,617,615,475,556,716,666,77,592,721,657,636,505,55,405,80,363,562,610,525,233,666,162,368,379,338,647,232,533,140,147,219,234,621,609,371,444,305,241,563,424,533,347,159,131,155,175,252,631,542,524,710,657,137,321,730,91,515,325,174,483,380,592,502,606,272,189,649,342,637,693,112,382,685,404,721,284,212,416,651,623,428,665,730,398,172,670,541,106,109,249,452,294,511,162,166,266,664,282,606,551,514,434,290,581,471,705,376,599,556,616,80,516,737,629,107,671,351,464,259,360,335,538,597,408,269,62,362,55,242,569,491,705,80,103,216,587,374,337,598,211,136,237,608,234,65,746,336,183,296,247,564,447,494,668,553,639,184,185,431,524,541,546,154,706,392,56,734,142,267,503,599,225,231,523,691,700,215,171,564,484,726,592,438,702,110,181,210,567,252,521,584,256,559,532,567,402,289,610,334,326,64,691,650,155,152,539,597,223,209,282,360,167,503,734,328,612,731,210,449,106,62,220,119,251,669,172,566,275,242,597,246,651,719,399,466,464,433,671,87,58,197,227,698,692,173,731,665,167,269,225,407,345,500,469,710,526,81,697,740,233,158,658,312,337,518,672,264,121,666,641,388,204,144,676,575,441,304,421,731,117,481,74,502,174,252,152,602,511,524,741,574,179,586,377,616,107,404,79,579,526,56,182,513,547,371,160,693,116,581,356,65,257,620,148,435,438,539,584,310,347,339,73,666,382,99,404,175,247,264,204,481,579,136,493,672,490,721,309,740,64,120,391,644,745,615,548,336,631,100,651,118,740,284,245,474,117,70,480,514,527,504,104,551,125,291,669,163,324,221,527,661,498,713,154,398,469,377,308,470,659,562,744,539,731,175,394,431,590,188,633,461,315,402,152,708,723,538,725,680,623,539,65,98,103,588,337,705,330,530,373,748,535,675,262,147,64,268,504,564,500,312,62,634,569,285,200,599,446,55,240,404,178,59,629,408,210,288,640,615,659,699,566,247,725,562,456,319,573,56,682,361,383,276,636,295,171,280,728,358,323,137,387,431,589,327,631,242,444,741,229,609,353,112,282,127,416,201,379,560,423,724,455,539,567,483,515,278,139,245,381,161,520,657,488,255,450,633,326,290,220,449,726,70,175,86,377,317,204,510,633,596,56,282,165,204,732,686,514,193,681,326,463,467,396,218,444,580,175,150,492,583,183,115,576,184,734,259,281,457,146,436,589,592,398,421,126,80,81,515,152,650,68,302,130,134,471,723,252,77,374,234,474,534,735,577,273,672,421,282,710,271,707,312,552,151,509,241,633,667,237,419,250,412,746,277,212,138,717,144,132,187,155,525,643,166,531,57,554,648,574,88,345,491,71,416,243,427,691,356,599,115,706,681,147,385,68,652,653,455,494,287,283,416,82,53,60,428,390,230,541,262,255,576,687,739,93,375,460,305,187,183,695,374,369,142,678,681,182,66,486,542,576,237,216,77,332,658,135,302,670,313,320,126,378,743,526,737,490,566,290,581,355,215,438,510,138,746,103,669,640,684,502,60,363,354,729,287,525,131,249,304,687,573,104,270,145,221,687,597,360,558,172,456,470,410,616,431,172,746,529,623,337,67,477,567,206,591,744,222,459,452,114,55,663,97,486,246,580,660,280,539,578,488,601,393,624,408,69,738,424,650,207,66,234,92,661,595,366,55,60,716,576,420,462,589,86,216,298,365,723,210,748,271,296,655,339,414,199,398,738,282,85,226,434,119,288,81,636,553,651,677,452,584,214,177,96,600,332,726,375,618,254,113,486,281,670,167,205,102,542,658,120,296,274,195,399,305,350,676,159,538,712,635,284,215,203,83,184,244,686,66,93,714,632,381,153,339,298,641,244,606,696,649,50,699,179,583,74,326,94,70,151,577,107,628,502,560,299,687,342,576,242,370,487,520,83,222,87,659,588,589,105,608,370,102,516,126,427,591,472,141,417,496,302,568,560,114,55,549,521,617,333,291,370,377,114,215,117,92,404,77,109,748,291,744,190,73,729,342,417,170,96,50,68,284,406,616,746,226,290,409,78,449,127,61,357,127,309,185,703,433,451,327,331,562,497,514,367,508,574,468,351,324,439,670,493,533,626,508,123,234,640,182,705,614,205,453,580,631,93,98,687,590,271,384,110,264,574,171,522,173,555,468,325,570,489,733,725,352,174,63,138,431,123,599,80,248,474,516,546,429,395,248,250,595,639,539,531,352,533,178,348,104,285,106,345,698,660,203,598,135,618,254,479,705,571,626,334,92,608,274,729,711,605,487,547,709,242,599,673,730,304,74,235,482,234,365,152,251,253,269,600,173,59,445,277,638,462,619,576,79,56,350,521,675,92,495,221,256,430,63,365,509,399,427,590,298,488,581,279,219,272,649,149,385,573,212,520,460,700,744,380,575,302,515,288,454,55,154,368,735,206,678,689,494,446,110,189,609,573,701,130,452,295,335,209,264,493,624,247,189,299,645,490,425,341,221,551,156,56,352,363,422,490,305,338,690,627,508,92,415,141,324,664,159,514,418,96,497,92,433,128,283,572,230,709,265,87,102,191,254,178,457,351,100,356,398,255,552,657,556,221,598,648,364,56,624,429,284,619,350,328,125,203,312,76,68,229,615,557,715,169,92,92,468,435,729,658,363,179,548,327,50,601,662,366,466,583,710,490,306,565,166,521,571,421,565,699,603,290,260,76,507,628,307,223,415,709,88,427,648,174,468,611,581,640,554,132,313,536,129,419,263,153,742,296,167,73,344,631,606,174,532,704,196,100,268,82,694,559,115,342,218,540,548,687,260,732,726,724,245,608,475,360,680,416,139,388,617,726,698,379,267,87,586,356,358,370,616,227,288,348,527,122,359,71,172,689,176,92,723,704,696,422,82,83,496,422,452,719,555,203,652,268,232,155,586,736,324,188,358,472,349,592,538,172,522,253,551,345,335,432,595,612,331,347,452,700,379,664,649,73,586,56,512,134,108,129,106,304,554,81,167,416,254,560,379,326,554,102,371,103,726,402,551,564,677,567,445,622,391,196,431,581,674,683,151,460,424,480,602,601,610,320,530,219,268,496,360,637,191,682,688,706,448,293,631,658,376,550,387,378,294,535,401,225,391,384,192,677,244,635,291,417,424,301,442,727,170,634,684,728,712,703,284,336,364,291,80,512,312,247,693,662,479,535,612,588,319,284,641,717,263,589,355,244,736,533,252,589,312,616,723,639,577,169,397,321,175,667,381,716,687,701,346,369,289,641,536,709,715,394,83,358,728,687,193,640,711,239,321,480,524,630,602,228,360,563,628,576,320,108,208,336,453,501,350,398,742,162,413,59,237,713,332,268,535,601,224,95,604,92,600,447,638,305,659,222,355,166,315,490,223,630,306,125,371,316,147,269,253,229,368,681,338,394,81,360,715,660,261,190,332,323,320,62,466,363,270,417,121,429,626,462,697,725,63,304,109,84,211,689,670,223,225,285,196,695,113,304,645,192,247,578,687,730,466,283,720,336,95,86,80,221,727,74,192,356,113,305,463,144,628,493,599,651,717,703,261,356,327,400,527,686,308,361,408,162,693,650,201,377,553,590,707,341,539,641,584,314,517,163,600,654,692,472,417,340,581,584,56,509,623,248,540,220,261,75,696,90,423,335,420,351,541,281,657,654,238,176,646,309,117,218,155,263,676,201,70,544,506,55,317,80,546,709,533,629,300,232,581,453,593,132,549,653,216,566,693,302,733,569,518,493,526,123,98,593,675,341,127,505,457,147,650,423,105,278,118,190,696,697,202,716,405,97,280,616,143,541,680,215,132,153,308,123,413,459,476,438,230,538,201,678,511,749,288,716,571,293,214,550,259,705,518,93,707,100,134,661,381,243,297,656,272,541,553,545,230,467,555,120,703,142,526,424,673,438,306,745,385,261,688,688,145,261,478,657,630,209,151,92,468,62,416,740,72,138,439,383,651,691,254,479,472,135,468,238,735,187,375,70,192,333,121,88,589,424,396,413,511,717,655,385,290,133,262,136,672,234,363,79,429,130,176,293,64,576,578,582,135,244,240,89,379,479,154,207,52,673,335,518,410,576,100,161,600,275,562,66,557,645,347,602,275,663,516,67,399,679,538,712,734,213,109,526,102,736,350,675,310,260,356,134,329,573,303,665,243,186,646,211,293,690,535,656,661,504,372,133,473,75,469,443,167,420,567,655,670,57,678,420,620,125,635,619,131,382,71,346,564,342,57,330,601,749,442,513,327,359,448,728,123,321,160,122,321,346,586,559,325,421,434,482,447,391,509,615,137,276,744,710,667,431,411,270,527,132,543,378,492,733,529,284,404,374,747,600,129,205,573,607,386,244,101,726,747,482,307,215,412,180,152,609,82,262,351,424,628,82,83,71,157,75,667,50,589,746,396,662,448,68,475,82,626,659,584,313,674,425,182,256,60,601,284,511,718,653,681,378,672,422,537,95,611,586,658,156,421,404,586,638,288,174,578,439,583,618,579,247,213,57,170,247,98,124,726,461,555,636,391,692,483,699,503,170,450,560,610,573,569,359,695,714,83,189,449,308,160,285,130,196,730,638,626,456,319,186,153,442,280,562,459,723,494,317,250,605,207,678,527,144,184,50,60,594,344,255,477,296,596,103,171,506,669,390,319,599,723,690,457,583,60,445,57,119,616,102,689,313,646,456,425,88,167,403,140,464,170,240,238,144,285,495,496,673,100,129,706,367,616,440,386,448,336,492,481,526,267,428,68,59,202,83,678,605,354,532,205,73,596,684,531,677,522,259,52,199,96,243,229,642,146,189,389,141,262,239,100,443,519,69,185,144,127,131,160,97,558,203,708,398,582,162,288,121,737,496,344,640,494,264,492,316,414,246,59,206,310,503,670,716,736,467,179,497,189,313,638,387,359,687,329,712,116,330,738,141,621,456,68,102,141,303,146,625,337,707,702,702,441,86,496,363,499,373,508,568,140,50,386,249,638,231,622,135,595,240,173,224,622,419,405,294,249,103,115,340,336,412,347,187,345,507,672,115,360,522,218,103,481,80,638,108,408,232,395,192,517,285,135,478,338,607,301,239,429,595,289,501,391,372,460,53,455,383,427,714,411,503,361,356,473,413,52,96,691,327,732,585,710,456,461,447,373,552,292,601,269,54,526,266,621,252,226,241,505,714,355,684,412,419,359,493,606,457,327,107,397,622,266,615,594,670,681,280,418,658,270,612,561,676,174,637,195,421,662,391,88,744,314,151,240,552,599,504,376,626,613,272,691,652,365,533,272,469,592,544,331,554,740,237,676,527,515,305,247,421,598,559,477,537,383,261,384,126,280,397,66,304,559,428,604,335,629,372,586,662,665,623,109,491,619,608,585,701,599,439,724,412,167,445,157,261,70,605,161,136,288,432,279,374,638,395,621,127,174,626,336,453,430,582,324,623,379,187,68,413,600,672,75,609,315,711,623,395,179,248,168,334,501,569,553,260,433,142,397,202,360,196,521,370,432,722,703,488,545,591,159,495,398,217,362,64,54,675,148,712,492,433,81,558,362,679,708,464,386,721,102,208,400,119,373,410,143,569,523,580,744,531,96,726,619,620,734,111,528,437,428,603,521,725,376,466,334,683,56,99,233,341,63,446,637,71,635,182,317,400,639,744,407,68,597,392,645,726,486,403,193,465,726,330,506,64,561,132,387,196,375,361,129,95,126,388,223,539,360,279,241,216,208,462,480,511,240,730,579,415,378,93,186,216,97,312,130,492,631,406,725,399,626,187,55,518,262,419,317,77,259,382,199,384,593,398,682,435,169,738,115,489,658,631,384,234,670,526,299,56,104,357,156,131,448,622,234,697,444,183,310,353,516,68,433,299,720,69,749,214,511,497,276,661,498,119,515,170,419,65,412,326,494,236,674,593,747,689,350,487,67,345,385,631,151,176,247,494,649,465,273,261,285,643,163,280,717,160,179,311,337,295,87,409,343,462,347,77,689,362,96,733,440,734,256,614,103,736,508,68,658,435,745,360,224,548,224,612,336,481,661,337,473,495,675,297,634,415,339,276,593,443,603,170,326,166,323,435,202,264,642,414,693,67,258,551,649,661,730,299,684,186,53,647,639,252,298,102,553,721,497,149,625,88,502,702,152,517,88,53,610,81,92,416,643,84,707,561,173,188,509,635,313,442,300,65,460,575,255,613,589,658,675,274,613,206,225,580,74,693,612,513,549,666,126,508,624,671,357,230,197,498,366,598,537,73,54,197,677,556,149,455,117,221,644,698,331,485,729,407,663,66,518,411,745,714,421,529,471,579,336,479,449,471,388,569,53,230,340,66,727,107,745,691,539,667,281,722,511,577,381,619,748,251,331,732,407,373,638,154,313,439,103,225,452,692,629,341,332,712,556,612,534,334,600,509,157,214,604,565,617,210,248,474,672,456,226,303,622,208,406,739,613,356,547,536,676,140,373,576,320,154,338,167,686,346,732,247,53,483,527,624,577,514,696,218,391,481,655,628,728,283,462,188,280,188,87,194,597,676,205,221,690,395,537,293,748,608,110,62,315,410,360,282,120,714,700,510,153,233,178,723,232,268,678,456,470,563,324,311,155,113,401,393,557,504,528,294,608,367,290,698,571,659,189,523,237,239,223,165,514,104,264,335,57,316,290,408,253,77,602,265,53,511,151,101,541,174,128,203,375,688,475,277,643,224,566,427,257,611,340,736,736,563,660,741,559,706,278,510,543,422,732,485,100,474,478,477,516,422,365,195,264,630,89,204,654,333,618,389,185,62,534,380,713,438,471,738,620,293,408,681,639,219,248,693,89,616,129,250,439,339,654,486,158,131,317,244,245,375,490,312,593,580,742,713,202,300,219,678,670,511,681,254,194,65,63,54,596,679,105,342,422,667,692,68,526,488,397,605,489,83,661,295,346,238,443,441,727,556,451,115,426,120,79,670,673,686,449,327,150,380,600,365,133,678,90,709,737,330,349,265,571,598,372,109,274,636,283,642,189,371,185,594,745,254,684,642,656,666,232,451,208,665,186,685,280,267,171,202,583,530,222,217,121,403,654,324,147,207,186,602,366,645,95,649,396,330,664,343,660,636,245,718,605,559,331,325,414,527,99,155,190,261,177,496,545,256,537,116,483,55,596,569,591,409,574,747,78,197,242,106,562,400,95,632,422,338,606,61,392,639,418,98,651,481,274,732,63,576,644,174,183,340,61,96,526,117,260,446,312,444,617,413,169,722,124,103,423,616,303,124,386,432,376,85,350,572,593,186,747,330,72,600,205,354,249,320,548,329,208,229,481,725,584,526,144,319,587,254,320,371,294,307,147,165,277,248,286,749,517,571,222,141,450,671,97,299,610,320,614,723,286,418,334,174,101,150,741,649,96,245,340,359,76,273,111,76,177,195,484,545,540,445,347,228,597,92,539,642,188,557,367,165,177,259,700,51,496,260,413,605,328,542,224,573,282,676,301,206,197,260,333,370,564,484,226,94,164,638,244,349,84,359,262,460,122,704,100,568,155,443,640,599,747,236,586,257,123,706,134,430,571,314,641,632,683,618,579,598,210,464,635,325,211,624,461,695,442,742,78,113,240,223,338,330,231,672,709,742,636,223,665,346,671,265,552,296,112,204,113,551,67,514,339,339,323,251,701,748,629,319,318,255,461,96,199,407,285,353,746,638,726,341,234,109,662,72,707,543,423,454,261,343,60,706,698,236,581,329,678,695,509,608,479,231,366,441,186,50,583,392,59,362,346,178,684,550,272,511,483,81,656,501,422,417,216,115,92,621,62,528,619,568,193,442,409,308,470,290,491,712,514,722,265,506,318,483,200,115,235,196,742,447,642,72,146,99,597,652,567,369,84,644,555,292,310,355,318,215,710,124,325,237,475,620,678,512,678,342,118,488,725,206,240,694,189,150,186,455,551,110,371,423,132,549,643,275,291,459,618,367,455,237,537,649,402,707,734,449,642,570,227,392,419,297,735,575,123,622,60,333,575,435,435,514,627,692,327,478,362,184,613,724,356,196,472,200,596,601,706,92,693,132,134,445,470,654,350,578,430,542,587,230,283,536,556,666,560,420,642,212,197,159,368,714,131,358,157,360,529,367,515,673,228,551,731,577,58,337,584,291,175,212,179,661,238,183,249,522,740,239,546,707,419,442,557,299,564,572,240,78,542,249,482,167,526,52,365,517,584,551,124,345,460,436,700,79,745,74,655,701,158,241,473,287,427,743,727,636,501,409,219,126,707,613,595,546,632,695,320,96,67,406,210,118,110,201,281,111,105,582,542,591,329,156,147,636,175,339,143,736,464,347,542,112,72,161,337,548,704,665,410,171,654,649,189,78,598,340,441,660,109,260,278,445,255,638,499,120,546,268,168,418,674,715,570,56,271,526,730,228,556,445,746,193,173,727,223,734,476,693,644,396,287,64,180,553,511,87,379,397,166,733,629,109,741,704,110,341,220,452,613,148,644,449,422,605,549,450,567,598,661,657,483,113,303,361,78,524,484,671,702,643,637,206,505,432,50,702,652,384,337,223,398,227,434,626,439,463,476,192,342,143,501,620,216,431,226,535,298,198,647,571,637,437,241,682,691,54,238,394,595,134,386,248,301,141,588,416,374,106,513,114,158,457,313,353,477,495,202,203,205,206,363,644,563,292,77,202,81,184,354,524,111,425,435,554,309,100,375,657,536,170,410,586,159,136,535,455,601,709,281,729,616,528,395,245,676,420,245,343,461,88,392,167,112,225,308,529,645,427,745,102,266,674,152,154,389,500,694,161,203,87,273,483,293,73,127,296,320,644,179,686,704,722,556,124,291,278,724,371,597,339,80,410,397,73,649,710,364,541,370,129,453,334,167,202,365,630,155,144,406,392,51,269,377,467,511,153,50,279,298,641,242,615,700,565,537,669,352,66,373,229,57,445,196,154,678,110,589,650,575,567,62,102,163,683,712,531,403,186,620,622,582,220,405,316,70,323,235,420,335,349,667,230,573,737,668,587,520,712,116,303,429,660,443,78,462,441,384,634,80,78,503,180,700,435,653,309,144,695,146,292,420,275,679,303,513,404,297,273,450,168,91,444,580,573,394,334,616,278,83,335,289,97,657,317,455,464,673,526,652,419,471,224,556,371,237,412,184,574,521,118,449,379,247,411,706,456,597,583,156,135,451,214,536,243,735,133,404,441,444,718,631,97,88,274,736,385,532,736,229,335,444,594,188,204,568,336,736,546,658,635,411,234,280,271,683,588,693,354,518,283,666,133,651,431,383,539,600,267,146,220,436,643,230,613,485,208,228,544,600,269,334,230,731,279,355,716,475,323,355,254,667,418,580,661,348,573,739,345,667,383,320,686,660,554,450,637,611,146,445,674,370,475,682,578,261,716,280,148,99,569,464,749,638,229,526,389,155,211,692,741,701,479,666,384,515,394,456,525,217,530,287,560,710,673,572,696,335,268,675,445,301,629,440,385,104,633,277,544,583,289,124,567,679,737,511,502,629,92,216,239,564,649,566,233,105,496,523,359,513,312,468,187,204,399,625,140,353,563,380,366,653,343,335,749,55,389,278,574,217,353,368,616,98,639,171,513,507,646,712,64,413,312,387,418,742,743,266,625,451,518,508,150,401,700,707,612,586,623,510,597,85,464,139,710,164,644,276,232,568,242,733,271,733,665,414,393,600,457,466,320,385,189,119,197,511,241,694,451,233,55,458,632,625,680,160,168,679,135,159,179,414,554,641,420,595,675,740,313,529,157,580,194,746,430,709,54,258,155,599,720,712,329,119,329,310,431,589,272,148,686,396,357,55,261,597,205,410,616,444,523,690,290,79,614,448,51,454,462,735,560,553,458,488,380,82,191,330,416,651,247,665,546,57,552,55,603,324,270,147,355,650,298,324,339,277,669,239,133,386,64,572,289,683,60,373,476,525,107,587,94,156,531,192,367,718,478,483,388,413,230,430,61,540,508,483,657,638,239,288,163,237,284,340,146,531,446,95,522,575,543,604,526,315,69,203,246,538,139,284,668,605,719,193,732,66,423,311,577,221,294,481,477,732,65,619,649,245,385,684,596,506,613,150,153,741,124,516,671,563,623,739,509,111,320,312,543,76,305,619,691,420,474,170,459,478,561,438,473,573,77,615,305,215,565,451,358,437,301,574,597,586,457,64,223,238,617,137,710,514,237,100,611,734,494,389,646,229,354,325,306,670,612,725,728,256,708,650,597,480,384,230,67,398,66,245,209,342,160,132,214,340,424,153,238,693,539,54,356,662,227,748,355,314,488,440,356,145,421,241,746,561,546,532,225,90,645,715,524,446,259,384,348,537,263,69,648,705,693,339,334,402,332,604,51,224,227,172,512,308,715,453,448,694,503,156,112,135,541,116,652,336,312,564,347,214,301,357,215,713,692,110,51,245,491,455,152,655,267,431,340,321,360,177,389,514,631,185,489,446,437,438,532,145,420,138,522,666,321,391,390,135,509,662,589,322,112,221,628,230,131,743,486,170,237,274,247,560,676,451,436,647,495,623,420,274,241,298,226,93,312,300,572,177,662,375,739,336,600,678,673,612,552,566,212,199,719,246,177,502,228,393,107,602,223,232,249,732,344,630,325,743,508,534,547,130,235,302,271,329,439,673,69,609,455,388,442,354,452,316,491,592,109,452,111,423,735,436,204,129,227,67,290,563,232,338,608,509,658,644,347,583,161,598,728,105,728,105,582,480,138,109,472,731,434,571,566,383,250,439,371,453,236,195,146,403,264,384,317,357,477,633,502,237,691,70,707,366,189,346,54,496,431,733,507,539,201,274,627,609,189,277,148,253,408,283,102,575,507,740,340,232,534,267,457,713,232,305,471,726,424,148,745,118,220,466,319,579,507,438,570,539,169,126,673,117,229,51,310,645,300,223,117,289,317,360,162,693,145,238,453,675,452,50,65,173,157,567,539,666,54,583,613,98,312,229,567,66,740,641,509,333,300,161,670,324,548,271,615,213,638,643,226,515,325,337,324,347,401,715,598,176,127,738,337,385,63,135,591,525,638,183,95,683,171,67,745,111,302,71,255,89,360,80,534,233,481,577,699,556,693,404,653,509,65,746,223,171,139,616,113,298,75,209,733,202,650,496,105,392,556,508,411,592,155,53,473,542,237,521,367,657,353,145,292,247,427,599,173,347,494,264,525,129,273,570,413,685,139,412,427,194,353,229,250,671,739,186,736,109,622,675,663,618,476,327,588,744,198,184,461,74,486,445,400,290,158,61,440,142,629,742,696,563,233,345,89,278,245,426,369,215,604,105,542,347,552,429,444,541,103,428,613,76,468,500,414,361,485,192,595,341,214,198,149,662,270,547,210,252,674,347,372,67,499,708,721,91,258,740,102,489,663,125,421,592,624,61,625,128,662,323,362,580,58,410,356,51,423,80,409,95,478,88,472,635,188,254,739,665,127,62,731,566,151,608,233,483,703,691,366,193,289,218,326,300,140,313,137,196,700,368,318,534,227,593,101,701,682,485,216,125,499,608,62,520,265,317,348,570,160,471,578,320,63,564,515,404,223,171,703,283,306,103,214,217,373,569,149,581,194,170,299,620,233,147,743,304,656,717,72,632,560,98,236,336,749,641,628,698,450,742,110,388,180,586,650,666,536,521,551,330,367,245,431,619,327,320,579,542,594,629,263,302,55,287,476,743,720,231,517,213,727,592,512,596,250,350,399,506,177,141,243,442,420,690,467,151,605,348,440,558,498,619,746,354,111,210,491,736,356,687,138,71,417,519,260,386,307,580,246,614,204,748,169,743,663,563,142,723,265,479,329,309,481,93,78,90,512,521,434,558,350,546,744,386,412,116,412,347,108,737,457,524,246,365,625,634,283,595,181,79,367,127,704,623,264,620,394,687,152,615,651,713,463,232,148,256,50,527,556,509,318,92,166,464,704,392,514,216,177,100,341,565,690,494,650,327,458,479,551,536,434,326,481,119,648,74,154,631,576,207,273,439,193,185,405,263,667,643,73,331,466,354,628,217,147,652,252,415,606,485,453,394,482,633,320,101,653,580,220,222,405,685,521,92,461,448,593,345,484,134,199,497,249,121,649,415,360,563,447,521,169,668,702,519,223,548,297,288,147,78,290,431,612,66,368,490,403,718,657,441,590,118,88,169,68,220,311,196,539,133,660,155,660,183,662,354,596,624,72,659,114,246,696,57,610,676,482,527,674,364,680,254,339,282,75,63,586,744,709,250,567,171,571,720,228,700,376,70,194,612,740,731,483,628,293,681,92,546,559,147,304,707,725,358,643,610,213,178,744,188,354,76,182,517,508,516,582,265,550,661,685,146,310,234,486,332,425,170,262,525,74,391,232,516,235,635,248,287,127,66,476,83,284,490,189,265,97,338,375,746,216,412,152,746,524,58,329,708,649,263,741,260,540,581,482,686,133,540,120,296,705,165,386,299,352,140,741,316,728,427,496,282,310,546,390,461,65,103,93,419,747,234,525,422,621,511,534,78,615,652,690,575,560,139,485,233,140,452,352,445,176,234,510,442,541,611,130,480,623,190,189,540,187,583,380,367,90,254,659,111,553,541,672,573,537,354,745,713,95,581,208,414,354,421,573,605,209,293,110,318,335,182,495,308,599,660,361,93,618,436,710,98,713,649,71,254,250,175,157,446,330,174,602,252,629,575,622,477,474,311,247,148,445,635,91,246,565,265,563,459,628,119,471,329,170,79,117,697,718,198,247,270,272,651,573,664,505,281,565,227,712,531,453,618,348,700,595,381,50,471,556,146,134,486,186,320,380,391,545,319,510,363,164,104,314,73,70,744,493,583,538,158,249,86,348,656,426,466,231,162,295,100,514,655,317,674,339,497,485,706,65,245,563,673,618,102,156,742,691,104,317,75,415,603,679,136,231,471,180,226,693,661,542,680,198,136,202,64,722,746,584,326,632,236,678,454,355,736,203,608,267,255,80,681,358,254,550,746,684,586,573,126,109,354,732,136,482,250,56,228,326,276,662,155,336,654,478,293,405,710,578,378,686,723,205,537,612,696,391,290,77,547,441,128,347,153,463,330,682,504,314,228,179,145,739,438,90,735,306,588,195,663,521,55,532,216,382,522,533,683,382,676,112,432,714,319,76,574,407,507,374,562,82,586,65,707,459,696,61,250,626,536,155,281,665,403,312,659,79,521,449,140,151,315,313,357,382,290,150,191,350,457,648,571,163,413,561,280,379,557,61,686,202,721,179,79,337,216,737,316,119,133,459,355,735,268,667,125,218,116,509,255,652,483,103,560,712,198,70,739,344,581,105,164,467,724,685,382,214,657,161,616,737,488,548,272,116,707,507,548,66,616,374,351,152,251,108,96,354,101,416,373,592,630,576,476,352,216,457,473,549,300,526,211,512,143,221,741,116,113,426,115,523,413,346,486,429,670,161,694,239,357,282,381,605,683,170,497,555,77,538,282,679,58,65,211,578,437,593,443,267,110,657,464,484,113,154,311,635,718,211,310,696,670,195,474,294,672,356,269,743,199,549,406,108,484,63,261,242,113,350,634,502,431,200,655,146,194,103,453,196,63,687,152,185,157,447,289,611,292,430,161,599,562,121,69,707,549,258,456,220,386,179,653,448,259,717,326,87,211,172,648,93,72,175,647,120,460,494,740,457,284,302,240,128,588,266,104,180,595,249,592,422,146,595,680,208,126,535,502,220,194,171,570,83,185,531,556,709,565,441,718,568,732,214,190,601,440,353,462,466,546,175,421,303,559,582,383,522,561,144,535,259,357,495,54,647,656,65,255,642,663,610,684,630,373,561,420,624,718,539,605,281,493,175,78,401,609,645,122,229,172,421,76,191,246,391,329,263,113,740,61,77,55,313,324,257,676,94,630,215,653,452,600,110,407,178,230,197,472,743,714,52,424,311,679,304,607,666,123,368,544,275,201,272,99,701,351,215,191,531,447,304,213,732,419,334,311,615,413,620,330,656,728,69,653,199,162,61,97,474,306,384,161,495,394,50,625,191,337,93,290,317,627,392,630,261,626,611,433,173,741,519,112,84,158,94,102,234,717,535,597,274,505,650,645,467,95,280,451,169,166,624,486,58,625,613,133,333,376,212,458,85,119,694,664,300,158,255,149,540,106,372,589,682,550,279,66,581,300,697,530,535,681,428,297,153,113,420,413,249,367,477,216,210,430,686,673,265,647,325,361,350,259,577,233,675,519,272,471,607,678,380,509,286,66,61,464,455,683,249,107,416,430,93,735,127,442,125,638,465,575,127,91,158,432,569,223,499,211,734,480,682,402,59,395,295,582,256,90,313,615,207,677,181,324,581,453,391,714,487,317,222,242,166,654,569,190,419,661,197,335,549,350,746,181,600,580,237,647,640,108,617,572,231,67,317,476,383,583,630,221,742,657,62,252,376,658,323,404,499,374,274,253,185,365,399,440,408,531,156,86,663,422,158,129,67,337,160,314,335,197,431,373,534,712,688,645,403,472,511,210,199,393,629,109,510,434,371,363,66,364,582,452,112,64,625,340,465,410,290,582,628,630,453,456,673,124,189,254,616,478,498,163,687,405,355,332,278,384,285,181,228,432,561,438,687,80,590,352,134,190,279,62,624,523,427,564,629,459,230,176,555,163,366,485,694,152,288,316,359,116,396,287,342,692,101,379,690,150,369,651,144,622,382,438,413,62,511,607,686,230,126,387,123,381,436,75,644,237,434,588,577,423,657,542,170,513,439,142,739,82,728,295,273,231,465,633,100,207,228,495,393,168,674,602,196,105,515,78,628,716,301,80,440,740,160,551,189,370,304,272,732,597,261,156,286,518,465,56,527,439,84,587,517,375,202,655,127,668,331,341,746,354,368,301,551,254,447,225,515,271,514,162,649,585,462,304,734,368,731,710,645,200,692,327,465,512,432,747,415,559,147,562,60,412,380,327,490,670,556,527,121,185,557,235,484,699,218,418,702,490,133,154,287,58,208,368,670,181,635,636,157,608,536,107,186,505,118,361,498,660,633,153,83,638,648,469,487,629,148,401,243,168,428,556,250,537,208,72,537,253,613,467,163,448,174,480,273,189,310,345,450,117,704,391,340,225,132,659,80,719,64,362,52,502,106,498,605,635,258,76,324,50,698,506,416,207,701,286,647,363,566,188,489,644,287,346,370,534,210,70,578,579,301,524,283,477,412,738,228,404,412,636,740,124,656,193,712,358,259,183,700,627,470,644,277,615,535,606,570,192,186,590,103,51,743,531,541,396,154,349,125,640,296,543,162,335,656,256,53,404,231,623,660,61,147,700,233,90,213,186,657,637,683,291,528,605,230,691,239,626,58,171,301,595,238,528,122,77,491,459,746,341,346,739,343,150,640,99,744,545,552,380,370,534,607,717,142,463,134,309,659,140,745,163,463,563,732,81,613,224,479,501,325,217,717,747,155,511,425,366,652,245,637,467,441,636,491,77,391,135,696,194,424,482,132,683,162,739,323,392,140,65,531,719,262,440,642,85,143,118,168,261,179,629,451,91,259,251,595,262,447,670,74,670,225,259,730,259,645,103,503,644,450,489,234,284,547,237,208,198,601,79,412,355,727,296,686,539,128,206,198,98,502,714,693,736,164,553,185,161,498,728,352,709,193,499,616,177,510,748,464,668,727,693,481,565,287,335,189,218,569,82,237,120,660,579,466,484,559,405,616,294,749,92,233,594,659,664,569,401,564,599,159,434,273,352,240,606,250,53,219,108,59,512,648,740,716,436,182,183,303,214,90,351,174,521,115,431,254,382,736,537,678,151,234,203,148,508,132,547,354,578,94,71,569,147,313,570,565,272,437,276,471,263,344,55,92,659,390,471,721,661,240,539,313,716,363,345,174,420,607,614,185,599,405,465,157,694,189,719,270,629,694,647,436,201,593,725,732,89,497,137,412,96,142,702,172,428,718,168,498,696,361,287,506,536,726,416,362,85,617,429,489,738,95,582,180,724,162,514,724,574,512,584,90,551,113,564,283,458,198,79,563,349,558,355,381,261,253,324,724,117,705,571,601,97,271,111,71,487,351,89,378,566,254,573,99,518,425,565,561,78,313,590,533,109,472,112,356,70,531,444,307,530,158,673,203,240,137,165,419,558,702,370,571,509,374,247,361,183,255,741,469,364,699,653,636,477,442,421,250,55,604,642,264,157,711,217,284,211,343,444,666,469,238,190,197,214,286,453,317,433,112,427,717,235,305,141,492,635,256,687,365,196,204,118,230,657,63,515,713,375,302,600,438,708,105,529,168,382,64,57,132,404,232,603,81,210,391,476,146,189,134,508,215,707,721,305,182,618,150,455,350,184,420,264,327,357,608,629,66,409,372,139,392,96,550,207,599,503,493,319,638,371,677,690,262,57,76,301,472,304,254,625,568,77,655,427,321,679,315,485,334,664,600,733,426,433,717,578,420,142,154,661,416,183,558,196,274,199,300,681,646,103,570,271,509,480,328,366,98,383,320,463,652,394,626,680,466,257,603,555,329,222,166,591,260,252,491,485,619,422,538,269,73,573,116,122,531,424,203,252,747,160,136,539,641,321,575,169,642,421,381,580,189,698,172,746,460,52,619,390,379,509,582,485,727,226,544,318,566,277,392,689,577,120,648,634,397,628,55,365,61,82,324,709,225,82,671,181,123,169,437,687,281,694,155,314,613,733,211,204,572,194,172,502,467,542,383,723,110,544,700,55,623,748,239,155,264,500,382,431,265,309,59,149,318,373,140,695,183,219,303,348,454,554,623,60,195,631,742,307,387,384,239,620,116,146,609,439,536,252,208,449,357,494,237,140,537,84,631,677,492,59,612,605,665,663,379,680,573,727,569,369,676,91,690,446,133,712,79,275,358,57,562,396,668,557,371,234,408,569,423,274,410,677,616,81,524,258,130,718,138,279,551,729,562,668,601,574,196,518,153,99,643,508,333,272,445,689,114,586,327,702,350,602,296,693,145,455,172,694,642,608,274,226,482,426,121,565,716,220,202,338,432,166,624,434,147,632,522,204,646,705,223,626,284,572,218,614,656,627,408,137,728,664,468,434,428,709,391,458,196,361,748,618,454,537,121,70,164,704,120,440,239,521,619,121,605,87,443,516,248,570,242,461,355,571,208,452,706,67,408,137,104,742,54,212,512,745,224,168,577,454,162,282,110,593,367,685,342,395,418,205,393,348,255,577,377,746,259,206,210,114,451,720,231,705,396,369,168,259,285,331,725,669,512,435,137,84,483,53,495,314,462,201,203,696,740,170,131,357,485,118,144,572,726,180,233,517,360,157,685,334,734,541,312,621,600,231,80,139,662,86,211,465,243,491,579,130,295,310,333,665,557,530,73,734,673,81,207,181,649,247,733,536,209,603,735,589,359,704,388,374,519,590,111,213,140,72,668,613,486,57,313,381,187,111,533,566,109,742,60,681,597,168,726,575,634,228,656,315,191,448,242,175,541,609,485,247,227,242,635,233,618,446,173,231,356,73,330,192,694,185,218,735,559,68,669,488,396,625,161,716,598,530,650,203,735,569,561,116,134,158,686,539,261,389,202,401,143,634,129,261,203,563,558,618,437,358,523,112,640,427,86,537,494,633,315,250,725,427,85,348,279,215,104,547,491,176,284,536,743,729,684,64,337,237,86,367,378,400,468,139,87,201,716,125,255,553,193,593,666,363,165,328,381,148,324,55,672,491,437,441,463,679,319,255,332,393,211,399,134,396,674,249,324,137,183,395,193,136,393,219,227,435,540,678,449,410,88,330,224,282,697,300,383,579,169,547,507,370,121,719,78,385,186,216,440,360,641,195,179,219,353,435,264,154,103,279,59,369,591,79,689,563,99,225,302,83,300,739,580,512,136,211,674,316,706,482,169,308,711,342,711,227,706,166,244,514,400,257,340,579,230,619,708,339,195,462,732,687,720,123,411,487,357,256,136,575,474,719,231,418,306,195,668,427,322,555,646,658,669,336,509,443,515,701,216,551,695,715,660,410,589,188,404,676,206,122,398,566,183,144,97,559,189,539,501,727,748,216,662,214,306,623,664,59,474,549,720,56,132,526,108,524,421,302,670,643,67,511,398,568,499,302,135,402,279,540,491,518,172,599,668,236,130,352,204,113,366,374,445,422,635,459,50,742,88,439,690,377,591,80,364,146,547,625,741,571,88,662,145,331,177,68,447,84,643,300,79,300,366,438,178,422,358,62,515,313,252,393,292,666,124,321,311,309,593,412,555,304,675,355,496,516,366,494,736,162,462,446,295,245,432,394,348,167,170,104,354,56,532,705,511,670,187,312,249,637,350,217,665,255,165,285,532,528,667,749,360,633,398,610,342,229,101,97,346,369,695,407,572,78,664,617,186,346,57,693,267,131,526,114,166,473,316,622,455,280,534,342,424,194,345,289,736,136,414,544,456,419,274,360,579,557,331,694,296,410,466,102,147,616,157,114,252,104,670,170,616,492,706,321,263,435,549,681,678,324,137,350,346,474,130,321,745,668,245,364,406,298,659,556,571,103,726,557,737,687,743,308,146,244,530,490,211,155,284,726,563,632,173,131,699,558,414,667,504,128,492,724,677,127,613,468,656,80,653,51,283,191,470,457,552,199,445,474,344,150,241,619,425,552,513,174,264,179,166,625,59,747,549,177,347,146,682,69,183,52,496,610,493,434,79,683,707,290,66,334,491,283,511,261,405,227,648,197,562,570,318,86,473,666,655,248,83,278,519,182,549,339,599,677,379,583,262,245,414,389,518,557,279,51,540,682,254,283,709,436,633,457,203,378,724,80,299,236,61,458,168,491,570,55,278,250,617,85,688,499,460,707,590,195,253,180,315,453,665,348,644,613,108,130,638,433,378,76,111,462,82,205,147,271,465,214,136,230,269,460,286,170,585,406,362,670,676,474,188,424,572,337,426,561,556,118,606,411,321,64,191,121,562,321,502,86,371,115,444,144,627,383,156,592,290,137,431,640,294,374,645,287,159,473,629,420,393,552,260,310,601,541,447,273,530,612,512,674,92,454,226,604,696,422,646,479,315,175,553,400,308,411,140,98,352,421,73,79,707,639,647,384,99,338,150,289,665,425,85,85,328,55,676,563,438,479,120,273,559,541,289,422,202,197,227,358,217,62,535,353,89,684,238,529,137,266,742,440,263,283,277,465,514,700,99,716,441,590,703,696,289,108,666,158,129,731,509,75,88,649,109,593,717,485,366,175,643,103,341,275,388,430,464,317,186,534,437,84,433,377,623,280,254,716,316,744,642,150,505,273,50,364,621,607,199,638,147,213,73,130,364,579,223,98,641,749,238,281,348,130,566,270,168,236,616,745,623,551,156,702,717,279,575,315,443,217,631,477,548,324,606,405,680,173,636,588,276,341,650,388,377,208,246,225,246,106,355,723,247,212,553,469,332,234,653,294,555,248,409,700,561,746,426,284,260,662,241,422,641,663,271,415,601,182,266,337,73,531,202,500,78,182,648,315,507,453,486,233,690,531,602,155,321,287,163,75,247,303,297,518,505,503,520,521,212,418,383,747,731,154,692,449,83,741,156,456,723,531,488,687,504,573,517,387,547,73,307,644,481,673,708,368,645,358,698,702,288,434,127,470,157,158,229,217,438,433,386,394,460,209,322,626,98,525,173,268,477,402,638,98,282,466,320,597,125,339,685,653,741,689,436,275,326,430,314,717,337,612,173,314,285,530,265,531,102,273,182,536,289,94,454,474,408,578,719,210,140,596,592,146,127,651,76,356,556,643,211,74,711,568,585,559,523,263,552,179,494,505,537,117,353,685,656,571,53,212,670,697,235,607,94,347,166,311,99,637,605,385,70,599,375,555,113,188,362,393,592,425,135,185,508,513,649,71,351,602,483,672,444,703,657,327,478,727,324,61,544,345,349,396,408,553,723,526,685,610,75,455,632,362,362,298,90,382,340,623,466,510,172,551,81,172,420,75,309,631,316,677,511,352,336,351,626,453,281,274,426,127,609,652,143,565,664,743,490,515,252,256,104,307,150,545,460,485,641,465,237,551,155,298,632,591,85,671,286,442,569,701,348,154,86,296,293,170,563,566,155,534,526,496,448,259,699,111,562,418,447,306,743,488,313,472,516,266,502,685}

	//fmt.Println(len(latencySLO_500s))
	//fmt.Println(len(latencySLO_1000s))
	//fmt.Println(len(latencySLO_2000s))
	//fmt.Println(len(latencySLO_5000s))
	//fmt.Println("scale=", len(latencySLO_100s))


	clusterCapConfig := InitBox(5000)
	cpuOverSell:= int32(0) //threads
	gpuOverSell:= int32(0) //SM percentage

	start := time.Now()

	cpuTotalConsumRate := float64(0)
	gpuTotalConsumRate := float64(0)
	memoryTotalConsumRate := float64(0)
	MaxResourceConsumRate := float64(0)

	tempGpuCoreQuota := float64(0)
	tempCpuQuota := float64(0)
	for l:=0; l<len(latencySLO_100s); l++ {

		resourcesConfigs :=  testEstimator(float64(latencySLO_100s[l]))
		if len(resourcesConfigs) == 0 {
			continue
		}

		/**
		 * calculate the probabilities
		 */
		if Greater(cpuTotalConsumRate, MaxResourceConsumRate){
			MaxResourceConsumRate = cpuTotalConsumRate
		}
		if Greater(gpuTotalConsumRate, MaxResourceConsumRate){
			MaxResourceConsumRate = gpuTotalConsumRate
		}
		if Greater(memoryTotalConsumRate, MaxResourceConsumRate){
			MaxResourceConsumRate = memoryTotalConsumRate
		}

		maxProb := float64(-1)
		maxNodeProbIndex := -1
		prob := float64(0)
		vce := math.Ceil(MaxResourceConsumRate)*math.E
		for j := 0; j < len(clusterCapConfig); j++ {
			if int(math.Ceil(MaxResourceConsumRate)) <= j && j < int(math.Floor(vce)) {
				prob = math.Log(float64(j+1)/float64(j))
			} else if j == int(math.Floor(vce)) {
				prob = math.Log(vce/float64(j))
			} else {
				prob = 0
			}
			if Greater(prob, maxProb) {
				maxProb = prob
				maxNodeProbIndex = j
			}
		}


		/**
		 * choose resource config based on throughput efficiency
		 */
		maxResourceConfigThrough := float64(-1)
		maxResourceConfigThroughIndex := -1
		tempThroughIntensity := float64(-1)
		for k :=0; k < len(resourcesConfigs); k++ {
			tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(20 + cpuOverSell)
			tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
			tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)

			if Greater(tempThroughIntensity, maxResourceConfigThrough) {
				maxResourceConfigThrough = tempThroughIntensity
				maxResourceConfigThroughIndex = k
			}
		}




		// update GPU memory allocation

		/**
		 * find a node to place function pod
		 */
		if clusterCapConfig[maxNodeProbIndex].CpuThreadsCap + cpuOverSell >= resourcesConfigs[maxResourceConfigThroughIndex].CpuThreads &&
			clusterCapConfig[maxNodeProbIndex].GpuCorePercentCap + gpuOverSell >= resourcesConfigs[maxResourceConfigThroughIndex].GpuCorePercent &&
			GreaterEqual(clusterCapConfig[maxNodeProbIndex].GpuMemoryRateCap, resourcesConfigs[maxResourceConfigThroughIndex].GpuMemoryRate) {

			clusterCapConfig[maxNodeProbIndex].CpuThreadsCap -= resourcesConfigs[maxResourceConfigThroughIndex].CpuThreads
			clusterCapConfig[maxNodeProbIndex].GpuCorePercentCap -= resourcesConfigs[maxResourceConfigThroughIndex].GpuCorePercent
			clusterCapConfig[maxNodeProbIndex].GpuMemoryRateCap -= resourcesConfigs[maxResourceConfigThroughIndex].GpuMemoryRate

			cpuTotalConsumRate += float64(resourcesConfigs[maxResourceConfigThroughIndex].CpuThreads) / float64(20 + cpuOverSell)
			gpuTotalConsumRate += float64(resourcesConfigs[maxResourceConfigThroughIndex].GpuCorePercent) / float64(100 + gpuOverSell)
			memoryTotalConsumRate += resourcesConfigs[maxResourceConfigThroughIndex].GpuMemoryRate
			//fmt.Printf("-----: place %dth Pod config %+v to %dth node\n", maxResourceConfigThroughIndex, resourcesConfigs[maxResourceConfigThroughIndex], maxNodeProbIndex)


		} else { //first fit
			for j := 0; j < len(clusterCapConfig); j++ {
				if clusterCapConfig[j].CpuThreadsCap + cpuOverSell >= resourcesConfigs[maxResourceConfigThroughIndex].CpuThreads &&
					clusterCapConfig[j].GpuCorePercentCap + gpuOverSell >= resourcesConfigs[maxResourceConfigThroughIndex].GpuCorePercent &&
					GreaterEqual(clusterCapConfig[j].GpuMemoryRateCap, resourcesConfigs[maxResourceConfigThroughIndex].GpuMemoryRate) {

					clusterCapConfig[j].CpuThreadsCap -= resourcesConfigs[maxResourceConfigThroughIndex].CpuThreads
					clusterCapConfig[j].GpuCorePercentCap -= resourcesConfigs[maxResourceConfigThroughIndex].GpuCorePercent
					clusterCapConfig[j].GpuMemoryRateCap -= resourcesConfigs[maxResourceConfigThroughIndex].GpuMemoryRate

					cpuTotalConsumRate += float64(resourcesConfigs[maxResourceConfigThroughIndex].CpuThreads) / float64(20 + cpuOverSell)
					gpuTotalConsumRate += float64(resourcesConfigs[maxResourceConfigThroughIndex].GpuCorePercent) / float64(100 + gpuOverSell)
					memoryTotalConsumRate += resourcesConfigs[maxResourceConfigThroughIndex].GpuMemoryRate
					//	fmt.Printf("place %dth Pod %+v to %dth node\n",l,resourcesConfigs[pickConfigIndex], i)
					//fmt.Printf("FT: place %dth Pod config %+v to %dth node\n", maxResourceConfigThroughIndex, resourcesConfigs[maxResourceConfigThroughIndex], maxNodeProbIndex)
					break
				}
			}


		} // check the next <CPU socket and GPU> to place function pod**/
	} // per socket

	fmt.Println("Solve Time: ", time.Since(start))

	boxNum :=0

	for j:=0; j< len(clusterCapConfig); j++ {
		if Equal(clusterCapConfig[j].GpuMemoryRateCap,1.0) &&
			(clusterCapConfig[j].GpuCorePercentCap == 100 ) &&
			clusterCapConfig[j].CpuThreadsCap == 20 {
		} else {
			boxNum++
			fmt.Printf("%f\t%f\t%f \n",
				float64(clusterCapConfig[j].CpuThreadsCap) / float64(20 + cpuOverSell),
				float64(clusterCapConfig[j].GpuCorePercentCap) / float64(100 +gpuOverSell),
				clusterCapConfig[j].GpuMemoryRateCap)
		}

	}
	fmt.Println("Total Box:", boxNum)

	fmt.Println("Optimized Box:")
	fmt.Println(math.Ceil(cpuTotalConsumRate/3))
	fmt.Println(math.Ceil(gpuTotalConsumRate/3))
	fmt.Println(math.Ceil(memoryTotalConsumRate/3))

}

/**
* input: latencySLO
* output: all available resources configurations under latency SLO
*  gpuTypes.FuncPodConfig <CpuThreads,GpuCorePercent,GpuMemoryRate>
 */
type FuncPodConfig struct {
	CpuThreads     int32
	GpuCorePercent int32
	GpuMemoryRate  float64
	ReqPerSecondMax int32
}
const ACC = 0.000001
func Equal(a, b float64) bool {
	return math.Abs(a-b) < ACC
}

func Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > ACC
}

func Less(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > ACC
}

func GreaterEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < ACC
}

func  LessEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < ACC
}
func testEstimator(latencySLO float64) []*FuncPodConfig{
	var availInstConfigs []*FuncPodConfig
	funcName := "resnet-50"
	reqPerSecondMax := int32(0)
	reqPerSecondMin := int32(0)
	for batchSize := int32(16); batchSize > 0; batchSize = batchSize / 2 {
		timeForExec := latencySLO/2
		initCpuThreads := batchSize
		if batchSize == 1 {
			timeForExec = latencySLO //no need to queue batch
			initCpuThreads = 2 // at least allocate 2 CPU threads
		}
		for cpuThreads := initCpuThreads; cpuThreads > 0; cpuThreads = cpuThreads-2 { //cpu threads decreases with 2
			expectTime := execTimeModelOnlyCPU(funcName, batchSize, cpuThreads)
			if LessEqual(expectTime, timeForExec) {
				reqPerSecondMax = int32(1000 / expectTime * float64(batchSize)) //no device idle time - queuing time equals execution time
				if batchSize == 1 {
					reqPerSecondMin = 1
				} else {
					reqPerSecondMin = int32(1000 / (latencySLO - expectTime) * float64(batchSize))
				}
				availInstConfigs = append(availInstConfigs, &FuncPodConfig {
					CpuThreads:     cpuThreads,
					GpuCorePercent: 0,
					GpuMemoryRate: 0,
					ReqPerSecondMax: reqPerSecondMax,
				})
				fmt.Printf("expectTime=%f, batchSize=%d, CpuThreads=%d GpuCorePercent=%d %d-%d\n",
					expectTime,
					batchSize,
					cpuThreads,
					0,
					reqPerSecondMin,
					reqPerSecondMax)

			}

			for gpuCorePercent := 50; gpuCorePercent > 0; gpuCorePercent = gpuCorePercent - 10 { //gpu cores decreases with 10%
				expectTime = execTimeModel(funcName, batchSize, cpuThreads, int32(gpuCorePercent))
				if LessEqual(expectTime, timeForExec) {
					reqPerSecondMax = int32(1000 / expectTime * float64(batchSize)) //no device idle time - queuing time equals execution time
					if batchSize == 1 {
						reqPerSecondMin = 1
					} else {
						reqPerSecondMin = int32(1000 / (latencySLO - expectTime) * float64(batchSize))
					}
					availInstConfigs = append(availInstConfigs, &FuncPodConfig {
						CpuThreads:   cpuThreads,
						GpuCorePercent: int32(gpuCorePercent),
						GpuMemoryRate: 0.1,
						ReqPerSecondMax: reqPerSecondMax,

					})

					fmt.Printf("batchSize=%d, CpuThreads=%d GpuCorePercent=%d expectTime=%f, %d-%d\n",
						batchSize,
						cpuThreads,
						gpuCorePercent,
						expectTime,
						reqPerSecondMin,
						reqPerSecondMax)

				}
			}
		}
	}
	return availInstConfigs
}
func execTimeModel(funcName string, batchSize int32, cpuThread int32, gpuCorePercent int32) float64{
	if funcName == "resnet50" {
		b := float64(batchSize)
		t := float64(cpuThread)
		g := float64(gpuCorePercent) / 100
		return float64(1423)*b/(t+40.09)*math.Pow(g,-0.3105)+17.92
	}
	log.Printf("estimator: could find exection time model for function %s", funcName)
	return 99999
}
func execTimeModelOnlyCPU(funcName string, batchSize int32, cpuThread int32) float64{
	if funcName == "resnet50" {
		b := float64(batchSize)
		t := float64(cpuThread)
		return 34.76*b/(math.Pow(t,0.341)-0.9926)+69.87+150
	}
	log.Printf("estimator: could find exection time model(only CPU) for function %s", funcName)
	return 99999

}

func inferResourceConfigsWithBatch(funcName string, latencySLO float64, batchSize int32, residualReq int32)(instanceConfig []*gpuTypes.FuncPodConfig, err error){
	var availInstConfigs []*gpuTypes.FuncPodConfig
	/**
	 * verify latencySLO is reasonable
	 */

	sloMeet := false
	/**
	 * deal with batch size=1
	 */
	timeForExec := latencySLO/2
	initCpuThreads := batchSize
	if batchSize == 1 {
		timeForExec = latencySLO //no need to queue batch
		initCpuThreads = 2 // at least allocate 2 CPU threads
	}
	//- supportBatchGroup[i]/reqArrivalRate
	//supportCPUthreadsGroup := [...]int32{16,8,4,2,1}
	reqPerSecondMax := int32(0)
	reqPerSecondMin := int32(0)
	for cpuThreads := initCpuThreads; cpuThreads > 0; cpuThreads = cpuThreads-2 { //cpu threads decreases with 2
		expectTime := execTimeModelOnlyCPU(funcName, batchSize, cpuThreads)
		if gpuTypes.LessEqual(expectTime, timeForExec) {
			sloMeet = true
			reqPerSecondMax = int32(1000/expectTime*float64(batchSize)) //no device idle time - queuing time equals execution time
			if batchSize == 1 {
				reqPerSecondMin = 1
			} else {
				reqPerSecondMin = int32(1000/(latencySLO-expectTime)*float64(batchSize))
			}
			if residualReq >= reqPerSecondMin {
				availInstConfigs = append(availInstConfigs, &gpuTypes.FuncPodConfig {
					BatchSize:      batchSize,
					CpuThreads:     cpuThreads,
					GpuCorePercent: 0,
					GpuMemoryRate: -1.0,
					ExecutionTime:  int32(expectTime),
					ReqPerSecondMax: reqPerSecondMax,
					ReqPerSecondMin: reqPerSecondMin,
				})
			}
		}
		for gpuCorePercent := 50; gpuCorePercent > 0; gpuCorePercent = gpuCorePercent - 5 { //gpu cores decreases with 10%
			expectTime = execTimeModel(funcName, batchSize, cpuThreads, int32(gpuCorePercent))
			if gpuTypes.LessEqual(expectTime, timeForExec) {
				sloMeet = true
				reqPerSecondMax = int32(1000/expectTime*float64(batchSize)) //no device idle time - queuing time equals execution time
				if batchSize == 1 {
					reqPerSecondMin = 1
				} else {
					reqPerSecondMin = int32(1000/(latencySLO-expectTime)*float64(batchSize))
				}
				if residualReq >= reqPerSecondMin {
					availInstConfigs = append(availInstConfigs, &gpuTypes.FuncPodConfig {
						BatchSize:      batchSize,
						CpuThreads:     cpuThreads,
						GpuCorePercent: int32(gpuCorePercent),
						GpuMemoryRate: -1.0,
						ExecutionTime:  int32(expectTime),
						ReqPerSecondMax: reqPerSecondMax,
						ReqPerSecondMin: reqPerSecondMin,
					})
				}
			}
		}
	}
	if len(availInstConfigs) > 0 {
		return availInstConfigs,nil
	} else {
		if sloMeet == true {
			err = fmt.Errorf("estimator: residualReq %d is too low to be met with (batchsize=%d)\n",residualReq, batchSize)
		} else {
			err = fmt.Errorf("estimator: latencySLO %f is too low to be met with (batchsize=%d)\n",latencySLO, batchSize)
		}

		return nil, err
	}
}

func CreatePreWarmPod(funcName string, namespace string, latencySLO float64, batchSize int32, clientset *kubernetes.Clientset){

	gpuMemAlloc, err := strconv.Atoi("2048")
	if err == nil {
		log.Printf("scheduler: warm reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
	} else {
		log.Println("scheduler: warm read memory error:", err.Error())
		return
	}


	resourcesConfigs, err := inferResourceConfigsWithBatch(funcName, latencySLO, batchSize, 1)
	if err != nil {
		log.Print(err.Error())
		wrappedErr := fmt.Errorf("scheduler: CreatePrewarmPod failed batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
			batchSize, funcName, latencySLO, 1, 1)
		log.Println(wrappedErr)
		return
	} else {
		/*for _, item := range resourcesConfigs {
			log.Printf("scheduler: warm resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
				funcName, latencySLO, item.ExecutionTime, batchSize, item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
		}*/
	}
	maxResourceQuotaNagDiffIndex := -1
	minResourceQuotaPosDiffIndex := -1
	pickConfigIndex := -1
	maxResourceQuotaNagDiff := float64(-999)
	minResourceQuotaPosDiff := float64(999)
	tempGpuCoreQuota := float64(0)
	tempCpuQuota := float64(0)
	tempDiffQuota := float64(0)

	cpuConsumedRate := float64(0)
	gpuMemConsumedRate := float64(0)
	gpuCoreConsumedRate := float64(0)
	tempThroughIntensity := float64(0)
	tempMinResourceQuotaPosThroughIntensity := float64(0)
	tempMaxResourceQuotaNagThroughIntensity := float64(0)


	cpuConsumedThreadsPerSocket := int(0)
	cpuTotalThreadsPerSocket := int(0)
	cpuOverSell := 0 //CPU threads overSell
	gpuOverSell := 0 //GPU SM percentage
	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
	clusterCapConfig := repository.GetClusterCapConfig()
	for i := 0; i < len(clusterCapConfig.ClusterCapacity); i++ { // per node
		cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
		gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
		gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate

		/** CPU GPU consumed rate **/
		cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
		for j := 0; j < len(cpuCapacity); j++ { // per CPU socket (aka per GPU device (j+1))
			/**
			 * calculate CPU and GPU physical consumption rate
			 */
			cpuConsumedThreadsPerSocket = 0
			cpuTotalThreadsPerSocket = 0
			cpuStatus := cpuCapacity[j].CpuStatus
			for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
				cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
				cpuTotalThreadsPerSocket++
			}
			cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
			cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
			cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j, normalized to 0-1
			gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate //normalized to 0-1
			gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
			log.Println()
			log.Printf("scheduler: warm current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
				i,
				j,
				j+1,
				float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
				clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
				clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)
			/**
			 * allocate resource
			 */
			maxResourceQuotaNagDiffIndex = -1
			minResourceQuotaPosDiffIndex = -1
			pickConfigIndex = -1
			maxResourceQuotaNagDiff = float64(-999)
			minResourceQuotaPosDiff = float64(999)
			if gpuTypes.LessEqual(cpuConsumedRate, gpuCoreConsumedRate) { // cpu is dominantly remained resource
				for k := 0; k < len(resourcesConfigs); k++ {
					if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
						gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
						gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
						log.Printf("scheduler: current node has no enough resources for %dth pod config\n",k)
						continue
					}
					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
					tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
					if gpuTypes.Greater(tempDiffQuota,0) {
						if gpuTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
							minResourceQuotaPosDiff = tempDiffQuota
							minResourceQuotaPosDiffIndex = k
						} else if gpuTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if gpuTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
								minResourceQuotaPosDiffIndex = k
							}
						}
					} else {
						if gpuTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
							maxResourceQuotaNagDiff = tempDiffQuota
							maxResourceQuotaNagDiffIndex = k
						} else if gpuTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
							tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if gpuTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
								maxResourceQuotaNagDiffIndex = k
							}
						}
					}

				}
				//log.Printf("scheduler: warm CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
			} else { // GPU core is dominantly remained resource
				for k := 0; k < len(resourcesConfigs); k++ {
					if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
						resourcesConfigs[k].GpuMemoryRate = 0
					} else {
						resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
					}
					if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
						gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
						gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
						log.Printf("scheduler: current node has no enough resources for %dth pod config\n",k)
						continue
					}
					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
					tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
					if gpuTypes.Greater(tempDiffQuota,0) {
						if gpuTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
							minResourceQuotaPosDiff = tempDiffQuota
							minResourceQuotaPosDiffIndex = k
						} else if gpuTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
							tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if gpuTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
								minResourceQuotaPosDiffIndex = k
							}
						}
					} else {
						if gpuTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
							maxResourceQuotaNagDiff = tempDiffQuota
							maxResourceQuotaNagDiffIndex = k
						} else if gpuTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
							tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
							if gpuTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
								maxResourceQuotaNagDiffIndex = k
							}
						}
					}
				}
				//log.Printf("scheduler: warm GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
			}

			if minResourceQuotaPosDiffIndex == -1 {
				pickConfigIndex = maxResourceQuotaNagDiffIndex
			} else {
				pickConfigIndex = minResourceQuotaPosDiffIndex
			}
			if pickConfigIndex == -1 {
				continue
			}

			// update GPU memory allocation
			cudaDeviceTh := j+1
			if resourcesConfigs[pickConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
				cudaDeviceTh = 0
			}

			if minResourceQuotaPosDiffIndex == -1 {
				log.Printf("scheduler: warm choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, maxResourceQuotaNagDiff=%f\n",
					pickConfigIndex,
					float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
					resourcesConfigs[pickConfigIndex].GpuMemoryRate,
					float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
					maxResourceQuotaNagDiff)
			} else {
				log.Printf("scheduler: warm choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, minResourceQuotaPosDiff=%f\n",
					pickConfigIndex,
					float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
					resourcesConfigs[pickConfigIndex].GpuMemoryRate,
					float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
					minResourceQuotaPosDiff)
			}

			/**
			 * find a node to place function pod
			 */


			var cpuCoreThList []int
			neededCores := resourcesConfigs[pickConfigIndex].CpuThreads >> 1 //hyper-threads
			for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
				if cpuStatus[k].TotalFuncInstance == 0 {
					cpuCoreThList = append(cpuCoreThList, k)
					neededCores--
				}
			}
			for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
				if cpuStatus[k].TotalFuncInstance != 0 {
					if cpuStatus[k].TotalFuncInstance < 3 && gpuTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate,0.8) {
						cpuCoreThList = append(cpuCoreThList, k)
						neededCores--
					}
				}
			}

			if neededCores > 0 {
				log.Printf("scheduler: warm failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
				continue
			}

			log.Printf("scheduler: warm decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
				i,
				j,
				cudaDeviceTh,
				cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
				cpuTotalThreadsPerSocket + cpuOverSell,
				gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
				1 + gpuMemOverSellRate,
				clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
				1 + float64(gpuOverSell)/100)

			resourcesConfigs[pickConfigIndex].NodeGpuCpuAllocation = &gpuTypes.NodeGpuCpuAllocation {
				NodeTh:       i,
				CudaDeviceTh: cudaDeviceTh,
				SocketTh: j,
				CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
			}
			funcPodConfig := resourcesConfigs[pickConfigIndex]
			funcPodConfig.FuncPodType = "p"
			funcPodConfig.InactiveCounter = 0
			repository.AddFuncPodConfig(funcName, funcPodConfig)




			log.Printf("scheduler: warm create prewarm function instance for function %s successfully\n", funcName)

			return
			// check the next <CPU socket and GPU> to place function pod
		} // per socket
	} // per node
}
func ScaleUpFuncCapacity(funcName string, namespace string, latencySLO float64, reqArrivalRate int32, supportBatchGroup []int32, clientset *kubernetes.Clientset) {
	//repository.UpdateFuncIsScalingIn(funcName,true)

	gpuMemAlloc, err := strconv.Atoi("2048")
	if err == nil {
		log.Printf("scheduler: reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
	} else {
		log.Println("scheduler: reading memory error:", err.Error())
		return
	}


	residualReq := reqArrivalRate
	residualFindFlag := false
	batchTryNum := 0

	maxResourceQuotaNagDiffIndex := -1
	minResourceQuotaPosDiffIndex := -1
	pickConfigIndex := -1
	maxResourceQuotaNagDiff := float64(-999)
	minResourceQuotaPosDiff := float64(999)
	tempGpuCoreQuota := float64(0)
	tempCpuQuota := float64(0)
	tempDiffQuota := float64(0)

	cpuConsumedRate := float64(0)
	gpuMemConsumedRate := float64(0)
	gpuCoreConsumedRate := float64(0)
	tempThroughIntensity := float64(0)
	tempMinResourceQuotaPosThroughIntensity := float64(0)
	tempMaxResourceQuotaNagThroughIntensity := float64(0)


	cpuConsumedThreadsPerSocket := int(0)
	cpuTotalThreadsPerSocket := int(0)
	cpuOverSell := 0 //CPU threads overSell
	gpuOverSell := 0 //GPU SM percentage
	gpuMemOverSellRate := float64(0) //GPU memory oversell rate

	for {
		if residualReq > 0 {
			residualFindFlag = false
			if batchTryNum == len(supportBatchGroup) {
				wrappedErr := fmt.Errorf("scheduler: failed to find suitable batchsize for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
					funcName, latencySLO, reqArrivalRate, residualReq)
				log.Println(wrappedErr)
				break
			}
			for batchIndex := 0; batchIndex < len(supportBatchGroup) && residualFindFlag == false; batchIndex++ {
				resourcesConfigs, errInfer := inferResourceConfigsWithBatch(funcName, latencySLO, supportBatchGroup[batchIndex], residualReq)
				if errInfer != nil {
					batchTryNum ++
					log.Print(errInfer.Error())
					wrappedErr := fmt.Errorf("scheduler: batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
						supportBatchGroup[batchIndex], funcName, latencySLO, reqArrivalRate, residualReq)
					log.Println(wrappedErr)
					continue
				} else {
					for _ , item := range resourcesConfigs {
						log.Printf("scheduler: resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
							funcName, latencySLO, item.ExecutionTime, supportBatchGroup[batchIndex], item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
					}
				}




				clusterCapConfig := repository.GetClusterCapConfig()
				for i := 0; i < len(clusterCapConfig.ClusterCapacity) && residualFindFlag == false; i++ { // per node
					cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
					gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
					gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate

					/** CPU GPU consumed rate **/
					cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
					for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
						/**
						 * calculate CPU and GPU physical consumption rate
						 */
						cpuConsumedThreadsPerSocket = 0
						cpuTotalThreadsPerSocket = 0
						cpuStatus := cpuCapacity[j].CpuStatus
						for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
							cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
							cpuTotalThreadsPerSocket++
						}
						cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
						cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
						cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j
						gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate
						gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100)

						log.Println()
						log.Printf("scheduler: current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
							i,
							j,
							j+1,
							float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
							clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
							clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)
						/**
						 * allocate resource
						 */
						maxResourceQuotaNagDiffIndex = -1
						minResourceQuotaPosDiffIndex = -1
						pickConfigIndex = -1
						maxResourceQuotaNagDiff = float64(-999)
						minResourceQuotaPosDiff = float64(999)
						if gpuTypes.LessEqual(cpuConsumedRate, gpuCoreConsumedRate) { // cpu is dominantly remained resource
							for k := 0; k < len(resourcesConfigs); k++ {
								if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
									resourcesConfigs[k].GpuMemoryRate = 0
								} else {
									resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
								}
								if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
									gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
									gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
									log.Printf("scheduler: current node has no enough resources for %dth pod config\n",k)
									continue
								}
								tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
								tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
								tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
								//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
								if gpuTypes.Greater(tempDiffQuota,0) {
									if gpuTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
										minResourceQuotaPosDiff = tempDiffQuota
										minResourceQuotaPosDiffIndex = k
									} else if gpuTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
										tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
											(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
												float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
										if gpuTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
											minResourceQuotaPosDiffIndex = k
										}
									}
								} else {
									if gpuTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
										maxResourceQuotaNagDiff = tempDiffQuota
										maxResourceQuotaNagDiffIndex = k
									} else if gpuTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
										tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
											(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
												float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
										if gpuTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
											maxResourceQuotaNagDiffIndex = k
										}
									}
								}

							}
							//log.Printf("scheduler: warm CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
							//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
						} else { // GPU core is dominantly remained resource
							for k := 0; k < len(resourcesConfigs); k++ {
								if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
									gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
									gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
									log.Printf("scheduler: current node has no enough resources for %dth pod config\n",k)
									continue
								}
								tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
								tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
								tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
								//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
								if gpuTypes.Greater(tempDiffQuota,0) {
									if gpuTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
										minResourceQuotaPosDiff = tempDiffQuota
										minResourceQuotaPosDiffIndex = k
									} else if gpuTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
										tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
											(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
												float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
										if gpuTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
											minResourceQuotaPosDiffIndex = k
										}
									}
								} else {
									if gpuTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
										maxResourceQuotaNagDiff = tempDiffQuota
										maxResourceQuotaNagDiffIndex = k
									} else if gpuTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
										tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
											(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
												float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
										if gpuTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
											maxResourceQuotaNagDiffIndex = k
										}
									}
								}
							}
							//log.Printf("scheduler: warm GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
							//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
						}
						if minResourceQuotaPosDiffIndex == -1 {
							pickConfigIndex = maxResourceQuotaNagDiffIndex
						} else {
							pickConfigIndex = minResourceQuotaPosDiffIndex
						}
						if pickConfigIndex == -1 {
							continue
						}
						// update GPU memory allocation
						cudaDeviceTh := j+1
						if resourcesConfigs[pickConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
							cudaDeviceTh = 0
						}

						if minResourceQuotaPosDiffIndex == -1 {
							log.Printf("scheduler: choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, maxResourceQuotaNagDiff=%f\n",
								pickConfigIndex,
								float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
								resourcesConfigs[pickConfigIndex].GpuMemoryRate,
								float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
								maxResourceQuotaNagDiff)
						} else {
							log.Printf("scheduler: choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, minResourceQuotaPosDiff=%f\n",
								pickConfigIndex,
								float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
								resourcesConfigs[pickConfigIndex].GpuMemoryRate,
								float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
								minResourceQuotaPosDiff)
						}

						/**
						 * find a node to place function pod
						 */

						var cpuCoreThList []int
						neededCores := resourcesConfigs[pickConfigIndex].CpuThreads >> 1 //hyper-threads
						for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
							if cpuStatus[k].TotalFuncInstance == 0 {
								cpuCoreThList = append(cpuCoreThList, k)
								neededCores--
							}
						}
						for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
							if cpuStatus[k].TotalFuncInstance != 0 {
								if cpuStatus[k].TotalFuncInstance < 3 && gpuTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate,0.8) {
									cpuCoreThList = append(cpuCoreThList, k)
									neededCores--
								}
							}
						}

						if neededCores > 0 {
							log.Printf("scheduler: failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
							continue
						}
						log.Printf("scheduler: decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
							i,
							j,
							cudaDeviceTh,
							cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
							cpuTotalThreadsPerSocket + cpuOverSell,
							gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
							1 + gpuMemOverSellRate,
							clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
							1 + float64(gpuOverSell)/100)

						resourcesConfigs[pickConfigIndex].NodeGpuCpuAllocation = &gpuTypes.NodeGpuCpuAllocation {
							NodeTh:       i,
							CudaDeviceTh: cudaDeviceTh,
							SocketTh: j,
							CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
						}




						// todo scaling functions
						funcPodConfig := resourcesConfigs[pickConfigIndex]

						funcPodConfig.FuncPodType = "i"
						funcPodConfig.InactiveCounter = 0
						repository.AddFuncPodConfig(funcName, funcPodConfig)
						log.Printf("scheduler: create function instance for function%s successfully, residualReq=%d-%d=%d \n",
							funcName, residualReq, resourcesConfigs[pickConfigIndex].ReqPerSecondMax, residualReq - resourcesConfigs[pickConfigIndex].ReqPerSecondMax)
						residualReq = residualReq - resourcesConfigs[pickConfigIndex].ReqPerSecondMax

						residualFindFlag = true
						batchTryNum = 0

					} // per socket
				} // per node
			} // per batch size
			if residualFindFlag == false {
				log.Printf("scheduler: failed to find suitable node for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
					funcName, latencySLO, reqArrivalRate, residualReq)
				break
			}

		} else {
			break // residualReq <= 0
		}
	}

	//repository.UpdateFuncIsScalingIn(funcName,false)
	return
}


func testScaleIn(){
	type Config struct {
		Lottary int32
		MinReq int32
		MaxReq int32
		PodName string
	}
	type Func struct {
		MinReqCap int32
		MaxReqCap int32
		LottarySum int32
		ConfigMap map[string]*Config
	}
	funcConfig := map[string]*Config{}
	funcConfig["pod1"]= &Config{
		Lottary: 0,
		MinReq:  60,
		MaxReq:  100,
		PodName: "pod1",
	}
	funcConfig["pod2"]= &Config{
		Lottary: 0,
		MinReq:  150,
		MaxReq:  200,
		PodName: "pod2",
	}
	funcConfig["pod3"]= &Config{
		Lottary: 0,
		MinReq:  180,
		MaxReq:  320,
		PodName: "pod3",
	}
	funcConfig["pod4"]= &Config{
		Lottary: 0,
		MinReq:  400,
		MaxReq:  500,
		PodName: "pod4",
	}
	funcConfig["pod5"]= &Config{
		Lottary: 0,
		MinReq:  10,
		MaxReq:  30,
		PodName: "pod5",
	}
	funcConfig["pod6"]= &Config{
		Lottary: 0,
		MinReq:  15,
		MaxReq:  25,
		PodName: "pod6",
	}
	funcConfig["pod7"]= &Config{
		Lottary: 0,
		MinReq:  10,
		MaxReq:  80,
		PodName: "pod7",
	}
	funcConfig["pod8"]= &Config{
		Lottary: 0,
		MinReq:  1,
		MaxReq:  2,
		PodName: "pod8",
	}
	funcConfig["pod9"]= &Config{
		Lottary: 0,
		MinReq:  5,
		MaxReq:  6,
		PodName: "pod9",
	}
	funcConfig["pod10"]= &Config{
		Lottary: 0,
		MinReq:  5,
		MaxReq:  8,
		PodName: "pod10",
	}
	funcConfig["pod11"]= &Config{
		Lottary: 0,
		MinReq:  12,
		MaxReq:  82,
		PodName: "pod11",
	}
	funcConfig["pod12"]= &Config{
		Lottary: 0,
		MinReq:  10,
		MaxReq:  30,
		PodName: "pod12",
	}
	funcConfig["pod13"]= &Config{
		Lottary: 0,
		MinReq:  15,
		MaxReq:  25,
		PodName: "pod13",
	}
	funcConfig["pod14"]= &Config{
		Lottary: 0,
		MinReq:  10,
		MaxReq:  80,
		PodName: "pod14",
	}
	funcConfig["pod15"]= &Config{
		Lottary: 0,
		MinReq:  1,
		MaxReq:  2,
		PodName: "pod15",
	}
	funcConfig["pod16"]= &Config{
		Lottary: 0,
		MinReq:  5,
		MaxReq:  6,
		PodName: "pod16",
	}
	funcConfig["pod17"]= &Config{
		Lottary: 0,
		MinReq:  5,
		MaxReq:  8,
		PodName: "pod17",
	}
	funcConfig["pod18"]= &Config{
		Lottary: 0,
		MinReq:  12,
		MaxReq:  82,
		PodName: "pod18",
	}
	funcConfig["pod19"]= &Config{
		Lottary: 0,
		MinReq:  10,
		MaxReq:  30,
		PodName: "pod19",
	}
	funcConfig["pod20"]= &Config{
		Lottary: 0,
		MinReq:  15,
		MaxReq:  25,
		PodName: "pod20",
	}
	funcConfig["pod21"]= &Config{
		Lottary: 0,
		MinReq:  10,
		MaxReq:  80,
		PodName: "pod21",
	}
	funcConfig["pod22"]= &Config{
		Lottary: 0,
		MinReq:  1,
		MaxReq:  2,
		PodName: "pod22",
	}
	funcConfig["pod23"]= &Config{
		Lottary: 0,
		MinReq:  5,
		MaxReq:  6,
		PodName: "pod23",
	}
	funcConfig["pod24"]= &Config{
		Lottary: 0,
		MinReq:  5,
		MaxReq:  8,
		PodName: "pod24",
	}
	funcConfig["pod25"]= &Config{
		Lottary: 0,
		MinReq:  12,
		MaxReq:  82,
		PodName: "pod26",
	}

	funcConfig["v"]= &Config{
		Lottary: 0,
		MinReq:  0,
		MaxReq:  0,
		PodName: "v",
	}
	var funcStatuss Func
	funcStatuss.ConfigMap = funcConfig
	funcStatuss.MaxReqCap = 0
	funcStatuss.MinReqCap = 0
	//funcStatuss.LottarySum = 0

	type Combine struct {
		Num int32
		MinCap int32
		MaxCap int32
		DeletePodNameList []string
		RemainPodNameList []string
	}

	var letter []*Config
	maxCap := int32(0)
	minCap := int32(0)
	for kk, vv := range funcStatuss.ConfigMap {
		if kk != "v" {
			letter = append(letter, vv)
		}
		maxCap += vv.MaxReq
		minCap += vv.MinReq
	}
	funcStatuss.MaxReqCap = maxCap
	funcStatuss.MinReqCap = minCap

	var comList []Combine

	n := uint(len(letter))
	var maxCount uint = 1 << n
	fmt.Println("maxCount=",maxCount," n=",n)
	var i uint
	var j uint
	for i = 1; i < maxCount-1; i++ {  //删除1个或删到只剩一个
		num := int32(0)
		minSum := int32(0)
		maxSum := int32(0)
		var deletePodNameList []string
		var remainPodNameList []string
		for j = 0; j < n; j++ {
			if (i & (1 << j)) != 0 { //在做位运算的时候需要注意数据类型为uint
				//
				num++
				minSum+=letter[j].MinReq
				maxSum+=letter[j].MaxReq
				deletePodNameList = append(deletePodNameList, letter[j].PodName)

			} else {
				remainPodNameList = append(remainPodNameList, letter[j].PodName)
			}
		}
		fmt.Printf("deleted num=%d", num)
		com := Combine{
			Num:    num,
			MinCap: minSum,
			MaxCap: maxSum,
			DeletePodNameList: deletePodNameList,
			RemainPodNameList: remainPodNameList,
		}
		comList = append(comList, com)
		fmt.Println()
	}
	//fmt.Println(comList)
	time.Sleep(time.Second*10)
	fmt.Printf("comList: %d\n",len(comList))

	//reqList := []int32{780,740,650,600,458,390,370,355,350,330,310,300,250,230,200,150,100,50,45,5,0}
	reqList := []int32{70, 150}
	for ii:=0;ii<len(reqList);ii++ {
		curReq := reqList[ii]
		foundFlag := false
		candidateNum := int32(n)-1
		for i := candidateNum; i > 0 && foundFlag == false; i-- {
			//fmt.Println(111)
			for j:= 0; j< len(comList); j++ {
				//fmt.Println(222)
				if comList[j].Num == int32(i) {
					tempFuncCapMin := funcStatuss.MinReqCap-comList[j].MinCap
					tempFuncCapMax := funcStatuss.MaxReqCap-comList[j].MaxCap
					if curReq >= tempFuncCapMin && curReq <= tempFuncCapMax {
						for k := 0; k < len(comList[j].DeletePodNameList); k++ {
							funcStatuss.ConfigMap[comList[j].DeletePodNameList[k]].Lottary = 0
						}
						for k := 0; k < len(comList[j].RemainPodNameList); k++ {
							scalingFactor := float32(funcStatuss.ConfigMap[comList[j].RemainPodNameList[k]].MaxReq-funcStatuss.ConfigMap[comList[j].RemainPodNameList[k]].MinReq)/float32(tempFuncCapMax-tempFuncCapMin)
							funcStatuss.ConfigMap[comList[j].RemainPodNameList[k]].Lottary = funcStatuss.ConfigMap[comList[j].RemainPodNameList[k]].MaxReq-int32(scalingFactor*float32(tempFuncCapMax-int32(curReq)))
						}
						funcStatuss.ConfigMap["v"].Lottary = 0
						foundFlag = true
						break
					}
				}
			}
		}
		if foundFlag == false {
			fmt.Println()
			fmt.Printf("%d faild to find one pod\n",curReq)

		}else {
			fmt.Println()
			fmt.Printf("%d \n",curReq)
			for _, v:= range funcStatuss.ConfigMap {
				fmt.Printf("%s, lottary=%d \n",v.PodName, v.Lottary)
			}
		}

	}
}

// NewFunctionScaler create a new scaler with the specified
// ScalingConfig
type ScalingConfig struct {
	// MaxPollCount attempts to query a function before giving up
	MaxPollCount uint

	// FunctionPollInterval delay or interval between polling a function's
	// readiness status
	FunctionPollInterval time.Duration

	// CacheExpiry life-time for a cache entry before considering invalid
	CacheExpiry time.Duration

	// ServiceQuery queries available/ready replicas for function

	// SetScaleRetries is the number of times to try scaling a function before
	// giving up due to errors
	SetScaleRetries uint
}

func NewFunctionScaler(config ScalingConfig) FunctionScaler {
	cache := FunctionCache {
		Cache:  make(map[string]*FunctionMeta),
		Expiry: config.CacheExpiry,
	}
	cacheUpdateFlag := FunctionCacheUpdateFlag{
		CacheUpdate: make(map[string]bool),
		Sync:        sync.RWMutex{},
	}

	return FunctionScaler {
		Cache:  &cache,
		Config: config,
		CacheUpdateFlag: &cacheUpdateFlag,
	}
}
type FunctionCache struct {
	Cache  map[string]*FunctionMeta
	Expiry time.Duration
	Sync   sync.RWMutex
}
type FunctionCacheUpdateFlag struct {
	CacheUpdate  map[string]bool
	Sync   sync.RWMutex
}
type FunctionMeta struct {
	LastRefresh          time.Time
	ServiceQueryResponse ServiceQueryResponse
}
type ServiceQueryResponse struct {
	Replicas          uint64
	MaxReplicas       uint64
	MinReplicas       uint64
	ScalingFactor     uint64
	AvailableReplicas uint64
}

// FunctionScaler scales from zero
type FunctionScaler struct {
	CacheUpdateFlag *FunctionCacheUpdateFlag
	Cache  *FunctionCache
	Config ScalingConfig
}

// FunctionScaleResult holds the result of scaling from zero
type FunctionScaleResult struct {
	Available bool
	Error     error
	Found     bool
	Duration  time.Duration
}


func (fc *FunctionCache) Set(funcNameKey string, serviceQueryResponse ServiceQueryResponse) {
	fc.Sync.Lock()
	defer fc.Sync.Unlock()

	if _, exists := fc.Cache[funcNameKey]; !exists {
		fc.Cache[funcNameKey] = &FunctionMeta{}
	}

	fc.Cache[funcNameKey].LastRefresh = time.Now()
	fc.Cache[funcNameKey].ServiceQueryResponse = serviceQueryResponse
	// entry.LastRefresh = time.Now()
	// entry.ServiceQueryResponse = serviceQueryResponse
}


// Get replica count for functionName
func (fc *FunctionCache) Get(funcNameKey string) (ServiceQueryResponse, bool) {
	replicas := ServiceQueryResponse{
		AvailableReplicas: 0,
	}

	hit := false
	fc.Sync.RLock()
	defer fc.Sync.RUnlock()

	if val, exists := fc.Cache[funcNameKey]; exists {
		replicas = val.ServiceQueryResponse
		hit = !val.Expired(fc.Expiry)
	}

	return replicas, hit
}

// Set replica count for functionName
func (fcuf *FunctionCacheUpdateFlag) SetFlag(funcNameKey string, flag bool) {
	fcuf.Sync.Lock()
	defer fcuf.Sync.Unlock()
	fcuf.CacheUpdate[funcNameKey] = flag
}

// Get replica count for functionName
func (fcuf *FunctionCacheUpdateFlag) GetFlag(funcNameKey string) (bool, bool) {
	fcuf.Sync.RLock()
	defer fcuf.Sync.RUnlock()

	val, exists := fcuf.CacheUpdate[funcNameKey]

	return val, exists
}
func (fm *FunctionMeta) Expired(expiry time.Duration) bool {
	return time.Now().After(fm.LastRefresh.Add(expiry))
}
func  (f *FunctionScaler) Scale(functionName, namespace string, user int) FunctionScaleResult {
	start := time.Now()
	funcNameKey := functionName+"."+namespace
	if cachedResponse, hit := f.Cache.Get(funcNameKey); hit &&
		cachedResponse.AvailableReplicas > 0 {
		return FunctionScaleResult{
			Error:     nil,
			Available: true,
			Found:     true,
			Duration:  time.Since(start),
		}
	}

	if val, exists := f.CacheUpdateFlag.GetFlag(funcNameKey); exists {
		if val == false { // nobody got in
			f.CacheUpdateFlag.SetFlag(funcNameKey,true) //  then get in and lock the door (flag)
			//cacheUpdateFlag[funcNameKey] = true
			queryResponse, err := GetReplicas(functionName, namespace, user)  // external.go implement
			if err != nil {
				f.CacheUpdateFlag.SetFlag(funcNameKey,false)
				return FunctionScaleResult {
					Error:     err,
					Available: false,
					Found:     false,
					Duration:  time.Since(start),
				}
			}

			f.Cache.Set(funcNameKey, queryResponse)

			if queryResponse.AvailableReplicas == 0 {
				minReplicas := uint64(1)
				if queryResponse.MinReplicas > 0 {
					minReplicas = queryResponse.MinReplicas
				}

				scaleResultErr := backoff(func(attempt int) error {
					if queryResponse.Replicas > 0 {
						return nil
					}
					log.Printf("[Scale %d] function=%s 0 => %d requested", attempt, functionName, minReplicas)
					setScaleErr := SetReplicas(functionName, namespace, minReplicas,user)
					if setScaleErr != nil {
						return fmt.Errorf("unable to scale function [%s], err: %s", functionName, setScaleErr)
					}

					return nil

				}, int(f.Config.SetScaleRetries), f.Config.FunctionPollInterval)

				if scaleResultErr != nil {
					f.CacheUpdateFlag.SetFlag(funcNameKey, false)
					return FunctionScaleResult{
						Error:     scaleResultErr,
						Available: false,
						Found:     true,
						Duration:  time.Since(start),
					}
				}

				for i := 0; i < int(f.Config.MaxPollCount); i++ {
					queryResponse, err := GetReplicas(functionName, namespace,user)
					if err == nil {
						f.Cache.Set(funcNameKey, queryResponse)
					}
					totalTime := time.Since(start)

					if err != nil {
						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
						return FunctionScaleResult{
							Error:     err,
							Available: false,
							Found:     true,
							Duration:  totalTime,
						}
					}

					if queryResponse.AvailableReplicas > 0 {
						log.Printf("user %d [Scale] function=%s 0 => %d successful - %fs", user,
							functionName, queryResponse.AvailableReplicas, totalTime.Seconds())
						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
						return FunctionScaleResult {
							Error:     nil,
							Available: true,
							Found:     true,
							Duration:  totalTime,
						}
					}
					time.Sleep(f.Config.FunctionPollInterval)
				}
			} else {
				f.CacheUpdateFlag.SetFlag(funcNameKey, false) // get out and unlock the door
			}
		} else {
			count := 0
			for{
				if count > 10 {
					break
				} else {
					count++
				}
				time.Sleep(time.Millisecond*100)
				log.Printf("user %d waiting the door count=%d \n",user, count)
				if  doorIsClosed, _ := f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false { // wait util someone unlock the door
					log.Printf("user %d door opens count=%d\n",user, count)
					break
				}
			}
		}
	} else {
		f.CacheUpdateFlag.SetFlag(funcNameKey,false)
		//cacheUpdateFlag[funcNameKey] = false
		//  then get in and lock the door (flag)
		if doorIsClosed, _ := f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {
			f.CacheUpdateFlag.SetFlag(funcNameKey,true)
			//cacheUpdateFlag[funcNameKey] = true
			queryResponse, err := GetReplicas(functionName, namespace, user)  // external.go implement
			if err != nil {
				f.CacheUpdateFlag.SetFlag(funcNameKey,false)
				return FunctionScaleResult{
					Error:     err,
					Available: false,
					Found:     false,
					Duration:  time.Since(start),
				}
			}

			f.Cache.Set(funcNameKey, queryResponse)

			if queryResponse.AvailableReplicas == 0 {
				minReplicas := uint64(1)
				if queryResponse.MinReplicas > 0 {
					minReplicas = queryResponse.MinReplicas
				}

				scaleResult := backoff(func(attempt int) error {
					if queryResponse.Replicas > 0 { // expected is not 0
						return nil
					}

					log.Printf("user %d [Scale %d] function=%s 0 => %d requested",user, attempt, functionName, minReplicas)
					setScaleErr := SetReplicas(functionName, namespace, minReplicas, user)
					if setScaleErr != nil {
						return fmt.Errorf("user %d unable to scale function [%s], err: %s", user, functionName, setScaleErr)
					}
					return nil

				}, int(f.Config.SetScaleRetries), f.Config.FunctionPollInterval)

				if scaleResult != nil {
					f.CacheUpdateFlag.SetFlag(funcNameKey,false)
					return FunctionScaleResult{
						Error:     scaleResult,
						Available: false,
						Found:     true,
						Duration:  time.Since(start),
					}
				}

				for i := 0; i < int(f.Config.MaxPollCount); i++ {
					queryResponse, err := GetReplicas(functionName, namespace, user)
					if err == nil {
						f.Cache.Set(funcNameKey, queryResponse)
					}
					totalTime := time.Since(start)

					if err != nil {
						f.CacheUpdateFlag.SetFlag(funcNameKey,false)
						return FunctionScaleResult{
							Error:     err,
							Available: false,
							Found:     true,
							Duration:  totalTime,
						}
					}

					if queryResponse.AvailableReplicas > 0 {
						log.Printf("user %d [Scale] function=%s 0 => %d successful - %fs", user,
							functionName, queryResponse.AvailableReplicas, totalTime.Seconds())
						f.CacheUpdateFlag.SetFlag(funcNameKey,false)
						return FunctionScaleResult{
							Error:     nil,
							Available: true,
							Found:     true,
							Duration:  totalTime,
						}
					}
					time.Sleep(f.Config.FunctionPollInterval)
				}
			} else {
				f.CacheUpdateFlag.SetFlag(funcNameKey,false) // get out and unlock the door
			}

		} else {
			count := 0
			for{
				if count > 10 {
					break
				} else {
					count++
				}
				time.Sleep(time.Millisecond*100)
				log.Printf("user %d waiting the door count=%d \n",user, count)
				if doorIsClosed, _ = f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {// wait util someone unlock the door
					log.Printf("user %d door opens count=%d\n",user, count)
					break
				}
			}
		}
	}

	return FunctionScaleResult{
		Error:     nil,
		Available: true,
		Found:     true,
		Duration:  time.Since(start),
	}
}
type routine func(attempt int) error
func backoff(r routine, attempts int, interval time.Duration) error {
	var err error

	for i := 0; i < attempts; i++ {
		res := r(i)
		if res != nil {
			err = res

			log.Printf("Attempt: %d, had error: %s\n", i, res)
		} else {
			err = nil
			break
		}
		time.Sleep(interval)
	}
	return err
}
func GetReplicas(serviceName, serviceNamespace string, user int) (ServiceQueryResponse, error) {
	start := time.Now()

	var err error
	var emptyServiceQueryResponse ServiceQueryResponse
	log.Printf("user %d %ssystem/function/%s?namespace=%s \n",user, "http://", serviceName, serviceNamespace)
	minReplicas := uint64(1)
	maxReplicas := uint64(20)
	scalingFactor := uint64(20)
	availableReplicas := rand.Intn(3)
	replicas := availableReplicas
	time.Sleep(time.Millisecond*200) // http request
	if 200 == http.StatusOK {
		log.Printf("user %d GetReplicas [%s.%s] =%d took: %fs",user, serviceName, serviceNamespace, availableReplicas, time.Since(start).Seconds())
	} else {
		log.Printf("user %d GetReplicas [%s.%s] =%d took: %fs, code: %d\n",user, serviceName, serviceNamespace, availableReplicas, time.Since(start).Seconds(), 202)
		return emptyServiceQueryResponse, fmt.Errorf("user %d server returned non-200 status code (%d) for function, %s", user,202, serviceName)
	}

	//log.Printf("GetReplicas [%s.%s] took: %fs", serviceName, serviceNamespace, time.Since(start).Seconds())

	return ServiceQueryResponse{
		Replicas:          uint64(replicas),
		MaxReplicas:       maxReplicas,
		MinReplicas:       minReplicas,
		ScalingFactor:     scalingFactor,
		AvailableReplicas: uint64(availableReplicas),
	}, err
}
type ScaleServiceRequest struct {
	ServiceName      string `json:"serviceName"`
	ServiceNamespace string `json:"serviceNamespace"`
	Replicas         uint64 `json:"replicas"`
}
// SetReplicas update the replica count
func SetReplicas(serviceName, serviceNamespace string, count uint64, user int) error {
	var err error

	log.Printf("user %d %ssystem/scale-function/%s?namespace=%s&replicas=%d \n",user, "http://", serviceName, serviceNamespace, count)

	start := time.Now()
	time.Sleep(time.Millisecond*150)
	log.Printf("user %d SetReplicas [%s.%s] took: %fs", user, serviceName, serviceNamespace, time.Since(start).Seconds())

	return err
}

func testDelete(lookupNamespace string, functionName string, clientset *kubernetes.Clientset){
	// build the label of the pods which will be deleted
	labelPod := labels.SelectorFromSet(map[string]string{"faas_function": functionName})
	listPodOptions := metav1.ListOptions{
		LabelSelector: labelPod.String(),
	}
	// This makes sure we don't delete non-labeled deployments
	podList, findPodsErr := clientset.CoreV1().Pods(lookupNamespace).List(listPodOptions)
	if findPodsErr != nil {
		if errors.IsNotFound(findPodsErr) {
			log.Println(http.StatusNotFound)
		} else {
			log.Println(http.StatusInternalServerError)
		}
		log.Println([]byte(findPodsErr.Error()))
		return
	}

	if podList != nil && len(podList.Items) > 0 {
		log.Printf("delete: find existing %d pods for function: %s \n", len(podList.Items), functionName)
		err := deleteFunctionPod(lookupNamespace, listPodOptions, clientset)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println(http.StatusBadRequest)
		log.Println([]byte("delete: can't find existing pods for function: " + functionName))
	}

	srvErr := deleteFunctionService(lookupNamespace, functionName, listPodOptions, clientset)
	if srvErr != nil {
		log.Println(srvErr)
		return
	} else {
		log.Printf("delete: find existing service for function: %s \n", functionName)
	}

	repository.DeleteFunc(functionName)
}
func testReplicas(lookupNamespace string, functionName string, clientset *kubernetes.Clientset){
	funcDeployStatus := repository.GetFunc(functionName)
	if funcDeployStatus == nil {
		fmt.Println(http.StatusInternalServerError)
		fmt.Println("replicas: Unable to lookup function deployment " + functionName)
		return
	}
	resourceLimits := &bootTypes.FunctionResources {
		Memory:     funcDeployStatus.FuncResources.Memory,
		CPU:        funcDeployStatus.FuncResources.CPU,
		GPU:        funcDeployStatus.FuncResources.GPU,
		GPU_Memory: funcDeployStatus.FuncResources.GPU_Memory,
	}

	differ := funcDeployStatus.ExpectedReplicas - funcDeployStatus.AvailReplicas
	if differ > 0 {
		scaleUpFunc(funcDeployStatus, lookupNamespace, resourceLimits, differ, clientset)
	} else if differ < 0 {
		scaleDownFunc(funcDeployStatus, lookupNamespace, -differ, clientset)
	} else {
		fmt.Println("-----------------expectedReplicas=availReplicas do nothing-----------------------")
		// expectedReplicas=availReplicas do nothing
	}
}



func testReaderList(functionNamespace string, clientset *kubernetes.Clientset){
	var functions []ptypes.FunctionStatus // init = nil
	var function *ptypes.FunctionStatus // init = nil

	// search service firstly
	listOpts := metav1.ListOptions{}
	srvs, srvErr := clientset.CoreV1().Services(functionNamespace).List(listOpts)
	if srvErr != nil {
		log.Println(srvErr.Error())
	}
	for _, srvItem := range srvs.Items {
		fmt.Println(srvItem.Name)
		// search pod secondly
		function = k8s.CreateFunctionPodStatus(srvItem.Name) // then read repository to get the pod information
		if function == nil {
			log.Printf("reader: function'pod %s not found \n", srvItem.Name)
		} else {
			//log.Printf("reader: create a func status for function %s from repository, ExpectedReplicas= %d, AvailReplicas= %d \n",srvItem.Name, function.Replicas, function.AvailableReplicas)
			functions = append(functions, *function)
		}
	}
}


func testReader(functionNamespace string, functionName string, clientset *kubernetes.Clientset){
	functionName="sleep"

	// search service firstly
	srvs, srvErr := clientset.CoreV1().Services(functionNamespace).Get(functionName,metav1.GetOptions{})
	if srvErr != nil {
		panic(srvErr.Error())
	}
	fmt.Println("===================================")
	if srvs.Name == functionName { //find this function's service
		// search pod secondly
		function := k8s.CreateFunctionPodStatus(functionName) // then read repository to get the pod information
		// result check
		if function == nil {
			log.Printf("reader: function's pod %s not found \n", functionName)
		} else {
			log.Printf("reader: create a func status for function %s from repository, ExpectedReplicas= %d, AvailReplicas= %d \n",functionName, function.Replicas, function.AvailableReplicas)
		}
	} else {
		log.Printf("reader: function's service %s not found \n", functionName)
	}

}


func testDeploy(clientset *kubernetes.Clientset, factory k8s.FunctionFactory) {
	serviceName := "sleep"
	var constaints []string
	initialReplicas := int32p(2)
	repository.UpdateFuncMinReplicas(serviceName, *initialReplicas)
	namespace := "test"
	resourceRequests := &bootTypes.FunctionResources {
		Memory:     "100Mi",
		CPU:        "2000m",
		GPU:        "0",
		GPU_Memory: "0.4",
	}
	lable := map[string]string{}
	lable["com.openfaas.cpu.bind"]="10,12"
	request := bootTypes.FunctionDeployment {
		Service:                serviceName,
		Image:                  "sleep:latest",
		Network:                "",
		EnvProcess:             "python index.py",
		EnvVars:                nil,
		RegistryAuth:           "",
		Constraints:            constaints,
		Secrets:                []string{},
		Labels:                 &lable,
		Annotations:            nil,
		Limits:                 nil,
		Requests:               resourceRequests,
		ReadOnlyRootFilesystem: false,
		Namespace:              "",
	}

	secrets := k8s.NewSecretsClient(factory.Client)
	existingSecrets, werr := secrets.GetSecrets(namespace, request.Secrets)
	if werr != nil {
		wrappedErr := fmt.Errorf("deploy: unable to fetch secrets: %s \n", werr.Error())
		log.Println(wrappedErr.Error())
		return
	}

	repository.UpdateFuncRequestResources(serviceName, request.Requests)
	repository.UpdateFuncConstrains(serviceName, constaints)
	repository.UpdateFuncExpectedReplicas(serviceName, *initialReplicas)

	var latestPodSepc *corev1.Pod
	for i := *int32p(0); i < *initialReplicas; i++ {
		pod, nodeGpuAlloc, specErr := makePodSpecTest(serviceName, constaints, request, factory, existingSecrets)
		if specErr != nil {
			wrappedErr := fmt.Errorf("deploy: failed make Pod spec for replica = %d: %s \n", i, specErr.Error())
			log.Println(wrappedErr)
			//http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}
		_, err := clientset.CoreV1().Pods(namespace).Create(pod)
		if err != nil {
			wrappedErr := fmt.Errorf("unable create Pod for replicas %d: %s \n", i, err.Error())
			log.Println(wrappedErr)
			//http.Error(w, wrappedErr.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("created-----------------")
		// search pod secondly
		//listOpts.LabelSelector = "faas_function=" + serviceName
		var ii =0
		go func() {
			for {
				ii++
				if ii > 30 {
					break
				}

				pods, podErr := clientset.CoreV1().Pods("test").Get(pod.Name, metav1.GetOptions{})
				if podErr != nil {
					log.Println(podErr.Error())
				}
				if pods.Status.PodIP == "" {
					fmt.Println("sleeping in go func-------------")
				} else {
					fmt.Println(pods.Status.PodIP)
					break
				}
				time.Sleep(time.Second * 1)
			}
		}()

		latestPodSepc = pod //update pointer
		fmt.Println("latestPodSepc-----------------")
		repository.UpdateFuncAvailReplicas(serviceName, i+1)


		funcPodConfig := gpuTypes.FuncPodConfig{
			FuncPodName:          "",
			BatchSize:            0,
			CpuThreads:           0,
			GpuCorePercent:       0,
			ExecutionTime:        0,
			ReqPerSecondMax:         0,
			FuncPodIp:            "",
			NodeGpuCpuAllocation: nodeGpuAlloc,
		}
		repository.AddFuncPodConfig(serviceName, &funcPodConfig)
		log.Printf("deploy: Deployment (pods with replicas = %d) created: %s.%s \n", i+1, serviceName, namespace)
	}
	repository.UpdateFuncSpec(serviceName, latestPodSepc, nil)
	// after that, deploy the service to find the pods with special label
	serviceSpec := makeServiceSpecTest(serviceName, request, factory)
	_, err := clientset.CoreV1().Services(namespace).Create(serviceSpec)
	if err != nil {
		wrappedErr := fmt.Errorf("deploy: failed create Service: %s \n", err.Error())
		log.Println(wrappedErr)
		//http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
	}
	repository.UpdateFuncSpec(serviceName, nil, serviceSpec) //update functionSpec map
	log.Printf("deploy: service created: %s.%s \n", serviceName, namespace)

	return
}

func int32p(i int32) *int32 {
	return &i
}



func makePodSpecTest(serviceName string, constaints []string, request bootTypes.FunctionDeployment, factory k8s.FunctionFactory, secret map[string]*corev1.Secret) (*corev1.Pod, *gpuTypes.NodeGpuCpuAllocation, error) {
	envVars := buildEnvVars(&request)

	labels := map[string]string{
		"faas_function": serviceName,
	}

	// GPU card selection start
	nodeSelector := map[string]string{} // init=map{}
	var nodeGpuAlloc *gpuTypes.NodeGpuCpuAllocation // init=nil
	if constaints != nil && len(constaints) > 0 {
		nodeSelector = createSelector(constaints) // user's defination first
	} else {
		nodeGpuAlloc = scheduler.FindGpuDeployNode(request.Requests,request.Constraints) // only for GPU and GPU_Memory
		if nodeGpuAlloc == nil {
			log.Println("deploy: no node for select")
			return nil, nil, nil
		}
		// build the node selector
		nodeLabelStrList := strings.Split(repository.GetClusterCapConfig().
			ClusterCapacity[nodeGpuAlloc.NodeTh].NodeLabel, "=")
		nodeSelector[nodeLabelStrList[0]] = nodeLabelStrList[1]

		envVars = append(envVars, corev1.EnvVar{
			Name:  "CUDA_VISIBLE_DEVICES",
			Value: strconv.Itoa(nodeGpuAlloc.CudaDeviceTh),
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "GPU_MEM_FRACTION",
			Value: request.Requests.GPU_Memory,
		})
		log.Println("deploy: GPU node selection = ", envVars[0])
	}
	// GPU card selection end

	resources, resourceErr := createResources(request) // only for CPU and memory

	if resourceErr != nil {
		return nil, nil, resourceErr
	}

	var imagePullPolicy corev1.PullPolicy
	switch "Never" {
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		imagePullPolicy = corev1.PullAlways
	}

	annotations := buildAnnotations(request)

	var serviceAccount string

	if request.Annotations != nil {
		annotations := *request.Annotations
		if val, ok := annotations["com.openfaas.serviceaccount"]; ok && len(val) > 0 {
			serviceAccount = val
		}
	}

	probes, err := factory.MakeProbes(request)
	if err != nil {
		return nil, nil, err
	}

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service + "-n"+strconv.Itoa(nodeGpuAlloc.NodeTh)+"-g"+strconv.Itoa(nodeGpuAlloc.CudaDeviceTh)+"pod-" + tools.RandomText(10),
			Annotations: annotations, //prometheus.io.scrape: false
			Labels: labels,
			//labels: com.openfaas.scale.max=15 com.openfaas.scale.min=1 com.openfaas.scale.zero=true
			//faas_function=mnist-test uid=44642818
		},
		Spec: corev1.PodSpec{
			NodeSelector: nodeSelector,
			Containers: []corev1.Container{
				{
					Name:  request.Service + "-con",
					Image: request.Image,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: factory.Config.RuntimeHTTPPort,
							Protocol: corev1.ProtocolTCP},
					},
					Env:             envVars,
					Resources:       *resources,
					ImagePullPolicy: imagePullPolicy,
					LivenessProbe:   probes.Liveness,
					ReadinessProbe:  probes.Readiness,
					SecurityContext: &corev1.SecurityContext{
						ReadOnlyRootFilesystem: &request.ReadOnlyRootFilesystem,
					},
				},
			},
			ServiceAccountName: serviceAccount,
			RestartPolicy:      corev1.RestartPolicyAlways,
			DNSPolicy:          corev1.DNSClusterFirst,
		},
	}

	factory.ConfigureReadOnlyRootFilesystem(request, pod)
	factory.ConfigureContainerUserID(pod)

	if err := factory.ConfigureSecrets(request, pod, secret); err != nil {
		return nil, nil, err
	}

	return pod, nodeGpuAlloc, nil
}


func makeServiceSpecTest(serviceName string, request bootTypes.FunctionDeployment,factory k8s.FunctionFactory) *corev1.Service {

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta {
			Name:        request.Service,
			Annotations: buildAnnotations(request),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string {
				"faas_function": request.Service,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     factory.Config.RuntimeHTTPPort,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: factory.Config.RuntimeHTTPPort,
					},
				},
			},
		},
	}

	return service
}

func createSelector(constraints []string) map[string]string {
	selector := make(map[string]string)

	if len(constraints) > 0 {
		for _, constraint := range constraints {
			parts := strings.Split(constraint, "=")
			if len(parts) == 2 {
				selector[parts[0]] = parts[1]
			}
		}
	}

	return selector
}

func createResources(request bootTypes.FunctionDeployment) (*corev1.ResourceRequirements, error) {
	resources := &corev1.ResourceRequirements{
		Limits:   corev1.ResourceList{},
		Requests: corev1.ResourceList{},
	}

	// Set Memory limits
	if request.Limits != nil && len(request.Limits.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.Memory)
		if err != nil {
			return resources, err
		}
		resources.Limits[corev1.ResourceMemory] = qty
	}

	if request.Requests != nil && len(request.Requests.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.Memory)
		if err != nil {
			return resources, err
		}
		resources.Requests[corev1.ResourceMemory] = qty
	}

	// Set CPU limits
	if request.Limits != nil && len(request.Limits.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.CPU)
		if err != nil {
			return resources, err
		}
		resources.Limits[corev1.ResourceCPU] = qty
	}

	if request.Requests != nil && len(request.Requests.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.CPU)
		if err != nil {
			return resources, err
		}
		resources.Requests[corev1.ResourceCPU] = qty
	}

	return resources, nil
}

func buildAnnotations(request bootTypes.FunctionDeployment) map[string]string {
	var annotations map[string]string
	if request.Annotations != nil {
		annotations = *request.Annotations
	} else {
		annotations = map[string]string{}
	}

	annotations["prometheus.io.scrape"] = "false"
	return annotations
}
func buildEnvVars(request *bootTypes.FunctionDeployment) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if len(request.EnvProcess) > 0 {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k8s.EnvProcessName,
			Value: request.EnvProcess,
		})
	}

	for k, v := range request.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}


	sort.SliceStable(envVars, func(i, j int) bool {
		return strings.Compare(envVars[i].Name, envVars[j].Name) == -1
	})

	return envVars
}

func scaleUpFunc(funcDeployStatus *gpuTypes.FuncDeployStatus, namespace string, resourceLimits *bootTypes.FunctionResources, differ int32, clientset *kubernetes.Clientset){
	nodeSelector := map[string]string{} // init=map{}
	var nodeGpuAlloc *gpuTypes.NodeGpuCpuAllocation // init=nil

	// there is need to decide the new pod's name and deploy node
	for i := int32(0); i < differ; i++ {
		if funcDeployStatus.FuncPlaceConstraints != nil && len(funcDeployStatus.FuncPlaceConstraints) > 0 {
			nodeSelector = createSelector(funcDeployStatus.FuncPlaceConstraints) // user's defination first
		} else {
			nodeGpuAlloc = scheduler.FindGpuDeployNode(resourceLimits,funcDeployStatus.FuncPlaceConstraints) // only for GPU and GPU_Memory
			if nodeGpuAlloc == nil {
				log.Println("replicas: no node for select")
				return
			}
			// build the node selector
			nodeLabelStrList := strings.Split(repository.GetClusterCapConfig().ClusterCapacity[nodeGpuAlloc.NodeTh].NodeLabel, "=")
			nodeSelector[nodeLabelStrList[0]] = nodeLabelStrList[1]
			// build the cuda device env str
			cudaDeviceIndexEnvStr := strconv.Itoa(nodeGpuAlloc.CudaDeviceTh)

			envItemSize := len(funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env)
			for i := 0; i < envItemSize; i++ {
				if funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[i].Name == "CUDA_VISIBLE_DEVICES" {
					funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[i].Value = cudaDeviceIndexEnvStr
					break
				}
			}
			for i := 0; i < envItemSize; i++ {
				if funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[i].Name == "GPU_MEM_FRACTION" {
					funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[i].Value = funcDeployStatus.FuncResources.GPU_Memory
					break
				}
			}
		}
		funcDeployStatus.FuncSpec.Pod.Name = funcDeployStatus.FunctionName + "-pod-" + tools.RandomText(10)
		funcDeployStatus.FuncSpec.Pod.Spec.NodeSelector = nodeSelector
		_, err := clientset.CoreV1().Pods(namespace).Create(funcDeployStatus.FuncSpec.Pod)
		if err != nil {
			wrappedErr := fmt.Errorf("replicas: scaleup function %s 's Pod for differ %d error: %s \n", funcDeployStatus.FunctionName, i+1, err.Error())
			log.Println(wrappedErr)
			return
		}
		repository.UpdateFuncAvailReplicas(funcDeployStatus.FunctionName, funcDeployStatus.AvailReplicas+1)
		funcPodConfig := gpuTypes.FuncPodConfig {
			FuncPodName:          "",
			BatchSize:            0,
			CpuThreads:           0,
			GpuCorePercent:       0,
			ExecutionTime:        0,
			ReqPerSecondMax:         0,
			FuncPodIp:            "",
			NodeGpuCpuAllocation: nodeGpuAlloc,
		}
		repository.AddFuncPodConfig(funcDeployStatus.FunctionName, &funcPodConfig)
		log.Printf("replicas: scaleup function %s 's Pod for differ %d successfully \n", funcDeployStatus.FunctionName, i+1)
	}
}
func scaleDownFunc(funcDeployStatus *gpuTypes.FuncDeployStatus, namespace string, differ int32, clientset *kubernetes.Clientset){

	if funcDeployStatus.AvailReplicas < differ {
		log.Printf("replicas: function %s does not has enough instances %d for differ %d \n", funcDeployStatus.FunctionName, funcDeployStatus.AvailReplicas, differ)
		return
	}
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	for i := int32(0); i < differ; i++ {
		podName := scheduler.FindGpuDeletePod(funcDeployStatus)
		err := clientset.CoreV1().Pods(namespace).Delete(podName, opts)
		if err != nil {
			log.Printf("replicas: function %s deleted pod %s error \n", funcDeployStatus.FunctionName, podName)
		}
		log.Printf("replicas: function %s deleted pod %s successfully \n", funcDeployStatus.FunctionName, podName)

		repository.UpdateFuncAvailReplicas(funcDeployStatus.FunctionName, funcDeployStatus.AvailReplicas-1)
		repository.DeleteFuncPodLocation(funcDeployStatus.FunctionName, podName)
	}

}

func deleteFunctionPod(functionNamespace string, listPodOptions metav1.ListOptions, clientset *kubernetes.Clientset) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	if  deletePodsErr := clientset.CoreV1().Pods(functionNamespace).DeleteCollection(opts, listPodOptions); deletePodsErr != nil {
		if errors.IsNotFound(deletePodsErr) {
			log.Println(http.StatusNotFound)
		} else {
			log.Println(http.StatusInternalServerError)
		}
		log.Println([]byte(deletePodsErr.Error()))
		return fmt.Errorf("delete: delete function %s pods error \n", listPodOptions.LabelSelector)
	}
	//log.Printf("delete: delete function %s pods successfully \n", listPodOptions.LabelSelector)

	return nil
}
func deleteFunctionService(functionNamespace string, functionName string, listPodOptions metav1.ListOptions, clientset *kubernetes.Clientset) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	if svcErr := clientset.CoreV1().Services(functionNamespace).Delete(functionName, opts); svcErr != nil {
		if errors.IsNotFound(svcErr) {
			log.Println(http.StatusNotFound)
		} else {
			log.Println(http.StatusInternalServerError)
		}

		log.Println([]byte(svcErr.Error()))
		return fmt.Errorf("delete: delete function %s service error \n", listPodOptions.LabelSelector)
	}
	log.Printf("delete: delete function %s service successfully \n", listPodOptions.LabelSelector)

	return nil
}


