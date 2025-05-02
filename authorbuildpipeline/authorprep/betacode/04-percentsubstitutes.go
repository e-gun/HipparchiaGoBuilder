//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var (
	percentsouter  = regexp.MustCompile(`%\d{1,4}`)
	percentsubsmap = map[int]string{
		// u'\u",
		// many of these early items can look like betacode accents, etc.
		// using small variants now; then lastpass can restore them
		1:   "﹖", // '?' (003f), but using small variant instead: fe56
		2:   "﹡", // '*' (u002a), but using small variant instead: fe61
		3:   "／", // '/' (u002f), but using fullwidth variant instead: ff0f
		4:   "﹗", // '!' (u0021), but using small variant instead: ff57
		5:   "│", // '|' (u007c), but using box drawings light vertical instead: 2502
		6:   "﹦", // '=' (u003d); but using small variant instead: fe66
		7:   "﹢", // '+' (u002b); but using small variant instead: fe62
		8:   "﹪", // '%' (u0025), but using small variant instead: fe6a
		9:   "﹠", // '&' (0026) is also a control character; using small version instead (fe60); can swap it out in the end
		10:  "﹕", // ':' (003a); samll variant instead: fe55
		11:  "\u2022",
		12:  "﹡", // '*' (u002a) might lead to future problems: small version instead (fe61); can swap it out in the end
		14:  "\u00a7",
		15:  "\u02c8",
		16:  "\u00a6",
		17:  "\u2016",
		18:  "\u0027", // look out for future problems: '
		19:  "\u2013",
		20:  "\u0301",
		21:  "\u0300",
		22:  "\u0302",
		23:  "\u0308",
		24:  "\u0342",
		25:  "\u0327",
		26:  "\u0304",
		27:  "\u0306",
		28:  "\u0308",
		29:  "\u0323\u0323",
		30:  "\u02bc",
		31:  "\u02bd",
		32:  "\u00b4", // look out for future problems: ´
		33:  "\u0060", // look out for future problems: `
		34:  "\u1fc0",
		35:  "\u1fce",
		36:  "\u1fde",
		37:  "\u1fdd",
		38:  "\u1fdf",
		39:  "\u00a8",
		40:  "\u23d1",
		41:  "\u2013",
		42:  "\u23d5",
		43:  "\u00d7",
		44:  "\u23d2",
		45:  "\u23d3",
		46:  "\u23d4",
		47:  "𐄑",
		48:  "u23d1\u23d1",
		49:  "\u23d1\u23d1\u23d1",
		50:  "<hb-pap-fract>½</hb-pap-fract>",
		51:  "<hb-pap-fract>¼</hb-pap-fract>",
		52:  "<hb-pap-fract>⅛</hb-pap-fract>",
		53:  "<hb-pap-fract>⅟<span class=\"denominator\">16</span></hb-pap-fract>",
		54:  "<hb-pap-fract>⅟<span class=\"denominator\">32</span></hb-pap-fract>",
		55:  "<hb-pap-fract>⅟<span class=\"denominator\">64</span></hb-pap-fract>",
		56:  "<hb-pap-fract>⅟<span class=\"denominator\">128</span></hb-pap-fract>",
		57:  "<hb-undoc-pct betacodeval=\"57\">⊚</hb-undoc-pct>",
		59:  "<hb-pap-fract>¾</hb-pap-fract>",
		60:  "<hb-pap-fract>⅓</hb-pap-fract>",
		61:  "<hb-pap-fract>⅙</hb-pap-fract>",
		62:  "<hb-pap-fract>⅟<span class=\"denominator\">12</span></hb-pap-fract>",
		63:  "<hb-pap-fract>⅟<span class=\"denominator\">24</span></hb-pap-fract>",
		64:  "<hb-pap-fract>⅟<span class=\"denominator\">48</span></hb-pap-fract>",
		65:  "<hb-pap-fract>⅟<span class=\"denominator\">96</span></hb-pap-fract>",
		69:  "\u03b2\u0338",
		70:  "<hb-pap-fract>⅟<span class=\"denominator\">50</span></hb-pap-fract>",
		71:  "<hb-pap-fract>⅟<span class=\"denominator\">100</span></hb-pap-fract>",
		72:  "<hb-pap-fract>⅟<span class=\"denominator\">100</span></hb-pap-fract>",
		73:  "<hb-pap-fract>⅟<span class=\"denominator\">100</span></hb-pap-fract>",
		75:  "<hb-undoc-pct betacodeval=\"75\">⊚</hb-undoc-pct>",
		79:  "<hb-undoc-pct betacodeval=\"79\">⊚</hb-undoc-pct>",
		80:  "<hb-sp-italic> v. </hb-sp-italic>",
		81:  "<hb-sp-italic> vac. </hb-sp-italic>",
		91:  "\u0485",
		92:  "\u0486",
		93:  "\u1dc0",
		94:  "\u0307",
		95:  "\u1dc1",
		96:  "\u035c",
		97:  "\u0308",
		98:  "\u0022",
		99:  "\u2248",
		100: "\u003b",
		// 101: "\u0023",  // had best do pounds before percents since this is '//'
		101: "﹟", // small number sign instead (ufe5f)
		102: "’", // single quotation mark
		// 103: "\u005c",  // backslash: careful
		103: "﹨", // small reverse solidus instead: ufe68
		105: "\u007c\u007c\u007c",
		106: "\u224c",
		107: "\u007e", // '~'
		108: "\u00b1", // '±'
		109: "\u00b7", // middle dot: '·'
		110: "\u25cb",
		127: "\u032f",
		128: "\u0302",
		129: "\u2020",
		130: "\u0307",
		132: "΅",
		133: "\u1fcd",
		134: "\u1fcf",
		140: "𐄒",
		141: "\u23d6",
		144: "\u23d1\u0036",
		145: "\u2013\u0301",
		146: "\u00b7",
		147: "\u030a",
		148: "\u030c",
		149: "\u0328",
		150: "\u007c",
		157: "\u2e38",
		159: "\u00d7", // multiplication sign: '×'
		160: "\u002d", // hyphen-minus: '-'
		162: "\u0338",
	}
)

func ReplacePercentSigns(ttc string) string {
	// Format % markup - Additional Punctuation and Characters
	// %50-79 Reserved for Greek documentary papyri
	// %80-89 Reserved for Greek inscriptions
	ttc = percentsouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return percentsubstitutes(match)
	})

	return ttc
}

func percentsubstitutes(match string) string {
	const (
		UHP = `<hb-build-error>％%s<hb-build-error>`
	)

	// match looks like "%14": drop that initial character
	match = strings.TrimPrefix(match, "%")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := percentsubsmap[val]
	if !exists {
		//if WARNINGS {
		//	fmt.Println("percentsubstitutes()\t", match)
		//}
		substitute = secondtrypercentsubstitutes(val)
	}
	return substitute
}

func secondtrypercentsubstitutes(firsttry int) string {
	// garbage can come up out of the annotations:
	// <hb-metadata_publicationinfo value="CIL 3.6545,17%67309" />
	// there is no "%67309"; presumably this is supposed to be "%6`7309"

	bigInt := big.NewInt(0)
	bigInt.SetInt64(int64(firsttry))

	ten := big.NewInt(10)

	// Divide the number by 10 until it's a single-digit number
	for bigInt.Cmp(ten) >= 0 {
		bigInt.Div(bigInt, ten)
	}
	secondtry := int(bigInt.Int64())
	substitute, exists := percentsubsmap[secondtry]
	if !exists {
		substitute = fmt.Sprintf("%d", firsttry)
		global.MSG(fmt.Sprintf("secondtrypercentsubstitutes() failed\t%d", firsttry))
	} else {
		// as far as I can tell it is always "="...
		// secondtrypercentsubstitutes() x into y:  6142 6 ﹦
		// fmt.Println("secondtrypercentsubstitutes() x into y:\t", firsttry, secondtry, substitute)
	}
	return substitute
}
