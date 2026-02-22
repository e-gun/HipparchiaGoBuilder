//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package latetidyups

import "strings"

func PHIAuthorExclamations(ttc string) string {
	// a LAT author tries to retain "!"
	// it is supposed to be "∙" in a GRK author (i.e, "missing/unknown character")
	// but this is going to be a real problem in a CHR text where there is no "excitement", only loss
	return strings.Replace(ttc, "﹗", "·", -1)
}
