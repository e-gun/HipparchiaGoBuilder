//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package latetidyups

func TidyUp(ttc string) string {
	ttc = LingeringMesses(ttc)
	ttc = PurgeHybrid(ttc)
	return ttc
}
