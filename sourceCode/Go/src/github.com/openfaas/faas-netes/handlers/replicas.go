// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas-netes/gpu/repository"
	ntypes "github.com/openfaas/faas-netes/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

/**
 * MakeReplicaUpdater updates desired count of replicas
 * For http calling for gateway cold start
 */


func MakeReplicaUpdater(defaultNamespace string, clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		functionName := vars["name"]
		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		req := ntypes.ScaleServiceRequest{}

		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := "replicas: Cannot parse request. Please pass valid JSON."
				w.Write([]byte(msg))
				log.Println(msg, marshalErr)
				return
			}
		}

		if len(req.ServiceNamespace) > 0 {
			lookupNamespace = req.ServiceNamespace
		}
		funcDeployStatus := repository.GetFunc(functionName)
		if funcDeployStatus == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("replicas: Unable to lookup function deployment " + functionName+lookupNamespace))
			return
		}
		if funcDeployStatus.AvailReplicas == 0 { // no avail ipod to use
			/*minBatchSize := int32(9999)
			var funcConfigWithMinBatch *gTypes.FuncPodConfig
			for _ , funcConfig := range funcDeployStatus.FuncPodConfigMap { //find a ppod to conver into ipod
				if funcConfig.FuncPodType == "p" {
					if funcConfig.BatchSize < minBatchSize {
						minBatchSize = funcConfig.BatchSize
						funcConfigWithMinBatch = funcConfig
					}
				}
			}
			// todo: deal with the ppod
			if funcConfigWithMinBatch != nil {
				funcConfigWithMinBatch.FuncPodType = "i"
				repository.UpdateFuncPodConfig(functionName, funcConfigWithMinBatch)
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfigWithMinBatch.FuncPodName, funcConfigWithMinBatch.ReqPerSecondMax)
				funcConfig, exist := funcDeployStatus.FuncPodConfigMap["v"]
				if exist {
					repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName,0)
				}
			}*/
			for _ , funcConfig := range funcDeployStatus.FuncPodConfigMap { //find a ppod to conver into ipod

				if funcConfig.FuncPodName == funcDeployStatus.FuncPrewarmPodName {
					repository.UpdateFuncExpectedReplicas(functionName, funcDeployStatus.ExpectedReplicas+1)
					//repository.UpdateFuncPrewarmPodName(functionName, funcConfig.FuncPodName)
					repository.UpdateFuncPodType(functionName, funcConfig.FuncPodName,"i")
					repository.UpdateFuncPodLottery(functionName, funcConfig.FuncPodName, funcConfig.ReqPerSecondMax)
					repository.UpdateFuncPodsTotalLotteryNoLog(functionName)
					repository.UpdateFuncAvailReplicas(functionName, funcDeployStatus.AvailReplicas+1)
					break
				}
			}
			funcConfig, exist := funcDeployStatus.FuncPodConfigMap["v"]
			if exist {
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName,0)
			}

			//log.Printf("replicas: conver ppod to ipod of function %s to scale up ......\n", functionName)
		} else {
			log.Printf("replicas: function %s is in scaling up, no need to convert ppod to ipod ......\n", functionName)
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}
}


/** MakeReplicaReader reads the amount of replicas for a deployment
 *  This function is invoked by gateway (scaler.Scale) when request arrivals
 */
func MakeReplicaReader(defaultNamespace string, clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		functionName := vars["name"]
		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		function, err := getService(lookupNamespace, functionName)
		if err != nil {
			log.Printf("replicas: unable to fetch service: %s %s\n", functionName, namespace)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if function == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		//log.Printf("replicas: read replicas - %s %s, %d/%d\n", functionName, lookupNamespace, function.AvailableReplicas, function.Replicas)


		functionBytes, _ := json.Marshal(function)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}

func createSelector(constraints []string) map[string]string {
	selector := make(map[string]string)

	if len(constraints) > 0 {
		for _, constraint := range constraints {
			parts := strings.Split(constraint, "=")
			if len(parts) == 2 {
				selector[parts[0]] = parts[1]
			}
		}
	}

	return selector
}

/**
 * For alert calling, no http
 */
//func ExecReplicaUpdater(namespace string, functionName string, instanceConfig *gpuTypes.FuncPodConfig, requiredReplicas int32, clientset *kubernetes.Clientset) error {
//
//	funcDeployStatus := repository.GetFunc(functionName)
//	repository.UpdateFuncExpectedReplicas(functionName, funcDeployStatus.ExpectedReplicas+1)
//	if funcDeployStatus == nil {
//		log.Println("replicas: Unable to lookup function deployment " + functionName)
//		return fmt.Errorf("replicas: Unable to lookup function deployment " + functionName)
//	}
//
//	differ := funcDeployStatus.ExpectedReplicas - funcDeployStatus.AvailReplicas
//	if differ > 0 {
//		err := scaleUpFunc(funcDeployStatus, namespace, instanceConfig, funcDeployStatus.FuncPlaceConstraints, differ, clientset)
//		if err != nil {
//			log.Println(err.Error())
//			log.Println("replicas: Unable to scaleUp function deployment " + functionName)
//			return err
//		}
//	} else if differ < 0 {
//		err := scaleDownFunc(funcDeployStatus, namespace, -differ, clientset)
//		if err != nil {
//			log.Println(err.Error())
//			log.Println("replicas: Unable to scaleDown function deployment " + functionName)
//			return fmt.Errorf("replicas: Unable to scaleDown function deployment " + functionName)
//		}
//	} else {
//		//log.Println("replicas: ---------expectedReplicas=availReplicas do nothing-----------")
//		// expectedReplicas=availReplicas do nothing
//	}
//	return nil
//}

//func ExecReplicaUpdater(namespace string, functionName string, requiredReplicas int32, clientset *kubernetes.Clientset) error {
//	repository.UpdateFuncExpectedReplicas(functionName, requiredReplicas)
//
//	funcDeployStatus := repository.GetFunc(functionName)
//	if funcDeployStatus == nil {
//		log.Println("replicas: Unable to lookup function deployment " + functionName)
//		return fmt.Errorf("replicas: Unable to lookup function deployment " + functionName)
//	}
//	resourceLimits := &ptypes.FunctionResources {
//		Memory:     funcDeployStatus.FuncResources.Memory,
//		CPU:        funcDeployStatus.FuncResources.CPU,
//		GPU:        funcDeployStatus.FuncResources.GPU,
//		GPU_Memory: funcDeployStatus.FuncResources.GPU_Memory,
//	}
//
//	differ := funcDeployStatus.ExpectedReplicas - funcDeployStatus.AvailReplicas
//	if differ > 0 {
//		err := scaleUpFunc(funcDeployStatus, namespace, resourceLimits, funcDeployStatus.FuncConstraints, differ, clientset)
//		if err != nil {
//			log.Println(err.Error())
//			log.Println("replicas: Unable to scaleUp function deployment " + functionName)
//			return err
//		}
//	} else if differ < 0 {
//		err := scaleDownFunc(funcDeployStatus, namespace, -differ, clientset)
//		if err != nil {
//			log.Println(err.Error())
//			log.Println("replicas: Unable to scaleDown function deployment " + functionName)
//			return fmt.Errorf("replicas: Unable to scaleDown function deployment " + functionName)
//		}
//	} else {
//		//log.Println("replicas: ---------expectedReplicas=availReplicas do nothing-----------")
//		// expectedReplicas=availReplicas do nothing
//	}
//	return nil
//}
//func scaleUpFunc(funcDeployStatus *gpuTypes.FuncDeployStatus, namespace string, resourceLimits *ptypes.FunctionResources, nodeConstrains []string, differ int32, clientset *kubernetes.Clientset) error{
//	nodeSelector := map[string]string{} // init=map{}
//	var nodeGpuAlloc *gpuTypes.NodeGpuCpuAllocation // init=nil
//
//	// there is need to decide the new pod's name and deploy node
//	for i := int32(0); i < differ; i++ {
//		//start := time.Now()
//		nodeGpuAlloc = scheduler.FindGpuDeployNode(resourceLimits, nodeConstrains) // only for GPU and GPU_Memory
//		if nodeGpuAlloc == nil || nodeGpuAlloc.NodeIndex == -1 {
//			log.Println("replicas: no available node in cluster for scale up")
//			return fmt.Errorf("replicas: no available node in cluster for scale up")
//		}
//		// build the node selector
//		nodeLabelStrList := strings.Split(repository.GetClusterCapConfig().ClusterCapacity[nodeGpuAlloc.NodeIndex].NodeLabel, "=")
//		nodeSelector[nodeLabelStrList[0]] = nodeLabelStrList[1]
//		// build the cuda device env str
//		cudaDeviceIndexEnvStr := strconv.Itoa(nodeGpuAlloc.CudaDeviceIndex)
//		if len(funcDeployStatus.FuncSpec.Pod.Spec.Containers)==0 {
//			return fmt.Errorf("replicas: funcSpec.pod.spec.container's length=0 error")
//		}
//		envItemSize := len(funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env)
//		for j := 0; j < envItemSize; j++ {
//			if funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[j].Name == "CUDA_VISIBLE_DEVICES" {
//				funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[j].Value = cudaDeviceIndexEnvStr
//				break
//			}
//		}
//		for j := 0; j < envItemSize; j++ {
//			if funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[j].Name == "GPU_MEM_FRACTION" {
//				funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env[j].Value = funcDeployStatus.FuncResources.GPU_Memory
//				break
//			}
//		}
//
//		funcDeployStatus.FuncSpec.Pod.Name = funcDeployStatus.FunctionName + "-pod-n" + strconv.Itoa(nodeGpuAlloc.NodeIndex) +"g"+ strconv.Itoa(nodeGpuAlloc.CudaDeviceIndex) + "-" +tools.RandomText(8)
//		funcDeployStatus.FuncSpec.Pod.Spec.NodeSelector = nodeSelector
//		_, err := clientset.CoreV1().Pods(namespace).Create(funcDeployStatus.FuncSpec.Pod)
//		if err != nil {
//			wrappedErr := fmt.Errorf("replicas: scaleup function %s 's Pod for differ %d error: %s \n", funcDeployStatus.FunctionName, i+1, err.Error())
//			log.Println(wrappedErr)
//			return err
//		}
//		//log.Printf("replicas: scale function %s took: %fs \n", funcDeployStatus.FunctionName, time.Since(start).Seconds())
//
//		/**
//		 * allocate cpu core
//		 */
//		if funcDeployStatus != nil && funcDeployStatus.FuncCpuCoreBind != ""{
//			coreBindErr := cpuRepository.AssignPodToCpuCore(clientset, funcDeployStatus.FuncSpec.Pod.Name, nodeGpuAlloc.NodeIndex, funcDeployStatus.FuncCpuCoreBind)
//			if coreBindErr != nil {
//				log.Println(coreBindErr.Error())
//				return fmt.Errorf(coreBindErr.Error())
//			}
//			nodeGpuAlloc.CpuCoreIdStr = funcDeployStatus.FuncCpuCoreBind
//		}
//		//log.Printf("replicas: scaleup function %s 's Pod for differ %d successfully \n", funcDeployStatus.FunctionName, i+1)
//		repository.UpdateFuncAvailReplicas(funcDeployStatus.FunctionName, funcDeployStatus.AvailReplicas+1)
//		funcPodConfig := gpuTypes.FuncPodConfig{
//			FuncPodName:          "",
//			BatchSize:            0,
//			CpuThreads:           0,
//			GpuCorePercent:       0,
//			ExecutionTime:        0,
//			ReqPerSecond:         0,
//			FuncPodIp:            "",
//			NodeGpuCpuAllocation: nodeGpuAlloc,
//		}
//		repository.UpdateFuncPodConfig(funcDeployStatus.FunctionName, funcDeployStatus.FuncSpec.Pod.Name, &funcPodConfig)
//	}
//	return nil
//}
//func scaleDownFunc(funcDeployStatus *gpuTypes.FuncDeployStatus, namespace string, differ int32, clientset *kubernetes.Clientset) error{
//	if funcDeployStatus.AvailReplicas < differ {
//		log.Printf("replicas: function %s does not has enough instances %d for differ %d \n", funcDeployStatus.FunctionName, funcDeployStatus.AvailReplicas, differ)
//		return fmt.Errorf("replicas: function %s does not has enough instances %d for differ %d \n", funcDeployStatus.FunctionName, funcDeployStatus.AvailReplicas, differ)
//	}
//	foregroundPolicy := metav1.DeletePropagationForeground
//	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}
//
//	for i := int32(0); i < differ; i++ {
//		//start := time.Now()
//		podName := scheduler.FindGpuDeletePod(funcDeployStatus)
//		err := clientset.CoreV1().Pods(namespace).Delete(podName, opts)
//		if err != nil {
//			log.Printf("replicas: function %s deleted pod %s error \n", funcDeployStatus.FunctionName, podName)
//			return err
//		}
//		//log.Printf("replicas: function %s deleted pod %s successfully \n", funcDeployStatus.FunctionName, podName)
//		repository.UpdateFuncAvailReplicas(funcDeployStatus.FunctionName, funcDeployStatus.AvailReplicas-1)
//		repository.DeleteFuncPodLocation(funcDeployStatus.FunctionName, podName)
//		//log.Printf("replicas: scale function %s took: %fs \n", funcDeployStatus.FunctionName, time.Since(start).Seconds())
//
//	}
//
//	return nil
//}
