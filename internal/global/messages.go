//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package global

import (
	"fmt"
	"runtime"
)

const (
	RED1    = "\033[38;5;160m" // Red3
	YELLOW1 = "\033[38;5;178m" // Gold3
	GREEN   = "\033[38;5;70m"  // Chartreuse3
	RESET   = "\033[0m"
	CTMPL   = "*** %s%s%s *** \n"
	CTMPL2  = "*** %s *** \n"
	VTEMPL  = "%s%s%s\n"
)

func MSG(s string) {
	const (
		TMPL = "\t%s\n"
	)
	if Config.QuietBuild {
		return
	}
	fmt.Printf(TMPL, s)
}

func HEAD(s string) {
	if runtime.GOOS == "windows" {
		fmt.Printf(CTMPL2, s)
	} else {
		fmt.Printf(CTMPL, YELLOW1, s, RESET)
	}
}

func SECT(s string) {
	if runtime.GOOS == "windows" {
		fmt.Println(s)
	} else {
		fmt.Printf(VTEMPL, GREEN, s, RESET)
	}
}

func BAD(s string) {
	if runtime.GOOS == "windows" {
		fmt.Printf(CTMPL2, s)
	} else {
		fmt.Printf(CTMPL, RED1, s, RESET)
	}
}
