//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"context"
	"fmt"
	"os"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/jackc/pgx/v5"
)

func BulkInsertWorks(works []structs.DbWork) {
	var (
		wktablerownames = []string{
			"universalid",
			"title",
			"language",
			"publication_info",
			"levellabels_00",
			"levellabels_01",
			"levellabels_02",
			"levellabels_03",
			"levellabels_04",
			"levellabels_05",
			"workgenre",
			"transmission",
			"worktype",
			"provenance",
			"recorded_date",
			"converted_date",
			"wordcount",
			"firstline",
			"lastline",
			"authentic",
		}
	)
	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tablename := "works"
	rows := make([][]interface{}, len(works))

	for i, w := range works {
		rows[i] = populateonewktablerow(w)
	}

	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, wktablerownames, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
		os.Exit(0)
	}

	fmt.Println("BulkInsertWorks() inserted", len(rows), "works")

	// fmt.Printf("BulkInsertWorks(): %s\tNumber of rows copied: %d\n", tablename, numbercopied)

	return
}

func populateonewktablerow(work structs.DbWork) []interface{} {
	const (
		NUMBEROFFIELDS = 20
	)

	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = work.UID
	row[1] = work.Title
	row[2] = work.Language
	row[3] = work.Pub
	row[4] = work.LL0
	row[5] = work.LL1
	row[6] = work.LL2
	row[7] = work.LL3
	row[8] = work.LL4
	row[9] = work.LL5
	row[10] = work.Genre
	row[11] = work.Xmit
	row[12] = work.Type
	row[13] = work.Prov
	row[14] = work.RecDate
	row[15] = work.ConvDate
	row[16] = work.WdCount
	row[17] = work.FirstLine
	row[18] = work.LastLine
	row[19] = work.Authentic
	return row
}
