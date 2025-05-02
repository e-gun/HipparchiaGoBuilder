//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"regexp"
	"strings"
)

const (
	GKLEMTESTDATA = `a(/bra	123337	a(/bra (fem nom/voc/acc dual) (fem nom/voc sg (attic doric aeolic))	a(/brai (fem nom/voc pl) (fem dat sg (attic doric aeolic))	a(/brais (fem dat pl)	a(/bran (fem acc sg (attic doric aeolic))	a(/bras (fem acc pl) (fem gen sg (attic doric aeolic))	a(/bra| (fem nom/voc pl) (fem dat sg (attic doric aeolic))	a(bra=n (fem gen pl (doric aeolic))	a(brw=n (fem gen pl)
a(/bruna	159444	a(/bruna (neut nom/voc/acc pl)
a(/dhma	1247181	a(/dhma (neut nom/voc/acc sg)
a(/dhn	1252572	a(/ddhn (indeclform (adverb))	a(/dhn (indeclform (adverb))	a)/ddhn (indeclform (adverb))	a)/dhn (indeclform (adverb))
a(/dicis	1394579	a(/dicis (fem nom sg)
a(/dos	1444306	a(/d' (masc voc sg)	a(/de (masc voc sg)	a(/dea (neut nom/voc/acc pl (epic ionic))	a(/dh (neut nom/voc/acc pl (attic epic doric)) (neut nom/voc/acc dual (doric aeolic))	a(/dhn (neut acc sg)	a(/doi (masc nom/voc pl)	a(/don (masc acc sg)	a(/dos (neut nom/voc/acc sg) (masc nom sg)	a(/dou (masc gen sg)	a(/dw (masc nom/voc/acc dual) (masc gen sg (doric aeolic))	a(/dwn (masc gen pl)	a(de/wn (neut gen pl (epic doric ionic aeolic))	a(dw=n (neut gen pl (attic epic doric))
a(/drunsis	1479614	a(/drunsis (fem nom sg)	a(dru/nsews (fem gen sg (attic))`
)

var (
	keywordfinder     = regexp.MustCompile(`^(.*?\t)(\d{1,})(.*?)$`)
	derivformcleaner1 = regexp.MustCompile(`\t\(.*?\)\t`)
	derivformcleaner2 = regexp.MustCompile(`(.*?)\s.*?$`)
)

func ParseGreekLemmata(entries []string) []structs.HeadwordAndForms {
	return ParseLemmata("greek", entries)
}

func ParseLemmata(lang string, entries []string) []structs.HeadwordAndForms {
	const (
		NOTIFYEVERY = 25000
	)

	//fmt.Println("ParseGreekLemmata() using GKLEMTESTDATA")
	//entries = strings.Split(GKLEMTESTDATA, "\n")

	lemmata := make([]structs.HeadwordAndForms, len(entries))

	for i, entry := range entries {
		groups := keywordfinder.FindStringSubmatch(entry)
		if len(groups) == 4 {
			lemmata[i].DictEntry = groups[1]
			lemmata[i].XREF = groups[2]
			lemmata[i].DerivFormsStr = groups[3]
		}
	}

	// cleaning the derivativeforms involves a LOSS OF INFORMATION
	// ζῳωδία           |    49601761 | ζῳωδίαϲ (fem acc pl) (fem gen sg (attic doric aeolic))  ζῳωδίᾳ (fem dat sg (attic doric aeolic))
	// becomes
	// ζῳωδία           |    49601761 | {ζῳωδίαϲ,ζῳωδίᾳ}
	// but nothing in HipparchiaServer currently needs anything other than the list of extant forms

	for i, lem := range lemmata {
		dfc := derivformcleaner1.ReplaceAllString(lem.DerivFormsStr, "\t")
		dfsegements := strings.Split(dfc, "\t")

		// now:
		//a(/brai (fem nom/voc pl) (fem dat sg (attic doric aeolic))
		//a(/brais (fem dat pl)
		//...

		for j, df := range dfsegements {
			dfsegements[j] = derivformcleaner2.ReplaceAllString(df, "$1")
		}

		// now:
		// a(/bra
		// a(/brai
		// a(/brais

		// but things still need cleaning: tab after DictEntry + one blank slice
		// 'sperata	'
		//	''
		//	'sperata'
		//	'speratae'
		//	'speratam'
		//	'sperataque'
		//	'speratarum'
		//	'speratis'
		//	'speratisne'
		//	'speratisque'

		lemmata[i].DictEntry = strings.TrimSpace(lemmata[i].DictEntry)

		for j := 0; j < len(dfsegements); j++ {
			if dfsegements[j] != "" {
				lemmata[i].DerivFormsSlc = append(lemmata[i].DerivFormsSlc, dfsegements[j])
			}
		}

		//if i%NOTIFYEVERY == 0 {
		//	fmt.Printf("ParseGreekLemmata() first pass on #%d of %d\n", i, len(lemmata))
		//	fmt.Printf("'%s'\n", lemmata[i].DictEntry)
		//	for j := 0; j < len(lemmata[i].DerivFormsSlc); j++ {
		//		fmt.Printf("\t'%s'\n", lemmata[i].DerivFormsSlc[j])
		//	}
		//}

	}

	if lang == "greek" {
		lemmata = betacodeforlemmata(lemmata)
	}

	// cannot have any empty data at time of insert: unable to encode "" into binary format for int4 (OID 23): cannot find encode plan (SQLSTATE 57014)
	// there should be only one of these
	for i, lem := range lemmata {
		if lem.DictEntry == "" {
			lemmata[i].DictEntry = " "
			lemmata[i].XREF = "0"
			lemmata[i].DerivFormsStr = " "
			lemmata[i].DerivFormsSlc = []string{" "}
		}
	}
	return lemmata
}

func betacodeforlemmata(lemmata []structs.HeadwordAndForms) []structs.HeadwordAndForms {
	// turn betacode into unicode for all of this...; look out for ValidUTF8 issues
	for i, lem := range lemmata {
		//if i%NOTIFYEVERY == 0 {
		//	fmt.Printf("ParseGreekLemmata() cleaning #%d of %d\n", i, len(lemmata))
		//}
		lemmata[i].DictEntry = strings.ToValidUTF8(generic.ConvertLCBetacode(lem.DictEntry), "")
		for j := 0; j < len(lemmata[i].DerivFormsSlc); j++ {
			lemmata[i].DerivFormsSlc[j] = strings.ToValidUTF8(generic.ConvertLCBetacode(lemmata[i].DerivFormsSlc[j]), "")
		}
		// shrink the things while we are at it...
		lemmata[i].DerivFormsStr = ""
	}
	return lemmata
}
