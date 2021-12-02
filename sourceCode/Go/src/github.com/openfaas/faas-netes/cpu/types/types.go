package types

import (
	"fmt"
)

type Container struct {
	CPUShares int64 
    CPUPeriod int64    
    CPUQuota  int64                
    CpusetCpus string       
	CpusetMems string     
	ContainerId string
}

type Pod struct {
	Containers []Container
	CpuUtilizationRate float32
	PodName string
}

type CpuCore struct {
	TotalCpuUtilizationRate float32
	PodNums int 
	Pods []Pod 
	CoreId int
	Idlerate float32
}

type CpuSocket struct {
	Core CpuCore
}

type Node struct {
	NodeId int
	NodeName string 
	CpuNum int 
	// CpuSockets []CpuSocket
	CpuCores []CpuCore
	AgentIp string
	NodeExporterIp string
}
func (n Node)ToString() string {
	return fmt.Sprintf("nodeId=%d, nodeName=%s, cpuNum= %d, AgentIp=%s, NodeExporterIp=%s", n.NodeId,n.NodeName,n.CpuNum,n.AgentIp,n.NodeExporterIp)
}

type ClusterNodes struct {
	Nodes []Node
	NodeNum int
}