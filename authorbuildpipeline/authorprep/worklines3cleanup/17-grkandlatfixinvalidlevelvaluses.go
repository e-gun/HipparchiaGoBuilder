//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines3cleanup

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func GRKandLATFixInvalidLevelValues(lines []structs.DbWorkline, works []structs.DbWork) []structs.DbWorkline {

	corpus := works[0].GetCorpus()

	validator := make(map[string]int)
	for _, w := range works {
		validator[w.UID] = w.CountLevels()
	}

	// should always be true...
	if corpus == global.TLGABBREV || corpus == global.LATABBREV {
		for i, line := range lines {
			lines[i] = authorvaluestripper(validator[line.WkUID], line)
		}
	}

	if corpus == global.CHRFIRSTPASS {
		lines = fixchrlvl5values(lines)
	}

	return lines
}

func authorvaluestripper(valid int, line structs.DbWorkline) structs.DbWorkline {
	switch valid {
	case 0:
	// impossible
	case 1:
		line.Lvl5Value = "-1"
		line.Lvl4Value = "-1"
		line.Lvl3Value = "-1"
		line.Lvl2Value = "-1"
		line.Lvl1Value = "-1"
	case 2:
		line.Lvl5Value = "-1"
		line.Lvl4Value = "-1"
		line.Lvl3Value = "-1"
		line.Lvl2Value = "-1"
	case 3:
		line.Lvl5Value = "-1"
		line.Lvl4Value = "-1"
		line.Lvl3Value = "-1"
	case 4:
		line.Lvl5Value = "-1"
		line.Lvl4Value = "-1"
	case 5:
		line.Lvl5Value = "-1"
	case 6:
	// do nothing
	default:
		fmt.Printf("valuestripper(): Invalid level value: %d\n", valid)

	}
	return line
}
