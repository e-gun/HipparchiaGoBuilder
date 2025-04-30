//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"testing"
)

func TestReplaceQuotationMarks(t *testing.T) {
	ttc := `Test of "6 00AB 00BB "6 "Left-Pointing Double Angle Quotation Mark" and "Right-Pointing Double Angle Quotation Mark"`
	result := ReplaceQuotationMarks(ttc)

	var want string
	if !SIMPLEQUOTES {
		want = `Test of « 00AB 00BB » “Left-Pointing Double Angle Quotation Mark” and “Right-Pointing Double Angle Quotation Mark”`
	} else {
		want = `Test of “ 00AB 00BB ” “Left-Pointing Double Angle Quotation Mark” and “Right-Pointing Double Angle Quotation Mark”`
	}

	if result != want {
		t.Errorf("ReplaceQuotationMarks(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestInsertNBSP(t *testing.T) {
	ttc := `nos ambae ^10 faciunt in hoc`
	result := InsertBlankQuarterSpaces(ttc)
	want := `nos ambae <hmu_blank_quarter_spaces quantity="10" />  faciunt in hoc`
	if result != want {
		t.Errorf("TestInsertNBSP(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplacePercentSigns(t *testing.T) {
	ttc := "Test of %14 00A7 § Section Sign"
	result := ReplacePercentSigns(ttc)
	want := "Test of § 00A7 § Section Sign"
	if result != want {
		t.Errorf("ReplacePercentSigns(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplacePoundSigns(t *testing.T) {
	ttc := "Test of #10 03FD Ͻ Greek Capital Reversed Lunate Sigma Symbol"
	result := ReplacePoundSigns(ttc)
	want := "Test of Ͻ 03FD Ͻ Greek Capital Reversed Lunate Sigma Symbol"
	if result != want {
		t.Errorf("ReplacePoundSigns(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceLeftSquareBrackets(t *testing.T) {
	ttc := "Test of [14 007C+003A |: Vertical Line and Colon ‣ Left Hymn Refrain Bracket"
	result := ReplaceLeftSquareBrackets(ttc)
	want := "Test of |: 007C+003A |: Vertical Line and Colon ‣ Left Hymn Refrain Bracket"
	if result != want {
		t.Errorf("ReplaceLeftSquareBrackets(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceRightSquareBrackets(t *testing.T) {
	ttc := "Test of ]14 003A+007C :| Colon and Vertical Line ‣ Right Hymn Refrain Bracket"
	result := ReplaceRightSquareBrackets(ttc)
	want := "Test of :| 003A+007C :| Colon and Vertical Line ‣ Right Hymn Refrain Bracket"
	if result != want {
		t.Errorf("ReplaceRightSquareBrackets(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceAtSigns(t *testing.T) {
	ttc := "Test of @6 Blank Line"
	result := ReplaceAtSigns(ttc)
	want := "Test of <br> Blank Line"
	if result != want {
		t.Errorf("ReplaceAtSigns(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceLeftCurlyBrackets(t *testing.T) {
	ttc := "Test of {41 Stage Direction"
	result := ReplaceLeftCurlyBrackets(ttc)
	want := "Test of <hb-sp-stagedirection> Stage Direction"
	if result != want {
		t.Errorf("ReplaceLeftSquareBrackets(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceRightCurlyBrackets(t *testing.T) {
	ttc := "Test of }41 Stage Direction"
	result := ReplaceRightCurlyBrackets(ttc)
	want := "Test of </hb-sp-stagedirection> Stage Direction"
	if result != want {
		t.Errorf("TestReplaceRightCurlyBrackets(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceLeftAngledBrackets(t *testing.T) {
	ttc := "<16 2035 2032 >16 : ‵ Left Interlinear MarkedUp Marker Bracket Right Interlinear MarkedUp Marker Bracket ′ "
	result := ReplaceLeftAngledBrackets(ttc)
	want := "‵ 2035 2032 >16 : ‵ Left Interlinear MarkedUp Marker Bracket Right Interlinear MarkedUp Marker Bracket ′ "
	if result != want {
		t.Errorf("ReplaceLeftAngledBrackets(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceRightAngledBrackets(t *testing.T) {
	ttc := "<16 2035 2032 >16 : ‵ Left Interlinear MarkedUp Marker Bracket Right Interlinear MarkedUp Marker Bracket ′ "
	result := ReplaceRightAngledBrackets(ttc)
	want := "<16 2035 2032 ′ : ‵ Left Interlinear MarkedUp Marker Bracket Right Interlinear MarkedUp Marker Bracket ′ "
	if result != want {
		t.Errorf("ReplaceLeftAngledBrackets(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceSingletons(t *testing.T) {
	ttc := "% dagger; { spkropen; } spkrclose"
	result := ReplaceSingletons(ttc)
	want := "† dagger; <speaker> spkropen; </speaker> spkrclose"
	if result != want {
		t.Errorf("ReplaceSingletons(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestBetaCodeCleanup(t *testing.T) {
	ttc := "% dagger; { spkropen; } spkrclose; <16 2035 2032 >16; {41 Stage Direction }41;  ]14 003A+007C"
	result := BetaCodeCleanup(ttc)
	want := "† dagger; <speaker> spkropen; </speaker> spkrclose; ‵ 2035 2032 ′; <hb-sp-stagedirection> Stage Direction </hb-sp-stagedirection>;  :| 003A+007C"
	if result != want {
		t.Errorf("TestBetaCodeCleanup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func BenchmarkBetaCodeCleanup(b *testing.B) {
	ttc := "% dagger; { spkropen; } spkrclose; <16 2035 2032 >16; {41 Stage Direction }41;  ]14 003A+007C"
	for i := 0; i < b.N; i++ {
		_ = BetaCodeCleanup(ttc)
	}
}

func TestLatinFontLineMarkupSubs(t *testing.T) {
	ttc := `⑧⓪ *AI)GUPTIAKW=N PRW/TW|. █⑨⓪ █⑧ⓕ █ⓕ④ █ⓕⓕ {1&E LIBRO SECUNDO.$}1 █⑧① @&Idem$%10 *NIKI/OU KW/MH *AI)GU/PTOU. *)AR. *AI)GUPTIAKW=N █⑧⓪ `
	result := AmpersandLatinFontMarkup(ttc)
	want := `⑧⓪ *AI)GUPTIAKW=N PRW/TW|. █⑨⓪ █⑧ⓕ █ⓕ④ █ⓕⓕ {1<hb-fs-l-normal>E LIBRO SECUNDO.</hb-fs-l-normal>}1 █⑧① @<hb-fs-l-normal>Idem</hb-fs-l-normal>%10 *NIKI/OU KW/MH *AI)GU/PTOU. *)AR. *AI)GUPTIAKW=N █⑧⓪ `
	if result != want {
		t.Errorf("AmpersandLatinFontMarkup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	// the next will fail; you need SimpleLatinSpanLATE() to handle it; otherwise all but the last {&7Ph.&} work while
	// the final one remains `{&7Ph.&}`

	//ttc = `{&7Pa.&} censeo. {&7Ph.&} sed heus tu. {&7Pa.&} quid vis? {&7Ph.&} censen posse me obfirmare et `
	//result = AmpersandLatinFontMarkup(ttc)
	//fmt.Println(result)

}

func TestSimplegktolatinshift(t *testing.T) {
	ttc := "<hmutitle>&E LIBRO SECUNDO.$</hmutitle> █⑧① <hmu_standalone_tabbedtext />&Idem$﹕"
	result := simplegktolatinspan(ttc)
	want := "<hmutitle><hb-fs-l-normal>E LIBRO SECUNDO.</hb-fs-l-normal></hmutitle> █⑧① <hmu_standalone_tabbedtext /><hb-fs-l-normal>Idem</hb-fs-l-normal>﹕"
	if result != want {
		t.Errorf("simplegktolatinspan(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `LUSAN. *TAU=TA KAI\ *)ARTEMI/DWRO/S FHSIN. █⑨⓪ <hmu_standalone_tabbedtext />&Plutarch. De Is. et Os. c. 5$﹕ *OI( DE\`
	result = simplegktolatinspan(ttc)
	want = `LUSAN. *TAU=TA KAI\ *)ARTEMI/DWRO/S FHSIN. █⑨⓪ <hmu_standalone_tabbedtext /><hb-fs-l-normal>Plutarch. De Is. et Os. c. 5</hb-fs-l-normal>﹕ *OI( DE\`
	if result != want {
		t.Errorf("simplegktolatinspan(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestLatinsubshift(t *testing.T) {
	ttc := `&3pserint [1&sc. de pyramidibus]1, &3sunt Herodotus`
	result := latinsubshiftspan(ttc)
	want := `<hb-fs-l-italic>pserint [1</hb-fs-l-italic><hb-fs-l-normal>sc. de pyramidibus]1, </hb-fs-l-normal>&3sunt Herodotus`
	if result != want {
		t.Errorf("latinsubshiftspan(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestComplexgktolatinspan(t *testing.T) {
	ttc := `greek &10stuff_in_shifted_latin_font$10 greek`
	result := complexgktolatinspan(ttc)
	want := `greek <hb-fs-l-smallerthannormal>stuff_in_shifted_latin_font</hb-fs-l-smallerthannormal> greek`
	if result != want {
		t.Errorf("complexgktolatinspan(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `<hmu_standalone_tabbedtext />&10PLATO Cratyl. 429 D$10 *SW. *)=ARA O(/TI YEUDH= LE/GEIN TO\ PARA/PAN OU)K E)/STIN`
	result = complexgktolatinspan(ttc)
	want = `<hmu_standalone_tabbedtext /><hb-fs-l-smallerthannormal>PLATO Cratyl. 429 D</hb-fs-l-smallerthannormal> *SW. *)=ARA O(/TI YEUDH= LE/GEIN TO\ PARA/PAN OU)K E)/STIN`
	if result != want {
		t.Errorf("complexgktolatinspan(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestDollarfinderswap(t *testing.T) {
	ttc := `█⑨⑨ █⑧② █ⓔ① $1DI/KH.$ █⑧⓪ H(`
	result := dollarfinderswap(ttc)
	want := "█⑨⑨ █⑧② █ⓔ① <hb-fs-g-bold>DI/KH.</hb-fs-g-bold> █⑧⓪ H("
	if result != want {
		t.Errorf("dollarfinderswap(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `█ⓔⓕ █⑧⓪ █ⓑ⑤ █ⓑ⓪ █ⓑ③ █ⓑ⑤ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ① █ⓕⓕ █ⓐⓕ █ⓒ⑤ █ⓕ⑤ █ⓕ④ █ⓔ⑧ █ⓕ⓪ █ⓔ⑧ █ⓕ② █ⓕⓕ █⑨ⓕ █ⓕ④ █ⓕⓕ &<hb-title>EIS TON EUQUFPONA</hb-title> █⑨⑨ █⑧② █ⓔ① $1DI/KH.$ █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧② █ⓔ① █ⓐⓒ █ⓔ② █ⓔ⑨ █ⓕ③ █ⓕⓕ $1GRAFH/N.$ █⑧⓪ H( U(PE\R DHMOSI/WN. █⑨⑨ █⑧② █ⓔ② $1A)GNW/S.`
	result = dollarfinderswap(ttc)
	want = `█ⓔⓕ █⑧⓪ █ⓑ⑤ █ⓑ⓪ █ⓑ③ █ⓑ⑤ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ① █ⓕⓕ █ⓐⓕ █ⓒ⑤ █ⓕ⑤ █ⓕ④ █ⓔ⑧ █ⓕ⓪ █ⓔ⑧ █ⓕ② █ⓕⓕ █⑨ⓕ █ⓕ④ █ⓕⓕ &<hb-title>EIS TON EUQUFPONA</hb-title> █⑨⑨ █⑧② █ⓔ① <hb-fs-g-bold>DI/KH.</hb-fs-g-bold> █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧② █ⓔ① █ⓐⓒ █ⓔ② █ⓔ⑨ █ⓕ③ █ⓕⓕ <hb-fs-g-bold>GRAFH/N.</hb-fs-g-bold> █⑧⓪ H( U(PE\R DHMOSI/WN. █⑨⑨ █⑧② █ⓔ② $1A)GNW/S.`
	if result != want {
		t.Errorf("dollarfinderswap(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceQuotationMarks2(t *testing.T) {
	q := "RI/AS E)PEXA/RATTON, OU) TH=S METRIO/THTOS A)LLA\\ TH=S DUNA/MEWS @1 \n TH=S *)ALECA/NDROU STOXAZO/MENOI, SKOPW=MEN: [1&Anth. Pal. \n &6XVI &`120 Preger Insc. gr. metr. 279$]1 █⑧⑥ @6 \n @@\"3AU)DASOU=NTI D' E)/OIKEN O( XA/LKEOS EI)S *DI/A LEU/SSWN: \n @@*GA=N U(P' E)MOI\\ TI/QEMAI: *ZEU=, SU\\ D' *)/OLUMPON E)/XE.\"3 █⑧⑧ █⑧⑧ @6 \n KAI\\ \"3*)ALE/CANDROS E)GW\\ *DIO\\S ME\\N UI(O/S\"3. TAU=TA ME\\N OU)=N. \n W(S E)/FHN, OI( POIHTAI\\ KOLAKEU/ONTES AU)TOU= TH\\N TU/XHN \n PROSEI=PON: TW=N D' A)LHQINW=N A)POFQEGMA/TWN *)ALECA/NDROU █⑨⓪ PRW=TON A)/N TIS TA\\ PAIDIKA\\ DIE/LQOI. PODWKE/STATOS GA\\R \n TW=N E)F' H(LIKI/AS NE/WN GENO/MENOS KAI\\ TW=N E(TAI/RWN AU)TO\\N \n E)P' *)OLU/MPIA PARORMW/NTWN, H)RW/THSEN EI) BASILEI=S \n A)GWNI/ZONTAI: TW=N D' \"3OU)/\"3 FAME/NWN A)/DIKON EI)=PEN EI)=NAI TH\\N \n A(/MILLAN, E)N H(=| NIKH/SEI ME\\N I)DIW/TAS NIKHQH/SETAI DE\\ BASI-\n LEU/S. TOU= DE\\ PATRO\\S *FILI/PPOU LO/GXH| TO\\N MHRO\\N E)N *TRIBAL-\n LOI=S DIAPARE/NTOS KAI\\ TO\\N ME\\N KI/NDUNON DIAFUGO/NTOS \n A)XQOME/NOU DE\\ TH=| XWLO/THTI, \"3QA/RREI PA/TER\"3 E)/FH \"3KAI\\ \n PRO/IQI FAIDRW=S, I(/NA TH=S A)RETH=S KATA\\ BH=MA MNHMONEU/H|S.\"3 \n TAU=T' OU)K E)/STI DIANOI/AS FILOSO/FOU KAI\\ DIA\\ TO\\N E)PI\\ TOI=S █⑨⓪ KALOI=S E)NQOUSIASMO\\N H)/DH TW=N TOU= SW/MATOS E)LATTWMA/TWN "
	r := ReplaceQuotationMarks(q)
	fmt.Println(r)
}

func TestTLG5038(t *testing.T) {
	ttc := `<hb-sp-marginaltext>&K&4g&U&4g&E&4g$</hb-sp-marginaltext>`
	r := complexgktolatinspan(ttc)
	fmt.Println(r)
}
