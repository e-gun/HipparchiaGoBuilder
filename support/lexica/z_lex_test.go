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

func TestFixfrontinus(t *testing.T) {
	x := `"Perseus:abo:phi,1245,001:Aquaed. 104"`
	y := fixfrontinus(x)
	z := `"Perseus:abo:phi,1245,002:104"`
	if y != z {
		t.Errorf("fixfrontinus(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixmartial(t *testing.T) {
	x := `bibl n="Perseus:abo:phi,1294,001:4:44:8" default`
	y := fixmartial(x)
	z := `bibl n="Perseus:abo:phi,1294,002:4:44:8" default`
	if y != z {
		t.Errorf("fixmartial(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixseneca(t *testing.T) {
	x := `"Perseus:abo:phi,1014,001:Ira. 3:18:3"`
	y := fixseneca(x)
	z := `"Perseus:abo:phi,1017,012:3:3:18:3"`
	fmt.Println(y)
	if y != z {
		t.Errorf("fixseneca(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
	x = `"Perseus:abo:phi,1014,001:Q. N. 5:16:5"`
	y = fixseneca(x)
	fmt.Println(y)
}

func TestFixVerrines(t *testing.T) {
	x := `abc <bibl n="Perseus:abo:phi,0474,005:2:58:section=142" default="NO" valid="yes"><author>id.</author> Verr. 2, 2, 58, § 142</bibl> abc`
	y := fixciceroverrines(x)
	z := `abc <bibl n="Perseus:abo:phi,0474,005:2:2:142" default="NO" valid="yes"><author>id.</author> Verr. 2, 2, 58, § 142</bibl> abc`
	fmt.Println(y)
	if y != z {
		t.Errorf("fixciceroverrines(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixCicSections(t *testing.T) {
	x := `abc <quote lang="la">omnes de tuā virtute commemorant,</quote> <bibl n="Perseus:abo:phi,0474,058:1:1:13:section=37" default="NO" valid="yes"><author>Cic.</author> abc`
	y := fixcicerosections(x)
	z := `abc <quote lang="la">omnes de tuā virtute commemorant,</quote> <bibl n="Perseus:abo:phi,0474,058:1:1:37" default="NO" valid="yes"><author>Cic.</author> abc`
	if y != z {
		t.Errorf("fixcicerosections(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestPurgesq(t *testing.T) {
	x := `abc <bibl n="Perseus:abo:phi,0472,001:61:130 sq" default="NO" valid="yes"><author>Cat.</author> 61, 130 sq.</bibl> abc`
	y := purgesq(x)
	z := ``
	if y != z {
		t.Errorf("purgesq(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixCicChh(t *testing.T) {
	x := `<bibl n="Perseus:abo:phi,0474,039:chapter=34" default="NO"><author>Cic.</author> Brut. 34</bibl>: accusatione desistere`
	y := fixcicerochapters(x)
	fmt.Println(y)
	z := `<bibl n=\"Perseus:abo:phi,0474,039:34\" default=\"NO\"><author>Cic.</author> Brut. 34</bibl>: accusatione desistere`
	if y != z {
		t.Errorf("fixcicerochapters(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestPurgelinklessbibls(t *testing.T) {
	x := `<bibl default="NO"><author>id.</author> ib. 4, 52</bibl> zzz <bibl n="Perseus:abo:phi,0474,005:2:58:section=142" default="NO" valid="yes"><author>id.</author> Verr. 2, 2, 58, § 142</bibl> <bibl z="InvalidPerseus:abo:phi,0474,039:34" default="NO"><author>Cic.</author> Brut. 34</bibl>: accusatione desistere`
	y := purgelinklessbibls(x)
	z := ``
	if y != z {
		t.Errorf("purgelinklessbibls(%q)\ngot\n%q\nwant\n%q", x, y, z)
	}
}

func TestFixVarro(t *testing.T) {
	x := `"Perseus:abo:phi,0684,001:L. L. 5:section=59"`
	y := fixvarro(x)
	z := `"Perseus:abo:phi,0684,001:5:59"`
	fmt.Println(y)
	if y != z {
		t.Errorf("fixvarro(%q)\ngot\n%q\nwant\n%q", x, y, z)
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
