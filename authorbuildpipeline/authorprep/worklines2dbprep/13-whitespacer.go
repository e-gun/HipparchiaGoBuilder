//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"regexp"
)

var (
	whitespacer = regexp.MustCompile(`(^\s|\s$)`)
)

func NoLeadingOrTrailingSpaces(lines []structs.DbWorkline) []structs.DbWorkline {
	// 	get rid of whitespace at ends of columns
	//	otherwise HipparchiaServer is constantly doing this

	for i := range lines {
		lines[i].MarkedUp = whitespacer.ReplaceAllString(lines[i].MarkedUp, "")
		lines[i].Accented = whitespacer.ReplaceAllString(lines[i].Accented, "")
		lines[i].Stripped = whitespacer.ReplaceAllString(lines[i].Stripped, "")
	}
	return lines
}
