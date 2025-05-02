//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package wordcounts

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"time"
)

// ParsedCountOfAllCorpora - count the N instances of parsed word W in corpus C ...
func ParsedCountOfAllCorpora() map[string]structs.DbHeadwordCounts {
	const (
		MSG = "ParsedCountOfAllCorpora(): All corpora processed. %.3fs"
	)
	global.SECT("ParsedCountOfAllCorpora()")
	start := time.Now()

	doparsing := true

	tlg := CountTLGWords(doparsing)
	lat := CountLATWords(doparsing)
	ins := CountINSWords(doparsing)
	ddp := CountDDPWords(doparsing)
	chr := CountCHRWords(doparsing)
	mm := MapMerger([]map[string]int{tlg, lat, ins, ddp, chr})

	ccounts := make(map[string]structs.DbHeadwordCounts, len(mm))
	for k, v := range mm {
		ccounts[k] = structs.DbHeadwordCounts{
			Word:  k,
			Total: v,
		}
	}
	for k, v := range tlg {
		// every k in tlg should already have an entry in ccounts since mm was comprehensive
		wc := ccounts[k]
		wc.TLG = v
		ccounts[k] = wc
	}
	for k, v := range lat {
		wc := ccounts[k]
		wc.LAT = v
		ccounts[k] = wc
	}
	for k, v := range ins {
		wc := ccounts[k]
		wc.INS = v
		ccounts[k] = wc
	}
	for k, v := range ddp {
		wc := ccounts[k]
		wc.DDP = v
		ccounts[k] = wc
	}
	for k, v := range chr {
		wc := ccounts[k]
		wc.CHR = v
		ccounts[k] = wc
	}

	d := fmt.Sprintf(MSG, time.Now().Sub(start).Seconds())
	fmt.Println(d)
	return ccounts
}
