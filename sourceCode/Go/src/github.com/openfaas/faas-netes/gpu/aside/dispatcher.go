//@file  : alert.go
//@author: Yanan Yang
//@date  : 2020/4/24
//@note: for scale in/out alert monitor, run as an aside thread
package aside

import (
	"github.com/openfaas/faas-netes/gpu/metrics"
	"github.com/openfaas/faas-netes/gpu/repository"
	gTypes "github.com/openfaas/faas-netes/gpu/types"
	"log"
	"net/http"
	"time"
)

const MonitorInterval = time.Second * 1
const PodLoadRate = 0.8
/**
 * update function pods's lotteries with workload changes and the cluster capacity
 */
func RpsDispatcherMonitor(funcNamespace string, loadGenHost string, loadGenPort int) {
	log.Printf("dispatcher: workload-aware request dispatcher starts successfully, loadGenHost=%s, loadGenPort=%d monitor interval=%ds",
		loadGenHost,
		loadGenPort,
		MonitorInterval/1000000000)
	ticker := time.NewTicker(MonitorInterval)
	quit := make(chan struct{})


	loadGenQuery := metrics.NewLoadGenQuery(loadGenHost, loadGenPort, &http.Client{})

	//expr := url.QueryEscape(`sum(rate(gateway_function_invocation_total{function_name=~".*", code=~".*", kubernetes_namespace="`+funcNamespace+`"}[10s])) by (function_name)`)
	//for i:=-50;i<150;i++ {
	//	fmt.Println(int32(math.Sin(float64(i)*math.Pi/100)*400+400))
	//}
	var funcDeployStatus *gTypes.FuncDeployStatus
	for {
		select {
		case <-ticker.C:
			results, fetchErr := loadGenQuery.Fetch()
			if fetchErr != nil {
				log.Printf("Error querying LoadGen API: %s \n", fetchErr.Error())
				time.Sleep(MonitorInterval)
				continue
			}
			//fmt.Println(results)
			for _ , item := range *results {
				funcDeployStatus = repository.GetFunc(item.LoaderName)
				if funcDeployStatus == nil {
					//log.Printf("dispatcher: function %s is not in repository, skip to allocate requests lottery \n", item.LoaderName)
					continue
				}
				if item.RealRps == 0 {
					if funcDeployStatus.FuncLastRealRps == 0 {
						repository.UpdateFuncRealRps(item.LoaderName,0)
						repository.UpdateFuncLastRealRps(item.LoaderName,0)
					} else {
						repository.UpdateFuncRealRps(item.LoaderName, funcDeployStatus.FuncLastRealRps)
						repository.UpdateFuncLastRealRps(item.LoaderName,0)
					}
				} else {
					repository.UpdateFuncLastRealRps(item.LoaderName, item.RealRps)
					repository.UpdateFuncRealRps(item.LoaderName, item.RealRps)
				}

				readLock := repository.GetFuncScalingLockState(item.LoaderName)
				readLock.Lock()
				updateFunctionsLottery(item.RealRps, funcDeployStatus)
				repository.UpdateFuncPodsTotalLottery(item.LoaderName)
				readLock.Unlock()

			}



			/*funcName := "resnet50"
			funcDeployStatus := repository.GetFunc(funcName)
			if funcDeployStatus == nil {
				log.Printf("dispatcher: function %s is not in repository, skip to allocate requests lottery \n", funcName)
				continue
			}
			curQps := math.Sin(float64(counter)*math.Pi/100)*400+400
			repository.UpdateFuncRealRps(funcName, int32(curQps))
			updateFunctionsLottery(int32(curQps), funcDeployStatus)
			counter+=1
			if counter == 150 {
				return
			}
			break*/
		case <-quit:
			return
		}
	}
	log.Printf("dispatcher: workload-aware request dispatcher exits -------------------------\n")
}
func updateFunctionsLottery(curQps int32, funcDeployStatus *gTypes.FuncDeployStatus) {
	if curQps > funcDeployStatus.FuncPodMaxCapacity {
		//log.Println("============over load=======================")
		for _, funcConfig := range funcDeployStatus.FuncPodConfigMap {
			if funcConfig.FuncPodType == "i" {
				// todo: deal with the ppod
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName, funcConfig.ReqPerSecondMax)
			} else if funcConfig.FuncPodType == "v" {
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName, curQps-funcDeployStatus.FuncPodMaxCapacity)
			}
		}
	} else if curQps >= (int32(float64(funcDeployStatus.FuncPodMaxCapacity-funcDeployStatus.FuncPodMinCapacity)*PodLoadRate)+funcDeployStatus.FuncPodMinCapacity) {
	//} else if curQps <= funcDeployStatus.FuncPodMaxCapacity && curQps >= funcDeployStatus.FuncPodMinCapacity {
		for _, funcConfig := range funcDeployStatus.FuncPodConfigMap {
			if funcConfig.FuncPodType == "i" { //int转换为下取整
				scalingFactor := float32(funcConfig.ReqPerSecondMax-funcConfig.ReqPerSecondMin) / float32(funcDeployStatus.FuncPodMaxCapacity-funcDeployStatus.FuncPodMinCapacity)
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName, funcConfig.ReqPerSecondMax-int32(scalingFactor*float32(funcDeployStatus.FuncPodMaxCapacity-curQps)))
			} else if funcConfig.FuncPodType == "v" {
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName, 0)
			}
		}

	} else { // less than the minCap
		if curQps == 0 {
			for _, funcConfig := range funcDeployStatus.FuncPodConfigMap {
				if funcConfig.FuncPodType == "i" { //int转换为下取整
					repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName, 0)
				}
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, "v", 0)
			}
		} else {
			changePodForLowWorkload(funcDeployStatus, curQps)
		}

	}
}


func changePodForLowWorkload(funcDeployStatus *gTypes.FuncDeployStatus, curReq int32){
	/**
	 * try to using last changedPodCombine Cache
	 */
	lastChangedPodCombine := funcDeployStatus.FuncLastChangedPodCombine
	if lastChangedPodCombine != nil {
		tempFuncCapMin := funcDeployStatus.FuncPodMinCapacity-lastChangedPodCombine.MinDeletedSumCap
		tempFuncCapMax := funcDeployStatus.FuncPodMaxCapacity-lastChangedPodCombine.MaxDeletedSumCap
		//if curReq <= tempFuncCapMax && curReq >= tempFuncCapMin &&
		//	curReq >= (int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate) >> 3) + tempFuncCapMin {
		if curReq >= tempFuncCapMin && curReq <= tempFuncCapMax {
			for _ , remItem := range lastChangedPodCombine.RemainPodNameList {
				scalingFactor := float32(funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMin)/float32(tempFuncCapMax-tempFuncCapMin)
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, remItem, funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-int32(scalingFactor*float32(tempFuncCapMax-curReq)))
			}
			for _ , delItem := range lastChangedPodCombine.DeletePodNameList {
				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, delItem,0)
			}
			repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName,"v", 0)
			//log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to using last cache successfully, combine= %+v\n",
			//	funcDeployStatus.FunctionName,
			//	curReq,
			//	lastChangedPodCombine)
			return
		} else {
			//log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to using last cache failed (req is not met), combine= %+v\n",
			//	funcDeployStatus.FunctionName,
			//	curReq,
			//	lastChangedPodCombine)
		}
	} else {
		//log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to using last cache failed (cache is nil), combine= %+v\n",
		//	funcDeployStatus.FunctionName,
		//	curReq,
		//	lastChangedPodCombine)
	}

	/**
	 * cache is missed
	 */
	var letter []*gTypes.FuncPodConfig
	for _, value := range funcDeployStatus.FuncPodConfigMap {
		if value.FuncPodType == "i" {
			letter = append(letter, value)  //convert into podConfigMap into list
		}
	}
	if len(letter) == 0 {
		return
	}
	//var changedPodCombineList []*gTypes.ChangedPodCombine

	var i,j uint
	//deletedPodCount := int32(0)
	minDeletedSumCap := int32(0)
	maxDeletedSumCap := int32(0)

	foundFlag := false
	tempFuncCapMin := int32(0)
	tempFuncCapMax := int32(0)
	var deletePodNameList []string
	var remainPodNameList []string

	n := uint(len(letter))
	var maxCount uint = (1 << n) - 1
	//log.Printf("dispatcher: function=%s podsLen=%d, maxCount=%d .........\n", funcDeployStatus.FunctionName, len(letter), maxCount)

	var firstCandidateChangedPodCombine *gTypes.ChangedPodCombine
	var secondCandidateChangedPodCombine *gTypes.ChangedPodCombine
	var thirdCandidateChangedPodCombine *gTypes.ChangedPodCombine
	var fourthCandidateChangedPodCombine *gTypes.ChangedPodCombine
	var changedPodCombine *gTypes.ChangedPodCombine

	for i = 1; i < maxCount; i++ { //if there has 5 pods, try to delete 1...4 pods every time
		if i > 1040000 { //2^15
			if secondCandidateChangedPodCombine != nil || thirdCandidateChangedPodCombine != nil || fourthCandidateChangedPodCombine != nil{
				//log.Println("break............",i)
				break
			}
		}
		//deletedPodCount = int32(0)
		minDeletedSumCap = 0
		maxDeletedSumCap = 0

		deletePodNameList = []string{}
		remainPodNameList = []string{}
		for j = 0; j < n; j++ {
			if (i & (1 << j)) != 0 { //在做位运算的时候需要注意数据类型为uint
				//deletedPodCount++
				minDeletedSumCap += letter[j].ReqPerSecondMin
				maxDeletedSumCap += letter[j].ReqPerSecondMax
				deletePodNameList = append(deletePodNameList, letter[j].FuncPodName)
				//fmt.Printf("%d",j)
			} else {
				remainPodNameList = append(remainPodNameList, letter[j].FuncPodName)
			}
		}
		//changedPodCombineList = append(changedPodCombineList, changedPodCombine)
		tempFuncCapMin = funcDeployStatus.FuncPodMinCapacity - minDeletedSumCap
		tempFuncCapMax = funcDeployStatus.FuncPodMaxCapacity - maxDeletedSumCap
		if curReq <= tempFuncCapMax && curReq >= tempFuncCapMin &&
			curReq >= (int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate) + tempFuncCapMin) {
			//if curReq >= tempFuncCapMin && curReq <= tempFuncCapMax {
			firstCandidateChangedPodCombine = &gTypes.ChangedPodCombine {
				//DeletedPodCount: deletedPodCount,
				MinDeletedSumCap:  minDeletedSumCap,
				MaxDeletedSumCap:  maxDeletedSumCap,
				DeletePodNameList: deletePodNameList,
				RemainPodNameList: remainPodNameList,
			}
			//log.Printf("first found, curReq=%d,tempFuncCapMax=%d, value=%d, tempFuncCapMin=%d, loadRate=%f,%+v\n",
			//	curReq,tempFuncCapMax, int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate)+tempFuncCapMin, tempFuncCapMin, float64(curReq-tempFuncCapMin)/float64(tempFuncCapMax-tempFuncCapMin),
			//	firstCandidateChangedPodCombine)
			break
		} else if secondCandidateChangedPodCombine == nil &&
			curReq <= tempFuncCapMax && curReq >= tempFuncCapMin &&
			curReq >= ((int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate) >> 1) + tempFuncCapMin) {
			secondCandidateChangedPodCombine = &gTypes.ChangedPodCombine {
				//DeletedPodCount: deletedPodCount,
				MinDeletedSumCap:  minDeletedSumCap,
				MaxDeletedSumCap:  maxDeletedSumCap,
				DeletePodNameList: deletePodNameList,
				RemainPodNameList: remainPodNameList,
			}
			//log.Printf("second found, curReq=%d,tempFuncCapMax=%d, value=%d, tempFuncCapMin=%d, loadRate=%f,%+v\n",
			//	curReq,tempFuncCapMax,int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate) >> 1+tempFuncCapMin, tempFuncCapMin, float64(curReq-tempFuncCapMin)/float64(tempFuncCapMax-tempFuncCapMin),
			//	secondCandidateChangedPodCombine)
		} else if secondCandidateChangedPodCombine == nil && thirdCandidateChangedPodCombine == nil &&
			curReq <= tempFuncCapMax && curReq >= tempFuncCapMin &&
			curReq >= ((int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate) >> 2) + tempFuncCapMin) {
			thirdCandidateChangedPodCombine = &gTypes.ChangedPodCombine {
				//DeletedPodCount: deletedPodCount,
				MinDeletedSumCap:  minDeletedSumCap,
				MaxDeletedSumCap:  maxDeletedSumCap,
				DeletePodNameList: deletePodNameList,
				RemainPodNameList: remainPodNameList,
			}
			//log.Printf("third found, curReq=%d,tempFuncCapMax=%d, value=%d, tempFuncCapMin=%d, loadRate=%f,%+v\n",
			//	curReq,tempFuncCapMax,int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate) >> 2+tempFuncCapMin, tempFuncCapMin, float64(curReq-tempFuncCapMin)/float64(tempFuncCapMax-tempFuncCapMin),
			//	thirdCandidateChangedPodCombine)
		} else if secondCandidateChangedPodCombine == nil && thirdCandidateChangedPodCombine == nil && fourthCandidateChangedPodCombine == nil &&
			curReq <= tempFuncCapMax && curReq >= tempFuncCapMin {
			fourthCandidateChangedPodCombine = &gTypes.ChangedPodCombine {
				//DeletedPodCount: deletedPodCount,
				MinDeletedSumCap:  minDeletedSumCap,
				MaxDeletedSumCap:  maxDeletedSumCap,
				DeletePodNameList: deletePodNameList,
				RemainPodNameList: remainPodNameList,
			}
			//log.Printf("fourth found, curReq=%d,tempFuncCapMax=%d, value=%d, tempFuncCapMin=%d, loadRate=%f,%+v\n",
			//	curReq,tempFuncCapMax,int32(float64(tempFuncCapMax-tempFuncCapMin)*PodLoadRate) >> 3+tempFuncCapMin, tempFuncCapMin, float64(curReq-tempFuncCapMin)/float64(tempFuncCapMax-tempFuncCapMin),
			//	fourthCandidateChangedPodCombine)
		}
	}

	if firstCandidateChangedPodCombine != nil {
		changedPodCombine = firstCandidateChangedPodCombine
		foundFlag = true
	} else if secondCandidateChangedPodCombine != nil {
		changedPodCombine = secondCandidateChangedPodCombine
		foundFlag = true
	} else if thirdCandidateChangedPodCombine != nil {
		changedPodCombine = thirdCandidateChangedPodCombine
		foundFlag = true
	} else if fourthCandidateChangedPodCombine != nil {
		changedPodCombine = fourthCandidateChangedPodCombine
		foundFlag = true
	} else {
		//log.Printf("dispatcher: no CandidateChangedPodCombine found\n")
	}

	//log.Printf("dispatcher: changedPodCombineList=%+v\n", changedPodCombineList)
	/**
	 * decide to scale down function pod according to their capacities and the current RPS
	 */
	if foundFlag == true {
		tempFuncCapMin = funcDeployStatus.FuncPodMinCapacity - changedPodCombine.MinDeletedSumCap
		tempFuncCapMax = funcDeployStatus.FuncPodMaxCapacity - changedPodCombine.MaxDeletedSumCap
		for _ , remItem := range changedPodCombine.RemainPodNameList {
			scalingFactor := float32(funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMin)/float32(tempFuncCapMax-tempFuncCapMin)
			repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, remItem, funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-int32(scalingFactor*float32(tempFuncCapMax-curReq)))
		}
		for _ , delItem := range changedPodCombine.DeletePodNameList {
			repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, delItem,0)
		}
		repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName,"v",0)
		repository.UpdateFuncLastChangedPodCombine(funcDeployStatus.FunctionName, changedPodCombine)
	} else { //we need to delete all existing function pods
		if curReq >= funcDeployStatus.FuncPodMinCapacity {
			for _, funcConfig := range funcDeployStatus.FuncPodConfigMap {
				if funcConfig.FuncPodType == "i" { //int转换为下取整
					scalingFactor := float32(funcConfig.ReqPerSecondMax-funcConfig.ReqPerSecondMin) / float32(funcDeployStatus.FuncPodMaxCapacity-funcDeployStatus.FuncPodMinCapacity)
					repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName, funcConfig.ReqPerSecondMax-int32(scalingFactor*float32(funcDeployStatus.FuncPodMaxCapacity-curReq)))

				} else if funcConfig.FuncPodType == "v" {
					repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName, 0)
				}
			}
			curReq = -1 //allocation finished
		} else {  //curReq < funcDeployStatus.FuncPodMinCapacity
			for _, item := range funcDeployStatus.FuncPodConfigMap {
				if item.FuncPodType == "i" {
					if item.ReqPerSecondMax >= curReq {
						repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, item.FuncPodName, curReq)
					} else {
						repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, item.FuncPodName, item.ReqPerSecondMax)
					}
					curReq = curReq - item.ReqPerSecondMax
					if curReq <= 0 {
						break
					}
				}
			}
		}

	}

	if curReq <= 0 {
		repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName,"v",0)
	} else {
		//log.Println("dispatcher: no luckcy dog found")
		repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName,"v", curReq)
	}
}









//changedPodCombine := &gTypes.ChangedPodCombine {
//	DeletedPodCount: deletedPodCount,
//	MinDeletedSumCap: minDeletedSumCap,
//	MaxDeletedSumCap: maxDeletedSumCap,
//	DeletePodNameList: deletePodNameList,
//	RemainPodNameList: nil,
//}
//repository.UpdateFuncLastChangedPodCombine(funcDeployStatus.FunctionName, changedPodCombine)
//fmt.Println()
//log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to 0\n", funcDeployStatus.FunctionName, curReq)


//fmt.Println()
//log.Printf("dispatcher: function=%s, curReq=%d and changed lottery is as follows:\n", funcDeployStatus.FunctionName, curReq)
//for _, v:= range funcDeployStatus.FuncPodConfigMap {
//	log.Printf("dispatcher: curRes=%d, funcPod=%s, Lottery=%d \n",curReq, v.FuncPodName, v.ReqPerSecondLottery)
//}


//func changeLotteryWithLowWorkload(funcDeployStatus *gTypes.FuncDeployStatus, curReq int32){
//	/**
//	 * try to using last changedPodCombine Cache
//	 */
//	lastChangedPodCombine := funcDeployStatus.FuncLastChangedPodCombine
//	if lastChangedPodCombine != nil {
//		tempFuncCapMin := funcDeployStatus.FuncPodMinCapacity-lastChangedPodCombine.MinDeletedSumCap
//		tempFuncCapMax := funcDeployStatus.FuncPodMaxCapacity-lastChangedPodCombine.MaxDeletedSumCap
//		if curReq >= tempFuncCapMin && curReq <= tempFuncCapMax {
//			for _ , delItem := range lastChangedPodCombine.DeletePodNameList {
//				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, delItem, 0)
//			}
//			for _ , remItem := range lastChangedPodCombine.RemainPodNameList {
//				scalingFactor := float32(funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMin)/float32(tempFuncCapMax-tempFuncCapMin)
//				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, remItem, funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-int32(scalingFactor*float32(tempFuncCapMax-curReq)))
//			}
//			log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to using last cache successfully, combine= %+v\n",
//				funcDeployStatus.FunctionName,
//				curReq,
//				lastChangedPodCombine)
//			return
//		} else {
//			log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to using last cache failed (req is not met), combine= %+v\n",
//				funcDeployStatus.FunctionName,
//				curReq,
//				lastChangedPodCombine)
//		}
//	} else {
//		log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to using last cache failed (cache is nil), combine= %+v\n",
//			funcDeployStatus.FunctionName,
//			curReq,
//			lastChangedPodCombine)
//	}
//
//	/**
//	 * cache is missed
//	 */
//	var letter []*gTypes.FuncPodConfig
//	for _, value := range funcDeployStatus.FuncPodConfigMap {
//		if value.FuncPodType == "i" {
//			letter = append(letter, value)  //convert into podConfigMap into list
//		}
//	}
//	if len(letter) == 0 {
//		return
//	}
//	var changedPodCombineList []*gTypes.ChangedPodCombine
//
//	var i,j uint
//	deletedPodCount := int32(0)
//	minDeletedSumCap := int32(0)
//	maxDeletedSumCap := int32(0)
//
//	n := uint(len(letter))
//	var maxCount uint = (1 << n) - 1
//	log.Printf("dispatcher: function=%s podsLen=%d, maxCount=%d\n", funcDeployStatus.FunctionName, len(letter), maxCount)
//	for i = 1; i < maxCount; i++ {  //if there has 5 pods, try to delete 1...4 pods every time
//		deletedPodCount = int32(0)
//		minDeletedSumCap = int32(0)
//		maxDeletedSumCap = int32(0)
//		var deletePodNameList []string
//		var remainPodNameList []string
//		for j = 0; j < n; j++ {
//			if (i & (1 << j)) != 0 { //在做位运算的时候需要注意数据类型为uint
//				deletedPodCount++
//				minDeletedSumCap += letter[j].ReqPerSecondMin
//				maxDeletedSumCap += letter[j].ReqPerSecondMax
//				deletePodNameList = append(deletePodNameList, letter[j].FuncPodName)
//				//fmt.Printf("%d",j)
//			} else {
//				remainPodNameList = append(remainPodNameList, letter[j].FuncPodName)
//			}
//		}
//		changedPodCombine := &gTypes.ChangedPodCombine {
//			DeletedPodCount: deletedPodCount,
//			MinDeletedSumCap: minDeletedSumCap,
//			MaxDeletedSumCap: maxDeletedSumCap,
//			DeletePodNameList: deletePodNameList,
//			RemainPodNameList: remainPodNameList,
//		}
//		changedPodCombineList = append(changedPodCombineList, changedPodCombine)
//	}
//	//log.Printf("dispatcher: changedPodCombineList=%+v\n", changedPodCombineList)
//	/**
//	 * decide to scale down function pod according to their capacities and the current RPS
//	 */
//	foundFlag := false
//	tempFuncCapMin := int32(0)
//	tempFuncCapMax := int32(0)
//	candidateNum := int32(n)-1
//
//	for c := candidateNum; c > 0 && foundFlag == false; c-- {  //Minimum pod-number left (MPL) principle
//		for _ , item := range changedPodCombineList {
//			if item.DeletedPodCount == c {
//				tempFuncCapMin = funcDeployStatus.FuncPodMinCapacity-item.MinDeletedSumCap
//				tempFuncCapMax = funcDeployStatus.FuncPodMaxCapacity-item.MaxDeletedSumCap
//				if curReq >= tempFuncCapMin && curReq <= tempFuncCapMax {
//					for _ , delItem := range item.DeletePodNameList {
//						repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, delItem, 0)
//						//funcDeployStatus.FuncPodConfigMap[changedPodCombineList[j].DeletePodNameList[k]].ReqPerSecondLottery = 0
//					}
//					for _ , remItem := range item.RemainPodNameList {
//						scalingFactor := float32(funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMin)/float32(tempFuncCapMax-tempFuncCapMin)
//						repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, remItem, funcDeployStatus.FuncPodConfigMap[remItem].ReqPerSecondMax-int32(scalingFactor*float32(tempFuncCapMax-curReq)))
//						//funcDeployStatus.FuncPodConfigMap[changedPodCombineList[j].RemainPodNameList[k]].ReqPerSecondLottery = funcDeployStatus.FuncPodConfigMap[changedPodCombineList[j].RemainPodNameList[k]].ReqPerSecondMax-int32(scalingFactor*float32(tempFuncCapMax-int32(curReq)))
//					}
//					foundFlag = true
//					repository.UpdateFuncLastChangedPodCombine(funcDeployStatus.FunctionName, item)
//					break
//				}
//			}
//		}
//	}
//	if foundFlag == false { //we need to delete all existing function pods
//		for _, v := range funcDeployStatus.FuncPodConfigMap {
//			if v.FuncPodType == "i" {
//				//deletePodNameList = append(deletePodNameList, v.FuncPodName)
//				repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, v.FuncPodName,0)
//			} else if v.FuncPodType == "v" {
//				if curReq == 0 {
//					repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, v.FuncPodName,999)
//				} else {
//					repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, v.FuncPodName, curReq)
//				}
//			}
//		}
//		//changedPodCombine := &gTypes.ChangedPodCombine {
//		//	DeletedPodCount: deletedPodCount,
//		//	MinDeletedSumCap: minDeletedSumCap,
//		//	MaxDeletedSumCap: maxDeletedSumCap,
//		//	DeletePodNameList: deletePodNameList,
//		//	RemainPodNameList: nil,
//		//}
//		//repository.UpdateFuncLastChangedPodCombine(funcDeployStatus.FunctionName, changedPodCombine)
//		//fmt.Println()
//		//log.Printf("dispatcher: function=%s, curReq= %d and update all func pods' lotteries to 0\n", funcDeployStatus.FunctionName, curReq)
//	} else {
//		funcConfig, exist := funcDeployStatus.FuncPodConfigMap["v"]
//		if exist {
//			repository.UpdateFuncPodLottery(funcDeployStatus.FunctionName, funcConfig.FuncPodName,0)
//		}
//	}
//
//	//fmt.Println()
//	//log.Printf("dispatcher: function=%s, curReq=%d and changed lottery is as follows:\n", funcDeployStatus.FunctionName, curReq)
//	//for _, v:= range funcDeployStatus.FuncPodConfigMap {
//	//	log.Printf("dispatcher: curRes=%d, funcPod=%s, Lottery=%d \n",curReq, v.FuncPodName, v.ReqPerSecondLottery)
//	//}
//}

