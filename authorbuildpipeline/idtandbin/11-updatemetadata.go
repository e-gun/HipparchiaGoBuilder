//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

func StoreUpdatedMetadata() {
	// only meaningful if you are building the TLG corpus
	// next 2 slated for deletion: this happens earlier/elsewhere
	//auu, wkk := LoadGreekCanon(global.Config.GreekDir)
	//StoreCanonInSharedMaps(auu, wkk)
	ReconcileAuMaps()
	ReconcileWkMaps()
	StoreAUChanges()
	BulkStoreWkChanges()
}
