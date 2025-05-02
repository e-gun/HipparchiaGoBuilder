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

var (
	gramentrycleaner = strings.NewReplacer("v", "u", "j", "i", "-", "", "#", "", "1", "¹", "2", "²", "3", "³", "4", "⁴")
)

func BuildAndLoadLatinGrammar(directory string, lemmatafilename string, analysisfilename string) {
	const (
		MSG1 = "Latin Lemmata Finished. %.3fs"
		MSG2 = "All grammar jobs processed. %.3fs"
	)
	start := time.Now()
	resetdb.InitializeLatinGrammarSupportTables()
	entries := LoadGrammarFile(directory, lemmatafilename)

	hwaf := FanoutLatinLemmata(entries)
	hwaf = cleanlatinlemmata(hwaf)

	err := insert.InsertLemmata("latin", hwaf)
	if err != nil {
		fmt.Println(err)
	}

	d := fmt.Sprintf(MSG1, time.Now().Sub(start).Seconds())
	fmt.Println(d)

	entries = LoadGrammarFile(directory, analysisfilename)

	gram := FanoutLatinAnalysis(entries)
	gram = cleanlatinalalyses(gram)

	err = insert.InsertMorphology("latin", gram)
	if err != nil {
		fmt.Println(err)
	}

	d = fmt.Sprintf(MSG2, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func cleanlatinlemmata(hwaf []structs.HeadwordAndForms) []structs.HeadwordAndForms {
	for i, h := range hwaf {
		hwaf[i].DictEntry = gramentrycleaner.Replace(h.DictEntry)
		for j, df := range h.DerivFormsSlc {
			h.DerivFormsSlc[j] = gramentrycleaner.Replace(df)
		}
		hwaf[i].DerivFormsSlc = h.DerivFormsSlc
		hwaf[i].DerivFormsStr = strings.Join(h.DerivFormsSlc, " ")
	}
	return hwaf
}

func cleanlatinalalyses(gram []structs.GramAnalysis) []structs.GramAnalysis {
	for i, g := range gram {
		gram[i].Observed = gramentrycleaner.Replace(g.Observed)
		var rhw []string
		for j, pos := range g.Possibilities {
			g.Possibilities[j].Headwd = gramentrycleaner.Replace(pos.Headwd)
			g.Possibilities[j].Scansion = gramentrycleaner.Replace(pos.Scansion)
			rhw = append(rhw, g.Possibilities[j].Headwd)
		}
		gram[i].Possibilities = g.Possibilities
		gram[i].Related = strings.Join(rhw, " ")
	}
	return gram
}
