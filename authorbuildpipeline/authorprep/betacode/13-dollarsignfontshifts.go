//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

// these are very dangerous and the source of a large number of debugging issues
// make sure that you do not destroy any hex runs like straydollars can...

// currently cannot process the following correctly: 1δίκη at 3 5035w001 and 1ἀγνώϲ at 7 5035w001

var (
	unlatungrk   = regexp.MustCompile(`&(\d{1,2})([^$█]*?)\$(\d{1,2})`)
	dollarfinder = regexp.MustCompile(`\$(\d{1,2})([^&$█]*?)(\$\d{0,2})`)
	straydollars = regexp.MustCompile(`\$(\d{1,2})([^█]*?)(\s?█)`)
	DollarSubMap = map[int][2]string{
		// $50-59 Reserved for Greek documentary papyri
		// $60-69 Reserved for Greek inscriptions
		70: {`<hb-fs-g-uncial>`, `</hb-fs-g-uncial>`},
		53: {`<hb-fs-g-hebrew>`, `</hb-fs-g-hebrew>`},
		52: {`<hb-fs-g-arabic>`, `</hb-fs-g-arabic>`},
		51: {`<hb-fs-g-demotic>`, `</hb-fs-g-demotic>`},
		50: {`<hb-fs-g-coptic>`, `</hb-fs-g-coptic>`},
		40: {`<hb-fs-g-extralarge>`, `</hb-fs-g-extralarge>`},
		30: {`<hb-fs-g-extrasmall>`, `</hb-fs-g-extrasmall>`},
		20: {`<hb-fs-g-largerthannormal>`, `</hb-fs-g-largerthannormal>`},
		18: {`<hb-fs-g-smaller>`, `</hb-fs-g-smaller>`},
		16: {`<hb-fs-g-smaller_supsc_bold>`, `</hb-fs-g-smaller_supsc_bold>`},
		15: {`<hb-fs-g-smaller_subsc>`, `</hb-fs-g-smaller_subsc>`},
		14: {`<hb-fs-g-smaller_supsc>`, `</hb-fs-g-smaller_supsc>`},
		13: {`<hb-fs-g-smaller_italic>`, `</hb-fs-g-smaller_italic>`},
		11: {`<hb-fs-g-smaller_bold>`, `</hb-fs-g-smaller_bold>`},
		10: {`<hb-fs-g-smaller>`, `</hb-fs-g-smaller>`},
		9:  {`<hb-fs-g-regular>`, `</hb-fs-g-regular>`},
		8:  {`<hb-fs-g-vertical>`, `</hb-fs-g-vertical>`},
		6:  {`<hb-fs-g-supsc_bold>`, `</hb-fs-g-supsc_bold>`},
		5:  {`<hb-fs-g-subsc>`, `</hb-fs-g-subsc>`},
		4:  {`<hb-fs-g-supsc>`, `</hb-fs-g-supsc>`},
		3:  {`<hb-fs-g-italic>`, `</hb-fs-g-italic>`},
		2:  {`<hb-fs-g-bold_italic>`, `</hb-fs-g-bold_italic>`},
		1:  {`<hb-fs-g-bold>`, `</hb-fs-g-bold>`},
		0:  {`<hb-fs-g-normal>`, `</hb-fs-g-normal>`},
	}
)

func DollarSignGreekFontMarkup(ttc string) string {
	// equivalent in python: replacegreekmarkup()
	ttc = unlatungrkswap(ttc)
	ttc = dollarfinderswap(ttc)
	ttc = straydollarswap(ttc)
	return ttc
}

func unlatungrkswap(ttc string) string {
	ttc = unlatungrk.ReplaceAllStringFunc(ttc, func(s string) string {
		// `&(\d{1,2})([^$█]*?)\$(\d{1,2})`
		// 1: the & digit
		// 2: the enclosed text
		// 3: the dollar digit

		matches := unlatungrk.FindStringSubmatch(s)
		if len(matches) == 4 {
			return removeshiftsymmetry(matches[1], matches[2], matches[3])
		}
		return s
	})
	return ttc
}

func dollarfinderswap(ttc string) string {
	ttc = dollarfinder.ReplaceAllStringFunc(ttc, func(s string) string {
		matches := dollarfinder.FindStringSubmatch(s)
		if len(matches) == 4 {
			return dollarssubstitutes(matches[1], matches[2])
		}
		return s
	})
	return ttc
}

func straydollarswap(ttc string) string {
	// adding "+ matches[3]" fixes a nybbler break: you cannot drop a "█"...

	ttc = straydollars.ReplaceAllStringFunc(ttc, func(s string) string {
		matches := straydollars.FindStringSubmatch(s)
		if len(matches) == 4 {
			return dollarssubstitutes(matches[1], matches[2]) + matches[3]
		}
		return s
	})
	return ttc
}

func removeshiftsymmetry(groupOne, groupTwo, groupThree string) string {
	// 	notice that we are requesting 'unlatin + unsmall' in:
	//
	//	*(ISTORIW=N [&10fr. 93 FHG I 215]$10 [2LE/GEI PROSISTORW=N]2
	//
	//	turn this into '&10... $'
	//
	//	otherwise you will get 'small latin span' + 'small greek span'

	// Create the base value by formatting groupOne and groupTwo
	value := fmt.Sprintf("&%s%s$", groupOne, groupTwo)

	// If groupOne is not equal to groupThree, append groupThree to the value
	if groupOne != groupThree {
		value += groupThree
	}

	return value
}

func dollarssubstitutes(match string, core string) string {
	const (
		UHP = `<hgb-build-error>＄%s<hgb-build-error>`
	)

	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	// Lookup the substitution for the given 'val'
	sub, exists := DollarSubMap[val]
	if !exists {
		sub[0] = fmt.Sprintf(UHP, match)
		sub[1] = ""
		global.MSG(fmt.Sprintf("dollarssubstitutes()\t%s", match))
	}

	return sub[0] + core + sub[1]
}
