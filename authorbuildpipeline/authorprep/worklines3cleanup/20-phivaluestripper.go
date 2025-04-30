//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines3cleanup

import "github.com/e-gun/HipparchiaGoBuilder/internal/structs"

func PhiValueStripper(lines []structs.DbWorkline) []structs.DbWorkline {
	// inscription level values are nuts and do all sorts of things with the values
	// lvl5 stores document numbers; but intermediate levels are often but not always blank too

	// let's try deleting lvl5 and all 1s above lvl0 (unless there is a not-1 higher than the 1)

	for i, line := range lines {
		line.Lvl5Value = "-1"
		if line.Lvl1Value == "1" && line.Lvl2Value == "1" && line.Lvl3Value == "1" && line.Lvl4Value == "1" {
			line.Lvl1Value = "-1"
			line.Lvl2Value = "-1"
			line.Lvl3Value = "-1"
			line.Lvl4Value = "-1"
		} else if line.Lvl2Value == "1" && line.Lvl3Value == "1" && line.Lvl4Value == "1" {
			line.Lvl2Value = "-1"
			line.Lvl3Value = "-1"
			line.Lvl4Value = "-1"
		} else if line.Lvl3Value == "1" && line.Lvl4Value == "1" {
			line.Lvl3Value = "-1"
			line.Lvl4Value = "-1"
		} else if line.Lvl4Value == "1" {
			line.Lvl4Value = "-1"
		}
		lines[i] = line
	}
	return lines
}
