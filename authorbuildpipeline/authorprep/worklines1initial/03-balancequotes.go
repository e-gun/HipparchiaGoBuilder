//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

import (
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
)

// unfortunately it looks like \w does not in fact work as a means of catching Greek
// brute force can leverage generic.GUC and GLC

var (
	smartsgquotepatt0 = regexp.MustCompile(`([` + generic.ALLGKANDLAT + `.,;])‘([\W])`)
	smartsgquotepatt1 = regexp.MustCompile(`([\W])’(` + generic.ALLGKANDLAT + `)`)
	smartsgquotepatt2 = regexp.MustCompile(`([aeiouαειουηωᾳῃῳᾶῖῦῆῶάέίόύήώὰὲὶὸὺὴὼἂἒἲὂὒἢὢᾃᾓᾣᾂᾒᾢ]\s)‘(\w)`)
	smartdbquotepatt1 = regexp.MustCompile(`(\W)”(` + generic.ALLGKANDLAT + `)`)
	smartdbquotepatt2 = regexp.MustCompile(`([` + generic.ALLGKANDLAT + `.,;])“([\W])`)
	smartpmquotepatt1 = regexp.MustCompile(`([` + generic.ALLGKANDLAT + `.,])‵([\W])`)
	smartpmquotepatt2 = regexp.MustCompile(`(\W)′(` + generic.ALLGKANDLAT + `)`)
	smartenfullwidth1 = regexp.MustCompile(`(\W)＇([` + generic.ALLGKANDLAT + `])`) // see very first of earlybirds
	smartenfullwidth2 = regexp.MustCompile(`([` + generic.ALLGKANDLAT + `.,;?:·])＇(\W)`)

	// 	smartsgquotepatt0 = regexp.MustCompile(`([\w.,;·:])‘([\W])`)
	//	smartsgquotepatt1 = regexp.MustCompile(`([\W])’(\w)`)
	//	smartsgquotepatt2 = regexp.MustCompile(`([aeiouαειουηωᾳῃῳᾶῖῦῆῶάέίόύήώὰὲὶὸὺὴὼἂἒἲὂὒἢὢᾃᾓᾣᾂᾒᾢ]\s)‘(\w)`)
	//	smartdbquotepatt1 = regexp.MustCompile(`(\W)”(\w)`)
	//	smartdbquotepatt2 = regexp.MustCompile(`([\w.,;·:])“([\W])`)
	//	smartpmquotepatt1 = regexp.MustCompile(`([\w.,])‵([\W])`)
	//	smartpmquotepatt2 = regexp.MustCompile(`(\W)′(\w)`)
	// 	smartenfullwidth1 = regexp.MustCompile(`(\W)＇(\w)`) // see very first of earlybirds
	//	smartenfullwidth2 = regexp.MustCompile(`(\w)＇(\W)`)
)

// BalanceQuotes - try to educate mismatched "smart" quotes and so forth; very fiddly...
func BalanceQuotes(ttc string) string {
	if SMARTSINGLEQUOTES {
		// ttc = smartersinglequotes(ttc)
	}

	// order of execution matters
	// many corner cases

	ttc = smartenfullwidth1.ReplaceAllString(ttc, "$1‘$2")
	ttc = smartenfullwidth2.ReplaceAllString(ttc, "$1’$2")

	// do this no matter what...
	ttc = smartsgquotepatt0.ReplaceAllString(ttc, "$1’$2")

	ttc = smarterdoublequotes(ttc)

	ttc = smarterprimequotes(ttc)

	// now just get rid of the primes...
	ttc = strings.ReplaceAll(ttc, "‵", "‘")
	ttc = strings.ReplaceAll(ttc, "′", "’")

	// <hb-tabbedtext />’οὐ γὰρ ἀναιμωτί --> <hb-tabbedtext />‘οὐ γὰρ ἀναιμωτί
	ttc = strings.ReplaceAll(ttc, `>’`, `>‘`)
	return ttc
}

func smartersinglequotes(ttc string) string {
	// if you enable the next a problem arises with initial elision: ‘κείνων instead of ’κείνων

	// `(\W)’(\w)`
	ttc = smartsgquotepatt1.ReplaceAllString(ttc, "$1‘$2")

	// now we try to undo the mess we just created by looking for vowel+space+quote+char
	// the assumption is that an actual quotation will have a punctuation mark that will invalidate this check
	// Latin is a mess, and you will get too many bad matchces: De uerbo ’quiesco’
	// but the following will still be wrong: τὰ ϲπέρματα· ‘κείνων γὰρ
	// it is unfixable? how do I know that a proper quote did not just start?

	// `([aeiouαειουηωᾳῃῳᾶῖῦῆῶάέίόύήώὰὲὶὸὺὴὼἂἒἲὂὒἢὢᾃᾓᾣᾂᾒᾢ]\s)‘(\w)`
	ttc = smartsgquotepatt2.ReplaceAllString(ttc, "$1‘$2")
	return ttc
}

func smarterdoublequotes(ttc string) string {
	ttc = smartdbquotepatt1.ReplaceAllString(ttc, "$1“$2") // open
	ttc = smartdbquotepatt2.ReplaceAllString(ttc, "$1”$2") // close
	return ttc
}

func smarterprimequotes(ttc string) string {
	ttc = smartpmquotepatt1.ReplaceAllString(ttc, "$1′$2") // close
	ttc = smartpmquotepatt2.ReplaceAllString(ttc, "$1‵$2") // open
	return ttc
}
