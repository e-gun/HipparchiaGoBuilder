//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import "regexp"

var (
	qsfinder = regexp.MustCompile(`\^(\d{1,2})+`)
)

func InsertBlankQuarterSpaces(ttc string) string {
	const (
		QSTAG = `<hb_blank_quarter_spaces quantity="$1" /> `
	)
	return qsfinder.ReplaceAllString(ttc, QSTAG)
}
