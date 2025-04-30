package lexica

import (
	"testing"
)

//func TestParseAndLoadLSJ(t *testing.T) {
//	ParseAndLoadLSJ()
//}

func TestExtractdiv2material(t *testing.T) {
	hm := `<div2 id="crossu(ywth/s" orig_id="n110025" key="u(ywth/s" type="main" opt="n"><head extent="suff"...`
	_, b := extractdiv2material(hm)
	bwant := `<head extent="suff"...`
	if b != bwant {
		t.Errorf("extractdiv2material(%q)\ngot\n%q\nwant\n%q", hm, b, bwant)
	}

}

//func TestExtractSenses(t *testing.T) {
//	a, b := extractdiv2material(TESTENTRY)
//	a = extractlsjsenses(a, b)
//
//	want := 3
//	as := a.Senses
//	if len(as) != want {
//		t.Errorf("extractlsjsenses()\ngot\n%d\nwant\n%d", len(as), want)
//	}
//	for _, s := range as {
//		fmt.Println(s)
//	}
//}

func TestGeneratehyperlinks(t *testing.T) {
	ttc := `<bibl n="Perseus:abo:tlg,1342,001:23:6">`
	got := generatehyperlinks(ttc)
	want := `<bibl id="perseus/gr1342/001/23:6">`
	if got != want {
		t.Errorf("generatehyperlinks(%q)\ngot\n%q\nwant\n%q", ttc, got, want)
	}
}

func TestFindLatinLexEntries(t *testing.T) {
	_ = FindLatinLexEntries(LTLEXTESTDATA)
}
