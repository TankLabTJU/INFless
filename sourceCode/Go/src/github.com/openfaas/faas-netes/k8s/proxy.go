// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package k8s

import (
	"fmt"
	corelister "k8s.io/client-go/listers/core/v1"
	"net/url"
	"strings"
	"sync"
)

// watchdogPort for the OpenFaaS function watchdog
const watchdogPort = 8080
const tensorflowServingPort = 8501
func NewFunctionLookup(ns string, lister corelister.EndpointsLister) *FunctionLookup {
	return &FunctionLookup {
		DefaultNamespace: ns,
		EndpointLister:   lister,
		Listers:          map[string]corelister.EndpointsNamespaceLister{},
		lock:             sync.RWMutex{},
	}
}

type FunctionLookup struct {
	DefaultNamespace string
	EndpointLister   corelister.EndpointsLister
	Listers          map[string]corelister.EndpointsNamespaceLister

	lock sync.RWMutex
}

func (f *FunctionLookup) GetLister(ns string) corelister.EndpointsNamespaceLister {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.Listers[ns]
}

func (f *FunctionLookup) SetLister(ns string, lister corelister.EndpointsNamespaceLister) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.Listers[ns] = lister
}

func getNamespace(name, defaultNamespace string) string {
	namespace := defaultNamespace
	if strings.Contains(name, ".") {
		namespace = name[strings.LastIndexAny(name, ".")+1:]
	}
	return namespace
}

func (l *FunctionLookup) Resolve(name string) (url.URL, error) {
	/*
		namespace := getNamespace(name, l.DefaultNamespace)
		if err := l.verifyNamespace(namespace); err != nil {
			return url.URL{}, err
		}

		nsEndpointLister := l.GetLister(namespace)

		if nsEndpointLister == nil {
			l.SetLister(namespace, l.EndpointLister.Endpoints(namespace))

			nsEndpointLister = l.GetLister(namespace)
		}

		svc, err := nsEndpointLister.Get(name)
		if err != nil {
			return url.URL{}, fmt.Errorf("error listing %s.%s %s", name, namespace, err.Error())
		}

		if len(svc.Subsets) == 0 {
			return url.URL{}, fmt.Errorf("no subsets available for %s.%s", name, namespace)
		}

		all := len(svc.Subsets[0].Addresses)
		if len(svc.Subsets[0].Addresses) == 0 {
			return url.URL{}, fmt.Errorf("no addresses in subset for %s.%s", name, namespace)
		}

		target := rand.Intn(all)

		serviceIP := svc.Subsets[0].Addresses[target].IP


		urlStr := fmt.Sprintf("http://%s:%d", serviceIP, watchdogPort)


		urlRes, err := url.Parse(urlStr)
		if err != nil {
			return url.URL{}, err
		}*/

	urlStr := fmt.Sprintf("http://%s:%d", reqDispatcher(name), tensorflowServingPort)
	//log.Printf("proxy: dispatche request to %s", urlStr)
	urlRes, err := url.Parse(urlStr)
	if err != nil {
		return url.URL{}, err
	}
	urlRes.Path = fmt.Sprintf("/v1/models/%s:predict", name)
	return *urlRes, nil
}

func (l *FunctionLookup) verifyNamespace(name string) error {
	if name != "kube-system" {
		return nil
	}
	// ToDo use global namepace parse and validation
	return fmt.Errorf("namespace not allowed")
}

/**
 * ipod and vpod could be dispatched requests
 * ppod does not provider service unless it is converted into ipod for cold start
 */

func reqDispatcher(funcName string) string {
	//counter := int32(0)
	ip := "127.3.1.1"
	//funcStatus := repository.GetFunc(funcName)
	//if funcStatus != nil {
	//	if funcStatus.FuncPodTotalLottery == 0 {
	//		//log.Printf("+++++++++++++++++++++++++++++++ FuncPodIp 127.1.1.1\n", )
	//		ip = "127.1.1.1"
	//	} else {
	//		winner := rand.Intn(int(funcStatus.FuncPodTotalLottery))
	//
	//
	//		//log.Printf("+++++++++++++++++++++++++++++++ winner=%d\n",winner)
	//		for _, v := range funcStatus.FuncSortedPodNameList {
	//			//log.Printf("+++++++++++++++++++++++++++++++ key=%s, value=%+v\n",k, v)
	//			//log.Printf("+++++++++++++++++++++++++++++++ counter=%d+%d\n", counter, v.ReqPerSecondLottery)
	//			counter = counter + funcStatus.FuncPodConfigMap[v].ReqPerSecondLottery
	//			//log.Printf("+++++++++++++++++++++++++++++++ counter > winnner %d %d\n", counter, winner)
	//			if counter > int32(winner) {
	//				//log.Printf("+++++++++++++++++++++++++++++++ FuncPodIp %s\n", funcStatus.FuncPodConfigMap[v].FuncPodIp)
	//				return funcStatus.FuncPodConfigMap[v].FuncPodIp
	//			}
	//		}
	//	}
	//}
	return ip
}
