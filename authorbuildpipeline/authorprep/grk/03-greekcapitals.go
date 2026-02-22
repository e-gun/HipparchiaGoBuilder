//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

// the regex
var (
	csga = regexp.MustCompile(`\*\)\\([AHW])\|`)
	crga = regexp.MustCompile(`\*\(\\([AHW])\|`)
	csaa = regexp.MustCompile(`\*\)/([AHW])\|`)
	craa = regexp.MustCompile(`\*\(/([AHW])\|`)
	csca = regexp.MustCompile(`\*\(=([AHW])\|`)
	crca = regexp.MustCompile(`\*\(=([AHW])\|`)

	csg = regexp.MustCompile(`\*\)\\([AEIOUHW])`)
	crg = regexp.MustCompile(`\*\(\\([AEIOUHW])`)
	csa = regexp.MustCompile(`\*\)/([AEIOUHW])`)
	cra = regexp.MustCompile(`\*\(/([AEIOUHW])`)
	csc = regexp.MustCompile(`\*\)=([AEIOUHW])`)
	crc = regexp.MustCompile(`\*\(=([AEIOUHW])`)

	cs = regexp.MustCompile(`\*\)([AEIOUHWR])`)
	cr = regexp.MustCompile(`\*\(([AEIOUHWR])`)

	cg = regexp.MustCompile(`\*([AEIOUHW])\\`)
	ca = regexp.MustCompile(`\*([AEIOUHW])/`)
	cc = regexp.MustCompile(`\*([AEIOUHW])=`)

	cad = regexp.MustCompile(`\*([AHW])\|`)

	cpt = regexp.MustCompile(`\*([A-Z])`)
)

// the lookup maps
var (
	capitalsigmassubsitutesmap = map[string]string{
		"3": "\u03f9",
	}
	capitalsmoothgraveadscriptmap = map[string]string{
		"A": "ᾊ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾚ",
		"W": "ᾪ",
	}
	capitalroughgraveadscriptmap = map[string]string{
		"A": "ᾋ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾛ",
		"W": "ᾫ",
	}
	capitalsmoothacuteadscriptmap = map[string]string{
		"A": "ᾌ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾜ",
		"W": "ᾬ",
	}
	capitalroughacuteadscriptmap = map[string]string{
		"A": "ᾍ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾝ",
		"W": "ᾭ",
	}
	capitalsmoothcircumflexadscriptmap = map[string]string{
		"A": "ᾎ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾞ",
		"W": "ᾮ",
	}
	capitalroughcircumflexadscriptmap = map[string]string{
		"A": "ᾏ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾟ",
		"W": "ᾯ",
	}
	capitalsmoothgravemap = map[string]string{
		"A": "Ἂ",
		"E": "Ἒ",
		"I": "Ἲ",
		"O": "Ὂ",
		"U": "",
		"H": "Ἢ",
		"W": "Ὢ",
	}
	capitalroughgravemap = map[string]string{
		"A": "Ἃ",
		"E": "Ἓ",
		"I": "Ἳ",
		"O": "Ὃ",
		"U": "Ὓ",
		"H": "Ἣ",
		"W": "Ὣ",
	}
	capitalsmoothacutemap = map[string]string{
		"A": "Ἄ",
		"E": "Ἔ",
		"I": "Ἴ",
		"O": "Ὄ",
		"U": "",
		"H": "Ἤ",
		"W": "Ὤ",
	}
	capitalroughacutemap = map[string]string{
		"A": "Ἅ",
		"E": "Ἕ",
		"I": "Ἵ",
		"O": "Ὅ",
		"U": "Ὕ",
		"H": "Ἥ",
		"W": "Ὥ",
	}
	capitalsmoothcircumflexmap = map[string]string{
		"A": "Ἆ",
		"E": "",
		"I": "Ἶ",
		"O": "",
		"U": "",
		"H": "Ἦ",
		"W": "Ὦ",
	}
	capitalroughcircumflexmap = map[string]string{
		"A": "Ἇ",
		"E": "",
		"I": "Ἷ",
		"O": "",
		"U": "Ὗ",
		"H": "Ἧ",
		"W": "Ὧ",
	}
	capitalsmoothmap = map[string]string{
		"A": "Ἀ",
		"E": "Ἐ",
		"I": "Ἰ",
		"O": "Ὀ",
		"U": "",
		"H": "Ἠ",
		"W": "Ὠ",
		"R": "Ρ",
	}
	capitalroughmap = map[string]string{
		"A": "Ἁ",
		"E": "Ἑ",
		"I": "Ἱ",
		"O": "Ὁ",
		"U": "Ὑ",
		"H": "Ἡ",
		"W": "Ὡ",
		"R": "Ῥ",
	}
	capitalgravemap = map[string]string{
		"A": "Α",
		"E": "Ε",
		"I": "Ι",
		"O": "Ο",
		"U": "Υ",
		"H": "Η",
		"W": "Ω",
	}
	capitalacutemap = map[string]string{
		"A": "Α",
		"E": "Ε",
		"I": "Ι",
		"O": "Ο",
		"U": "Υ",
		"H": "Η",
		"W": "Ω",
	}
	capitalcircumflexmap = map[string]string{
		"A": "Α",
		"E": "Ε",
		"I": "Ι",
		"O": "Ο",
		"U": "Υ",
		"H": "Η",
		"W": "Ω",
	}
	capitaladscriptmap = map[string]string{
		"A": "ᾼ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ῌ",
		"W": "ῼ",
	}
	capitalsmap = map[string]string{
		"A": "Α",
		"B": "Β",
		"C": "Ξ",
		"D": "Δ",
		"E": "Ε",
		"F": "Φ",
		"G": "Γ",
		"H": "Η",
		"I": "Ι",
		"J": "⒣", // need the unused J because Roman characters are present; see IG I3 1-2 [1-500 501-1517]) - 370: fr d-f, line 69 - Jέ[κτει Jεμέραι -- ⒣εμέραι is what we want
		"K": "Κ",
		"L": "Λ",
		"M": "Μ",
		"N": "Ν",
		"O": "Ο",
		"P": "Π",
		"Q": "Θ",
		"R": "Ρ",
		"S": "Ϲ",
		"T": "Τ",
		"U": "Υ",
		"V": "Ϝ",
		"W": "Ω",
		"X": "Χ",
		"Y": "Ψ",
		"Z": "Ζ",
	}
)

// notice the lack of error checking in the map lookups: we are interested in panicing instead of warning...

func ConvertGreekCapitals(betacode string) string {
	// needs to be done in order of length of regex string
	// capital + breathing + accent + adscript
	// '*)/W|ETO GOU=N KAI\ O( *DIONU/SIOS...'
	unicode := csga.ReplaceAllStringFunc(betacode, capitalsmoothgraveadscript)
	unicode = crga.ReplaceAllStringFunc(unicode, capitalroughgraveadscript)
	unicode = csaa.ReplaceAllStringFunc(unicode, capitalsmoothacuteadscript)
	unicode = craa.ReplaceAllStringFunc(unicode, capitalroughacuteadscript)
	unicode = csca.ReplaceAllStringFunc(unicode, capitalsmoothcircumflexadscript)
	unicode = crca.ReplaceAllStringFunc(unicode, capitalroughcircumflexadscript)

	unicode = csg.ReplaceAllStringFunc(unicode, capitalsmoothgrave)
	unicode = crg.ReplaceAllStringFunc(unicode, capitalroughgrave)
	unicode = csa.ReplaceAllStringFunc(unicode, capitalsmoothacute)
	unicode = cra.ReplaceAllStringFunc(unicode, capitalroughacute)
	unicode = csc.ReplaceAllStringFunc(unicode, capitalsmoothcircumflex)
	unicode = crc.ReplaceAllStringFunc(unicode, capitalroughcircumflex)

	unicode = cs.ReplaceAllStringFunc(unicode, capitalsmooth)
	unicode = cr.ReplaceAllStringFunc(unicode, capitalrough)

	unicode = cg.ReplaceAllStringFunc(unicode, capitalgrave)
	unicode = ca.ReplaceAllStringFunc(unicode, capitalacute)
	unicode = cc.ReplaceAllStringFunc(unicode, capitalcircumflex)

	unicode = cad.ReplaceAllStringFunc(unicode, capitaladscript)

	if !LUNATE {
		sig := regexp.MustCompile(`[*]S([1-3]){0,1}`)
		unicode = sig.ReplaceAllStringFunc(unicode, capitalsigmassubsitutes)
	} else {
		sig := regexp.MustCompile(`[*]S[1-3]{0,1}`)
		unicode = sig.ReplaceAllString(unicode, "\u03f9")
	}

	unicode = cpt.ReplaceAllStringFunc(unicode, capitals)

	return unicode
}

func capitalsigmassubsitutes(match string) string {
	match = match[len(match)-1:]
	substitute, ok := capitalsigmassubsitutesmap[match]
	if !ok {
		substitute = "Σ"
	}
	return substitute
}

//
// adscripts
//

func capitalsmoothgraveadscript(match string) string {
	// in something like "*)/W|" you want the penultimate character
	match = match[len(match)-2 : len(match)-1]
	m, ok := capitalsmoothgraveadscriptmap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalsmoothgraveadscript() failed for : %s", match))
	}
	return m
}

func capitalroughgraveadscript(match string) string {
	match = match[len(match)-2 : len(match)-1]
	m, ok := capitalroughgraveadscriptmap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalroughgraveadscript() failed for : %s", match))
	}
	return m
}

func capitalsmoothacuteadscript(match string) string {
	match = match[len(match)-2 : len(match)-1]
	m, ok := capitalsmoothacuteadscriptmap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalsmoothacuteadscript() failed for : %s", match))
	}
	return m
}

func capitalroughacuteadscript(match string) string {
	match = match[len(match)-2 : len(match)-1]
	m, ok := capitalroughacuteadscriptmap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalroughacuteadscript() failed for :  %s", match))
	}
	return m
}

func capitalsmoothcircumflexadscript(match string) string {
	match = match[len(match)-2 : len(match)-1]
	m, ok := capitalsmoothcircumflexadscriptmap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalsmoothcircumflexadscript() failed for : %s", match))
	}
	return m
}

func capitalroughcircumflexadscript(match string) string {
	match = match[len(match)-2 : len(match)-1]
	m, ok := capitalroughcircumflexadscriptmap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalroughcircumflexadscript() failed for : %s", match))
	}
	return m
}

func capitalsmoothgrave(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalsmoothgravemap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalsmoothgrave() failed for : %s", match))
	}
	return m
}

func capitalroughgrave(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalroughgravemap[match]
	if !ok {
		global.MSG(fmt.Sprintf("capitalroughgrave() failed for : %s", match))
	}
	return m
}

func capitalsmoothacute(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalsmoothacutemap[match]
	if !ok && WARNINGS {
		global.MSG(fmt.Sprintf("capitalsmoothacute() failed for : %s", match))
	}
	return m
}

func capitalroughacute(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalroughacutemap[match]
	if !ok && WARNINGS {
		global.MSG(fmt.Sprintf("capitalroughacute() failed for : %s", match))
	}
	return m
}

func capitalsmoothcircumflex(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalsmoothcircumflexmap[match]
	if !ok && WARNINGS {
		fmt.Println("capitalsmoothcircumflex() failed for ", match)
	}
	return m
}

func capitalroughcircumflex(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalroughcircumflexmap[match]
	if !ok && WARNINGS {
		fmt.Println("capitalroughcircumflex() failed for ", match)
	}
	return m
}

func capitalsmooth(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalsmoothmap[match]
	if !ok && WARNINGS {
		fmt.Println("capitalsmooth() failed for ", match)
	}
	return m
}

func capitalrough(match string) string {
	match = match[len(match)-1:]
	m, ok := capitalroughmap[match]
	if !ok && WARNINGS {
		fmt.Println("capitalrough() failed for ", match)
	}
	return m
}

func capitalgrave(match string) string {
	// 	this sort of thing will get requested, but it looks nuts:
	//		ΠΕΡῚ ἈΝΑΓΝΏϹΕΩϹ
	//		ΠΕΡῚ ΓΡΑΜΜΑΤΙΚΗ῀Ϲ
	//
	//	just return the capitals instead
	//
	//	NB, you are still going to get stuck with ΠΕΡΙ ἈΝΑΓΝΩϹΕΩϹ because you can't turn off
	//	the 'Ἀ' to handle this situation without producing chaos throughout every other text
	//	(unless you are willing to waste precious cycles in order to do some sort of BLANK+CAP+CAP check)

	// with something like "*(A\" the key is in the penultimate place
	match = match[len(match)-2 : len(match)-1]

	//substitutions := map[string]string{
	//	"A": "Ὰ",
	//	"E": "Ὲ",
	//	"I": "Ὶ",
	//	"O": "Ὸ",
	//	"U": "Ὺ",
	//	"H": "Ὴ",
	//	"W": "Ὼ",
	//}

	m, ok := capitalgravemap[match]
	if !ok && WARNINGS {
		fmt.Println("capitalgrave() failed for ", match)
	}
	return m
}

func capitalacute(match string) string {
	// 	this sort of thing will get requested, but it looks nuts:
	//		ΠΕΡῚ ἈΝΑΓΝΏϹΕΩϹ
	//		ΠΕΡῚ ΓΡΑΜΜΑΤΙΚΗ῀Ϲ
	//
	//	just return the capitals instead

	// with something like "*(A\" the key is in the penultimate place
	match = match[len(match)-2 : len(match)-1]

	//substitutions := map[string]string{
	//	"A": "Ά",
	//	"E": "Έ",
	//	"I": "Ί",
	//	"O": "Ό",
	//	"U": "Ύ",
	//	"H": "Ή",
	//	"W": "Ώ",
	//}

	m, ok := capitalacutemap[match]
	if !ok && WARNINGS {
		fmt.Println("capitalacute() failed for ", match)
	}
	return m
}

func capitalcircumflex(match string) string {
	// 	this sort of thing will get requested, but it looks nuts:
	//		ΠΕΡῚ ἈΝΑΓΝΏϹΕΩϹ
	//		ΠΕΡῚ ΓΡΑΜΜΑΤΙΚΗ῀Ϲ
	//
	//	just return the capitals instead

	// with something like "*(A=" the key is in the penultimate place
	match = match[len(match)-2 : len(match)-1]

	//substitutions := map[string]string{
	//	"A": "Α\u1fc0",
	//	"E": "Ε\u1fc0",
	//	"I": "Ι\u1fc0",
	//	"O": "Ο\u1fc0",
	//	"U": "Υ\u1fc0",
	//	"H": "Η\u1fc0",
	//	"W": "Ω\u1fc0",
	//}
	m, ok := capitalcircumflexmap[match]
	if !ok && WARNINGS {
		fmt.Println("capitalcircumflex() failed for ", match)
	}
	return m
}

func capitaladscript(match string) string {
	// with something like "*(A\" the key is in the penultimate place
	match = match[len(match)-2 : len(match)-1]
	m, ok := capitaladscriptmap[match]
	if !ok && WARNINGS {
		fmt.Println("capitaladscript() failed for ", match)
	}
	return m
}

func capitals(match string) string {
	// the match arrives as "*A"
	// this is all 8-bat so it is safe to just grab the last character
	match = match[len(match)-1:]
	m, ok := capitalsmap[match]
	if !ok && WARNINGS {
		fmt.Println("capitals() failed for ", match)
	}
	return m
}

func getallknowngucchars() string {
	// ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΤΥΦΧΨΩϜϹἈἉἊἋἌἍἎἏἘἙἚἛἜἝἨἩἪἫἬἭἮἯἸἹἺἻἼἽἾἿὈὉὊὋὌὍὙὛὝὟὨὩὪὫὬὭὮὯᾊᾋᾌᾍᾎᾏᾚᾛᾜᾝᾞᾟᾪᾫᾬᾭᾮᾯᾼῌῬῼ⒣
	charmaps := []map[string]string{
		capitalsigmassubsitutesmap,
		capitalsmoothgraveadscriptmap,
		capitalroughgraveadscriptmap,
		capitalsmoothacuteadscriptmap,
		capitalroughacuteadscriptmap,
		capitalsmoothcircumflexadscriptmap,
		capitalroughcircumflexadscriptmap,
		capitalsmoothgravemap,
		capitalroughgravemap,
		capitalsmoothacutemap,
		capitalroughacutemap,
		capitalsmoothcircumflexmap,
		capitalroughcircumflexmap,
		capitalsmoothmap,
		capitalroughmap,
		capitalgravemap,
		capitalacutemap,
		capitalcircumflexmap,
		capitaladscriptmap,
		capitalsmap,
	}
	var aggregate string
	for _, charmap := range charmaps {
		for _, v := range charmap {
			if v != "" {
				aggregate += v
			}
		}
	}
	rn := []rune(aggregate)
	sort.Slice(rn, func(i, j int) bool {
		return rn[i] < rn[j]
	})
	var dedup []rune
	var prev rune
	for _, char := range rn {
		if char != prev {
			dedup = append(dedup, char)
		}
		prev = char
	}
	return string(dedup)
}
