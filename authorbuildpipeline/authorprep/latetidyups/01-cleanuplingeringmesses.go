//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package latetidyups

import (
	"strings"
)

func LingeringMesses(ttc string) string {
	// 	we've made it to the bitter end but there is something ugly in the results
	//	here we can clean things up that we are too lazy/stupid/afraid-of-worse to prevent from ending up
	// all the way here at the end of the line

	// closequote at opening:
	// μέλει τὸν ἀκροατήν, ”κἂν ἄμουϲοϲ ᾖ” παντάπαϲι, --> μέλει τὸν ἀκροατήν, “κἂν ἄμουϲοϲ ᾖ” παντάπαϲι,
	ttc = strings.ReplaceAll(ttc, ` ”`, ` “`)

	return ttc
}
