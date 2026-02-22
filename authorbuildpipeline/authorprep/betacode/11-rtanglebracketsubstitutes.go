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
	rabouter    = regexp.MustCompile(`>\d{1,2}`)
	rabrsubsmap = map[int]string{
		//1: "\u0332",
		1:  "⟩", // see note in ltanglebracketsubstitutes()
		2:  "\u2032",
		3:  "\u0361",         // Combining Inverted Breve
		4:  "\u035c",         // Combining Breve Below
		5:  "\u035d",         // Combining Breve
		6:  "</hb-sp-supsc>", // hmu_shift_font_to_superscript
		7:  "</hb-sp-subsc>", // hmu_shift_font_to_subscript
		8:  "\u0333",
		9:  "</hb-sp-textual_lemma>",        // hmu_textual_lemma
		10: "</hb-sp-stacked_text_lower>",   // hmu_stacked_text_lower
		11: "</hb-sp-stacked_text_upper>",   // hmu_stacked_text_upper
		12: "</hb-sp-nonstandarddirection>", // nonstandarddirection
		13: "<hb-singlelinespacing_in_doublespacedtext />",
		14: "</hb-sp-interlineartext>",       // interlineartext
		15: "</hb-sp-interlinearmarginalia>", // hmu_interlinear_marginalia
		16: "\u2032",
		17: "u\0333",
		19: "\u2032",
		20: "</hb-sp-expanded_text>",
		21: "</hb-sp-latin_expanded_text>",
		22: "</hb-unk-anglebracketspan22>",
		24: "</hb-unk-anglebracketspan24>",
		30: "</hb-sp-overline>",               // Combining Overline and Dependent Vertical Bars
		31: "</hb-sp-strikethrough>",          // strikethrough
		32: "</hb-sp-overline_and_underline>", // hmu_overline_and_underline
		34: "",                                // fractions
		48: "<hb-unk-anglebracketspan48>",
		50: "</hb-unk-anglebracketspan50>",
		51: "</hb-unk-anglebracketspan51>",
		52: "</hb-unk-anglebracketspan52>",
		53: "</hb-unk-anglebracketspan53>",
		60: "</hb-sp-preferred_epigraphical_text_used>",         // hmu_preferred_epigraphical_text_used
		61: "</hb-sp-epigraphical_text_inserted_after_erasure>", // hb-epig-text_inserted_after_erasure
		62: "</hb-sp-lineover>",                                 // epigraphical line over letters
		63: "</hb-sp-epigraphical_text_after_correction>",       // hb-epig-text_after_correction
		64: "</hb-sp-letterbox>",                                // letters in a box
		65: "</hb-epig-letters_enclosed_in_wreath>",
		66: "</hb-epig-project_escape_66>",
		67: "</hb-epig-project_escape_67>",
		68: "</hb-epig-project_escape_68>",
		69: "</hb-epig-project_escape_69>",
		70: "</hb-sp-diagram>",
		71: "</hb-sp-diagramsection>",
		72: "</hb-sp-diagramrelation>",
		73: "</hb-sp-diagramlvl03>",
		74: "</hb-sp-diagramlvl04>",
		82: "</hb-unk-anglebracketspan82>",
		91: "</hb-unk-anglebracketspan91>",
		96: "</hb-unk-anglebracketspan96>",
	}
)

func ReplaceRightAngledBrackets(ttc string) string {
	// Purge > markup
	ttc = rabouter.ReplaceAllStringFunc(ttc, func(match string) string {
		return rabsubstitutes(match)
	})

	return ttc
}

func rabsubstitutes(match string) string {
	const (
		UHP = `<hgb-build-error>＞%s<hgb-build-error>`
	)

	// match looks like ">14": drop that initial character
	match = strings.TrimPrefix(match, ">")
	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(UHP, match)
	}

	substitute, exists := rabrsubsmap[val]
	if !exists {
		substitute = fmt.Sprintf(UHP, match)
		global.MSG(fmt.Sprintf("rabsubstitutes()\t%s", match))
	}
	return substitute
}
