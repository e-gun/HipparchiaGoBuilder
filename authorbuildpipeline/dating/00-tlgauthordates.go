//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dating

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseTLGDate(d string) int {
	d = strings.ReplaceAll(d, "?", "")
	containsBC := strings.Contains(d, "B.C.")
	containsAD := strings.Contains(d, "A.D.")
	containsBCandAD := containsBC && containsAD

	date := 9999
	if containsBCandAD {
		date = bcandad(d)
	} else if containsBC {
		date = justbc(d)
	} else if containsAD {
		date = justad(d)
	}
	return date
}

func justbc(d string) int {
	d = strings.Replace(d, "／", "–", 1)
	d = strings.ReplaceAll(d, "B.C.", "")
	ff := 0
	ff, d = postfudge(d)
	ff, d = antefudge(d)

	halves := strings.Split(d, "–")
	h1, e := strconv.Atoi(strings.TrimSpace(halves[0]))
	if e != nil {
		fmt.Println("justbca", d, e)
		return -9999
	}

	if len(halves) != 1 {
		h2, e2 := strconv.Atoi(strings.TrimSpace(halves[1]))
		if e2 != nil {
			fmt.Println("justbcb", d, e)
			return -9999
		}
		return int(-100*float32((h1+h2)/2) + 50 + float32(ff))
	}
	return (-100 * h1) + 50 + ff
}

func justad(d string) int {
	d = strings.Replace(d, "／", "–", 1)
	d = strings.ReplaceAll(d, "A.D.", "")
	ff := 0
	ff, d = postfudge(d)
	ff, d = antefudge(d)

	halves := strings.Split(d, "–")
	h1, e := strconv.Atoi(strings.TrimSpace(halves[0]))
	if e != nil {
		fmt.Println("dating - justad():", d, e)
		// next is not a dating error: it is a canon file parsing error
		// 4 key 4124 007 strconv.Atoi: parsing "4 key 4124 007": invalid syntax
		return -9999
	}

	if len(halves) != 1 {
		h2, e2 := strconv.Atoi(strings.TrimSpace(halves[1]))
		if e2 != nil {
			fmt.Println("justadb", d, e)
			return -9999
		}
		return int(100*float32((h1+h2)/2) - 50 + float32(ff))
	}
	return (100 * h1) - 50 + ff
}

func bcandad(d string) int {
	d = strings.Replace(d, "／", "–", 1)
	ff := 0
	halves := strings.Split(d, "–")
	h1 := halves[0]
	h2 := halves[1]

	var ff1 int
	var ff2 int

	ff1, h1 = postfudge(h1)
	ff1, h1 = antefudge(h1)

	ff2, h2 = postfudge(h2)
	ff2, h2 = antefudge(h2)
	ff = ff1 + ff2

	bc := nbc(h1)
	ad := nad(h2)

	date := bc + ad + ff
	return date
}

// postfudge - `p. 2 B.C.` --> -100 instead of -150
func postfudge(d string) (int, string) {
	pf := 0
	if strings.Contains(d, "p. ") {
		pf = 50
		d = strings.Replace(d, "p. ", "", 1)
	}
	return pf, d
}

// antefudge - `a. 1 B.C.` --> -100 instead of -50
func antefudge(d string) (int, string) {
	af := 0
	if strings.Contains(d, "a. ") {
		af = -50
		d = strings.Replace(d, "a. ", "", 1)
	}
	return af, d
}

// nbc - `4 B.C.` --> 350
func nbc(d string) int {
	d = strings.ReplaceAll(d, "B.C.", "")
	n, e := strconv.Atoi(strings.TrimSpace(d))
	if e != nil {
		fmt.Println("nbc", d, e)
		n = -9999
	}
	return (-100 * n) + 50
}

// nad - `A.D. 11` --> 1050
func nad(d string) int {
	d = strings.ReplaceAll(d, "A.D.", "")
	n, e := strconv.Atoi(strings.TrimSpace(d))
	if e != nil {
		fmt.Println("nad", d, e)
		n = -9999
	}
	return (100 * n) - 50
}

var (
	thedates = []string{
		"1 B.C.",
		"1 B.C.?",
		"1 B.C.–A.D. 1",
		"1 B.C.–A.D. 1?",
		"1 B.C.／A.D. 1",
		"1 B.C.／A.D. 1?",
		"1 B.C.／A.D. 4?",
		"2 B.C.",
		"2 B.C.?",
		"2 B.C.–A.D. 4",
		"2 B.C.／A.D. 2",
		"2 B.C.／A.D. 2?",
		"2 B.C.／A.D. 3",
		"2–1 B.C.",
		"2–1 B.C.?",
		"2／1 B.C.",
		"2／1 B.C.?",
		"3 B.C.",
		"3 B.C.?",
		"3 B.C.?／A.D. 1",
		"3 B.C.／A.D. 1",
		"3 B.C.／A.D. 2",
		"3–2 B.C.",
		"3–2 B.C.?",
		"3／1 B.C.",
		"3／2 B.C.",
		"3／2 B.C.?",
		"4 B.C.",
		"4 B.C.?",
		"4 B.C.／A.D. 1",
		"4 B.C.／A.D. 2",
		"4?／2 B.C.",
		"4–2 B.C.",
		"4–3 B.C.",
		"4–3 B.C.?",
		"4／1 B.C.",
		"4／1 B.C.?",
		"4／2 B.C.",
		"4／2 B.C.?",
		"4／3 B.C.",
		"4／3 B.C.?",
		"5 B.C.",
		"5 B.C.?",
		"5–4 B.C.",
		"5–4 B.C.?",
		"5／3 B.C.",
		"5／3 B.C.?",
		"5／4 B.C.",
		"5／4 B.C.?",
		"6 B.C.",
		"6 B.C.?",
		"6–5 B.C.",
		"6／5 B.C.",
		"7 B.C.",
		"7–6 B.C.",
		"7／6 B.C.",
		"7／6 B.C.?",
		"8 B.C.",
		"8–6 B.C.",
		"8／6 B.C.",
		"8／6 B.C.?",
		"8／7 B.C.",
		"8／7 B.C.?",
		"A.D. 1",
		"A.D. 10",
		"A.D. 10／15",
		"A.D. 11",
		"A.D. 11–12",
		"A.D. 12",
		"A.D. 12?",
		"A.D. 12–13",
		"A.D. 13",
		"A.D. 13–14",
		"A.D. 14",
		"A.D. 14–15",
		"A.D. 15",
		"A.D. 15–16",
		"A.D. 1?",
		"A.D. 1?／6",
		"A.D. 1–2",
		"A.D. 1–2?",
		"A.D. 1–7",
		"A.D. 1／2",
		"A.D. 1／2?",
		"A.D. 1／3",
		"A.D. 2",
		"A.D. 2?",
		"A.D. 2?／4",
		"A.D. 2–3",
		"A.D. 2–3?",
		"A.D. 2／3",
		"A.D. 2／3?",
		"A.D. 2／4",
		"A.D. 2／4?",
		"A.D. 3",
		"A.D. 3?",
		"A.D. 3–4",
		"A.D. 3／10",
		"A.D. 3／4",
		"A.D. 3／4?",
		"A.D. 3／5",
		"A.D. 4",
		"A.D. 4?",
		"A.D. 4–5",
		"A.D. 4–5?",
		"A.D. 4／5",
		"A.D. 4／5?",
		"A.D. 4／6",
		"A.D. 4／6?",
		"A.D. 5",
		"A.D. 5?",
		"A.D. 5–6",
		"A.D. 5–6?",
		"A.D. 5／10",
		"A.D. 5／6",
		"A.D. 5／6?",
		"A.D. 5／7",
		"A.D. 6",
		"A.D. 6?",
		"A.D. 6–10",
		"A.D. 6–7",
		"A.D. 6／13",
		"A.D. 6／7",
		"A.D. 7",
		"A.D. 7?",
		"A.D. 7–8",
		"A.D. 7／8",
		"A.D. 7／9",
		"A.D. 8",
		"A.D. 8–9",
		"A.D. 8–9?",
		"A.D. 8／10",
		"A.D. 9",
		"A.D. 9?",
		"A.D. 9–10",
		"A.D. 9／10",
		"Incertum",
		"Varia",
		"a. 1 B.C.",
		"a. 1 B.C.?",
		"a. 2 B.C.",
		"a. 2 B.C.?",
		"a. 3 B.C.",
		"a. 3 B.C.?",
		"a. 4 B.C.",
		"a. 4 B.C.?",
		"a. 5 B.C.",
		"a. 6 B.C.",
		"a. A.D. 1",
		"a. A.D. 10",
		"a. A.D. 11",
		"a. A.D. 12",
		"a. A.D. 14",
		"a. A.D. 14／15",
		"a. A.D. 15",
		"a. A.D. 1?",
		"a. A.D. 1／2",
		"a. A.D. 2",
		"a. A.D. 2／3",
		"a. A.D. 3",
		"a. A.D. 3?",
		"a. A.D. 4",
		"a. A.D. 5",
		"a. A.D. 5?",
		"a. A.D. 8",
		"p. 1 B.C.",
		"p. 1 B.C.?",
		"p. 2 B.C.",
		"p. 3 B.C.",
		"p. 4 B.C.",
		"p. 4 B.C.?",
		"p. 4 B.C.／a. A.D. 2",
		"p. 5 B.C.",
		"p. 7 B.C.",
		"p. A.D. 1",
		"p. A.D. 10",
		"p. A.D. 2",
		"p. A.D. 3",
		"p. A.D. 4",
		"p. A.D. 5",
		"p. A.D. 6",
		"p. A.D. 7",
		"p. A.D. 9／10",
	}
)
