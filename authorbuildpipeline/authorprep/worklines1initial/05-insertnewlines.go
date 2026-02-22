//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

import "regexp"

var (
	newlinefinder = regexp.MustCompile(`(<hb-set_l|<hb-metadata_new)`)
)

func InsertNewLinesAtNewLevels(ttc string) string {
	// 	break up the file into something you can walk through line-by-line
	return newlinefinder.ReplaceAllString(ttc, "\n$1")
}
