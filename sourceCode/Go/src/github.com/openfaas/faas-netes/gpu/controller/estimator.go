//@file: estimator.go
//@author: Yanan Yang
//@date: 2020/11/9
//@note:
package controller

import (
	"fmt"
	gTypes "github.com/openfaas/faas-netes/gpu/types"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var maxThroughputEfficiencyMap map[string]float64
var resnet50ProfilesMap map[string]float64
var catdogProfilesMap  map[string]float64
var lstm2365ProfilesMap map[string]float64
var textcnn69ProfilesMap map[string]float64
var mobilenetProfilesMap map[string]float64
var ssdProfilesMap map[string]float64


func InitProfiler() {
	maxThroughputEfficiencyMap = make(map[string]float64)
	resnet50ProfilesMap = make(map[string]float64)
	catdogProfilesMap = make(map[string]float64)
	lstm2365ProfilesMap = make(map[string]float64)
	textcnn69ProfilesMap = make(map[string]float64)
	mobilenetProfilesMap = make(map[string]float64)
	ssdProfilesMap = make(map[string]float64)

	initModel(resnet50ProfilesMap,"resnet-50","./yaml/profiler/resnet-50-profile-results.txt",maxThroughputEfficiencyMap)
	initModel(catdogProfilesMap,"catdog","./yaml/profiler/catdog-profile-results.txt",maxThroughputEfficiencyMap)
	initModel(lstm2365ProfilesMap,"lstm-maxclass-2365","./yaml/profiler/lstm-maxclass-2365-profile-results.txt",maxThroughputEfficiencyMap)
	initModel(textcnn69ProfilesMap,"textcnn-69","./yaml/profiler/textcnn-69-profile-results.txt",maxThroughputEfficiencyMap)
	initModel(mobilenetProfilesMap,"mobilenet","./yaml/profiler/mobilenet-profile-results.txt",maxThroughputEfficiencyMap)
	initModel(ssdProfilesMap,"ssd","./yaml/profiler/ssd-profile-results.txt",maxThroughputEfficiencyMap)
}

func initModel(profilesMap map[string]float64, modelName, filePath string, maxThroughputEfficiencyMap map[string]float64) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}

	maxThroughputEfficiency := -1.0
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	for _, value := range lines {
		params := strings.Split(value, " ")
		cpuCores, err := strconv.ParseFloat(params[1], 64)
		if err != nil {
			cpuCores = 999
		}
		gpuCores, err := strconv.ParseFloat(params[2], 64)
		if err == nil {
			gpuCores = 999
		}
		batchSize, err := strconv.ParseFloat(params[3], 64)
		if err == nil {
			batchSize = 1
		}
		latency, err := strconv.ParseFloat(params[4], 64)
		if err == nil {
			profilesMap[params[1]+"_"+params[2]+"_"+params[3]] = latency
		}
		// calculate the maximum throughput efficiency per model  i.e., throughput/resource
		tempThroughputEfficiency := 1000/latency*batchSize/(cpuCores*64+gpuCores*142)   // 64 GFLOPs per CPU core and 142 GFLOPs per GPU SM core
		if tempThroughputEfficiency > maxThroughputEfficiency {
			maxThroughputEfficiency = tempThroughputEfficiency
		}
	}
	maxThroughputEfficiencyMap[modelName] = maxThroughputEfficiency
	if len(profilesMap) > 0 {
		log.Printf("estimator: read model %s profiling data successfully, number of combinations=%d, maxThroughputEfficieny=%f\n",filePath, len(profilesMap), maxThroughputEfficiency)
	}

}

func inferResourceConfigsWithBatch(funcName string, latencySLO float64, batchSize int32, residualReq int32)(instanceConfig []*gTypes.FuncPodConfig, err error){
	var availInstConfigs []*gTypes.FuncPodConfig
	/**
	 * verify latencySLO is reasonable
	 */
	 /*
	minExecTimeWithGpu := int32(execTimeModel(funcName,1, 2,100))
	minExecTimeOnlyCpu := int32(execTimeModelOnlyCPU(funcName,1, 2))
	if latencySLO < minExecTimeOnlyCpu && latencySLO < minExecTimeWithGpu {
		//log.Printf("estimator: latencySLO %d is too low to be met with minExecTimeOnlyCPU=%d and minExecTimeWithGpu=%d(cpuThread=%d)\n",latencySLO, minExecTimeOnlyCpu, minExecTimeWithGpu, batchSize)
		err = fmt.Errorf("estimator:w latencySLO %d is too low to be met with minExecTimeOnlyCPU=%d and minExecTimeWithGpu=%d(batchSize=%d)\n",latencySLO, minExecTimeOnlyCpu, minExecTimeWithGpu, batchSize)
		return nil, err
	}*/

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
	//supportCPUthreadsGroup := [...]int32{16,14,12,...,2}
	reqPerSecondMax := int32(0)
	reqPerSecondMin := int32(0)
	batchTimeOut := int32(0)
	for cpuThreads := initCpuThreads; cpuThreads > 0; cpuThreads = cpuThreads-2 { //cpu threads decreases with 2
		expectTime := getExecTimeModel(funcName, batchSize, cpuThreads, 0)
		if gTypes.LessEqual(expectTime, timeForExec) {
			sloMeet = true
			reqPerSecondMax = int32(1000/expectTime*float64(batchSize)) //no device idle time - queuing time equals execution time
			if batchSize == 1 {
				reqPerSecondMin = 1
				batchTimeOut = 0
			} else {
				reqPerSecondMin = int32(1000/(latencySLO-expectTime)*float64(batchSize))
				batchTimeOut = int32(latencySLO-expectTime)*1000
				if batchTimeOut < 0 {
					batchTimeOut = 0
				}
			}

			if residualReq >= reqPerSecondMin {
				availInstConfigs = append(availInstConfigs, &gTypes.FuncPodConfig {
					BatchSize:      batchSize,
					CpuThreads:     cpuThreads,
					GpuCorePercent: 0,
					GpuMemoryRate: -1,
					ExecutionTime:  int32(expectTime),
					BatchTimeOut: batchTimeOut,
					ReqPerSecondMax: reqPerSecondMax,
					ReqPerSecondMin: reqPerSecondMin,
				})
				//log.Printf("estimator: function=%s, batch=%d, cpuThread=%d, gpuPercent=%d, batchTimeOut=%d, value=%d\n", funcName,batchSize,cpuThreads,int32(0),batchTimeOut,int32(expectTime))
			}
		}
		for gpuCorePercent := 50; gpuCorePercent > 0; gpuCorePercent = gpuCorePercent - 10 { //gpu cores decreases with 10%
			expectTime = getExecTimeModel(funcName, batchSize, cpuThreads, int32(gpuCorePercent))
			if gTypes.LessEqual(expectTime, timeForExec) {
				sloMeet = true
				reqPerSecondMax = int32(1000/expectTime*float64(batchSize)) //no device idle time - queuing time equals execution time
				if batchSize == 1 {
					reqPerSecondMin = 1
					batchTimeOut = 0
				} else {
					reqPerSecondMin = int32(1000/(latencySLO-expectTime)*float64(batchSize))
					batchTimeOut = int32(latencySLO-expectTime)*1000
					if batchTimeOut < 0 {
						batchTimeOut = 0
					}
				}

				if residualReq >= reqPerSecondMin {
					availInstConfigs = append(availInstConfigs, &gTypes.FuncPodConfig {
						BatchSize:      batchSize,
						CpuThreads:     cpuThreads,
						GpuCorePercent: int32(gpuCorePercent),
						GpuMemoryRate: -1,
						ExecutionTime:  int32(expectTime),
						BatchTimeOut: batchTimeOut,
						ReqPerSecondMax: reqPerSecondMax,
						ReqPerSecondMin: reqPerSecondMin,
					})
					//log.Printf("estimator: function=%s, batch=%d, cpuThread=%d, gpuPercent=%d, batchTimeOut=%d, value=%d\n", funcName,batchSize,cpuThreads,int32(gpuCorePercent),batchTimeOut,int32(expectTime))
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

/**
 * An simulator for estimating execution time with fitting model
 */
//func execTimeModel(funcName string, batchSize int32, cpuThread int32, gpuCorePercent int32) float64{
//	if funcName == "resnet50" {
//		b := float64(batchSize)
//		t := float64(cpuThread)
//		g := float64(gpuCorePercent) / 100
//		return float64(1423)*b/(t+40.09)*math.Pow(g,-0.3105)+17.92
//	}
//	log.Printf("estimator: could find exection time model for function %s", funcName)
//	return 99999
//}
//func execTimeModelOnlyCPU(funcName string, batchSize int32, cpuThread int32) float64{
//	if funcName == "resnet50" {
//		b := float64(batchSize)
//		t := float64(cpuThread)
//		return 34.76*b/(math.Pow(t,0.341)-0.9926)+69.87+150
//	}
//	log.Printf("estimator: could find exection time model(only CPU) for function %s", funcName)
//	return 99999
//
//}



func getExecTimeModel(funcName string, batchSize int32, cpuThread int32, gpuCorePercent int32) float64{
	key := strconv.Itoa(int(cpuThread)) + "_" + strconv.Itoa(int(gpuCorePercent)) + "_" + strconv.Itoa(int(batchSize))
	value := float64(0)
	ok := false
	if funcName == "resnet-50" {
		value, ok = resnet50ProfilesMap[key]
	} else if funcName == "catdog" {
		value, ok = catdogProfilesMap[key]
	} else if funcName == "lstm-maxclass-2365" {
		value, ok = lstm2365ProfilesMap[key]
	} else if funcName == "textcnn-69" {
		value, ok = textcnn69ProfilesMap[key]
	} else if funcName == "mobilenet" {
		value, ok = mobilenetProfilesMap[key]
	} else if funcName == "ssd" {
		value, ok = ssdProfilesMap[key]
	}
	if ok {
		//log.Printf("estimator: function=%s, batch=%d, cpuThread=%d, gpuPercent=%d, value=%f", funcName,batchSize,cpuThread,gpuCorePercent,value)
		return value
	} else {
		//log.Printf("estimator: could not find exection time model for function=%s, batch=%d, cpuThread=%d, gpuPercent=%d", funcName,batchSize,cpuThread,gpuCorePercent)
		return 99999
	}
}

/**
 * query the maximum throughput efficiency of a function
 */
func getMaxThroughputEfficiency(funcName string) float64{
	value, ok := maxThroughputEfficiencyMap[funcName]
	if ok {
		//log.Printf("estimator: function=%s, batch=%d, cpuThread=%d, gpuPercent=%d, value=%f", funcName,batchSize,cpuThread,gpuCorePercent,value)
		return value
	} else {
		//log.Printf("estimator: could not find exection time model for function=%s, batch=%d, cpuThread=%d, gpuPercent=%d", funcName,batchSize,cpuThread,gpuCorePercent)
		return 99999
	}
}
