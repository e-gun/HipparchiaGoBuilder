package wordcounts

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"regexp"
	"time"
)

const (
	GREEKRAWCOUNTROW = `__greekunparsedwordcounttotalsstoredamongheadwordcounts`
	LATINRAWCOUNTROW = `__latinunparsedwordcounttotalsstoredamongheadwordcounts`
)

// CountGenreRawWords - in order to do the weights, you need raw counts
func CountGenreRawWords() {
	// weighting calculations should be based off of the raw word count and not parsed word counts
	// this is *not* headword data; but it is being stored in a 'headword' table because only
	// this table knows about genres and TheEras

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
		lat := countlatin(counts)
		insert.InsertOneRawCountIntoHeadwordWordcounts(genre, LATINRAWCOUNTROW, lat)
		grk := countgreek(counts)
		insert.InsertOneRawCountIntoHeadwordWordcounts(genre, GREEKRAWCOUNTROW, grk)
	}
	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func CountEraRawWords() {
	const (
		MSG = "CountEraRawWords(): All TheEras processed. %.3fs"
	)
	global.SECT("CountEraRawWords()")
	start := time.Now()
	for era, _ := range TheErasAlt {
		fmt.Printf("\tEra wordcount working on '%s'\n", era)
		wbhh := getrelevanttimespanworks(era)
		doparsing := false
		counts := FanoutGenreAndTimeCounter(wbhh, doparsing)
		lat := countlatin(counts)
		insert.InsertOneRawCountIntoHeadwordWordcounts(era, LATINRAWCOUNTROW, lat)
		grk := countgreek(counts)
		insert.InsertOneRawCountIntoHeadwordWordcounts(era, GREEKRAWCOUNTROW, grk)
	}
	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func countgreek(counts map[string]int) int {
	isgreek := regexp.MustCompile("^[^a-z]")
	total := 0
	// now we just grab the total
	for k, v := range counts {
		if isgreek.MatchString(k) {
			total = total + v
		}
	}
	return total
}

func countlatin(counts map[string]int) int {
	islatin := regexp.MustCompile("^[a-z]")
	total := 0
	// now we just grab the total
	for k, v := range counts {
		if islatin.MatchString(k) {
			total = total + v
		}
	}
	return total
}
