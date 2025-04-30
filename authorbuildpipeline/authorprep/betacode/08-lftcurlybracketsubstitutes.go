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
)

var (
	lcbouter    = regexp.MustCompile(`\{\d{1,2}`)
	lcbrsubsmap = map[int]string{
		1: "<hb-title>",
		2: "<hb-sp-marginaltext>", // note that we are going to pull this into Annotations and out of the main line
		3: "<hb-sp-scholium>",     // TLG4090 has this; ex: <hb-sp-scholium><hb-fs-l-normal>Mt 6, 16 — 17</hb-fs-l-normal></hb-sp-reference_in_scholium> (but not a candidate for moving to Annotations)
		4: "<hb-sp-unconventional_form_written_by_scribe>",
		5: "<hb-sp-form_altered_by_scribe>",
		6: "<hb-sp-discarded_form>",
		7: "<hb-sp-reading_discarded_in_another_source>",
		8: "<hb-sp-numerical_equivalent>",
		9: "<hb-sp-alternative_reading>",
		// 10: u"\u0332",
		10: "⟨", // the inactive version is what the betacode manual says to do, but in the inscriptions we just want brackets and not a combining underline
		26: "<hb-sp-rectified_form>",
		27: "\u0359",
		28: "<hb-sp-date_or_numeric_equivalent_of_date>",
		29: "<hb-sp-emendation_by_editor_of_text_not_obviously_incorrect>",
		30: "<hb-inscr-non_text_characters_30 />",
		31: "<hb-inscr-non_text_characters_31 />",
		32: "<hb-inscr-non_text_characters_32 />",
		33: "<hb-inscr-non_text_characters_33 />",
		34: "<hb-inscr-non_text_characters_34 />",
		35: "<hb-inscr-non_text_characters_35 />",
		36: "<hb-inscr-non_text_characters_36 />",
		37: "<hb-inscr-non_text_characters_37 />",
		38: "<hb-inscr-non_text_characters_38 />",
		39: "<hb-inscr-non_text_characters_39 />",
		40: "<hb-speaker>",
		41: "<hb-sp-stagedirection>", // not in TLG...
		43: "<hb-serviusformatting>", // not the same as most "spans" since this is often multi-line; "span" will confuse the paragraph catcher in HipparchiaServer
	}
)

func ReplaceLeftCurlyBrackets(ttc string) string {
	// Purge { markup -  Textual Mark-Up
	//{4-24 Reserved for Greek documentary papyri
	//{25-39 Reserved for Greek inscriptions
	//{40-69 Reserved for PHI (Latin)

	ttc = lcbouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return lcbsubstitutes(match)
	})

	return ttc
}

func lcbsubstitutes(match string) string {
	const (
		UHP = `<hgb-build-error>{%s<hgb-build-error>`
	)

	// match looks like "{14": drop that initial character
	match = strings.TrimPrefix(match, "{")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := lcbrsubsmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		if WARNINGS {
			fmt.Println("lcbsubstitutes()\t", match)
		}

	}
	return substitute
}
