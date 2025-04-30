//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

// PrepareForDB - move from strings to []DbWorkline; AssignCitationValues() and ExtractNotesAndMarginalText() are the key fncs
func PrepareForDB(wklinesasstring string) []structs.DbWorkline {
	// placeholder; feed me the end of the 'worklines1initial' authorbuildpipeline

	lines := AssignCitationValues(wklinesasstring)
	lines = LastLinefix(lines)
	lines = ExtractNotesAndMarginalText(lines)
	lines = CleanBlanks(lines)
	lines = DeIncrement(lines)
	lines = BuildStrippedAndAccentedlines(lines)
	lines = NoLeadingOrTrailingSpaces(lines)
	lines = FixHyphens(lines)
	lines = InsertNBSP(lines)

	//
	// VARIOUS DEBUGGING TOOLS....
	//

	// authorprep.WriteWorklineProgress(lines)
	//authorprep.WriteWorklineAnnotationsProgress(lines)
	// os.Exit(1)

	// ViewWorklines("InsertNBSP", false, true, lines)

	//
	// manually abort on a target line (useful for figuring out which INS file is an original source
	//

	//for i, line := range lines {
	//	if strings.Contains(line.MarkedUp, "πάντων δὲ τῶν δανειϲάντων") {
	//		fmt.Println("PrepareForDB() aborting")
	//		fmt.Println(lines[i-1])
	//		fmt.Println(lines[i])
	//		fmt.Println(lines[i+1])
	//		os.Exit(1)
	//	}
	//}

	return lines
}
