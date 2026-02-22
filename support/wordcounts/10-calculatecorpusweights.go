//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import (
	"context"
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
)

func CalculateUnparsedWordcountWeights() {
	const (
		EN1 = `__wordcounttotals`
		Q1  = `DELETE FROM unparsed_wordcounts WHERE entry_name = '%s';`
		Q2  = `select sum (%s) from unparsed_wordcounts;`
	)

	var (
		//toget = []string{
		//	"total_count",
		//	"gr_count",
		//	"lt_count",
		//	"dp_count",
		//	"in_count",
		//	"ch_count",
		//}
		toget = insert.UnparsedWordCountsTableRows[1:len(insert.UnparsedWordCountsTableRows)]
	)

	global.SECT("CalculateUnparsedWordcountWeights()")

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	// reset the old data
	_, err := dbconn.Exec(context.Background(), fmt.Sprintf(Q1, EN1))
	if err != nil {
		fmt.Println(err)
	}

	// get the new data
	var allsums []int
	for _, get := range toget {
		var thissum int
		s := dbconn.QueryRow(context.Background(), fmt.Sprintf(Q2, get))
		e := s.Scan(&thissum)
		if e != nil {
			fmt.Println(e)
		}
		allsums = append(allsums, thissum)
	}

	wct := structs.DbUnparsedWordCounts{
		Word:  EN1,
		Total: allsums[0],
		TLG:   allsums[1],
		LAT:   allsums[2],
		DDP:   allsums[3],
		INS:   allsums[4],
		CHR:   allsums[5],
	}

	wct.PrintOut()
	// __wordcounttotals	t: 88470869	g: 72722744	l: 7382222	i: 3493172	d: 4009321	c: 863410

	// i.e. weights are
	// 72722744	1
	// 7382222	9.85106435433668
	// 4009321	18.1384189492435
	// 3493172	20.8185408562762
	// 863410	84.2273589604012

	// set the new data
	towrite := map[string]structs.DbUnparsedWordCounts{EN1: wct}
	insert.InsertUnparsedWordCountsIntoTable(towrite)

}
