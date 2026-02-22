//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dating

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	stripalmostallstrings = regexp.MustCompile(`[^\d\-／\s]`)
	stripallstrings       = regexp.MustCompile(`\D`)
	swapspanner           = regexp.MustCompile(`[\-／]`)
	stripforromans1       = regexp.MustCompile(`[^IVX\-／\s]`)
	stripforromans0       = regexp.MustCompile(`[^IVX\s]`)
)

// onearabicsimple: `101 bc` --> -101
func onearabicsimple(fp structs.FingerPrint) structs.FingerPrint {
	// fmt.Println("onearabicsimple", fp.OrigDateString)
	cleaned := stripallstrings.ReplaceAllString(fp.OrigDateString, "")
	cleaned = strings.TrimSpace(cleaned)
	d, e := strconv.Atoi(cleaned)
	if e != nil {
		global.MSG(fmt.Sprintf("onearabicsimple parser failed '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}
	fp.Calculated = d
	fp.ApplyBCE()
	fp.ApplySimpleFudges()
	// fmt.Println("onearabicsimple", fp.Calculated)
	return fp
}

// twoarabicsimple -  `1-50 ac` --> 25
func twoarabicsimple(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripalmostallstrings.ReplaceAllString(fp.OrigDateString, "")
	cleaned = swapspanner.ReplaceAllString(cleaned, " ")
	cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	halves := strings.Split(strings.TrimSpace(cleaned), " ")
	if len(halves) != 2 {
		global.MSG(fmt.Sprintf("twoarabicsimple parser failed ptA '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}
	d1, e1 := strconv.Atoi(halves[0])
	d2, e2 := strconv.Atoi(halves[1])
	if e1 != nil || e2 != nil {
		global.MSG(fmt.Sprintf("twoarabicsimple parser failed ptB '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}

	// `before 337／6 bc` does not involve the average of two numbers...
	if len(halves[0]) != 3 && len(halves[1]) != 1 {
		mid := (d1 + d2) / 2
		fp.Calculated = mid
	} else {
		fp.Calculated = d1
	}
	fp.ApplyBCE()
	fp.ApplySimpleFudges()
	return fp
}

// twoarabiccomplex - `30 BC-AD 68` --> 49
func twoarabiccomplex(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripalmostallstrings.ReplaceAllString(fp.OrigDateString, "")
	cleaned = swapspanner.ReplaceAllString(cleaned, " ")
	chunks := strings.Split(strings.TrimSpace(cleaned), " ")
	var halves []string
	for _, chunk := range chunks {
		if chunk != "" {
			halves = append(halves, chunk)
		}
	}
	if len(halves) != 2 {
		global.MSG(fmt.Sprintf("twoarabiccomplex parser failed '%s' & '%s'", fp.OrigDateString, cleaned))
		fp.ParserFailed = true
		return fp
	}
	d1, e1 := strconv.Atoi(halves[0])
	d2, e2 := strconv.Atoi(halves[1])
	if e1 != nil || e2 != nil {
		global.MSG(fmt.Sprintf("twoarabiccomplex parser failed ptB '%s' & '%s'", halves[0], halves[1]))
		fp.ParserFailed = true
		return fp
	}

	mid := (d1 + d2) / 2
	fp.Calculated = mid
	fp.ApplySimpleFudges()
	return fp
}

// tworomannumcenturiescomplex - `s I a／s I p`
func tworomannumcenturiescomplex(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripforromans1.ReplaceAllString(fp.OrigDateString, "")
	cleaned = strings.ReplaceAll(cleaned, "／", "-")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	halves := strings.Split(strings.TrimSpace(cleaned), "-")
	if len(halves) != 2 {
		global.MSG(fmt.Sprintf("tworomannumcenturiescomplex parser failed split '%s' & '%s'", fp.OrigDateString, cleaned))
		fp.ParserFailed = true
		return fp
	}
	d1, e1 := romannumerals[halves[0]]
	d2, e2 := romannumerals[halves[1]]
	if e1 != true || e2 != true {
		global.MSG(fmt.Sprintf("tworomannumcenturiescomplex parser failed ptB '%s' : %s\n", fp.OrigDateString, halves))
		fp.ParserFailed = true
		return fp
	}
	d1 = (-100 * d1) + 50
	d2 = (100 * d2) - 50

	fp.Calculated = (d1 + d2) / 2
	fp.ApplySimpleFudges()
	return fp
}

// arabiccenturybceandcedate - `1 BC／AD 1` --> 0
func arabiccenturybceandcedate(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripalmostallstrings.ReplaceAllString(fp.OrigDateString, "")
	cleaned = strings.ReplaceAll(cleaned, "／", "-")
	halves := strings.Split(cleaned, "-")
	if len(halves) != 2 {
		global.MSG(fmt.Sprintf("romanbcedate parser failed split '%s' & '%s'", fp.OrigDateString, cleaned))
		fp.ParserFailed = true
		return fp
	}
	halves[0] = strings.TrimSpace(halves[0])
	halves[1] = strings.TrimSpace(halves[1])
	d1, e1 := strconv.Atoi(halves[0])
	d2, e2 := strconv.Atoi(halves[1])
	if e1 != nil || e2 != nil {
		global.MSG(fmt.Sprintf("romanbcedate parser failed ptB '%s' & '%s'", halves[0], halves[1]))
		fp.ParserFailed = true
		return fp
	}
	d1 = (d1 * -100) + 50
	d2 = (d2 * 100) - 50
	mid := (d1 + d2) / 2
	fp.Calculated = mid
	fp.ApplySimpleFudges()
	return fp
}

// twoarabicthcenturies -  `7th-9th ac` --> 750
func twoarabicthcenturies(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripalmostallstrings.ReplaceAllString(fp.OrigDateString, "")
	cleaned = swapspanner.ReplaceAllString(cleaned, " ")
	cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	halves := strings.Split(strings.TrimSpace(cleaned), " ")
	// fmt.Println(halves)
	if len(halves) != 2 {
		global.MSG(fmt.Sprintf("twoarabicthcenturies parser failed ptA '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}
	d1, e1 := strconv.Atoi(halves[0])
	d2, e2 := strconv.Atoi(halves[1])
	if e1 != nil || e2 != nil {
		global.MSG(fmt.Sprintf("twoarabicthcenturies parser failed ptB '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}
	mid := (float32(d1) + float32(d2)) / 2
	mid = mid * 100
	fp.Calculated = int(mid)
	fp.ApplyBCE()
	if fp.Calculated < 0 {
		// 5th bce is not -500 but -450...
		fp.Calculated = fp.Calculated + 50
	} else {
		// 5th ce is not 500 but 450...
		fp.Calculated = fp.Calculated - 50
	}
	fp.ApplySimpleFudges()
	return fp
}

// onearabiccentury : `early 5th bc` --> -475
func onearabiccentury(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripalmostallstrings.ReplaceAllString(fp.OrigDateString, "")
	cleaned = strings.TrimSpace(cleaned)
	d, e := strconv.Atoi(cleaned)
	if e != nil {
		global.MSG(fmt.Sprintf("onearabiccentury parser failed '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}
	fp.Calculated = d * 100
	fp.ApplyBCE()
	if fp.Calculated < 0 {
		// 5th bce is not -500 but -450...
		fp.Calculated = fp.Calculated + 50
	} else {
		// 5th ce is not 500 but 450...
		fp.Calculated = fp.Calculated - 50
	}
	fp.ApplySimpleFudges()
	return fp
}

// tworomannumcenturies - `c II／IIIp` --> 200
func tworomannumcenturies(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripforromans1.ReplaceAllString(fp.OrigDateString, "")
	cleaned = swapspanner.ReplaceAllString(cleaned, " ")
	cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	halves := strings.Split(strings.TrimSpace(cleaned), " ")
	// fmt.Println(halves)
	if len(halves) != 2 {
		global.MSG(fmt.Sprintf("tworomannumcenturies parser failed ptA '%s' : %s", fp.OrigDateString, cleaned))
		fp.ParserFailed = true
		return fp
	}
	d1, e1 := romannumerals[halves[0]]
	d2, e2 := romannumerals[halves[1]]
	if e1 != true || e2 != true {
		global.MSG(fmt.Sprintf("tworomannumcenturies parser failed ptB '%s' : %s", fp.OrigDateString, cleaned))
		fp.ParserFailed = true
		return fp
	}
	mid := (float32(d1) + float32(d2)) / 2
	mid = mid * 100
	fp.Calculated = int(mid)
	fp.ApplyBCE()
	if fp.Calculated < 0 {
		// 5th bce is not -500 but -450...
		fp.Calculated = fp.Calculated + 50
	} else {
		// 5th ce is not 500 but 450...
		fp.Calculated = fp.Calculated - 50
	}
	fp.ApplySimpleFudges()
	return fp
}

// oneromancentury - `VIII bc?` -->   -750
func oneromancentury(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := stripforromans0.ReplaceAllString(fp.OrigDateString, "")
	cleaned = strings.TrimSpace(cleaned)
	d, e := romannumerals[cleaned]
	if e != true {
		global.MSG(fmt.Sprintf("oneromancentury parser failed '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}
	fp.Calculated = d * 100
	fp.ApplyBCE()
	if fp.Calculated < 0 {
		// 5th bce is not -500 but -450...
		fp.Calculated = fp.Calculated + 50
	} else {
		// 5th ce is not 500 but 450...
		fp.Calculated = fp.Calculated - 50
	}
	fp.ApplySimpleFudges()
	return fp
}

// slashdated - `c.63／2-51／0 bc` --> -57
func slashdated(fp structs.FingerPrint) structs.FingerPrint {
	// note that there is an infinite loop possibility: pickandrunparser() is how you got here
	// if that "／" does not disappear, you could return

	// `(\d+)／\d(\D)`
	cleaned := hasslashdate1.ReplaceAllString(fp.OrigDateString, "$1$2")
	// fmt.Println("slashdated", cleaned)

	newfp := TakeFingerprint(cleaned)
	newfp = pickandrunparser(newfp)
	newfp.OrigDateString = fp.OrigDateString
	return newfp
}

// eitherordate - `618 or 633 ac` --> 618
func eitherordate(fp structs.FingerPrint) structs.FingerPrint {
	// note that there is an infinite loop possibility: pickandrunparser() is how you got here
	// `(.*)( or \d.*\s)
	cleaned := strings.Split(fp.OrigDateString, " or ")
	// fmt.Println("eitherordate cleaned to:", cleaned[0])
	newfp := TakeFingerprint(cleaned[0])
	newfp.HasOR = false
	newfp = pickandrunparser(newfp)
	newfp.OrigDateString = fp.OrigDateString
	newfp.HasOR = true
	return newfp
}

// multislashdate - 109／108／106／105 BC --> -109
func multislashdate(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := strings.Split(fp.OrigDateString, "／")
	cleaned[0] = stripalmostallstrings.ReplaceAllString(cleaned[0], "")
	d, e := strconv.Atoi(cleaned[0])
	if e != nil {
		// it might be 'AD 652／3／667／8' instead
		c2 := strings.Split(cleaned[0], " ")
		if len(c2) < 2 {
			global.MSG(fmt.Sprintf("multislashdate parser failed '%s'", fp.OrigDateString))
			fp.ParserFailed = true
			return fp
		}
		d2, e2 := strconv.Atoi(c2[1])
		if e2 != nil {
			global.MSG(fmt.Sprintf("multislashdate parser failed '%s'", fp.OrigDateString))
			fp.ParserFailed = true
			return fp
		}
		d = d2
	}
	fp.Calculated = d
	fp.ApplyBCE()
	fp.ApplySimpleFudges()
	return fp
}

// multicommadate - 114, 116, ﹠ 156 ac --> 114
func multicommadate(fp structs.FingerPrint) structs.FingerPrint {
	cleaned := strings.Split(fp.OrigDateString, ",")
	cleaned[0] = stripalmostallstrings.ReplaceAllString(cleaned[0], "")
	d, e := strconv.Atoi(cleaned[0])
	if e != nil {
		global.MSG(fmt.Sprintf("multicommadate parser failed '%s'", fp.OrigDateString))
		fp.ParserFailed = true
		return fp
	}
	fp.Calculated = d
	fp.ApplyBCE()
	fp.ApplySimpleFudges()
	return fp
}

// andsigndate - 1299 ﹠ 1344 ac --> 1299
func andsigndate(fp structs.FingerPrint) structs.FingerPrint {
	oldstr := fp.OrigDateString
	fp.OrigDateString = strings.Split(fp.OrigDateString, "﹠")[0]
	fp = onearabicsimple(fp)
	fp.OrigDateString = oldstr
	return multicommadate(fp)
}

// mixedspans - `245-244／220-219 BC` --> 245
func mixedspans(fp structs.FingerPrint) structs.FingerPrint {
	// note that there is an infinite loop possibility: pickandrunparser() is how you got here
	cleaned := strings.Split(fp.OrigDateString, "／")
	newfp := TakeFingerprint(cleaned[0])
	newfp = pickandrunparser(newfp)
	newfp.OrigDateString = fp.OrigDateString
	return newfp
}

// bracketdate - `med Ia [K.87(14)]` --> 50
func bracketdate(fp structs.FingerPrint) structs.FingerPrint {
	// note that there is an infinite loop possibility: pickandrunparser() is how you got here
	cleaned := hasbracket.ReplaceAllString(fp.OrigDateString, "")
	newfp := TakeFingerprint(cleaned)
	newfp = pickandrunparser(newfp)
	newfp.OrigDateString = fp.OrigDateString
	return newfp
}
