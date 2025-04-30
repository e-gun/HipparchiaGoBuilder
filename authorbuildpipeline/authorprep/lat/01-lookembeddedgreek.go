//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lat

import (
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/grk"
	"regexp"
	"strconv"
	"strings"
)

var (
	grabpsuedoline    = regexp.MustCompile(`(.*?)(\s?█)`)
	greekspan1        = regexp.MustCompile(`(\$\d{0,2})([^$]*?)(<hb-fs-l-normal>)`)
	greekspan2        = regexp.MustCompile(`(</hb-fs-l-normal>)([^$<]*?)(<hb-fs-l-normal>)`)
	greekspan3        = regexp.MustCompile(`(\$\d{0,2})([^$]*?)$`)
	greekspan4        = regexp.MustCompile(`(</hb-fs-l-normal>)([^$<]*?)$`)
	betacodecandidate = regexp.MustCompile(`(<hb-fs-g-[^>]*?>)([^<]*?)(</hb-fs-g-[^>]*?>)`)
)

func GreekFontshiftsInLatinAuthor(ttc string) string {
	ttc = grabpsuedoline.ReplaceAllStringFunc(ttc, rewritegreekfontshift)
	ttc = betacodecandidate.ReplaceAllStringFunc(ttc, applybetacodeconversion)
	return ttc
}

func rewritegreekfontshift(line string) string {
	if !strings.Contains(line, "$") {
		// there is no greek fontshift in this line
		return line
	}

	groups := grabpsuedoline.FindStringSubmatch(line)
	linebody := groups[1]
	lineend := groups[2]

	// after DollarSignGreekFontMarkup and AmpersandLatinFontMarkup run there are two basic patterns left in a Latin
	// author: (1) the part line; (2) the whole line; (3) the tail of a line

	// (1) part lines:
	// (1a) L. Aelium magistrum suum in $E)TUMOLOGI/A|<hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal>
	// (1b) $KU/BON<hb-fs-l-normal>, quid </hb-fs-l-normal>GRAMMH/N<hb-fs-l-normal>; quibusque ista omnia Latinis uoca-</hb-fs-l-normal>
	// note the lack of a '$' before GRAMMH/N

	// (2) whole lines:
	// <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext />$EU)FHMEI=N XRH\ KA)CI/STASQAI TOI=S H(METE/ROISI XOROI=SIN,
	// <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext />$O(/STIS A)/PEIROS TOIW=NDE LO/GWN H)\ GNW/MH| MH\ KAQAREU/EI

	// (3) tail of a line (esp after punctuation sparks a very brief "latin normal"
	// </hb-fs-l-normal>LE/GE MOI, TOU\S QEOU/S SOI, A(\ PRW/|HN E)/LEGES,

	// you need to do the part lines before the whole lines

	// (1a) L. Aelium magistrum suum in $E)TUMOLOGI/A|<hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal>
	dollargroups := greekspan1.FindStringSubmatch(linebody)
	if len(dollargroups) == 4 {
		dsv, e := strconv.Atoi(strings.ReplaceAll(dollargroups[0], "$", ""))
		if e != nil {
			dsv = 0
		}
		brk := getdollarmapval(dsv)
		subst := brk[0] + dollargroups[2] + brk[1] + dollargroups[3]
		linebody = greekspan1.ReplaceAllString(linebody, subst)
	}

	// (1b) $KU/BON<hb-fs-l-normal>, quid </hb-fs-l-normal>GRAMMH/N<hb-fs-l-normal>; quibusque ista omnia Latinis uoca-</hb-fs-l-normal>
	dollargroups = greekspan2.FindStringSubmatch(linebody)
	if len(dollargroups) == 4 {
		brk := getdollarmapval(0)
		subst := dollargroups[1] + brk[0] + dollargroups[2] + brk[1] + dollargroups[3]
		linebody = greekspan2.ReplaceAllString(linebody, subst)
	}

	// (2) <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext />$EU)FHMEI=N XRH\ KA)CI/STASQAI TOI=S H(METE/ROISI XOROI=SIN,
	dollargroups = greekspan3.FindStringSubmatch(linebody)
	if len(dollargroups) == 3 {
		dsv, e := strconv.Atoi(strings.ReplaceAll(dollargroups[0], "$", ""))
		if e != nil {
			dsv = 0
		}
		brk := getdollarmapval(dsv)
		subst := brk[0] + dollargroups[2] + brk[1]
		linebody = greekspan3.ReplaceAllString(linebody, subst)
	}

	// (3) </hb-fs-l-normal>LE/GE MOI, TOU\S QEOU/S SOI, A(\ PRW/|HN E)/LEGES,
	dollargroups = greekspan4.FindStringSubmatch(linebody)
	if len(dollargroups) == 3 {
		brk := getdollarmapval(0)
		subst := dollargroups[1] + brk[0] + dollargroups[2] + brk[1]
		linebody = greekspan4.ReplaceAllString(linebody, subst)
	}
	return linebody + lineend
}

func getdollarmapval(val int) [2]string {
	// 0:  {`<hb-fs-g-normal>`, `</hb-fs-g-normal>`},
	sub, exists := betacode.DollarSubMap[val]
	if !exists {
		sub = betacode.DollarSubMap[0]
	}
	return sub
}

func applybetacodeconversion(gfs string) string {
	groups := betacodecandidate.FindStringSubmatch(gfs)
	if len(groups) == 4 {
		unicodegrk := grk.ConvertGreekLowers(grk.ConvertGreekCapitals(groups[2]))
		subst := groups[1] + unicodegrk + groups[3]
		gfs = betacodecandidate.ReplaceAllString(gfs, subst)
	}
	return gfs
}
