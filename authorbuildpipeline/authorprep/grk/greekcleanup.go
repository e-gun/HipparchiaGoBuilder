//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

func GreekCleanup(ttc string) string {
	// the order of execution can matter a lot...
	// DollarSignGreekFontMarkup() has to happen before RestoreLatinWithinGreek(), for example

	ttc = GreekColons(ttc)
	ttc = InstertCoptic(ttc)
	ttc = ConvertGreekCapitals(ttc)
	ttc = ConvertGreekLowers(ttc)
	ttc = GreekExclamation(ttc)
	ttc = GreekCombiningDot(ttc)
	ttc = RestoreLatinWithinGreek(ttc)

	// ttc = authorprep.ViewAndFilter("GreekCleanup", "RestoreLatinWithinGreek", true, "hb-set_l", ttc)
	// ttc = authorprep.DipIntoResults("GreekCleanup", "RestoreLatinWithinGreek", 0, true, ttc)
	// authorprep.WriteTTCProgress(ttc)
	// os.Exit(1)
	return ttc
}
