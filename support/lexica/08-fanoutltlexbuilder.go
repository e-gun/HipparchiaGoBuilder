package lexica

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"strings"
	"sync"
	"time"
)

func FanoutLatinLexBuilder(dir string, fn string) {
	const (
		MSG = "FanoutGkLexBuilder(): All jobs processed. %.3fs"
	)

	start := time.Now()

	data := LoadLexFile(dir, fn)

	splitdata := strings.Split(data, "\n")
	bundles := generic.SplitIntoNSlices(splitdata, global.Config.WorkerCount)

	// Create a channel to send jobs to workers
	bc := make(chan []string, len(bundles))

	// Use WaitGroup to wait for all workers to complete their tasks
	var wg sync.WaitGroup

	// Start runtime.NumCPU() workers
	for i := 0; i <= global.Config.WorkerCount; i++ {
		// for i := 1; i <= 1; i++ {
		wg.Add(1)
		go latlexworker(i, bc, &wg)
	}

	// Send jobs to the channel
	for _, a := range bundles {
		bc <- a
	}

	// Close the channel since no more jobs will be sent
	close(bc)

	// Wait for all workers to finish processing
	wg.Wait()

	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func latlexworker(id int, ac <-chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for a := range ac {
		// fmt.Printf("%d\t%s\n", id, a)
		ParseAndLoadLatinLex(id, a)
	}
}
