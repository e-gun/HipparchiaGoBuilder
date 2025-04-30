//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
	"strconv"
	"strings"
)

var (
	rcbouter    = regexp.MustCompile(`}\d{1,2}`)
	rcbrsubsmap = map[int]string{
		1: "</hb-title>",
		2: "</hb-sp-marginaltext>",
		3: "</hb-sp-reference_in_scholium>",
		4: "</hb-sp-unconventional_form_written_by_scribe>",
		5: "</hb-sp-form_altered_by_scribe>",
		6: "</hb-sp-discarded_form>",
		7: "</hb-sp-reading_discarded_in_another_source>",
		8: "</hb-sp-numerical_equivalent>",
		9: "</hb-sp-alternative_reading>",
		// 10: "\u0332",
		// 10: "\u0332",
		// cf. ltanglebracketsubstitutes() //1
		// Diogenes seems to have decided that this is the way to go; I wonder how often you will be sorry that you do not have \u0332 instead...
		10: "⟩", // the inactive version is what the betacode manual says to do, but in the inscriptions we just want brackets and not a combining underline
		26: "</hb-sp-rectified_form>",
		27: "\u0359",
		28: "</hb-sp-date_or_numeric_equivalent_of_date>",
		29: "</hb-sp-emendation_by_editor_of_text_not_obviously_incorrect>",
		40: "</hb-speaker>",
		41: "</hb-sp-stagedirection>",
		43: "</hb-serviusformatting>",
	}
)

func ReplaceRightCurlyBrackets(ttc string) string {
	// Purge { markup
	ttc = rcbouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return rcbsubstitutes(match)
	})

	return ttc
}

func rcbsubstitutes(match string) string {
	const (
		UHP = `<hgb-build-error>}%s<hgb-build-error>`
	)

	// match looks like "}14": drop that initial character
	match = strings.TrimPrefix(match, "}")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := rcbrsubsmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("rcbsubstitutes()\t%s", match))
	}
	return substitute
}
