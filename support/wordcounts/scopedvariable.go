//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

// HeadwordLookupMap could produce a race condition; but it is read-only after write-once

var (
	HeadwordLookupMap = make(map[string][]string)
)
