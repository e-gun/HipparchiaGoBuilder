//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"regexp"
)

var (
	findsetters  = regexp.MustCompile(`<hb-set_l_\d_to_.*?\s/>`)
	findadders   = regexp.MustCompile(`<hb-incr_l_\d_by_1\s/>`)
	findnewauth  = regexp.MustCompile(`<hb-metadata_newauthor value="\d\d\d\d"`)
	findnewwk    = regexp.MustCompile(`<hb-metadata_newwork value="\d\d\d" />`)
	findwkabbrev = regexp.MustCompile(`workabbrev: ([^·]*)`)
)

func DeIncrement(lines []structs.DbWorkline) []structs.DbWorkline {
	//  strip out the "hmu_increment_" stuff
	//
	//	sample in:
	//		['1', [('0', '4'), ('1', '3'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_increment_level_0_by_1 /><hmu_standalone_tabbedtext />ταύτηϲ γὰρ κεῖνοι δάμονέϲ εἰϲι μάχηϲ ']
	//
	//	sample out:
	//		['1', [('0', '4'), ('1', '3'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_standalone_tabbedtext />ταύτηϲ γὰρ κεῖνοι δάμονέϲ εἰϲι μάχηϲ ']
	//

	// there is an interesting problem with inscriptions
	// we want to remap them and so are attentive to "newwork", etc
	// but that material is on "blank" lines that will be pruned because they are empty but for markup
	// so we need to capture this data and insert it later: the first lines of INS130:
	//        <hb-metadata_newauthor value="0130" />
	//        <hb-metadata_newwork value="001" /><hb-metadata_workabbrev value="IC 1:i" /><hb-metadata_authabbrev value="Crete" />
	//        <hb-set_l_5_to_1 />
	//        <hb-set_l_0_to_1 /><hb-metadata_region value="C. Crete" /><hb-metadata_city value="Aigaion Antron" /><hb-metadata_date value="aet Rom" />ΕΠΙ( — ).

	// by the time you get here that second line now looks like this in line.Annotations:
	//         workabbrev: IC 1:i · authabbrev: Crete

	skips := 0
	index := -1

	var storedworkabbrev string

	for _, line := range lines {
		// initial cleanup
		newline := findsetters.ReplaceAllString(line.MarkedUp, "")
		newline = findadders.ReplaceAllString(newline, "")
		newline = findnewwk.ReplaceAllString(newline, "")
		newline = findnewauth.ReplaceAllString(newline, "")
		if findwkabbrev.MatchString(line.Annotations) {
			storedworkabbrev = findwkabbrev.FindStringSubmatch(line.Annotations)[1]
		}

		// rebuild the line
		if newline == "" {
			skips++
		} else {
			if line.Annotations != "" && storedworkabbrev != "" {
				line.Annotations = line.Annotations + ` · workabbrev: ` + storedworkabbrev
				storedworkabbrev = ""
			}
			index++
			line.MarkedUp = newline
			lines[index] = line
		}
	}

	trim := len(lines) - skips
	if trim < 0 {
		trim = 0
	}

	lines = lines[0:trim]

	// fmt.Println("DeIncrement() skipped", skips)
	return lines
}
