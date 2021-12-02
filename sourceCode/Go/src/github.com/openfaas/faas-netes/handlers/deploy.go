// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	scheduler "github.com/openfaas/faas-netes/gpu/controller"
	"github.com/openfaas/faas-netes/gpu/repository"
	gTypes "github.com/openfaas/faas-netes/gpu/types"
	"github.com/openfaas/faas-netes/k8s"
	ptypes "github.com/openfaas/faas-provider/types"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// initialReplicasCount how many replicas to start of creating for a function
const initialReplicasCount = 1
const initialMaxReplicasCount = 20
const initialScaleToZero = "false"
const initialMaxBatchSize = 1
const initialInactiveNum = 3
const MonitorInterval = time.Second * 10

// MakeDeployHandler creates a handler to create new functions in the cluster
func MakeDeployHandler(functionNamespace string, factory k8s.FunctionFactory, clientset *kubernetes.Clientset) http.HandlerFunc {
	secrets := k8s.NewSecretsClient(factory.Client)

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)

		request := ptypes.FunctionDeployment{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			wrappedErr := fmt.Errorf("deploy: failed to unmarshal request: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		if err := ValidateDeployRequest(&request); err != nil {
			wrappedErr := fmt.Errorf("deploy: validation failed: %s \n", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		namespace := functionNamespace
		if len(request.Namespace) > 0 {
			namespace = request.Namespace
		}

		existingSecrets, err := secrets.GetSecrets(namespace, request.Secrets)
		if err != nil {
			wrappedErr := fmt.Errorf("deploy: unable to fetch secrets: %s \n", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}
		repository.RegisterFuncDeploy(request.Service)
		repository.UpdateFuncProfileCache(&gTypes.FuncProfile {
			FunctionName:    request.Service,
			MaxCpuCoreUsage: 0.9,
			MinCpuCoreUsage: 0.1,
			AvgCpuCoreUsage: rand.Float64(),
		})

		// get the init replicas of function default=1
		initialReplicas := int32p(initialReplicasCount)
		if request.Labels != nil {
			// min
			/*if min := getMinReplicaCount(*request.Labels); min != nil {
				initialReplicas = min
				repository.UpdateFuncMinReplicas(request.Service, *min)
			} else {
				repository.UpdateFuncMinReplicas(request.Service, *initialReplicas)
			}
			// max
			if max := getMaxReplicaCount(*request.Labels); max != nil {
				repository.UpdateFuncMaxReplicas(request.Service, *max)
			} else {
				repository.UpdateFuncMaxReplicas(request.Service, *int32p(initialMaxReplicasCount))
			}*/
			// scale to zero
			if scaleToZero := getScaleToZero(*request.Labels); scaleToZero != nil {
				repository.UpdateFuncScaleToZero(request.Service, *scaleToZero)
			} else {
				repository.UpdateFuncScaleToZero(request.Service, initialScaleToZero)
			}
			// CPU core bind
			//if cpuCoreBind := getCpuCoreBind(*request.Labels); cpuCoreBind != nil {
			//	repository.UpdateFuncCpuCoreBind(request.Service, *cpuCoreBind)
			//} else {
			//	repository.UpdateFuncCpuCoreBind(request.Service, initialCpuCoreBind)
			//}
			// supported batch size
			if inactiveNum := getInactiveNum(*request.Labels); inactiveNum != nil {
				repository.UpdateFuncInactiveNum(request.Service, *inactiveNum)
			} else {
				repository.UpdateFuncInactiveNum(request.Service, initialInactiveNum)
			}
			// supported batch size
			if maxBatchSize := getMaxBatchSize(*request.Labels); maxBatchSize != nil {
				repository.UpdateFuncSupportBatchSize(request.Service, *maxBatchSize)
			} else {
				repository.UpdateFuncSupportBatchSize(request.Service, initialMaxBatchSize)
			}

		}
		qpsPerInstances, err := strconv.ParseFloat(request.QpsPerInstance, 64)
		if err != nil {
			log.Println(err.Error())
			return
		} else {
			repository.UpdateFuncQpsPerInstance(request.Service, qpsPerInstances)
		}

		repository.UpdateFuncRequestResources(request.Service, request.Requests)
		repository.UpdateFuncConstrains(request.Service, request.Constraints)
		repository.UpdateFuncExpectedReplicas(request.Service, 0)
		repository.UpdateFuncAvailReplicas(request.Service, 0)

		//var latestPodSpec *corev1.Pod
		for i := *int32p(0); i < *initialReplicas; i++ {
			podSpec, specErr := makePodSpec(request, existingSecrets, factory)
			if specErr != nil {
				wrappedErr := fmt.Errorf("deploy: failed make Pod spec for replica = %d: %s \n", i+1, specErr.Error())
				log.Println(wrappedErr.Error())
				http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
				return
			}
			//log.Printf("deploy: Deployment (pods with replicas = %d) created: funcName= %s, namespace= %s \n", i+1, request.Service, namespace)

			// after that, deploy the service to find the pods with special label
			serviceSpec := makeServiceSpec(request, factory)
			_, err = factory.Client.CoreV1().Services(namespace).Create(serviceSpec)
			if err != nil {
				wrappedErr := fmt.Errorf("deploy: failed create Service: %s \n", err.Error())
				log.Println(wrappedErr)
				http.Error(w, wrappedErr.Error(), http.StatusInternalServerError)
				return
			}
			//log.Printf("deploy: service created: %s.%s \n", request.Service, namespace)
			repository.UpdateFuncSpec(request.Service, podSpec, serviceSpec)
			repository.AddVirtualFuncPodConfig(request.Service) // add one virtual pod

		}
		go scheduler.CreatePreWarmPod(request.Service, namespace, qpsPerInstances,1, clientset)
		go instanceScaleMonitor(request.Service, functionNamespace, clientset) // Launch a go routine to scaling function instances

		w.WriteHeader(http.StatusAccepted)
		return
	}
}


func makePodSpec(request ptypes.FunctionDeployment, existingSecrets map[string]*corev1.Secret, factory k8s.FunctionFactory) (*corev1.Pod, error) {
	envVars := buildEnvVars(&request)  // prase self-defined environments in faas-cli deploy yaml

	labels := map[string]string {
		"faas_function": request.Service,
	}

	if request.Labels != nil {
		for k, v := range *request.Labels {
			labels[k] = v
		}
	}

	// GPU card selection start
	nodeSelector := map[string]string{} // init=map{}

	// build the node selector
	//nodeLabelStrList := strings.Split(repository.GetClusterCapConfig().
	//	ClusterCapacity[0].NodeLabel, "=")
	//nodeSelector[nodeLabelStrList[0]] = nodeLabelStrList[1]

	envVars = append(envVars, corev1.EnvVar {  // env parameter
		Name:  "CUDA_VISIBLE_DEVICES",
		Value: "-1",
	})
	envVars = append(envVars, corev1.EnvVar { // env parameter
		Name:  "CUDA_MPS_ACTIVE_THREAD_PERCENTAGE",
		Value: "0",
	})

	envVars = append(envVars, corev1.EnvVar { // tfserving exec command parameter
		Name:  "GPU_MEM_FRACTION",
		Value: "0",
	})

	envVars = append(envVars, corev1.EnvVar { // tfserving exec command parameter
		Name:  "BATCH_SIZE",
		Value: "0",
	})

	envVars = append(envVars, corev1.EnvVar { // tfserving exec command parameter
		Name:  "BATCH_TIMEOUT",
		Value: "0",
	})


	//log.Println("deploy: GPU resource envVars = ", envVars)

	resources, resourceErr := createResources(request) // only for CPU and host memory
	if resourceErr != nil {
		return nil, resourceErr
	}

	var imagePullPolicy corev1.PullPolicy
	switch factory.Config.ImagePullPolicy {
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		imagePullPolicy = corev1.PullAlways
	}

	annotations := buildAnnotations(request)
	var serviceAccount string

	if request.Annotations != nil {
		annotations = *request.Annotations
		if val, ok := annotations["com.openfaas.serviceaccount"]; ok && len(val) > 0 {
			serviceAccount = val
		}
	}

	probes, err := factory.MakeProbes(request)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			//sleep2-pod.1.-1.951225,
			Name: request.Service + "-pod-n0g-1-00000000",
			Annotations: annotations, //prometheus.io.scrape: false
			Labels: labels,
			//{labels: com.openfaas.scale.max=15 com.openfaas.scale.min=1 com.openfaas.scale.zero=true
			//faas_function=mnist-test uid=44642818}
		},
		Spec: corev1.PodSpec {
			//HostIPC: true,
			NodeSelector: nodeSelector,
			Containers: []corev1.Container {
				{
					Name:  request.Service + "-con",
					Image: request.Image,
					Ports: []corev1.ContainerPort {
						{
							ContainerPort: factory.Config.RuntimeHTTPPort,
							Protocol: corev1.ProtocolTCP},
					},
					Env:             envVars,
					Resources:       *resources,
					ImagePullPolicy: imagePullPolicy,
					LivenessProbe:   probes.Liveness,
					ReadinessProbe:  probes.Readiness,
					SecurityContext: &corev1.SecurityContext {
						ReadOnlyRootFilesystem: &request.ReadOnlyRootFilesystem,
					},
				},
			},
			ServiceAccountName: serviceAccount,
			RestartPolicy:      corev1.RestartPolicyAlways,
			DNSPolicy:          corev1.DNSClusterFirst,
		},
	}

	factory.ConfigureReadOnlyRootFilesystem(request, pod)
	factory.ConfigureContainerUserID(pod)

	if err = factory.ConfigureSecrets(request, pod, existingSecrets); err != nil {
		return nil, err
	}

	return pod, nil
}


func makeServiceSpec(request ptypes.FunctionDeployment, factory k8s.FunctionFactory) *corev1.Service {
	service := &corev1.Service {
		TypeMeta: metav1.TypeMeta {
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta {
			Name:        request.Service,
			Annotations: buildAnnotations(request),
		},
		Spec: corev1.ServiceSpec {

			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string {
				"faas_function": request.Service,
			},
			Ports: []corev1.ServicePort {
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     factory.Config.RuntimeHTTPPort,
					TargetPort: intstr.IntOrString {
						Type:   intstr.Int,
						IntVal: factory.Config.RuntimeHTTPPort,
					},
				},
			},
		},
	}

	return service
}

func buildAnnotations(request ptypes.FunctionDeployment) map[string]string {
	var annotations map[string]string
	if request.Annotations != nil {
		annotations = *request.Annotations
	} else {
		annotations = map[string]string{}
	}

	annotations["prometheus.io.scrape"] = "false"
	return annotations
}

func buildEnvVars(request *ptypes.FunctionDeployment) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if len(request.EnvProcess) > 0 {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k8s.EnvProcessName,
			Value: request.EnvProcess,
		})
	}

	for k, v := range request.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	sort.SliceStable(envVars, func(i, j int) bool {
		return strings.Compare(envVars[i].Name, envVars[j].Name) == -1
	})

	return envVars
}


func int32p(i int32) *int32 {
	return &i
}
func int64p(i int64) *int64 {
	return &i
}

func createResources(request ptypes.FunctionDeployment) (*corev1.ResourceRequirements, error) {
	resources := &corev1.ResourceRequirements {
		Limits:   corev1.ResourceList{},
		Requests: corev1.ResourceList{},
	}

	// Set Memory limits
	if request.Limits != nil && len(request.Limits.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.Memory)
		if err != nil {
			return resources, err
		}
		resources.Limits[corev1.ResourceMemory] = qty
	}

	if request.Requests != nil && len(request.Requests.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.Memory)
		if err != nil {
			return resources, err
		}
		resources.Requests[corev1.ResourceMemory] = qty
	}

	// Set CPU limits
	if request.Limits != nil && len(request.Limits.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.CPU)
		if err != nil {
			return resources, err
		}
		resources.Limits[corev1.ResourceCPU] = qty
	}

	if request.Requests != nil && len(request.Requests.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.CPU)
		if err != nil {
			return resources, err
		}
		resources.Requests[corev1.ResourceCPU] = qty
	}

	return resources, nil
}

func getMinReplicaCount(labels map[string]string) *int32 {
	if value, exists := labels["com.openfaas.scale.min"]; exists {
		minReplicas, err := strconv.Atoi(value)
		if err == nil && minReplicas > 0 {
			return int32p(int32(minReplicas))
		}
		log.Println(err," minReplicas <= 0")
	}
	return nil
}
func getMaxReplicaCount(labels map[string]string) *int32 {
	if value, exists := labels["com.openfaas.scale.max"]; exists {
		maxReplicas, err := strconv.Atoi(value)
		if err == nil {
			return int32p(int32(maxReplicas))
		}
		log.Println(err)
	}
	return nil
}
func getScaleToZero(labels map[string]string) *string {
	if value, exists := labels["com.openfaas.scale.zero"]; exists {
		scaleToZero := value
		return &scaleToZero
	}
	return nil
}
/*func getCpuCoreBind(labels map[string]string) *string {
	if value, exists := labels["com.openfaas.cpu.bind"]; exists {
		cpuCoreBindStr := strings.ReplaceAll(value,".",",")
		return &cpuCoreBindStr
	}
	return nil
}*/
func getInactiveNum(labels map[string]string) *int32 {
	if value, exists := labels["com.openfaas.inactive.num"]; exists {
		inactiveNum, err := strconv.Atoi(value)
		if err == nil && inactiveNum > 0 {
			return int32p(int32(inactiveNum))
		}
	}
	return nil
}
func getMaxBatchSize(labels map[string]string) *int32 {
	if value, exists := labels["com.openfaas.max.batch"]; exists {
		maxBatchSize, err := strconv.Atoi(value)
		if err == nil && maxBatchSize > 0 {
			return int32p(int32(maxBatchSize))
		}
	}
	return nil
}
func instanceScaleMonitor(funcName string, namespace string, clientset *kubernetes.Clientset) {
	log.Printf("deploy: function scaling go rountine starts namespace=%s, functioName=%s, monitor interval=%ds ......", namespace, funcName, MonitorInterval/1000000000)
	ticker := time.NewTicker(MonitorInterval)
	quit := make(chan struct{})
	writeLock := repository.GetFuncScalingLockState(funcName)
	for {
		select {
		case <-ticker.C:
			funcDeployStatus := repository.GetFunc(funcName)
			if funcDeployStatus == nil {
				log.Printf("deploy: function %s is not in repository, exist scaling in go routine \n", funcName)
				return
			}

			funcConfig, exist := funcDeployStatus.FuncPodConfigMap["v"]
			if exist {
				//if float32(funcConfig.ReqPerSecondLottery) > 0 {

				if funcDeployStatus.FuncRealRps > 0 && float32(funcConfig.ReqPerSecondLottery)/float32(funcDeployStatus.FuncRealRps) > 0.03 {
					//log.Printf("deploy: scale up function pods with vpod lottery= %d, FuncRealRps=%d for function %s ......\n",
					//	funcConfig.ReqPerSecondLottery, funcDeployStatus.FuncRealRps, funcName)
					writeLock.Lock()
					scheduler.ScaleUpFuncCapacity(funcName, namespace, funcDeployStatus.FuncQpsPerInstance, funcConfig.ReqPerSecondLottery, funcDeployStatus.FuncSupportBatchSize, clientset)
					writeLock.Unlock()
					//}
				} else {
					//log.Printf("deploy: no need to scale up and check to scale down function pods with vpod lottery= %d, FuncRealRps=%d for function %s......\n", funcConfig.ReqPerSecondLottery, funcDeployStatus.FuncRealRps, funcName)
					var deletedFuncPodConfig []*gTypes.FuncPodConfig
					for _, v := range funcDeployStatus.FuncPodConfigMap {
						if v.FuncPodType == "i" && v.ReqPerSecondLottery == 0 {
							if v.FuncPodName == funcDeployStatus.FuncPrewarmPodName {
								repository.UpdateFuncExpectedReplicas(funcName, funcDeployStatus.ExpectedReplicas-1)
								repository.UpdateFuncPodType(funcName, v.FuncPodName,"p")
								repository.UpdateFuncAvailReplicas(funcName, funcDeployStatus.AvailReplicas-1)
							} else {
								if v.InactiveCounter >= funcDeployStatus.FunctionInactiveNum {
									deletedFuncPodConfig = append(deletedFuncPodConfig, v)
								} else {
									repository.UpdateFuncPodInactiveCounter(funcName, v.FuncPodName, v.InactiveCounter+1)
								}
							}
						}
					}
					if len(deletedFuncPodConfig) > 0 {
						//if funcDeployStatus.FunctionIsScalingIn == false {
						//log.Printf("deploy: scaling in %d function pods for function %s......\n", len(deletedFuncPodConfig), funcName)

						writeLock.Lock()
						scheduler.ScaleDownFuncCapacity(funcName, namespace, deletedFuncPodConfig, clientset)
						writeLock.Unlock()

						repository.UpdateFuncLastChangedPodCombine(funcName,nil)
						//} else {
						//	log.Printf("deploy: waitting scaling in %d function pods for function %s......\n", len(deletedFuncPodConfig), funcName)
						//}
					} else {
						//log.Printf("deploy: no need to scale in %d function pods for function %s......\n", len(deletedFuncPodConfig), funcName)
					}
				}
			}
			continue
		case <-quit:
			return
		}
	}
	log.Printf("deploy: function scaling go rountine exits --------------------------- \n")
}
/*
func makeDeploymentSpec(request types.FunctionDeployment, existingSecrets map[string]*corev1.Secret, factory k8s.FunctionFactory) (*appsv1.Deployment, error) {
	envVars := buildEnvVars(&request)

	initialReplicas := int32p(initialReplicasCount)
	labels := map[string]string{
		"faas_function": request.Service,
	}

	if request.Labels != nil {
		if min := getMinReplicaCount(*request.Labels); min != nil {
			initialReplicas = min
		}
		for k, v := range *request.Labels {
			labels[k] = v
		}
	}
	// GPU card selection start
	//nodeSelector := make(map[string]string)
	//cudaDeviceIndex := ""
	nodeSelector, cudaDeviceIndex := findGpuNode(request.Limits) // only for GPU and GPU_Memory
	envVars = append(envVars, corev1.EnvVar{
		Name:  "CUDA_VISIBLE_DEVICES",
		Value: cudaDeviceIndex,
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  "GPU_MEM_FRACTION",
		Value: request.Limits.GPU_Memory,
	})
	if request.Constraints != nil && len(request.Constraints) > 0 {
		nodeSelector = createSelector(request.Constraints) // user's defination first
	}
	log.Println("GPU selection: ",envVars)
	// GPU card selection end

	resources, resourceErr := createResources(request) // only for CPU and memory

	if resourceErr != nil {
		return nil, resourceErr
	}

	var imagePullPolicy corev1.PullPolicy
	switch factory.Config.ImagePullPolicy {
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		imagePullPolicy = corev1.PullAlways
	}

	annotations := buildAnnotations(request)


	var serviceAccount string

	if request.Annotations != nil {
		annotations := *request.Annotations
		if val, ok := annotations["com.openfaas.serviceaccount"]; ok && len(val) > 0 {
			serviceAccount = val
		}
	}

	probes, err := factory.MakeProbes(request)
	if err != nil {
		return nil, err
	}

	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service,
			Annotations: annotations,
			Labels: map[string]string{
				"faas_function": request.Service,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"faas_function": request.Service,
				},
			},
			Replicas: initialReplicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			RevisionHistoryLimit: int32p(10),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        request.Service,
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelector,
					Containers: []corev1.Container{
						{
							Name:  request.Service,
							Image: request.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: factory.Config.RuntimeHTTPPort,
									Protocol: corev1.ProtocolTCP},
							},
							Env:             envVars,
							Resources:       *resources,
							ImagePullPolicy: imagePullPolicy,
							LivenessProbe:   probes.Liveness,
							ReadinessProbe:  probes.Readiness,
							SecurityContext: &corev1.SecurityContext{
								ReadOnlyRootFilesystem: &request.ReadOnlyRootFilesystem,
							},
						},
					},
					ServiceAccountName: serviceAccount,
					RestartPolicy:      corev1.RestartPolicyAlways,
					DNSPolicy:          corev1.DNSClusterFirst,
				},
			},
		},
	}

	factory.ConfigureReadOnlyRootFilesystem(request, deploymentSpec)
	factory.ConfigureContainerUserID(deploymentSpec)

	if err := factory.ConfigureSecrets(request, deploymentSpec, existingSecrets); err != nil {
		return nil, err
	}

	return deploymentSpec, nil
}
*/
