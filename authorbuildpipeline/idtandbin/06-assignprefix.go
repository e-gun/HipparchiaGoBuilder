//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	structs2 "github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

// AssignPrefixToAuthor - because the IDX and BIN files do not know about 'gr' or 'lt'
func AssignPrefixToAuthor(pfx string, author *structs2.DbAuthor) {
	author.UID = pfx + author.UID
}

// AssignPrefixToWork - because the IDX and BIN files do not know about 'gr' or 'lt'
func AssignPrefixToWork(pfx string, work *structs2.DbWork) {
	work.UID = pfx + work.UID
}
