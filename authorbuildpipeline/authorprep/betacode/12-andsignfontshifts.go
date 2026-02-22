//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

var (
	cplxgktolatfshift      = regexp.MustCompile(`&(\d{0,2})([^$█]*)\$(\d{0,2})`) // keep/drop `█`?; \d{0,2} or \d{1,2} up front?
	latshiftinsidelatshift = regexp.MustCompile(`&(\d{0,2})([^$█]*)&(\d{0,2})`)  // will catch runs up to a `█`
	basicgktolatfshift     = regexp.MustCompile(`&([^\d][^$&█]*)\$`)             // don't use "[\D]": only "[^\d]" works (!)
	wholelinelatftshift    = regexp.MustCompile(`&(\d{0,2})([^$█]*)█`)
	simplelatinshift       = regexp.MustCompile(`&(\d{1,2})([^$&█]*)&([^\d])`) // `&4g&` needs to be handled differently from `&3Poet. lyr. gr. III&4`
	vanillalatinshift      = regexp.MustCompile(`&(\d{0,2})(.*?)&`)            // can only handle `&4g&` vel sim

	latinfontmap = map[string][2]string{
		//&50-59 Reserved for Greek documentary papyri
		//&60-69 Reserved for Greek inscriptions
		// unknown span errors are typically parsing errors that fused characters/numbers
		// BUT CHR has 81 and 91 for real?
		// AND INS has 41 and 42 for real? and 83?
		// comment the next out to watch it happen...
		"91": {`<hb-fs-l-undoc91>`, `</hb-fs-l-undoc91>`},
		"90": {`<hb-fs-l-undoc90>`, `</hb-fs-l-undoc90>`},
		"83": {`<hb-fs-l-undoc83>`, `</hb-fs-l-undoc83>`},
		"82": {`<hb-fs-l-undoc82>`, `</hb-fs-l-undoc82>`},
		"81": {`<hb-fs-l-undoc81>`, `</hb-fs-l-undoc81>`},
		"43": {`<hb-fs-l-undoc43>`, `</hb-fs-l-undoc43>`},
		"42": {`<hb-fs-l-undoc42>`, `</hb-fs-l-undoc42>`},
		"41": {`<hb-fs-l-undoc41>`, `</hb-fs-l-undoc41>`},
		"20": {`<hb-fs-l-largerthannormal>`, `</hb-fs-l-largerthannormal>`},
		"15": {`<hb-fs-l-smaller_subsc>`, `</hb-fs-l-smaller_subsc>`},
		"14": {`<hb-fs-l-smaller_supsc>`, `</hb-fs-l-smaller_supsc>`},
		"13": {`<hb-fs-l-smaller_italic>`, `</hb-fs-l-smaller_italic>`},
		"10": {`<hb-fs-l-smaller>`, `</hb-fs-l-smaller>`},
		"9":  {`<hb-fs-l-normal>`, `</hb-fs-l-normal>`},
		"8":  {`<hb-fs-l-smallcapitals_italic>`, `</hb-fs-l-smallcapitals_italic>`},
		"7":  {`<hb-fs-l-smallcapitals>`, `</hb-fs-l-smallcapitals>`},
		"6":  {`<hb-fs-l-romannumerals>`, `</hb-fs-l-romannumerals>`},
		"5":  {`<hb-fs-l-subsc>`, `</hb-fs-l-subsc>`},
		"4":  {`<hb-fs-l-supsc>`, `</hb-fs-l-supsc>`},
		"3":  {`<hb-fs-l-italic>`, `</hb-fs-l-italic>`},
		"2":  {`<hb-fs-l-bold_italic>`, `</hb-fs-l-bold_italic>`},
		"1":  {`<hb-fs-l-bold>`, `</hb-fs-l-bold>`},
		"0":  {`<hb-fs-l-normal>`, `</hb-fs-l-normal>`},
	}
)

// todo: tricky and not yet handled (and also botched by the old builder): `&10 — Metaphys.$10 *G&10 5. 1010a 7$10`
// note that you go L10, G10, L10, G10
// at the moment we do L10 G L10 G.
// a smarter cplxgktolatfshift could capture the last digits and do something...

// TLG5038 annotations contain some of the most intense series of shifts and are ideal for checking the parser

// AmpersandLatinFontMarkup processes line-like segments of text and performs substitutions if '&' escapes are found.
func AmpersandLatinFontMarkup(ttc string) string {
	// ttc = SimpleLatinSpanLATE(ttc)
	// the order of execution matters because the rewrite in B might interfere with the test in A
	// ttc = marginaliashifts(ttc)
	ttc = complexgktolatinspan(ttc)
	ttc = latinsubshiftspan(ttc) // to pick up strays after complexgktolatinspan()
	ttc = simplegktolatinspan(ttc)
	ttc = severallinesoflatin(ttc)

	// this next is problem; it is not producing unwanted results...
	// ttc = SimpleLatinSpanLATE(ttc)
	// nevertheless, a Latin author does need this check to happen eventually; see the notes on SimpleLatinSpanLATE() below

	return ttc
}

// complexgktolatinspan - `greek &10stuff_in_shifted_latin_font$10 greek`
func complexgktolatinspan(ttc string) string {
	return cplxgktolatfshift.ReplaceAllStringFunc(ttc, complexmatchshifter)
}

func complexmatchshifter(match string) string {
	// 0: whole match
	// 1: shift number
	// 2: the line of text
	// 3: the next number (if it exists)

	// `&(\d{1,2})([^$█]*)\$(\d{0,2})`

	// the issue here is that you might be sent `&4g&GP&4g$`
	// that needs sub-processing

	groups := cplxgktolatfshift.FindStringSubmatch(match)
	// [&4g&U&4g&E&4g$ 4 g&U&4g&E&4g ]
	// care about: 4g&U&4g&E&4g
	reparse := groups[1] + groups[2]

	var components []string
	subunits := strings.Split(reparse, "&")
	for _, unit := range subunits {
		if unit != "" {
			toshift := "&" + unit + "&"
			components = append(components, vanillalatinshifter(toshift))
		}
	}

	// test case: TLG2319
	// &10CLEM. Str. II 130 [II 184, 14 St.%100 s. oben II 133, 12]$10

	// that `$10` is not helpful...
	// you already get `<hb-fs-l-smaller>CLEM. Str. II 130 [II 184, 14 St.%100 s. oben II 133, 12]</hb-fs-l-smaller>`
	// so we will discard it rather than adding it back on. If it is added back on, at the moment nothing catches
	// the `$10` instruction; nor, really, should anything happen: the font has been reset, and we are back in action

	//suppl := ""
	//if groups[3] != "" {
	//	fmt.Println("suppl", groups[3])
	//	suppl = "$" + groups[3]
	//}

	// fmt.Println("complexmatchshifter in: ", len(match), "\t", match)
	// fmt.Println("complexmatchshifter out:", strings.Join(components, "")+groups[3])
	return strings.Join(components, "")
}

// latinsubshiftspan - `&3pserint [1&sc. de pyramidibus]1, &3sunt Herodotus`
func latinsubshiftspan(ttc string) string {
	// the example above is interesting because you are about to unbalance some html look-and-feel...
	// nevertheless, a literal reading of the betacode instructions has italic "[" then roman up to "], "

	// complexmatchshifter() will have dealt with a lot of this already...

	return latshiftinsidelatshift.ReplaceAllStringFunc(ttc, latinsubshifter)
}

// latinsubshifter - virtually identical to complexmatchshifter(); easier to debug if `...$` and `...&` are kept apart
func latinsubshifter(match string) string {
	//  in:  &❨fgm. 49 Bergk &3Poet. lyr. gr. III&4
	//  out: <hb-fs-l-normal>❨fgm. 49 Bergk </hb-fs-l-normal><hb-fs-l-italic>Poet. lyr. gr. III</hb-fs-l-italic>&4
	//  note the trailing `&4` which means more shifts are coming but that we reached a `█`

	// 0: whole match
	// 1: shift number
	// 2: the line of text
	// 3: the next number (if it exists)

	// `&(\d{0,2})([^$█]*)&(\d{0,2})`

	// the issue here is that you might be sent `&4g&GP&4g$`
	// that needs sub-processing

	groups := latshiftinsidelatshift.FindStringSubmatch(match)
	// [&❨fgm. 827 Nauck &3Trag. gr. fgm.&4  ❨fgm. 827 Nauck &3Trag. gr. fgm. 4]

	reparse := groups[1] + groups[2]

	var components []string
	subunits := strings.Split(reparse, "&")
	for _, unit := range subunits {
		if unit != "" {
			toshift := "&" + unit + "&"
			components = append(components, vanillalatinshifter(toshift))
		}
	}

	//fmt.Println("latinsubshifter in: ", len(match), "\t", match)
	//fmt.Println("latinsubshifter out:", strings.Join(components, "")+"&"+groups[3])
	return strings.Join(components, "") + "&" + groups[3]
}

// simplegktolatinspan - `greek &stuff_in_latin_font$ greek`
func simplegktolatinspan(ttc string) string {
	const (
		LFSN = `<hb-fs-l-normal>$1</hb-fs-l-normal>`
	)
	// dangerous if there are balance problems...
	return basicgktolatfshift.ReplaceAllString(ttc, LFSN)
}

// severallinesoflatin - `greek &3latin ... █⑨⓪ &3more latin... █⑨⓪ &3stuff_in_shifted_latin_font$ greek`
func severallinesoflatin(ttc string) string {
	// example:
	//  *)ARISTAGO/RAS: OU)X O(MOLOGOU=SI DE\ AU)TOI=S *AI)GU/PTIOI. █⑨⓪ <hmu_standalone_tabbedtext />&Plinius H. N. XXXVI, 12, 17﹕ &3Qui de iis scri-
	// &3pserint ❨&sc. de pyramidibus❩, &3sunt Herodotus, Eu-
	// &3hemerus, Duris Samius, Aristagoras, Dionysius,
	// ...
	// &3ac caepas mille sexcenta talenta erogata.$ <hmu_standalone_endofpage /> █⑨⑧ █⑧ⓑ <hmu_standalone_tabbedtext />&Diog. L. I, 72$﹕ *BRAXULO/GOS TE H)=N ❨&sc$. O( *XI/LWN❩:
	return wholelinelatftshift.ReplaceAllStringFunc(ttc, severallinesoflatinshifter)
}

func severallinesoflatinshifter(match string) string {
	// &(\d{0,2})([^$█]*)█
	// 0: whole match
	// 1: shift number
	// 2: the line of text

	const (
		UHP    = `<hgb-build-error>＆%s<hgb-build-error>`
		KLUDGE = `%s`
	)

	groups := wholelinelatftshift.FindStringSubmatch(match)

	// the match might be for "" if just "&" sent you here
	if groups[1] == "" {
		groups[1] = "0"
	}

	tags, exists := latinfontmap[groups[1]]
	if !exists {
		if WARNINGS {
			// choking on lots of things that look like: "&87</hb-fs-l-romannumerals>❩, █"
			// and here it always looks like you just want the number "87"
			// so the underlying request was <rnon><normon><off>?
			// but KLUDGE is likely a problem if we are also failing on other patterns...
			// fmt.Println("severallinesoflatinshifter() error\t", groups[1], "\t", groups[0], "\t", groups[2])
		}
		// fmt.Println(fmt.Sprintf(KLUDGE, groups[1]) + groups[2] + "█")
		return fmt.Sprintf(KLUDGE, groups[1]) + groups[2] + "█"
	}

	return tags[0] + groups[2] + tags[1] + "█"
}

func SimpleLatinSpanLATE(ttc string) string {
	// if you call "simplelatinshift" in the normal course of things, you will break a bunch of stuff
	// so we are only invoking it as a Late item in a latin-only build: "lat.LatinCleanup()"
	return vanillalatinshift.ReplaceAllStringFunc(ttc, vanillalatinshifter)
}

func vanillalatinshifter(match string) string {
	// `&(\d{1,2})(.*?)&`
	// 0: whole match
	// 1: shift number
	// 2: the text span

	groups := vanillalatinshift.FindStringSubmatch(match)
	if groups[1] == "" {
		groups[1] = "0"
	}
	tags, exists := latinfontmap[groups[1]]
	if !exists {
		global.MSG(fmt.Sprintf("vanillalatinshifter() could not find()\t%s", groups[1]))
		return match
	}
	//fmt.Println("vanillalatinshifter in: ", match)
	//fmt.Println("vanillalatinshifter out:", tags[0]+groups[2]+tags[1])
	return tags[0] + groups[2] + tags[1]
}
