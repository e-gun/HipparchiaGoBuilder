//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import (
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

const (
	GKANALYSISTESTDATA = `!a/bais	{46079165 9 a_(/bais,h(/bh	youthful prime, youth	fem dat pl (doric)}
!a/brois	{28324761 9 e)/brois,e)/bros	he-goat	masc dat pl}
!a/can	{1026524 9 a)/can,a)/gw	lead, carry, fetch, bring	aor part act neut nom/voc/acc sg}
!a/cin	{35696967 9 e)/cin,e)/ceimi2	issue	imperf ind act 3rd pl (epic)}
!a/cios	{10280715 9 a)/cios,a)/cios	counterbalancing	masc nom sg}
!a/da	{1444306 9 a(/da_,a(/dos	satiety, loathing	neut nom/voc/acc pl (doric aeolic)}
!a/dh	{1444306 9 a(/dh,a(/dos	satiety, loathing	neut nom/voc/acc pl (attic epic doric)}{1444306 9 a(/dh,a(/dos	satiety, loathing	neut nom/voc/acc dual (doric aeolic)}
!a/dhn	{1252572 9 a(/dhn,a(/dhn	to oneʼs fill	indeclform (adverb)}{1444306 9 a(/dhn,a(/dos	satiety, loathing	neut acc sg}`
)

// note that a leading "!" in this data has an unknown purpose; it is undocumented relative to the betacode standard
// 		{"!", "∙"}, // missing letter dots (!); but not defined elsewhere?

// 	an analysis line looks like this:
//		!a/sdwn {32430564 9 e)/sdwn,ei)sdi/dwmi flow into       aor ind act 3rd pl (epic doric aeolic)}{32430564 9 e)/sdwn,ei)sdi/dwmi  flow into       aor ind act 1st sg (epic)}
//
//		beluasque	{8991758 9 be_lua_s,belua	a beast	fem acc pl}{9006674 9 be_lua_s,beluus	 	fem acc pl}
//
//		word(TAB){analysis 1}{analysis 2}{...}
//
//	each inset analysis is:
//		{xrefnumber digit ancientform1,ancientform2(TAB)translation(TAB)parsinginfo}
//
//	possible_dictionary_forms is going to be JSON:
//		{'NUMBER': {
//				'headword': 'WORD',
//				'scansion': 'SCAN',
//				'xref_value': 'VAL',
//				'xref_kind': 'KIND',
//				'transl': 'TRANS',
//				'analysis': 'ANAL'
//			},
//		'NUMBER': {
//				'headword': 'WORD',
//				'scansion': 'SCAN',
//				'xref_value': 'VAL',
//				'xref_kind': 'KIND',
//				'transl': 'TRANS',
//				'analysis': 'ANAL'
//			},
//		...}
//
//	note that scansion is meaningless for the greek words
//
//	 observed_form |   xrefs   | prefixrefs |                                                                      possible_dictionary_forms
//	 βάδην         | 18068282 |            | {"1": {"scansion": "", "headword": "\u03b2\u03ac\u03b4\u03b7\u03bd", "xref_value": "18068282", "xref_kind": "9", "transl": "step by step", "analysis": "indeclform (adverb)"}}
//	 βάδιζ'        | 18070850 |            | {"1": {"scansion": "\u03b2\u03ac\u03b4\u03b9\u03b6\u03b5", "headword": "\u03b2\u03b1\u03b4\u03af\u03b6\u03c9", "xref_value": "18070850", "xref_kind": "9", "transl": "walk", "analysis": "pres imperat act 2nd sg"}, "2": {"scansion": "\u03b2\u03ac\u03b4\u03b9\u03b6\u03b5", "headword": "\u03b2\u03b1\u03b4\u03af\u03b6\u03c9", "xref_value": "18070850", "xref_kind": "9", "transl": "walk", "analysis": "imperf ind act 3rd sg (homeric ionic)"}}
//	 βάδιζε        | 18070850 |            | {"1": {"scansion": "", "headword": "\u03b2\u03b1\u03b4\u03af\u03b6\u03c9", "xref_value": "18070850", "xref_kind": "9", "transl": "walk", "analysis": "pres imperat act 2nd sg"}, "2": {"scansion": "", "headword": "\u03b2\u03b1\u03b4\u03af\u03b6\u03c9", "xref_value": "18070850", "xref_kind": "9", "transl": "walk", "analysis": "imperf ind act 3rd sg (homeric ionic)"}}
//
//	 cubabunt      | 19418734 |            | {"1": {"scansion": "cuba\u0304bunt", "headword": "cubo", "xref_value": "19418734", "xref_kind": "9", "transl": " ", "analysis": "fut ind act 3rd pl"}}
//	 cubandi       | 19418734 |            | {"1": {"scansion": "cubandi\u0304", "headword": "cubo", "xref_value": "19418734", "xref_kind": "9", "transl": " ", "analysis": "gerundive masc nom/voc pl"}, "2": {"scansion": "cubandi\u0304", "headword": "cubo", "xref_value": "19418734", "xref_kind": "9", "transl": " ", "analysis": "gerundive neut gen sg"}, "3": {"scansion": "cubandi\u0304", "headword": "cubo", "xref_value": "19418734", "xref_kind": "9", "transl": " ", "analysis": "gerundive masc gen sg"}}

// 	most end with '}', but some end with a bracketed number or numbers
//	w)ph/sasqai	{49798258 9 a)ph/sasqai,a)po/-h(/domai	swād-	aor inf mid (ionic)}[12711488]
//	a)mfiperih|w/rhntai	{3398393 9 a)mfiperi+h|w/rhntai,a)mfi/,peri/-ai)wre/w	lift up	perf ind mp 3rd pl}[6238652][88377399]

// 	hunting down the bracketrefs will get you to things like: ἀπό and ἀμφί
//	that is, these mark words that have a prefix and THAT information can be used as part of the solution to the comma conundrum
//	namely, ὑπό,ἀνά,ἀπό-νέω needs to be recomposed by attending to the commas, but the commas are not only associated with prefixes
//	cf. ἠχήϲαϲα: <possibility_1>ἠχθόμαν, ἄχθομαι<xref_value>19480186</xref_value>
//
//	the number of items in bracketrefs corresponds to the number of prefix checks you will need to make to recompose the verb

var (
	analysisfinder1 = regexp.MustCompile(`^(.*?)\t(\{.*.?})$`)
	analysisfinder2 = regexp.MustCompile(`(\d+)\s(\d)\s(.*?)\t(.*?)\t(.*?)(}|)$`) // {xrefnumber digit ancientform1,ancientform2(TAB)translation(TAB)parsinginfo}
)

func ParseGreekAnalyses(entries []string) []structs.GramAnalysis {
	return ParseAnalyses("greek", entries)
}

func ParseAnalyses(lang string, entries []string) []structs.GramAnalysis {
	const (
		NOTIFYEVERY = 100000
	)

	// fmt.Println("ParseAnalyses() using GKANALYSISTESTDATA")
	//entries = strings.Split(GKANALYSISTESTDATA, "\n")

	gramanal := make([]structs.GramAnalysis, len(entries))

	for i, entry := range entries {
		groups := analysisfinder1.FindStringSubmatch(entry)
		if len(groups) == 3 {
			gramanal[i].Observed = groups[1]
			gramanal[i].Scratchpad = groups[2]
		}
	}
	for i := range gramanal {
		poss := strings.Split(gramanal[i].Scratchpad, "}{")
		var collectedpossibilities []structs.MorphPossib

		for _, pos := range poss {
			groups := analysisfinder2.FindStringSubmatch(pos)
			if len(groups) == 7 {
				var mp structs.MorphPossib
				mp.Xrefval = groups[1]
				mp.Xrefkind = groups[2]

				mp.Transl = groups[4]
				mp.Analysis = groups[5]

				hw := strings.Split(groups[3], ",") // might be "ἅλλοντο,ἅλλομαι" vel sim; needs more processing
				mp.Headwd = hw[len(hw)-1]           // take the last item...
				gramanal[i].XREFS = append(gramanal[i].XREFS, mp.Xrefval)
				collectedpossibilities = append(collectedpossibilities, mp)
			}
		}

		//if i%NOTIFYEVERY == 0 {
		//	fmt.Printf("ParseGreekAnalyses() first pass on #%d of %d\n", i, len(gramanal))
		//	// fmt.Println(collectedpossibilities)
		//}

		gramanal[i].Possibilities = collectedpossibilities
	}

	if lang == "greek" {
		gramanal = betacodeforanalyses(gramanal)
	}

	return gramanal
}

func betacodeforanalyses(gramanal []structs.GramAnalysis) []structs.GramAnalysis {
	// turn betacode into unicode for all of this...; look out for ValidUTF8 issues
	for i, ga := range gramanal {
		//if i%NOTIFYEVERY == 0 {
		//	fmt.Printf("ParseGreekAnalyses() cleaning #%d of %d\n", i, len(gramanal))
		//}
		gramanal[i].Observed = strings.ToValidUTF8(generic.ConvertLCBetacode(ga.Observed), "") // no vowel length info should be in a headword
		gramanal[i].XREFS = generic.Unique(gramanal[i].XREFS)
		for j := 0; j < len(gramanal[i].Possibilities); j++ {
			// allow vowel length info; but "a_(/bais" will choke because there is no test for short-a + "(/"
			// p := generic.HandleVowelLengths(gramanal[i].Possibilities[j].Headwd)
			p := generic.SuperScriptNumbers(gramanal[i].Possibilities[j].Headwd)
			gramanal[i].Possibilities[j].Headwd = strings.ToValidUTF8(generic.ConvertLCBetacode(p), "")
			// shrink the things while we are at it...
			gramanal[i].Scratchpad = ""
		}
	}
	return gramanal
}
