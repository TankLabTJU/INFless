//File  : func_config.go
//Author: Yanan Yang
//Date  : 2020/4/7
package types

import (
	ptypes "github.com/openfaas/faas-provider/types"
	corev1 "k8s.io/api/core/v1"
	"log"
	"sync"
)

type FuncDeployStatus struct {
	FunctionName string
	FuncQpsPerInstance float64 //SLO
	FunctionScalingLock *sync.RWMutex
	FunctionSortedPodLock *sync.RWMutex
	FuncSpec *FuncSpec
	ExpectedReplicas int32
	AvailReplicas int32
	MinReplicas int32
	MaxReplicas int32
	ScaleToZero string
	FunctionInactiveNum int32

	FuncResources *ptypes.FunctionResources
	FuncPlaceConstraints []string

	FuncPrewarmPodName string
	FuncPodConfigMap map[string]*FuncPodConfig // key:podName value:InstanceConfig{}
	FuncPodMaxCapacity int32  // the max request processing capacity of function pods
	FuncPodMinCapacity int32  // the min request processing capacity of function pods
	FuncPodTotalLottery int32
	FuncRealRps int32
	FuncLastRealRps int32
	FuncSupportBatchSize []int32
	FuncLastChangedPodCombine *ChangedPodCombine
	FuncSortedPodNameList []string
}
type FuncSpec struct {
	Pod *corev1.Pod
	Service *corev1.Service
}
type FuncPodConfig struct {
	FuncPodType    string
	FuncPodName    string
	BatchSize      int32
	CpuThreads     int32
	GpuCorePercent int32
	GpuMemoryRate  float64
	ExecutionTime  int32
	BatchTimeOut int32
	ReqPerSecondLottery int32
	ReqPerSecondMax  int32
	ReqPerSecondMin  int32
	InactiveCounter int32
	FuncPodIp string  // end-point ip
	NodeGpuCpuAllocation *NodeGpuCpuAllocation
}
type NodeGpuCpuAllocation struct {
	NodeTh int
	CudaDeviceTh int
	SocketTh int
	CpuCoreThList []int
}

type FuncProfile struct {
	FunctionName string
	MaxCpuCoreUsage float64
	MinCpuCoreUsage float64
	AvgCpuCoreUsage float64
}

type ChangedPodCombine struct {
	DeletedPodCount int32
	MinDeletedSumCap int32
	MaxDeletedSumCap int32
	DeletePodNameList []string
	RemainPodNameList []string
}
func (funcPodConfig *FuncPodConfig) ToString() {
	log.Printf("FuncPodName = %s, FuncPodType= %s, ReqPerSecondLottery= %d\n",
		funcPodConfig.FuncPodName,
		funcPodConfig.FuncPodType,
		funcPodConfig.ReqPerSecondLottery)
}

func (funcDeployStatus *FuncDeployStatus) ToString() {
	log.Printf("Function = %s, ExpectedReplicas = %d, AvailReplicas = %d,", funcDeployStatus.FunctionName, funcDeployStatus.ExpectedReplicas, funcDeployStatus.AvailReplicas)
	for podName, podConfig := range funcDeployStatus.FuncPodConfigMap {
		log.Println()
		log.Println(" pod:", podName, "--> node: ", podConfig)
	}
}
