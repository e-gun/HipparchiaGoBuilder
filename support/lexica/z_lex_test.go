package lexica

import (
	"fmt"
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

func TestFixeuripides(t *testing.T) {
	x := `abc <bibl n="Perseus:abo:tlg,0006,018:601"><author>E.</author> <title>IA</title> 601</bibl> abc`
	y := fixeuripides(x)
	z := `abc <bibl n="Perseus:abo:tlg,0006,051:601"><author>E.</author> <title>IA</title> 601</bibl> abc`
	if y != z {
		t.Errorf("fixeuripides(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixnepos(t *testing.T) {
	x := `abc <bibl n="Perseus:abo:phi,0588,001:Alcib. 11:4" default="NO" valid="yes"><author>Nep.</author> Alcib. 11, 4</bibl></cit> abc`
	y := fixnepos(x)
	z := `abc <bibl n="Perseus:abo:phi,0588,001:Alc:11:4" default="NO" valid="yes"><author>Nep.</author> Alcib. 11, 4</bibl></cit> abc`
	if y != z {
		t.Errorf("fixnepos(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixsallust(t *testing.T) {
	x := `abc <bibl n="Perseus:abo:phi,0631,001:J. 57 sq" default="NO" valid="yes"><author>Sall.</author> J. 57 sq.</bibl> abc`
	y := fixsallust(x)
	z := `abc <bibl n="Perseus:abo:phi,0631,002:57" default="NO" valid="yes"><author>Sall.</author> J. 57 sq.</bibl> abc`
	if y != z {
		t.Errorf("fixsallust(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixsuetonius(t *testing.T) {
	x := `"Perseus:abo:phi,1348,001:life=vesp.:16"`
	y := fixsuetonius(x)
	z := `"Perseus:abo:phi,1348,001:Ves:16"`
	if y != z {
		t.Errorf("fixsuetonius(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixibidem(t *testing.T) {
	x := `"Perseus:abo:phi,1014,001:ib. 76:2"`
	y := fixibidem(x)
	z := `"Perseus:abo:phi,1014,001:76:2"`
	fmt.Println(y)
	if y != z {
		t.Errorf("fixibidem(%q)\ngot\n%q\nwant\n%q", x, y, z)
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
