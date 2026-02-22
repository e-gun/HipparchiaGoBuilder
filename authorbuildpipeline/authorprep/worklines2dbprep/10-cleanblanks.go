//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"regexp"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	blanklinepattern = regexp.MustCompile(`^<hb-set_l_0_to_.*? /><br> $`)
)

func CleanBlanks(lines []structs.DbWorkline) []structs.DbWorkline {
	// 	multiple sets and level shifts in rapid succession sometimes leaves blank lines that are 'numbered': zap them
	//
	//	note that this is an old, subtle issue. and it tends to shift as with changes to the functions that come befor it.
	//
	//	the defunct halflinecleanup() was a flawed fix for a slightly different version of the problem.
	//
	//	it may be necessary to keep debugging this: lyric is the best place to look for the issue
	//
	//	aristophanes:
	//
	//		['6', [('0', '1099'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_increment_level_0_by_1 /><hmu_blank_quarter_spaces quantity="63" /> ἠρινά τε βοϲκόμεθα παρθένια '],
	//		['6', [('0', '1100-1101'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_set_level_0_to_1100-1101 /><hmu_blank_quarter_spaces quantity="63" /> λευκότροφα μύρτα Χαρίτων τε κηπεύματα. '],
	//		['6', [('0', '1100-1101'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_set_level_0_to_1100-1101 /><br /> '],
	//		['6', [('0', '1102'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_set_level_0_to_1102 /><hmu_blank_quarter_spaces quantity="38" /> Τοῖϲ κριταῖϲ εἰπεῖν τι βουλόμεϲθα τῆϲ νίκηϲ πέρι, '],
	//		['6', [('0', '1103'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_increment_level_0_by_1 /><hmu_blank_quarter_spaces quantity="38" /> ὅϲ’ ἀγάθ’, ἢν κρίνωϲιν ἡμᾶϲ, πᾶϲιν αὐτοῖϲ δώϲομεν, '],
	//
	//	aeschylus:
	//
	//		['1', [('0', '175e'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_increment_level_0_by_1 /><hmu_blank_quarter_spaces quantity="27" /> χαλεποῦ γὰρ ἐκ ']
	//		['1', [('0', '175f'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_increment_level_0_by_1 /><hmu_blank_quarter_spaces quantity="43" /> πνεύματοϲ εἶϲι χειμών.⟩ ']
	//		['1', [('0', '175f'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_set_level_0_to_175f /><br /> ']
	//		['1', [('0', '176'), ('1', '1'), ('2', '1'), ('3', '1'), ('4', '1'), ('5', '1')], '<hmu_set_level_0_to_176 /><hmu_blank_quarter_spaces quantity="4" /> <span class="speaker"><span class="smallerthannormal">ΔΑΝΑΟϹ</span></span> ']

	skips := 0
	index := -1
	for _, line := range lines {
		if blanklinepattern.MatchString(line.MarkedUp) {
			skips++
		} else {
			index++
			lines[index] = line
		}
	}

	if skips > 0 {
		lines = lines[0 : len(lines)-skips]
	}

	// fmt.Println("CleanBlanks() skipped", skips)
	return lines
}
