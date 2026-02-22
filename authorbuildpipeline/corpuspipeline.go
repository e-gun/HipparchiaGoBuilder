//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/idtandbin"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
)

func RunCorpusPipeline(datadir string, prefix string) {
	// the next matters if you try to build greek and latin in the same run
	global.TheIdtWkMap.Purge()
	global.TheCanonAuMap.Purge()
	global.TheMasterWkMap.Purge()
	global.TheIdtAuMap.Purge()
	global.TheCanonAuMap.Purge()
	global.TheMasterAuMap.Purge()

	if prefix == "TLG" {
		// this will load TheCanonAuMap, etc
		dba, dbw := idtandbin.LoadGreekCanon(global.Config.GreekDir)
		idtandbin.StoreCanonInSharedMaps(dba, dbw)
	}

	DeleteOldDataFromAuthorsTable(prefix)
	DeleteOldDataFromWorksTable(prefix)

	auu := BuildAuthorsSlice(datadir, prefix)
	FanoutBuilder(auu, datadir)

	idtandbin.MapIdtAuMapOntoMasterAuMap()

	if prefix == "TLG" {
		tlgmetadata := idtandbin.GatherTLGAuthorMetadata()
		AssignTLGMetadataToAuthors(tlgmetadata)
	}

	if prefix == "LAT" {
		LoadLatinGenresIntoWorkMap()
		LoadLatinCanonPubInfoIntoWorkMap(datadir)
	}

	idtandbin.ReconcileAuMaps()
	idtandbin.ReconcileWkMaps()
	idtandbin.StoreAUChanges()
	idtandbin.BulkStoreWkChanges()

	// free of real cost to do this as a goroutine?
	go insert.BuildTrigramIndices()

	notes := ""
	if prefix != "TLG" && prefix != "LAT" {
		notes = fmt.Sprintf("Reproducible: %t", global.Config.Reproducible)
	}
	insert.InsertBuildMetadata(prefix, notes)
}
