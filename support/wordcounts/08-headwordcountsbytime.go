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
	"github.com/jackc/pgx/v5"
	"time"
)

// Headwords are sorted into three eras. Early or ⓔ is 900BCE to 300BCE. Middle or ⓜ is 299BCE to 300CE. Late or ⓛ is 301CE to 1600CE.

var (
	eras = map[string][2]int{
		"early":  {-900, -300},
		"middle": {-299, 300},
		"late":   {301, 1600},
	}
)

// ParsedCountOfAllEras - count the N instances of parsed word W in era ⓔ ...
func ParsedCountOfAllEras() map[string]structs.DbHeadwordCounts {
	const (
		MSG = "ParsedCountOfAllEras(): All eras processed. %.3fs"
	)
	global.SECT("ParsedCountOfAllEras()")

	start := time.Now()

	alleras := make(map[string]map[string]int, len(global.KnownGenres))

	for k, _ := range eras {
		fmt.Printf("\tEra wordcount working on '%s'\n", k)
		thisera := getoneera(k)
		alleras[k] = thisera
	}

	wordsbyera := convertallerasmaptowordgenremap(alleras)
	wcstructs := converterawcmaptohwcstruct(wordsbyera)

	// wcstructs["αἰϲχρόϲ"].PrintOut("αἰϲχρόϲ")

	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
	return wcstructs
}

func getoneera(era string) map[string]int {
	wbhh := getrelevanttimespanworks(era)
	doparsing := true
	thisera := FanoutGenreAndTimeCounter(wbhh, doparsing)
	return thisera
}

func getrelevanttimespanworks(era string) []WorksAndBoundsHolder {
	const (
		Q = `SELECT universalid,firstline,lastline FROM works WHERE converted_date BETWEEN $1 AND $2;`
	)

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	lowandhigh, _ := eras[era]
	query := fmt.Sprintf(Q)

	foundrows, err := dbconn.Query(context.Background(), query, lowandhigh[0], lowandhigh[1])
	if err != nil {
		fmt.Println(err)
	}

	thesefinds, err := pgx.CollectRows(foundrows, pgx.RowToStructByPos[WorksAndBoundsHolder])
	if err != nil {
		fmt.Println(err)
	}

	// at this point you have
	// Acta:
	// [{gr2038w002 3438 3604} {gr0031w005 11155 14338} {gr0388w004 658 1053} {gr2948w003 2267 2603} ...
	// but the queries are against 'gr2038' and not 'gr2038w002'; fix that
	for i, f := range thesefinds {
		thesefinds[i].T = f.T[0:6]
	}

	// now:
	// [{gr2038 3438 3604} {gr0031 11155 14338} {gr0388 658 1053} {gr2948 2267 2603} {gr2038 3728 4275} ...

	return thesefinds
}

func convertallerasmaptowordgenremap(alleras map[string]map[string]int) map[string]map[string]int {
	// see convertallgenremaptowordgenremap() comments
	wordsbyera := make(map[string]map[string]int)
	for era, ewc := range alleras {
		for word, count := range ewc {
			if _, ok := wordsbyera[word]; !ok {
				wordsbyera[word] = map[string]int{era: count}
			} else {
				wordsbyera[word][era] = count
			}
		}
	}
	return wordsbyera
}

func converterawcmaptohwcstruct(wordsbyera map[string]map[string]int) map[string]structs.DbHeadwordCounts {
	allwords := make(map[string]structs.DbHeadwordCounts, len(wordsbyera))
	for word, ec := range wordsbyera {
		thishwc := structs.DbHeadwordCounts{
			Word: word,
		}
		for era, count := range ec {
			switch era {
			case "early":
				thishwc.Early = count
			case "middle":
				thishwc.Middle = count
			case "late":
				thishwc.Late = count
			default:
				fmt.Println("Unknown era")
			}
		}
		allwords[word] = thishwc
	}
	return allwords
}
