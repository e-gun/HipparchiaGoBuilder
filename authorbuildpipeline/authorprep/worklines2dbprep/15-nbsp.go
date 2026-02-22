//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	bqsfinder = regexp.MustCompile(`<hb_blank_quarter_spaces quantity="(\d+)" />`)
)

const (
	TABSVIASPACES          = `&nbsp;&nbsp;&nbsp;`
	STANDALONE             = `<hb-tabbedtext />`
	THISCOUNTSASONEQUARTER = 5
)

// InsertNBSP - pseudo-markup re spaces into actual html for spaces
func InsertNBSP(lines []structs.DbWorkline) []structs.DbWorkline {
	//	this is a place where some commitments that have been deferred will get made
	//	obviously you could put all of this in earlier for the sake of efficiency,
	//	but handling this here and now is convenient: one-stop shopping w/out tripping up the parser earlier

	for i, line := range lines {
		newtext := strings.ReplaceAll(line.MarkedUp, STANDALONE, TABSVIASPACES)
		newtext = bqsfinder.ReplaceAllStringFunc(newtext, insertnspaces)
		lines[i].MarkedUp = newtext
	}
	return lines
}

func insertnspaces(match string) string {
	m := bqsfinder.ReplaceAllString(match, "$1")
	d, e := strconv.Atoi(m)
	if e != nil {
		fmt.Println("insertnspaces() could not convert string to int:", m[1])
		panic(e)
	}
	sp := d / THISCOUNTSASONEQUARTER

	var spacer string
	for i := 0; i < sp; i++ {
		spacer += "&nbsp;"
	}
	return spacer
}
