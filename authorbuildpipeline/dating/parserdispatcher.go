//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dating

import (
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func pickandrunparser(fp structs.FingerPrint) structs.FingerPrint {
	// order of the tests matters
	if fp.HasTwoArabic && fp.HasOneSpan && (fp.Has2dArabic || fp.Has3dArabic || fp.Has4Arabic) && !(fp.HasCE && fp.HasBCE) && !fp.HasTH {
		// twoarabicsimple -  `1-50 ac` --> 25
		// fmt.Println("twoarabicsimple")
		return twoarabicsimple(fp)
	}

	if fp.HasTwoArabic && fp.HasOneSpan && (fp.Has2dArabic || fp.Has3dArabic || fp.Has4Arabic) && fp.HasCE && fp.HasBCE {
		// twoarabiccomplex - `30 BC-AD 68` --> 49
		// fmt.Println("twoarabiccomplex")
		return twoarabiccomplex(fp)
	}

	if fp.HasOneArabic && (fp.Has2dArabic || fp.Has3dArabic || fp.Has4Arabic) {
		// onearabicsimple: `101 bc` --> -101
		return onearabicsimple(fp)
	}
	if fp.HasTwoArabic && fp.HasCE && fp.HasBCE && fp.Has1dArabic && !fp.HasOR {
		// arabiccenturybceandcedate - `1 BC／AD 1` --> 0
		return arabiccenturybceandcedate(fp)
	}
	if (fp.HasTH || fp.HasGerman) && fp.HasTwoArabic && (fp.Has1dArabic || fp.Has2dArabic) && !fp.HasOR {
		// twoarabicthcenturies -  `7th-9th ac` --> 750
		// fmt.Println("twoarabicthcenturies")
		return twoarabicthcenturies(fp)
	}
	if !fp.HasOR && (fp.HasTH || fp.HasGerman) && fp.HasOneArabic && (fp.Has1dArabic || fp.Has2dArabic) {
		// onearabiccentury : `early 5th bc` --> -475
		// fmt.Println("onearabiccentury")
		return onearabiccentury(fp)
	}
	if fp.HasTwoRoman && fp.HasOneSpan && fp.HasCE && fp.HasBCE {
		// tworomannumcenturiescomplex -
		return tworomannumcenturiescomplex(fp)
	}
	if fp.HasTwoRoman && fp.HasOneSpan {
		// tworomannumcenturies - `c II／IIIp` --> 200
		return tworomannumcenturies(fp)
	}
	if !fp.HasOR && fp.Has1dArabic && !fp.HasOneSpan {
		// more dangerous, but the way to do 'c 5p
		return onearabiccentury(fp)
	}
	if fp.HasOneRoman {
		// oneromancentury - `VIII bc?` -->   -750
		return oneromancentury(fp)
	}
	if fp.HasMultiSlashDate {
		// multislashdate - 109／108／106／105 BC --> -109
		return multislashdate(fp)
	}
	if len(strings.Split(fp.OrigDateString, ",")) > 1 {
		// multicommadate - 114, 116, ﹠ 156 ac --> 114
		return multicommadate(fp)
	}
	if len(strings.Split(fp.OrigDateString, "﹠")) > 1 {
		// andsigndate - 1299 ﹠ 1344 ac --> 1299
		return andsigndate(fp)
	}
	// this seems like it should be a late test...
	if fp.HasTwoArabic && fp.HasOneSpan && fp.Has1dArabic && !fp.HasOR {
		// `6-5 BC` --> -500
		return twoarabicthcenturies(fp)
	}

	// now we are in the zone where recursive calls might be made; look out for infinite loops
	if fp.HasSlashDate {
		// slashdated - `c.63／2-51／0 bc` --> -57
		return slashdated(fp)
	}
	if fp.HasOR {
		// eitherordate - `618 or 633 ac` --> 618
		return eitherordate(fp)
	}
	if fp.HasBracket {
		// bracketdate - `med Ia [K.87(14)]` --> 50
		return bracketdate(fp)
	}
	if fp.HasMixedSpans && fp.HasMultiDigitArabic {
		// mixedspans - `245-244／220-219 BC` --> 245
		return mixedspans(fp)
	}

	// desperate people should try a lookup: we might have "byzantinisch", vel sim
	if fp.ContainsNoDigits {
		return nodigitdispatcher(fp)
	}

	fp.LacksParser = true
	// fmt.Printf("no parser for '%s'\n", fp.OrigDateString)
	return fp
}
