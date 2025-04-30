//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
)

func BulkStoreWkChanges() {
	keys := global.TheMasterWkMap.Keys()
	var works []structs.DbWork
	for _, key := range keys {
		wk, ok := global.TheMasterWkMap.Get(key)
		if ok {
			works = append(works, wk)
		}
	}
	insert.BulkInsertWorks(works)
}
