////File  : scheduler.go
////Author: Yanan Yang
////Date  : 2020/4/7
package controller
//
//import (
//	"fmt"
//	"github.com/openfaas/faas-netes/gpu/repository"
//	gTypes "github.com/openfaas/faas-netes/gpu/types"
//	"k8s.io/client-go/kubernetes"
//	"log"
//	"strconv"
//	"sync"
//)
//var lock sync.Mutex
//
//const CpuUsageRateThreshold = 0.8
//const CpuFuncInstanceThreshold = 3
//
//
//func CreatePreWarmPod(funcName string, namespace string, latencySLO float64, batchSize int32, clientset *kubernetes.Clientset){
//	funcStatus := repository.GetFunc(funcName)
//	if funcStatus == nil {
//		log.Printf("scheduler: warm function %s is nil in repository, error to read GPU memory", funcName)
//		return
//	}
//	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
//	if err == nil {
//		//log.Printf("scheduler: warm reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
//	} else {
//		log.Println("scheduler: warm read memory error:", err.Error())
//		return
//	}
//
//	resourcesConfigs, err := inferResourceConfigsWithBatch(funcName, latencySLO, batchSize,1)
//	if err != nil {
//		/*log.Print(err.Error())
//		wrappedErr := fmt.Errorf("scheduler: CreatePrewarmPod failed batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
//			batchSize, funcName, latencySLO, 1, 1)
//		log.Println(wrappedErr)*/
//		return
//	} else {
//		/*for _, item := range resourcesConfigs {
//			log.Printf("scheduler: warm resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
//				funcName, latencySLO, item.ExecutionTime, batchSize, item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
//		}*/
//	}
//
//	cpuConsumedThreadsPerSocket := int(0)
//	cpuTotalThreadsPerSocket := int(0)
//	cpuOverSell := 0 //CPU threads overSell
//	gpuOverSell := 0 //GPU SM percentage
//	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
//	residualFindFlag := false
//
//	gpuMemConsumedRate := float64(0)
//	gpuCoreConsumedRate := float64(0)
//	cpuConsumedRate := float64(0)
//	gpuConsumedRate := float64(0)
//
//	tempGpuCoreQuotaRate := float64(0)
//	tempGpuMemQuotaRate := float64(0)
//	tempCpuQuotaRate := float64(0)
//	tempGpuQuotaRate := float64(0)
//
//	tempCRE := float64(0)
//	maxCRE := float64(-1)
//	maxCREConfigIndex := -1
//
//	lock.Lock()
//	//log.Println("scheduler: warm locked---------------------------")
//	clusterCapConfig := repository.GetClusterCapConfig()
//	for i := 0; i < len(clusterCapConfig.ClusterCapacity) && residualFindFlag == false; i++ { // per node
//		cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
//		gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
//		gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate
//
//		/** CPU GPU consumed rate **/
//		cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
//		for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
//			/**
//			 * calculate CPU and GPU core, memory physical resource consumption rate for each slot
//			 */
//			cpuConsumedThreadsPerSocket = 0
//			cpuTotalThreadsPerSocket = 0
//			cpuStatus := cpuCapacity[j].CpuStatus
//			for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
//				cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
//				cpuTotalThreadsPerSocket++
//			}
//			cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
//			cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
//			cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j, normalized to 0-1
//			gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
//			gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate / (1.0 + gpuMemOverSellRate)//normalized to 0-1
//			/**
//			 * comparison of GPU core and memory
//			 */
//			if gTypes.GreaterEqual(gpuCoreConsumedRate, gpuMemConsumedRate) {
//				gpuConsumedRate = gpuCoreConsumedRate
//			} else {
//				gpuConsumedRate = gpuMemConsumedRate
//			}
//
//			//log.Println()
//			//			//log.Printf("scheduler: warm current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
//			//			//	i,
//			//			//	j,
//			//			//	j+1,
//			//			//	float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
//			//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
//			//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)
//
//			maxCRE = float64(-1)
//			maxCREConfigIndex = -1
//			if gTypes.LessEqual(cpuConsumedRate, gpuConsumedRate) { // cpu is dominantly remained resource
//				for k := 0; k < len(resourcesConfigs); k++ {
//					if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
//						resourcesConfigs[k].GpuMemoryRate = 0
//					} else {
//						resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
//					}
//					/**
//					 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
//					 */
//					tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
//					tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
//					tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)
//
//					/**
//					 * check if the slot has enough resource
//					 */
//					if cpuConsumedRate + tempCpuQuotaRate > 1.01 ||
//						gpuCoreConsumedRate + tempGpuCoreQuotaRate > 1.01 ||
//						gpuMemConsumedRate + tempGpuMemQuotaRate > 1.01  {
//						//log.Printf("scheduler: warm current node has no enough resources for %dth pod config, skip to next pod config\n",k)
//						continue
//					} else {
//						//log.Printf("scheduler: warm current node has enough resources for %dth pod config, skip to next pod config\n",k)
//					}
//					/**
//					 * comparison of GPU core and memory
//					 */
//					if gTypes.GreaterEqual(tempGpuCoreQuotaRate, tempGpuMemQuotaRate) {
//						tempGpuQuotaRate = tempGpuCoreQuotaRate
//					} else {
//						tempGpuQuotaRate = tempGpuMemQuotaRate
//					}
//					/**
//					 * match the dominated resource
//					 */
//					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
//					if gTypes.GreaterEqual(tempCpuQuotaRate, tempGpuQuotaRate) {
//						tempCRE = float64(resourcesConfigs[k].ReqPerSecondMax) / (float64(resourcesConfigs[k].CpuThreads)*64 + float64(resourcesConfigs[k].GpuCorePercent)*142)
//						if gTypes.Greater(tempCRE, maxCRE) {
//							maxCRE = tempCRE
//							maxCREConfigIndex = k
//						}
//					}
//				}
//
//				//log.Printf("scheduler: warm CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
//				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
//			} else { // GPU core is dominantly remained resource
//				for k := 0; k < len(resourcesConfigs); k++ {
//					if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
//						resourcesConfigs[k].GpuMemoryRate = 0
//					} else {
//						resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
//					}
//					/**
//					 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
//					 */
//					tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
//					tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
//					tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)
//
//					/**
//					 * check if the slot has enough resource
//					 */
//					if cpuConsumedRate + tempCpuQuotaRate > 1.01 ||
//						gpuCoreConsumedRate + tempGpuCoreQuotaRate > 1.01 ||
//						gpuMemConsumedRate + tempGpuMemQuotaRate > 1.01  {
//						//log.Printf("scheduler: warm current node has no enough resources for %dth pod config, skip to next pod config\n",k)
//						continue
//					} else {
//						//log.Printf("scheduler: warm current node has enough resources for %dth pod config, skip to next pod config\n",k)
//					}
//					/**
//					 * comparison of GPU core and memory
//					 */
//					if gTypes.GreaterEqual(tempGpuCoreQuotaRate, tempGpuMemQuotaRate) {
//						tempGpuQuotaRate = tempGpuCoreQuotaRate
//					} else {
//						tempGpuQuotaRate = tempGpuMemQuotaRate
//					}
//					/**
//					 * match the dominated resource
//					 */
//					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
//					if gTypes.LessEqual(tempCpuQuotaRate, tempGpuQuotaRate) {
//						tempCRE = float64(resourcesConfigs[k].ReqPerSecondMax) / (float64(resourcesConfigs[k].CpuThreads)*64 + float64(resourcesConfigs[k].GpuCorePercent)*142)
//						if gTypes.Greater(tempCRE, maxCRE) {
//							maxCRE = tempCRE
//							maxCREConfigIndex = k
//						}
//					}
//				}
//				//log.Printf("scheduler: warm GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
//				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
//			}
//
//			if maxCREConfigIndex == -1 {
//				continue // no pod resource config can be placed into this socket
//			}
//
//			// update GPU memory allocation
//			cudaDeviceTh := j+1
//			if resourcesConfigs[maxCREConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
//				cudaDeviceTh = 0
//			}
//
//			//if minResourceQuotaPosDiffIndex == -1 {
//			//	log.Printf("scheduler: warm choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, maxResourceQuotaNagDiff=%f\n",
//			//		pickConfigIndex,
//			//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
//			//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
//			//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
//			//		maxResourceQuotaNagDiff)
//			//} else {
//			//	log.Printf("scheduler: warm choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, minResourceQuotaPosDiff=%f\n",
//			//		pickConfigIndex,
//			//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
//			//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
//			//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
//			//		minResourceQuotaPosDiff)
//			//}
//
//			/**
//			 * allocate CPU threads
//			 */
//			var cpuCoreThList []int
//			neededCores := resourcesConfigs[maxCREConfigIndex].CpuThreads >> 1 //hyper-threads
//			for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//				if cpuStatus[k].TotalFuncInstance == 0 {
//					cpuCoreThList = append(cpuCoreThList, k)
//					neededCores--
//				}
//			}
//			for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//				if cpuStatus[k].TotalFuncInstance != 0 {
//					if cpuStatus[k].TotalFuncInstance < CpuFuncInstanceThreshold && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate, CpuUsageRateThreshold) {
//						cpuCoreThList = append(cpuCoreThList, k)
//						neededCores--
//					}
//				}
//			}
//
//			if neededCores > 0 {
//				//log.Printf("scheduler: warm failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
//				continue
//			}
//
//			//log.Printf("scheduler: warm decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
//			//	i,
//			//	j,
//			//	cudaDeviceTh,
//			//	cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
//			//	cpuTotalThreadsPerSocket + cpuOverSell,
//			//	gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
//			//	1 + gpuMemOverSellRate,
//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
//			//	1 + float64(gpuOverSell)/100)
//
//			resourcesConfigs[maxCREConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation {
//				NodeTh:       i,
//				CudaDeviceTh: cudaDeviceTh,
//				SocketTh: j,
//				CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
//			}
//
//			createErr := createFuncInstance(funcName, namespace, resourcesConfigs[maxCREConfigIndex],"p", clientset)
//			if createErr != nil {
//				log.Println("scheduler: warm create prewarm function instance failed", createErr.Error())
//			} else {
//				//log.Printf("scheduler: warm create prewarm function instance for function %s successfully\n", funcName)
//				residualFindFlag = true
//			}
//		} // per socket
//	} // per node
//	lock.Unlock()
//	if residualFindFlag == false {
//		supportBatchGroup := []int32{1}
//		residualReq := createPreWarmPodWithoutFilter(funcName, namespace, latencySLO,1, supportBatchGroup, clientset)
//		if residualReq > 0 {
//			log.Printf("scheduler-withoutFilter: warm create prewarm function instance for function %s failed\n", funcName)
//		}
//	}
//	//log.Println("scheduler: warm unlocked---------------------------")
//	return
//}
//
//
//
//func createPreWarmPodWithoutFilter(funcName string, namespace string, latencySLO float64, residualReq int32, supportBatchGroup []int32, clientset *kubernetes.Clientset) int32 {
//
//	funcStatus := repository.GetFunc(funcName)
//	if funcStatus == nil {
//		log.Printf("scheduler-withoutFilter: function %s is nil in repository, error to scale up", funcName)
//		return residualReq
//	}
//	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
//	if err == nil {
//		//log.Printf("scheduler: reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
//	} else {
//		log.Println("scheduler-withoutFilter: reading memory error:", err.Error())
//		return residualReq
//	}
//
//	cpuConsumedThreadsPerSocket := 0
//	cpuTotalThreadsPerSocket := 0
//	cpuOverSell := 0 //CPU threads overSell
//	gpuOverSell := 0 //GPU SM percentage
//	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
//
//	batchTryNum := 0
//	residualFindFlag := false
//
//	cpuConsumedRate := float64(0)
//	gpuMemConsumedRate := float64(0)
//	gpuCoreConsumedRate := float64(0)
//
//	tempGpuCoreQuotaRate := float64(0)
//	tempGpuMemQuotaRate := float64(0)
//	tempCpuQuotaRate := float64(0)
//
//	tempCRE := float64(0)
//	maxCRE := float64(-1)
//	maxCREConfigIndex := -1
//
//	if residualReq > 0 {
//		for batchIndex := 0; batchIndex < len(supportBatchGroup) && residualFindFlag == false; batchIndex++ {
//			resourcesConfigs, errInfer := inferResourceConfigsWithBatch(funcName, latencySLO, supportBatchGroup[batchIndex], residualReq)
//			if errInfer != nil {
//				batchTryNum ++
//				if batchTryNum >= len(supportBatchGroup) {
//					wrappedErr := fmt.Errorf("scheduler-withoutFilter: failed to find suitable batchsize for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
//						funcName, latencySLO, residualReq, residualReq)
//					log.Println(wrappedErr)
//					return residualReq // stop explore batch
//				} else {
//					continue // try next smaller batch size
//				}
//			} else {
//				//try the allocation, do the following code
//			}
//
//			lock.Lock()
//			clusterCapConfig := repository.GetClusterCapConfig()
//			for i := 0; i < len(clusterCapConfig.ClusterCapacity) && residualFindFlag == false; i++ { // per node
//				cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
//				gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
//				gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate
//
//				/** CPU GPU consumed rate **/
//				cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
//				for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
//					/**
//					 * calculate CPU and GPU core, memory physical resource consumption rate for each slot
//					 */
//					cpuConsumedThreadsPerSocket = 0
//					cpuTotalThreadsPerSocket = 0
//					cpuStatus := cpuCapacity[j].CpuStatus
//					for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
//						cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
//						cpuTotalThreadsPerSocket++
//					}
//					cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
//					cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
//					cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j, normalized to 0-1
//					gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
//					gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate / (1.0 + gpuMemOverSellRate)//normalized to 0-1
//
//					// choose the instance configuration
//					maxCRE = float64(-1)  //reset maxCRE
//					maxCREConfigIndex = -1
//					for k := 0; k < len(resourcesConfigs); k++ {
//						if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
//							resourcesConfigs[k].GpuMemoryRate = 0
//						} else {
//							resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
//						}
//						/**
//						 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
//						 */
//						tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
//						tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
//						tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)
//
//						/**
//						 * check if the slot has enough resource
//						 */
//						if cpuConsumedRate + tempCpuQuotaRate > 1.01 ||
//							gpuCoreConsumedRate + tempGpuCoreQuotaRate > 1.01 ||
//							gpuMemConsumedRate + tempGpuMemQuotaRate > 1.01  {
//							//log.Printf("scheduler: current node has no enough resources for %dth pod config, skip to next pod config\n",k)
//							continue
//						} else {
//							tempCRE = float64(resourcesConfigs[k].ReqPerSecondMax) / (float64(resourcesConfigs[k].CpuThreads)*64 + float64(resourcesConfigs[k].GpuCorePercent)*142)
//							if gTypes.Greater(tempCRE, maxCRE) {
//								maxCRE = tempCRE
//								maxCREConfigIndex = k
//							}
//						}
//					}
//					if maxCREConfigIndex == -1 {
//						continue // no pod resource config can be placed into this socket
//					}
//
//					// update GPU memory allocation
//					cudaDeviceTh := j+1
//					if resourcesConfigs[maxCREConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
//						cudaDeviceTh = 0
//					}
//
//					/**
//					 * find a node to place function pod
//					 */
//
//					var cpuCoreThList []int
//					neededCores := resourcesConfigs[maxCREConfigIndex].CpuThreads >> 1 //hyper-threads
//					for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//						if cpuStatus[k].TotalFuncInstance == 0 {
//							cpuCoreThList = append(cpuCoreThList, k)
//							neededCores--
//						}
//					}
//					for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//						if cpuStatus[k].TotalFuncInstance != 0 {
//							if cpuStatus[k].TotalFuncInstance < CpuFuncInstanceThreshold && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate, CpuUsageRateThreshold) {
//								cpuCoreThList = append(cpuCoreThList, k)
//								neededCores--
//							}
//						}
//					}
//
//					if neededCores > 0 {
//						//log.Printf("scheduler: failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
//						continue
//					}
//
//					resourcesConfigs[maxCREConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation {
//						NodeTh: i,
//						CudaDeviceTh: cudaDeviceTh,
//						SocketTh: j,
//						CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
//					}
//
//					// todo scaling functions
//					createErr := createFuncInstance(funcName, namespace, resourcesConfigs[maxCREConfigIndex], "p", clientset)
//					if createErr != nil {
//						log.Println("scheduler-withoutFilter: create function instance failed ", createErr.Error())
//						// don't update the residualReq and execute for loop again
//					} else {
//						residualReq = residualReq - resourcesConfigs[maxCREConfigIndex].ReqPerSecondMax
//						residualFindFlag = true
//					}
//				} // per socket
//			} // per node
//			lock.Unlock()
//		} // per batch size
//	}
//	return residualReq
//}
//
//
//func ScaleUpFuncCapacity(funcName string, namespace string, latencySLO float64, reqArrivalRate int32, supportBatchGroup []int32, clientset *kubernetes.Clientset) {
//	//repository.UpdateFuncIsScalingIn(funcName,true)
//	funcStatus := repository.GetFunc(funcName)
//	if funcStatus == nil {
//		log.Printf("scheduler: function %s is nil in repository, error to scale up", funcName)
//		return
//	}
//	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
//	if err == nil {
//		//log.Printf("scheduler: reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
//	} else {
//		log.Println("scheduler: reading memory error:", err.Error())
//		return
//	}
//
//	cpuConsumedThreadsPerSocket := 0
//	cpuTotalThreadsPerSocket := 0
//	cpuOverSell := 0 //CPU threads overSell
//	gpuOverSell := 0 //GPU SM percentage
//	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
//	residualReq := reqArrivalRate
//	residualFindFlag := false
//	batchTryNum := 0
//
//	gpuMemConsumedRate := float64(0)
//	gpuCoreConsumedRate := float64(0)
//	cpuConsumedRate := float64(0)
//	gpuConsumedRate := float64(0)
//
//	tempGpuCoreQuotaRate := float64(0)
//	tempGpuMemQuotaRate := float64(0)
//	tempCpuQuotaRate := float64(0)
//	tempGpuQuotaRate := float64(0)
//
//	tempCRE := float64(0)
//	maxCRE := float64(-1)
//	maxCREConfigIndex := -1
//
//	for {
//		if residualReq > 0 {
//			residualFindFlag = false
//			if batchTryNum >= len(supportBatchGroup) {
//				wrappedErr := fmt.Errorf("scheduler: failed to find suitable batchsize for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
//					funcName, latencySLO, reqArrivalRate, residualReq)
//				log.Println(wrappedErr)
//				break
//			}
//			for batchIndex := 0; batchIndex < len(supportBatchGroup) && residualFindFlag == false; batchIndex++ {
//				resourcesConfigs, errInfer := inferResourceConfigsWithBatch(funcName, latencySLO, supportBatchGroup[batchIndex], residualReq)
//				if errInfer != nil {
//					batchTryNum ++
//					/*log.Print(errInfer.Error())
//					wrappedErr := fmt.Errorf("scheduler: batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
//						supportBatchGroup[batchIndex], funcName, latencySLO, reqArrivalRate, residualReq)
//					log.Println(wrappedErr)*/
//					continue
//				} else {
//					/*for _ , item := range resourcesConfigs {
//						log.Printf("scheduler: resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
//							funcName, latencySLO, item.ExecutionTime, supportBatchGroup[batchIndex], item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
//					}*/
//				}
//
//				lock.Lock()
//				//log.Println("scheduler: scale up locked---------------------------")
//				clusterCapConfig := repository.GetClusterCapConfig()
//				for i := 0; i < len(clusterCapConfig.ClusterCapacity) && residualFindFlag == false; i++ { // per node
//					cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
//					gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
//					gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate
//
//					/** CPU GPU consumed rate **/
//					cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
//					for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
//						/**
//						 * calculate CPU and GPU core, memory physical resource consumption rate for each slot
//						 */
//						cpuConsumedThreadsPerSocket = 0
//						cpuTotalThreadsPerSocket = 0
//						cpuStatus := cpuCapacity[j].CpuStatus
//						for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
//							cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
//							cpuTotalThreadsPerSocket++
//						}
//						cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
//						cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
//						cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j, normalized to 0-1
//						gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
//						gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate / (1.0 + gpuMemOverSellRate)//normalized to 0-1
//						/**
//						 * comparison of GPU core and memory
//						 */
//						if gTypes.GreaterEqual(gpuCoreConsumedRate, gpuMemConsumedRate) {
//							gpuConsumedRate = gpuCoreConsumedRate
//						} else {
//							gpuConsumedRate = gpuMemConsumedRate
//						}
//
//						//log.Println()
//						//			//log.Printf("scheduler: current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
//						//			//	i,
//						//			//	j,
//						//			//	j+1,
//						//			//	float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
//						//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
//						//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)
//
//						maxCRE = float64(-1)  //reset maxCRE
//						maxCREConfigIndex = -1
//						if gTypes.LessEqual(cpuConsumedRate, gpuConsumedRate) { // cpu is dominantly remained resource
//							for k := 0; k < len(resourcesConfigs); k++ {
//								if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
//									resourcesConfigs[k].GpuMemoryRate = 0
//								} else {
//									resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
//								}
//								/**
//								 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
//								 */
//								tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
//								tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
//								tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)
//
//								/**
//								 * check if the slot has enough resource
//								 */
//								if cpuConsumedRate + tempCpuQuotaRate > 1.01 ||
//									gpuCoreConsumedRate + tempGpuCoreQuotaRate > 1.01 ||
//									gpuMemConsumedRate + tempGpuMemQuotaRate > 1.01  {
//									//log.Printf("scheduler: current node has no enough resources for %dth pod config, skip to next pod config\n",k)
//									continue
//								} else {
//									//log.Printf("scheduler: current node has enough resources for %dth pod config, skip to next pod config\n",k)
//								}
//								/**
//								 * comparison of GPU core and memory
//								 */
//								if gTypes.GreaterEqual(tempGpuCoreQuotaRate, tempGpuMemQuotaRate) {
//									tempGpuQuotaRate = tempGpuCoreQuotaRate
//								} else {
//									tempGpuQuotaRate = tempGpuMemQuotaRate
//								}
//								/**
//								 * match the dominated resource
//								 */
//								//log.Printf("scheduler: k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
//								if gTypes.GreaterEqual(tempCpuQuotaRate, tempGpuQuotaRate) {
//									tempCRE = float64(resourcesConfigs[k].ReqPerSecondMax) / (float64(resourcesConfigs[k].CpuThreads)*64 + float64(resourcesConfigs[k].GpuCorePercent)*142)
//									if gTypes.Greater(tempCRE, maxCRE) {
//										maxCRE = tempCRE
//										maxCREConfigIndex = k
//									}
//								}
//							}
//
//							//log.Printf("scheduler: CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
//							//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
//						} else { // GPU core is dominantly remained resource
//							for k := 0; k < len(resourcesConfigs); k++ {
//								if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
//									resourcesConfigs[k].GpuMemoryRate = 0
//								} else {
//									resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
//								}
//								/**
//								 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
//								 */
//								tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
//								tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
//								tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)
//
//								/**
//								 * check if the slot has enough resource
//								 */
//								if cpuConsumedRate + tempCpuQuotaRate > 1.01 ||
//									gpuCoreConsumedRate + tempGpuCoreQuotaRate > 1.01 ||
//									gpuMemConsumedRate + tempGpuMemQuotaRate > 1.01  {
//									//log.Printf("scheduler: current node has no enough resources for %dth pod config, skip to next pod config\n",k)
//									continue
//								} else {
//									//log.Printf("scheduler: current node has enough resources for %dth pod config, skip to next pod config\n",k)
//								}
//								/**
//								 * comparison of GPU core and memory
//								 */
//								if gTypes.GreaterEqual(tempGpuCoreQuotaRate, tempGpuMemQuotaRate) {
//									tempGpuQuotaRate = tempGpuCoreQuotaRate
//								} else {
//									tempGpuQuotaRate = tempGpuMemQuotaRate
//								}
//								/**
//								 * match the dominated resource
//								 */
//								//log.Printf("scheduler: k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
//								if gTypes.LessEqual(tempCpuQuotaRate, tempGpuQuotaRate) {
//									tempCRE = float64(resourcesConfigs[k].ReqPerSecondMax) / (float64(resourcesConfigs[k].CpuThreads)*64 + float64(resourcesConfigs[k].GpuCorePercent)*142)
//									if gTypes.Greater(tempCRE, maxCRE) {
//										maxCRE = tempCRE
//										maxCREConfigIndex = k
//									}
//								}
//							}
//							//log.Printf("scheduler: GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
//							//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
//						}
//
//						if maxCREConfigIndex == -1 {
//							continue // no pod resource config can be placed into this socket
//						}
//						// update GPU memory allocation
//						cudaDeviceTh := j+1
//						if resourcesConfigs[maxCREConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
//							cudaDeviceTh = 0
//						}
//
//						//if minResourceQuotaPosDiffIndex == -1 {
//						//	log.Printf("scheduler: choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, maxResourceQuotaNagDiff=%f\n",
//						//		pickConfigIndex,
//						//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
//						//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
//						//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
//						//		maxResourceQuotaNagDiff)
//						//} else {
//						//	log.Printf("scheduler: choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, minResourceQuotaPosDiff=%f\n",
//						//		pickConfigIndex,
//						//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
//						//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
//						//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
//						//		minResourceQuotaPosDiff)
//						//}
//
//						/**
//						 * find a node to place function pod
//						 */
//
//						var cpuCoreThList []int
//						neededCores := resourcesConfigs[maxCREConfigIndex].CpuThreads >> 1 //hyper-threads
//						for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//							if cpuStatus[k].TotalFuncInstance == 0 {
//								cpuCoreThList = append(cpuCoreThList, k)
//								neededCores--
//							}
//						}
//						for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//							if cpuStatus[k].TotalFuncInstance != 0 {
//								if cpuStatus[k].TotalFuncInstance < CpuFuncInstanceThreshold && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate, CpuUsageRateThreshold) {
//									cpuCoreThList = append(cpuCoreThList, k)
//									neededCores--
//								}
//							}
//						}
//
//						if neededCores > 0 {
//							//log.Printf("scheduler: failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
//							continue
//						}
//						//log.Printf("scheduler: decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
//						//	i,
//						//	j,
//						//	cudaDeviceTh,
//						//	cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
//						//	cpuTotalThreadsPerSocket + cpuOverSell,
//						//	gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
//						//	1 + gpuMemOverSellRate,
//						//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
//						//	1 + float64(gpuOverSell)/100)
//
//						resourcesConfigs[maxCREConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation {
//							NodeTh: i,
//							CudaDeviceTh: cudaDeviceTh,
//							SocketTh: j,
//							CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
//						}
//
//						// todo scaling functions
//						createErr := createFuncInstance(funcName, namespace, resourcesConfigs[maxCREConfigIndex],"i", clientset)
//						if createErr != nil {
//							log.Println("scheduler: create function instance failed ", createErr.Error())
//							// don't update the residualReq and execute for loop again
//						} else {
//							//log.Printf("scheduler: create function instance for function%s successfully, residualReq=%d-%d=%d \n",
//							//	funcName, residualReq, resourcesConfigs[pickConfigIndex].ReqPerSecondMax, residualReq - resourcesConfigs[pickConfigIndex].ReqPerSecondMax)
//							residualReq = residualReq - resourcesConfigs[maxCREConfigIndex].ReqPerSecondMax
//						}
//						residualFindFlag = true
//						batchTryNum = 0
//					} // per socket
//				} // per node
//				lock.Unlock()
//				//repository.UpdateFuncScalingLock(funcName,false)
//				//log.Println("scheduler: scale up unlocked---------------------------")
//			} // per batch size
//			if residualFindFlag == false {
//				log.Printf("scheduler: failed to find suitable node for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d and try the without filter\n",
//					funcName, latencySLO, reqArrivalRate, residualReq)
//				residualReq = scaleUpFuncCapacityWithoutFilter(funcName, namespace, latencySLO, reqArrivalRate, supportBatchGroup, clientset)
//			}
//		} else {
//			break // break the for loop, residualReq <= 0
//		}
//	}
//	return
//}
//
//
//
//
//func scaleUpFuncCapacityWithoutFilter(funcName string, namespace string, latencySLO float64, residualReq int32, supportBatchGroup []int32, clientset *kubernetes.Clientset) int32 {
//
//	funcStatus := repository.GetFunc(funcName)
//	if funcStatus == nil {
//		log.Printf("scheduler-withoutFilter: function %s is nil in repository, error to scale up", funcName)
//		return residualReq
//	}
//	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
//	if err == nil {
//		//log.Printf("scheduler: reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
//	} else {
//		log.Println("scheduler-withoutFilter: reading memory error:", err.Error())
//		return residualReq
//	}
//
//	cpuConsumedThreadsPerSocket := 0
//	cpuTotalThreadsPerSocket := 0
//	cpuOverSell := 0 //CPU threads overSell
//	gpuOverSell := 0 //GPU SM percentage
//	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
//
//	batchTryNum := 0
//	residualFindFlag := false
//
//	cpuConsumedRate := float64(0)
//	gpuMemConsumedRate := float64(0)
//	gpuCoreConsumedRate := float64(0)
//
//	tempGpuCoreQuotaRate := float64(0)
//	tempGpuMemQuotaRate := float64(0)
//	tempCpuQuotaRate := float64(0)
//
//	tempCRE := float64(0)
//	maxCRE := float64(-1)
//	maxCREConfigIndex := -1
//
//	if residualReq > 0 {
//		for batchIndex := 0; batchIndex < len(supportBatchGroup) && residualFindFlag == false; batchIndex++ {
//			resourcesConfigs, errInfer := inferResourceConfigsWithBatch(funcName, latencySLO, supportBatchGroup[batchIndex], residualReq)
//			if errInfer != nil {
//				batchTryNum ++
//				if batchTryNum >= len(supportBatchGroup) {
//					wrappedErr := fmt.Errorf("scheduler-withoutFilter: failed to find suitable batchsize for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
//						funcName, latencySLO, residualReq, residualReq)
//					log.Println(wrappedErr)
//					return residualReq // stop explore batch
//				} else {
//					continue // try next smaller batch size
//				}
//			} else {
//				//try the allocation, do the following code
//			}
//
//			lock.Lock()
//			clusterCapConfig := repository.GetClusterCapConfig()
//			for i := 0; i < len(clusterCapConfig.ClusterCapacity) && residualFindFlag == false; i++ { // per node
//				cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
//				gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
//				gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate
//
//				/** CPU GPU consumed rate **/
//				cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
//				for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
//					/**
//					 * calculate CPU and GPU core, memory physical resource consumption rate for each slot
//					 */
//					cpuConsumedThreadsPerSocket = 0
//					cpuTotalThreadsPerSocket = 0
//					cpuStatus := cpuCapacity[j].CpuStatus
//					for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
//						cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
//						cpuTotalThreadsPerSocket++
//					}
//					cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
//					cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
//					cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j, normalized to 0-1
//					gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
//					gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate / (1.0 + gpuMemOverSellRate)//normalized to 0-1
//
//					// choose the instance configuration
//					maxCRE = float64(-1)  //reset maxCRE
//					maxCREConfigIndex = -1
//					for k := 0; k < len(resourcesConfigs); k++ {
//						if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
//							resourcesConfigs[k].GpuMemoryRate = 0
//						} else {
//							resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
//						}
//						/**
//						 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
//						 */
//						tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
//						tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
//						tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)
//
//						/**
//						 * check if the slot has enough resource
//						 */
//						if cpuConsumedRate + tempCpuQuotaRate > 1.01 ||
//							gpuCoreConsumedRate + tempGpuCoreQuotaRate > 1.01 ||
//							gpuMemConsumedRate + tempGpuMemQuotaRate > 1.01  {
//							//log.Printf("scheduler: current node has no enough resources for %dth pod config, skip to next pod config\n",k)
//							continue
//						} else {
//							tempCRE = float64(resourcesConfigs[k].ReqPerSecondMax) / (float64(resourcesConfigs[k].CpuThreads)*64 + float64(resourcesConfigs[k].GpuCorePercent)*142)
//							if gTypes.Greater(tempCRE, maxCRE) {
//								maxCRE = tempCRE
//								maxCREConfigIndex = k
//							}
//						}
//					}
//					if maxCREConfigIndex == -1 {
//						continue // no pod resource config can be placed into this socket
//					}
//
//					// update GPU memory allocation
//					cudaDeviceTh := j+1
//					if resourcesConfigs[maxCREConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
//						cudaDeviceTh = 0
//					}
//
//					/**
//					 * find a node to place function pod
//					 */
//
//					var cpuCoreThList []int
//					neededCores := resourcesConfigs[maxCREConfigIndex].CpuThreads >> 1 //hyper-threads
//					for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//						if cpuStatus[k].TotalFuncInstance == 0 {
//							cpuCoreThList = append(cpuCoreThList, k)
//							neededCores--
//						}
//					}
//					for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
//						if cpuStatus[k].TotalFuncInstance != 0 {
//							if cpuStatus[k].TotalFuncInstance < CpuFuncInstanceThreshold && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate, CpuUsageRateThreshold) {
//								cpuCoreThList = append(cpuCoreThList, k)
//								neededCores--
//							}
//						}
//					}
//
//					if neededCores > 0 {
//						//log.Printf("scheduler: failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
//						continue
//					}
//
//					resourcesConfigs[maxCREConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation {
//						NodeTh: i,
//						CudaDeviceTh: cudaDeviceTh,
//						SocketTh: j,
//						CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
//					}
//
//					// todo scaling functions
//					createErr := createFuncInstance(funcName, namespace, resourcesConfigs[maxCREConfigIndex],"i", clientset)
//					if createErr != nil {
//						log.Println("scheduler-withoutFilter: create function instance failed ", createErr.Error())
//						// don't update the residualReq and execute for loop again
//					} else {
//						residualReq = residualReq - resourcesConfigs[maxCREConfigIndex].ReqPerSecondMax
//					}
//					residualFindFlag = true
//				} // per socket
//			} // per node
//			lock.Unlock()
//		} // per batch size
//	}
//	return residualReq
//}
//
//
//func ScaleDownFuncCapacity(funcName string, namespace string, deletedFuncPodConfig []*gTypes.FuncPodConfig, clientset *kubernetes.Clientset) {
//	lock.Lock()
//	//log.Println("scheduler: scale down locked---------------------------")
//	err := deleteFuncInstance(funcName, namespace, deletedFuncPodConfig, clientset)
//	if err != nil {
//		log.Println("scheduler: delete function instance failed ", err.Error())
//	} else {
//		//log.Println("scheduler: delete function instance successfully")
//	}
//	//log.Println("scheduler: scale down unlocked---------------------------")
//	lock.Unlock()
//	return
//}
//
//
//
//
//////File  : scheduler.go
//////Author: Yanan Yang
//////Date  : 2020/4/7
//////Desc  : based on the least abs diff for the pod configuration selection
//
////package controller
////
////import (
////	"fmt"
////	"github.com/openfaas/faas-netes/gpu/repository"
////	gTypes "github.com/openfaas/faas-netes/gpu/types"
////	"k8s.io/client-go/kubernetes"
////	"log"
////	"strconv"
////	"sync"
////	"time"
////)
////var lock sync.Mutex
/////*const cpuOverSell = int32(0) //threads overSell
////const gpuOverSell= int32(20) //SM percentage
////const gpuMemOverSellRate = float64(-0.1)*/
////func ScaleUpFuncCapacity(funcName string, namespace string, latencySLO float64, reqArrivalRate int32, supportBatchGroup []int32, clientset *kubernetes.Clientset) {
////	//repository.UpdateFuncIsScalingIn(funcName,true)
////	funcStatus := repository.GetFunc(funcName)
////	if funcStatus == nil {
////		log.Printf("scheduler: function %s is nil in repository, error to scale up", funcName)
////		return
////	}
////	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
////	if err == nil {
////		//log.Printf("scheduler: reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
////	} else {
////		log.Println("scheduler: reading memory error:", err.Error())
////		return
////	}
////
////
////	residualReq := reqArrivalRate
////	residualFindFlag := false
////	batchTryNum := 0
////
////	maxResourceQuotaNagDiffIndex := -1
////	minResourceQuotaPosDiffIndex := -1
////	pickConfigIndex := -1
////	maxResourceQuotaNagDiff := float64(-999)
////	minResourceQuotaPosDiff := float64(999)
////	tempGpuCoreQuota := float64(0)
////	tempCpuQuota := float64(0)
////	tempDiffQuota := float64(0)
////
////	cpuConsumedRate := float64(0)
////	gpuMemConsumedRate := float64(0)
////	gpuCoreConsumedRate := float64(0)
////	tempThroughIntensity := float64(0)
////	tempMinResourceQuotaPosThroughIntensity := float64(0)
////	tempMaxResourceQuotaNagThroughIntensity := float64(0)
////
////
////	cpuConsumedThreadsPerSocket := int(0)
////	cpuTotalThreadsPerSocket := int(0)
////	cpuOverSell := 0 //CPU threads overSell
////	gpuOverSell := 0 //GPU SM percentage
////	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
////
////	//repository.UpdateFuncScalingLock(funcName,false) //init the scaling lock
////	for {
////		//if funcStatus.FunctionScalingLock == true {
////		//	time.Sleep(time.Millisecond*200)
////		//	continue
////		//}
////		if residualReq > 0 {
////
////			//repository.UpdateFuncScalingLock(funcName,true) //lock
////			residualFindFlag = false
////			if batchTryNum == len(supportBatchGroup) {
////				wrappedErr := fmt.Errorf("scheduler: failed to find suitable batchsize for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
////					funcName, latencySLO, reqArrivalRate, residualReq)
////				log.Println(wrappedErr)
////				break
////			}
////			for batchIndex := 0; batchIndex < len(supportBatchGroup) && residualFindFlag == false; batchIndex++ {
////				resourcesConfigs, errInfer := inferResourceConfigsWithBatch(funcName, latencySLO, supportBatchGroup[batchIndex], residualReq)
////				if errInfer != nil {
////					batchTryNum ++
////					/*log.Print(errInfer.Error())
////					wrappedErr := fmt.Errorf("scheduler: batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
////						supportBatchGroup[batchIndex], funcName, latencySLO, reqArrivalRate, residualReq)
////					log.Println(wrappedErr)*/
////					continue
////				} else {
////					/*for _ , item := range resourcesConfigs {
////						log.Printf("scheduler: resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
////							funcName, latencySLO, item.ExecutionTime, supportBatchGroup[batchIndex], item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
////					}*/
////				}
////
////				lock.Lock()
////				//log.Println("scheduler: scale up locked---------------------------")
////				clusterCapConfig := repository.GetClusterCapConfig()
////				for i := 0; i < len(clusterCapConfig.ClusterCapacity) && residualFindFlag == false; i++ { // per node
////					cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
////					gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
////					gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate
////
////					/** CPU GPU consumed rate **/
////					cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
////					for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
////						/**
////						 * calculate CPU and GPU physical consumption rate
////						 */
////						cpuConsumedThreadsPerSocket = 0
////						cpuTotalThreadsPerSocket = 0
////						cpuStatus := cpuCapacity[j].CpuStatus
////						for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
////							cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
////							cpuTotalThreadsPerSocket++
////						}
////						cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
////						cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
////						cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j
////						gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate
////						gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100)
////
////						//log.Println()
////						//log.Printf("scheduler: current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
////						//	i,
////						//	j,
////						//	j+1,
////						//	float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
////						//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
////						//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)
////						/**
////						 * allocate resource
////						 */
////						maxResourceQuotaNagDiffIndex = -1
////						minResourceQuotaPosDiffIndex = -1
////						pickConfigIndex = -1
////						maxResourceQuotaNagDiff = float64(-999)
////						minResourceQuotaPosDiff = float64(999)
////						if gTypes.LessEqual(cpuConsumedRate, gpuCoreConsumedRate) { // cpu is dominantly remained resource
////							for k := 0; k < len(resourcesConfigs); k++ {
////								if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
////									resourcesConfigs[k].GpuMemoryRate = 0
////								} else {
////									resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
////								}
////								if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
////									gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
////									gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
////									//log.Printf("scheduler: current node has no enough resources for %dth pod config, skip to next pod config\n",k)
////									continue
////								} else {
////									//log.Printf("scheduler: current node has enough resources for %dth pod config, skip to next pod config\n",k)
////								}
////								tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
////								tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
////								tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
////								//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
////								if gTypes.Greater(tempDiffQuota,0) {
////									if gTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
////										minResourceQuotaPosDiff = tempDiffQuota
////										minResourceQuotaPosDiffIndex = k
////									} else if gTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
////										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
////										tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
////											(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////												float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////										if gTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
////											minResourceQuotaPosDiffIndex = k
////										}
////									}
////								} else {
////									if gTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
////										maxResourceQuotaNagDiff = tempDiffQuota
////										maxResourceQuotaNagDiffIndex = k
////									} else if gTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
////										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
////										tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
////											(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////												float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////										if gTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
////											maxResourceQuotaNagDiffIndex = k
////										}
////									}
////								}
////
////							}
////							//log.Printf("scheduler: warm CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
////							//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
////						} else { // GPU core is dominantly remained resource
////							for k := 0; k < len(resourcesConfigs); k++ {
////								if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
////									resourcesConfigs[k].GpuMemoryRate = 0
////								} else {
////									resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
////								}
////								if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
////									gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
////									gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
////									//log.Printf("scheduler: current node has no enough resources for %dth pod config, skip to next pod config\n",k)
////									continue
////								} else {
////									//log.Printf("scheduler: current node has enough resources for %dth pod config, skip to next pod config\n",k)
////								}
////								tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
////								tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
////								tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
////								//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
////								if gTypes.Greater(tempDiffQuota,0) {
////									if gTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
////										minResourceQuotaPosDiff = tempDiffQuota
////										minResourceQuotaPosDiffIndex = k
////									} else if gTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
////										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
////										tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
////											(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////												float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////										if gTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
////											minResourceQuotaPosDiffIndex = k
////										}
////									}
////								} else {
////									if gTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
////										maxResourceQuotaNagDiff = tempDiffQuota
////										maxResourceQuotaNagDiffIndex = k
////									} else if gTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
////										tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
////										tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
////											(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////												float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////										if gTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
////											maxResourceQuotaNagDiffIndex = k
////										}
////									}
////								}
////							}
////							//log.Printf("scheduler: warm GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
////							//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
////						}
////						if minResourceQuotaPosDiffIndex == -1 {
////							pickConfigIndex = maxResourceQuotaNagDiffIndex
////						} else {
////							pickConfigIndex = minResourceQuotaPosDiffIndex
////						}
////						if pickConfigIndex == -1 {
////							continue
////						}
////						// update GPU memory allocation
////						cudaDeviceTh := j+1
////						if resourcesConfigs[pickConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
////							cudaDeviceTh = 0
////						}
////
////						//if minResourceQuotaPosDiffIndex == -1 {
////						//	log.Printf("scheduler: choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, maxResourceQuotaNagDiff=%f\n",
////						//		pickConfigIndex,
////						//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
////						//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
////						//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
////						//		maxResourceQuotaNagDiff)
////						//} else {
////						//	log.Printf("scheduler: choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, minResourceQuotaPosDiff=%f\n",
////						//		pickConfigIndex,
////						//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
////						//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
////						//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
////						//		minResourceQuotaPosDiff)
////						//}
////
////						/**
////						 * find a node to place function pod
////						 */
////
////						var cpuCoreThList []int
////						neededCores := resourcesConfigs[pickConfigIndex].CpuThreads >> 1 //hyper-threads
////						for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
////							if cpuStatus[k].TotalFuncInstance == 0 {
////								cpuCoreThList = append(cpuCoreThList, k)
////								neededCores--
////							}
////						}
////						for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
////							if cpuStatus[k].TotalFuncInstance != 0 {
////								if cpuStatus[k].TotalFuncInstance < 3 && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate,0.8) {
////									cpuCoreThList = append(cpuCoreThList, k)
////									neededCores--
////								}
////							}
////						}
////
////						if neededCores > 0 {
////							//log.Printf("scheduler: failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
////							continue
////						}
////						//log.Printf("scheduler: decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
////						//	i,
////						//	j,
////						//	cudaDeviceTh,
////						//	cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
////						//	cpuTotalThreadsPerSocket + cpuOverSell,
////						//	gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
////						//	1 + gpuMemOverSellRate,
////						//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
////						//	1 + float64(gpuOverSell)/100)
////
////						resourcesConfigs[pickConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation {
////							NodeTh:       i,
////							CudaDeviceTh: cudaDeviceTh,
////							SocketTh: j,
////							CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
////						}
////
////						// todo scaling functions
////						createErr := createFuncInstance(funcName, namespace, resourcesConfigs[pickConfigIndex],"i", clientset)
////						if createErr != nil {
////							log.Println("scheduler: create function instance failed ", createErr.Error())
////							// don't update the residualReq and execute for loop again
////						} else {
////							//log.Printf("scheduler: create function instance for function%s successfully, residualReq=%d-%d=%d \n",
////							//	funcName, residualReq, resourcesConfigs[pickConfigIndex].ReqPerSecondMax, residualReq - resourcesConfigs[pickConfigIndex].ReqPerSecondMax)
////							residualReq = residualReq - resourcesConfigs[pickConfigIndex].ReqPerSecondMax
////						}
////						residualFindFlag = true
////						batchTryNum = 0
////					} // per socket
////				} // per node
////				lock.Unlock()
////				//repository.UpdateFuncScalingLock(funcName,false)
////				//log.Println("scheduler: scale up unlocked---------------------------")
////			} // per batch size
////			if residualFindFlag == false {
////				log.Printf("scheduler: failed to find suitable node for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
////					funcName, latencySLO, reqArrivalRate, residualReq)
////				time.Sleep(time.Second*30)
////				break
////			}
////		} else {
////			break // residualReq <= 0
////		}
////	}
////	//repository.UpdateFuncIsScalingIn(funcName,false)
////	return
////}
////func CreatePreWarmPod(funcName string, namespace string, latencySLO float64, batchSize int32, clientset *kubernetes.Clientset){
////	funcStatus := repository.GetFunc(funcName)
////	if funcStatus == nil {
////		log.Printf("scheduler: warm function %s is nil in repository, error to read GPU memory", funcName)
////		return
////	}
////	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
////	if err == nil {
////		//log.Printf("scheduler: warm reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
////	} else {
////		log.Println("scheduler: warm read memory error:", err.Error())
////		return
////	}
////
////	resourcesConfigs, err := inferResourceConfigsWithBatch(funcName, latencySLO, batchSize, 1)
////	if err != nil {
////		/*log.Print(err.Error())
////		wrappedErr := fmt.Errorf("scheduler: CreatePrewarmPod failed batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
////			batchSize, funcName, latencySLO, 1, 1)
////		log.Println(wrappedErr)*/
////		return
////	} else {
////		/*for _, item := range resourcesConfigs {
////			log.Printf("scheduler: warm resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
////				funcName, latencySLO, item.ExecutionTime, batchSize, item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
////		}*/
////	}
////	maxResourceQuotaNagDiffIndex := -1
////	minResourceQuotaPosDiffIndex := -1
////	pickConfigIndex := -1
////	maxResourceQuotaNagDiff := float64(-999)
////	minResourceQuotaPosDiff := float64(999)
////	tempGpuCoreQuota := float64(0)
////	tempCpuQuota := float64(0)
////	tempDiffQuota := float64(0)
////
////	cpuConsumedRate := float64(0)
////	gpuMemConsumedRate := float64(0)
////	gpuCoreConsumedRate := float64(0)
////	tempThroughIntensity := float64(0)
////	tempMinResourceQuotaPosThroughIntensity := float64(0)
////	tempMaxResourceQuotaNagThroughIntensity := float64(0)
////
////
////	cpuConsumedThreadsPerSocket := int(0)
////	cpuTotalThreadsPerSocket := int(0)
////	cpuOverSell := 0 //CPU threads overSell
////	gpuOverSell := 0 //GPU SM percentage
////	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
////	residualFindFlag := false
////	//for {
////	//	if repository.ClusterCapConfigLock == true {
////	//		time.Sleep(time.Millisecond * 50)
////	//	} else {
////	//		repository.ClusterCapConfigLock = true
////	//		break
////	//	}
////	//}
////
////	lock.Lock()
////	//log.Println("scheduler: warm locked---------------------------")
////	clusterCapConfig := repository.GetClusterCapConfig()
////	for i := 0; i < len(clusterCapConfig.ClusterCapacity) && residualFindFlag == false; i++ { // per node
////		cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
////		gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
////		gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate
////
////		/** CPU GPU consumed rate **/
////		cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
////		for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
////			/**
////			 * calculate CPU and GPU physical consumption rate
////			 */
////			cpuConsumedThreadsPerSocket = 0
////			cpuTotalThreadsPerSocket = 0
////			cpuStatus := cpuCapacity[j].CpuStatus
////			for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
////				cpuConsumedThreadsPerSocket+=cpuStatus[k].TotalFuncInstance
////				cpuTotalThreadsPerSocket++
////			}
////			cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
////			cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
////			cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket + cpuOverSell) // cpu usage rate in node i socket j, normalized to 0-1
////			gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate //normalized to 0-1
////			gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
////			//log.Println()
////			//			//log.Printf("scheduler: warm current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
////			//			//	i,
////			//			//	j,
////			//			//	j+1,
////			//			//	float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
////			//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
////			//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)
////			/**
////			 * allocate resource
////			 */
////			maxResourceQuotaNagDiffIndex = -1
////			minResourceQuotaPosDiffIndex = -1
////			pickConfigIndex = -1
////			maxResourceQuotaNagDiff = float64(-999)
////			minResourceQuotaPosDiff = float64(999)
////			if gTypes.LessEqual(cpuConsumedRate, gpuCoreConsumedRate) { // cpu is dominantly remained resource
////				for k := 0; k < len(resourcesConfigs); k++ {
////					if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
////						resourcesConfigs[k].GpuMemoryRate = 0
////					} else {
////						resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
////					}
////					if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
////						gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
////						gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
////						//log.Printf("scheduler: warm current node has no enough resources for %dth pod config, skip to next pod config\n",k)
////						continue
////					} else {
////						//log.Printf("scheduler: warm current node has enough resources for %dth pod config, skip to next pod config\n",k)
////					}
////
////					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
////					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
////					tempDiffQuota = tempCpuQuota - tempGpuCoreQuota
////					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
////					if gTypes.Greater(tempDiffQuota,0) {
////						if gTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
////							minResourceQuotaPosDiff = tempDiffQuota
////							minResourceQuotaPosDiffIndex = k
////						} else if gTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
////							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
////							tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
////								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////							if gTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
////								minResourceQuotaPosDiffIndex = k
////							}
////						}
////					} else {
////						if gTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
////							maxResourceQuotaNagDiff = tempDiffQuota
////							maxResourceQuotaNagDiffIndex = k
////						} else if gTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
////							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota+tempGpuCoreQuota)
////							tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
////								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////							if gTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
////								maxResourceQuotaNagDiffIndex = k
////							}
////						}
////					}
////
////				}
////				//log.Printf("scheduler: warm CPU is in lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
////				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
////			} else { // GPU core is dominantly remained resource
////				for k := 0; k < len(resourcesConfigs); k++ {
////					if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
////						resourcesConfigs[k].GpuMemoryRate = 0
////					} else {
////						resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc)/float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
////					}
////					if cpuConsumedThreadsPerSocket + int(resourcesConfigs[k].CpuThreads) > (cpuTotalThreadsPerSocket + cpuOverSell) ||
////						gpuCoreConsumedRate + float64(resourcesConfigs[k].GpuCorePercent)/float64(gpuOverSell+100) > 1.01 ||
////						gpuMemConsumedRate + resourcesConfigs[k].GpuMemoryRate > (1.01 + gpuMemOverSellRate) {
////						//log.Printf("scheduler: current node has no enough resources for %dth pod config, skip to next pod config\n",k)
////						continue
////					} else {
////						//log.Printf("scheduler: current node has enough resources for %dth pod config, skip to next pod config\n",k)
////					}
////
////					tempCpuQuota = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell)
////					tempGpuCoreQuota = float64(resourcesConfigs[k].GpuCorePercent) / float64(100 + gpuOverSell)
////					tempDiffQuota = tempGpuCoreQuota - tempCpuQuota
////					//log.Printf("scheduler: warm k=%d, resourceConfig=%+v, diffQuota=%f\n", k, resourcesConfigs[k], tempDiffQuota)
////					if gTypes.Greater(tempDiffQuota,0) {
////						if gTypes.Less(tempDiffQuota, minResourceQuotaPosDiff) {
////							minResourceQuotaPosDiff = tempDiffQuota
////							minResourceQuotaPosDiffIndex = k
////						} else if gTypes.Equal(tempDiffQuota, minResourceQuotaPosDiff) {
////							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
////							tempMinResourceQuotaPosThroughIntensity = float64(resourcesConfigs[minResourceQuotaPosDiffIndex].ReqPerSecondMax)/
////								(float64(resourcesConfigs[minResourceQuotaPosDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////									float64(resourcesConfigs[minResourceQuotaPosDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////							if gTypes.Greater(tempThroughIntensity, tempMinResourceQuotaPosThroughIntensity) {
////								minResourceQuotaPosDiffIndex = k
////							}
////						}
////					} else {
////						if gTypes.Greater(tempDiffQuota, maxResourceQuotaNagDiff) {
////							maxResourceQuotaNagDiff = tempDiffQuota
////							maxResourceQuotaNagDiffIndex = k
////						} else if gTypes.Equal(tempDiffQuota, maxResourceQuotaNagDiff) {
////							tempThroughIntensity = float64(resourcesConfigs[k].ReqPerSecondMax)/(tempCpuQuota + tempGpuCoreQuota)
////							tempMaxResourceQuotaNagThroughIntensity = float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].ReqPerSecondMax)/
////								(float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket + cpuOverSell) +
////									float64(resourcesConfigs[maxResourceQuotaNagDiffIndex].GpuCorePercent) / float64(100 + gpuOverSell))
////							if gTypes.Greater(tempThroughIntensity, tempMaxResourceQuotaNagThroughIntensity) {
////								maxResourceQuotaNagDiffIndex = k
////							}
////						}
////					}
////				}
////				//log.Printf("scheduler: warm GPU is lowest consumed rate, resourceConfigs: minResourceQuotaPosDiff=%f, index=%d, maxResourceQuotaNagDiff=%f, index=%d\n",
////				//	minResourceQuotaPosDiff, minResourceQuotaPosDiffIndex, maxResourceQuotaNagDiff, maxResourceQuotaNagDiffIndex)
////			}
////
////			if minResourceQuotaPosDiffIndex == -1 {
////				pickConfigIndex = maxResourceQuotaNagDiffIndex
////			} else {
////				pickConfigIndex = minResourceQuotaPosDiffIndex
////			}
////			if pickConfigIndex == -1 {
////				continue // no pod resource config can be placed into this socket
////			}
////
////			// update GPU memory allocation
////			cudaDeviceTh := j+1
////			if resourcesConfigs[pickConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
////				cudaDeviceTh = 0
////			}
////
////			//if minResourceQuotaPosDiffIndex == -1 {
////			//	log.Printf("scheduler: warm choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, maxResourceQuotaNagDiff=%f\n",
////			//		pickConfigIndex,
////			//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
////			//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
////			//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
////			//		maxResourceQuotaNagDiff)
////			//} else {
////			//	log.Printf("scheduler: warm choosed %dth resourceConfigs with physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f, minResourceQuotaPosDiff=%f\n",
////			//		pickConfigIndex,
////			//		float64(resourcesConfigs[pickConfigIndex].CpuThreads) / float64(cpuTotalThreadsPerSocket),
////			//		resourcesConfigs[pickConfigIndex].GpuMemoryRate,
////			//		float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
////			//		minResourceQuotaPosDiff)
////			//}
////
////			/**
////			 * allocate CPU threads
////			 */
////			var cpuCoreThList []int
////			neededCores := resourcesConfigs[pickConfigIndex].CpuThreads >> 1 //hyper-threads
////			for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
////				if cpuStatus[k].TotalFuncInstance == 0 {
////					cpuCoreThList = append(cpuCoreThList, k)
////					neededCores--
////				}
////			}
////			for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
////				if cpuStatus[k].TotalFuncInstance != 0 {
////					if cpuStatus[k].TotalFuncInstance < 3 && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate,0.8) {
////						cpuCoreThList = append(cpuCoreThList, k)
////						neededCores--
////					}
////				}
////			}
////
////			if neededCores > 0 {
////				//log.Printf("scheduler: warm failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
////				continue
////			}
////
////			//log.Printf("scheduler: warm decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
////			//	i,
////			//	j,
////			//	cudaDeviceTh,
////			//	cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
////			//	cpuTotalThreadsPerSocket + cpuOverSell,
////			//	gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
////			//	1 + gpuMemOverSellRate,
////			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
////			//	1 + float64(gpuOverSell)/100)
////
////			resourcesConfigs[pickConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation {
////				NodeTh:       i,
////				CudaDeviceTh: cudaDeviceTh,
////				SocketTh: j,
////				CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
////			}
////
////			createErr := createFuncInstance(funcName, namespace, resourcesConfigs[pickConfigIndex], "p", clientset)
////			if createErr != nil {
////				log.Println("scheduler: warm create prewarm function instance failed ", createErr.Error())
////			} else {
////				//log.Printf("scheduler: warm create prewarm function instance for function %s successfully\n", funcName)
////			}
////			residualFindFlag = true
////		} // per socket
////	} // per node
////
////	lock.Unlock()
////	//log.Println("scheduler: warm unlocked---------------------------")
////	return
////}
////func ScaleDownFuncCapacity(funcName string, namespace string, deletedFuncPodConfig []*gTypes.FuncPodConfig, clientset *kubernetes.Clientset) {
////	lock.Lock()
////	//log.Println("scheduler: scale down locked---------------------------")
////	err := deleteFuncInstance(funcName, namespace, deletedFuncPodConfig, clientset)
////	if err != nil {
////		log.Println("scheduler: delete function instance failed ", err.Error())
////	} else {
////		//log.Println("scheduler: delete function instance successfully")
////	}
////	//log.Println("scheduler: scale down unlocked---------------------------")
////	lock.Unlock()
////	return
////}
////
////
