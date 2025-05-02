//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import (
	"strings"
	"testing"
)

func TestParseGreekLemmata(t *testing.T) {
	entries := strings.Split(GKLEMTESTDATA, "\n")
	ll := ParseGreekLemmata(entries)
	l1 := ll[0].DictEntry
	l2 := ll[0].DerivFormsSlc[3]
	w1 := "ἅβρα"
	w2 := "ἅβραιϲ"
	if l1 != w1 && l2 != w2 {
		t.Errorf("ParseGreekLemmata()\ngot\n%s\nwant\n%s", l1+" "+l2, w1+" "+w2)
	}
}

func TestParseLatinLemmata(t *testing.T) {
	entries := strings.Split(LATLEMTESTDATA, "\n")
	ll := ParseLatinLemmata(entries)
	l1 := ll[0].DictEntry
	l2 := ll[0].DerivFormsSlc[3]
	w1 := "Aaron"
	w2 := "Aaroni"
	if l1 != w1 && l2 != w2 {
		t.Errorf("ParseLatinLemmata()\ngot\n%s\nwant\n%s", l1+" "+l2, w1+" "+w2)
	}
}

func TestParseAnalyses(t *testing.T) {
	ParseGreekAnalyses([]string{})
}

func TestParseLatinAnalyses(t *testing.T) {
	entries := strings.Split(LTD2, "\n")
	ann := ParseLatinAnalyses(entries)
	for _, a := range ann {
		//fmt.Println(a.Observed)
		//fmt.Println(a.XREFS)
		//fmt.Println("# of poss:", len(a.Possibilities))
		//fmt.Println(a.Possibilities)
		a.PrintOut()
	}
}
