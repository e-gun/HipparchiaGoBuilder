//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
)

func StoreAUChanges() {
	keys := global.TheMasterAuMap.Keys()
	auu := make([]structs.DbAuthor, len(keys))
	for i, key := range keys {
		au, e := global.TheMasterAuMap.Get(key)
		if !e {
			fmt.Println("StoreAUChanges() TheMasterAuMap is missing", key)
		}
		// insert.InsertOneAuthorIntoAuthorTable(&au)
		auu[i] = au
	}
	insert.BulkInsertIntoAuthorsTable(auu)

	fmt.Printf("StoreAUChanges() modified %d authors\n", len(keys))
}
