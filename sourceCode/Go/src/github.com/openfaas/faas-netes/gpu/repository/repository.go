//File  : repository.go
//Author: Yanan Yang
//Date  : 2020/4/7
package repository

import (
	gTypes "github.com/openfaas/faas-netes/gpu/types"
	ptypes "github.com/openfaas/faas-provider/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"log"
	"sort"
	"sync"
)


var clusterCapConfig = gTypes.ClusterCapConfig{}
var funcDeployStatusMap = map[string]*gTypes.FuncDeployStatus{}
var funcProileCache = map[string]*gTypes.FuncProfile{}
//var clusterCapConfigLock *gTypes.ClusterCapConfigLock

// init the cluster configuration /root/yaml/clusterConfig.yml
func Init() {
	parseYAML(&clusterCapConfig)
	//clusterCapConfigLock = &gTypes.ClusterCapConfigLock {
	//	LockerName: "",
	//	LockState: false,
	//}
}

func GetClusterCapConfig() *gTypes.ClusterCapConfig {
	return &clusterCapConfig
}
/*
func LockClusterCapConfig(funcName string) bool {
	if clusterCapConfigLock.LockState == false {
		clusterCapConfigLock.LockState = true
		clusterCapConfigLock.LockerName = funcName
		log.Printf("repository: Lock successfully, locker= %s \n", clusterCapConfigLock.LockerName)
		return true
	} else if clusterCapConfigLock.LockerName == funcName {
		log.Printf("repository: Lock successfully, locker= %s \n", clusterCapConfigLock.LockerName)
		return true
	}
	log.Printf("repository: Lock failed, some one %s is handling the lock \n", clusterCapConfigLock.LockerName)
	return false
}
func UnlockClusterCapConfig(funName string) bool {
	if clusterCapConfigLock.LockState == false {
		log.Printf("repository: Unlock successfully, locker state is false\n")
		return true
	} else if clusterCapConfigLock.LockerName == funName {
		clusterCapConfigLock.LockState = false
		log.Printf("repository: Unlock successfully, unlocker= %s \n", clusterCapConfigLock.LockerName)
		return true
	}
	log.Printf("repository: Unlock failed, lock state= %+v and locker is %s while unlocker is %s\n",
		clusterCapConfigLock.LockState,
		clusterCapConfigLock.LockerName,
		funName)
	return false
}*/

func displayClusterCapConfig() {
	var resource gTypes.ResourceDisplay
	resource = &clusterCapConfig
	resource.ToString()
}

/**
 * storage the function pod spec and service spec for scaling up/down
 */
func GetFunc(funcName string) *gTypes.FuncDeployStatus {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		////log.Printf("repository: function %s exists and get successfully \n", funcName)
		return value
	} else {
		////log.Printf("repository: function %s doesnot exists and get nil \n", funcName)
	}
	return nil
}

func GetFuncProfileCache(funcName string) *gTypes.FuncProfile {
	value, exist := funcProileCache[funcName]
	if exist {
		//log.Printf("repository: funcProfileCache for function=%s hit \n", funcName)
		return value
	}else{
		//log.Printf("repository: funcProfileCache for function=%s miss \n", funcName)
	}
	return nil
}

func UpdateFuncProfileCache(funcProfile *gTypes.FuncProfile) {
	funcProileCache[funcProfile.FunctionName] = funcProfile
	log.Printf("repository: funcProfileCache %+v update successful\n", funcProfile)
}
func RegisterFuncDeploy(funcName string) {
	_ , exist := funcDeployStatusMap[funcName]
	if exist {
		log.Printf("repository: functionStatus=%+v register information exists and init failed\n", funcDeployStatusMap[funcName])
	} else {
		funcDeployStatusMap[funcName] = &gTypes.FuncDeployStatus {
			FunctionName:      funcName,
			FuncSpec:          &gTypes.FuncSpec {
				Pod:     &corev1.Pod{},
				Service: &corev1.Service{},
			},
			ExpectedReplicas: 0,
			AvailReplicas: 0,
			MinReplicas: 0,
			MaxReplicas: 0,
			ScaleToZero: "false",
			FuncQpsPerInstance: 0,
			FuncPodConfigMap: map[string]*gTypes.FuncPodConfig{},
			FuncResources: &ptypes.FunctionResources{},
			FuncPlaceConstraints: []string{},
			FunctionInactiveNum: 0,
			FuncSupportBatchSize: []int32{},
			FuncPodMaxCapacity: 0,
			FuncPodMinCapacity: 0,
			FuncPodTotalLottery: 0,
			FuncRealRps: 0,
			FuncLastRealRps: 0,
			FuncPrewarmPodName: "",
			FunctionScalingLock: &sync.RWMutex{},
			FunctionSortedPodLock: &sync.RWMutex{},
		}
		log.Printf("repository: init functionStatus=%+v deploy information into register map\n",
			funcDeployStatusMap[funcName])
	}
}

func UpdateFuncExpectedReplicas(funcName string, expectedReplicas int32) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.ExpectedReplicas = expectedReplicas
	} else {
		log.Printf("repository: update function %s expectedReplicas= %d failed\n", funcName, expectedReplicas)
		return
	}
	//log.Printf("repository: update function %s expectedReplicas= %d \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].ExpectedReplicas)
}
func UpdateFuncAvailReplicas(funcName string, availReplicas int32) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.AvailReplicas = availReplicas
	} else {
		log.Printf("repository: update function %s availReplicas= %d failed\n", funcName, availReplicas)
		return
	}
	//log.Printf("repository: update function %s availReplicas= %d \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].AvailReplicas)
}
func UpdateFuncMinReplicas(funcName string, minReplicas int32) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.MinReplicas = minReplicas
	} else {
		log.Printf("repository: update function %s minReplicas= %d failed\n", funcName, minReplicas)
		return
	}
	//log.Printf("repository: update function %s minReplicas= %d \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].MinReplicas)
}
func UpdateFuncMaxReplicas(funcName string, maxReplicas int32) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.MaxReplicas = maxReplicas
	} else {
		log.Printf("repository: update function %s maxReplicas= %d failed\n", funcName, maxReplicas)
		return
	}
	//log.Printf("repository: update function %s maxReplicas= %d \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].MaxReplicas)
}
func UpdateFuncScaleToZero(funcName string, scaleToZero string) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.ScaleToZero = scaleToZero
	} else {
		log.Printf("repository: update function %s scaleToZero= %s failed\n", funcName, scaleToZero)
		return
	}
	//log.Printf("repository: update function %s scaleToZero= %s \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].ScaleToZero)
}
/*
func UpdateFuncCpuCoreBind(funcName string, cpuCoreBind string) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncCpuCoreBind = cpuCoreBind
	}else{
		//log.Printf("repository: update function %s //FuncCpuCoreBind= %s failed\n", funcName, funcDeployStatusMap[funcName].FuncCpuCoreBind)
	}

	//log.Printf("repository: update function %s //FuncCpuCoreBind= %s \n", funcName, funcDeployStatusMap[funcName].FuncCpuCoreBind)
}*/

func UpdateFuncQpsPerInstance(funcName string, funcQpsPerInstance float64) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncQpsPerInstance = funcQpsPerInstance
	}else{
		log.Printf("repository: update function %s funcSLO= %f failed\n", funcName, funcQpsPerInstance)
		return
	}
	//log.Printf("repository: update function %s funcSLO= %f \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].FuncQpsPerInstance)
}

func UpdateFuncRealRps(funcName string, funcRealRps int32) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncRealRps = funcRealRps
	} else {
		log.Printf("repository: update function %s funcRealRps= %d failed\n", funcName, funcRealRps)
		return
	}
	log.Printf("repository: update function %s funcRealRps= %d \n",
		funcName,
		funcDeployStatusMap[funcName].FuncRealRps)
}

func UpdateFuncLastRealRps(funcName string, funcLastRealRps int32) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncLastRealRps = funcLastRealRps
	} else {
		log.Printf("repository: update function %s funcLastRealRps= %d failed\n", funcName, funcLastRealRps)
		return
	}
	//log.Printf("repository: update function %s funcLastRealRps= %d \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].FuncLastRealRps)
}

func UpdateFuncInactiveNum(funcName string, funcInactiveNum int32) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FunctionInactiveNum = funcInactiveNum
	} else {
		log.Printf("repository: update function %s functionInactiveNum= %d failed \n", funcName, funcInactiveNum)
		return
	}
	//log.Printf("repository: update function %s functionInactiveNum= %d \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].FunctionInactiveNum)
}

func UpdateFuncRequestResources(funcName string, resources *ptypes.FunctionResources) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncResources = resources
	} else {
		log.Printf("repository: update function %s request resource = %+v failed\n", funcName, resources)
		return
	}
	//log.Printf("repository: update function %s request resource = %+v \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].FuncResources)
}
func UpdateFuncSupportBatchSize(funcName string, maxBatchSize int32) {
	var batchSizeGroup []int32
	for i:= maxBatchSize; i > 0; i = i >> 1 {
		batchSizeGroup = append(batchSizeGroup, i)
	}
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncSupportBatchSize = batchSizeGroup
	}else{
		log.Printf("repository: update function %s max batchSizeGroup = %+v failed\n", funcName, batchSizeGroup)
		return
	}
	//log.Printf("repository: update function %s max batchSizeGroup = %+v \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].FuncSupportBatchSize)
}
func UpdateFuncConstrains(funcName string, constrains []string) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncPlaceConstraints = constrains
	}else{
		log.Printf("repository: update function %s constraints= %+v failed\n", funcName, constrains)
		return
	}
	//log.Printf("repository: update function %s constraints= %+v \n",
	//	funcName,
	//	funcDeployStatusMap[funcName].FuncPlaceConstraints)
}

//func UpdateFuncPrewarmPodName(funcName string, funcPrewarmPodName string) {
//	value, exist := funcDeployStatusMap[funcName]
//	if exist {
//		value.FuncPrewarmPodName = funcPrewarmPodName
//	}else{
//		log.Printf("repository: update function %s funcPrewarmPodName= %+v failed\n", funcName, funcPrewarmPodName)
//	}
//	log.Printf("repository: update function %s funcPrewarmPodName= %+v \n", funcName, funcPrewarmPodName)
//}

func GetFuncScalingLockState(funcName string) *sync.RWMutex {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		return value.FunctionScalingLock
	}else{
		log.Printf("repository: get function %s ScalingLockState failed since function does not exist\n", funcName)
	}
	return nil
}
func GetFunctionSortedPodLockState(funcName string) *sync.RWMutex {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		return value.FunctionSortedPodLock
	}else{
		log.Printf("repository: get function %s ScalingLockState failed since function does not exist\n", funcName)
	}
	return nil
}

func UpdateFuncSpec(funcName string, podSpec *corev1.Pod, srvSpec *corev1.Service){
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		if podSpec != nil {
			value.FuncSpec.Pod = podSpec
		}
		if srvSpec != nil {
			value.FuncSpec.Service = srvSpec
		}
	} else {
		log.Printf("repository: update function %s in repository podSec= %+v failed\n", funcName, podSpec)
		log.Printf("repository: update function %s in repository srvSec= %+v failed\n", funcName, srvSpec)
		return
	}
	//if podSpec != nil {
	//	log.Printf("repository: update function %s in repository podSec= %+v \n",
	//		funcName,
	//		funcDeployStatusMap[funcName].FuncSpec.Pod)
	//}
	//if srvSpec != nil {
	//	log.Printf("repository: update function %s in repository srvSec= %+v \n",
	//		funcName,
	//		funcDeployStatusMap[funcName].FuncSpec.Service)
	//}
}

/**
 * only be executed when a new function is deployed
 */
func AddVirtualFuncPodConfig(funcName string){
	funcPodConfig := &gTypes.FuncPodConfig {
		FuncPodType:          "v",
		FuncPodName:          "v",
		ReqPerSecondLottery:  0,
		ReqPerSecondMin: 0,
		ReqPerSecondMax: 0,
		FuncPodIp:            "127.0.0.1",
		NodeGpuCpuAllocation: nil,
	}
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncPodConfigMap[funcPodConfig.FuncPodName] = funcPodConfig
	} else {
		log.Printf("repository: add virtual pod for function %s failed\n", funcName)
		return
	}
	log.Printf("repository: add virtual pod for function %s FuncPodConfigMap= %+v \n",
		funcName,
		funcDeployStatusMap[funcName].FuncPodConfigMap)
	//funcPodConfig.ToString()
}

func AddFuncPodConfig(funcName string, funcPodConfig *gTypes.FuncPodConfig) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncPodConfigMap[funcPodConfig.FuncPodName] = funcPodConfig
		if funcPodConfig.FuncPodType == "i" {
			value.FuncPodMaxCapacity += funcPodConfig.ReqPerSecondMax
			value.FuncPodMinCapacity += funcPodConfig.ReqPerSecondMin
		} else if funcPodConfig.FuncPodType == "p" {
			value.FuncPrewarmPodName = funcPodConfig.FuncPodName
		}
		log.Printf("repository: add function pod configuration =%+v totalPod=%d\n",
			funcPodConfig,
			len(value.FuncPodConfigMap))
		//log.Printf("repository: add function pod configuration =%+v (NodeGpuCpuAlloc=%+v) of function %s and func preWarmPod=%s, maxCap=%d, minCap=%d \n",
		//	funcDeployStatusMap[funcName].FuncPodConfigMap[funcPodConfig.FuncPodName],
		//	funcDeployStatusMap[funcName].FuncPodConfigMap[funcPodConfig.FuncPodName].NodeGpuCpuAllocation,
		//	funcName,
		//	funcDeployStatusMap[funcName].FuncPrewarmPodName,
		//	funcDeployStatusMap[funcName].FuncPodMaxCapacity,
		//	funcDeployStatusMap[funcName].FuncPodMinCapacity)

		/**
		 * update the function instance (pod) in GPU card and the memory left
		 */
		nodeCap := clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh]
		gpuDevice := nodeCap.GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh] //node0:{gpu<-1>,gpu<0>,gpu<1>,gpu<2>}
		cpuSocket := nodeCap.CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh]     //node0:{socket0:<0,1,2>, socket1<3,4,5>, socket2<6,7,8>}
		// add CPU allocation
		cpuCoreThList := funcPodConfig.NodeGpuCpuAllocation.CpuCoreThList
		for _, coreTh := range cpuCoreThList { // k th cpu cores of selected socket in deployed node
			cpuSocket.CpuStatus[coreTh].TotalFuncInstance++
			funcProfile := GetFuncProfileCache(funcName)
			if funcProfile == nil {
				cpuSocket.CpuStatus[coreTh].TotalCpuUsageRate += 0.5 //todo: update the cpu usage rate with profiling cache
			} else {
				cpuSocket.CpuStatus[coreTh].TotalCpuUsageRate += funcProileCache[funcName].AvgCpuCoreUsage
			}
			//log.Printf("repository: add CPU allocation <DataStruct: node=%dth, socket=%dth, cpuCoreAlloc=%+v th> <OS indexNum: node=%s, socket=%d, core=%d, totalFuncInstance=%d, coreUsageRate=%f> for pod=%s, function=%s\n",
			//	funcPodConfig.NodeGpuCpuAllocation.NodeTh,
			//	funcPodConfig.NodeGpuCpuAllocation.SocketTh,
			//	funcPodConfig.NodeGpuCpuAllocation.CpuCoreThList,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].NodeLabel,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuSocketIndex,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuStatus[coreTh].CpuCoreIndex,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuStatus[coreTh].TotalFuncInstance,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuStatus[coreTh].TotalCpuUsageRate,
			//	funcPodConfig.FuncPodName,
			//	funcName)
		}
		// add GPU allocation
		gpuDevice.TotalFuncInstance++ //here must be calculated even when gpuDeviceIndex == -1
		gpuDevice.TotalGpuMemUsageRate += funcPodConfig.GpuMemoryRate
		gpuDevice.TotalGpuCoreUsageRate += float64(funcPodConfig.GpuCorePercent) / 100

		//log.Printf("repository: add GPU allocation <DataStruct indexNum: node=%dth, cudaDevice=%dth> <OS indexNum: node=%s, cudaDevice=%d, totalFuncInstances=%d, totalGpuMemUsageRate=%f, totalGpuCoreUsageRage=%f> for pod=%s, function=%s\n",
		//	funcPodConfig.NodeGpuCpuAllocation.NodeTh,
		//	funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh,
		//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].NodeLabel,
		//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].CudaDeviceIndex,
		//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].TotalFuncInstance,
		//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].TotalGpuMemUsageRate,
		//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].TotalGpuCoreUsageRate,
		//	funcPodConfig.FuncPodName,
		//	funcName)
		//displayClusterCapConfig()
	} else {
		log.Println("repository: add new pod config failed because function does not exist in repository")
	}
}

/**
 * if a ipod is converted into ppod, then remove it from the funcPodCap
 * if a ppod is converted into ipod, then add it to the funcPodCap
 */
func UpdateFuncPodType(funcName string, funcPodName string, newPodType string){
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		funcPodConfig, has := value.FuncPodConfigMap[funcPodName]
		if has {
			funcPodConfig.FuncPodType = newPodType
			if funcPodConfig.FuncPodType == "i" { // conv p to i, need to update cap
				value.FuncPodMaxCapacity += funcPodConfig.ReqPerSecondMax
				value.FuncPodMinCapacity += funcPodConfig.ReqPerSecondMin
			} else if funcPodConfig.FuncPodType == "p" { // conv i to p, need to update cap
				value.FuncPodMaxCapacity -= funcPodConfig.ReqPerSecondMax
				value.FuncPodMinCapacity -= funcPodConfig.ReqPerSecondMin
			}
			log.Printf("repository: update podType to %s of pod=%s in function=%s and func maxCap=%d, minCap=%d \n",
				funcDeployStatusMap[funcName].FuncPodConfigMap[funcPodName].FuncPodType,
				funcPodName,
				funcName,
				funcDeployStatusMap[funcName].FuncPodMaxCapacity,
				funcDeployStatusMap[funcName].FuncPodMinCapacity)
		} else {
			log.Printf("repository: update podType to %s failed of pod=%s in function %s since pod does not exist in repository\n",
				newPodType, funcPodName, funcName)
		}
	} else {
		log.Printf("repository: update podType to %s failed of pod=%s since function %s does not exist in repository\n",
			newPodType, funcPodName, funcName)
	}
}

func UpdateFuncLastChangedPodCombine(funcName string, changedPodCombine *gTypes.ChangedPodCombine){
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		value.FuncLastChangedPodCombine = changedPodCombine
		//log.Printf("repository: update lastChangedPodCombine %+v of function %s successfully\n",
		//	changedPodCombine,
		//	funcName)

	} else {
		log.Printf("repository: update lastChangedPodCombine %+v failed since function %s does not exist in repository\n",
			changedPodCombine,
			funcName)
	}
}
/**
 * inactiveCounter of a pod is set to 0 when it is created
 * when it is going to be deleted, the inactiveCounter++, and this value is reset to 0 when it's lottery allocated > 0
 */
func UpdateFuncPodInactiveCounter(funcName string, funcPodName string, podInactiveCounter int32){
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		funcPodConfig, has := value.FuncPodConfigMap[funcPodName]
		if has {
			funcPodConfig.InactiveCounter = podInactiveCounter
			//if podInactiveCounter != 0 {
			//	log.Printf("repository: update podInactiveCounter to %d of pod=%s in function=%s\n",
			//		funcDeployStatusMap[funcName].FuncPodConfigMap[funcPodName].InactiveCounter,
			//		funcPodName,
			//		funcName)
			//}
		} else {
			log.Printf("repository: update podInactiveCounter to %d of pod=%s in function=%s failed since pod does not exist in repository\n",
				podInactiveCounter,
				funcPodName,
				funcName)
		}
	} else {
		log.Printf("repository: update podInactiveCounter to %d of pod=%s in function=%s failed since function does not exist in repository\n",
			podInactiveCounter,
			funcPodName,
			funcName)
	}
}
/**
 * only i and v pod could be updated
 */
func UpdateFuncPodLottery(funcName string, funcPodName string, podLottery int32){
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		podConf, has := value.FuncPodConfigMap[funcPodName]
		if has {
			podConf.ReqPerSecondLottery = podLottery
			if podLottery > 0 {
				UpdateFuncPodInactiveCounter(funcName, funcPodName,0) //reset deleteInactiveNum
			}
			//log.Printf("repository: update function %s pod %s lottery=%d successed\n",
			//	funcName,
			//	funcPodName,
			//	funcDeployStatusMap[funcName].FuncPodConfigMap[funcPodName].ReqPerSecondLottery)
		} else {
			log.Printf("repository: update function %s pod %s lottery=%d failed, pod doesnot exist \n",funcName, funcPodName, podLottery)
		}
	} else {
		log.Printf("repository: update function %s pod %s lottery=%d failed, function doesnot exist \n",funcName, funcPodName, podLottery)
	}
}

func UpdateFuncPodsTotalLotteryNoLog(funcName string) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		funcPodTotalLottery := int32(0)
		var keys []string
		for _ , v := range value.FuncPodConfigMap {
			if v.FuncPodType != "p" {
				funcPodTotalLottery+= v.ReqPerSecondLottery
				keys = append(keys, v.FuncPodName)
			}
		}
		sort.Strings(keys)
		lock := GetFunctionSortedPodLockState(funcName)
		lock.Lock()
		value.FuncPodTotalLottery = funcPodTotalLottery
		value.FuncSortedPodNameList = keys
		lock.Unlock()

	} else {
		log.Printf("repository: update function %s pods total lottery failed, function doesnot exist \n",funcName)
	}
}
func UpdateFuncPodsTotalLottery(funcName string) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		funcPodTotalLottery := int32(0)
		var keys []string
		for _ , v := range value.FuncPodConfigMap {
			if v.FuncPodType != "p" {
				funcPodTotalLottery+= v.ReqPerSecondLottery
				keys = append(keys, v.FuncPodName)
			}
		}
		sort.Strings(keys)

		value.FuncPodTotalLottery = funcPodTotalLottery
		value.FuncSortedPodNameList = keys


		for _, v := range value.FuncPodConfigMap {

			//if v.FuncPodType == "v" {
			log.Printf("repository: updated function %s type %s pod=%s lottery=%d, (%.1f%%) funcPodTotalLottery=%d, funcMaxCap=%d, funcMinCap=%d, SortedListLen=%d\n",
				funcName,
				v.FuncPodType,
				v.FuncPodName,
				funcDeployStatusMap[funcName].FuncPodConfigMap[v.FuncPodName].ReqPerSecondLottery,
				float64(funcDeployStatusMap[funcName].FuncPodConfigMap[v.FuncPodName].ReqPerSecondLottery)/float64(funcDeployStatusMap[funcName].FuncPodTotalLottery)*100,
				funcDeployStatusMap[funcName].FuncPodTotalLottery,
				funcDeployStatusMap[funcName].FuncPodMaxCapacity,
				funcDeployStatusMap[funcName].FuncPodMinCapacity,
				len(keys))
			//}

		}





	} else {
		log.Printf("repository: update function %s pods total lottery failed, function doesnot exist \n",funcName)
	}
}

/**
 * only ipod will be deleted
 */
func DeleteFuncPodLocation(funcName string, podName string){
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		funcPodConfig, has := value.FuncPodConfigMap[podName]
		if has {
			nodeCap := clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh]
			gpuDevice := nodeCap.GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh]
			cpuSocket := nodeCap.CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh]
			/**
			 * delete CPU and CPU allocation
			 */
			cpuCoreThList := funcPodConfig.NodeGpuCpuAllocation.CpuCoreThList
			for _ , coreTh := range cpuCoreThList {  // k th cpu cores of selected socket in deployed node
				cpuSocket.CpuStatus[coreTh].TotalFuncInstance--
				funcProfile := GetFuncProfileCache(funcName)
				if funcProfile == nil {
					cpuSocket.CpuStatus[coreTh].TotalCpuUsageRate-= 0.5 //todo: update the cpu usage rate with profiling cache
				} else {
					cpuSocket.CpuStatus[coreTh].TotalCpuUsageRate-= funcProileCache[funcName].AvgCpuCoreUsage
				}
				//log.Printf("repository: delete CPU allocation <DataStruct: node=%dth, socket=%dth, cpuCoreAlloc=%+v th> <OS indexNum: node=%s, socket=%d, core=%d, totalFuncInstance=%d, coreUsageRate=%f> for pod=%s, function=%s\n",
				//	funcPodConfig.NodeGpuCpuAllocation.NodeTh,
				//	funcPodConfig.NodeGpuCpuAllocation.SocketTh,
				//	funcPodConfig.NodeGpuCpuAllocation.CpuCoreThList,
				//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].NodeLabel,
				//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuSocketIndex,
				//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuStatus[coreTh].CpuCoreIndex,
				//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuStatus[coreTh].TotalFuncInstance,
				//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].CpuCapacity[funcPodConfig.NodeGpuCpuAllocation.SocketTh].CpuStatus[coreTh].TotalCpuUsageRate,
				//	funcPodConfig.FuncPodName,
				//	funcName)
			}
			//delete GPU allocation
			gpuDevice.TotalFuncInstance-- //here must be calculated even when gpuDevice=-1
			gpuDevice.TotalGpuMemUsageRate -= funcPodConfig.GpuMemoryRate
			gpuDevice.TotalGpuCoreUsageRate -= float64(funcPodConfig.GpuCorePercent)/100

			delete(value.FuncPodConfigMap, funcPodConfig.FuncPodName)
			if funcPodConfig.FuncPodType == "i" {
				value.FuncPodMaxCapacity-= funcPodConfig.ReqPerSecondMax
				value.FuncPodMinCapacity-= funcPodConfig.ReqPerSecondMin
			}
			//log.Printf("repository: delete GPU allocation <DataStruct indexNum: node=%dth, cudaDevice=%dth> <OS indexNum: node=%s, cudaDevice=%d, totalFuncInstances=%d, totalGpuMemUsageRate=%f, totalGpuCoreUsageRage=%f> for pod=%s, function=%s\n",
			//	funcPodConfig.NodeGpuCpuAllocation.NodeTh,
			//	funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].NodeLabel,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].CudaDeviceIndex,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].TotalFuncInstance,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].TotalGpuMemUsageRate,
			//	clusterCapConfig.ClusterCapacity[funcPodConfig.NodeGpuCpuAllocation.NodeTh].GpuCapacity[funcPodConfig.NodeGpuCpuAllocation.CudaDeviceTh].TotalGpuCoreUsageRate,
			//	funcPodConfig.FuncPodName,
			//	funcName)
			//displayClusterCapConfig()
		}
	} else {
		log.Printf("repository: pod %s of function %s does not exist in repository and delete failed", podName, funcName)
	}
}

/**
 * vpod is no need to release resources
 */
func DeleteFunc(funcName string) {
	value, exist := funcDeployStatusMap[funcName]
	if exist {
		for _, funcPodConfig := range value.FuncPodConfigMap {
			if funcPodConfig.FuncPodType != "v" {
				DeleteFuncPodLocation(funcName, funcPodConfig.FuncPodName)
			}
		}
		delete(funcDeployStatusMap, funcName)
		log.Printf("repository: function %s exists and is deleted successfully\n", funcName)
	} else {
		log.Printf("repository: function %s does not exist in repository and can not be deleted \n", funcName)
	}
}

func parseYAML(clusterCapConfig *gTypes.ClusterCapConfig) {
	data, _ := ioutil.ReadFile("./yaml/clusterCapConfig-dev.yml")
	//fmt.Println(string(data))
	//把yaml形式的字符串解析成struct类型
	err := yaml.Unmarshal(data, clusterCapConfig)
	if err != nil {
		log.Println("repository: reading yaml error " + err.Error())
	} else {
		log.Println("repository: read clusterCapConfig.yml successfully")
	}
}