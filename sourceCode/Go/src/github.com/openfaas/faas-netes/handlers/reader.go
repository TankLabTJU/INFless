// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	ptypes "github.com/openfaas/faas-provider/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"

	"github.com/openfaas/faas-netes/k8s"
)

/** MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
 *  This function is invoked by gateway metrics server
 */
func MakeFunctionReader(defaultNamespace string, clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		if lookupNamespace == "kube-system" {
			http.Error(w, "reader: unable to list within the kube-system namespace", http.StatusUnauthorized)
			return
		}

		functions, err := getServiceList(lookupNamespace, clientset)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}

/**
 * invoked by makeFunctionReader
 */
func getServiceList(functionNamespace string, clientset *kubernetes.Clientset) ([]ptypes.FunctionStatus, error) {
	var functions []ptypes.FunctionStatus // init = nil
	var function *ptypes.FunctionStatus // init = nil

	// search service firstly
	listOpts := metav1.ListOptions{}
	srvList, srvErr := clientset.CoreV1().Services(functionNamespace).List(listOpts)
	if srvErr != nil {
		log.Println(srvErr.Error())
		return nil, fmt.Errorf("reader: funtions's services in namespace %s list error \n", functionNamespace)
	}
	for _, srvItem := range srvList.Items {
		function = k8s.CreateFunctionPodStatus(srvItem.Name) // then read repository to get the pod information
		if function == nil {
			log.Printf("reader: error function'pod %s not found in repository while service still exists\n", srvItem.Name)
			return nil, fmt.Errorf("reader: error function'pod %s not found in repository while service still exists\n", srvItem.Name)
		} else {
			//log.Printf("reader: create a func status for function %s in namespace %s from repository, ExpectedReplicas= %d, AvailReplicas= %d \n", functionNamespace, srvItem.Name, function.Replicas, function.AvailableReplicas)
			functions = append(functions, *function)
		}
	}
	//for _, srvItem := range srvList.Items {
	//	// search pod secondly
	//	listOpts.LabelSelector = "faas_function=" + srvItem.Name
	//	podsList, podErr := clientset.CoreV1().Pods(functionNamespace).List(listOpts)
	//	if podErr != nil {
	//		log.Println(podErr.Error())
	//		return nil, fmt.Errorf("reader: funtions's pods in namespace %s list error \n", functionNamespace)
	//	}
	//	if len(podsList.Items) > 0 {
	//		function = k8s.AsFunctionPodStatus(srvItem.Name, podsList.Items[0])
	//	}
	//	if function == nil { //means service exists, but pod has been scaled in to 0
	//		log.Printf("reader: service %s exists in namespace=%s, but pod has been scaled in to 0 \n",srvItem.Name, functionNamespace)
	//		function = k8s.CreateFunctionPodStatus(srvItem.Name) // then read repository to get the pod information
	//		if function == nil {
	//			log.Printf("reader: error function'pod %s not found in repository while service still exists\n", srvItem.Name)
	//			return nil, fmt.Errorf("reader: error function'pod %s not found in repository while service still exists\n", srvItem.Name)
	//		} else {
	//			log.Printf("reader: create a func status for function %s in namespace %s from repository, ExpectedReplicas= %d, AvailReplicas= %d \n", functionNamespace, srvItem.Name, function.Replicas, function.AvailableReplicas)
	//			functions = append(functions, *function)
	//		}
	//	} else {
	//		log.Printf("reader: as a func status for function %s in namespace %s from repository, ExpectedReplicas= %d, AvailReplicas= %d \n", functionNamespace, srvItem.Name, function.Replicas, function.AvailableReplicas)
	//		functions = append(functions, *function)
	//	}
	//}
	return functions, nil
}


/**
 * getService returns a function/service or nil if not found
 * called by MakeReplicaReader()
 */

func getService(functionNamespace string, functionName string) (*ptypes.FunctionStatus, error) {
	function := k8s.CreateFunctionPodStatus(functionName) // then read repository to get the pod information
	if function == nil {
		log.Printf("reader: error function'pod %s not found in repository while service still exists\n", functionName)
		return nil, fmt.Errorf("reader: error function'pod %s not found in repository while service still exists\n", functionName)
	} else {
		//log.Printf("reader: create a func status for function %s in namespace %s from repository, ExpectedReplicas= %d, AvailReplicas= %d \n", functionNamespace, functionName, function.Replicas, function.AvailableReplicas)
		return function, nil // create found
	}
}


//func getService(functionNamespace string, functionName string, clientset *kubernetes.Clientset) (*ptypes.FunctionStatus, error) {
//	var function *ptypes.FunctionStatus  // init = nil
//	funcDeployStatus := repository.GetFunc(functionName)
//	if funcDeployStatus == nil {
//		log.Printf("reader: function %s in namespace %s has been deleted, repository is nil\n", functionName, functionNamespace)
//		return nil, nil
//	}
//	// search service firstly
//	srvs, srvErr := clientset.CoreV1().Services(functionNamespace).Get(functionName,metav1.GetOptions{})
//	if srvErr != nil {
//		log.Println(srvErr.Error())
//		return nil, fmt.Errorf("reader: service of function %s in namespace %s Get error \n", functionName, functionNamespace)
//	}
//	if srvs.Name == functionName { //find this function's service
//		// search pod secondly
//		listOpts := metav1.ListOptions {
//			LabelSelector: "faas_function=" + functionName,
//		}
//		pods, podErr := clientset.CoreV1().Pods(functionNamespace).List(listOpts)
//		if podErr != nil {
//			log.Println(podErr.Error())
//			return nil, fmt.Errorf("reader: function's pod %s in namespace %s list error \n", functionName, functionNamespace)
//		}
//		if len(pods.Items) > 0 {
//			function = k8s.AsFunctionPodStatus(functionName, pods.Items[0])
//		}
//		if function == nil { //means service exists, but pod has been scaled in to 0
//			log.Printf("reader: service %s exists in namespace=%s, but pod has been scaled in to 0 \n", functionName, functionNamespace)
//			function = k8s.CreateFunctionPodStatus(functionName) // then read repository to get the pod information
//			if function == nil {
//				log.Printf("reader: error function'pod %s not found in repository while service still exists\n", functionName)
//				return nil, fmt.Errorf("reader: error function'pod %s not found in repository while service still exists\n", functionName)
//			} else {
//				log.Printf("reader: create a func status for function %s in namespace %s from repository, ExpectedReplicas= %d, AvailReplicas= %d \n", functionNamespace, functionName, function.Replicas, function.AvailableReplicas)
//				return function, nil // create found
//			}
//		} else {
//			return function, nil // pod list found
//		}
//
//	} else {
//		log.Printf("reader: function's service %s not found \n", functionName)
//		return function, fmt.Errorf("reader: function's service %s not found \n", functionName)
//	}
//}
/*
// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	types "github.com/openfaas/faas-provider/types"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openfaas/faas-netes/k8s"
)

// MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
func MakeFunctionReader(defaultNamespace string, clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		if lookupNamespace == "kube-system" {
			http.Error(w, "unable to list within the kube-system namespace", http.StatusUnauthorized)
			return
		}

		functions, err := getServiceList(lookupNamespace, clientset)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}

func getServiceList(functionNamespace string, clientset *kubernetes.Clientset) ([]types.FunctionStatus, error) {
	functions := []types.FunctionStatus{}

	listOpts := metav1.ListOptions{
		LabelSelector: "faas_function",
	}

	res, err := clientset.AppsV1().Deployments(functionNamespace).List(listOpts)

	if err != nil {
		return nil, err
	}

	for _, item := range res.Items {
		function := k8s.AsFunctionStatus(item)
		if function != nil {
			functions = append(functions, *function)
		}
	}
	return functions, nil
}

// getService returns a function/service or nil if not found
func getService(functionNamespace string, functionName string, clientset *kubernetes.Clientset) (*types.FunctionStatus, error) {

	getOpts := metav1.GetOptions{}

	item, err := clientset.AppsV1().Deployments(functionNamespace).Get(functionName, getOpts)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	if item != nil {
		function := k8s.AsFunctionStatus(*item)
		if function != nil {
			return function, nil
		}
	}

	return nil, fmt.Errorf("function: %s not found", functionName)
}

*/