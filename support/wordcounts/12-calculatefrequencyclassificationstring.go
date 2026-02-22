//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import "regexp"

type tmpstruct struct {
	I int
	S string
}

var (
	islatin = regexp.MustCompile("^[a-z]")
	gfrq    = []tmpstruct{
		{56841, "top 250"},
		{4712, "top 2500"},
		{101, "core"},
		{5, "rare"},
		{0, "very rare"},
	}
	lfrq = []tmpstruct{
		{10305, "top 250"},
		{978, "top 2500"},
		{101, "core"},
		{5, "rare"},
		{0, "very rare"},
	}
)

func CalculateFrequencyClassification(word string, total int) string {
	// see HGS "helpdictionaries.html" on the cutoffs and for the sql to get the constants

	mapper := gfrq
	if islatin.MatchString(word) {
		mapper = lfrq
	}

	for _, m := range mapper {
		if total > m.I {
			return m.S
		}
	}

	return "(unknown)"

	//250th most common greek headword: ἀρτάβη (56,842)
	//250th most common latin headword: prima (10,306)
	//
	//2500th most common greek headword: προλέγω² (4,712)
	//2500th most common latin headword: aedo¹ (977)

	// GREEK
	//core (not in top 2500; >100 occurrences)
	//    count	19,211
	//    high	4,711
	//    low	    101
	//    average	770.45
	//    median	370

	//rare (between 100 and 6 occurrences)
	//    count	34,830
	//    high	100
	//    low	    5
	//    average	28.70
	//    median	19
	//
	//very rare (5 or fewer occurrences)
	//    count	30,933
	//    high	5
	//    low	    1
	//    average	2.40
	//    median	2

	// LATIN
	// top 250 latin headwords
	//    count	250
	//    high	424,851 // qui²
	//    low	    10,306
	//    average	37,152.97
	//    median	19,524
	//
	//top 2500 latin headwords (less the top 250)
	//    count	2250
	//    high	10,302
	//    low	    978
	//    average	2,877.87
	//    median	2,087
	//
	//core (not in top 2500; >100 occurrences)
	//    count	6049
	//    high	978
	//    low	    101
	//    average	344.86
	//    median	265
	//
	//rare (between 100 and 6 occurrences)
	//    count	11,694
	//    high	100
	//    low	    6
	//    average	30.12
	//    median	20
	//
	//very rare (fewer than 5 occurrences)
	//    count	9025
	//    high	5
	//    low	    1
	//    average	2.44
	//    median	2
}
