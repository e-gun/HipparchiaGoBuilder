//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/resetdb"
)

func DoAllWordcounts() {
	resetdb.InitializeHeadwordWordCountTable()
	resetdb.InitializeUnparsedWordCountTable()

	UnparsedCountOfAllCorpora()
	pcc := ParsedCountOfAllCorpora()
	pgc := ParsedCountOfAllGenres()
	ptc := ParsedCountOfAllEras()

	merged := make(map[string]structs.DbHeadwordCounts, len(pcc))
	for k, v := range pcc {
		// use the genre counts as the base and merge onto it
		wc, ok := pgc[k]
		if !ok {
			wc = v
		}

		wc.Total = v.Total
		wc.TLG = v.TLG
		wc.LAT = v.LAT
		wc.INS = v.INS
		wc.CHR = v.CHR
		wc.DDP = v.DDP

		tc, ok := ptc[k]
		if !ok {
			tc = v
		}

		wc.Early = tc.Early
		wc.Middle = tc.Middle
		wc.Late = tc.Late
		merged[k] = wc
	}

	for k, v := range merged {
		v.FrqClas = CalculateFrequencyClassification(k, v.Total)
		merged[k] = v
	}

	insert.InsertHeadwordWordcountsIntoTable(merged)

	// see the head of CalculateParsedWordcountTotals()
	// if you comment out the next two you will still have viable weights, but they will have a slightly different meaning
	CountEraRawWords()
	CountGenreRawWords()

	CalculateUnparsedWordcountWeights()
	CalculateParsedWordcountTotals()

	insert.InsertBuildMetadata("counts", fmt.Sprintf("%d headwords counted", len(merged)))
}

// hgdb=> select * from wordcounts where entry_name = 'παρθένοϲ';
// entry_name | total_count | gr_count | lt_count | dp_count | in_count | ch_count
//------------+-------------+----------+----------+----------+----------+----------
// παρθένοϲ   |        2820 |     2728 |        1 |        3 |       53 |       35
//(1 row)

// hgdb=> select * from wordcounts where entry_name = 'παρθένοιϲ';
// entry_name | total_count | gr_count | lt_count | dp_count | in_count | ch_count
//------------+-------------+----------+----------+----------+----------+----------
// παρθένοιϲ  |         342 |      315 |        0 |        0 |       27 |        0
//(1 row)

// hgdb=> select total_count,early_occurrences, middle_occurrences, late_occurrences  from headword_wordcounts where entry_name = 'παρθένοϲ';
// total_count | early_occurrences | middle_occurrences | late_occurrences
//-------------+-------------------+--------------------+------------------
//       12533 |               442 |               2769 |             7289
//(1 row)

// hgdb=> select total_count,comic,epic,lyr,trag  from headword_wordcounts where entry_name = 'παρθένοϲ';
// total_count | comic | epic | lyr | trag
//-------------+-------+------+-----+------
//       12533 |    61 |  180 |  41 |   82
//(1 row)
