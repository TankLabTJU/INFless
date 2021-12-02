package repository

import (
	"errors"
	"github.com/openfaas/faas-netes/cpu/fifo"
	cpuMetrics "github.com/openfaas/faas-netes/cpu/metrics"
	"github.com/openfaas/faas-netes/cpu/tools"
	cpuTypes "github.com/openfaas/faas-netes/cpu/types"
	gpuMetrics "github.com/openfaas/faas-netes/gpu/metrics"
	gpuRepository "github.com/openfaas/faas-netes/gpu/repository"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var GlobalClusterNodes cpuTypes.ClusterNodes
var GlobalClusterNodesMutex sync.RWMutex

var GlobalClusterNodeFifoMap = make(map[int]*fifo.FIFO)

const Namespace = "openfaasdev"

type CoreBindingRequest struct {
	Podname string
	Coreids string
}

func GetIds(clientset *kubernetes.Clientset,podName string) ([]string,error) {
	container_ids,err := tools.GetPodContainersWithPodName(clientset, Namespace+"-fn", podName)
	if len(container_ids)==0 {
		return []string{},errors.New("Containers are not ready")
	}
	return container_ids,err
}

func DoRequest(req *CoreBindingRequest,container_ids []string, node_id int) {
	_,err := tools.BindCpuCore(req.Coreids, GlobalClusterNodes.Nodes[node_id].AgentIp, container_ids)
	if err != nil {
		log.Printf("cpuRepository: do request failed %s",err.Error())
		//panic(err.Error())
	}
	//log.Println("do request success")
	// fmt.Println(resp)
}

func CoreBindingWorkerThread(clientset *kubernetes.Clientset,fifo *fifo.FIFO,node_id int) {
	log.Printf("core binding worker thread running %d",node_id)
	var process_items []*CoreBindingRequest
	// ticker := time.NewTicker(time.Second)
	buffer_chs := make(chan interface{},10)
	go func() {
		for {
			It,err := fifo.Pop()
			//log.Printf("fifo poped %d",node_id)
			if err != nil {
				log.Println("cpuRepository: ", err.Error())
			} else {
				buffer_chs<-It
			}
		}
	}()
	for {
		select {
		// case <-ticker.C:
		case <-time.After(time.Second):
			if len(process_items)>0 {
				left_items := []*CoreBindingRequest{}
				for _,req := range process_items {
					ids,err := GetIds(clientset,req.Podname)
					if err == nil {
						DoRequest(req,ids,node_id)
					} else {
						left_items = append(left_items,req)
					}
				}
				process_items = left_items
			}
		case It:=<-buffer_chs:
			//log.Println("req poped from chan")
			req := It.(*CoreBindingRequest)
			ids,err := GetIds(clientset,req.Podname)
			if err == nil {
				//log.Println("do request chan")
				DoRequest(req,ids,node_id)
			} else {
				process_items = append(process_items,req)
			}

		}
	}
}

//With Lock
func InitializeCluster(clientset *kubernetes.Clientset) {
	// GetPods(clientset,"openfaas-fn")
	//GlobalClusterNodesMutex.Lock()
	//nodes := tools.GetNodes(clientset)
	nodesConfig := gpuRepository.GetClusterCapConfig()
//	namespace := "openfaasdev"
	var nodeName string
	for i := 0; i < len(nodesConfig.ClusterCapacity); i++ {
		parts := strings.Split(nodesConfig.ClusterCapacity[i].NodeLabel, "=")
		if len(parts) == 2 {
			nodeName = parts[1]
		}

		deployname := "cpuagentcontroller-deploy-"+strconv.Itoa(i)
		servicename := "cpuagentcontroller-service-"+strconv.Itoa(i)
		service_ip := tools.CreateCpuAgentController(clientset, Namespace, deployname,servicename,nodeName)
		log.Printf("cpuAgent controller created in node %s, ip: %s \n", nodeName, service_ip)

		/*deployname = "node-exporter-deploy-"+strconv.Itoa(i)
		servicename = "node-exporter-service-"+strconv.Itoa(i)
		node_exporter_ip := tools.CreateNodeExporter(clientset,"kube-system",deployname,servicename,nodeName)
		log.Printf("node-Exporter created in node %s, ip: %s \n", nodeName, node_exporter_ip)*/

		var node cpuTypes.Node
		node.NodeName = nodeName
		//s_cpu_num := fmt.Sprintf("%s",nd.Status.Capacity.Cpu())
		node.CpuNum = 80
		node.NodeId = i
		node.AgentIp = service_ip
		//node.NodeExporterIp = node_exporter_ip
		for j:=0; j<node.CpuNum; j++ {
			// var cpu_socket CpuSocket
			// cpu_socket.Core.CoreId = j
			// node.CpuSockets = append(node.CpuSockets,cpu_socket)
			var cpu_core cpuTypes.CpuCore
			cpu_core.CoreId = j
			node.CpuCores = append(node.CpuCores,cpu_core)
		}
		GlobalClusterNodes.Nodes = append(GlobalClusterNodes.Nodes,node)
		log.Println("------node information-----", node.ToString())
	}
	// deployname := "cpuagentcontroller-deploy-1"
	// nodename := "jelix-virtual-machine"
	// namespace := "cpuagentcontroller"
	// servicename := "cpuagentcontroller-service-1"

	//initializeFifo(clientset,1)
	//log.Printf("cpuRepository: initialize Fifo, per node thread num=%d \n",1)
	//defer GlobalClusterNodesMutex.Unlock()
}

//Without Lock(because this function is called inside InitializeCluster(...))
func initializeFifo(clientset *kubernetes.Clientset,per_node_thread_num int) {
	num_nodes := len(GlobalClusterNodes.Nodes)
	log.Printf("num nodes in initialize fifo %d",num_nodes)
	for i:=0; i<num_nodes; i++ {
		GlobalClusterNodeFifoMap[i] = fifo.NewFIFO(
			func(obj interface{}) (string,error) {
				req,ok := obj.(*CoreBindingRequest)
				if ok {
					return req.Podname,nil
				} else {
					return "",errors.New("Cannot convert interface{} to *CoreBindingRequest")
				}
			},
		)
	}
	for i:=0; i<num_nodes; i++ {
		for j:=0; j<per_node_thread_num; j++ {
			go CoreBindingWorkerThread(clientset,GlobalClusterNodeFifoMap[i],i)
		}
	}
}

//With Lock
func UpdateClusterCpuIdleRate(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()

	//prometheusip := tools.GetPrometheusServiceIp(clientset)
	prometheusQuery := cpuMetrics.NewPrometheusQuery("prometheus",9090, &http.Client{})

	for node_id,node := range GlobalClusterNodes.Nodes {
		query_string := `rate(node_cpu_seconds_total{mode="idle",kubernetes_pod_name=~"node-exporter-deploy-` + strconv.Itoa(node_id) +`.*"}[1m])`
		// query_string := `rate(node_cpu_seconds_total{mode="idle",instance="` + node.NodeExporterIp + `:9100"}[1m])`
		// fmt.Println(query_string)
		expr := url.QueryEscape(query_string)

		results, fetchErr := prometheusQuery.Fetch(expr)
		if fetchErr != nil {
			// fmt.Println("Error querying Prometheus API: %s\n", fetchErr.Error())
			panic(fetchErr.Error())
		}
		for _, v := range results.Data.Result {
			// fmt.Println(v.Metric.Cpu)
			// fmt.Println(strings.Split(v.Metric.Instance,":")[0])
			// fmt.Println(v.Value[1])
			cpu_core_id,_ := strconv.Atoi(v.Metric.Cpu)
			cpu_idle_rate, _ := strconv.ParseFloat(v.Value[1].(string),32)
			node.CpuCores[cpu_core_id].Idlerate = float32(cpu_idle_rate)
			// fmt.Println(node.CpuCores[cpu_core_id].Idlerate)
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func UpdatePodCpuRate(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()

	//prometheusip := tools.GetPrometheusServiceIp(clientset)
	prometheusQuery := gpuMetrics.NewPrometheusQuery("prometheus",9090, &http.Client{})

	expr := url.QueryEscape(`sum(rate(container_cpu_usage_seconds_total{namespace=~`+Namespace+`}[1m])) by (pod)`)

	results, fetchErr := prometheusQuery.Fetch(expr)
	if fetchErr != nil {
		panic(fetchErr.Error())
	}
	for _, v := range results.Data.Result {
		// podNameSplit := strings.Split(v.Metric.PodName, "-")
		podName := v.Metric.PodName
		podCpuUsage, _ := strconv.ParseFloat(v.Value[1].(string),32)
		// fmt.Println(podName)
		// fmt.Println(podCpuUsage)
		for _,node := range GlobalClusterNodes.Nodes {
			for _,core := range node.CpuCores {
				for _,pd := range core.Pods {
					if pd.PodName == podName {
						pd.CpuUtilizationRate = float32(podCpuUsage)
						// fmt.Println(pd.CpuUtilizationRate)
					}
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}


type PodSortInfo struct {
	PodName string
	CpuUtilizationRate float32
	NewCpuId int
	OldCpuId int
	OldCpuPodId int
}

type PodSortInfos []PodSortInfo

func (a PodSortInfos) Len() int           { return len(a) }
func (a PodSortInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PodSortInfos) Less(i, j int) bool { return a[i].CpuUtilizationRate > a[j].CpuUtilizationRate }


//With Lock
func ReSchedulePods_Resources_Greedy(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()
	cut_off_busy_rate := float32(0.9)
	for _,node := range GlobalClusterNodes.Nodes {
		if len(node.CpuCores) > 1 {
			var pods_sorted PodSortInfos
			for core_id,core := range node.CpuCores {
				for pod_id,pd := range core.Pods {
					var pd_info PodSortInfo
					pd_info.PodName = pd.PodName
					pd_info.CpuUtilizationRate = pd.CpuUtilizationRate
					pd_info.OldCpuId = core_id
					pd_info.OldCpuPodId = pod_id
					pd_info.NewCpuId = -1
					pods_sorted = append(pods_sorted,pd_info)
				}
			}
			if len(pods_sorted) > 1 {
				sort.Sort(pods_sorted)
				var cpu_rate_record map[int]float32
				for i := 0; i < len(node.CpuCores); i++ {
					cpu_rate_record[i] = 0.0
				}
				for _,pd := range pods_sorted {
					for i := 0; i < len(node.CpuCores); i++ {
						if (pd.CpuUtilizationRate + cpu_rate_record[i]) < cut_off_busy_rate {
							pd.NewCpuId = i
							cpu_rate_record[i] += pd.CpuUtilizationRate
							break
						}
					}
				}
				for _,pd := range pods_sorted {
					if pd.NewCpuId == -1 {
						var cur_min_rate_core_id int
						cur_min_rate := float32(-1.0)
						for i := 0; i < len(node.CpuCores); i++ {
							if cur_min_rate < 0 || cpu_rate_record[i] < cur_min_rate {
								cur_min_rate = cpu_rate_record[i]
								cur_min_rate_core_id = i
							}
						}
						pd.NewCpuId = cur_min_rate_core_id
						cpu_rate_record[cur_min_rate_core_id] += pd.CpuUtilizationRate
					}
				}
				var new_cpu_core_pods_record map[int][]cpuTypes.Pod
				for _,pd := range pods_sorted {
					new_cpu_core_pods_record[pd.NewCpuId] = append(new_cpu_core_pods_record[pd.NewCpuId],node.CpuCores[pd.OldCpuId].Pods[pd.OldCpuPodId])
				}
				for core_id,core := range node.CpuCores {
					core.Pods = new_cpu_core_pods_record[core_id]
				}
				for _,pd := range pods_sorted {
					AssignPodToCpuCore(pd.PodName, node.NodeId, strconv.Itoa(pd.NewCpuId))
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func ReSchedulePods_Disturb_Greedy(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()

	for _,node := range GlobalClusterNodes.Nodes {
		if len(node.CpuCores) > 1 {
			var pods_sorted PodSortInfos
			for core_id,core := range node.CpuCores {
				for pod_id,pd := range core.Pods {
					var pd_info PodSortInfo
					pd_info.PodName = pd.PodName
					pd_info.CpuUtilizationRate = pd.CpuUtilizationRate
					pd_info.OldCpuId = core_id
					pd_info.OldCpuPodId = pod_id
					pods_sorted = append(pods_sorted,pd_info)
				}
			}
			if len(pods_sorted) > 1 {
				sort.Sort(pods_sorted)
				var cpu_rate_record map[int]float32
				for i := 0; i < len(node.CpuCores); i++ {
					cpu_rate_record[i] = 0.0
				}
				for _,pd := range pods_sorted {
					var cur_min_rate_core_id int
					cur_min_rate := float32(-1.0)
					for i := 0; i < len(node.CpuCores); i++ {
						if cur_min_rate < 0 || cpu_rate_record[i] < cur_min_rate {
							cur_min_rate = cpu_rate_record[i]
							cur_min_rate_core_id = i
						}
					}
					pd.NewCpuId = cur_min_rate_core_id
					cpu_rate_record[cur_min_rate_core_id] += pd.CpuUtilizationRate
				}
				var new_cpu_core_pods_record map[int][]cpuTypes.Pod
				for _,pd := range pods_sorted {
					new_cpu_core_pods_record[pd.NewCpuId] = append(new_cpu_core_pods_record[pd.NewCpuId],node.CpuCores[pd.OldCpuId].Pods[pd.OldCpuPodId])
				}
				for core_id,core := range node.CpuCores {
					core.Pods = new_cpu_core_pods_record[core_id]
				}
				for _,pd := range pods_sorted {
					AssignPodToCpuCore(pd.PodName, node.NodeId, strconv.Itoa(pd.NewCpuId))
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func ReSchedulePods_MinMax(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()
	cut_off_idle_rate := float32(0.1)
	for _,node := range GlobalClusterNodes.Nodes {
		if len(node.CpuCores) > 1 {
			var cur_max_idle_core_id int
			cur_max_idle_rate := float32(-1.0)
			var cur_min_idle_core_id int
			cur_min_idle_rate := float32(-1.0)
			for core_id,core := range node.CpuCores {
				if cur_max_idle_rate < 0 || core.Idlerate > cur_max_idle_rate {
					cur_max_idle_rate = core.Idlerate
					cur_max_idle_core_id = core_id
				}
				if cur_min_idle_rate < 0 || core.Idlerate < cur_min_idle_rate {
					cur_min_idle_rate = core.Idlerate
					cur_min_idle_core_id = core_id
				}
			}
			if cur_min_idle_rate < cut_off_idle_rate && cur_max_idle_rate > cut_off_idle_rate {
				// MovePodBetweenCores(cur_min_idle_core_id,cur_max_idle_core_id)
				from_core_id := cur_min_idle_core_id
				from_core_idle_rate := cur_min_idle_rate
				to_core_id := cur_max_idle_core_id
				to_core_idle_rate := cur_max_idle_rate
				if len(node.CpuCores[from_core_id].Pods) > 1 {
					var aviliable_pod_ids []int
					for pd_id,pd := range node.CpuCores[from_core_id].Pods {
						if (from_core_idle_rate + pd.CpuUtilizationRate) > cut_off_idle_rate {
							if (to_core_idle_rate - pd.CpuUtilizationRate) > cut_off_idle_rate {
								aviliable_pod_ids = append(aviliable_pod_ids,pd_id)
							}
						}
					}
					if len(aviliable_pod_ids) > 0 {
						var final_pod_id int
						min_rate := float32(-1.0)
						for _,apid := range aviliable_pod_ids {
							if min_rate < 0 || node.CpuCores[from_core_id].Pods[apid].CpuUtilizationRate < min_rate {
								min_rate = node.CpuCores[from_core_id].Pods[apid].CpuUtilizationRate
								final_pod_id = apid
							}
						}
						AssignPodToCpuCore(node.CpuCores[from_core_id].Pods[final_pod_id].PodName,node.NodeId,strconv.Itoa(to_core_id))
					}
				}
			}
		}
	}
	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func AllocateCpu(podname string) (nodename string,cpu_core_id int) {
	GlobalClusterNodesMutex.Lock()

	cur_node_id := -1
	cur_core_id := -1
	cur_max_idle_rate := float32(-1.0)
	for i,node := range GlobalClusterNodes.Nodes {
		for j,core := range node.CpuCores {
			if core.Idlerate > cur_max_idle_rate {
				cur_max_idle_rate = core.Idlerate
				cur_core_id = j
				cur_node_id = i
			}
		}
	}

	pods := &GlobalClusterNodes.Nodes[cur_node_id].CpuCores[cur_core_id].Pods
	var pd cpuTypes.Pod
	pd.PodName = podname
	*pods = append(*pods,pd)

	cur_node_name := GlobalClusterNodes.Nodes[cur_node_id].NodeName

	defer GlobalClusterNodesMutex.Unlock()

	return cur_node_name,cur_core_id
}

//With Lock
func AllocateCpuOnSpecNode(podname string,nodename string) (cpu_core_id int) {
	GlobalClusterNodesMutex.Lock()

	cur_core_id := -1
	cur_max_idle_rate := float32(-1.0)

	for i,node := range GlobalClusterNodes.Nodes {
		if node.NodeName == nodename {
			for j,core := range node.CpuCores {
				if core.Idlerate > cur_max_idle_rate {
					cur_max_idle_rate = core.Idlerate
					cur_core_id = j
				}
			}
			pods := &GlobalClusterNodes.Nodes[i].CpuCores[cur_core_id].Pods
			var pd cpuTypes.Pod
			pd.PodName = podname
			*pods = append(*pods,pd)
			break
		}
	}

	defer GlobalClusterNodesMutex.Unlock()

	return cur_core_id
}

//With Lock
func DeletePod(podname string) {
	GlobalClusterNodesMutex.Lock()

	for _,node := range GlobalClusterNodes.Nodes {
		for _,core := range node.CpuCores {
			for i,pod := range core.Pods {
				if pod.PodName == podname {
					core.Pods = append(core.Pods[:i],core.Pods[i+1:]...)
					return
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

func getNodeIdWithName(nodename string) (nodeid int) {
	for i,node := range GlobalClusterNodes.Nodes {
		if node.NodeName == nodename {
			return i
		}
	}
	return -1
}
/*
func AssignPodToCpuCore(clientSet *kubernetes.Clientset,  podName string, nodeIndex int, cpuCoreIdStr string) error {
	err := GlobalClusterNodeFifoMap[nodeIndex].Add(&CoreBindingRequest{podName,cpuCoreIdStr})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}*/

/*
import (
	"github.com/openfaas/faas-netes/cpu/tools"
	. "github.com/openfaas/faas-netes/cpu/types"
	gpuRepository "github.com/openfaas/faas-netes/gpu/repository"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	// "strings"
	cpuMetrics "github.com/openfaas/faas-netes/cpu/metrics"
	gpuMetrics "github.com/openfaas/faas-netes/gpu/metrics"
	"sort"
)

var GlobalClusterNodes ClusterNodes
//  RLock()  defer RUnlock()  read
//  Lock()   defer UnLock()   write
var GlobalClusterNodesMutex sync.RWMutex


//With Lock
func InitializeCluster(clientset *kubernetes.Clientset) {
	// GetPods(clientset,"openfaas-fn")
	GlobalClusterNodesMutex.Lock()
	//nodes := tools.GetNodes(clientset)
	nodesConfig := gpuRepository.GetClusterCapConfig()
	namespace := "openfaas"
	var nodeName string
	for i := 0; i < len(nodesConfig.ClusterCapacity); i++ {
		parts := strings.Split(nodesConfig.ClusterCapacity[i].NodeLabel, "=")
		if len(parts) == 2 {
			nodeName = parts[1]
		}

		deployname := "cpuagentcontroller-deploy-"+strconv.Itoa(i)
		servicename := "cpuagentcontroller-service-"+strconv.Itoa(i)
		service_ip := tools.CreateCpuAgentController(clientset,namespace,deployname,servicename,nodeName)
		log.Printf("cpuAgent controller created in node %s, ip: %s \n", nodeName, service_ip)

		//deployname = "node-exporter-deploy-"+strconv.Itoa(i)
		//servicename = "node-exporter-service-"+strconv.Itoa(i)
		//node_exporter_ip := tools.CreateNodeExporter(clientset,"kube-system",deployname,servicename,nodeName)
		//log.Printf("node-Exporter created in node %s, ip: %s \n", nodeName, node_exporter_ip)

		var node Node
		node.NodeName = nodeName
		//s_cpu_num := fmt.Sprintf("%s",nd.Status.Capacity.Cpu())
		node.CpuNum = 80
		node.NodeId = i
		node.AgentIp = service_ip
		//node.NodeExporterIp = node_exporter_ip
		for j:=0; j<node.CpuNum; j++ {
			// var cpu_socket CpuSocket
			// cpu_socket.Core.CoreId = j
			// node.CpuSockets = append(node.CpuSockets,cpu_socket)
			var cpu_core CpuCore
			cpu_core.CoreId = j
			node.CpuCores = append(node.CpuCores,cpu_core)
		}
		GlobalClusterNodes.Nodes = append(GlobalClusterNodes.Nodes,node)
		log.Println("------node information-----", node.ToString())
	}
	// deployname := "cpuagentcontroller-deploy-1"
	// nodename := "jelix-virtual-machine"
	// namespace := "cpuagentcontroller"
	// servicename := "cpuagentcontroller-service-1"

	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func UpdateClusterCpuIdleRate(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()

	prometheusip := tools.GetPrometheusServiceIp(clientset)
	prometheusQuery := cpuMetrics.NewPrometheusQuery(prometheusip,9090, &http.Client{})

	for node_id,node := range GlobalClusterNodes.Nodes {
		query_string := `rate(node_cpu_seconds_total{mode="idle",kubernetes_pod_name=~"node-exporter-deploy-` + strconv.Itoa(node_id) +`.*"}[1m])`
		// query_string := `rate(node_cpu_seconds_total{mode="idle",instance="` + node.NodeExporterIp + `:9100"}[1m])`
		// fmt.Println(query_string)
		expr := url.QueryEscape(query_string)

		results, fetchErr := prometheusQuery.Fetch(expr)
		if fetchErr != nil {
			// fmt.Println("Error querying Prometheus API: %s\n", fetchErr.Error())
			panic(fetchErr.Error())
		}
		for _, v := range results.Data.Result {
			// fmt.Println(v.Metric.Cpu)
			// fmt.Println(strings.Split(v.Metric.Instance,":")[0])
			// fmt.Println(v.Value[1])
			cpu_core_id,_ := strconv.Atoi(v.Metric.Cpu)
			cpu_idle_rate, _ := strconv.ParseFloat(v.Value[1].(string),32)
			node.CpuCores[cpu_core_id].Idlerate = float32(cpu_idle_rate)
			// fmt.Println(node.CpuCores[cpu_core_id].Idlerate)
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func UpdatePodCpuRate(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()

	//prometheusip := tools.GetPrometheusServiceIp(clientset)
	prometheusQuery := gpuMetrics.NewPrometheusQuery("prometheus",9090, &http.Client{})

	expr := url.QueryEscape(`sum(rate(container_cpu_usage_seconds_total{namespace=~"openfaas-fn"}[1m])) by (pod)`)

	results, fetchErr := prometheusQuery.Fetch(expr)
	if fetchErr != nil {
		panic(fetchErr.Error())
	}
	for _, v := range results.Data.Result {
		// podNameSplit := strings.Split(v.Metric.PodName, "-")
		podName := v.Metric.PodName
		podCpuUsage, _ := strconv.ParseFloat(v.Value[1].(string),32)
		// fmt.Println(podName)
		// fmt.Println(podCpuUsage)
		for _,node := range GlobalClusterNodes.Nodes {
			for _,core := range node.CpuCores {
				for _,pd := range core.Pods {
					if pd.PodName == podName {
						pd.CpuUtilizationRate = float32(podCpuUsage)
						// fmt.Println(pd.CpuUtilizationRate)
					}
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}


type PodSortInfo struct {
	PodName string
	CpuUtilizationRate float32
	NewCpuId int
	OldCpuId int
	OldCpuPodId int
}

type PodSortInfos []PodSortInfo

func (a PodSortInfos) Len() int           { return len(a) }
func (a PodSortInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PodSortInfos) Less(i, j int) bool { return a[i].CpuUtilizationRate > a[j].CpuUtilizationRate }


//With Lock
func ReSchedulePods_Resources_Greedy(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()
	cut_off_busy_rate := float32(0.9)
	for _,node := range GlobalClusterNodes.Nodes {
		if len(node.CpuCores) > 1 {
			var pods_sorted PodSortInfos
			for core_id,core := range node.CpuCores {
				for pod_id,pd := range core.Pods {
					var pd_info PodSortInfo
					pd_info.PodName = pd.PodName
					pd_info.CpuUtilizationRate = pd.CpuUtilizationRate
					pd_info.OldCpuId = core_id
					pd_info.OldCpuPodId = pod_id
					pd_info.NewCpuId = -1
					pods_sorted = append(pods_sorted,pd_info)
				}
			}
			if len(pods_sorted) > 1 {
				sort.Sort(pods_sorted)
				var cpu_rate_record map[int]float32
				for i := 0; i < len(node.CpuCores); i++ {
					cpu_rate_record[i] = 0.0
				}
				for _,pd := range pods_sorted {
					for i := 0; i < len(node.CpuCores); i++ {
						if (pd.CpuUtilizationRate + cpu_rate_record[i]) < cut_off_busy_rate {
							pd.NewCpuId = i
							cpu_rate_record[i] += pd.CpuUtilizationRate
							break
						}
					}
				}
				for _,pd := range pods_sorted {
					if pd.NewCpuId == -1 {
						var cur_min_rate_core_id int
						cur_min_rate := float32(-1.0)
						for i := 0; i < len(node.CpuCores); i++ {
							if cur_min_rate < 0 || cpu_rate_record[i] < cur_min_rate {
								cur_min_rate = cpu_rate_record[i]
								cur_min_rate_core_id = i
							}
						}
						pd.NewCpuId = cur_min_rate_core_id
						cpu_rate_record[cur_min_rate_core_id] += pd.CpuUtilizationRate
					}
				}
				var new_cpu_core_pods_record map[int][]Pod
				for _,pd := range pods_sorted {
					new_cpu_core_pods_record[pd.NewCpuId] = append(new_cpu_core_pods_record[pd.NewCpuId],node.CpuCores[pd.OldCpuId].Pods[pd.OldCpuPodId])
				}
				for core_id,core := range node.CpuCores {
					core.Pods = new_cpu_core_pods_record[core_id]
				}
				for _,pd := range pods_sorted {
					AssignPodToCpuCore(clientset,pd.PodName, node.NodeId, strconv.Itoa(pd.NewCpuId))
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func ReSchedulePods_Disturb_Greedy(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()

	for _,node := range GlobalClusterNodes.Nodes {
		if len(node.CpuCores) > 1 {
			var pods_sorted PodSortInfos
			for core_id,core := range node.CpuCores {
				for pod_id,pd := range core.Pods {
					var pd_info PodSortInfo
					pd_info.PodName = pd.PodName
					pd_info.CpuUtilizationRate = pd.CpuUtilizationRate
					pd_info.OldCpuId = core_id
					pd_info.OldCpuPodId = pod_id
					pods_sorted = append(pods_sorted,pd_info)
				}
			}
			if len(pods_sorted) > 1 {
				sort.Sort(pods_sorted)
				var cpu_rate_record map[int]float32
				for i := 0; i < len(node.CpuCores); i++ {
					cpu_rate_record[i] = 0.0
				}
				for _,pd := range pods_sorted {
					var cur_min_rate_core_id int
					cur_min_rate := float32(-1.0)
					for i := 0; i < len(node.CpuCores); i++ {
						if cur_min_rate < 0 || cpu_rate_record[i] < cur_min_rate {
							cur_min_rate = cpu_rate_record[i]
							cur_min_rate_core_id = i
						}
					}
					pd.NewCpuId = cur_min_rate_core_id
					cpu_rate_record[cur_min_rate_core_id] += pd.CpuUtilizationRate
				}
				var new_cpu_core_pods_record map[int][]Pod
				for _,pd := range pods_sorted {
					new_cpu_core_pods_record[pd.NewCpuId] = append(new_cpu_core_pods_record[pd.NewCpuId],node.CpuCores[pd.OldCpuId].Pods[pd.OldCpuPodId])
				}
				for core_id,core := range node.CpuCores {
					core.Pods = new_cpu_core_pods_record[core_id]
				}
				for _,pd := range pods_sorted {
					AssignPodToCpuCore(clientset,pd.PodName,node.NodeId, strconv.Itoa(pd.NewCpuId))
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func ReSchedulePods_MinMax(clientset *kubernetes.Clientset) {
	GlobalClusterNodesMutex.Lock()
	cut_off_idle_rate := float32(0.1)
	for _,node := range GlobalClusterNodes.Nodes {
		if len(node.CpuCores) > 1 {
			var cur_max_idle_core_id int
			cur_max_idle_rate := float32(-1.0)
			var cur_min_idle_core_id int
			cur_min_idle_rate := float32(-1.0)
			for core_id,core := range node.CpuCores {
				if cur_max_idle_rate < 0 || core.Idlerate > cur_max_idle_rate {
					cur_max_idle_rate = core.Idlerate
					cur_max_idle_core_id = core_id
				}
				if cur_min_idle_rate < 0 || core.Idlerate < cur_min_idle_rate {
					cur_min_idle_rate = core.Idlerate
					cur_min_idle_core_id = core_id
				}
			}
			if cur_min_idle_rate < cut_off_idle_rate && cur_max_idle_rate > cut_off_idle_rate {
				// MovePodBetweenCores(cur_min_idle_core_id,cur_max_idle_core_id)
				from_core_id := cur_min_idle_core_id
				from_core_idle_rate := cur_min_idle_rate
				to_core_id := cur_max_idle_core_id
				to_core_idle_rate := cur_max_idle_rate
				if len(node.CpuCores[from_core_id].Pods) > 1 {
					var aviliable_pod_ids []int
					for pd_id,pd := range node.CpuCores[from_core_id].Pods {
						if (from_core_idle_rate + pd.CpuUtilizationRate) > cut_off_idle_rate {
							if (to_core_idle_rate - pd.CpuUtilizationRate) > cut_off_idle_rate {
								aviliable_pod_ids = append(aviliable_pod_ids,pd_id)
							}
						}
					}
					if len(aviliable_pod_ids) > 0 {
						var final_pod_id int
						min_rate := float32(-1.0)
						for _,apid := range aviliable_pod_ids {
							if min_rate < 0 || node.CpuCores[from_core_id].Pods[apid].CpuUtilizationRate < min_rate {
								min_rate = node.CpuCores[from_core_id].Pods[apid].CpuUtilizationRate
								final_pod_id = apid
							}
						}
						AssignPodToCpuCore(clientset,node.CpuCores[from_core_id].Pods[final_pod_id].PodName,node.NodeId,strconv.Itoa(to_core_id))
					}
				}
			}
		}
	}
	defer GlobalClusterNodesMutex.Unlock()
}

//With Lock
func AllocateCpu(podname string) (nodename string,cpu_core_id int) {
	GlobalClusterNodesMutex.Lock()

	cur_node_id := -1
	cur_core_id := -1
	cur_max_idle_rate := float32(-1.0)
	for i,node := range GlobalClusterNodes.Nodes {
		for j,core := range node.CpuCores {
			if core.Idlerate > cur_max_idle_rate {
				cur_max_idle_rate = core.Idlerate
				cur_core_id = j
				cur_node_id = i
			}
		}
	}

	pods := &GlobalClusterNodes.Nodes[cur_node_id].CpuCores[cur_core_id].Pods
	var pd Pod
	pd.PodName = podname
	*pods = append(*pods,pd)

	cur_node_name := GlobalClusterNodes.Nodes[cur_node_id].NodeName

	defer GlobalClusterNodesMutex.Unlock()

	return cur_node_name,cur_core_id
}

//With Lock
func AllocateCpuOnSpecNode(podname string,nodename string) (cpu_core_id int) {
	GlobalClusterNodesMutex.Lock()

	cur_core_id := -1
	cur_max_idle_rate := float32(-1.0)

	for i,node := range GlobalClusterNodes.Nodes {
		if node.NodeName == nodename {
			for j,core := range node.CpuCores {
				if core.Idlerate > cur_max_idle_rate {
					cur_max_idle_rate = core.Idlerate
					cur_core_id = j
				}
			}
			pods := &GlobalClusterNodes.Nodes[i].CpuCores[cur_core_id].Pods
			var pd Pod
			pd.PodName = podname
			*pods = append(*pods,pd)
			break
		}
	}

	defer GlobalClusterNodesMutex.Unlock()

	return cur_core_id
}

//With Lock
func DeletePod(podname string) {
	GlobalClusterNodesMutex.Lock()

	for _,node := range GlobalClusterNodes.Nodes {
		for _,core := range node.CpuCores {
			for i,pod := range core.Pods {
				if pod.PodName == podname {
					core.Pods = append(core.Pods[:i],core.Pods[i+1:]...)
					return
				}
			}
		}
	}

	defer GlobalClusterNodesMutex.Unlock()
}

func getNodeIdWithName(nodename string) (nodeid int) {
	for i,node := range GlobalClusterNodes.Nodes {
		if node.NodeName == nodename {
			return i
		}
	}
	return -1
}
*/

/**
 * allocate cpu core for pod using bind
 */

func AssignPodToCpuCoreSync(clientSet *kubernetes.Clientset, namespace string, podName string, nodeIndex int, cpuCoreIdStr string) error {
	containerIds, err := tools.GetPodContainersWithPodName(clientSet, namespace, podName)
	if err != nil {
		log.Println("cpuReposiory:", err.Error())
		return err
	}
	_ , err = tools.BindCpuCore(cpuCoreIdStr, GlobalClusterNodes.Nodes[nodeIndex].AgentIp, containerIds)
	if err != nil {
		log.Println("cpuReposiory:", err.Error())
		return err
	}
	//log.Printf("cpuRepository: cpu bind state %s, for pod %s, container %s, in node %d with IP %s, core %s",
	//	res,
	//	podName,
	//	containerIds[0],
	//	nodeIndex,
	//	GlobalClusterNodes.Nodes[nodeIndex].AgentIp,
	//	cpuCoreIdStr)
	return nil


}


/**
 *
 */
func AssignPodToCpuCore(podName string, nodeTh int, cpuCoreStr string) error {
	err := GlobalClusterNodeFifoMap[nodeTh].Add(&CoreBindingRequest{podName,cpuCoreStr})
	if err != nil {
		log.Println("cpuReposiory:", err.Error())
		//panic(err.Error())
	}
	return err
}
