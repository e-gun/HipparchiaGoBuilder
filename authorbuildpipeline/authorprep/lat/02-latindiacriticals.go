//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lat

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
)

var (
	findlatindiacrit = regexp.MustCompile(`[aeiouyAEIOUV][+\\=/]`)
	latinswapmap     = map[string]string{
		"a/":  "\u00e1",
		"e/":  "\u00e9",
		"i/":  "\u00ed",
		"o/":  "\u00f3",
		"u/":  "\u00fa",
		"y/":  "\u00fd",
		"A/":  "\u00c1",
		"E/":  "\u00c9",
		"I/":  "\u00cd",
		"O/":  "\u00d3",
		"U/":  "\u00da",
		"V/":  "\u00da",
		"a+":  "ä",
		"A+":  "Ä",
		"e+":  "ë",
		"E+":  "Ë",
		"i+":  "ï",
		"I+":  "Ï",
		"o+":  "ö",
		"O+":  "Ö",
		"u+":  "ü",
		"U+":  "Ü",
		"a=":  "â",
		"A=":  "Â",
		"e=":  "ê",
		"E=":  "Ê",
		"i=":  "î",
		"I=":  "Î",
		"o=":  "ô",
		"O=":  "Ô",
		"u=":  "û",
		"U=":  "Û",
		"V=":  "Û",
		"a\\": "à",
		"A\\": "À",
		"e\\": "è",
		"E\\": "È",
		"i\\": "ì",
		"I\\": "Ì",
		"o\\": "ò",
		"O\\": "Ò",
		"u\\": "ù",
		"U\\": "Ù",
		"V\\": "Ù",
	}
)

func ConvertLatinDiacriticals(betacode string) string {
	// any greek betacode that remains in the text will turn into a big problem: BÉRWNI for BE/RWNI for βέρωνι
	return findlatindiacrit.ReplaceAllStringFunc(betacode, latinswap)
}

func latinswap(match string) string {
	m, ok := latinswapmap[match]
	if !ok {
		global.MSG(fmt.Sprintf("latinswap() failed for '%s'", match))
	}
	return m
}
