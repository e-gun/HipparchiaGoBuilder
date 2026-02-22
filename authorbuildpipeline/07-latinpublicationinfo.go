//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/idtandbin"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

func LoadLatinCanonPubInfoIntoWorkMap(datadir string) {
	workpubinfo := idtandbin.LoadLatinCanon(datadir)
	// wcan.Pub is authoritative in workrmerger() [called by ReconcileWkMaps()]
	// but we will just set the idt version...
	for k, v := range workpubinfo {
		widx, ok := global.TheIdtWkMap.Get(k)
		if !ok {
			fmt.Println("LoadLatinCanonPubInfoIntoWorkMap() cound not find key:", k)
		} else {
			widx.Pub = v
			global.TheIdtWkMap.Set(k, widx)
		}
	}
}
