//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package global

import (
	"sync"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	TheMasterWkMap = NewSafeWkMap()
	TheIdtWkMap    = NewSafeWkMap()
	TheCanonWkMap  = NewSafeWkMap()
)

type SafeWkMap struct {
	mu   sync.Mutex
	data map[string]structs.DbWork
}

func NewSafeWkMap() *SafeWkMap {
	return &SafeWkMap{
		data: make(map[string]structs.DbWork),
	}
}

func (sm *SafeWkMap) Set(key string, value structs.DbWork) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeWkMap) Get(key string) (structs.DbWork, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	val, ok := sm.data[key]
	return val, ok
}

func (sm *SafeWkMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

func (sm *SafeWkMap) Keys() []string {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	keys := make([]string, 0, len(sm.data))
	for k := range sm.data {
		keys = append(keys, k)
	}
	return keys
}

func (sm *SafeWkMap) Purge() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data = make(map[string]structs.DbWork)
}
