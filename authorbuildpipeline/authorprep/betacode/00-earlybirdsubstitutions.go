//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

// lookaheads not allowed in go's regex
// https://github.com/dlclark/regexp2 has them

var (
	// '<' is worth taking care of before you deal with '<2', '<3', etc. Plus various difficult punctuation: ?, }, {, ...

	earlybirdreplacer1 = strings.NewReplacer(
		"'", "\uFF07",
		"_", "— ",
		"!", "·",
	)
	earlybirdreplacer2 = strings.NewReplacer(
		"'", "\uFF07",
		"_", "— ",
	)
	SmartSingleQuotes = false
	betacodetuples    = [][2]string{
		{`<([^0-9])`, `‹$1`}, // '<': this one is super-dangerous: triple-check
		{`>([^0-9])`, `›$1`}, // '>': this one is super-dangerous: triple-check

		// the papyri exposed an interesting problem with '?'
		// let's try to deal with this at EarlyBirdSubstitutions() because if you let '?' turn into '\u0323' it seems impossible to undo that
		//
		// many papyrus lines start like: '[ &c ? ]$' (cf. '[ &c ? $TO\ PRA=]GMA')
		// this will end up as: '[ <hmu_latin_normal>c ̣ ]</hmu_latin_normal>'
		// the space after '?' is not always there
		// 	'[ &c ?]$! KEKEI/NHKA DI/KH PERI\ U(/BREWS [4!!!!!!!!!![ &c ?]4 ]$'
		// also get a version of the pattern that does not have '[' early because we are not starting a line:
		//	'&{10m4}10 [ c ? ]$IASNI#80 *)EZIKEH\ M[ARTURW= &c ? ]$'
		// this one also fails to have '&c' because the '&' came earlier
		// here's hoping there is no other way to achieve this pattern...
		{`&c\s\?(.*?)\$`, `&c ﹖$1$`},      // the question mark needs to be preserved, so we substitute a small question mark
		{`\[\sc\s\?(.*?)\$`, `[ c ﹖$1$`},  // try to catch '&{10m4}10 [ c ? ]$I' without doing any damage
		{`&\?(.*?)\](.*?)\$`, `&﹖$1]$2$`}, // some stray lonely '?' cases remain
	}
	EarybirdTuples = GetEarlyBirdTuples()
	dotswap        = regexp.MustCompile(`!(\d+)`)
)

func GetEarlyBirdTuples() [][2]string {
	var thetuples [][2]string
	thetuples = append(thetuples, betacodetuples...)
	return thetuples
}

// EarlyBirdSubstitutions - attempt to get out in front of some of the trickiest bits that cause problems for others
func EarlyBirdSubstitutions(ttc string) string {
	ttc = dotswap.ReplaceAllStringFunc(ttc, func(match string) string {
		return dotter(match)
	})
	if global.WorkingOnCorpus != "LAT" {
		ttc = earlybirdreplacer1.Replace(ttc)
	} else {
		ttc = earlybirdreplacer2.Replace(ttc)
	}
	for _, tuple := range EarybirdTuples {
		re := regexp.MustCompile(tuple[0])
		ttc = re.ReplaceAllString(ttc, tuple[1])
	}
	return ttc
}

func dotter(match string) string {
	// [!10`22!10 A)PO] --> [10 dots + 22 + 10 dots + A)PO]
	// i.e, `[∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙22∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙ ἀπο]`
	match = strings.TrimPrefix(match, "!")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		fmt.Printf("dotter() error converting %s to int", match)
		return match
	}
	dots := make([]string, val)
	for i := 0; i < val; i++ {
		dots[i] = "·"
	}
	return strings.Join(dots, " ")
}
