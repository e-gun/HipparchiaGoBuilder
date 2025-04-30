//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"testing"
)

func TestAmpersandLatinFontsInAGreekText(t *testing.T) {

	// ttc = `█ⓕⓕ █⑨ⓕ █ⓕ④ █ⓕⓕ &<hb-title>EIS TON EUQUFPONA</hb-title> █⑨⑨ █⑧② █ⓔ① $1DI/KH.$ █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧②`
	ttc := `█ⓕⓕ &<hb-title>EIS TON EUQUFPONA</hb-title> █⑨⑨ █⑧② █ⓔ① $1DI/KH.$ █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧② █ⓔ① █ⓐⓒ █ⓔ② █ⓔ⑨ █ⓕ③ █ⓕⓕ $1GRAFH/N.$ █⑧⓪ H( U(PE\R DHMOSI/WN. █⑨⑨ █⑧② █ⓔ② $1A)GNW/S.`
	result := betacode.AmpersandLatinFontMarkup(ttc)
	want := `█ⓕⓕ <hb-fs-l-normal><hb-title>EIS TON EUQUFPONA</hb-title> </hb-fs-l-normal>█⑨⑨ █⑧② █ⓔ① $1DI/KH.$ █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧② █ⓔ① █ⓐⓒ █ⓔ② █ⓔ⑨ █ⓕ③ █ⓕⓕ $1GRAFH/N.$ █⑧⓪ H( U(PE\R DHMOSI/WN. █⑨⑨ █⑧② █ⓔ② $1A)GNW/S.`
	if result != want {
		t.Errorf("AmpersandLatinFontMarkup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestReplaceGreekFontMarkup(t *testing.T) {
	ttc := "μάλιϲτα δὲ ἀπεδέχετο [$9Epikur Epicurea p. 365, 16 Us. vgl. dessen Ind. S. 400$], φηϲὶ "
	result := betacode.DollarSignGreekFontMarkup(ttc)
	want := "μάλιϲτα δὲ ἀπεδέχετο [<hb-fs-g-regular>Epikur Epicurea p. 365, 16 Us. vgl. dessen Ind. S. 400</hb-fs-g-regular>], φηϲὶ "
	if result != want {
		t.Errorf("DollarSignGreekFontMarkup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}

	// ttc = `█ⓕⓕ █⑨ⓕ █ⓕ④ █ⓕⓕ &<hb-title>EIS TON EUQUFPONA</hb-title> █⑨⑨ █⑧② █ⓔ① $1DI/KH.$ █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧②`
	ttc = `█ⓔⓕ █⑧⓪ █ⓑ⑤ █ⓑ⓪ █ⓑ③ █ⓑ⑤ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ① █ⓕⓕ █ⓐⓕ █ⓒ⑤ █ⓕ⑤ █ⓕ④ █ⓔ⑧ █ⓕ⓪ █ⓔ⑧ █ⓕ② █ⓕⓕ █⑨ⓕ █ⓕ④ █ⓕⓕ &<hb-title>EIS TON EUQUFPONA</hb-title> █⑨⑨ █⑧② █ⓔ① $1DI/KH.$ █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧② █ⓔ① █ⓐⓒ █ⓔ② █ⓔ⑨ █ⓕ③ █ⓕⓕ $1GRAFH/N.$ █⑧⓪ H( U(PE\R DHMOSI/WN. █⑨⑨ █⑧② █ⓔ② $1A)GNW/S.`
	result = betacode.DollarSignGreekFontMarkup(ttc)
	want = `█ⓔⓕ █⑧⓪ █ⓑ⑤ █ⓑ⓪ █ⓑ③ █ⓑ⑤ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ① █ⓕⓕ █ⓐⓕ █ⓒ⑤ █ⓕ⑤ █ⓕ④ █ⓔ⑧ █ⓕ⓪ █ⓔ⑧ █ⓕ② █ⓕⓕ █⑨ⓕ █ⓕ④ █ⓕⓕ &<hb-title>EIS TON EUQUFPONA</hb-title> █⑨⑨ █⑧② █ⓔ① <hb-fs-g-bold>DI/KH.</hb-fs-g-bold> █⑧⓪ H( U(PE\R I)DIWTIKW=N A)DIKHMA/TWN KRI/SIS. █⑨ⓐ █⑧② █ⓔ① █ⓐⓒ █ⓔ② █ⓔ⑨ █ⓕ③ █ⓕⓕ <hb-fs-g-bold>GRAFH/N.</hb-fs-g-bold> █⑧⓪ H( U(PE\R DHMOSI/WN. █⑨⑨ █⑧② █ⓔ② $1A)GNW/S.`
	if result != want {
		t.Errorf("DollarSignGreekFontMarkup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestInstertCoptic(t *testing.T) {
	ttc := "παρὰ τῷ <hb-fs-g-coptic>*KEfFf*MfW</hb-fs-g-coptic> παρὰ τῷ"
	result := InstertCoptic(ttc)
	want := "παρὰ τῷ <hb-fs-g-coptic>ⲔⲉϥⲫϥⲘϥⲱ</hb-fs-g-coptic> παρὰ τῷ"
	if result != want {
		t.Errorf("InstertCoptic(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestConvertGreekCapitals(t *testing.T) {
	ttc := `*A*I*G*U*P*T*I*A*K*A. &E LIBRO PRIMO.$ *(ERMOTUMBIEI=S, MOI=RA TW=N MAXI/MWN �E)N *AI)GU/PTW|... *)/W|ETO GOU=N KAI\ O( *DIONU/SIOS`
	result := ConvertGreekCapitals(ttc)
	want := `ΑΙΓΥΠΤΙΑΚΑ. &E LIBRO PRIMO.$ ἙRMOTUMBIEI=S, MOI=RA TW=N MAXI/MWN �E)N ΑI)GU/PTW|... ᾬETO GOU=N KAI\ O( ΔIONU/SIOS`
	if result != want {
		t.Errorf("ConvertGreekCapitals(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestConvertGreekLowers(t *testing.T) {
	ttc := `*A*I*G*U*P*T*I*A*K*A. &E LIBRO PRIMO.$ *(ERMOTUMBIEI=S, MOI=RA TW=N MAXI/MWN �E)N *AI)GU/PTW|... *)/W|ETO GOU=N KAI\ O( *DIONU/SIOS`
	result := ConvertGreekLowers(ttc)
	want := `*α*ι*γ*υ*π*τ*ι*α*κ*α. &ε λιβρο πριμο.$ *(ερμοτυμβιεῖϲ, μοῖρα τῶν μαχίμων �ἐν *αἰγύπτῳ... *)/ῳετο γοῦν καὶ ὁ *διονύϲιοϲ`
	if result != want {
		t.Errorf("ConvertGreekLowers(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `<hb-tabbedtext />*TI/NES OU)=N H)=SAN ⟨AI(⟩ E)LPI/DES, E)F＇ AI(=S DIE/BAINEN`
	result = ConvertGreekLowers(ttc)
	want = `<hb-tabbedtext />*τίνεϲ οὖν ἦϲαν ⟨αἱ⟩ ἐλπίδεϲ, ἐφ＇ αἷϲ διέβαινεν`
	if result != want {
		t.Errorf("ConvertGreekLowers(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestConvertGreekUpperAndLower(t *testing.T) {
	ttc := `*A*I*G*U*P*T*I*A*K*A. &E LIBRO PRIMO.$ *(ERMOTUMBIEI=S, MOI=RA TW=N MAXI/MWN �E)N *AI)GU/PTW|... *)/W|ETO GOU=N KAI\ O( *DIONU/SIOS`
	result := ConvertGreekLowers(ConvertGreekCapitals(ttc))
	want := `ΑΙΓΥΠΤΙΑΚΑ. &ε λιβρο πριμο.$ Ἑρμοτυμβιεῖϲ, μοῖρα τῶν μαχίμων �ἐν Αἰγύπτῳ... ᾬετο γοῦν καὶ ὁ Διονύϲιοϲ`
	if result != want {
		t.Errorf("ConvertGreekUpperAndLower(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestGreekExclamation(t *testing.T) {
	ttc := `something !10'45 something`
	result := GreekExclamation(ttc)
	want := `something ∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙ ∙'45 something`
	if result != want {
		t.Errorf("GreekExclamation(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestGreekCombiningDot(t *testing.T) {
	ttc := `a?b?cd?ef`
	result := GreekCombiningDot(ttc)
	want := "a\u0323b\u0323cd\u0323ef"
	if result != want {
		t.Errorf("GreekCombiningDot(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestRestoreLatinWithinGreek(t *testing.T) {
	ttc := `ΑΙΓΥΠΤΙΑΚΑ. <hb-fs-l-normal>ε λιβρο πριμο.</hb-fs-l-normal> Ἑρμοτυμβιεῖϲ`
	result := RestoreLatinWithinGreek(ttc)
	want := "ΑΙΓΥΠΤΙΑΚΑ. <hb-fs-l-normal>E LIBRO PRIMO.</hb-fs-l-normal> Ἑρμοτυμβιεῖϲ"
	if result != want {
		t.Errorf("RestoreLatinWithinGreek(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestINS0090Block(t *testing.T) {
	ttc := ` <hb-fs-l-smallcapitals_italic> lines completely effaced █⑧⑧ █⑨⓪  — N?E? —  
  — KO —  
  — EI?S —  
  — ROI —  █ⓓ⓪ █⑨ⓕ █ⓔ⑥ █ⓔ① █ⓔ③ █ⓔ⑤ █ⓐ⓪ █ⓒ① █ⓕⓕ █⑧① █ⓔⓕ █ⓔ① █ⓒ④ █ⓔ⑤ █ⓔⓒ █ⓔⓕ █ⓕ③ █ⓕⓕ █ⓔⓕ █ⓔ④ █ⓔ③ █ⓐ⓪ █ⓑ① █ⓑ⑨ █ⓑ⓪ █ⓐ⓪ █ⓔ① █ⓑⓕ █ⓕⓕ █ⓔⓕ █ⓔ⑤ █ⓕⓕ █ⓔⓑ █ⓕⓐ █⑧① █⑧② <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /> — ON?∙3IL —  
 <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /> — ATATO∙∙TW —  
 <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /> —  [STA/MNON  —   — , E)]PIGRAFH/: A)PO\ TH=S  —`
	result := GreekCleanup(ttc)
	fmt.Println(result)
}

func TestAggregate(t *testing.T) {
	result := getallknownglcchars()
	fmt.Println(result)
	result = getallknowngucchars()
	fmt.Println(result)
}
