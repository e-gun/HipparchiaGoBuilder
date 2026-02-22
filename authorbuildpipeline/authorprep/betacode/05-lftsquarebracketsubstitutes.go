//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

var (
	lsbrouter   = regexp.MustCompile(`\[\d{1,2}`)
	lsbrsubsmap = map[int]string{
		// "\",
		1: "❨", // supposed to be parenthesis "("; but can interfere with betacode parsing; either swap here or change order of execution
		// 2: "\u2329", // "〈"
		2:  "⟨",
		3:  "❴",
		4:  "⟦",
		5:  "⌊",
		6:  "⌈",
		7:  "⌈",
		8:  "⌊",
		9:  "\u2027",
		10: "<hb-sp-largerthannormal>[</hb-sp-largerthannormal>",
		11: "\u208d",
		12: "\u2192",
		13: "<hb-sp-italic>\u005b", // supposed to be italic as well
		14: "\u007c\u003a",
		17: "\u230a\u230a",
		18: "27ea",
		20: "\u23a7",
		21: "\u23aa",
		22: "\u23a8",
		23: "\u23a9",
		30: "\u239b",
		31: "\u239c",
		32: "\u239d",
		// this one is odd: you will "open" it 17x in dp0002 and "close" it 1x; what is really going on?
		33: "<hb-parenthesis_ancient_punctuation>｟",
		34: "<hb-sp-parenthesis_deletion_marker> ⸨",
		35: "<hb-pap-lt_bracket_35 />",
		49: "<hb-pap-lt_bracket_49 />", // 49-35
		51: "<hb-sp-erasedepigraphicaltext>",
		52: "<hb-sp-text-before-correction>", // // non-TLG; Text Before Correction; not an erasure; epigraphical (but in INS0150)
		53: "❨",                              // non-TLG; Parenthesis modern Punctuation (but in INS0150)
	}
)

func ReplaceLeftSquareBrackets(ttc string) string {
	// Format [ markup

	// a square bracket stands for the various kinds of brackets which appear in a text: a
	// square bracket by itself ([) indicates a normal square bracket ([]); a square bracket followed by
	// numeral one ([1) indicates a parenthesis (()); when followed by numeral two ([2), an angle
	// bracket (< >); when followed by numeral three ([3), a brace ({ }).

	ttc = lsbrouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return lsbsubstitutes(match)
	})

	return ttc
}

func lsbsubstitutes(match string) string {
	const (
		UHP = `<hgb-build-error>［%s<hgb-build-error>`
	)

	// match looks like "[14": drop that initial character
	match = strings.TrimPrefix(match, "[")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := lsbrsubsmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("lsbsubstitutes()\t%s", match))
	}
	return substitute
}
