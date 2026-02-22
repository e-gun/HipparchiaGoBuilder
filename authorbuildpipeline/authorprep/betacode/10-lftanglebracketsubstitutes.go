//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
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
	labouter    = regexp.MustCompile(`<\d{1,2}`)
	labrsubsmap = map[int]string{
		// 1: "",
		// Diogenes seems to have decided that this is the way to go; I wonder how often you will be sorry that you do not have \u0332 instead...
		1: "⟨",
		// the inactive version is what the betacode manual says to do, but in the inscriptions we just want brackets and not a combining underline
		2:  "\u2035",
		3:  "",
		4:  "",
		5:  "",
		6:  "<hb-sp-supsc>", // hmu_shift_font_to_superscript
		7:  "<hb-sp-subsc>", // hmu_shift_font_to_subscript
		8:  "",
		9:  "<hb-sp-textual_lemma>",      // hmu_textual_lemma
		10: "<hb-sp-stacked_text_lower>", // hmu_stacked_text_lower
		11: "<hb-sp-stacked_text_upper>", // hmu_stacked_text_upper
		12: "<hb-sp-nonstandarddirection>",
		13: "<hmu_standalone_singlelinespacing_in_doublespacedtext />",
		14: "<hb-sp-interlineartext>",
		15: "<hb-sp-interlinearmarginalia>", // hmu_interlinear_marginalia
		16: "\u2035",
		17: "", // Combining Double Underline
		19: "\u2035",
		20: "<hb-sp-expanded_text>",       // hmu_expanded_text
		21: "<hb-sp-latin_expanded_text>", // hmu_latin_expanded_text
		22: "<hb-unk-anglebracketspan22>",
		24: "<hb-unk-anglebracketspan24>",
		30: "<hb-sp-overline>", // Combining Overline and Dependent Vertical Bars
		31: "<hb-sp-strikethrough>",
		32: "<hb-sp-overline_and_underline>", // hmu_overline_and_underline
		34: "⁄",                              // fractions (which have balanced sets of markup...)
		48: "<hb-unk-anglebracketspan48>",
		50: "<hb-unk-anglebracketspan50>",
		51: "<hb-unk-anglebracketspan51>",
		52: "<hb-unk-anglebracketspan52>",
		53: "<hb-unk-anglebracketspan53>",
		60: "<hb-sp-preferred_epigraphical_text_used>",
		61: "<hb-sp-epigraphical_text_inserted_after_erasure>",
		62: "<hb-sp-lineover>",
		63: "<hb-sp-epigraphical_text_after_correction>",
		64: "<hb-sp-letterbox>",
		65: "<hmu_epigraphical_letters_enclosed_in_wreath>",
		66: "<hmu_epigraphical_project_escape_66>",
		67: "<hmu_epigraphical_project_escape_67>",
		68: "<hmu_epigraphical_project_escape_68>",
		69: "<hmu_epigraphical_project_escape_69>",
		70: "<hb-sp-diagram>",         // hmu_inset_diagram
		71: "<hb-sp-diagramsection>",  // hmu_inset_diagram
		72: "<hb-sp-diagramrelation>", // hmu_logical_relationship_in_diagram
		73: "<hb-sp-diagramlvl03>",
		74: "<hb-sp-diagramlvl04>",
		82: "<hb-unk-anglebracketspan82>",
		91: "<hb-unk-anglebracketspan91>",
		96: "<hb-unk-anglebracketspan96>",
	}
)

func ReplaceLeftAngledBrackets(ttc string) string {
	// Purge < markup - Text Formatting
	// <50-59 Reserved for Greek documentary papyri
	// <60-69 Reserved for Greek inscriptions

	ttc = labouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return labsubstitutes(match)
	})

	return ttc
}

func labsubstitutes(match string) string {
	const (
		UHP = `<hgb-build-error>＜%s<hgb-build-error>`
	)

	// match looks like "<14": drop that initial character
	match = strings.TrimPrefix(match, "<")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := labrsubsmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("labsubstitutes()\t%s", match))
	}
	return substitute
}
