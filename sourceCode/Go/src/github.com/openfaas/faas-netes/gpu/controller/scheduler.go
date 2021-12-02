//File  : scheduler.go
//Author: Yanan Yang
//Date  : 2020/4/7
package controller

import (
	"fmt"
	"github.com/openfaas/faas-netes/gpu/repository"
	gTypes "github.com/openfaas/faas-netes/gpu/types"
	"k8s.io/client-go/kubernetes"
	"log"
	"strconv"
	"sync"
)
var lock sync.Mutex

const CpuUsageRateThreshold = 0.8
const CpuFuncInstanceThreshold = 3


func CreatePreWarmPod(funcName string, namespace string, latencySLO float64, batchSize int32, clientset *kubernetes.Clientset){
	funcStatus := repository.GetFunc(funcName)
	if funcStatus == nil {
		log.Printf("scheduler: warm function %s is nil in repository, error to read GPU memory", funcName)
		return
	}
	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
	if err == nil {
		//log.Printf("scheduler: warm reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
	} else {
		log.Println("scheduler: warm read memory error:", err.Error())
		return
	}

	resourcesConfigs, err := inferResourceConfigsWithBatch(funcName, latencySLO, batchSize,1)
	if err != nil {
		/*log.Print(err.Error())
		wrappedErr := fmt.Errorf("scheduler: CreatePrewarmPod failed batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
			batchSize, funcName, latencySLO, 1, 1)
		log.Println(wrappedErr)*/
		return
	} else {
		/*for _, item := range resourcesConfigs {
			log.Printf("scheduler: warm resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
				funcName, latencySLO, item.ExecutionTime, batchSize, item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
		}*/
	}

	maxThroughputEfficiency := getMaxThroughputEfficiency(funcName)

	cpuConsumedThreadsPerSocket := int(0)
	cpuTotalThreadsPerSocket := int(0)
	cpuOverSell := 0 //CPU threads overSell
	gpuOverSell := 0 //GPU SM percentage
	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
	residualFindFlag := false

	gpuMemConsumedRate := float64(0)
	gpuCoreConsumedRate := float64(0)
	cpuConsumedRate := float64(0)
	slotCpuCapacity := float64(0)
	slotGpuCapacity := float64(0)
	slotUnitCapacity := float64(0)

	tempGpuCoreQuotaRate := float64(0)
	tempGpuMemQuotaRate := float64(0)
	tempCpuQuotaRate := float64(0)

	tempCRE := float64(0)
	maxCRE := float64(-1)
	maxCREConfigIndex := -1
	maxCRENodeIndex := -1
	maxCRESlotIndex := -1

	lock.Lock()
	defer lock.Unlock()
	//log.Println("scheduler: warm locked---------------------------")
	clusterCapConfig := repository.GetClusterCapConfig()

	for i := 0; i < len(clusterCapConfig.ClusterCapacity); i++ { // per node
		cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
		gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
		gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate

		/** CPU GPU consumed rate **/
		cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
		for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
			/**
			 * calculate CPU and GPU core, memory physical resource consumption rate for each slot
			 */
			cpuConsumedThreadsPerSocket = 0
			cpuTotalThreadsPerSocket = 0
			cpuStatus := cpuCapacity[j].CpuStatus
			for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
				cpuConsumedThreadsPerSocket += cpuStatus[k].TotalFuncInstance
				cpuTotalThreadsPerSocket++
			}
			cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
			cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
			cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket+cpuOverSell)                              // cpu usage rate in node i socket j, normalized to 0-1
			gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
			gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate / (1.0 + gpuMemOverSellRate)         //normalized to 0-1

			//log.Println()
			//			//log.Printf("scheduler: warm current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
			//			//	i,
			//			//	j,
			//			//	j+1,
			//			//	float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
			//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
			//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)

			slotCpuCapacity = float64(cpuTotalThreadsPerSocket+cpuOverSell) * 64
			slotGpuCapacity = float64(100+gpuOverSell) * 142
			slotUnitCapacity = slotCpuCapacity + slotGpuCapacity
			for k := 0; k < len(resourcesConfigs); k++ {
				if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
					resourcesConfigs[k].GpuMemoryRate = 0
				} else {
					resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc) / float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
				}
				/**
				 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
				 */
				tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket+cpuOverSell)
				tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100+gpuOverSell)
				tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)

				/**
				 * check if the slot has enough resource
				 */
				targetCpuConsumedRate := cpuConsumedRate + tempCpuQuotaRate
				targetGpuConsumedRate := gpuCoreConsumedRate + tempGpuCoreQuotaRate
				targetGpuMemConsumedRate := gpuMemConsumedRate + tempGpuMemQuotaRate

				if targetCpuConsumedRate > 1.001 ||
					targetGpuConsumedRate > 1.001 ||
					targetGpuMemConsumedRate > 1.001 {
					//log.Printf("scheduler: warm current node has no enough resources for %dth pod config, skip to next pod config\n",k)
					continue
				} else {
					fragmentResource := slotCpuCapacity*(1.001-targetCpuConsumedRate) + slotGpuCapacity*(1.001-targetGpuConsumedRate)
					tempCRE = (float64(resourcesConfigs[k].ReqPerSecondMax) / maxThroughputEfficiency) / (fragmentResource / slotUnitCapacity)
					if tempCRE > maxCRE {
						maxCRE = tempCRE
						maxCRENodeIndex = i
						maxCRESlotIndex = j
						maxCREConfigIndex = k
					}
					//log.Printf("scheduler: warm current node has enough resources for %dth pod config, skip to next pod config\n",k)
				}
			} // per configuration
		} // per socket
	} // per node

	if maxCREConfigIndex == -1 || maxCRESlotIndex == -1 {
		log.Printf("scheduler: error! no pod resource config can be placed in the cluster\n")
	} else {
		/**
		 * allocate CPU threads
		 */
		var cpuCoreThList []int
		cpuStatus := clusterCapConfig.ClusterCapacity[maxCRENodeIndex].CpuCapacity[maxCRESlotIndex].CpuStatus
		neededCores := resourcesConfigs[maxCREConfigIndex].CpuThreads >> 1 //hyper-threads
		for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
			if cpuStatus[k].TotalFuncInstance == 0 {
				cpuCoreThList = append(cpuCoreThList, k)
				neededCores--
			}
		}
		for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
			if cpuStatus[k].TotalFuncInstance != 0 {
				if cpuStatus[k].TotalFuncInstance < CpuFuncInstanceThreshold && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate, CpuUsageRateThreshold) {
					cpuCoreThList = append(cpuCoreThList, k)
					neededCores--
				}
			}
		}

		if neededCores > 0 {
			log.Printf("scheduler: error! warm failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
		} else {
			//log.Printf("scheduler: warm decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
			//	i,
			//	j,
			//	cudaDeviceTh,
			//	cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
			//	cpuTotalThreadsPerSocket + cpuOverSell,
			//	gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
			//	1 + gpuMemOverSellRate,
			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
			//	1 + float64(gpuOverSell)/100)

			// update GPU memory allocation
			cudaDeviceTh := maxCRESlotIndex + 1
			if resourcesConfigs[maxCREConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
				cudaDeviceTh = 0
			}

			resourcesConfigs[maxCREConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation{
				NodeTh:        maxCRENodeIndex,
				SocketTh:      maxCRESlotIndex,
				CudaDeviceTh:  cudaDeviceTh,
				CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
			}

			createErr := createFuncInstance(funcName, namespace, resourcesConfigs[maxCREConfigIndex], "p", clientset)
			if createErr != nil {
				log.Println("scheduler: warm create prewarm function instance failed", createErr.Error())
			} else {
				//log.Printf("scheduler: warm create prewarm function instance for function %s successfully\n", funcName)
				residualFindFlag = true
			}
		}
	}
	//log.Println("scheduler: warm unlocked---------------------------")
	return
}


func ScaleUpFuncCapacity(funcName string, namespace string, latencySLO float64, reqArrivalRate int32, supportBatchGroup []int32, clientset *kubernetes.Clientset) {
	//repository.UpdateFuncIsScalingIn(funcName,true)
	funcStatus := repository.GetFunc(funcName)
	if funcStatus == nil {
		log.Printf("scheduler: function %s is nil in repository, error to scale up", funcName)
		return
	}
	gpuMemAlloc, err := strconv.Atoi(funcStatus.FuncResources.GPU_Memory)
	if err == nil {
		//log.Printf("scheduler: reading GPU memory alloc of function %s = %d\n", funcName, gpuMemAlloc)
	} else {
		log.Println("scheduler: reading memory error:", err.Error())
		return
	}

	maxThroughputEfficiency := getMaxThroughputEfficiency(funcName)

	cpuConsumedThreadsPerSocket := 0
	cpuTotalThreadsPerSocket := 0
	cpuOverSell := 0 //CPU threads overSell
	gpuOverSell := 0 //GPU SM percentage
	gpuMemOverSellRate := float64(0) //GPU memory oversell rate
	residualReq := reqArrivalRate
	residualFindFlag := false
	batchTryNum := 0

	gpuMemConsumedRate := float64(0)
	gpuCoreConsumedRate := float64(0)
	cpuConsumedRate := float64(0)
	slotCpuCapacity := float64(0)
	slotGpuCapacity := float64(0)
	slotUnitCapacity := float64(0)

	tempGpuCoreQuotaRate := float64(0)
	tempGpuMemQuotaRate := float64(0)
	tempCpuQuotaRate := float64(0)

	tempCRE := float64(0)
	maxCRE := float64(-1)
	maxCREConfigIndex := -1
	maxCRENodeIndex := -1
	maxCRESlotIndex := -1

	lock.Lock()
	defer lock.Unlock()
	for {
		if residualReq > 0 {
			residualFindFlag = false
			if batchTryNum >= len(supportBatchGroup) {
				wrappedErr := fmt.Errorf("scheduler: failed to find suitable batchsize for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
					funcName, latencySLO, reqArrivalRate, residualReq)
				log.Println(wrappedErr)
				break
			}
			for batchIndex := 0; batchIndex < len(supportBatchGroup) && residualFindFlag == false; batchIndex++ {
				resourcesConfigs, errInfer := inferResourceConfigsWithBatch(funcName, latencySLO, supportBatchGroup[batchIndex], residualReq)
				if errInfer != nil {
					batchTryNum++
					/*log.Print(errInfer.Error())
					wrappedErr := fmt.Errorf("scheduler: batch=%d cannot meet for function=%s, SLO=%f, reqArrivalRate=%d, residualReq=%d\n",
						supportBatchGroup[batchIndex], funcName, latencySLO, reqArrivalRate, residualReq)
					log.Println(wrappedErr)*/
					continue
				} else {
					/*for _ , item := range resourcesConfigs {
						log.Printf("scheduler: resourcesConfigs={funcName=%s, latencySLO=%f, expectTime=%d, batchSize=%d, cpuThreads=%d, gpuCorePercent=%d, maxCap=%d, minCap=%d}\n",
							funcName, latencySLO, item.ExecutionTime, supportBatchGroup[batchIndex], item.CpuThreads, item.GpuCorePercent, item.ReqPerSecondMax, item.ReqPerSecondMin)
					}*/
				}
				//log.Println("scheduler: warm locked---------------------------")
				clusterCapConfig := repository.GetClusterCapConfig()
				for i := 0; i < len(clusterCapConfig.ClusterCapacity); i++ { // per node
					cpuOverSell = clusterCapConfig.ClusterCapacity[i].CpuCoreOversell
					gpuOverSell = clusterCapConfig.ClusterCapacity[i].GpuCoreOversellPercentage
					gpuMemOverSellRate = clusterCapConfig.ClusterCapacity[i].GpuMemOversellRate

					/** CPU GPU consumed rate **/
					cpuCapacity := clusterCapConfig.ClusterCapacity[i].CpuCapacity
					for j := 0; j < len(cpuCapacity) && residualFindFlag == false; j++ { // per CPU socket (aka per GPU device (j+1))
						/**
						 * calculate CPU and GPU core, memory physical resource consumption rate for each slot
						 */
						cpuConsumedThreadsPerSocket = 0
						cpuTotalThreadsPerSocket = 0
						cpuStatus := cpuCapacity[j].CpuStatus
						for k := 0; k < len(cpuStatus); k++ { // per CPU core in each socket
							cpuConsumedThreadsPerSocket += cpuStatus[k].TotalFuncInstance
							cpuTotalThreadsPerSocket++
						}
						cpuConsumedThreadsPerSocket = cpuConsumedThreadsPerSocket << 1
						cpuTotalThreadsPerSocket = cpuTotalThreadsPerSocket << 1
						cpuConsumedRate = float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket+cpuOverSell)                              // cpu usage rate in node i socket j, normalized to 0-1
						gpuCoreConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate / (1.0 + float64(gpuOverSell)/100) //normalized to 0-1
						gpuMemConsumedRate = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate / (1.0 + gpuMemOverSellRate)         //normalized to 0-1

						//log.Println()
						//			//log.Printf("scheduler: warm current node=%dth, socket=%dth, GPU=%dth, physical CpuConsumedRate=%f, GpuMemConsumedRate=%f, GpuCoreConsumedRate=%f",
						//			//	i,
						//			//	j,
						//			//	j+1,
						//			//	float64(cpuConsumedThreadsPerSocket) / float64(cpuTotalThreadsPerSocket),
						//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemUsageRate,
						//			//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate)

						slotCpuCapacity = float64(cpuTotalThreadsPerSocket+cpuOverSell) * 64
						slotGpuCapacity = float64(100+gpuOverSell) * 142
						slotUnitCapacity = slotCpuCapacity + slotGpuCapacity
						for k := 0; k < len(resourcesConfigs); k++ {
							if resourcesConfigs[k].GpuCorePercent == 0 { //if only CPU are allocated
								resourcesConfigs[k].GpuMemoryRate = 0
							} else {
								resourcesConfigs[k].GpuMemoryRate = float64(gpuMemAlloc) / float64(clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuMemory)
							}
							/**
							 * calculate the CPU, GPU core and GPU memory quota rate of each resource configuration
							 */
							tempCpuQuotaRate = float64(resourcesConfigs[k].CpuThreads) / float64(cpuTotalThreadsPerSocket+cpuOverSell)
							tempGpuCoreQuotaRate = float64(resourcesConfigs[k].GpuCorePercent) / float64(100+gpuOverSell)
							tempGpuMemQuotaRate = resourcesConfigs[k].GpuMemoryRate / (1.0 + gpuMemOverSellRate)

							/**
							 * check if the slot has enough resource
							 */
							targetCpuConsumedRate := cpuConsumedRate + tempCpuQuotaRate
							targetGpuConsumedRate := gpuCoreConsumedRate + tempGpuCoreQuotaRate
							targetGpuMemConsumedRate := gpuMemConsumedRate + tempGpuMemQuotaRate

							if targetCpuConsumedRate > 1.001 ||
								targetGpuConsumedRate > 1.001 ||
								targetGpuMemConsumedRate > 1.001 {
								//log.Printf("scheduler: warm current node has no enough resources for %dth pod config, skip to next pod config\n",k)
								continue
							} else {
								fragmentResource := slotCpuCapacity*(1.001-targetCpuConsumedRate) + slotGpuCapacity*(1.001-targetGpuConsumedRate)
								tempCRE = (float64(resourcesConfigs[k].ReqPerSecondMax) / maxThroughputEfficiency) / (fragmentResource / slotUnitCapacity)
								if tempCRE > maxCRE {
									maxCRE = tempCRE
									maxCRENodeIndex = i
									maxCRESlotIndex = j
									maxCREConfigIndex = k
								}
								//log.Printf("scheduler: warm current node has enough resources for %dth pod config, skip to next pod config\n",k)
							}
						} // per configuration
					} // per socket
				} // per node

				if maxCREConfigIndex == -1 || maxCRESlotIndex == -1 {
					log.Printf("scheduler: error! no pod resource config can be placed in the cluster\n")
				} else {
					/**
					 * allocate CPU threads
					 */
					var cpuCoreThList []int
					cpuStatus := clusterCapConfig.ClusterCapacity[maxCRENodeIndex].CpuCapacity[maxCRESlotIndex].CpuStatus
					neededCores := resourcesConfigs[maxCREConfigIndex].CpuThreads >> 1 //hyper-threads
					for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
						if cpuStatus[k].TotalFuncInstance == 0 {
							cpuCoreThList = append(cpuCoreThList, k)
							neededCores--
						}
					}
					for k := 0; k < len(cpuStatus) && neededCores > 0; k++ {
						if cpuStatus[k].TotalFuncInstance != 0 {
							if cpuStatus[k].TotalFuncInstance < CpuFuncInstanceThreshold && gTypes.LessEqual(cpuStatus[k].TotalCpuUsageRate, CpuUsageRateThreshold) {
								cpuCoreThList = append(cpuCoreThList, k)
								neededCores--
							}
						}
					}

					if neededCores > 0 {
						log.Printf("scheduler: error! warm failed to find enough CPU cores in current socket for residual neededCores=%d", neededCores)
					} else {
						//log.Printf("scheduler: warm decide to schedule pod on node=%dth, socket=%dth, GPU=%dth, physical cpuExpectConsumedThreads=%d (oversell=%d threads), gpuMemExpectConsumedRate=%f (oversell=%f), gpuCoreExpectConsumedRate=%f (oversell=%f)",
						//	i,
						//	j,
						//	cudaDeviceTh,
						//	cpuConsumedThreadsPerSocket + int(resourcesConfigs[pickConfigIndex].CpuThreads),
						//	cpuTotalThreadsPerSocket + cpuOverSell,
						//	gpuMemConsumedRate + resourcesConfigs[pickConfigIndex].GpuMemoryRate,
						//	1 + gpuMemOverSellRate,
						//	clusterCapConfig.ClusterCapacity[i].GpuCapacity[j+1].TotalGpuCoreUsageRate + float64(resourcesConfigs[pickConfigIndex].GpuCorePercent) / 100,
						//	1 + float64(gpuOverSell)/100)

						// update GPU memory allocation
						cudaDeviceTh := maxCRESlotIndex + 1
						if resourcesConfigs[maxCREConfigIndex].GpuCorePercent == 0 { //if only CPU are allocated
							cudaDeviceTh = 0
						}

						resourcesConfigs[maxCREConfigIndex].NodeGpuCpuAllocation = &gTypes.NodeGpuCpuAllocation{
							NodeTh:        maxCRENodeIndex,
							SocketTh:      maxCRESlotIndex,
							CudaDeviceTh:  cudaDeviceTh,
							CpuCoreThList: cpuCoreThList, //no need to check length since cpu must be allocated at least one core
						}

						createErr := createFuncInstance(funcName, namespace, resourcesConfigs[maxCREConfigIndex], "i", clientset)
						if createErr != nil {
							log.Println("scheduler: create function instance failed ", createErr.Error())
							// don't update the residualReq and execute for loop again
						} else {
							//log.Printf("scheduler: create function instance for function%s successfully, residualReq=%d-%d=%d \n",
							//	funcName, residualReq, resourcesConfigs[pickConfigIndex].ReqPerSecondMax, residualReq - resourcesConfigs[pickConfigIndex].ReqPerSecondMax)
							residualReq = residualReq - resourcesConfigs[maxCREConfigIndex].ReqPerSecondMax
						}
						residualFindFlag = true
						batchTryNum = 0
					}
				}
			} // batchsize
		} // residual RPS>0 ?
	} // while true
	return
}



func ScaleDownFuncCapacity(funcName string, namespace string, deletedFuncPodConfig []*gTypes.FuncPodConfig, clientset *kubernetes.Clientset) {
	lock.Lock()
	//log.Println("scheduler: scale down locked---------------------------")
	err := deleteFuncInstance(funcName, namespace, deletedFuncPodConfig, clientset)
	if err != nil {
		log.Println("scheduler: delete function instance failed ", err.Error())
	} else {
		//log.Println("scheduler: delete function instance successfully")
	}
	//log.Println("scheduler: scale down unlocked---------------------------")
	lock.Unlock()
	return
}
