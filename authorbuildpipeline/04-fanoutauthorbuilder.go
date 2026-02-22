//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

func FanoutBuilder(auu []string, datadir string) {
	const (
		MSG = "FanoutBuilder() All jobs processed. %.3fs"
	)

	start := time.Now()

	auu = SortWorkpileByFilesize(auu, datadir)

	// TESTING...
	//fmt.Println("TESTING; EDIT FanoutBuilder")
	//auu = auu[250:275]

	// Create a channel to send jobs to workers
	ac := make(chan string, len(auu))

	// Use WaitGroup to wait for all workers to complete their tasks
	var wg sync.WaitGroup

	// INS, etc call UniqueAuthorNames; the only way to make those names always map to the same people is to
	// cut the core count to 1...

	// after N cores you hit diminishing returns and that N is far less than the 20 available on the devel machine...
	// for example, there is *little* difference between 14 and 20: 322.296s vs 346.133s
	// (I/O saturation + the overhead of N pgsql clients for N workers)

	workercount := global.Config.WorkerCount
	if global.Config.Reproducible {
		workercount = 1
	}

	for i := 1; i <= workercount; i++ {
		wg.Add(1)
		go worker(i, datadir, ac, &wg)
	}

	// Send jobs to the channel
	for _, a := range auu {
		ac <- a
	}

	// Close the channel since no more jobs will be sent
	close(ac)

	// Wait for all workers to finish processing
	wg.Wait()

	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	global.DONE(d)
}

func worker(id int, datadir string, ac <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for a := range ac {
		BuildOneAuthor(datadir, a)
	}
}

func SortWorkpileByFilesize(auu []string, datadir string) []string {
	// this cuts 20s from the greek corpus build time and drops it to 107s
	fmt.Println("sorting work pile by length: longest authors go first")
	var ps stringintpairs
	for _, a := range auu {
		size := 0
		fi, err := os.Stat(datadir + a + ".TXT")
		if err != nil {
			// do nothing as size is already 0
		} else {
			size = int(fi.Size())
		}
		np := sipair{a, size}
		ps.Pairs = append(ps.Pairs, np)
	}
	sort.Sort(ps)
	for i, p := range ps.Pairs {
		auu[i] = p.Str
	}
	slices.Reverse(auu)
	return auu
}

type sipair struct {
	Str string
	Int int
}
type stringintpairs struct {
	Pairs []sipair
}

// Implement the sort.Interface for stringintpairs

func (p stringintpairs) Len() int           { return len(p.Pairs) }
func (p stringintpairs) Less(i, j int) bool { return p.Pairs[i].Int < p.Pairs[j].Int }
func (p stringintpairs) Swap(i, j int)      { p.Pairs[i], p.Pairs[j] = p.Pairs[j], p.Pairs[i] }
