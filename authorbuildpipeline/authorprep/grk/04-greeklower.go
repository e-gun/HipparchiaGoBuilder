//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

import (
	"fmt"
	"regexp"
	"sort"
)

// the regex
var (
	// lowercase + breathing + accent + subscript
	lsga = regexp.MustCompile(`([AHW])\)\\\|`)
	lrga = regexp.MustCompile(`([AHW])\(\\\|`)
	lsaa = regexp.MustCompile(`([AHW])\)/\|`)
	lraa = regexp.MustCompile(`([AHW])\(/\|`)
	lsca = regexp.MustCompile(`([AHW])\)=\|`)
	lrca = regexp.MustCompile(`([AHW])\(=\|`)

	// lowercase + breathing + accent
	lsg = regexp.MustCompile(`([AEIOUHW])\)\\`)
	lrg = regexp.MustCompile(`([AEIOUHW])\(\\`)
	lsa = regexp.MustCompile(`([AEIOUHW])\)/`)
	lra = regexp.MustCompile(`([AEIOUHW])\(/`)
	lsc = regexp.MustCompile(`([AEIOUHW])\)=`)
	lrc = regexp.MustCompile(`([AEIOUHW])\(=`)

	// lowercase + accent + subscript
	lgas = regexp.MustCompile(`([AHW])\\\|`)
	laas = regexp.MustCompile(`([AHW])/\|`)
	lcas = regexp.MustCompile(`([AHW])=\|`)

	// lowercase + breathing + subscript
	lss = regexp.MustCompile(`([AHW])\)\|`)
	lrs = regexp.MustCompile(`([AHW])\(\|`)

	// lowercase + accent + diaresis
	lgd = regexp.MustCompile(`([IU])\\\+`)
	lad = regexp.MustCompile(`([IU])/\+`)
	lcd = regexp.MustCompile(`([U])=\+`)

	// lowercase + breathing
	lsb  = regexp.MustCompile(`([AEIOUHWR])\)`)
	lrb  = regexp.MustCompile(`([AEIOUHWR])\([^0-9]`)
	lrbe = regexp.MustCompile(`([AEIOUHWR])\($`)

	// lowercase + accent
	lga = regexp.MustCompile(`([AEIOUHW])\\`)
	laa = regexp.MustCompile(`([AEIOUHW])/`)
	lca = regexp.MustCompile(`([AEIOUHW])=`)

	lcad = regexp.MustCompile(`([AHW])\|`)
)

// the lookup maps; note that failures to lookup should be impossible and the checks out to be irrelevant...
var (
	lowercasesmoothgravesubscriptmap = map[string]string{
		"A": "ᾂ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾒ",
		"W": "ᾢ",
	}
	lowercaseroughgravesubscriptmap = map[string]string{
		"A": "ᾃ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾓ",
		"W": "ᾣ",
	}
	lowercasesmoothacutesubscriptmap = map[string]string{
		"A": "ᾄ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾔ",
		"W": "ᾤ",
	}
	lowercaseroughacutesubscriptmap = map[string]string{
		"A": "ᾅ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾕ",
		"W": "ᾥ",
	}
	lowercasesmoothcircumflexsubscriptmap = map[string]string{
		"A": "ᾆ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾖ",
		"W": "ᾦ",
	}
	lowercaseroughcircumflexsubscriptmap = map[string]string{
		"A": "ᾇ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ᾗ",
		"W": "ᾧ",
	}
	lowercasesmoothgravemap = map[string]string{
		"A": "ἂ",
		"E": "ἒ",
		"I": "ἲ",
		"O": "ὂ",
		"U": "ὒ",
		"H": "ἢ",
		"W": "ὢ",
	}
	lowercaseroughgravemap = map[string]string{
		"A": "ἃ",
		"E": "ἓ",
		"I": "ἳ",
		"O": "ὃ",
		"U": "ὓ",
		"H": "ἣ",
		"W": "ὣ",
	}
	lowercasesmoothacutemap = map[string]string{
		"A": "ἄ",
		"E": "ἔ",
		"I": "ἴ",
		"O": "ὄ",
		"U": "ὔ",
		"H": "ἤ",
		"W": "ὤ",
	}
	lowercaseroughacutemap = map[string]string{
		"A": "ἅ",
		"E": "ἕ",
		"I": "ἵ",
		"O": "ὅ",
		"U": "ὕ",
		"H": "ἥ",
		"W": "ὥ",
	}
	lowercasesmoothcircumflexmap = map[string]string{
		"A": "ἆ",
		"E": "ἐ͂",
		"I": "ἶ",
		"O": "ὀ͂",
		"U": "ὖ",
		"H": "ἦ",
		"W": "ὦ",
	}
	lowercaseroughcircumflexmap = map[string]string{
		"A": "ἇ",
		"E": "ἑ͂",
		"I": "ἷ",
		"O": "ὁ͂",
		"U": "ὗ",
		"H": "ἧ",
		"W": "ὧ",
	}
	lowercasegravesubmap = map[string]string{
		"A": "ᾲ",
		"H": "ῂ",
		"W": "ῲ",
	}
	lowercaseacutedsubmap = map[string]string{
		"A": "ᾴ",
		"H": "ῄ",
		"W": "ῴ",
	}
	lowercasesircumflexsubmap = map[string]string{
		"A": "ᾷ",
		"H": "ῇ",
		"W": "ῷ",
	}
	lowercasesmoothsubmap = map[string]string{
		"A": "ᾀ",
		"H": "ᾐ",
		"W": "ᾠ",
	}
	lowercaseroughsubmap = map[string]string{
		"A": "ᾁ",
		"H": "ᾑ",
		"W": "ᾡ",
	}
	lowercasegravediaresismap = map[string]string{
		"A": "",
		"E": "",
		"I": "ῒ",
		"O": "",
		"U": "ῢ",
		"H": "",
		"W": "",
	}
	lowercaseacutediaresismap = map[string]string{
		"A": "",
		"E": "",
		"I": "ΐ",
		"O": "",
		"U": "ΰ",
		"H": "",
		"W": "",
	}
	lowercasesircumflexdiaresismap = map[string]string{
		"A": "",
		"E": "",
		"I": "",
		"O": "",
		"U": "ῧ",
		"H": "",
		"W": "",
	}
	lowercasesmoothmap = map[string]string{
		"A": "ἀ",
		"E": "ἐ",
		"I": "ἰ",
		"O": "ὀ",
		"U": "ὐ",
		"H": "ἠ",
		"W": "ὠ",
		"R": "ῤ",
	}
	lowercaseroughmap = map[string]string{
		"A": "ἁ",
		"E": "ἑ",
		"I": "ἱ",
		"O": "ὁ",
		"U": "ὑ",
		"H": "ἡ",
		"W": "ὡ",
		"R": "ῥ",
	}
	lowercasegravemap = map[string]string{
		"A": "ὰ",
		"E": "ὲ",
		"I": "ὶ",
		"O": "ὸ",
		"U": "ὺ",
		"H": "ὴ",
		"W": "ὼ",
	}
	lowercaseacutemap = map[string]string{
		"A": "ά",
		"E": "έ",
		"I": "ί",
		"O": "ό",
		"U": "ύ",
		"H": "ή",
		"W": "ώ",
	}
	lowercascircumflexmap = map[string]string{
		"A": "ᾶ",
		"E": "\u03b5\u0342",
		"I": "ῖ",
		"O": "\u03bf\u0342",
		"U": "ῦ",
		"H": "ῆ",
		"W": "ῶ",
	}
	lowercasediaresismap = map[string]string{
		"A": "",
		"E": "",
		"I": "ϊ",
		"O": "",
		"U": "ϋ",
		"H": "",
		"W": "",
	}
	lowercasesubscriptmap = map[string]string{
		"A": "ᾳ",
		"E": "",
		"I": "",
		"O": "",
		"U": "",
		"H": "ῃ",
		"W": "ῳ",
	}
	lowercasesmap = map[string]string{
		"A": "α",
		"B": "β",
		"C": "ξ",
		"D": "δ",
		"E": "ε",
		"F": "φ",
		"G": "γ",
		"H": "η",
		"I": "ι",
		"J": "J",
		"K": "κ",
		"L": "λ",
		"M": "μ",
		"N": "ν",
		"O": "ο",
		"P": "π",
		"Q": "θ",
		"R": "ρ",
		"S": "ϲ",
		"T": "τ",
		"U": "υ",
		"V": "ϝ",
		"W": "ω",
		"X": "χ",
		"Y": "ψ",
		"Z": "ζ",
	}
)

func ConvertGreekLowers(betacode string) string {
	// lowercase + breathing + accent + subscript
	unicode := lsga.ReplaceAllStringFunc(betacode, lowercasesmoothgravesubscript)
	unicode = lrga.ReplaceAllStringFunc(unicode, lowercaseroughgravesubscript)
	unicode = lsaa.ReplaceAllStringFunc(unicode, lowercasesmoothacutesubscript)
	unicode = lraa.ReplaceAllStringFunc(unicode, lowercaseroughacutesubscript)
	unicode = lsca.ReplaceAllStringFunc(unicode, lowercasesmoothcircumflexsubscript)
	unicode = lrca.ReplaceAllStringFunc(unicode, lowercaseroughcircumflexsubscript)

	// lowercase + breathing + accent
	unicode = lsg.ReplaceAllStringFunc(unicode, lowercasesmoothgrave)
	unicode = lrg.ReplaceAllStringFunc(unicode, lowercaseroughgrave)
	unicode = lsa.ReplaceAllStringFunc(unicode, lowercasesmoothacute)
	unicode = lra.ReplaceAllStringFunc(unicode, lowercaseroughacute)
	unicode = lsc.ReplaceAllStringFunc(unicode, lowercasesmoothcircumflex)
	unicode = lrc.ReplaceAllStringFunc(unicode, lowercaseroughcircumflex)

	// lowercase + accent + subscript
	unicode = lgas.ReplaceAllStringFunc(unicode, lowercasegravesub)
	unicode = laas.ReplaceAllStringFunc(unicode, lowercaseacutedsub)
	unicode = lcas.ReplaceAllStringFunc(unicode, lowercasesircumflexsub)

	// lowercase + breathing + subscript
	unicode = lss.ReplaceAllStringFunc(unicode, lowercasesmoothsub)
	unicode = lrs.ReplaceAllStringFunc(unicode, lowercaseroughsub)

	// lowercase + accent + diaresis
	unicode = lgd.ReplaceAllStringFunc(unicode, lowercasegravediaresis)
	unicode = lad.ReplaceAllStringFunc(unicode, lowercaseacutediaresis)
	unicode = lcd.ReplaceAllStringFunc(unicode, lowercasesircumflexdiaresis)

	// lowercase + breathing
	unicode = lsb.ReplaceAllStringFunc(unicode, lowercasesmooth)
	unicode = lrb.ReplaceAllStringFunc(unicode, lowercaserough)
	unicode = lrbe.ReplaceAllStringFunc(unicode, lowercaseroughend)

	// lowercase + accent
	unicode = lga.ReplaceAllStringFunc(unicode, lowercasegrave)
	unicode = laa.ReplaceAllStringFunc(unicode, lowercaseacute)
	unicode = lca.ReplaceAllStringFunc(unicode, lowercascircumflex)

	// lowercase + diaresis
	ld := regexp.MustCompile(`([IU])\+`)

	unicode = ld.ReplaceAllStringFunc(unicode, lowercasediaresis)

	// lowercase + subscript
	unicode = lcad.ReplaceAllStringFunc(unicode, lowercasesubscript)

	if !LUNATE {
		sig := regexp.MustCompile(`S([1-3]){0,1}`)
		unicode = sig.ReplaceAllStringFunc(unicode, lowercasesigmassubsitutes)
		straypunct := `\<\>\{\}\[\]\(\)⟨⟩₍₎\.\?\!⌉⎜͙✳※¶§͜﹖→𐄂𝕔;:ˈ＇,‚‛‘“”„·‧∣`
		combininglowerdot := `\u0323`
		boundaries := `([` + combininglowerdot + straypunct + `\s]|$)`
		terminalsigma := regexp.MustCompile(`σ` + boundaries)
		unicode = terminalsigma.ReplaceAllString(unicode, `ς$1`)
	} else {
		sig := regexp.MustCompile(`S[1-3]{0,1}`)
		unicode = sig.ReplaceAllString(unicode, `ϲ`)
	}

	// lowercases
	lap := regexp.MustCompile(`([A-Z])`)
	unicode = lap.ReplaceAllStringFunc(unicode, lowercases)

	return unicode
}

func lowercasesigmassubsitutes(match string) string {
	substitutions := map[string]string{
		"1": "σ",
		"2": "ς",
		"3": "ϲ",
	}

	substitute, ok := substitutions[match]
	if !ok {
		substitute = "σ"
	}

	return substitute
}

func lowercasesmoothgravesubscript(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothgravesubscriptmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmoothgravesubscript() failed for ", match)
	}
	return m
}

func lowercaseroughgravesubscript(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughgravesubscriptmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughgravesubscript() failed for ", match)
	}
	return m
}

func lowercasesmoothacutesubscript(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothacutesubscriptmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmoothacutesubscript() failed for ", match)
	}
	return m
}

func lowercaseroughacutesubscript(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughacutesubscriptmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughacutesubscript() failed for ", match)
	}
	return m
}

func lowercasesmoothcircumflexsubscript(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothcircumflexsubscriptmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmoothcircumflexsubscript() failed for ", match)
	}
	return m
}

func lowercaseroughcircumflexsubscript(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughcircumflexsubscriptmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughcircumflexsubscript() failed for ", match)
	}
	return m
}

func lowercasesmoothgrave(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothgravemap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmoothgrave() failed for ", match)
	}
	return m
}

func lowercaseroughgrave(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughgravemap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughgrave() failed for ", match)
	}
	return m
}

func lowercasesmoothacute(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothacutemap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmoothacute() failed for ", match)
	}
	return m
}

func lowercaseroughacute(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughacutemap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughacute() failed for ", match)
	}
	return m
}

func lowercasesmoothcircumflex(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothcircumflexmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmoothcircumflex() failed for ", match)
	}
	return m
}

func lowercaseroughcircumflex(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughcircumflexmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughcircumflex() failed for ", match)
	}
	return m
}

func lowercasegravesub(match string) string {
	match = match[0:1]
	m, ok := lowercasegravesubmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasegravesub() failed for ", match)
	}
	return m
}

func lowercaseacutedsub(match string) string {
	match = match[0:1]
	m, ok := lowercaseacutedsubmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseacutedsub() failed for ", match)
	}
	return m
}

func lowercasesircumflexsub(match string) string {
	match = match[0:1]
	m, ok := lowercasesircumflexsubmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesircumflexsub() failed for ", match)
	}
	return m
}

func lowercasesmoothsub(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothsubmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmoothsub() failed for ", match)
	}
	return m
}

func lowercaseroughsub(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughsubmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughsub() failed for ", match)
	}
	return m
}

func lowercasegravediaresis(match string) string {
	match = match[0:1]
	m, ok := lowercasegravediaresismap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasegravediaresis() failed for ", match)
	}
	return m
}

func lowercaseacutediaresis(match string) string {
	match = match[0:1]
	m, ok := lowercaseacutediaresismap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseacutediaresis() failed for ", match)
	}
	return m
}

func lowercasesircumflexdiaresis(match string) string {
	match = match[0:1]
	m, ok := lowercasesircumflexdiaresismap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesircumflexdiaresis() failed for ", match)
	}
	return m
}

func lowercasesmooth(match string) string {
	match = match[0:1]
	m, ok := lowercasesmoothmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesmooth() failed for ", match)
	}
	return m
}

func lowercaserough(match string) string {
	// "([AEIOUHWR])\([^0-9]" will capture one extra character...
	// note that ⟨AI(⟩ will fail w/out runes here because that '⟩' is not ascii
	rns := []rune(match)
	tail := rns[2:3]
	match = string(rns[0:1])
	m, ok := lowercaseroughmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaserough() failed for ", match)
	}
	return m + string(tail)
}

// lowercaseroughend - for when lowercaserough() does not have that one extra character...
func lowercaseroughend(match string) string {
	match = match[0:1]
	m, ok := lowercaseroughmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseroughend() failed for ", match)
	}
	return m
}

func lowercasegrave(match string) string {
	match = match[0:1]
	m, ok := lowercasegravemap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasegrave() failed for ", match)
	}
	return m
}

func lowercaseacute(match string) string {
	match = match[0:1]
	m, ok := lowercaseacutemap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercaseacute() failed for ", match)
	}
	return m
}

func lowercascircumflex(match string) string {
	match = match[0:1]
	m, ok := lowercascircumflexmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercascircumflex() failed for ", match)
	}
	return m
}

func lowercasediaresis(match string) string {
	match = match[0:1]
	m, ok := lowercasediaresismap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasediaresis() failed for ", match)
	}
	return m
}

func lowercasesubscript(match string) string {
	match = match[0:1]
	m, ok := lowercasesubscriptmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercasesubscript() failed for ", match)
	}
	return m
}

func lowercases(match string) string {
	// match = match[0:1]
	m, ok := lowercasesmap[match]
	if WARNINGS && !ok {
		fmt.Println("lowercases() failed for ", match)
	}
	return m
}

func getallknownglcchars() string {
	// J͂ΐάέήίΰαβγδεζηθικλμνξοπρτυφχψωϊϋόύώϝϲἀἁἂἃἄἅἆἇἐἑἒἓἔἕἠἡἢἣἤἥἦἧἰἱἲἳἴἵἶἷὀὁὂὃὄὅὐὑὒὓὔὕὖὗὠὡὢὣὤὥὦὧὰὲὴὶὸὺὼᾀᾁᾂᾃᾄᾅᾆᾇᾐᾑᾒᾓᾔᾕᾖᾗᾠᾡᾢᾣᾤᾥᾦᾧᾲᾳᾴᾶᾷῂῃῄῆῇῒῖῢῤῥῦῧῲῳῴῶῷ
	charmaps := []map[string]string{
		lowercasesmoothgravesubscriptmap,
		lowercaseroughgravesubscriptmap,
		lowercasesmoothacutesubscriptmap,
		lowercaseroughacutesubscriptmap,
		lowercasesmoothcircumflexsubscriptmap,
		lowercaseroughcircumflexsubscriptmap,
		lowercasesmoothgravemap,
		lowercaseroughgravemap,
		lowercasesmoothacutemap,
		lowercaseroughacutemap,
		lowercasesmoothcircumflexmap,
		lowercaseroughcircumflexmap,
		lowercasegravesubmap,
		lowercaseacutedsubmap,
		lowercasesircumflexsubmap,
		lowercasesmoothsubmap,
		lowercaseroughsubmap,
		lowercasegravediaresismap,
		lowercaseacutediaresismap,
		lowercasesircumflexdiaresismap,
		lowercasesmoothmap,
		lowercaseroughmap,
		lowercasegravemap,
		lowercaseacutemap,
		lowercascircumflexmap,
		lowercasediaresismap,
		lowercasesubscriptmap,
		lowercasesmap,
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
