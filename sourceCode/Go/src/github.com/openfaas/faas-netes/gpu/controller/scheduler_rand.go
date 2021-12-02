//@file: schedulerOld.go
//@author: Yanan Yang
//@date: 2020/11/9
//@note:
package controller

import (
	"fmt"
	"github.com/openfaas/faas-netes/gpu/repository"
	gTypes "github.com/openfaas/faas-netes/gpu/types"
	ptypes "github.com/openfaas/faas-provider/types"
	"log"
	"math/rand"
	"strconv"
)

func FindGpuDeployNode(funcResource *ptypes.FunctionResources, nodeConstrains []string) (nodeGpuAlloc *gTypes.NodeGpuCpuAllocation) {
	// parse resource apply
	var numberGpu = 0
	var GpuMemoryFraction = 0.0

	if funcResource != nil && len(funcResource.GPU) > 0 {
		ng, err := strconv.Atoi(funcResource.GPU)
		if err != nil {
			fmt.Errorf("wrong number of GPU: %s", err.Error())
		}else{
			numberGpu = ng
		}
	}
	if funcResource != nil && len(funcResource.GPU_Memory) > 0 {
		gmf, err := strconv.ParseFloat(funcResource.GPU_Memory, 64)
		if err != nil {
			fmt.Errorf("wrong GPU memory fraction: %s", err.Error())
		}else{
			GpuMemoryFraction = gmf
		}
	}

	// parse self-defined deployment node
	var nodeIndex = -1
	var cudaDeviceIndex = -1
	clusterCapConfig := repository.GetClusterCapConfig()
	if nodeConstrains != nil && len(nodeConstrains) > 0 {
		for i:=0; i<len(clusterCapConfig.ClusterCapacity); i++ {
			if clusterCapConfig.ClusterCapacity[i].NodeLabel == nodeConstrains[0] {
				nodeIndex = i
				break
			}
		}
	}
	if nodeIndex != -1 { //user self-defined node
		if numberGpu == 0 || GpuMemoryFraction == 0.0 { // only use CPU
			nodeGpuAlloc = &gTypes.NodeGpuCpuAllocation {
				NodeTh:       nodeIndex,
				CudaDeviceTh: -1,
			}
			//log.Printf("scheduler: decide to deploy self-defined nodeIndex = %d, cudaDeviceIndex = %d \n", nodeIndex, -1)
			return nodeGpuAlloc
		} else {
			gpuNum := len(clusterCapConfig.ClusterCapacity[nodeIndex].GpuCapacity)
			for j := 0; j < gpuNum; j++ {
				if clusterCapConfig.ClusterCapacity[nodeIndex].GpuCapacity[j].CudaDeviceIndex == -1 {
					continue
				}
				if availMemFraction := 1 - clusterCapConfig.ClusterCapacity[nodeIndex].GpuCapacity[j].TotalGpuMemUsageRate;
					availMemFraction > GpuMemoryFraction {
					cudaDeviceIndex = clusterCapConfig.ClusterCapacity[nodeIndex].GpuCapacity[j].CudaDeviceIndex
					if cudaDeviceIndex == -1 {
						log.Printf("scheduler: find no node to deploy with self-defined nodeIndex = %d, GpuMemoryFraction = %f \n", nodeIndex, GpuMemoryFraction)
					} else {
						//log.Printf("scheduler: decide to deploy self-defined nodeIndex = %d, cudaDeviceIndex = %d \n", nodeIndex, cudaDeviceIndex)
					}
					break
				}
			}

		}
	} else {  //no user self-defined node
		if numberGpu == 0 || GpuMemoryFraction == 0.0 { // only use CPU
			nodeIndex = rand.Intn(len(clusterCapConfig.ClusterCapacity)) // random selection
			nodeGpuAlloc = &gTypes.NodeGpuCpuAllocation {
				NodeTh:       nodeIndex,
				CudaDeviceTh: -1,
			}
			//log.Printf("scheduler: decide to deploy nodeIndex = %d, cudaDeviceIndex = %d \n", nodeIndex, -1)
			return nodeGpuAlloc
		} else {
			numberGpu = 1 // function pod can only be deployed in one GPU in current version
			for i := 0; i < numberGpu; i++ {
				nodeIndex, cudaDeviceIndex = findGpuNode(GpuMemoryFraction, clusterCapConfig)
				if nodeIndex == -1 || cudaDeviceIndex == -1 {
					log.Printf("scheduler: find no node to deploy with numberGpu = %d, GpuMemoryFraction = %f \n", numberGpu, GpuMemoryFraction)
				}else {
					//log.Printf("scheduler: decide to deploy nodeIndex = %d, cudaDeviceIndex = %d \n", nodeIndex, cudaDeviceIndex)
				}
				break
			}
		}
	}
	// return result
	nodeGpuAlloc = &gTypes.NodeGpuCpuAllocation {
		NodeTh:       nodeIndex,
		CudaDeviceTh: cudaDeviceIndex,
	}
	return nodeGpuAlloc
}
func findGpuNode(GpuMemoryFraction float64, clusterCapConfig *gTypes.ClusterCapConfig) (nodeIndex int, cudaDeviceIndex int) {
	nodeIndex = -1
	cudaDeviceIndex = -1
	clusterNum := len(clusterCapConfig.ClusterCapacity)
	for i := 0; i < clusterNum; i++ {
		gpuNum := len(clusterCapConfig.ClusterCapacity[i].GpuCapacity)
		for j := 0; j < gpuNum; j++ {
			if clusterCapConfig.ClusterCapacity[i].GpuCapacity[j].CudaDeviceIndex == -1 {
				continue
			}
			if availMemFraction := 1 - clusterCapConfig.ClusterCapacity[i].GpuCapacity[j].TotalGpuMemUsageRate;
				availMemFraction > GpuMemoryFraction {
				nodeIndex = i
				cudaDeviceIndex = clusterCapConfig.ClusterCapacity[i].GpuCapacity[j].CudaDeviceIndex
				break
			}
		}
	}
	return nodeIndex, cudaDeviceIndex
}

/**
 * @author yanan
 * @desc public function
 * @date 2020/5/6
 * @param
 * @return
 **/
func FindGpuDeletePod(funcDeployStatus *gTypes.FuncDeployStatus) (podName string) {
	//clusterCapConfig := repository.GetClusterCapConfig()
	return findGpuDeletePod(funcDeployStatus)
}
func findGpuDeletePod(funcDeployStatus *gTypes.FuncDeployStatus) (candidatePodName string) {
	score := 0
	for k, v := range funcDeployStatus.FuncPodConfigMap {
		if temp := (v.NodeGpuCpuAllocation.NodeTh + 1)*100 + v.NodeGpuCpuAllocation.CudaDeviceTh; temp > score {
			score = temp
			candidatePodName = k
		}
	}
	//log.Printf("scheduler: decide to delete podName = %s, nodeIndex = %d, cudaDeviceIndex = %d \n", candidatePodName, funcDeployStatus.FuncPodLocationMap[candidatePodName].NodeIndex, funcDeployStatus.FuncPodLocationMap[candidatePodName].CudaDeviceIndex)
	return candidatePodName
}
func findCpuCore(cpuSocket []gTypes.CpuSocket, cpuCoreUsage float64) (cpuCoreIndex int){
	/*
	   TO DO
	*/

	return cpuCoreIndex
}