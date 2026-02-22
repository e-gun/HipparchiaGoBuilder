//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
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

func OneAuthorIntoAuthorTable(au *structs.DbAuthor) {
	const (
		FAIL = `"InsertOneAuthorIntoAuthorTable() failed to insert %s into authors"`
		INS  = `
			INSERT INTO authors
				(universalid, 
				 language, 
				 idxname, 
				 akaname, 
				 shortname,
				 cleanname,
				 genres,
				 recorded_date,
				 converted_date,
				 location)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (universalid)
			DO UPDATE SET
			     universalid = EXCLUDED.universalid,
				 language = EXCLUDED.language, 
				 idxname = EXCLUDED.idxname, 
				 akaname = EXCLUDED.akaname, 
				 shortname = EXCLUDED.shortname,
				 cleanname = EXCLUDED.cleanname,
				 genres = EXCLUDED.genres,
				 recorded_date = EXCLUDED.recorded_date,
				 converted_date = EXCLUDED.converted_date,
				 location = EXCLUDED.location
			`
	)

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	// UPSERT: update if the author already exists
	_, err := dbconn.Exec(context.Background(), INS,
		au.UID,
		au.Language,
		au.IDXname,
		au.Name,
		au.Shortname,
		au.Cleaname,
		au.Genres,
		au.RecDate,
		au.ConvDate,
		au.Location)

	if err != nil {
		// note that you probably just broke the ability of HGS to use this build
		// there will be works that cannot be mapped to an author and this will yield
		// null pointer problems
		fmt.Println(fmt.Sprintf(FAIL, au.UID))
		fmt.Println(err)
		au.PrintOut()
	}
}

func BulkInsertIntoAuthorsTable(auu []structs.DbAuthor) {
	var (
		authorsrownames = []string{
			"universalid",
			"language",
			"idxname",
			"akaname",
			"shortname",
			"cleanname",
			"genres",
			"recorded_date",
			"converted_date",
			"location",
		}
	)
	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tablename := "authors"
	rows := make([][]interface{}, len(auu))

	for i, a := range auu {
		rows[i] = populateoneatr(a)
	}
	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, authorsrownames, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
		os.Exit(0)
	} else {
		fmt.Println("BulkInsertIntoAuthorsTable() inserted", len(rows), "authors")
	}

	// fmt.Printf("BulkInsertWorks(): %s\tNumber of rows copied: %d\n", tablename, numbercopied)

	return
}

func populateoneatr(au structs.DbAuthor) []interface{} {
	const (
		NUMBEROFFIELDS = 10
	)

	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = au.UID
	row[1] = au.Language
	row[2] = au.IDXname
	row[3] = au.Name
	row[4] = au.Shortname
	row[5] = au.Cleaname
	row[6] = au.Genres
	row[7] = au.RecDate
	row[8] = au.ConvDate
	row[9] = au.Location

	return row
}
