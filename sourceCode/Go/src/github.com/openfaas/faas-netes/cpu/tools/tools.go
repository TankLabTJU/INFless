package tools

import (
	"encoding/json"
	"fmt"
	"github.com/openfaas/faas-netes/gpu/repository"
	"log"
	"recallsong/httpc"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"strings"

	// "time"

	// "net/http"

)
const Namespace = "openfaasdev"
func int32p(i int32) *int32 {
	return &i
}

func boolp(i bool) *bool {
	return &i
}

func hostTypeP(typename string) *corev1.HostPathType{
	hostType := corev1.HostPathType(typename)
	return &hostType
}

func makeCpuAgentControllerDeploymentSpec(deploymentname string,nodename string) *appsv1.Deployment {
	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deploymentname,
			Labels: map[string]string{
				"CpuAgentController": deploymentname,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"CpuAgentController": deploymentname,
				},
			},
			Replicas: int32p(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        deploymentname,
					Labels:      map[string]string{
						"CpuAgentController": deploymentname,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{"kubernetes.io/hostname":nodename},
					Containers: []corev1.Container{
						{
							Name:  deploymentname,
							Image: "cpuagent/controller:v1",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Protocol: corev1.ProtocolTCP,
								},
							},
							ImagePullPolicy: "IfNotPresent",
							SecurityContext: &corev1.SecurityContext{Privileged: boolp(true)},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "docker-sock-mount",
									MountPath: "/var/run/docker.sock",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "docker-sock-mount",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/run/docker.sock",
									Type: hostTypeP(""),
								},
							},
						},
					},
				},
			},
		},
	}
	return deploymentSpec
}


func makeCpuAgentControllerServiceSpec(deploymentname string,servicename string) *corev1.Service {
	serviceSpec := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        servicename,
			Labels: map[string]string{
				"CpuAgentController": servicename,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"CpuAgentController": deploymentname,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
					},
				},
			},
		},
	}

	return serviceSpec
}

type CpuRequest struct {
	Container_id string
	CPUShares string
	CPUPeriod string
	CPUQuota  string
	CpusetCpus string
	CpusetMems string
}

func CreateCpuAgentController(clientset *kubernetes.Clientset,namespace string,deployname string,servicename string,nodename string) (service_ip string){
	deploylabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"CpuAgentController":deployname},
	}
	deploylistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(deploylabelSelector.MatchLabels).String(),
	}
	deploys, err := clientset.AppsV1().Deployments(namespace).List(deploylistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(deploys.Items) > 0 {
		err = clientset.AppsV1().Deployments(namespace).Delete(deployname,&metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
	}
	deployment := makeCpuAgentControllerDeploymentSpec(deployname,nodename)
	_, err = clientset.AppsV1().Deployments(namespace).Create(deployment)
	if err != nil {
		panic(err.Error())
	}

	servicelabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"CpuAgentController":servicename},
	}
	servicelistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(servicelabelSelector.MatchLabels).String(),
	}
	services, err := clientset.CoreV1().Services(namespace).List(servicelistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(services.Items) > 0 {
		err = clientset.CoreV1().Services(namespace).Delete(servicename,&metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
	}
	service := makeCpuAgentControllerServiceSpec(deployname,servicename)
	new_service,err := clientset.CoreV1().Services(namespace).Create(service)
	if err != nil {
		panic(err.Error())
	}
	return new_service.Spec.ClusterIP
	// service_ip := new_service.Spec.ClusterIP
	// time.Sleep(1000000000)
	// // resp, err := http.Get("http://"+service_ip+":80")
	// // fmt.Println(err,resp)
	// cpu_request := []CpuRequest{
	// 	CpuRequest{
	// 		Container_id : "abefe43fc952320717ce5c1b0fb94f7d9a58e84a8a506e45f5f5e6fb54824444",
	// 		CPUShares : "",
	// 		CPUPeriod : "",
	// 		CPUQuota  : "",
	// 		CpusetCpus : "3",
	// 		CpusetMems : "",
	// 	},
	// 	CpuRequest{
	// 		Container_id : "abefe43fc952320717ce5c1b0fb94f7d9a58e84a8a506e45f5f5e6fb54824444",
	// 		CPUShares : "",
	// 		CPUPeriod : "",
	// 		CPUQuota  : "",
	// 		CpusetCpus : "",
	// 		CpusetMems : "0",
	// 	},
	// }
	// content,_ := json.Marshal(cpu_request)
	// new_content := string(content)
	// var resp string
	// err = httpc.New("http://"+service_ip+":80").
	// 			Path("cpu").
	// 			Query("content", new_content).
	// 			Post(&resp,httpc.TypeApplicationJson)
	// fmt.Println(err,resp)

}

func BindCpuCore(cpuCoreIdStr string, serviceIp string, containers []string) (string,error){
	start := time.Now()
	var cpuRequests []CpuRequest
		for _,id := range containers {
		var single_cpu_request CpuRequest
		single_cpu_request.Container_id = id
		single_cpu_request.CpusetCpus = cpuCoreIdStr
			cpuRequests = append(cpuRequests, single_cpu_request)
	}

	content,_ := json.Marshal(cpuRequests)
	new_content := string(content)
	var resp string
	err := httpc.New("http://"+serviceIp+":80").
		Path("cpu").
		Query("content", new_content).
		Post(&resp,httpc.TypeApplicationJson)
	if err != nil {
		log.Printf("cpuTools: bind cpu core failed took %fs, %s",time.Since(start).Seconds(), err.Error())
	} else {
		//log.Printf("cpuTools: bind cpu core successfully took %fs",time.Since(start).Seconds())
	}
	return resp,err
}

// func makeCpuAgentControllerPodSpec(podname string,nodename string) *corev1.Pod {
// 	pod := &corev1.Pod{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Pod",
// 			APIVersion: "v1",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:   podname,
// 		},
// 		Spec: corev1.PodSpec{
// 			NodeSelector: map[string]string{"hostname":nodename},
// 			Containers: []corev1.Container{
// 				{
// 					Name:  "cpuagentcontroller",
// 					Image: "cpuagentcontroller:v1",
// 					Ports: []corev1.ContainerPort{
// 						{
// 							ContainerPort: 8080,
// 							Protocol: corev1.ProtocolTCP,
// 						},
// 					},
// 					ImagePullPolicy: "IfNotPresent",
// 				},
// 			},
// 		},
// 	}
// 	return pod
// }


// func makeServiceSpec(request ptypes.FunctionDeployment, factory k8s.FunctionFactory) *corev1.Service {
// 	service := &corev1.Service{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Service",
// 			APIVersion: "v1",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:        request.Service,
// 			Annotations: buildAnnotations(request),
// 		},
// 		Spec: corev1.ServiceSpec{
// 			Type: corev1.ServiceTypeClusterIP,
// 			Selector: map[string]string{
// 				"faas_function": request.Service,
// 			},
// 			Ports: []corev1.ServicePort{
// 				{
// 					Name:     "http",
// 					Protocol: corev1.ProtocolTCP,
// 					Port:     factory.Config.RuntimeHTTPPort,
// 					TargetPort: intstr.IntOrString{
// 						Type:   intstr.Int,
// 						IntVal: factory.Config.RuntimeHTTPPort,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	return service
// }

func GetPodsWithFunctionName(clientset *kubernetes.Clientset,namespace string,function_name string)([]corev1.Pod) {
	podlabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"faas_function":function_name},
	}
	podlistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(podlabelSelector.MatchLabels).String(),
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(podlistOpts)
	if err != nil {
		panic(err.Error())
	}
	return pods.Items
}
/*
func GetPodContainersWithPodName(clientset *kubernetes.Clientset,namespace string,podName string) ([]string, error){
	parts := strings.Split(podName, "-")
	fucntionName := parts[0]
	for {
		time.Sleep(time.Millisecond*200)
		getPod, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		//log.Println(len(getPod.Status.ContainerStatuses) == 0, repository.GetFunc(request.Service) != nil)

		if len(getPod.Status.ContainerStatuses) == 0 && repository.GetFunc(fucntionName) != nil {
			continue
		} else {
			if getPod.Status.ContainerStatuses[0].Ready == true {
				break
			}
		}
	}
	pod, err := clientset.CoreV1().Pods(namespace).Get(podName,metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	var containerIds []string
	for _,item := range pod.Status.ContainerStatuses {
		if item.Ready {
			containerIds = append(containerIds,strings.Split(item.ContainerID,"docker://")[1])
		} else {
			return []string{""},fmt.Errorf("cpuTools: Containers are not ready")
		}
	}
	return containerIds, nil
}*/


func GetPodContainersWithPodName(clientset *kubernetes.Clientset, namespace string, podName string) ([]string, error){
	count := 0

	functionName := strings.Split(podName,"-n")[0]
	pod, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		log.Println("tools: ", err.Error())
		return nil, err
	}

	for {
		count++
		if count > 30 { // 60s time out
			break
		}
		if len(pod.Status.ContainerStatuses) == 0 {
			if repository.GetFunc(functionName) == nil {  // func is deleted, then break
				return nil, fmt.Errorf("cpuTools: function is deleted, break cpu bind error")
			}
		} else if pod.Status.ContainerStatuses[0].Ready == true {
			break
		}

		time.Sleep(time.Second*1)
		pod, err = clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			log.Println("tools: ", err.Error())
			return nil, err
		}
		//log.Println(len(getPod.Status.ContainerStatuses) == 0, repository.GetFunc(request.Service) != nil)
	}

	var containerIds []string
	for _,item := range pod.Status.ContainerStatuses {
		if item.Ready {
			containerIds = append(containerIds,strings.Split(item.ContainerID,"docker://")[1])
		} else {
			return []string{""},fmt.Errorf("cpuTools: Containers are not ready, timeout error")
		}
	}
	return containerIds, nil
}

func GetNodes(clientset *kubernetes.Clientset) ([]corev1.Node){
	listOpts := metav1.ListOptions{}
	nodes, err := clientset.CoreV1().Nodes().List(listOpts)
	if err != nil {
		panic(err.Error())
	}
	return nodes.Items
}

func makeNodeExporterDeploymentSpec(deploymentname string,nodename string) *appsv1.Deployment {
	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deploymentname,
			Labels: map[string]string{
				"k8s-app": "node-exporter",
				"node-exporter": deploymentname,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "node-exporter",
				},
			},
			Replicas: int32p(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"prometheus.io.scrape": "true",
						"prometheus.io.port": "9100",
						"prometheus.io.path": "metrics",
					},
					Labels:      map[string]string{
						"k8s-app": "node-exporter",
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{"kubernetes.io/hostname":nodename},
					Containers: []corev1.Container{
						{
							Name:  deploymentname,
							Image: "prom/node-exporter",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9100,
									Protocol: corev1.ProtocolTCP,
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},
	}
	return deploymentSpec
}

func makeNodeExporterServiceSpec(deploymentname string,servicename string) *corev1.Service {
	serviceSpec := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        servicename,
			Labels: map[string]string{
				"k8s-app": "node-exporter",
				"node-exporter": servicename,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"k8s-app": "node-exporter",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     9100,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 31673,
					},
				},
			},
		},
	}

	return serviceSpec
}

func CreateNodeExporter(clientset *kubernetes.Clientset,namespace string,deployname string,servicename string,nodename string) (service_ip string){
	deploylabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"node-exporter":deployname},
	}
	deploylistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(deploylabelSelector.MatchLabels).String(),
	}
	deploys, err := clientset.AppsV1().Deployments(namespace).List(deploylistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(deploys.Items) == 0 {
		deployment := makeNodeExporterDeploymentSpec(deployname,nodename)
		_, err = clientset.AppsV1().Deployments(namespace).Create(deployment)
		if err != nil {
			panic(err.Error())
		}
	}

	servicelabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"node-exporter":servicename},
	}
	servicelistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(servicelabelSelector.MatchLabels).String(),
	}
	services, err := clientset.CoreV1().Services(namespace).List(servicelistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(services.Items) == 0 {
		service := makeNodeExporterServiceSpec(deployname,servicename)
		new_service,err := clientset.CoreV1().Services(namespace).Create(service)
		if err != nil {
			panic(err.Error())
		}
		return new_service.Spec.ClusterIP
	}
	return services.Items[0].Spec.ClusterIP
}

func GetPrometheusServiceIp(clientset *kubernetes.Clientset)  (service_ip string){
	service_name := "prometheus"

	service, err := clientset.CoreV1().Services(Namespace).Get(service_name,metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return service.Spec.ClusterIP
}

/*package tools

import (
	"encoding/json"
	"fmt"
	"github.com/openfaas/faas-netes/gpu/repository"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"recallsong/httpc"
	"strings"

	// "time"

	// "net/http"

)

func int32p(i int32) *int32 {
	return &i
}

func boolp(i bool) *bool {
	return &i
}

func hostTypeP(typename string) *corev1.HostPathType{
	hostType := corev1.HostPathType(typename)
	return &hostType
}

func makeCpuAgentControllerDeploymentSpec(deploymentname string,nodename string) *appsv1.Deployment {
	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deploymentname,
			Labels: map[string]string{
				"CpuAgentController": deploymentname,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"CpuAgentController": deploymentname,
				},
			},
			Replicas: int32p(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        deploymentname,
					Labels:      map[string]string{
						"CpuAgentController": deploymentname,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{"kubernetes.io/hostname":nodename},
					Containers: []corev1.Container{
						{
							Name:  deploymentname,
							Image: "cpuagent/controller:v1",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Protocol: corev1.ProtocolTCP,
								},
							},
							ImagePullPolicy: "IfNotPresent",
							SecurityContext: &corev1.SecurityContext{Privileged: boolp(true)},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "docker-sock-mount",
									MountPath: "/var/run/docker.sock",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "docker-sock-mount",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/run/docker.sock",
									Type: hostTypeP(""),
								},
							},
						},
					},
				},
			},
		},
	}
	return deploymentSpec
}


func makeCpuAgentControllerServiceSpec(deploymentname string,servicename string) *corev1.Service {
	serviceSpec := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        servicename,
			Labels: map[string]string{
				"CpuAgentController": servicename,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"CpuAgentController": deploymentname,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
					},
				},
			},
		},
	}

	return serviceSpec
}

type CpuRequest struct {
	Container_id string
	CPUShares string
	CPUPeriod string
	CPUQuota  string
	CpusetCpus string
	CpusetMems string
}

func CreateCpuAgentController(clientset *kubernetes.Clientset,namespace string,deployname string,servicename string,nodename string) (service_ip string){
	deploylabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"CpuAgentController":deployname},
	}
	deploylistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(deploylabelSelector.MatchLabels).String(),
	}
	deploys, err := clientset.AppsV1().Deployments(namespace).List(deploylistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(deploys.Items) > 0 {
		err = clientset.AppsV1().Deployments(namespace).Delete(deployname,&metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
	}
	deployment := makeCpuAgentControllerDeploymentSpec(deployname,nodename)
	_, err = clientset.AppsV1().Deployments(namespace).Create(deployment)
	if err != nil {
		panic(err.Error())
	}

	servicelabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"CpuAgentController":servicename},
	}
	servicelistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(servicelabelSelector.MatchLabels).String(),
	}
	services, err := clientset.CoreV1().Services(namespace).List(servicelistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(services.Items) > 0 {
		err = clientset.CoreV1().Services(namespace).Delete(servicename,&metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
	}
	service := makeCpuAgentControllerServiceSpec(deployname,servicename)
	new_service,err := clientset.CoreV1().Services(namespace).Create(service)
	if err != nil {
		panic(err.Error())
	}
	return new_service.Spec.ClusterIP
}

func BindCpuCore(cpuCoreIdStr string, serviceIp string,containers []string) (string,error){
	// resp, err := http.Get("http://"+service_ip+":80")
	// fmt.Println(err,resp)
	var cpu_request []CpuRequest
	for _,id := range containers {
		var single_cpu_request CpuRequest
		single_cpu_request.Container_id = id
		single_cpu_request.CpusetCpus = cpuCoreIdStr
		cpu_request = append(cpu_request,single_cpu_request)
	}

	content,_ := json.Marshal(cpu_request)
	new_content := string(content)
	var resp string
	err := httpc.New("http://"+serviceIp+":80").
		Path("cpu").
		Query("content", new_content).
		Post(&resp,httpc.TypeApplicationJson)

	return resp,err
}

func GetPodsWithFunctionName(clientset *kubernetes.Clientset, namespace string, function_name string)([]corev1.Pod) {
	podlabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"faas_function":function_name},
	}
	podlistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(podlabelSelector.MatchLabels).String(),
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(podlistOpts)
	if err != nil {
		panic(err.Error())
	}
	return pods.Items
}

func GetPodContainersWithPodName(clientset *kubernetes.Clientset, namespace string, podName string) ([]string, error){
	count := 0
	fucntionName := strings.Split(podName, "-")[0]

	pod, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for {
		count++
		if count > 300 { // 20s time out
			break
		}
		if len(pod.Status.ContainerStatuses) == 0 {
			if repository.GetFunc(fucntionName) == nil {  // func is deleted, then break
				return nil, fmt.Errorf("cpuTools: function is deleted, break cpu bind error")
			}
		} else if pod.Status.ContainerStatuses[0].Ready == true {
			break
		}

		time.Sleep(time.Millisecond*200)
		pod, err = clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		//log.Println(len(getPod.Status.ContainerStatuses) == 0, repository.GetFunc(request.Service) != nil)
	}

	var containerIds []string
	for _,item := range pod.Status.ContainerStatuses {
		if item.Ready {
			containerIds = append(containerIds,strings.Split(item.ContainerID,"docker://")[1])
		} else {
			return []string{""},fmt.Errorf("cpuTools: Containers are not ready")
		}
	}
	return containerIds, nil
}

func GetNodes(clientset *kubernetes.Clientset) ([]corev1.Node){
	listOpts := metav1.ListOptions{}
	nodes, err := clientset.CoreV1().Nodes().List(listOpts)
	if err != nil {
		panic(err.Error())
	}
	return nodes.Items
}

func makeNodeExporterDeploymentSpec(deploymentname string,nodename string) *appsv1.Deployment {
	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deploymentname,
			Labels: map[string]string{
				"k8s-app": "node-exporter",
				"node-exporter": deploymentname,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "node-exporter",
				},
			},
			Replicas: int32p(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"prometheus.io.scrape": "true",
						"prometheus.io.port": "9100",
						"prometheus.io.path": "metrics",
					},
					Labels:      map[string]string{
						"k8s-app": "node-exporter",
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{"kubernetes.io/hostname":nodename},
					Containers: []corev1.Container{
						{
							Name:  deploymentname,
							Image: "prom/node-exporter",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9100,
									Protocol: corev1.ProtocolTCP,
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},
	}
	return deploymentSpec
}

func makeNodeExporterServiceSpec(deploymentname string,servicename string) *corev1.Service {
	serviceSpec := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        servicename,
			Labels: map[string]string{
				"k8s-app": "node-exporter",
				"node-exporter": servicename,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"k8s-app": "node-exporter",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     9100,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 31673,
					},
				},
			},
		},
	}

	return serviceSpec
}

func CreateNodeExporter(clientset *kubernetes.Clientset,namespace string,deployname string,servicename string,nodename string) (service_ip string){
	deploylabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"node-exporter":deployname},
	}
	deploylistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(deploylabelSelector.MatchLabels).String(),
	}
	deploys, err := clientset.AppsV1().Deployments(namespace).List(deploylistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(deploys.Items) == 0 {
		deployment := makeNodeExporterDeploymentSpec(deployname,nodename)
		_, err = clientset.AppsV1().Deployments(namespace).Create(deployment)
		if err != nil {
			panic(err.Error())
		}
	}

	servicelabelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{"node-exporter":servicename},
	}
	servicelistOpts := metav1.ListOptions{
		LabelSelector: labels.Set(servicelabelSelector.MatchLabels).String(),
	}
	services, err := clientset.CoreV1().Services(namespace).List(servicelistOpts)
	if err != nil {
		panic(err.Error())
	}
	if len(services.Items) == 0 {
		service := makeNodeExporterServiceSpec(deployname,servicename)
		new_service,err := clientset.CoreV1().Services(namespace).Create(service)
		if err != nil {
			panic(err.Error())
		}
		return new_service.Spec.ClusterIP
	}
	return services.Items[0].Spec.ClusterIP
}

func GetPrometheusServiceIp(clientset *kubernetes.Clientset)  (service_ip string){
	service_name := "prometheus"
	namespace := "openfaas"
	service, err := clientset.CoreV1().Services(namespace).Get(service_name,metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return service.Spec.ClusterIP
}*/

