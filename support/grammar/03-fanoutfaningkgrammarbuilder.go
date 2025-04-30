package grammar

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"sync"
)

// gkanalysisworker processes an individual value
func gkanalysisworker(id int, entries <-chan []string, results chan<- []structs.GramAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	result := ParseGreekAnalyses(<-entries)
	results <- result
}

func FanoutGreekAnalysis(entries []string) []structs.GramAnalysis {
	return FanoutAnalysis(entries, gkanalysisworker)
}

// gklemmataworker processes an individual value
func gklemmataworker(id int, entries <-chan []string, results chan<- []structs.HeadwordAndForms, wg *sync.WaitGroup) {
	defer wg.Done()
	result := ParseGreekLemmata(<-entries)
	results <- result
}

func FanoutGreekLemmata(entries []string) []structs.HeadwordAndForms {
	return FanoutLemmata(entries, gklemmataworker)
}
