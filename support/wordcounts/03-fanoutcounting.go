//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"sync"
)

func FanoutCorpusCount(doparsing bool, authorpile []string) map[string]int {
	// Input slice of bundles to be processed
	// authorpile := FindAuthorsToGrab(corpus)

	// Channels for fan-out (tasks) and fan-in (results)
	tasks := make(chan string, len(authorpile))
	results := make(chan map[string]int, global.Config.WorkerCount)

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	var myworkertype func(id int, tables <-chan string, results chan<- map[string]int, wg *sync.WaitGroup)
	if doparsing {
		myworkertype = parseandcountworker
	} else {
		myworkertype = countworker
	}

	// Start the workers
	for i := 0; i <= global.Config.WorkerCount; i++ {
		wg.Add(1)
		go myworkertype(i, tasks, results, &wg)
	}

	// Fan-out: Distribute tasks to workers
	go func() {
		for _, value := range authorpile {
			tasks <- value
		}
		close(tasks)
	}()

	// Wait for all workers to finish and collect results
	var aggregateresults []map[string]int
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan-in: Collect the processed results from workers
	for result := range results {
		aggregateresults = append(aggregateresults, result)
	}

	// these aggregateresults still need merging into a single map[int]
	mastermap := MapMerger(aggregateresults)

	return mastermap
}

func FanoutGenreAndTimeCounter(worksandbounds []WorksAndBoundsHolder, doparsing bool) map[string]int {
	// Input slice of bundles to be processed
	// authorpile := FindAuthorsToGrab(corpus)

	// Channels for fan-out (tasks) and fan-in (results)
	tasks := make(chan WorksAndBoundsHolder, len(worksandbounds))
	results := make(chan map[string]int, global.Config.WorkerCount)

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start the workers
	for i := 0; i <= global.Config.WorkerCount; i++ {
		wg.Add(1)
		if doparsing {
			go parseandcountsubetworker(i, tasks, results, &wg)
		} else {
			go countsubetworker(i, tasks, results, &wg)
		}
	}

	// Fan-out: Distribute tasks to workers
	go func() {
		for _, wbh := range worksandbounds {
			tasks <- wbh
		}
		close(tasks)
	}()

	// Wait for all workers to finish and collect results
	var aggregateresults []map[string]int
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan-in: Collect the processed results from workers
	for result := range results {
		aggregateresults = append(aggregateresults, result)
	}

	// these aggregateresults still need merging into a single map[int]
	mastermap := MapMerger(aggregateresults)

	return mastermap
}

// MapMerger - []map[string]int -> map[string]int
func MapMerger(aggregateresults []map[string]int) map[string]int {
	mastermap := make(map[string]int)
	for _, oneresultmap := range aggregateresults {
		for k, v := range oneresultmap {
			if _, ok := mastermap[k]; !ok {
				mastermap[k] = v
			} else {
				mastermap[k] = mastermap[k] + v
			}
		}
	}
	return mastermap
}

// countworker processes an individual table; do not parse the words
func countworker(id int, tables <-chan string, results chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	localmap := make(map[string]int)
	for t := range tables {
		// fmt.Println(id, t)
		lines := GetAllLinesFrom(t)
		localmap = CountThisBundle(localmap, lines)
	}
	results <- localmap
}

// parseandcountworker processes an individual table; parse the words
func parseandcountworker(id int, tables <-chan string, results chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	localmap := make(map[string]int)
	for t := range tables {
		// fmt.Println(id, t)
		lines := GetAllLinesFrom(t)
		localmap = ParseAndCountThisBundle(localmap, lines)
	}
	results <- localmap
}

// parseandcountsubetworker processes a subset of an author table; parse the words
func parseandcountsubetworker(id int, worksandbounds <-chan WorksAndBoundsHolder, results chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	localmap := make(map[string]int)
	for wb := range worksandbounds {
		// fmt.Println(id, t)
		lines := GetSelectLinesFrom(wb.T, wb.F, wb.L)
		localmap = ParseAndCountThisBundle(localmap, lines)
	}
	results <- localmap
}

// countsubetworker processes a subset of an author table; do not parse the words
func countsubetworker(id int, worksandbounds <-chan WorksAndBoundsHolder, results chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	localmap := make(map[string]int)
	for wb := range worksandbounds {
		// fmt.Println(id, t)
		lines := GetSelectLinesFrom(wb.T, wb.F, wb.L)
		localmap = CountThisBundle(localmap, lines)
	}
	results <- localmap
}
