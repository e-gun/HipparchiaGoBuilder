//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/jackc/pgx/v5"
)

var (
	morphtablerownames = []string{
		"observed_form",
		"xrefs",
		"prefixrefs",
		"possible_dictionary_forms",
		"related_headwords",
	}
	lemmatatablerownames = []string{
		"dictionary_entry",
		"xref_number",
		"derivative_forms",
	}
)

func InsertLemmata(lang string, entries []structs.HeadwordAndForms) error {
	if len(entries) == 0 {
		return errors.New("insertLemmata(): No entries to insert")
		// return nil
	}
	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tablename := lang + "_lemmata"

	rows := make([][]interface{}, len(entries))

	// need to avoid: greek_lemmata ERROR: COPY from stdin failed: unable to encode "" into binary format for int4 (OID 23): cannot find encode plan (SQLSTATE 57014)
	// a single row will empty if you are not careful

	for i, ent := range entries {
		rows[i] = populateonelemmarow(ent)
	}

	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, lemmatatablerownames, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
	} else {
		fmt.Println("InsertLemmata() Inserted ", len(rows), " rows")
	}

	// fmt.Printf("InsertLemmata(): %s\tNumber of rows copied: %d\n", tablename, numbercopied)

	return nil
}

func populateonelemmarow(entry structs.HeadwordAndForms) []interface{} {
	const (
		NUMBEROFFIELDS = 3
	)

	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = entry.DictEntry
	row[1] = entry.XREF
	row[2] = entry.DFNotEmpty()

	return row
}

func InsertMorphology(lang string, entries []structs.GramAnalysis) error {
	if len(entries) == 0 {
		return errors.New("insertMorphology(): No entries to insert")
	}

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tablename := lang + "_morphology"

	rows := make([][]interface{}, len(entries))

	for i, ent := range entries {
		rows[i] = populateonemorphrow(ent)
	}

	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, morphtablerownames, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
	} else {
		fmt.Println("InsertMorphology Inserted ", len(rows), " rows")
	}

	// fmt.Printf("InsertMorphology(): %s\tNumber of rows copied: %d\n", tablename, numbercopied)

	return nil
}

func populateonemorphrow(entry structs.GramAnalysis) []interface{} {
	const (
		NUMBEROFFIELDS = 5
	)

	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = entry.Observed
	row[1] = strings.Join(entry.XREFS, " ")
	row[2] = entry.PFXXREFS
	row[3] = entry.PossibilitiesJSONOut()
	row[4] = entry.Related

	return row
}
