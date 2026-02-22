//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package global

import (
	"sync"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	TheMasterAuMap     = NewSafeAuMap()
	TheIdtAuMap        = NewSafeAuMap()
	TheCanonAuMap      = NewSafeAuMap()
	TheDefunctIdtAuMap = NewSafeAuMap()
)

type SafeAuMap struct {
	mu   sync.Mutex
	data map[string]structs.DbAuthor
}

func NewSafeAuMap() *SafeAuMap {
	return &SafeAuMap{
		data: make(map[string]structs.DbAuthor),
	}
}

func (sm *SafeAuMap) Set(key string, value structs.DbAuthor) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeAuMap) Get(key string) (structs.DbAuthor, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	val, ok := sm.data[key]
	return val, ok
}

func (sm *SafeAuMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

func (sm *SafeAuMap) Keys() []string {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	keys := make([]string, 0, len(sm.data))
	for k := range sm.data {
		keys = append(keys, k)
	}
	return keys
}

func (sm *SafeAuMap) Purge() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data = make(map[string]structs.DbAuthor)
}
