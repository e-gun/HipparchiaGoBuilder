//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"encoding/json"
	"fmt"
	"log"
)

type HeadwordAndForms struct {
	DictEntry     string
	XREF          string
	DerivFormsStr string
	DerivFormsSlc []string
}

func (hw *HeadwordAndForms) DFNotEmpty() []string {
	var ne []string
	for _, f := range hw.DerivFormsSlc {
		if f != "" {
			ne = append(ne, f)
		}
	}
	return ne
}

type GramAnalysis struct {
	Observed      string
	XREFS         []string
	PFXXREFS      string // currently always ""; fix or remove...
	Possibilities []MorphPossib
	Related       string
	Scratchpad    string
}

func (ga *GramAnalysis) PrintOut() {
	fmt.Printf("GramAnalysis.Printout() for '%s':\n", ga.Observed)
	fmt.Printf("\t%d Possibilies\n", len(ga.Possibilities))
	for _, f := range ga.Possibilities {
		fmt.Printf("\t%s (%s)\n", f.Headwd, f.Analysis)
	}
}

func (ga *GramAnalysis) PossibilitiesJSONOut() []byte {
	// on other end see rt-lexica.go for extractmorphpossibilities(raw string) []str.MorphPossib
	//	// Input:     {"1": {"transl": "A.I. stem, tree; II. shaft of a spear", "analysis": "neut nom/voc/acc sg", "headword": "δόρυ", "scansion": "", "xref_kind": "9", "xref_value": "26874791"}}
	//	// Unmarshal: map[1:{A.I. stem, tree; II. shaft of a spear neut nom/voc/acc sg δόρυ  9 26874791}]

	// need to match that

	pmap := make(map[int]MorphPossib)
	for i, p := range ga.Possibilities {
		pmap[i] = p
	}

	js, err := json.Marshal(pmap)
	if err != nil {
		log.Fatal(err)
	}
	return js
}

type MorphPossib struct {
	Transl   string `json:"transl"`
	Analysis string `json:"analysis"`
	Headwd   string `json:"headword"`
	Scansion string `json:"scansion"`
	Xrefkind string `json:"xref_kind"`
	Xrefval  string `json:"xref_value"`
}
