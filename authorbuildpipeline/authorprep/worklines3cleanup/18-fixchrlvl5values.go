//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines3cleanup

import (
	"strconv"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func fixchrlvl5values(lines []structs.DbWorkline) []structs.DbWorkline {
	previous := structs.DbWorkline{}
	lastsensibleline := structs.DbWorkline{}

	for i, line := range lines {
		// fmt.Println(line.TbIndex, line.WkUID)
		if line.Lvl5Value != previous.Lvl5Value {
			// maybe you do nothing: 414 -->  415 is uninteresting
			// but 417 --> 1 is what we are worried about
			// if 1 follows 417, this is 417A
			// if 2 follows 417, this is 417B
			// not clear that you *ever* see 3 or 4
			l5int, _ := strconv.Atoi(line.Lvl5Value)
			p5int, _ := strconv.Atoi(previous.Lvl5Value)

			if l5int < p5int {
				switch l5int {
				case 1:
					lines[i].Lvl5Value = lastsensibleline.Lvl5Value + "A"
				case 2:
					lines[i].Lvl5Value = lastsensibleline.Lvl5Value + "B"
				case 3:
					lines[i].Lvl5Value = lastsensibleline.Lvl5Value + "C"
				case 4:
					lines[i].Lvl5Value = lastsensibleline.Lvl5Value + "D"
				default:
					//fmt.Println("fixchrlvl5values() was confused by a shift")
					//fmt.Println("l5int < p5int", l5int, p5int)
					//fmt.Println(previous)
					//fmt.Println(line)

					// the following is ok and seems to be typical

					// l5int < p5int 24 520
					// {79861 520 1 1 1 1 12 ch0130w034 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;di.    }
					// {79862 24 1 1 1 1 1 ch0130w035 d(omi)ne d(eu)s﹖. domine deus domine deus  region: Orkney [Sct.] · city: Papa Stronsay Isl. · publicationinfo: H.214 · documentnumber: og · workabbrev: ECMScot III }
				}
				// fmt.Println("fixchrlvl5values()", lines[i].TbIndex, lines[i].WkUID, lines[i].Lvl5Value)

				// the following is ok and seems to be typical

				// fixchrlvl5values() 78842 ch0130w032 46A
				// fixchrlvl5values() 78848 ch0130w032 54A
				// fixchrlvl5values() 79449 ch0130w032 403A
				// fixchrlvl5values() 79463 ch0130w033 410A
				// fixchrlvl5values() 79676 ch0130w034 19A
				// fixchrlvl5values() 79723 ch0130w034 472A

			} else {
				lastsensibleline = line
			}
		}
		previous = line
	}
	return lines
}
