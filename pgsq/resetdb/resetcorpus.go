package resetdb

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

func ResetCorpus(corp string) {
	gv := ""
	switch corp {
	case "LAT":
		gv = global.LATABBREV
	case "TLG":
		gv = global.TLGABBREV
	case "CHR":
		gv = global.CHRABBREV
	case "DDP":
		gv = global.DDPABBREV
	case "INS":
		gv = global.INSABBREV
	default:
		fmt.Println("Unknown corpus")
		gv = "find_nothing"

	}
	global.SECT("Resetting " + corp + " Corpus")
	DropAuthorMultipleTables(gv)
	ResetAuthorsAndWorks(gv)
}
