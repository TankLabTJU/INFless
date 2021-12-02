package controller

import (
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// corev1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/api/errors"
	// . "github.com/openfaas/faas-netes/cpu/types"
	// "fmt"
	// "log"
	// "os"
	"time"
	// "strconv"
	"k8s.io/client-go/kubernetes"
	// "github.com/openfaas/faas-netes/cpu/tools"	
	"github.com/openfaas/faas-netes/cpu/repository"	
	// "encoding/json"
	// "sync"
)

func Initialize(clientset *kubernetes.Clientset) {
	repository.InitializeCluster(clientset)
	repository.UpdateClusterCpuIdleRate(clientset)
	repository.UpdatePodCpuRate(clientset)
}

func ReScheduler(clientset *kubernetes.Clientset) {
	for {
		time.Sleep(time.Minute)
		// time.Sleep(time.Second)
		repository.UpdateClusterCpuIdleRate(clientset)
		repository.UpdatePodCpuRate(clientset)
		repository.ReSchedulePods_MinMax(clientset)
		// repository.ReSchedulePods_Disturb_Greedy(clientset)
		// repository.ReSchedulePods_Resources_Greedy(clientset)
	}
}

func Scheduler(podname string) (nodename string,cpu_core_id int) {
	return repository.AllocateCpu(podname)
}

func SchedulerGPU(podname string,nodename string) (cpu_core_id int) {
	return repository.AllocateCpuOnSpecNode(podname,nodename)
}

func StartScheduler(clientset *kubernetes.Clientset) {
	Initialize(clientset)
	go ReScheduler(clientset)
}

// time.Sleep(1000000000)
// pods := tools.GetPodsWithFunctionName(clientset,"openfaas-fn","hello-python")
// for _,pds := range pods {
// 	// fmt.Printf(string(json.Marshal(pds.Spec.Containers)))
// 	fmt.Printf("%s\n",string(pds.Name))
// 	container_ids,err := tools.GetPodContainersWithPodName(clientset,"openfaas-fn",pds.Name)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	for _,id := range container_ids {
// 		fmt.Printf("%s\n",id)
// 	}
// 	resp,err := tools.BindCpuCore(0,GlobalClusterNodes.Nodes[0].AgentIp,container_ids)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println(resp)
// }


// {
// 	"metadata": {
// 		"name": "hello-python-869cb68fc9-ck929",
// 		"generateName": "hello-python-869cb68fc9-",
// 		"namespace": "openfaas-fn",
// 		"selfLink": "/api/v1/namespaces/openfaas-fn/pods/hello-python-869cb68fc9-ck929",
// 		"uid": "c523b383-e71b-4792-8a01-836b102345f5",
// 		"resourceVersion": "355194",
// 		"creationTimestamp": "2020-04-26T04:11:11Z",
// 		"labels": {
// 			"faas_function": "hello-python",
// 			"pod-template-hash": "869cb68fc9"
// 		},
// 		"annotations": {
// 			"prometheus.io.scrape": "false"
// 		},
// 		"ownerReferences": [{
// 			"apiVersion": "apps/v1",
// 			"kind": "ReplicaSet",
// 			"name": "hello-python-869cb68fc9",
// 			"uid": "61901ff9-c389-45fb-88a8-fbc29f296bed",
// 			"controller": true,
// 			"blockOwnerDeletion": true
// 		}]
// 	},
// 	"spec": {
// 		"volumes": [{
// 			"name": "default-token-4lc5q",
// 			"secret": {
// 				"secretName": "default-token-4lc5q",
// 				"defaultMode": 420
// 			}
// 		}],
// 		"containers": [{
// 			"name": "hello-python",
// 			"image": "hello-python:latest",
// 			"ports": [{
// 				"containerPort": 8080,
// 				"protocol": "TCP"
// 			}],
// 			"env": [{
// 				"name": "fprocess",
// 				"value": "python index.py"
// 			}],
// 			"resources": {},
// 			"volumeMounts": [{
// 				"name": "default-token-4lc5q",
// 				"readOnly": true,
// 				"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
// 			}],
// 			"livenessProbe": {
// 				"httpGet": {
// 					"path": "/_/health",
// 					"port": 8080,
// 					"scheme": "HTTP"
// 				},
// 				"initialDelaySeconds": 2,
// 				"timeoutSeconds": 1,
// 				"periodSeconds": 2,
// 				"successThreshold": 1,
// 				"failureThreshold": 3
// 			},
// 			"readinessProbe": {
// 				"httpGet": {
// 					"path": "/_/health",
// 					"port": 8080,
// 					"scheme": "HTTP"
// 				},
// 				"initialDelaySeconds": 2,
// 				"timeoutSeconds": 1,
// 				"periodSeconds": 2,
// 				"successThreshold": 1,
// 				"failureThreshold": 3
// 			},
// 			"terminationMessagePath": "/dev/termination-log",
// 			"terminationMessagePolicy": "File",
// 			"imagePullPolicy": "IfNotPresent",
// 			"securityContext": {
// 				"readOnlyRootFilesystem": false
// 			}
// 		}],
// 		"restartPolicy": "Always",
// 		"terminationGracePeriodSeconds": 30,
// 		"dnsPolicy": "ClusterFirst",
// 		"serviceAccountName": "default",
// 		"serviceAccount": "default",
// 		"nodeName": "jelix-virtual-machine",
// 		"securityContext": {},
// 		"schedulerName": "default-scheduler",
// 		"tolerations": [{
// 			"key": "node.kubernetes.io/not-ready",
// 			"operator": "Exists",
// 			"effect": "NoExecute",
// 			"tolerationSeconds": 300
// 		}, {
// 			"key": "node.kubernetes.io/unreachable",
// 			"operator": "Exists",
// 			"effect": "NoExecute",
// 			"tolerationSeconds": 300
// 		}],
// 		"priority": 0,
// 		"enableServiceLinks": true
// 	},
// 	"status": {
// 		"phase": "Running",
// 		"conditions": [{
// 			"type": "Initialized",
// 			"status": "True",
// 			"lastProbeTime": null,
// 			"lastTransitionTime": "2020-04-26T04:11:11Z"
// 		}, {
// 			"type": "Ready",
// 			"status": "True",
// 			"lastProbeTime": null,
// 			"lastTransitionTime": "2020-06-06T11:31:45Z"
// 		}, {
// 			"type": "ContainersReady",
// 			"status": "True",
// 			"lastProbeTime": null,
// 			"lastTransitionTime": "2020-06-06T11:31:45Z"
// 		}, {
// 			"type": "PodScheduled",
// 			"status": "True",
// 			"lastProbeTime": null,
// 			"lastTransitionTime": "2020-04-26T04:11:11Z"
// 		}],
// 		"hostIP": "192.168.183.147",
// 		"podIP": "10.42.0.155",
// 		"podIPs": [{
// 			"ip": "10.42.0.155"
// 		}],
// 		"startTime": "2020-04-26T04:11:11Z",
// 		"containerStatuses": [{
// 			"name": "hello-python",
// 			"state": {
// 				"running": {
// 					"startedAt": "2020-06-06T11:31:43Z"
// 				}
// 			},
// 			"lastState": {
// 				"terminated": {
// 					"exitCode": 255,
// 					"reason": "Error",
// 					"startedAt": "2020-05-25T11:52:18Z",
// 					"finishedAt": "2020-06-06T11:31:22Z",
// 					"containerID": "docker://fbc48ad25a0741671490144c67f5a039f553130ff397717fadd8b741a5d07536"
// 				}
// 			},
// 			"ready": true,
// 			"restartCount": 4,
// 			"image": "hello-python:latest",
// 			"imageID": "docker://sha256:8ff6faf796007cef4e907ade0573b803bdafa944a856707da49c72c3b453812d",
// 			"containerID": "docker://e4895845f908cd0926dce26ddd82422d281742757a3b777ce5c9347f0531727d",
// 			"started": true
// 		}],
// 		"qosClass": "BestEffort"
// 	}
// }