// Copyright (c) OpenFaaS Author(s). All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package scaling

import (
	"sync"
	"time"
)

// FunctionMeta holds the last refresh and any other
// meta-data needed for caching.
type FunctionMeta struct {
	LastRefresh          time.Time
	ServiceQueryResponse ServiceQueryResponse
}

// Expired find out whether the cache item has expired with
// the given expiry duration from when it was stored.
func (fm *FunctionMeta) Expired(expiry time.Duration) bool {
	return time.Now().After(fm.LastRefresh.Add(expiry))
}

// FunctionCache provides a cache of Function replica counts
type FunctionCache struct {
	Cache  map[string]*FunctionMeta
	Expiry time.Duration
	Sync   sync.RWMutex
}

type FunctionCacheUpdateFlag struct {
	CacheUpdate  map[string]bool
	Sync   sync.RWMutex
}

func (fc *FunctionCache) Set(funcNameKey string, serviceQueryResponse ServiceQueryResponse) {
	fc.Sync.Lock()
	defer fc.Sync.Unlock()

	if _, exists := fc.Cache[funcNameKey]; !exists {
		fc.Cache[funcNameKey] = &FunctionMeta{}
	}

	fc.Cache[funcNameKey].LastRefresh = time.Now()
	fc.Cache[funcNameKey].ServiceQueryResponse = serviceQueryResponse
}


// Get replica count for functionName
func (fc *FunctionCache) Get(funcNameKey string) (ServiceQueryResponse, bool) {
	replicas := ServiceQueryResponse{
		AvailableReplicas: 0,
	}

	hit := false
	fc.Sync.RLock()
	defer fc.Sync.RUnlock()

	if val, exists := fc.Cache[funcNameKey]; exists {
		replicas = val.ServiceQueryResponse
		hit = !val.Expired(fc.Expiry)
	}

	return replicas, hit
}

// Set replica count for functionName
func (fcuf *FunctionCacheUpdateFlag) SetFlag(funcNameKey string, flag bool) {
	fcuf.Sync.Lock()
	defer fcuf.Sync.Unlock()
	fcuf.CacheUpdate[funcNameKey] = flag
}

// Get replica count for functionName
func (fcuf *FunctionCacheUpdateFlag) GetFlag(funcNameKey string) (bool, bool) {
	fcuf.Sync.RLock()
	defer fcuf.Sync.RUnlock()

	val, exists := fcuf.CacheUpdate[funcNameKey]

	return val, exists
}