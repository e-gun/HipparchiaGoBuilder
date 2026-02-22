//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

var (
	singletonouter = regexp.MustCompile(`[#%@{}]`)
	singletonswap  = map[string]string{
		"#": "\u0374", // Greek Numeral Sign: ʹ
		"%": "\u2020", // dagger: †
		"@": "<hb-tabbedtext />",
		"{": "<hb-speaker>",
		"}": "</hb-speaker>",
	}
)

// ReplaceSingletons - swap out single characters; "`" is the one that *has* to be done very late
func ReplaceSingletons(ttc string) string {
	ttc = singletonouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return singsubstitutes(match)
	})

	// The null character (`) is used to separate a Beta escape code when a numeral follows
	//immediately. So &4`1& is used to represent the Indic-Arabic numeral one in superscript Roman
	//font.

	// this means you cannot run this replacement before andsignfontshifts or dollarsignfontshifts
	ttc = strings.ReplaceAll(ttc, "`", "")

	return ttc
}

func singsubstitutes(match string) string {
	const (
		UHP = `<hgb-build-error>%s<hgb-build-error>`
	)

	substitute, exists := singletonswap[match]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("singsubstitutes()\t%s", match))
	}
	return substitute
}
