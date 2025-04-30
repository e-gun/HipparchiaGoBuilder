//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"context"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/jackc/pgx/v5"
)

var (
	UnparsedWordCountsTableRows = []string{
		"entry_name",
		"total_count",
		"gr_count",
		"lt_count",
		"dp_count",
		"in_count",
		"ch_count",
	}
	HeadwordWordCountsRows = []string{
		"entry_name",
		"total_count",
		"gr_count",
		"lt_count",
		"dp_count",
		"in_count",
		"ch_count",
		"frequency_classification",
		"early_occurrences",
		"middle_occurrences",
		"late_occurrences",
		"acta",
		"agric",
		"alchem",
		"anthol",
		"apocalyp",
		"apocryph",
		"apol",
		"astrol",
		"astron",
		"biogr",
		"bucol",
		"caten",
		"chronogr",
		"comic",
		"comm",
		"concil",
		"coq",
		"dialog",
		"docu",
		"doxogr",
		"eccl",
		"eleg",
		"encom",
		"epic",
		"epigr",
		"epist",
		"evangel",
		"exeget",
		"fab",
		"geogr",
		"gnom",
		"gramm",
		"hagiogr",
		"hexametr",
		"hist",
		"homilet",
		"hymn",
		"hypoth",
		"iamb",
		"ignotum",
		"invectiv",
		"inscr",
		"jurisprud",
		"lexicogr",
		"liturg",
		"lyr",
		"magica",
		"math",
		"mech",
		"med",
		"metrolog",
		"mim",
		"mus",
		"myth",
		"narrfict",
		"nathist",
		"onir",
		"orac",
		"orat",
		"papyrus",
		"paradox",
		"parod",
		"paroem",
		"perieg",
		"phil",
		"physiognom",
		"poem",
		"polyhist",
		"prophet",
		"pseudepigr",
		"rhet",
		"satura",
		"satyr",
		"schol",
		"tact",
		"test",
		"theol",
		"trag",
		"allrhet",
		"allrel",
	}
)

func InsertUnparsedWordCountsIntoTable(wordcounts map[string]structs.DbUnparsedWordCounts) {

	tablename := "unparsed_wordcounts"

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rows := make([][]interface{}, len(wordcounts))

	count := 0
	for k, _ := range wordcounts {
		rows[count] = populateonwctablerow(wordcounts[k])
		count++
	}

	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, UnparsedWordCountsTableRows, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
	} else {
		fmt.Println("InsertUnparsedWordCounts() inserted", len(rows), "rows")
	}

	// fmt.Printf("InsertUnparsedWordCountsIntoTable(): %s\tNumber of rows copied: %d\n", tablename, numbercopied)
}

func InsertHeadwordWordcountsIntoTable(wordcounts map[string]structs.DbHeadwordCounts) {

	tablename := "headword_wordcounts"

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rows := make([][]interface{}, len(wordcounts))

	count := 0
	for k, v := range wordcounts {
		rows[count] = populateoneheadwordcounttablerow(k, v)
		count++
	}

	_, err := dbconn.CopyFrom(ctx, pgx.Identifier{tablename}, HeadwordWordCountsRows, pgx.CopyFromRows(rows))
	if err != nil {
		fmt.Println(tablename, err)
	} else {
		fmt.Println("InsertHeadwordWordcounts inserted", len(rows), "rows")
	}

	// fmt.Printf("InsertUnparsedWordCountsIntoTable(): %s\tNumber of rows copied: %d\n", tablename, numbercopied)
}

// populateonwctablerow - build a slice to feed to dbconn.CopyFrom()
func populateonwctablerow(word structs.DbUnparsedWordCounts) []interface{} {
	const (
		NUMBEROFFIELDS = 7
	)

	// should be disabled after parser bugs get ironed out
	// checklinesforerrors(word)

	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = word.Word
	row[1] = word.Total
	row[2] = word.TLG
	row[3] = word.LAT
	row[4] = word.DDP
	row[5] = word.INS
	row[6] = word.CHR
	return row
}

func populateoneheadwordcounttablerow(word string, wcs structs.DbHeadwordCounts) []interface{} {
	const (
		NUMBEROFFIELDS = 91
	)
	row := make([]interface{}, NUMBEROFFIELDS)
	row[0] = word
	row[1] = wcs.Total
	row[2] = wcs.TLG
	row[3] = wcs.LAT
	row[4] = wcs.DDP
	row[5] = wcs.INS
	row[6] = wcs.CHR
	row[7] = wcs.FrqClas
	row[8] = wcs.Early
	row[9] = wcs.Middle
	row[10] = wcs.Late
	row[11] = wcs.Acta
	row[12] = wcs.Agric
	row[13] = wcs.Alchem
	row[14] = wcs.Anthol
	row[15] = wcs.Apocal
	row[16] = wcs.Apocry
	row[17] = wcs.Apol
	row[18] = wcs.Astrol
	row[19] = wcs.Astron
	row[20] = wcs.Biogr
	row[21] = wcs.Bucol
	row[22] = wcs.Caten
	row[23] = wcs.Chron
	row[24] = wcs.Comic
	row[25] = wcs.Comm
	row[26] = wcs.Concil
	row[27] = wcs.Coq
	row[28] = wcs.Dial
	row[29] = wcs.Docu
	row[30] = wcs.Doxog
	row[31] = wcs.Eccl
	row[32] = wcs.Eleg
	row[33] = wcs.Encom
	row[34] = wcs.Epic
	row[35] = wcs.Epigr
	row[36] = wcs.Epist
	row[37] = wcs.Evang
	row[38] = wcs.Exeg
	row[39] = wcs.Fab
	row[40] = wcs.Geog
	row[41] = wcs.Gnom
	row[42] = wcs.Gram
	row[43] = wcs.Hagiog
	row[44] = wcs.Hexam
	row[45] = wcs.Hist
	row[46] = wcs.Homil
	row[47] = wcs.Hymn
	row[48] = wcs.Hypoth
	row[49] = wcs.Iamb
	row[50] = wcs.Ignot
	row[51] = wcs.Inscr
	row[52] = wcs.Invectiv
	row[53] = wcs.Liturg
	row[54] = wcs.Juris
	row[55] = wcs.Lexic
	row[56] = wcs.Lyr
	row[57] = wcs.Magica
	row[58] = wcs.Math
	row[59] = wcs.Mech
	row[60] = wcs.Med
	row[61] = wcs.Meteor
	row[62] = wcs.Mim
	row[63] = wcs.Mus
	row[64] = wcs.Myth
	row[65] = wcs.NarrFic
	row[66] = wcs.NatHis
	row[67] = wcs.Onir
	row[68] = wcs.Orac
	row[69] = wcs.Orat
	row[70] = wcs.Papyrus
	row[71] = wcs.Paradox
	row[72] = wcs.Parod
	row[73] = wcs.Paroem
	row[74] = wcs.Perig
	row[75] = wcs.Phil
	row[76] = wcs.Physiog
	row[77] = wcs.Poem
	row[78] = wcs.Polyhist
	row[79] = wcs.Proph
	row[80] = wcs.Pseud
	row[81] = wcs.Rhet
	row[82] = wcs.Satura
	row[83] = wcs.Satyr
	row[84] = wcs.Schol
	row[85] = wcs.Tact
	row[86] = wcs.Test
	row[87] = wcs.Theol
	row[88] = wcs.Trag
	row[89] = wcs.AllRhet
	row[90] = wcs.AllRelig
	return row
}
