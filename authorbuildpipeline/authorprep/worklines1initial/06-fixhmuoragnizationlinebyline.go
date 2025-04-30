//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	fontspanopener = regexp.MustCompile(`<hb-(sp|fs)-(.*?)>`)
	fontspancloser = regexp.MustCompile(`</hb-(sp|fs)-(.*?)>`)
)

func FixIrrationalTags(ttc string) string {
	lines := strings.Split(ttc, "\n")
	for i, l := range lines {
		lines[i] = fixhmuirrationaloragnization(l)
	}
	return strings.Join(lines, "\n")
}

type modification struct {
	Item     int
	Position int
	CloseTag string
	OpenTag  string
}

func fixhmuirrationaloragnization(workline string) string {
	// this is a line-by-line function, not a textblock function

	// 		Note the irrationality (for HTML) of the following (which is masked by the 'spanner'):
	//		[have 'EX_ON' + 'SM_ON' + 'EX_OFF' + 'SM_OFF']
	//		[need 'EX_ON' + 'SM_ON' + 'SM_OFF' + 'EX_OFF' + 'SM_ON' + 'SM_OFF' ]
	//
	//		hipparchiaDB=# SELECT index, marked_up_line FROM gr0085 where index = 14697;
	//		 index |                                                                           marked_up_line
	//		-------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	//		 14697 | <hmu_span_expanded_text><hmu_fontshift_greek_smallerthannormal>τίϲ ἡ τάραξιϲ</hmu_span_expanded_text> τοῦ βίου; τί βάρβιτοϲ</hmu_fontshift_greek_smallerthannormal>
	//		(1 row)
	//
	//
	//		hipparchiaDB=> SELECT index, marked_up_line FROM gr0085 where index = 14697;
	//		 index |                                                marked_up_line
	//		-------+---------------------------------------------------------------------------------------------------------------
	//		 14697 | <span class="expanded_text"><span class="smallerthannormal">τίϲ ἡ τάραξιϲ</span> τοῦ βίου; τί βάρβιτοϲ</span>
	//		(1 row)
	//
	//
	//		fixing this is an interesting question; it seems likely that I have missed some way of doing it wrong...
	//		but note 'b' below: this is pretty mangled and the output is roughly right...
	//
	//		invalidline = '<hmu_span_expanded_text><hmu_fontshift_greek_smallerthannormal>τίϲ ἡ τάραξιϲ</hmu_span_expanded_text> τοῦ βίου; τί βάρβιτοϲ</hmu_fontshift_greek_smallerthannormal>'
	//		openspans {0: 'span_expanded_text', 24: 'fontshift_greek_smallerthannormal'}
	//		closedspans {76: 'span_expanded_text', 123: 'fontshift_greek_smallerthannormal'}
	//		balancetest [(False, False, True)]
	//
	//		validline = '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<hmu_fontshift_latin_smallcapitals>errantes</hmu_fontshift_latin_smallcapitals><hmu_fontshift_latin_normal> pascentes, ut alibi “mille meae Siculis</hmu_fontshift_latin_normal>'
	//		openspans {36: 'fontshift_latin_smallcapitals', 115: 'fontshift_latin_normal'}
	//		closedspans {79: 'fontshift_latin_smallcapitals', 183: 'fontshift_latin_normal'}
	//		balancetest [(False, True, False)]
	//
	//		# need a third check: or not (open[okeys[x]] == closed[ckeys[x]])
	//		z = '&nbsp;&nbsp;&nbsp;<hmu_fontshift_latin_normal>II 47.</hmu_fontshift_latin_normal><hmu_fontshift_latin_italic> prognosticorum causas persecuti sunt et <hmu_span_latin_expanded_text>Boëthus Stoicus</hmu_span_latin_expanded_text>,</hmu_fontshift_latin_italic>'
	//		openspans {18: 'fontshift_latin_normal', 81: 'fontshift_latin_italic', 150: 'span_latin_expanded_text'}
	//		closedspans {52: 'fontshift_latin_normal', 195: 'span_latin_expanded_text', 227: 'fontshift_latin_italic'}
	//		balancetest [(False, True, False), (True, False, True)]
	//
	//		a = '[]κακ<hmu_span_superscript>η</hmu_span_superscript> βου<hmu_span_superscript>λ</hmu_span_superscript>'
	//		openspans {5: 'span_superscript', 55: 'span_superscript'}
	//		closedspans {28: 'span_superscript', 78: 'span_superscript'}
	//		balancetest [(False, True, False)]
	//
	//		b = []κακ<hmu_span_superscript>η</hmu_span_superscript> β<hmu_span_x>ο<hmu_span_y>υab</hmu_span_x>c<hmu_span_superscript>λ</hmu_span_y></hmu_span_superscript>
	//		testresult (False, True, False)
	//		testresult (False, False, True)
	//		testresult (False, False, True)
	//		balanced to:
	//			[]κακ<hmu_span_superscript>η</hmu_span_superscript> β<hmu_span_x>ο<hmu_span_y>υab</hmu_span_y></hmu_span_x><hmu_span_y>c<hmu_span_superscript>λ</hmu_span_superscript></hmu_span_y><hmu_span_superscript></hmu_span_superscript>

	openings := fontspanopener.FindAllStringSubmatchIndex(workline, -1)
	closings := fontspancloser.FindAllStringSubmatchIndex(workline, -1)

	openspans := make(map[int]string)
	for _, o := range openings {
		openspans[o[0]] = fmt.Sprintf("%s_%s", workline[o[2]:o[3]], workline[o[4]:o[5]])
	}

	closedspans := make(map[int]string)
	for _, c := range closings {
		closedspans[c[0]] = fmt.Sprintf("%s_%s", workline[c[2]:c[3]], workline[c[4]:c[5]])
	}

	balancetest := [][]bool{}
	invalidpattern := []bool{false, false, true}

	okeys := make([]int, 0, len(openspans))
	ckeys := make([]int, 0, len(closedspans))

	if len(openspans) == len(closedspans) && len(openspans) > 1 {
		// test 1: a problem if the next open ≠ this close and next open position comes before this close position
		//   	open: {0: 'span_expanded_text', 24: 'fontshift_greek_smallerthannormal'}
		// 		closed: {76: 'span_expanded_text', 123: 'fontshift_greek_smallerthannormal'}
		// test 2: succeed if the next open comes after the this close AND the this set of tags match
		//       open {18: 'fontshift_latin_normal', 81: 'fontshift_latin_italic', 150: 'span_latin_expanded_text'}
		// 		closed {52: 'fontshift_latin_normal', 195: 'span_latin_expanded_text', 227: 'fontshift_latin_italic'}
		// test 3: succeed if the next open comes before the previous close

		rng := len(openspans) - 1

		for k := range openspans {
			okeys = append(okeys, k)
		}
		for k := range closedspans {
			ckeys = append(ckeys, k)
		}
		sort.Ints(okeys)
		sort.Ints(ckeys)

		testone := make([]bool, rng)
		testtwo := make([]bool, rng)
		testthree := make([]bool, rng)
		for x := 0; x < rng; x++ {
			testone[x] = !(openspans[okeys[x+1]] != closedspans[ckeys[x]]) && (okeys[x+1] < ckeys[x])
			testtwo[x] = okeys[x+1] > ckeys[x] && openspans[okeys[x]] == closedspans[ckeys[x]]
			testthree[x] = okeys[x+1] < ckeys[x]
		}

		for x := 0; x < rng; x++ {
			balancetest = append(balancetest, []bool{testone[x], testtwo[x], testthree[x]})
		}
	}

	if boolslicecontains(invalidpattern, balancetest) {
		modifications := []modification{}
		boolslicereverse(balancetest)
		itemnumber := 0
		for len(balancetest) > 0 {
			testresult := balancetest[len(balancetest)-1]
			balancetest = balancetest[:len(balancetest)-1]
			if boolsliceequal(testresult, invalidpattern) {
				needInsertionAt := ckeys[itemnumber]
				insertionReopenTag := workline[openings[itemnumber+1][0]:openings[itemnumber+1][1]]
				insertionCloseTag := strings.Replace(insertionReopenTag, "<", "</", 1)
				modifications = append(modifications, modification{
					Item:     itemnumber,
					Position: needInsertionAt,
					CloseTag: insertionCloseTag,
					OpenTag:  insertionReopenTag,
				})
			}
			itemnumber++
		}

		newline := ""
		placeholder := 0
		for _, m := range modifications {
			newline += workline[placeholder:m.Position]
			newline += m.CloseTag
			newline += workline[closings[m.Item][0]:closings[m.Item][1]]
			newline += m.OpenTag
			placeholder = m.Position + len(workline[closings[m.Item][0]:closings[m.Item][1]])
		}
		newline += workline[placeholder:]

		workline = newline
	}

	return workline
}

func boolslicecontains(pattern []bool, tests [][]bool) bool {
	for _, test := range tests {
		if boolsliceequal(pattern, test) {
			return true
		}
	}
	return false
}

func boolsliceequal(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func boolslicereverse(arr [][]bool) {
	for i := len(arr)/2 - 1; i >= 0; i-- {
		opp := len(arr) - 1 - i
		arr[i], arr[opp] = arr[opp], arr[i]
	}
}
