// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"github.com/openfaas/faas-netes/gpu/repository"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"net/http"

	"github.com/openfaas/faas/gateway/requests"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// MakeDeleteHandler delete a function
func MakeDeleteHandler(defaultNamespace string, clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		if lookupNamespace == "kube-system" {
			http.Error(w, "delete: unable to list within the kube-system namespace", http.StatusUnauthorized)
			return
		}

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.DeleteFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(request.FunctionName) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// build the label of the pods which will be deleted
		labelPod := labels.SelectorFromSet(map[string]string{"faas_function": request.FunctionName})
		listPodOptions := metav1.ListOptions {
			LabelSelector: labelPod.String(),
		}
		// This makes sure we don't delete non-labeled deployments
		podList, findPodsErr := clientset.CoreV1().Pods(lookupNamespace).List(listPodOptions)
		if findPodsErr != nil {
			if errors.IsNotFound(findPodsErr) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println(findPodsErr.Error())
			w.Write([]byte(findPodsErr.Error()))
			return
		}

		if podList != nil && len(podList.Items) > 0 {
			//log.Printf("delete: find existing %d pods for function: %s \n", len(podList.Items), request.FunctionName)
			deleteFuncErr := deleteFunctionPod(lookupNamespace, listPodOptions, clientset)
			if deleteFuncErr != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("delete: deleting existing pods for function: " + request.FunctionName+ " error"))
				return
			}
		} else { // pod list is nil
			funcDeployStatus := repository.GetFunc(request.FunctionName)
			if funcDeployStatus == nil {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("delete: can't find existing pods for function: " + request.FunctionName+ "may be it has been deleted"))
			} else if funcDeployStatus.AvailReplicas == 0 {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("delete: can't find existing pods for function: " + request.FunctionName+ "may be it has been scaled to 0"))
			} else {
				log.Println("delete: can't find existing pods for function: " + request.FunctionName)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("delete: can't find existing pods for function: " + request.FunctionName))
				return
			}
		}
		// This makes sure we don't delete non-labeled deployments
		listOpts := metav1.ListOptions{}
		srvsList, findServErr := clientset.CoreV1().Services(lookupNamespace).List(listOpts)
		if findServErr != nil {
			if errors.IsNotFound(findServErr) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println(findServErr.Error())
			w.Write([]byte(findServErr.Error()))
			return
		}
		serviceFound := false
		for _, srvItem := range srvsList.Items {
			if srvItem.Name == request.FunctionName {
				serviceFound = true
				break
			}
		}
		if serviceFound == true {
			deleteServErr := deleteFunctionService(lookupNamespace, listPodOptions, clientset, request.FunctionName)
			if deleteServErr != nil {
				log.Println(deleteServErr)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("delete: deleting existing service for function: " + request.FunctionName+ " error"))
				return
			} else {
				//log.Printf("delete: find existing service for function: %s \n", request.FunctionName)
			}
		} else {
			funcDeployStatus := repository.GetFunc(request.FunctionName)
			if funcDeployStatus == nil {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("delete: can't find existing service for function: " + request.FunctionName+ "may be it has been deleted"))
			} else {
				log.Println("delete: can't find existing service for function: " + request.FunctionName)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("delete: can't find existing service for function: " + request.FunctionName))
				return
			}
		}
		repository.DeleteFunc(request.FunctionName)
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func deleteFunctionPod(functionNamespace string, listPodOptions metav1.ListOptions, clientset *kubernetes.Clientset) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	deletePodsErr := clientset.CoreV1().Pods(functionNamespace).DeleteCollection(opts, listPodOptions)
	if deletePodsErr != nil {
		return deletePodsErr
	}
	//log.Printf("delete: delete function %s pods successfully \n", listPodOptions.LabelSelector)
	return nil
}
func deleteFunctionService(functionNamespace string, listPodOptions metav1.ListOptions, clientset *kubernetes.Clientset, functionName string) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions {
		PropagationPolicy: &foregroundPolicy}
	svcErr := clientset.CoreV1().Services(functionNamespace).Delete(functionName, opts)
	if svcErr != nil {
		return svcErr
	}
	//log.Printf("delete: delete function %s service successfully \n", listPodOptions.LabelSelector)
	return nil
}
/*
// MakeDeleteHandler delete a function
func MakeDeleteHandler(defaultNamespace string, clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

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

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.DeleteFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(request.FunctionName) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		getOpts := metav1.GetOptions{}

		// This makes sure we don't delete non-labelled deployments
		deployment, findDeployErr := clientset.AppsV1().Deployments(lookupNamespace).Get(request.FunctionName, getOpts)

		if findDeployErr != nil {
			if errors.IsNotFound(findDeployErr) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			w.Write([]byte(findDeployErr.Error()))
			return
		}

		if isFunction(deployment) {
			err := deleteFunction(lookupNamespace, clientset, request, w)
			if err != nil {
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)

			w.Write([]byte("Not a function: " + request.FunctionName))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func isFunction(deployment *appsv1.Deployment) bool {
	if deployment != nil {
		if _, found := deployment.Labels["faas_function"]; found {
			return true
		}
	}
	return false
}

func deleteFunction(functionNamespace string, clientset *kubernetes.Clientset, request requests.DeleteFunctionRequest, w http.ResponseWriter) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	if deployErr := clientset.AppsV1().Deployments(functionNamespace).
		Delete(request.FunctionName, opts); deployErr != nil {

		if errors.IsNotFound(deployErr) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(deployErr.Error()))
		return fmt.Errorf("error deleting function's deployment")
	}

	if svcErr := clientset.CoreV1().
		Services(functionNamespace).
		Delete(request.FunctionName, opts); svcErr != nil {

		if errors.IsNotFound(svcErr) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(svcErr.Error()))
		return fmt.Errorf("error deleting function's service")
	}
	return nil
}
*/