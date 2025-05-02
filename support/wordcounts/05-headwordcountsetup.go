//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import (
	"context"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/jackc/pgx/v5"
	"strings"
)

type tmpgramholder struct {
	Observed      string
	Possibilities string
}

// PrepareAnalysisLookupMap - map all observed forms to slices of headword possibilities; will fill HeadwordLookupMap
func PrepareAnalysisLookupMap() map[string][]string {
	const (
		Q1 = `SELECT observed_form, related_headwords FROM greek_morphology`
		Q2 = `SELECT observed_form, related_headwords FROM latin_morphology`
	)

	global.SECT("PrepareAnalysisLookupMap()")

	// greek_morphology holds GramAnalyses
	// i.e.
	// 	Observed      string
	//	XREFS         []string
	//	PFXXREFS      string // currently always ""; fix or remove...
	//	Possibilities []MorphPossib
	//	Related       string
	//	Scratchpad    string

	// Related is a " " join of the Headwords in Possibilities []MorphPossib

	parsemap := make(map[string][]string)
	grk := getanalysisfinds(Q1)
	parsemap = populateparsemap(parsemap, grk)
	lat := getanalysisfinds(Q2)
	parsemap = populateparsemap(parsemap, lat)

	return parsemap
}

func getanalysisfinds(q string) []tmpgramholder {
	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	foundrows, err := dbconn.Query(context.Background(), q)
	if err != nil {
		fmt.Println(err)
	}

	thesefinds, err := pgx.CollectRows(foundrows, pgx.RowToStructByPos[tmpgramholder])
	if err != nil {
		fmt.Println("PrepareAnalysisLookupMap()", err)
	}
	return thesefinds
}

func populateparsemap(parsemap map[string][]string, poss []tmpgramholder) map[string][]string {
	// should only see unique values in poss, so we are not going to have key collisions/overwrites in parsemap
	// we send the blank map here; fill it with greek; then fill it with latin
	for _, pos := range poss {
		// ultimately goes back to a " " join of the Headwords in Possibilities []MorphPossib
		parsemap[pos.Observed] = strings.Split(pos.Possibilities, " ")
	}
	return parsemap
}
