package types

import (
	"log"
)

type ClusterCapConfig struct {
	Version string `yaml:"version"`
	ClusterCapacity []*NodeCapacity `yaml:"cluster-capacity"`
}

type NodeCapacity struct {
	NodeLabel string `yaml:"node-label"`
	HyperThreadOffset int `yaml:"hyper-thread-offset"`
	CpuCoreOversell int `yaml:"cpu-thread-oversell"`
	GpuCoreOversellPercentage int `yaml:"gpu-core-oversell-percentage"`
	GpuMemOversellRate float64 `yaml:"gpu-memory-oversell-rate"`
	GpuCapacity []*GpuDevice `yaml:"gpu-capacity"`
	CpuCapacity []*CpuSocket `yaml:"cpu-capacity"`
}
type GpuDevice struct {
	CudaDeviceIndex int `yaml:"cuda-device"`
	TotalGpuMemory int `yaml:"gpu-memory"`
	TotalFuncInstance int `yaml:"total-func-instance"`
	TotalGpuMemUsageRate float64 `yaml:"total-gpu-memory-usage-rate"`
	TotalGpuCoreUsageRate float64 `yaml:"total-gpu-core-usage-rate"`
}
type CpuSocket struct {
	CpuSocketIndex int `yaml:"cpu-socket"`
	CpuStatus []*CpuCore `yaml:"socket-core"`
}
type CpuCore struct {
	CpuCoreIndex int `yaml:"cpu-core-index"`
	TotalFuncInstance int `yaml:"total-func-instance"`
	TotalCpuUsageRate float64 `yaml:"total-cpu-core-usage-rate"`
}
//type ClusterCapConfigLock struct {
//	LockerName string
//	LockState bool
//}
func (clusterCapConfig *ClusterCapConfig) ToString() {
	for _ , node := range clusterCapConfig.ClusterCapacity {
		log.Printf("cluster_config: =================node%s=================\n", node.NodeLabel)
		for i:=0; i< len(node.GpuCapacity); i++ {
			log.Printf("cluster_config: --------------------------------socket%d---------------------------------\n", i)
			log.Printf("cluster_config: GPU %dth device, memoryUsage=%f, coreUsage=%f, instance=%d\n",
				i,
				node.GpuCapacity[i].TotalGpuMemUsageRate,
				node.GpuCapacity[i].TotalGpuCoreUsageRate,
				node.GpuCapacity[i].TotalFuncInstance)
			if i != 0 {
				for _ , cpu := range node.CpuCapacity[i-1].CpuStatus {
					if cpu.TotalFuncInstance > 0 {
						log.Printf("cluster_config: CPU %dth socket, core=%d, TotalFuncInstance=%d, coreUsage=%f\n",
							i-1,
							cpu.CpuCoreIndex,
							cpu.TotalFuncInstance,
							cpu.TotalCpuUsageRate)
					}
				}
			}
			//log.Printf("cluster_config: -----------------------------------------------------------\n")

		}
		log.Printf("cluster_config: ========================================================================\n")

	}
}

