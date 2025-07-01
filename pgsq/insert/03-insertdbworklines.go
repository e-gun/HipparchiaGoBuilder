//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"context"
	"errors"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/jackc/pgx/v5"
	"os"
)

var (
	autablerownames = []string{
		"index",
		"wkuniversalid",
		"level_05_value",
		"level_04_value",
		"level_03_value",
		"level_02_value",
		"level_01_value",
		"level_00_value",
		"marked_up_line",
		"accented_line",
		"stripped_line",
		"hyphenated_words",
		"annotations",
	}
)

func InsertWorklinesIntoTable(lines []structs.DbWorkline) error {
	if len(lines) == 0 {
		return errors.New("No worklines to insert.")
	}

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tablename := lines[0].AuID()
	CreateAuthorTable(tablename)

	rows := make([][]interface{}, len(lines))

	for i, line := range lines {
		rows[i] = populateoneautablerow(line)
	}

	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, autablerownames, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
		os.Exit(0)
	} else {
		// fmt.Println("Inserted ", len(rows), " rows")
	}
	return nil
}

// populateoneautablerow - build a slice to feed to dbconn.CopyFrom()
func populateoneautablerow(line structs.DbWorkline) []interface{} {
	const (
		NUMBEROFFIELDS = 13
	)

	// should be disabled after parser bugs get ironed out
	// checklinesforerrors(line)

	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = line.TbIndex
	row[1] = line.WkUID
	row[2] = line.Lvl5Value
	row[3] = line.Lvl4Value
	row[4] = line.Lvl3Value
	row[5] = line.Lvl2Value
	row[6] = line.Lvl1Value
	row[7] = line.Lvl0Value
	row[8] = line.MarkedUp
	row[9] = line.Accented
	row[10] = line.Stripped
	row[11] = line.Hyphenated
	row[12] = line.Annotations
	return row
}

// checklinesforerrors - are you about to insert lines that will not fit? you will fail the whole table; this lets you see exactly why/where
func checklinesforerrors(line structs.DbWorkline) {
	if len(line.Lvl0Value) > 16 {
		fmt.Println("L0", line.WkUID, len(line.Lvl0Value), line.Lvl0Value)
	}
	if len(line.Lvl1Value) > 16 {
		fmt.Println("L1", line.WkUID, len(line.Lvl1Value), line.Lvl1Value)
	}
	if len(line.Lvl2Value) > 16 {
		fmt.Println("L2", line.WkUID, len(line.Lvl2Value), line.Lvl2Value)
	}
	if len(line.Lvl3Value) > 16 {
		fmt.Println("L3", line.WkUID, len(line.Lvl3Value), line.Lvl3Value)
	}
	if len(line.Lvl4Value) > 16 {
		fmt.Println("L4", line.WkUID, len(line.Lvl4Value), line.Lvl4Value)
	}
	if len(line.Lvl5Value) > 16 {
		fmt.Println("L5", line.WkUID, len(line.Lvl5Value), line.Lvl5Value)
	}
	if len(line.Accented) > 350 {
		fmt.Println("Stripped", line.WkUID, len(line.Accented), line.Accented)
	}
	if len(line.Annotations) > 600 {
		fmt.Println("Annotations", line.WkUID, len(line.Annotations), line.Annotations)
	}
}

// CreateAuthorTable - DROP/CREATE and author table
func CreateAuthorTable(tablename string) {
	const (
		DROP   = `DROP TABLE IF EXISTS %s;` // but it should not exist because of DropAuthorTables() executed up top
		CREATE = `
CREATE TABLE public.%s (
    index integer DEFAULT nextval('public.%s'::regclass) NOT NULL,
    wkuniversalid character varying(11),
    level_05_value character varying(24),
    level_04_value character varying(24),
    level_03_value character varying(24),
    level_02_value character varying(24),
    level_01_value character varying(24),
    level_00_value character varying(24),
    marked_up_line text,
    accented_line character varying(350),
    stripped_line character varying(350),
    hyphenated_words character varying(64),
    annotations character varying(680),
    UNIQUE(index)
);`
	)

	// that UNIQUE(index) constraint gives you a 5x speedup...

	// marked_up_line can get fairly long because of the things in the brackets (nb: we do not search this line...)
	// but accented_line and stripped_line will seldom be especially long; 256 chars will always work for GRK and LAT
	// 128 nearly works on GRK and LAT there is one failure...; 192 is supposed to be safe on GRK and LAT

	// but the CHR inscriptions have a small number of long lines: 337 is longest: καὶ τὴν βαϲιλείαν μου κατὰ το παρὸν...
	// DDP has many over 256, but longest is <325
	// INS has 330 as its max: θεῶν καὶ πτεροφόραι καὶ ἱερογραμματεῖϲ...

	// longest annotation is 669: dpz0qbw1g0 `corrected: <hb-sp-unconventional_form_written_by_scribe>(αλθεαιουϲ · ϲυντεύξειϲ ...`
	// moving unconventional form markup, etc into Annotations makes some of them get very big

	// level values that require more than 32 chars are suspect...

	// fmt.Println("CreateAuthorTable", tablename)

	queries := []string{
		fmt.Sprintf(DROP, tablename),
		fmt.Sprintf(CREATE, tablename, tablename),
	}

	dbc.DBCCommandSequence(queries)
}
