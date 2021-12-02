package scaling

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// NewFunctionScaler create a new scaler with the specified
// ScalingConfig
func NewFunctionScaler(config ScalingConfig) FunctionScaler {
	cache := FunctionCache {
		Cache:  make(map[string]*FunctionMeta),
		Expiry: config.CacheExpiry,
	}
	cacheUpdateFlag := FunctionCacheUpdateFlag {
		CacheUpdate: make(map[string]bool),
		Sync:        sync.RWMutex{},
	}

	return FunctionScaler {
		Cache:  &cache,
		Config: config,
		CacheUpdateFlag: &cacheUpdateFlag,
	}
}


// FunctionScaler scales from zero
type FunctionScaler struct {
	CacheUpdateFlag *FunctionCacheUpdateFlag
	Cache  *FunctionCache
	Config ScalingConfig
}

// FunctionScaleResult holds the result of scaling from zero
type FunctionScaleResult struct {
	Available bool
	Error     error
	Found     bool
	Duration  time.Duration
}

// Scale scales a function from zero replicas to 1 or the value set in
// the minimum replicas metadata
func (f *FunctionScaler) Scale(functionName, namespace string) {
	funcNameKey := functionName+"."+namespace
	if cachedResponse, hit := f.Cache.Get(funcNameKey); hit && cachedResponse.AvailableReplicas > 0 {
		return
	}
	if val, exists := f.CacheUpdateFlag.GetFlag(funcNameKey); exists {
		if val == false { // nobody got in
			f.CacheUpdateFlag.SetFlag(funcNameKey,true) //  then get in and lock the door (flag)
			queryResponse, err := f.Config.ServiceQuery.GetReplicas(functionName, namespace)  // external.go implement
			if err != nil {
				f.CacheUpdateFlag.SetFlag(funcNameKey,false)
				log.Printf("getwayScaler: firstly get replicas error %s\n",err.Error())
				return
			}
			f.Cache.Set(funcNameKey, queryResponse)
			if queryResponse.AvailableReplicas == 0 { //check to scale out from zero
				if queryResponse.Replicas == 0 { // expectedReplicas is not zero
					minReplicas := uint64(1)
					scaleResultErr := backoff(func(attempt int) error {
						//log.Printf("[Scale %d] function=%s 0 => %d requested", attempt, functionName, minReplicas)
						setScaleErr := f.Config.ServiceQuery.SetReplicas(functionName, namespace, minReplicas)
						if setScaleErr != nil {
							return fmt.Errorf("getwayScaler: unable to scale function [%s], err: %s", functionName, setScaleErr)
						}
						return nil
					}, int(f.Config.SetScaleRetries), f.Config.FunctionPollInterval)
					if scaleResultErr != nil {
						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
						log.Printf("getwayScaler: set replicas error %s\n",scaleResultErr.Error())
						return
					}
				}
				// if replicas > 0 or at least one setReplicas in backoff is not error
				for i := 0; i < int(f.Config.MaxPollCount); i++ {
					time.Sleep(f.Config.FunctionPollInterval)
					queryResponse, err = f.Config.ServiceQuery.GetReplicas(functionName, namespace)
					if err == nil {
						f.Cache.Set(funcNameKey, queryResponse)
					} else {
						log.Printf("getwayScaler: get replicas error, try again ...... %s\n",err.Error())
						continue
					}
					if queryResponse.AvailableReplicas > 0 {
						//log.Printf("[Scale] function=%s 0 => %d successful - %fs", functionName, queryResponse.AvailableReplicas, totalTime.Seconds())
						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
						return
					}
				}
			}
			// if (availReplicas==0 && maxMollCount trying number out)
			f.CacheUpdateFlag.SetFlag(funcNameKey, false) // get out and unlock the door
		} else {
			count := uint(0)
			for {
				time.Sleep(f.Config.FunctionPollInterval)
				if count > f.Config.MaxPollCount {
					break
				}
				count++
				if doorIsClosed, _ := f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {// wait util someone unlock the door
					break
				}
			}
		}
	} else {
		f.CacheUpdateFlag.SetFlag(funcNameKey,false)  //create a door firstly and open it to come in
		if doorIsClosed, _ := f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {
			f.CacheUpdateFlag.SetFlag(funcNameKey,true) //  then get in and lock the door (flag)
			queryResponse, err := f.Config.ServiceQuery.GetReplicas(functionName, namespace)  // external.go implement
			if err != nil {
				f.CacheUpdateFlag.SetFlag(funcNameKey,false)
				log.Printf("getwayScaler: firstly get replicas error %s\n",err.Error())
				return
			}
			f.Cache.Set(funcNameKey, queryResponse)
			if queryResponse.AvailableReplicas == 0 { //check to scale out from zero
				if queryResponse.Replicas == 0 { // expectedReplicas is not zero
					minReplicas := uint64(1)
					scaleResultErr := backoff(func(attempt int) error {
						//log.Printf("[Scale %d] function=%s 0 => %d requested", attempt, functionName, minReplicas)
						setScaleErr := f.Config.ServiceQuery.SetReplicas(functionName, namespace, minReplicas)
						if setScaleErr != nil {
							return fmt.Errorf("getwayScaler: unable to scale function [%s], err: %s", functionName, setScaleErr)
						}
						return nil
					}, int(f.Config.SetScaleRetries), f.Config.FunctionPollInterval)
					if scaleResultErr != nil {
						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
						log.Printf("getwayScaler: set replicas error %s\n",scaleResultErr.Error())
						return
					}
				}
				// if replicas > 0 or at least one setReplicas in backoff is not error
				for i := 0; i < int(f.Config.MaxPollCount); i++ {
					time.Sleep(f.Config.FunctionPollInterval)
					queryResponse, err = f.Config.ServiceQuery.GetReplicas(functionName, namespace)
					if err == nil {
						f.Cache.Set(funcNameKey, queryResponse)
					} else {
						log.Printf("getwayScaler: get replicas error, try again ...... %s\n",err.Error())
						continue
					}
					if queryResponse.AvailableReplicas > 0 {
						//log.Printf("[Scale] function=%s 0 => %d successful - %fs", functionName, queryResponse.AvailableReplicas, totalTime.Seconds())
						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
						return
					}
				}
			}
			// if (availReplicas==0 && maxMollCount trying number out)
			f.CacheUpdateFlag.SetFlag(funcNameKey, false) // get out and unlock the door
		} else {
			count := uint(0)
			for {
				time.Sleep(f.Config.FunctionPollInterval)
				if count > f.Config.MaxPollCount {
					break
				}
				count++
				if doorIsClosed, _ = f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {// wait util someone unlock the door
					break
				}
			}
		}
	}
	return
}

type routine func(attempt int) error

func backoff(r routine, attempts int, interval time.Duration) error {
	var err error

	for i := 0; i < attempts; i++ {
		res := r(i)
		if res != nil {
			err = res

			log.Printf("Attempt: %d, had error: %s\n", i, res)
		} else {
			err = nil
			break
		}
		time.Sleep(interval)
	}
	return err
}

//
//func (f *FunctionScaler) Scale(functionName, namespace string) FunctionScaleResult {
//	start := time.Now()
//	funcNameKey := functionName+"."+namespace
//	if cachedResponse, hit := f.Cache.Get(funcNameKey); hit &&
//		cachedResponse.AvailableReplicas > 0 {
//		return FunctionScaleResult {
//			Error:     nil,
//			Available: true,
//			Found:     true,
//			//Duration:  time.Since(start),
//		}
//	}
//	if val, exists := f.CacheUpdateFlag.GetFlag(funcNameKey); exists {
//		if val == false { // nobody got in
//			f.CacheUpdateFlag.SetFlag(funcNameKey,true) //  then get in and lock the door (flag)
//			queryResponse, err := f.Config.ServiceQuery.GetReplicas(functionName, namespace)  // external.go implement
//			if err != nil {
//				f.CacheUpdateFlag.SetFlag(funcNameKey,false)
//				return FunctionScaleResult {
//					Error:     err,
//					Available: false,
//					Found:     false,
//					//Duration:  time.Since(start),
//				}
//			}
//
//			f.Cache.Set(funcNameKey, queryResponse)
//
//			if queryResponse.AvailableReplicas == 0 { //scale out from zero
//				minReplicas := uint64(1)
//				if queryResponse.MinReplicas > 0 {
//					minReplicas = queryResponse.MinReplicas
//				}
//
//				scaleResultErr := backoff(func(attempt int) error {
//					if queryResponse.Replicas > 0 { // expectedReplicas is not zero
//						return nil
//					}
//					//log.Printf("[Scale %d] function=%s 0 => %d requested", attempt, functionName, minReplicas)
//					setScaleErr := f.Config.ServiceQuery.SetReplicas(functionName, namespace, minReplicas)
//					if setScaleErr != nil {
//						return fmt.Errorf("unable to scale function [%s], err: %s", functionName, setScaleErr)
//					}
//					return nil
//				}, int(f.Config.SetScaleRetries), f.Config.FunctionPollInterval)
//
//				if scaleResultErr != nil {
//					f.CacheUpdateFlag.SetFlag(funcNameKey, false)
//					return FunctionScaleResult {
//						Error:     scaleResultErr,
//						Available: false,
//						Found:     true,
//						Duration:  time.Since(start),
//					}
//				}
//
//				for i := 0; i < int(f.Config.MaxPollCount); i++ {
//					time.Sleep(f.Config.FunctionPollInterval)
//					queryResponse, err = f.Config.ServiceQuery.GetReplicas(functionName, namespace)
//					if err == nil {
//						f.Cache.Set(funcNameKey, queryResponse)
//					}
//					totalTime := time.Since(start)
//
//					if err != nil {
//						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
//						return FunctionScaleResult{
//							Error:     err,
//							Available: false,
//							Found:     true,
//							Duration:  totalTime,
//						}
//					}
//
//					if queryResponse.AvailableReplicas > 0 {
//						//log.Printf("[Scale] function=%s 0 => %d successful - %fs", functionName, queryResponse.AvailableReplicas, totalTime.Seconds())
//						f.CacheUpdateFlag.SetFlag(funcNameKey, false)
//						return FunctionScaleResult {
//							Error:     nil,
//							Available: true,
//							Found:     true,
//							Duration:  totalTime,
//						}
//					}
//				}
//			}
//			f.CacheUpdateFlag.SetFlag(funcNameKey, false) // get out and unlock the door
//		} else {
//			count := uint(0)
//			for {
//				time.Sleep(f.Config.FunctionPollInterval)
//				if count > f.Config.MaxPollCount {
//					break
//				}
//				count++
//				if doorIsClosed, _ := f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {// wait util someone unlock the door
//					break
//				}
//			}
//		}
//	} else {
//		f.CacheUpdateFlag.SetFlag(funcNameKey,false)  //create a door firstly and open it to come in
//		if doorIsClosed, _ := f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {
//			f.CacheUpdateFlag.SetFlag(funcNameKey,true)  // then get in and lock the door (flag)
//			queryResponse, err := f.Config.ServiceQuery.GetReplicas(functionName, namespace)  // external.go implement
//			if err != nil {
//				f.CacheUpdateFlag.SetFlag(funcNameKey,false)
//				return FunctionScaleResult{
//					Error:     err,
//					Available: false,
//					Found:     false,
//					//Duration:  time.Since(start),
//				}
//			}
//
//			f.Cache.Set(funcNameKey, queryResponse)
//
//			if queryResponse.AvailableReplicas == 0 {
//				minReplicas := uint64(1)
//				if queryResponse.MinReplicas > 0 {
//					minReplicas = queryResponse.MinReplicas
//				}
//
//				scaleResult := backoff(func(attempt int) error {
//					if queryResponse.Replicas > 0 { // expected is not 0
//						return nil
//					}
//
//					//log.Printf("[Scale %d] function=%s 0 => %d requested", attempt, functionName, minReplicas)
//					setScaleErr := f.Config.ServiceQuery.SetReplicas(functionName, namespace, minReplicas)
//					if setScaleErr != nil {
//						return fmt.Errorf("unable to scale function [%s], err: %s", functionName, setScaleErr)
//					}
//					return nil
//
//				}, int(f.Config.SetScaleRetries), f.Config.FunctionPollInterval)
//
//				if scaleResult != nil {
//					f.CacheUpdateFlag.SetFlag(funcNameKey,false)
//					return FunctionScaleResult{
//						Error:     scaleResult,
//						Available: false,
//						Found:     true,
//						//Duration:  time.Since(start),
//					}
//				}
//
//				for i := 0; i < int(f.Config.MaxPollCount); i++ {
//					time.Sleep(f.Config.FunctionPollInterval)
//					queryResponse, err = f.Config.ServiceQuery.GetReplicas(functionName, namespace)
//					if err == nil {
//						f.Cache.Set(funcNameKey, queryResponse)
//					}
//					//totalTime := time.Since(start)
//
//					if err != nil {
//						f.CacheUpdateFlag.SetFlag(funcNameKey,false)
//						return FunctionScaleResult{
//							Error:     err,
//							Available: false,
//							Found:     true,
//							//Duration:  totalTime,
//						}
//					}
//
//					if queryResponse.AvailableReplicas > 0 {
//						//log.Printf("[Scale] function=%s 0 => %d successful - %fs", functionName, queryResponse.AvailableReplicas, totalTime.Seconds())
//						f.CacheUpdateFlag.SetFlag(funcNameKey,false)
//						return FunctionScaleResult{
//							Error:     nil,
//							Available: true,
//							Found:     true,
//							//Duration:  totalTime,
//						}
//					}
//				}
//			}
//			f.CacheUpdateFlag.SetFlag(funcNameKey,false) // get out and unlock the door
//		} else {
//			count := uint(0)
//			for{
//				time.Sleep(f.Config.FunctionPollInterval)
//				if count > f.Config.MaxPollCount {
//					break
//				}
//				count++
//				if doorIsClosed, _ = f.CacheUpdateFlag.GetFlag(funcNameKey); doorIsClosed == false {// wait util someone unlock the door
//					break
//				}
//			}
//		}
//	}
//
//	return FunctionScaleResult {
//		Error:     nil,
//		Available: true,
//		Found:     true,
//		//Duration:  time.Since(start),
//	}
//}
