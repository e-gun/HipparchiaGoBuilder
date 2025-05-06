package wordcounts

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"time"
)

// CountGenreRawWords - in order to do the weights, you need raw counts
func CountGenreRawWords() {
	// weighting calculations should be based off of the raw word count and not parsed word counts
	// this is *not* headword data; but it is being stored in a 'headword' table because only
	// this table knows about genres and eras
	const (
		MSG = "CountGenreRawWords(): All genres processed. %.3fs"
	)
	global.SECT("CountGenreRawWords()")
	start := time.Now()
	for i, genre := range global.KnownGenres {
		if i%20 == 0 {
			fmt.Printf("\tCountGenreRawWords() working on %d of %d genres\n", i, len(global.KnownGenres))
		}
		wbhh := getrelevantgenreworks(genre)
		doparsing := false
		counts := FanoutGenreAndTimeCounter(wbhh, doparsing)
		genretotalcount := 0
		// now we just grab the total
		for _, v := range counts {
			genretotalcount = genretotalcount + v
		}
		insert.InsertOneRawCountIntoHeadwordWordcounts(genre, genretotalcount)
	}
	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func CountEraRawWords() {
	const (
		MSG = "CountEraRawWords(): All eras processed. %.3fs"
	)
	global.SECT("CountEraRawWords()")
	start := time.Now()
	for era, _ := range eras {
		fmt.Printf("\tEra wordcount working on '%s'\n", era)
		wbhh := getrelevanttimespanworks(era)
		doparsing := false
		counts := FanoutGenreAndTimeCounter(wbhh, doparsing)
		eratotalcount := 0
		// now we just grab the total
		for _, v := range counts {
			eratotalcount = eratotalcount + v
		}
		insert.InsertOneRawCountIntoHeadwordWordcounts(era, eratotalcount)
	}
	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}
