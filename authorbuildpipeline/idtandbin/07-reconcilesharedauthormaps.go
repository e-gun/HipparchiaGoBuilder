//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"fmt"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func ReconcileAuMaps() {
	// the information stored in TheIdtAuMap and TheCanonAuMap is typically out of sync
	// the IDT map has better data in some cases; the Canon map has better data in other cases
	// after loading new author data tables, consolidate this information, then update the author tables on the DB
	// specifically the names are better in the IDX and the genres, dates, and locations are in the Canon

	// a build might just do a few authors, so the IDX set is the smaller/correct one

	//fmt.Println(global.TheMasterAuMap.Keys())
	//fmt.Println(global.TheIdtAuMap.Keys())
	//fmt.Println(global.TheCanonAuMap.Keys())

	keys := global.TheIdtAuMap.Keys()
	for _, key := range keys {
		aidx, _ := global.TheIdtAuMap.Get(key)

		if !strings.HasPrefix(aidx.UID, global.TLGABBREV) {
			global.TheMasterAuMap.Set(key, aidx)
			continue
		}

		acan, e := global.TheCanonAuMap.Get(key)
		if !e {
			fmt.Println("ReconcileAuMaps() TheCanonAuMap is missing", key)
			CleanDBAuthorNames(&aidx)
			aidx.PrintOut()
			// plutarch (!?) is broken; should fix the idxname, etx before setting thing
			global.TheMasterAuMap.Set(key, aidx)
			continue
		}

		mg := authormerger(aidx, acan)
		global.TheMasterAuMap.Set(key, mg)
	}
}

func MapIdtAuMapOntoMasterAuMap() {
	keys := global.TheIdtAuMap.Keys()
	for _, key := range keys {
		aidx, _ := global.TheIdtAuMap.Get(key)
		global.TheMasterAuMap.Set(key, aidx)
	}
}

func authormerger(aidx structs.DbAuthor, acan structs.DbAuthor) structs.DbAuthor {
	return structs.DbAuthor{
		UID:       aidx.UID,
		Language:  acan.Language,
		IDXname:   aidx.IDXname,
		Name:      aidx.Name,
		Shortname: aidx.Shortname,
		Cleaname:  aidx.Cleaname,
		Genres:    acan.Genres,
		RecDate:   acan.RecDate,
		ConvDate:  acan.ConvDate, // set by a call to dating.ParseTLGDate() by idtandbin.LoadGreekCanon()
		Location:  acan.Location,
		WorkList:  nil,
	}
}
