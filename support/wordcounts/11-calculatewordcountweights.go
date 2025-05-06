//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import (
	"context"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"github.com/jackc/pgx/v5/pgxpool"
)

// there is going to be distortion here because of the GRK/LAT issue

// that is, the weight of Lat genres should be different, but the counts are not sensitive to that
// should only sum `WHERE entry_name ~* '[α-ωϲ]'`
// note that `[^a-z]` will catch `adeo²` because of the ² and so it is not a "not greek" pattern

// Satura quidem tota nostra est...

// hgdb=> SELECT sum (satura) FROM headword_wordcounts where entry_name ~* '[α-ωϲ]';
//  sum
//-------
// 50604

//
//hgdb=> SELECT sum (satura) FROM headword_wordcounts where entry_name ~* '[^α-ωϲ]';
//  sum
//--------
// 303616

// CalculateParsedWordcountTotals - calculate and load weighting for genres and TheEras: 1 word of Trag is equiv to X words of Phil...
// the logic for that hinges on the use of mps.ParsedGreekWeightsCorpora and str.WeightedFieldValuePair in HGS
// the high count value is set to "1" and all others are "1/X" to give a multiplier that is applied to the value you see
// so if there is 10x more Phil than Epic, 3 instances of a word in Epic is worth 30 in Phil and so if there are 3 in
// Epic and 10 in Phil, the word will be scored as "much more typical of epic" since 30 > 10.
// NB: the *parsed* headwords are weighed against *unparsed* totals, and so you have a fiction because the total number
// of parsed words is greater than the total number of unparsed words owing to homonyms. This might need to be reverted,
// but it is not clear that the 'statistics' mean either more or less depending on which route you use to derive this
// number; note also that despite the above the `total_count`, `gr_count`, ... represent *parsed* headwords; genres and
// eras are *unparsed* totals; comment out the relevant lines of DoAllWordcounts() to go to pure parsed counts
func CalculateParsedWordcountTotals() {
	const (
		EN1 = `__wordcounttotals`
		EN2 = GREEKRAWCOUNTROW
		EN3 = LATINRAWCOUNTROW
		Q1  = `DELETE FROM headword_wordcounts WHERE entry_name = $1;`
		Q2  = `SELECT sum (%s) FROM headword_wordcounts;`
		Q3  = `SELECT sum (%s) FROM headword_wordcounts where entry_name ~* '[` + generic.GLC + `]'`
		Q4  = `SELECT sum (%s) FROM headword_wordcounts where entry_name ~* '[^` + generic.GLC + `]'`
	)

	global.SECT("CalculateParsedWordcountTotals()")

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	// reset the old data
	for _, tot := range []string{EN1, EN2, EN3} {
		_, err := dbconn.Exec(context.Background(), Q1, tot)
		if err != nil {
			fmt.Println(err)
		}
	}

	// get the new data for __wordcounttotals
	allsums1 := fillsums(Q2, dbconn)

	// get the new data for __greekwordcounttotals
	allsums2 := fillsums(Q3, dbconn)

	// get the new data for __latinwordcounttotals
	allsums3 := fillsums(Q4, dbconn)

	// can only insert at end: otherwise you will count the prior counts
	insertone(EN1, allsums1)
	insertone(EN2, allsums2)
	insertone(EN3, allsums3)
}

func fillsums(q string, dbconn *pgxpool.Conn) []int {
	var (
		skip = map[string]bool{
			"entry_name":               true,
			"frequency_classification": true,
		}
	)
	var toget []string
	for _, row := range insert.HeadwordWordCountsRows {
		_, shouldskip := skip[row]
		if !shouldskip {
			toget = append(toget, row)
		}
	}
	var allsums []int
	for _, get := range toget {
		var thissum int
		s := dbconn.QueryRow(context.Background(), fmt.Sprintf(q, get))
		e := s.Scan(&thissum)
		if e != nil {
			fmt.Println(e)
		}
		//m := message.NewPrinter(language.English)
		//fmt.Println(get, m.Sprintf("%d", thissum))
		allsums = append(allsums, thissum)
	}
	return allsums
}

func insertone(category string, sums []int) {
	wct := structs.DbHeadwordCounts{
		Word:     category,
		Total:    sums[0],
		TLG:      sums[1],
		LAT:      sums[2],
		DDP:      sums[3],
		INS:      sums[4],
		CHR:      sums[5],
		FrqClas:  "",
		Early:    sums[6],
		Middle:   sums[7],
		Late:     sums[8],
		Acta:     sums[9],
		Agric:    sums[10],
		Alchem:   sums[11],
		Anthol:   sums[12],
		Apocal:   sums[13],
		Apocry:   sums[14],
		Apol:     sums[15],
		Astrol:   sums[16],
		Astron:   sums[17],
		Biogr:    sums[18],
		Bucol:    sums[19],
		Caten:    sums[20],
		Chron:    sums[21],
		Comic:    sums[22],
		Comm:     sums[23],
		Concil:   sums[24],
		Coq:      sums[25],
		Dial:     sums[26],
		Docu:     sums[27],
		Doxog:    sums[28],
		Eccl:     sums[29],
		Eleg:     sums[30],
		Encom:    sums[31],
		Epic:     sums[32],
		Epigr:    sums[33],
		Epist:    sums[34],
		Evang:    sums[35],
		Exeg:     sums[36],
		Fab:      sums[37],
		Geog:     sums[38],
		Gnom:     sums[39],
		Gram:     sums[40],
		Hagiog:   sums[41],
		Hexam:    sums[42],
		Hist:     sums[43],
		Homil:    sums[44],
		Hymn:     sums[45],
		Hypoth:   sums[46],
		Iamb:     sums[47],
		Ignot:    sums[48],
		Inscr:    sums[49],
		Invectiv: sums[50],
		Liturg:   sums[51],
		Juris:    sums[52],
		Lexic:    sums[53],
		Lyr:      sums[54],
		Magica:   sums[55],
		Math:     sums[56],
		Mech:     sums[57],
		Med:      sums[58],
		Meteor:   sums[59],
		Mim:      sums[60],
		Mus:      sums[61],
		Myth:     sums[62],
		NarrFic:  sums[63],
		NatHis:   sums[64],
		Onir:     sums[65],
		Orac:     sums[66],
		Orat:     sums[67],
		Papyrus:  sums[68],
		Paradox:  sums[69],
		Parod:    sums[70],
		Paroem:   sums[71],
		Perig:    sums[72],
		Phil:     sums[73],
		Physiog:  sums[74],
		Poem:     sums[75],
		Polyhist: sums[76],
		Proph:    sums[77],
		Pseud:    sums[78],
		Rhet:     sums[79],
		Satura:   sums[80],
		Satyr:    sums[81],
		Schol:    sums[82],
		Tact:     sums[83],
		Test:     sums[84],
		Theol:    sums[85],
		Trag:     sums[86],
		AllRhet:  sums[87],
		AllRelig: sums[88],
	}

	//fmt.Println(category)
	//wct.PrintOut()

	//DbHeadwordCounts: __wordcounttotals
	//Epic: 1,782,790
	//Phil: 16,093,727
	//Trag: 457,516
	//AllRhet: 9,565,004
	//AllRelig: 19,881,637
	//Total: 137,774,681

	//DbHeadwordCounts: __greekwordcounttotals
	//Epic: 755,782
	//Phil: 14,549,749
	//Trag: 288,359
	//AllRhet: 7,104,369
	//AllRelig: 19,865,609
	//Total: 119,386,701
	//InsertHeadwordWordcounts inserted 1 rows

	// __latinwordcounttotals
	//DbHeadwordCounts: __latinwordcounttotals
	//Epic: 1,074,868
	//Phil: 4,209,536
	//Trag: 201,647
	//AllRhet: 3,642,682
	//AllRelig: 3,237,452
	//Total: 37,421,670

	insertme := map[string]structs.DbHeadwordCounts{
		wct.Word: wct,
	}
	insert.InsertHeadwordWordcountsIntoTable(insertme)
}

// hgdb=> select universalid,title,workgenre,wordcount from works where workgenre ~* 'Trag(|.|;)' and universalid ~* '^gr' order by wordcount desc;
// universalid |                            title                            |                   workgenre                   | wordcount
//-------------+-------------------------------------------------------------+-----------------------------------------------+-----------
// gr0085w008  | Fragmenta                                                   | Satyr.; Trag.                                 |     34166
// gr2022w003  | Christus patiens [Dub.] (fort. auctore Constantino Manasse) | Trag.                                         |     15934
// gr0059w009  | Parmenides                                                  | Phil.; Epist.; Epigr.; Trag.; Dialog.         |     15261
// gr0011w008  | Fragmenta                                                   | Satyr.; Trag.                                 |     12629
// gr1738w003  | Fragmenta                                                   | Satyr.; Trag.                                 |     11105
// gr0011w004  | Oedipus tyrannus                                            | Trag.                                         |      9771
// gr0011w006  | Philoctetes                                                 | Trag.                                         |      9264
// gr0011w005  | Electra                                                     | Trag.                                         |      9142
// gr0085w005  | Agamemnon                                                   | Trag.                                         |      8491
// gr0011w003  | Ajax                                                        | Trag.                                         |      8234
// gr0011w002  | Antigone                                                    | Trag.                                         |      7684
// gr0011w001  | Trachiniae                                                  | Trag.                                         |      7570
// gr0341w002  | Alexandra                                                   | Trag.                                         |      7462
// gr0085w003  | Prometheus vinctus                                          | Trag.                                         |      6158
// gr0085w006  | Choephoroe                                                  | Trag.                                         |      5745
// ...

// hgdb=> select universalid,title,workgenre,wordcount from works where workgenre ~* 'Trag.' and universalid ~* '^gr' order by wordcount desc;
//
// hgdb=> SELECT sum (wordcount) FROM works where workgenre ~* 'Trag(|.|;)' and universalid ~* '^gr';
//  sum
//--------
// 205270
//(1 row)

// The homonymns issue lets the 'works' count be smaller than the headwords count. But this is a big gap...

// hgdb=> SELECT sum (wordcount) FROM works where workgenre ~* 'Phil.' and universalid ~* '^gr';
//   sum
//---------
// 9091119
//(1 row)
//
//hgdb=> SELECT sum (wordcount) FROM works where workgenre ~* 'Phil.' and universalid ~* '^lt';
//  sum
//--------
// 604747
//(1 row)
