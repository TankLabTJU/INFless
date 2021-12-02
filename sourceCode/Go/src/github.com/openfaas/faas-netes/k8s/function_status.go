// Copyright 2019 OpenFaaS Authors
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package k8s

import (
	"github.com/openfaas/faas-netes/gpu/repository"
	ptypes "github.com/openfaas/faas-provider/types"
	"log"
)

// EnvProcessName is the name of the env variable containing the function process
const EnvProcessName = "fprocess"

// AsFunctionStatus reads a Deployment object into an OpenFaaS FunctionStatus, parsing the
// Deployment and Container spec into a simplified summary of the Function
/*func AsFunctionStatus(item appsv1.Deployment) *types.FunctionStatus {
	var replicas uint64
	if item.Spec.Replicas != nil {
		replicas = uint64(*item.Spec.Replicas)
	}

	functionContainer := item.Spec.Template.Spec.Containers[0]

	labels := item.Spec.Template.Labels
	function := types.FunctionStatus{
		Name:              item.Name,
		Replicas:          replicas,
		Image:             functionContainer.Image,
		AvailableReplicas: uint64(item.Status.AvailableReplicas),
		InvocationCount:   0,
		Labels:            &labels,
		Annotations:       &item.Spec.Template.Annotations,
		Namespace:         item.Namespace,
	}

	for _, v := range functionContainer.Env {
		if EnvProcessName == v.Name {
			function.EnvProcess = v.Value
		}
	}

	return &function
}*/
// AsFunctionStatus reads a Deployment object into an OpenFaaS FunctionStatus, parsing the
// Deployment and Container spec into a simplified summary of the Function
//func AsFunctionPodStatus(funcName string, item corev1.Pod) *ptypes.FunctionStatus {
//	//funcName := strings.Split(item.Name, "-")[0]
//	funcDeployStatus := repository.GetFunc(funcName)
//	if funcDeployStatus == nil {  // when the pods are in deleting, funcDeployStatus may be nil
//		log.Printf("func_status: pods is not deleted yet, repository has no record for function %s \n", funcName)
//		return nil
//	}
//
//	functionContainer := item.Spec.Containers[0]
//	function := &ptypes.FunctionStatus {
//		Name:              funcName,
//		Replicas:          uint64(funcDeployStatus.ExpectedReplicas),
//		Image:             functionContainer.Image,
//		AvailableReplicas: uint64(funcDeployStatus.AvailReplicas),
//		InvocationCount:   0, //gateway will add this item
//		Labels:            &item.Labels,
//		Annotations:       &item.Annotations,
//		Namespace:         item.Namespace,
//	}
//	for _, v := range functionContainer.Env {
//		if EnvProcessName == v.Name {
//			function.EnvProcess = v.Value
//			break
//		}
//	}
//	return function
//}

/**
 * {"name":"sleep",
	"image":"sleep:latest",
	"invocationCount":0,
	"replicas":2,
	"envProcess":"python index.py","availableReplicas":2,
 "labels":{"com.openfaas.scale.max":"10","com.openfaas.scale.min":"2","com.openfaas.scale.zero":"true","faas_function":"sleep"},
 "annotations":{"prometheus.io.scrape":"false"},"namespace":"openfaas-fn"}
 */
func CreateFunctionPodStatus(functionName string) *ptypes.FunctionStatus {
	funcDeployStatus := repository.GetFunc(functionName)
	if funcDeployStatus == nil {  // when the pods are in deleting, funcDeployStatus may be nil
		log.Printf("func_status: pods is deleted, repository has no record for function %s \n", functionName)
		return nil
	}
	if len(funcDeployStatus.FuncSpec.Pod.Spec.Containers) > 0 {
		function := &ptypes.FunctionStatus {
			Name:              functionName,
			Replicas:          uint64(funcDeployStatus.ExpectedReplicas),
			Image:             funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Image,
			AvailableReplicas: uint64(funcDeployStatus.AvailReplicas),
			InvocationCount:   0,  //gateway will add this item
			Labels:            &funcDeployStatus.FuncSpec.Pod.Labels,
			Annotations:       &funcDeployStatus.FuncSpec.Pod.Annotations,
			Namespace:         funcDeployStatus.FuncSpec.Pod.Namespace,
		}
		for _, v := range 	funcDeployStatus.FuncSpec.Pod.Spec.Containers[0].Env {
			if EnvProcessName == v.Name {
				function.EnvProcess = v.Value
				break
			}
		}
		return function
	} else {
		log.Printf("func_status: function= %s pods's container length is 0\n", functionName)
		return nil
	}
}