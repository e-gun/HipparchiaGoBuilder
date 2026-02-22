//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func ReconcileWkMaps() {
	// the information stored in TheIdtWkMap and TheCanonWkMap is typically out of sync
	// the IDT map has better data in some cases; the Canon map has better data in other cases
	// after loading new author data tables, consolidate this information, then update the author tables on the DB
	// specifically the names are better in the IDX and the genres, dates, and locations are in the Canon

	// a build might just do a few authors, so the IDX set is the smaller/correct one

	keys := global.TheIdtWkMap.Keys()

	//fmt.Printf("TheIdtWkMap Keys: %d\n", len(keys))
	//fmt.Printf("TheCanonWkMap Keys: %d\n", len(global.TheCanonWkMap.Keys()))

	for _, key := range keys {
		widx, e := global.TheIdtWkMap.Get(key)
		if !e {
			// this should be impossible
			fmt.Println("ReconcileWkMaps() TheIdtWkMap is missing", key)
		}
		wcan, e := global.TheCanonWkMap.Get(key)
		if !e {
			// spammy: many are in fact missing
			// fmt.Println("ReconcileWkMaps() TheCanonWkMap is missing", key)
			global.TheMasterWkMap.Set(key, widx)
			continue
		}
		mg := workrmerger(widx, wcan)
		global.TheMasterWkMap.Set(key, mg)
	}
}

func workrmerger(widx structs.DbWork, wcan structs.DbWork) structs.DbWork {
	// widx has the right LL0-5 values
	// widx.PrintOut()
	widx.Language = wcan.Language
	widx.Pub = wcan.Pub
	widx.Genre = wcan.Genre
	widx.Xmit = wcan.Xmit
	return widx
}
