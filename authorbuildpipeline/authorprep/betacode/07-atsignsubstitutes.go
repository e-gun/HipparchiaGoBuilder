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

var (
	atsignouter  = regexp.MustCompile(`@\d{1,2}`)
	atsignsubmap = map[int]string{
		1:  "<hb-nb-endofpage />",
		2:  "<hb-nb-column_end />",
		3:  "<hb-nb-omitted_graphic_marker />",
		4:  "<hb-nb-table_starts />",
		5:  "<hb-nb-table_ends />",
		6:  "<br>",
		7:  "——————",
		8:  "———", // hb-nb-mid_line_citation_boundary
		9:  "<hb-sp-breakintext>[break in text for unknown length]</hb-sp-breakintext>",
		10: "<hb-nb-linetoolongforscreen />",
		11: "<hb-nb-tablecellindicator />",
		12: "<hb-nb-subcellindicator />",
		20: "<hb-nb-start_of_columnar_text />",
		21: "</hb-nb-end_of_columnar_text>",
		30: "<hb-nb-start_of_stanza />",
		50: "<hb-nb-writing_perpendicular_to_main_text />",
		51: "<hb-nb-writing_inverse_to_main_text />",
		70: "<hb-sp-quotedtext>",
		71: "</hb-sp-quotedtext>",
		73: "<hb-sp-poetictext>",
		74: "</hb-sp-poetictext>",
	}
)

func ReplaceAtSigns(ttc string) string {
	// Purge @ markup - Page Formatting
	// @50-59 Reserved for Greek documentary papyri
	// @60-69 Reserved for Greek inscriptions
	ttc = atsignouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return atsubstitutes(match)
	})

	return ttc
}

func atsubstitutes(match string) string {
	const (
		UHP = `<hb-build-error>＠%s<hb-build-error>`
	)

	// match looks like "@6": drop that initial character
	match = strings.TrimPrefix(match, "@")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := atsignsubmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("atsubstitutes()\t%s", match))
	}
	return substitute
}
