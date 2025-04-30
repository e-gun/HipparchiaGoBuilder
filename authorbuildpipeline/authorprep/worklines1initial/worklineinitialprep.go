//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

// WorklinePrep - the last segment with the full text block; the key function is HexRunner() which decodes the binary blocks
func WorklinePrep(ttc string) string {

	ttc = HexRunner(ttc)
	ttc = BracketFixer(ttc)
	ttc = BalanceQuotes(ttc)
	ttc = DebugHostileSubstitutions(ttc)
	ttc = InsertNewLinesAtNewLevels(ttc)
	ttc = FixIrrationalTags(ttc)

	//
	// VARIOUS DEBUGGING TOOLS....
	//

	// ttc = authorprep.WorkLineViewAndFilter("WorklinePrep", "FixIrrationalTags", true, "", ttc)
	// ttc = authorprep.DipIntoResults("WorklinePrep", "HexRunner", 2000, true, ttc)
	//authorprep.WriteTTCProgress(ttc)
	//os.Exit(1)

	return ttc
}
