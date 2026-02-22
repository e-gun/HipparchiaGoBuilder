//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import (
	"sync"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

// gkanalysisworker processes an individual value
func ltanalysisworker(id int, entries <-chan []string, results chan<- []structs.GramAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	result := ParseLatinAnalyses(<-entries)
	results <- result
}

func FanoutLatinAnalysis(entries []string) []structs.GramAnalysis {
	return FanoutAnalysis(entries, ltanalysisworker)
}

// latlemmataworker processes an individual value
func latlemmataworker(id int, entries <-chan []string, results chan<- []structs.HeadwordAndForms, wg *sync.WaitGroup) {
	defer wg.Done()
	result := ParseLatinLemmata(<-entries)
	results <- result
}

func FanoutLatinLemmata(entries []string) []structs.HeadwordAndForms {
	return FanoutLemmata(entries, latlemmataworker)
}
