//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

import (
	"regexp"
	"strings"
)

var (
	swapbreves           = regexp.MustCompile("(.)([\u035c\u035d\u0361])")
	fixbracketplusnumber = regexp.MustCompile(`(\d)`)
	fixbracketplusspancl = regexp.MustCompile(`([❨❩❴❵⟦⟧⟪⟫《》‹›⦅⦆₍₎⟨⟩\[\](){}])(</span>)`)
	fixbracketplusspanop = regexp.MustCompile(`(<span class="[^"]*?">)([❨❩❴❵⟦⟧⟪⟫《》‹›⦅⦆₍₎⟨⟩\[\](){}])`)
	fixemptyspans        = regexp.MustCompile(`<span class="[^"]*?"></span> `)
	findbrackets         = regexp.MustCompile(`[❨❩❴❵⟦⟧⟪⟫《》‹›⦅⦆₍₎]`)
	bracketmap           = map[string]string{
		"❨": "(",
		"❩": ")",
		"❴": "{",
		"❵": "}",
		"⟦": "[",
		"⟧": "]",
		"⦅": "(",
		"⦆": ")",
		"⸨": "(",
		"⸩": ")",
		// "₍": "(", // "[11" (enclose missing letter dots (!), expressing doubt whether there is a letter there at all)
		// "₎": ")", // "11]"
		// various angled brackets all set to "mathematical left/right angle bracket" (u+27e8, u+27e9)
		// alternately one could consider small versions instead of the full-sized versions (u+fe64, u+fe65)
		// the main issue is that "<" and ">" are being kept out of the text data because of the HTML problem
		// "⟪": "⟨", // but these are all asserted in the betacode
		// "⟫": "⟩", // but these are all asserted in the betacode
		"《": "⟨",
		"》": "⟩",
		"‹": "⟨",
		"›": "⟩",
	}
)

func BracketFixer(ttc string) string {
	// regex work that for some reason or other needs to be put off until the very last second
	// mostly this is all about brackets which are, of course, a tricky issue

	// gr2762 and chr0012 will fail the COPY TO command because of '\\'
	ttc = strings.ReplaceAll(ttc, `\\`, "⑊")

	// a format shift code like '[3' if followed by a number that is supposed to print has an intervening ` to stop the TLG parser
	// if you do this prematurely you will generate spurious codes by joining numbers that should be kept apart
	ttc = fixbracketplusnumber.ReplaceAllString(ttc, "$1")

	ttc = strings.ReplaceAll(ttc, `\(`, `(`)
	ttc = strings.ReplaceAll(ttc, `\)`, `)`)

	if SIMPLIFYBRACKETS {
		ttc = findbrackets.ReplaceAllStringFunc(ttc, bracketsimplifier)
	}

	ttc = swapitemordersuite(ttc)

	ttc = fixemptyspans.ReplaceAllString(ttc, "")

	// there are a ton of these, really...
	return ttc
}

func bracketsimplifier(match string) string {
	b, ok := bracketmap[match]
	if ok {
		return b
	}
	return match
}

func swapitemordersuite(ttc string) string {
	// 	change:
	//    <span class="latin smallerthannormal">Gnom. Vatic. 743 [</span>
	//	into:
	//	  <span class="latin smallerthannormal">Gnom. Vatic. 743 </span>[

	ttc = fixbracketplusspancl.ReplaceAllString(ttc, "$2$1")
	ttc = fixbracketplusspanop.ReplaceAllString(ttc, "$2$1")
	ttc = swapbreves.ReplaceAllString(ttc, "$2$1")
	return ttc
}
