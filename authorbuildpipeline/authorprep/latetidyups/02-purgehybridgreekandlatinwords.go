//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package latetidyups

import (
	"fmt"
	"regexp"
)

var (
	mixfinder1      = regexp.MustCompile(`(^|\s)([α-ωϲϝ])([a-z])`)
	mixfinder2      = regexp.MustCompile(`([a-z])([αβψδεφγηιξκλμνοπϙρστυωχθζϲϝ])`)
	badbreathfinder = regexp.MustCompile("([ἁἑἱὁὑἡὡῥ])([a-z])")
	punctfinder1    = regexp.MustCompile(`(\s)([α-ωϲϝ])(\.\s[A-Z])`)
	punctfinder2    = regexp.MustCompile(`(\s)([α-ωϲϝ])(\[[a-z]|\([a-z])`)
	errantbreathing = map[rune]string{
		'ἁ': "(A",
		'ἑ': "(E",
		'ἱ': "(I",
		'ὁ': "(O",
		'ὑ': "(U",
		'ἡ': "(H",
		'ὡ': "(W",
		'ῥ': "(R",
	}
	unmapgreekcap = map[rune]string{
		'α': "A",
		'β': "B",
		'ξ': "C",
		'δ': "D",
		'ε': "E",
		'φ': "F",
		'γ': "G",
		'η': "H",
		'ι': "I",
		'⒣': "J",
		'κ': "K",
		'λ': "L",
		'μ': "M",
		'ν': "N",
		'ο': "O",
		'π': "P",
		'ρ': "R",
		'ϲ': "S",
		'τ': "T",
		'υ': "U",
		'ϝ': "V",
		'ω': "W",
		'χ': "X",
		'ζ': "Z",
	}
	unmapgreeklc = map[rune]string{
		'α': "a",
		'β': "b",
		'ξ': "c",
		'δ': "d",
		'ε': "e",
		'φ': "f",
		'γ': "g",
		'η': "h",
		'ι': "i",
		'⒣': "j",
		'κ': "k",
		'λ': "l",
		'μ': "m",
		'ν': "n",
		'ο': "o",
		'π': "p",
		'ρ': "r",
		'ϲ': "s",
		'τ': "t",
		'υ': "u",
		'ϝ': "v",
		'ω': "w",
		'χ': "x",
		'ζ': "z",
	}
)

func PurgeHybrid(ttc string) string {
	//  the CHR files really need this...
	// 	LAT files that have both greek and latin can be naughty about flagging the swaps in language
	//  most of that has been caught, but the
	//	the result can be stuff like this:
	//		"domno piissimo αugusto λeone anno χϝι et ξonstantino"
	//
	//	this will look for g+latin and turn it into Latin
	//	the roman numeral problem will remain: so the real fix is to dig into this elsewhere/earlier

	// 	fixes: ξonstantino
	ttc = mixfinder1.ReplaceAllStringFunc(ttc, func(m string) string {
		return unmixer1(mixfinder1.FindStringSubmatch(m))
	})

	// 	fixes: χϝι (todo: not working)
	//ttc = mixfinder2.ReplaceAllStringFunc(ttc, func(m string) string {
	//	return unmixer2(mixfinder2.FindStringSubmatch(m))
	//})

	//	fixes: s(anctae) ῥomanae
	ttc = badbreathfinder.ReplaceAllStringFunc(ttc, func(m string) string {
		return unbreather(badbreathfinder.FindStringSubmatch(m))
	})

	//	fixes: λ. Pomponius
	ttc = punctfinder1.ReplaceAllStringFunc(ttc, func(m string) string {
		return unpunctuated(punctfinder1.FindStringSubmatch(m))
	})

	//	fixes: ξ[apitoli]o or ξ(apitolio)
	ttc = punctfinder2.ReplaceAllStringFunc(ttc, func(m string) string {
		return unpunctuated(punctfinder2.FindStringSubmatch(m))
	})

	return ttc
}

func unmixer1(matchslice []string) string {
	// fixes: ξonstantino
	greekChar := matchslice[1]
	latinChar := matchslice[2]
	transformed, ok := unmapgreekcap[[]rune(greekChar)[0]]
	if !ok {
		transformed = matchslice[1]
	}
	return " " + transformed + latinChar
}

func unmixer2(matchslice []string) string {
	// fixes: χϝι
	fmt.Println("unmixer2:", matchslice)
	latinChar1 := matchslice[1]
	greekChar := matchslice[2]
	transformed, ok := unmapgreeklc[[]rune(greekChar)[0]]
	if !ok {
		transformed = matchslice[1]
	}
	return latinChar1 + transformed
}

func unbreather(matchGroup []string) string {
	greekChar := []rune(matchGroup[1])[0]
	latinChar := matchGroup[2]
	return errantbreathing[greekChar] + latinChar
}

func unpunctuated(matchGroup []string) string {
	initial := matchGroup[1]
	greekChar := matchGroup[2]
	punct := matchGroup[3]
	transformed, ok := unmapgreekcap[[]rune(greekChar)[0]]
	if !ok {
		transformed = matchGroup[2]
	}
	return initial + transformed + punct
}
