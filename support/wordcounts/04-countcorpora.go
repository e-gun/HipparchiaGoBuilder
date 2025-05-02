//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import (
	"context"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/resetdb"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

// UnparsedCountOfAllCorpora - count the N instances of unparsed word W in corpus C ...
func UnparsedCountOfAllCorpora() {
	const (
		MSG = "UnparsedCountOfAllCorpora(): All corpora processed. %.3fs"
	)

	global.SECT("UnparsedCountOfAllCorpora()")

	doparsing := false

	start := time.Now()

	tlg := CountTLGWords(doparsing)
	lat := CountLATWords(doparsing)
	ins := CountINSWords(doparsing)
	ddp := CountDDPWords(doparsing)
	chr := CountCHRWords(doparsing)
	mm := MapMerger([]map[string]int{tlg, lat, ins, ddp, chr})
	fmt.Println("size of full merged map of words:", len(mm))

	ccounts := make(map[string]structs.DbUnparsedWordCounts, len(mm))
	for k, v := range mm {
		ccounts[k] = structs.DbUnparsedWordCounts{
			Word:  k,
			Total: v,
		}
	}
	for k, v := range tlg {
		// every k in tlg should already have an entry in ccounts since mm was comprehensive
		wc := ccounts[k]
		wc.TLG = v
		ccounts[k] = wc
	}
	for k, v := range lat {
		wc := ccounts[k]
		wc.LAT = v
		ccounts[k] = wc
	}
	for k, v := range ins {
		wc := ccounts[k]
		wc.INS = v
		ccounts[k] = wc
	}
	for k, v := range ddp {
		wc := ccounts[k]
		wc.DDP = v
		ccounts[k] = wc
	}
	for k, v := range chr {
		wc := ccounts[k]
		wc.CHR = v
		ccounts[k] = wc
	}

	// ccounts["ϲτόματι"].PrintOut()

	resetdb.InitializeUnparsedWordCountTable()
	insert.InsertUnparsedWordCountsIntoTable(ccounts)

	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func CountTLGWords(needsparsing bool) map[string]int {
	auu := authorbuildpipeline.BuildAuthorsSlice(global.Config.GreekDir, "TLG")
	auu = authorbuildpipeline.SortWorkpileByFilesize(auu, global.Config.GreekDir)
	var cleanedauu []string
	for _, a := range auu {
		cleanedauu = append(cleanedauu, strings.ReplaceAll(a, "TLG", "gr"))
	}

	fmt.Println("\tCounting the words in the TLG author tables", len(cleanedauu))
	results := FanoutCorpusCount(needsparsing, cleanedauu)
	return results
}

func CountLATWords(needsparsing bool) map[string]int {
	auu := authorbuildpipeline.BuildAuthorsSlice(global.Config.LatDir, "LAT")
	auu = authorbuildpipeline.SortWorkpileByFilesize(auu, global.Config.GreekDir)
	var cleanedauu []string
	for _, a := range auu {
		cleanedauu = append(cleanedauu, strings.ReplaceAll(a, "LAT", "lt"))
	}

	fmt.Println("\tCounting the words in the LAT author tables", len(cleanedauu))
	results := FanoutCorpusCount(needsparsing, cleanedauu)
	return results
}

func CountINSWords(needsparsing bool) map[string]int {
	const (
		Q = `SELECT universalid FROM authors where universalid ~* '^in'`
	)

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()
	foundrows, err := dbconn.Query(context.Background(), Q)
	if err != nil {
		panic(err)
	}
	auu, err := pgx.CollectRows(foundrows, pgx.RowTo[string])
	fmt.Println("\tCounting the words in the INS author tables", len(auu))
	results := FanoutCorpusCount(needsparsing, auu)
	return results
}

func CountDDPWords(needsparsing bool) map[string]int {
	const (
		Q = `SELECT universalid FROM authors where universalid ~* '^dp'`
	)

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()
	foundrows, err := dbconn.Query(context.Background(), Q)
	if err != nil {
		panic(err)
	}
	auu, err := pgx.CollectRows(foundrows, pgx.RowTo[string])
	fmt.Println("\tCounting the words in the DDP author tables", len(auu))
	results := FanoutCorpusCount(needsparsing, auu)
	return results
}

func CountCHRWords(needsparsing bool) map[string]int {
	const (
		Q = `SELECT universalid FROM authors where universalid ~* '^ch'`
	)

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()
	foundrows, err := dbconn.Query(context.Background(), Q)
	if err != nil {
		panic(err)
	}
	auu, err := pgx.CollectRows(foundrows, pgx.RowTo[string])
	fmt.Println("\tCounting the words in the CHR author tables", len(auu))
	results := FanoutCorpusCount(needsparsing, auu)
	return results
}
