//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lat

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/grk"
)

var (
	grabpsuedoline     = regexp.MustCompile(`(.*?)(\s?█)`)
	greekspan1         = regexp.MustCompile(`(\$\d{0,2})([^$█]+?)(<hb-fs-l-normal>)`)
	greekspan2         = regexp.MustCompile(`(</hb-fs-l-normal>)([^$<█]+?)(<hb-fs-l-normal>)`)
	greekspan3         = regexp.MustCompile(`(\$\d{0,2})([^$█]+?)$`)
	greekspan4         = regexp.MustCompile(`(</hb-fs-l-normal>)([^$<█]+?)$`)
	betacodecandidate1 = regexp.MustCompile(`(<hb-fs-g-[^>]+?>)([^█]+?)(</hb-fs-g-[^>]+?>)`)
	betacodecandidate2 = regexp.MustCompile(`(<hb-fs-g-[^>]+?>)([^█]+)(<hb-fs-l-[^>]+?>)`)
	betacodecandidate3 = regexp.MustCompile(`(</hb-fs-g-[^>]+?><hb-fs-l-[^>]+?>)([^█]+)(</hb-fs-l-[^>]+?>)([^█]+)(█[^a-z]+)(<hb-fs-g-[^>]+?>)`)
)

func GreekFontshiftsInLatinAuthor(ttc string) string {
	ttc = grabpsuedoline.ReplaceAllStringFunc(ttc, rewritegreekfontshift)

	ttc = betacodecandidate1.ReplaceAllStringFunc(ttc, applybetacodeconversion1)

	// Res Gestae fails betacodecandidate1 at lines where `<hb-sp-alternative_reading>` is present; `[^<]` in "betacodecandidate1" is the problem
	// $BE/RWNI, § TH=S [TE S1]⟨UNKLH/TOU⟩ ⟨KAI\⟩ ⟨TOU=⟩ ⟨DH/MOU⟩ ⟨T⟩W=N <hb-sp-alternative_reading>TOU= <hb-fs-l-normal>Apoll.</hb-sp-alternative_reading> </hb-fs-l-normal>

	ttc = betacodecandidate2.ReplaceAllStringFunc(ttc, applybetacodeconversion2)

	// Res Gestae fails betacodecandidate2 in the following
	// OU)DEI\S █⑧⓪ $E)/NPROS1⟨QEN⟩ ⟨I(STO/RHS1＇⟩ <hb-fs-l-normal><hb-sp-alternative_reading>Apoll. </hb-fs-l-normal>I(STO/RHSEN</hb-sp-alternative_reading> ⟨E)PI\⟩ ⟨*(RW/MHS⟩ G⟨EGONE/NAI⟩, ⟨*PO⟩-█⑧⓪ $⟨PLI/W⟩I
	// █⑧⓪ <hb-fs-g-normal>ἔνπροϲ⟨θεν⟩ ⟨ἱϲτόρηϲ＇⟩ </hb-fs-g-normal><hb-fs-l-normal><hb-sp-alternative_reading>Apoll. </hb-fs-l-normal>I(STÓRHSEN</hb-sp-alternative_reading> ⟨E)PÌ⟩ ⟨*(RW/MHS⟩ G⟨EGONÉNAI⟩, ⟨*PO⟩-█⑧⓪
	// the 'aternative reading' switches greek off and it never comes back on again until the next line

	ttc = betacodecandidate3.ReplaceAllStringFunc(ttc, applybetacodeconversion3)

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

func applybetacodeconversion1(gfs string) string {
	groups := betacodecandidate1.FindStringSubmatch(gfs)
	if len(groups) == 4 {
		unicodegrk := grk.ConvertGreekLowers(grk.ConvertGreekCapitals(groups[2]))
		subst := groups[1] + unicodegrk + groups[3]
		gfs = betacodecandidate1.ReplaceAllString(gfs, subst)
	}
	return gfs
}

func applybetacodeconversion2(gfs string) string {
	groups := betacodecandidate2.FindStringSubmatch(gfs)
	if len(groups) == 4 {
		unicodegrk := grk.ConvertGreekLowers(grk.ConvertGreekCapitals(groups[2]))
		subst := groups[1] + unicodegrk + groups[3]
		gfs = betacodecandidate2.ReplaceAllString(gfs, subst)
	}
	return gfs
}

func applybetacodeconversion3(gfs string) string {
	// `(</hb-fs-g-[^>]+?><hb-fs-l-[^>]+?>)([^█]+)(</hb-fs-l-[^>]+?>)([^█]+)(█[^a-z]+)(<hb-fs-g-[^>]+?>)`
	// 0: whole
	// 1: greek-off+latin-on
	// 2: latincontents
	// 3: latin-off
	// 4: betacode candidate
	// 5: highhexrun
	// 6: greek-on at start of next line

	groups := betacodecandidate3.FindStringSubmatch(gfs)
	if len(groups) == 7 {
		unicodegrk := grk.ConvertGreekLowers(grk.ConvertGreekCapitals(groups[4]))
		subst := groups[1] + groups[2] + groups[3] + unicodegrk + groups[5] + groups[6]
		gfs = betacodecandidate3.ReplaceAllString(gfs, subst)
	}
	return gfs
}
