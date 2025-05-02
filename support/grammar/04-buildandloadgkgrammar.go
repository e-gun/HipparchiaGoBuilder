//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/resetdb"
	"strings"
	"time"
)

func BuildAndLoadGreekGrammar(directory string, lemmatafilename string, analysisfilename string) {
	const (
		MSG1 = "Greek Lemmata Finished. %.3fs"
		MSG2 = "All grammar jobs processed. %.3fs"
	)
	start := time.Now()
	resetdb.InitializeGreekGrammarSupportTables()
	entries := LoadGrammarFile(directory, lemmatafilename)

	hwaf := FanoutGreekLemmata(entries)

	err := insert.InsertLemmata("greek", hwaf)
	if err != nil {
		fmt.Println(err)
	}

	d := fmt.Sprintf(MSG1, time.Now().Sub(start).Seconds())
	fmt.Println(d)

	entries = LoadGrammarFile(directory, analysisfilename)

	gram := FanoutGreekAnalysis(entries)
	gram = cleangreekalalyses(gram)

	err = insert.InsertMorphology("greek", gram)
	if err != nil {
		fmt.Println(err)
	}

	d = fmt.Sprintf(MSG2, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func cleangreekalalyses(gram []structs.GramAnalysis) []structs.GramAnalysis {
	for i, g := range gram {
		var rhw []string
		for j, _ := range g.Possibilities {
			rhw = append(rhw, g.Possibilities[j].Headwd)
		}
		gram[i].Related = strings.Join(rhw, " ")
	}
	return gram
}
