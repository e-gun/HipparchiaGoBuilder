//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorprep

import (
	"fmt"
	"os"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func ViewWorklines(lastfunc string, all bool, halt bool, wll []structs.DbWorkline) []structs.DbWorkline {
	fmt.Println("ViewWorklines() for", lastfunc)
	stop := STOPAFTER
	if all {
		stop = len(wll)
	}

	start := 0
	for _, v := range wll {
		fmt.Println(v.TbIndex, v.WkUID, v.FindLocus(), "\tmu\t", v.MarkedUp)
		start++
		if start >= stop {
			break
		}
	}

	fmt.Println("ViewWorklines() for", lastfunc)
	if halt {
		os.Exit(1)
	}

	return wll
}

func ViewWLAnnotations(lastfunc string, all bool, halt bool, wll []structs.DbWorkline) []structs.DbWorkline {
	fmt.Println("ViewWLAnnotations() for", lastfunc)
	stop := STOPAFTER
	if all {
		stop = len(wll)
	}

	start := 0
	for _, v := range wll {
		if v.Annotations != "" {
			fmt.Println(v.TbIndex, v.WkUID, v.FindLocus(), "\tan\t", v.Annotations)
			start++
		}
		if start >= stop {
			break
		}
	}

	fmt.Println("ViewWLAnnotations() for", lastfunc)
	if halt {
		os.Exit(1)
	}

	return wll
}

func WriteWorklineProgress(wll []structs.DbWorkline) {
	const (
		TMPL = "[%d] %s %s \t %s \t Notes: %s"
	)

	var linesout []string
	for _, v := range wll {
		linesout = append(linesout, fmt.Sprintf(TMPL, v.TbIndex, v.WkUID, v.FindLocus(), v.MarkedUp, v.Annotations))
	}

	// Write the string to the file (creates or overwrites the file)
	err := os.WriteFile("hgb-debug-output.txt", []byte(strings.Join(linesout, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File written successfully!")
	// os.Exit(0)
}

func WriteWorklineAnnotationsProgress(wll []structs.DbWorkline) {
	const (
		TMPL = "[%d] %s %s \t %s"
	)

	var linesout []string
	for _, v := range wll {
		if v.Annotations != "" {
			linesout = append(linesout, fmt.Sprintf(TMPL, v.TbIndex, v.WkUID, v.FindLocus(), v.Annotations))
		}
	}

	// Write the string to the file (creates or overwrites the file)
	err := os.WriteFile("hgb-debug-annotations-output.txt", []byte(strings.Join(linesout, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File written successfully!")
	//os.Exit(0)
}
