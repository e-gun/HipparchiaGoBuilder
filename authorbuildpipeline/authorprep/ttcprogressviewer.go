//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorprep

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	SPLITON = `█⑧⓪`
	// SPLITON   = `█⓪ █⓪ █⓪`
	SPLITALT  = `<hb-nb-endofpage />`
	STOPAFTER = 50
)

// ViewAndFilter - view the first N lines of a function's output and possibly halt execution here
func ViewAndFilter(folder string, lastfunc string, halt bool, filter string, ttc string) string {
	fmt.Println("START ViewAndFilter() for", folder, ":", lastfunc)
	ttc = vwandfilt(filter, ttc, SPLITON, 0)
	fmt.Println("END ViewAndFilter() for", folder, ":", lastfunc)

	if halt {
		os.Exit(1)
	}

	return ttc
}

func DipIntoResults(folder string, lastfunc string, startat int, halt bool, ttc string) string {
	fmt.Println("DipIntoResuls() for", folder, ":", lastfunc, "at", startat)
	ttc = vwandfilt("", ttc, SPLITON, startat)
	fmt.Println("DipIntoResuls() for", folder, ":", lastfunc, "at", startat)
	if halt {
		os.Exit(1)
	}
	return ttc
}

func WorkLineViewAndFilter(folder string, lastfunc string, halt bool, filter string, ttc string) string {
	fmt.Println("WorkLineViewAndFilter() for", folder, ":", lastfunc)
	ttc = vwandfilt(filter, ttc, "\n", 0)
	fmt.Println("WorkLineViewAndFilter() for", folder, ":", lastfunc)

	if halt {
		os.Exit(1)
	}

	return ttc
}

func vwandfilt(filter string, ttc string, spl string, skipto int) string {

	split := strings.Split(ttc, spl)

	count := 0
	for i, v := range split {
		if i < skipto {
			continue
		}
		if strings.Contains(v, filter) {
			count++
			if strings.Contains(v, "█⓪") {
				continue
			}
			fmt.Printf("\t" + split[i] + "\n")
		}
		if count > STOPAFTER {
			break
		}
	}

	return ttc
}

func WriteTTCProgress(content string) {

	content = strings.Replace(content, SPLITON, "\n", -1)

	// Write the string to the file (creates or overwrites the file)
	fname := fmt.Sprintf("hgb-debug-output-%v.txt", time.Now().UnixNano())
	// fname := "hgb-debug-output.txt"
	err := os.WriteFile(fname, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println(fname, "written successfully!")
	// os.Exit(0)
}
