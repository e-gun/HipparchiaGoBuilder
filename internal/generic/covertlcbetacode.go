//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package generic

import (
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/grk"
)

var (
	findinitialcap = regexp.MustCompile("^([A-Z])")
	superscr       = strings.NewReplacer("1", "¹", "2", "²", "3", "³", "4", "⁴")
	lunates        = strings.NewReplacer("v", "u", "j", "i", "σ", "ϲ", "ς", "ϲ", "Σ", "Ϲ", "U", "V")
)

func ConvertLCBetacode(s string) string {
	if findinitialcap.MatchString(s) {
		s = "*" + s
	}
	// greek vowel quantities need to be stripped for the keys; they are ok in the body
	// note that HandleVowelLengths() would need to run before you call ConvertLCBetacode()
	// if you care about them...
	s = strings.ReplaceAll(s, "^", "")
	s = strings.ReplaceAll(s, "_", "")
	return grk.ConvertGreekLowers(grk.ConvertGreekCapitals(strings.ToUpper(s)))
}

func HandleVowelLengths(s string) string {
	// the non-documented long-and-short codes
	s = strings.ReplaceAll(s, "^", "\u0306")
	s = strings.ReplaceAll(s, "_", "\u0304")
	return s
}

func StripVowelLengths(s string) string {
	// the non-documented long-and-short codes
	s = strings.ReplaceAll(s, "^", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}

func SuperScriptNumbers(s string) string {
	return superscr.Replace(s)
}

func LunatesAndUV(s string) string {
	return lunates.Replace(s)
}
