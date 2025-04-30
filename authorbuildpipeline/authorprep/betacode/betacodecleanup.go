//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

func BetaCodeCleanup(ttc string) string {
	// the order of execution can matter a lot...
	// early has to happen early
	// singletons have to come last

	// NOTE any "halt" when probing below will stop the program at CleanDBAuthorNames() [i.e before you have a real author text]

	// fmt.Println(ttc)
	// authorprep.WriteTTCProgress(ttc)

	ttc = EarlyBirdSubstitutions(ttc)
	ttc = ReplaceQuotationMarks(ttc)
	ttc = InsertBlankQuarterSpaces(ttc)
	ttc = ReplacePoundSigns(ttc)
	ttc = ReplacePercentSigns(ttc)
	ttc = ReplaceLeftSquareBrackets(ttc)
	ttc = ReplaceRightSquareBrackets(ttc)
	ttc = ReplaceAtSigns(ttc)
	ttc = ReplaceLeftCurlyBrackets(ttc)
	ttc = ReplaceRightCurlyBrackets(ttc)
	ttc = ReplaceLeftAngledBrackets(ttc)
	ttc = ReplaceRightAngledBrackets(ttc)
	ttc = AmpersandLatinFontMarkup(ttc)
	ttc = DollarSignGreekFontMarkup(ttc)
	ttc = ReplaceSingletons(ttc)

	// ttc = authorprep.ViewAndFilter("BetaCodeCleanup", "ReplaceSingletons", false, "", ttc)
	// ttc = authorprep.DipIntoResults("BetaCodeCleanup", "ReplaceSingletons", 1, false, ttc)
	// authorprep.WriteTTCProgress(ttc)
	return ttc
}
