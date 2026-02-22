//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

import "regexp"

var (
	dbhsub1 = regexp.MustCompile(`[$]`)
)

func DebugHostileSubstitutions(ttc string) string {
	// 	all sorts of things will be hard to figure out if you run this suite
	//	but it does make many things 'look better' even if there are underlying problems.
	//
	//	see latinfontlinemarkupparser() for notes on what the problems are/look like
	//
	//	if the $ is part of an irrational 'on-without-off' Greek font toggle, then we don't care
	//	it is anything that does not fit that pattern that is the problem
	//
	//	the hard part is churning through lots of texts looking for ones that do not fit that pattern
	//
	//	at the moment few texts seem to have even the benign toggle issue; still looking for places
	//	where there is a genuine problem

	// note that '&' will return to the text via the hexrunner: it can be embedded in the annotations
	// and you will want it later in order to format that material when it hits HipparchiaServer:
	// in 'Gel. &3N.A.& 20.3.2' the '&3' turns on italics and stripping & leaves you with 3N.A. (which is hard to deal with)

	// $ is still a problem:
	// e.g., 0085:
	//   Der Antiatt. p. 115, 3 Bekk.: ‘ὑδρηλοὺϲ’ $πίθουϲ καὶ ‘οἰνηροὺϲ’
	//   @&Der Antiatt. p. 115, 3 Bekk.%10 $8’U(DRHLOU\S‘ $PI/QOUS KAI\ $8’OI)NHROU\S‘$

	if !HIDEKNOWNBLEMISHES {
		return ttc
	}

	ttc = dbhsub1.ReplaceAllString(ttc, "")
	return ttc
}
