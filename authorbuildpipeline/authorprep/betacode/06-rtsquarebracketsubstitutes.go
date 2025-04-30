//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
	"strconv"
	"strings"
)

var (
	rsbrouter   = regexp.MustCompile(`]\d{1,2}`)
	rsbrsubsmap = map[int]string{
		// "\",
		1: "❩", // swapped for ")"
		// 2: "\u232a", // "〉"
		2:  "⟩",
		3:  "❵",
		4:  "⟧",
		5:  "⌋",
		6:  "⌉",
		7:  "⌉",
		8:  "⌋",
		9:  "\u2027",
		10: "<hb-sp-largerthannormal>]</hb-sp-largerthannormal>",
		11: "\u208e",
		12: "\u2190",
		13: "\u005d</hb-sp-italic>", // supposed to be italic as well
		14: "\u003a\u007c",          // ":|"
		17: "\u230b\u230b",          // "⌋⌋"
		18: "\u27eb",                // "⟫"
		20: "\u23ab",                // "⎫"
		21: "\u23aa",                // "⎪"
		22: "\u23ac",                // "⎬"
		23: "\u23ad",                // "⎭"
		30: "\u239e",                // "⎞" typo in betacode manual: should say 239e and not 329e (㊞)
		31: "\u239f",
		32: "\u23a0", // "⎠" typo in betacode manual: should say 23a0 and not 32a0 (㊠)
		33: "｠</hb-parenthesis_ancient_punctuation>",
		34: "⸩</hb-sp-parenthesis_deletion_marker>",
		35: "<hb-pap-rt_bracket_35 />",
		49: "<hb-pap-rt_bracket_49 />", // 49-35
		51: "</hb-sp-erasedepigraphicaltext>",
		52: "</hb-sp-text-before-correction>", // // non-TLG; Text Before Correction; not an erasure; epigraphical (but in INS0150)
		53: "❩",                               // non-TLG; Parenthesis modern Punctuation (but in INS0150)
	}
)

func ReplaceRightSquareBrackets(ttc string) string {
	// Purge ] markup
	ttc = rsbrouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return rlsbsubstitutes(match)
	})

	return ttc
}

func rlsbsubstitutes(match string) string {
	const (
		UHP = `<hb-build-error>%s⟧<hb-build-error>`
	)

	// match looks like "]14": drop that initial character
	match = strings.TrimPrefix(match, "]")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := rsbrsubsmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("lsbsubstitutes()\t%s", match))
	}
	return substitute
}
