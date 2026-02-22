//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines3cleanup

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

// expected to reach PrepareForDB() before you are ready to do this

func FixWorklineMetadata(prefix string, lines []structs.DbWorkline) []structs.DbWorkline {
	for i, line := range lines {
		line.WkUID = prefix + line.WkUID
		lines[i] = line
	}
	return lines
}
