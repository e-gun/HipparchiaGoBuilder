//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dating

import (
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func nodigitdispatcher(fp structs.FingerPrint) structs.FingerPrint {
	// parser failures: 23 of 513
	newfp := refingerprint(fp)
	newfp = readbackwards(newfp)

	newfp.ApplySimpleFudges()
	newfp.OrigDateString = fp.OrigDateString
	fp.Calculated = newfp.Calculated
	return fp
}

func refingerprint(fp structs.FingerPrint) structs.FingerPrint {
	ds := fp.OrigDateString
	ds = strings.ReplaceAll(ds, "?", "")
	if strings.Contains(ds, "a aet") {
		fp.ManualFudgeFactor = -100
		ds = strings.ReplaceAll(ds, "a aet", "")
	}
	if strings.Contains(ds, "p aet") {
		fp.ManualFudgeFactor = 100
		ds = strings.ReplaceAll(ds, "p aet", "")
	}
	if strings.Contains(ds, "init aet") {
		fp.ManualFudgeFactor = -50
		ds = strings.ReplaceAll(ds, "a aet", "")
	}
	if strings.Contains(ds, "late") {
		fp.ManualFudgeFactor = 50
		ds = strings.ReplaceAll(ds, "a aet", "")
	}
	if strings.Contains(ds, "post") {
		fp.ManualFudgeFactor = 50
	}
	if strings.Contains(ds, "ante") {
		fp.ManualFudgeFactor = -50
	}
	if strings.Contains(ds, "fin") {
		// all emperors...
		fp.ManualFudgeFactor = 15
		ds = strings.ReplaceAll(ds, "a aet", "")
	}
	if strings.Contains(ds, "init") {
		// almost all emperors...
		fp.ManualFudgeFactor = -15
		ds = strings.ReplaceAll(ds, "a aet", "")
	}

	fp.OrigDateString = ds
	return fp
}

// readbackwards - take a date string, split it, and look for a match in the map...
func readbackwards(fp structs.FingerPrint) structs.FingerPrint {
	elements := strings.Split(fp.OrigDateString, " ")
	for _, element := range elements {
		val, found := aetates[element]
		if found {
			fp.Calculated = val
			break
		}
	}
	if fp.Calculated == 9999 {
		fp.ParserFailed = true
	}
	return fp
}
