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
	"strings"
	"time"
)

type WorksAndBoundsHolder struct {
	T string
	F int
	L int
}

// ParsedCountOfAllGenres - count the N instances of parsed word W in genre G ...
func ParsedCountOfAllGenres() map[string]structs.DbHeadwordCounts {
	const (
		MSG = "ParsedCountOfAllGenres(): All genres processed. %.3fs"
	)
	global.SECT("ParsedCountOfAllGenres()")

	start := time.Now()

	allgenres := make(map[string]map[string]int, len(global.KnownGenres))

	for i, genre := range global.KnownGenres {
		if i%20 == 0 {
			fmt.Printf("\tGenre wordcount working on %d of %d genres\n", i, len(global.KnownGenres))
		}
		thisgenre := getonegenre(genre)
		delete(thisgenre, "")
		allgenres[genre] = thisgenre
	}

	wordsbygenre := convertallgenremaptowordgenremap(allgenres)
	wcstructs := convertwcmaptohwcstruct(wordsbygenre)

	// wcstructs["αἰϲχρόϲ"].PrintOut("αἰϲχρόϲ")

	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
	return wcstructs
}

func getonegenre(genre string) map[string]int {
	genre = strings.Title(genre)
	wbhh := getrelevantgenreworks(genre)
	thisgenre := FanoutGenreAndTimeCounter(wbhh)
	return thisgenre
}

func getrelevantgenreworks(cat string) []WorksAndBoundsHolder {
	const (
		Q   = `SELECT universalid,firstline,lastline FROM works WHERE workgenre ~* $1;`
		RE  = `%s(|.|;)`
		Q2  = `SELECT universalid,firstline,lastline FROM works WHERE workgenre = $1;`
		RE2 = `%s.`
	)

	// Q2 & RE2 let you vaguely simulate the old wordcount situation where works had only one genre
	// they should not be used except for testing and for rough comparisons against the hipparchiaDB data

	// DbHeadwordCounts: αἰϲχρόϲ
	// AllRhet: 3375
	// AllRelig: 2111
	// Trag: 232
	// Epic: 31
	// Phil: 4360

	// but in hipparchiaDBthere are only 13 in epic, 160 in trag, and 723 in phil...
	// but works are now multi-genre, so this makes a big difference...

	// if you use Q2
	// DbHeadwordCounts: αἰϲχρόϲ
	// AllRhet: 985
	// AllRelig: 688
	// Trag: 141
	// Epic: 25
	// Phil: 1048

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	query := fmt.Sprintf(Q)
	genre := fmt.Sprintf(RE, cat)
	foundrows, err := dbconn.Query(context.Background(), query, genre)
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
	// but the quieries are against 'gr2038' and not 'gr2038w002'; fix that
	for i, f := range thesefinds {
		thesefinds[i].T = f.T[0:6]
	}

	// now:
	// [{gr2038 3438 3604} {gr0031 11155 14338} {gr0388 658 1053} {gr2948 2267 2603} {gr2038 3728 4275} ...

	return thesefinds
}

func convertallgenremaptowordgenremap(allgenres map[string]map[string]int) map[string]map[string]int {
	// at this point you have something like ...
	// allgenres: map[acta:map[:220 amitto:6 amplus:2 atque:1 b:1 bibliotheca:6 capitum:2 ...], epic:map[...] ... ]
	// but what you need is to know what genres γυνή, e.g., appears in; so it is time to remap
	// wordsbygenre: γυνή:[acta: 12, epic: 100, ...]
	wordsbygenre := make(map[string]map[string]int)
	for genre, gwc := range allgenres {
		for word, count := range gwc {
			if _, ok := wordsbygenre[word]; !ok {
				wordsbygenre[word] = map[string]int{genre: count}
			} else {
				wordsbygenre[word][genre] = count
			}
		}
	}
	return wordsbygenre
}

func convertwcmaptohwcstruct(wordsbygenre map[string]map[string]int) map[string]structs.DbHeadwordCounts {
	allwords := make(map[string]structs.DbHeadwordCounts, len(wordsbygenre))
	for word, gc := range wordsbygenre {
		thishwc := structs.DbHeadwordCounts{
			Word: word,
		}
		for genre, count := range gc {
			switch genre {
			case "acta":
				thishwc.AllRelig += count
				thishwc.Acta = count
			case "agric":
				thishwc.Agric = count
			case "alchem":
				thishwc.Alchem = count
			case "anthol":
				thishwc.Anthol = count
			case "apocalyp":
				thishwc.AllRelig += count
				thishwc.Apocal = count
			case "apocryph":
				thishwc.AllRelig += count
				thishwc.Apocry = count
			case "apol":
				thishwc.AllRelig += count
				thishwc.Apol = count
			case "astrol":
				thishwc.Astrol = count
			case "astron":
				thishwc.Astron = count
			case "biogr":
				thishwc.Biogr = count
			case "bucol":
				thishwc.Bucol = count
			case "caten":
				thishwc.AllRelig += count
				thishwc.Caten = count
			case "chronogr":
				thishwc.Chron = count
			case "comic":
				thishwc.Comic = count
			case "comm":
				thishwc.Comm = count
			case "concil":
				thishwc.AllRelig += count
				thishwc.Concil = count
			case "coq":
				thishwc.Coq = count
			case "dialog":
				thishwc.Dial = count
			case "docu":
				thishwc.Docu = count
			case "doxogr":
				thishwc.Doxog = count
			case "eccl":
				thishwc.AllRelig += count
				thishwc.Eccl = count
			case "eleg":
				thishwc.Eleg = count
			case "encom":
				thishwc.AllRhet += count
				thishwc.Encom = count
			case "epic":
				thishwc.Epic = count
			case "epigr":
				thishwc.Epigr = count
			case "epist":
				thishwc.Epist = count
			case "evangel":
				thishwc.AllRelig += count
				thishwc.Evang = count
			case "exeget":
				thishwc.AllRelig += count
				thishwc.Exeg = count
			case "fab":
				thishwc.Fab = count
			case "geogr":
				thishwc.Geog = count
			case "gnom":
				thishwc.Gnom = count
			case "gramm":
				thishwc.Gram = count
			case "hagiogr":
				thishwc.AllRelig += count
				thishwc.Hagiog = count
			case "hexametr":
				thishwc.Hexam = count
			case "hist":
				thishwc.Hist = count
			case "homilet":
				thishwc.AllRelig += count
				thishwc.Homil = count
			case "hymn":
				thishwc.Hymn = count
			case "hypoth":
				thishwc.Hypoth = count
			case "iamb":
				thishwc.Iamb = count
			case "ignotum":
				thishwc.Ignot = count
			case "invectiv":
				thishwc.Invectiv = count
			case "inscr":
				thishwc.Inscr = count
			case "jurisprud":
				thishwc.Juris = count
			case "lexicogr":
				thishwc.Lexic = count
			case "liturg":
				thishwc.AllRelig += count
				thishwc.Liturg = count
			case "lyr":
				thishwc.Lyr = count
			case "magica":
				thishwc.Magica = count
			case "math":
				thishwc.Math = count
			case "mech":
				thishwc.Mech = count
			case "med":
				thishwc.Med = count
			case "metrolog":
				thishwc.Meteor = count
			case "mim":
				thishwc.Mim = count
			case "mus":
				thishwc.Mus = count
			case "myth":
				thishwc.Myth = count
			case "narrfict":
				thishwc.NarrFic = count
			case "nathist":
				thishwc.NatHis = count
			case "onir":
				thishwc.Onir = count
			case "orac":
				thishwc.Orac = count
			case "orat":
				thishwc.AllRhet += count
				thishwc.Orat = count
			case "papyrus":
				thishwc.Papyrus = count
			case "paradox":
				thishwc.Paradox = count
			case "parod":
				thishwc.Parod = count
			case "paroem":
				thishwc.Paroem = count
			case "perieg":
				thishwc.Perig = count
			case "phil":
				thishwc.Phil = count
			case "physiognom":
				thishwc.Physiog = count
			case "poem":
				thishwc.Poem = count
			case "polyhist":
				thishwc.Polyhist = count
			case "prophet":
				thishwc.AllRelig += count
				thishwc.Proph = count
			case "pseudepigr":
				thishwc.Pseud = count
			case "rhet":
				thishwc.AllRhet += count
				thishwc.Rhet = count
			case "satura":
				thishwc.Satura = count
			case "satyr":
				thishwc.Satyr = count
			case "schol":
				thishwc.Schol = count
			case "tact":
				thishwc.Tact = count
			case "test":
				thishwc.Test = count
			case "theol":
				thishwc.AllRelig += count
				thishwc.Theol = count
			case "trag":
				thishwc.Trag = count
			default:
				fmt.Println("Unknown genre")
			}
		}
		allwords[word] = thishwc
	}
	return allwords
}
