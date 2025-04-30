//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
	"strings"
)

const (
	GENRESFILE   = "LIST4CLA.BIN"
	GENRESCLXF   = "LIST4CLX.BIN"
	RECDATEFILE  = "LIST4DAT.BIN"
	EPITHETFILE  = "LIST4EPI.BIN"
	LOCATIONFILE = "LIST4GEO.BIN"
	GENDERFILE   = "LIST4FEM.BIN"
)

func GatherTLGAuthorMetadata() map[string]map[string][]string {
	gnn := BINDataExtraction(BINFileLoad(global.Config.GreekDir + GENRESFILE))
	gnx := BINDataExtraction(BINFileLoad(global.Config.GreekDir + GENRESCLXF))
	dtg := BINDataExtraction(BINFileLoad(global.Config.GreekDir + RECDATEFILE))
	ept := BINDataExtraction(BINFileLoad(global.Config.GreekDir + EPITHETFILE))
	loc := BINDataExtraction(BINFileLoad(global.Config.GreekDir + LOCATIONFILE))

	// looks like: 	Tactici: [0058 0546 0556 0648 3075 3181]
	// the locations have a lot of "Tarsos [vel Tarsus]" entries; this is the easiest moment to clean that
	cleaner := regexp.MustCompile(`(.*?) \[vel .*?]`)
	clnloc := make(map[string][]string)
	for k, v := range loc {
		ck := cleaner.ReplaceAllString(k, "$1")
		clnloc[ck] = v
	}

	// the epithets have a lot of "Historici／–ae" entries; this is the easiest moment to clean that
	clnept := make(map[string][]string)
	for k, v := range ept {
		ck := strings.ReplaceAll(k, "／–ae", "")
		clnept[ck] = v
	}

	// allmetadata := make(map[string]map[string][]string)
	allmetadata := map[string]map[string][]string{
		"gnn": gnn,
		"gnx": gnx,
		"dtg": dtg,
		"ept": clnept,
		"loc": clnloc,
	}

	for k, v := range allmetadata {
		allmetadata[k] = renamemapkeys(invertmap(v))
	}

	//for k, v := range allmetadata["gnn"] {
	//	fmt.Println(k, v)
	//}

	//gr9985 [Theologica  Comica  Astrologica]
	//gr1833 [Lyrica ]
	//gr2583 [Astrologica Hexametrica]
	//gr0333 [Tragica Satyrae]

	return allmetadata
}

// invertmap - turn entries like {Poetae Medici : [0280 0281 0676 0888]} into { 2080 : [Poetae Medici], 0281 : [Poetae Medici] ...}
func invertmap(ssl map[string][]string) map[string][]string {
	inverted := make(map[string][]string)

	for k, stringslice := range ssl {
		for _, s := range stringslice {
			if _, ok := inverted[s]; !ok {
				inverted[s] = []string{k}
			} else {
				inverted[s] = append(inverted[s], k)
			}
		}
	}

	// cleanup
	for k, v := range inverted {
		inverted[k] = generic.TrimAllStrings(generic.DropEmptyStrings(generic.Unique(v)))
	}
	return inverted
}

// renamemapkeys - turn 2080 into gr2080
func renamemapkeys(ssl map[string][]string) map[string][]string {
	const (
		CORPUSPREFIX = "gr"
	)
	renamed := make(map[string][]string)
	for k, v := range ssl {
		renamed[CORPUSPREFIX+k] = v
	}
	return renamed
}
