package lexica

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"sync"
	"time"
)

func FanoutGkLexBuilder(xmls []string, datadir string) {
	const (
		MSG = "FanoutGkLexBuilder(): All jobs processed. %.3fs"
	)

	start := time.Now()

	// Create a channel to send jobs to workers
	xc := make(chan string, len(xmls))

	// Use WaitGroup to wait for all workers to complete their tasks
	var wg sync.WaitGroup

	// Start runtime.NumCPU() workers
	for i := 0; i <= global.Config.WorkerCount; i++ {
		// for i := 1; i <= 1; i++ {
		wg.Add(1)
		go gklexworker(i, datadir, xc, &wg)
	}

	// Send jobs to the channel
	for _, a := range xmls {
		xc <- a
	}

	// Close the channel since no more jobs will be sent
	close(xc)

	// Wait for all workers to finish processing
	wg.Wait()

	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func gklexworker(id int, datadir string, ac <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for a := range ac {
		// fmt.Printf("%d\t%s\n", id, a)
		// 1	greatscott01.xml
		// 15	greatscott20.xml
		// 10	greatscott11.xml
		// ...
		ParseAndLoadLSJ(datadir, a)
	}
}
