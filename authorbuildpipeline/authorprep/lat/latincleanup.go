//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lat

import (
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/latetidyups"
)

func LatinCleanup(ttc string) string {

	ttc = betacode.DollarSignGreekFontMarkup(ttc)
	ttc = betacode.AmpersandLatinFontMarkup(ttc)
	ttc = GreekFontshiftsInLatinAuthor(ttc)
	ttc = latetidyups.PurgeHybrid(ttc)
	ttc = ConvertLatinDiacriticals(ttc)
	ttc = betacode.SimpleLatinSpanLATE(ttc)

	//authorprep.WriteTTCProgress(ttc)
	//os.Exit(1)
	return ttc
}
