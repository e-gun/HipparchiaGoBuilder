//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dating

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"regexp"
	"strings"
)

var (
	hasonearabic   = regexp.MustCompile(`\d+`)
	hastwoarabic   = regexp.MustCompile(`\d+\D+\d`)
	hasthreearabic = regexp.MustCompile(`\d+\D+\d+\D+\d`)
	has4darabic    = regexp.MustCompile(`\d\d\d\d`)
	has3darabic    = regexp.MustCompile(`\d\d\d`)
	has2darabic    = regexp.MustCompile(`\d\d`)
	has1darabic    = regexp.MustCompile(`\d`)
	hasoneroman    = regexp.MustCompile(`[IXV]+`)
	hastworoman    = regexp.MustCompile(`[IXV]+[^IXV]+[IXV]+`)
	germandate     = regexp.MustCompile(`(nach|vor|Jhdt|Jh)`)
	hasbce         = regexp.MustCompile(`(vor|BC|B\.C\.|bc|v.Chr|a$| a／| a-)`)
	hasce          = regexp.MustCompile(`(nach|AD|A\.D\.|ad|n.Chr|ac$|p$)`)
	has1sthalf     = regexp.MustCompile(`(1st half|1. H.|1. Hälfte)`)
	has2ndhalf     = regexp.MustCompile(`(2nd half|2. H.|2. Hälfte)`)
	hasbeg         = regexp.MustCompile(`(^in |init |init s |／in |- in |early|beg|Anf\.|Early)`)
	hasmid         = regexp.MustCompile(`mid`)
	hasend         = regexp.MustCompile(`(late|end|Late|fin )`)
	hasthird       = regexp.MustCompile(`Drittel`)
	hasquarter     = regexp.MustCompile(`Viertel`)
	hasth          = regexp.MustCompile(`(\dth|2nd|1st|3rd)`)
	hasonespan     = regexp.MustCompile(`[\-／]`)
	hastwospans    = regexp.MustCompile(`[\-／][^\-／][\-／]`)
	hasbefore      = regexp.MustCompile(`(before|bef\.)`)
	hasafter       = regexp.MustCompile(`(^p |after|aft\.)`)
	hasslashdate1  = regexp.MustCompile(`(\d+)／\d(\D)`)  // a year like 234/5 *not* a range/span
	hasor          = regexp.MustCompile(`(\d+)( or .*)`) // dangerous because we will end up discarding the rest
	hasorlike      = regexp.MustCompile(`(.*)( or \d.*\s)`)
	hasbracket     = regexp.MustCompile(`\[.*?]`)
	hasaet         = regexp.MustCompile(`aet`)
	nodigits       = regexp.MustCompile(`^[^0-9]*$`)
)

func TakeFingerprint(text string) structs.FingerPrint {
	text = strings.ReplaceAll(text, "<1", "")
	text = strings.ReplaceAll(text, ">1", "")
	fp := structs.FingerPrint{
		OrigDateString: text,
		Calculated:     9999,
	}
	if hasonearabic.MatchString(text) {
		fp.HasOneArabic = true
	}
	if hastwoarabic.MatchString(text) {
		fp.HasTwoArabic = true
	}
	if hasthreearabic.MatchString(text) {
		fp.HasThreeArabic = true
	}
	if has4darabic.MatchString(text) {
		fp.Has4Arabic = true
	}
	if has3darabic.MatchString(text) {
		fp.Has3dArabic = true
	}
	if has2darabic.MatchString(text) {
		fp.Has2dArabic = true
	}
	if has1darabic.MatchString(text) {
		fp.Has1dArabic = true
	}
	if hasoneroman.MatchString(text) {
		fp.HasOneRoman = true
	}
	if hastworoman.MatchString(text) {
		fp.HasTwoRoman = true
	}
	if germandate.MatchString(text) {
		fp.HasGerman = true
	}
	if hasbce.MatchString(text) {
		fp.HasBCE = true
	}
	if hasce.MatchString(text) {
		fp.HasCE = true
	}
	if hasonespan.MatchString(text) {
		fp.HasOneSpan = true
	}
	if hastwospans.MatchString(text) {
		fp.HasTwoSpans = true
	}
	if hasth.MatchString(text) {
		fp.HasTH = true
	}
	if hasend.MatchString(text) || has2ndhalf.MatchString(text) {
		fp.HasEnd = true
	}
	if hasbeg.MatchString(text) || has1sthalf.MatchString(text) {
		fp.HasBeginning = true
	}
	if hasbefore.MatchString(text) {
		fp.HasAnte = true
	}
	if hasafter.MatchString(text) {
		fp.HasPost = true
	}
	if hasor.MatchString(text) || hasorlike.MatchString(text) {
		fp.HasOR = true
	}
	if hasslashdate1.MatchString(text) {
		fp.HasSlashDate = true
	}
	if hasbracket.MatchString(text) {
		fp.HasBracket = true
	}
	if strings.Count(text, "／") > 1 {
		fp.HasMultiSlashDate = true
	}
	if strings.Count(text, "／") > 0 && strings.Count(text, "-") > 1 {
		// 245-244／220-219 BC
		fp.HasMixedSpans = true
	}
	if nodigits.MatchString(text) {
		fp.ContainsNoDigits = true
	}

	fp.Rationalize()
	return fp
}
