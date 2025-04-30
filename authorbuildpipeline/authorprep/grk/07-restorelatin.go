//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// needsrestoring is tricky; `(.*?)` could grab huge swathes; a span was unbalanced...; but will `([^█]*?)` be too small?
	needsrestoring   = regexp.MustCompile(`(<hb-fs-l-.*?>)([^█]*?)(</hb-fs-l-.*?>)`)
	latindiacritical = regexp.MustCompile(`[aeiouyAEIOU][+\\=/]`)
	lcdsubs          = map[string]string{
		"a/":  "á",
		"e/":  "é",
		"i/":  "í",
		"o/":  "ó",
		"u/":  "ú",
		"y/":  "ý",
		"A/":  "Á",
		"E/":  "É",
		"I/":  "Í",
		"O/":  "Ó",
		"U/":  "Ú",
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
	}
	restoreinvals  = []rune("αβξδεφγηι⒣κλμνοπθρϲστυϝωχψζϋ·̣")
	restoreoutvals = []rune("ABCDEFGHIJKLMNOPQRSSTUVWXYZü:?")
)

func RestoreLatinWithinGreek(ttc string) string {
	// "Ex Libro Primo" looks like "εx λibro πrimo" at this point
	// the capitals all got turned into lc Greek

	// First replacement using parseromaninsidegreek
	ttc = needsrestoring.ReplaceAllStringFunc(ttc, parseromaninsidegreek)

	// Second replacement using ldc
	ttc = needsrestoring.ReplaceAllStringFunc(ttc, latindiacriticalswaps)

	return ttc
}

func parseromaninsidegreek(match string) string {
	groups := needsrestoring.FindStringSubmatch(match)

	if len(groups) < 4 {
		fmt.Printf("parseromaninsidegreek() bad match:\n\t%q\n", match)
		return match
	}

	openTag := groups[1]      // The opening tag
	mangledRoman := groups[2] // The content to transform
	closeTag := groups[3]     // The closing tag

	// Translate characters
	transformed := strings.Map(func(r rune) rune {
		for i, char := range restoreinvals {
			if r == char {
				return restoreoutvals[i]
			}
		}
		return r
	}, mangledRoman)

	// Construct the transformed span
	span := fmt.Sprintf("%s%s%s", openTag, transformed, closeTag)

	return span
}

func latindiacriticalswaps(match string) string {
	return latindiacritical.ReplaceAllStringFunc(match, latinsubstitutes)
}

// latinsubstitutes replaces found diacritical marks with their correct Unicode characters.
func latinsubstitutes(match string) string {
	if substitute, exists := lcdsubs[match]; exists {
		return substitute
	}
	return "" // Return empty string if no match
}
