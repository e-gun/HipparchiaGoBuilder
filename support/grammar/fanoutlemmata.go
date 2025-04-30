package grammar

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"sync"
)

func FanoutLemmata(entries []string, workerfnc func(id int, entries <-chan []string, results chan<- []structs.HeadwordAndForms, wg *sync.WaitGroup)) []structs.HeadwordAndForms {
	// Input slice of bundles to be processed
	bundles := generic.SplitIntoNSlices(entries, global.Config.WorkerCount)

	// Channels for fan-out (tasks) and fan-in (results)
	tasks := make(chan []string, len(bundles))
	results := make(chan []structs.HeadwordAndForms, len(bundles))

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start the workers
	for i := 0; i <= global.Config.WorkerCount-1; i++ {
		wg.Add(1)
		go workerfnc(i, tasks, results, &wg)
	}

	// Fan-out: Distribute tasks to workers
	go func() {
		for _, value := range bundles {
			tasks <- value
		}
		close(tasks)
	}()

	// Wait for all workers to finish and collect results
	var aggregateresults []structs.HeadwordAndForms
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan-in: Collect the processed results from workers
	for result := range results {
		aggregateresults = append(aggregateresults, result...)
	}

	// Output the processed results
	return aggregateresults
}
