//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	manydods = regexp.MustCompile(`!\d+`)
)

func GreekExclamation(ttc string) string {
	// 	exclamation point not properly documented
	//	note that runs of individual '!' can be found: '∙∙∙'
	//	but if you see '!21`45' you are supposed to print
	//	21 of them and then '45'
	return manydods.ReplaceAllStringFunc(ttc, multipledots)
}

func multipledots(match string) string {
	count := strings.TrimPrefix(match, "!")
	val, err := strconv.Atoi(count)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf("multipledots() failed with %s", count)
	}
	var dots []string
	for i := 0; i < val; i++ {
		dots = append(dots, "∙")
	}
	return strings.Join(dots, " ")
}
