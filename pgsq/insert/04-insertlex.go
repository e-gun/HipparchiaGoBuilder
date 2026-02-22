//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"context"
	"fmt"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/jackc/pgx/v5"
)

var (
	lextablerownames = []string{
		"idval",
		"entry_name",
		"entry_metr",
		"idname",
		"entry_type",
		"translations",
		"usedby",
		"prelim_info",
		"sense_ids",
		"senses",
	}
)

func EntriesIntoGkLexicon(entries []structs.DbLexicon) error {
	tablename := "greek_dictionary"

	_ = EntriesIntoLexicon(tablename, entries)
	return nil
}

func EntriesIntoLatinLexicon(entries []structs.DbLexicon) error {
	tablename := "latin_dictionary"

	_ = EntriesIntoLexicon(tablename, entries)
	return nil
}

func EntriesIntoLexicon(tablename string, entries []structs.DbLexicon) error {
	if len(entries) == 0 {
		// return errors.New("InsertEntriesIntoGkLexicon(): No entries to insert.")
		return nil
	}

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rows := make([][]interface{}, len(entries))

	for i, ent := range entries {
		rows[i] = populateonelexrow(ent)
	}

	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, lextablerownames, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
	} else {
		// fmt.Println("Inserted ", len(rows), " rows")
	}

	// fmt.Printf("InsertEntriesIntoGkLexicon(): %s\tNumber of rows copied: %d\n", tablename, numbercopied)

	return nil
}

func populateonelexrow(entry structs.DbLexicon) []interface{} {
	const (
		NUMBEROFFIELDS = 10
	)
	// num, _ := strconv.Atoi(stripalpha.ReplaceAllString(entry.IdString, ""))

	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = entry.IdFloat
	row[1] = entry.EntryName
	row[2] = entry.EntryMetr
	row[3] = entry.IdString
	row[4] = entry.EntryType
	row[5] = entry.Transl
	row[6] = entry.Usedby
	row[7] = entry.PrelimInfo
	row[8] = strings.Join(entry.SenseIDs, " ")
	row[9] = entry.SensesJSONOut()
	return row
}
