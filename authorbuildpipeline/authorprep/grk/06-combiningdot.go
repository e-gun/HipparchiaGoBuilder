//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

import "strings"

func GreekCombiningDot(ttc string) string {
	return strings.ReplaceAll(ttc, "?", "\u0323")
}
