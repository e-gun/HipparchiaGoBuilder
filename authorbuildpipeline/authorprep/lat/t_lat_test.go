//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lat

import "testing"

func TestConvertLatinDiacriticals(t *testing.T) {
	ttc := `u=V=i+e/o\`
	result := ConvertLatinDiacriticals(ttc)
	want := "ûÛïéò"
	if result != want {
		t.Errorf("ConvertLatinDiacriticals(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestRewritegreekfontshift(t *testing.T) {
	ttc := `L. Aelium magistrum suum in $E)TUMOLOGI/A|<hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	result := rewritegreekfontshift(ttc)
	want := `L. Aelium magistrum suum in <hb-fs-g-normal>E)TUMOLOGI/A|</hb-fs-g-normal><hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	if result != want {
		t.Errorf("rewritegreekfontshift(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `$KU/BON<hb-fs-l-normal>, quid </hb-fs-l-normal>GRAMMH/N<hb-fs-l-normal>; quibusque ista omnia Latinis uoca-</hb-fs-l-normal> █`
	result = rewritegreekfontshift(ttc)
	want = `<hb-fs-g-normal>KU/BON</hb-fs-g-normal><hb-fs-l-normal>, quid </hb-fs-l-normal><hb-fs-g-normal>GRAMMH/N</hb-fs-g-normal><hb-fs-l-normal>; quibusque ista omnia Latinis uoca-</hb-fs-l-normal> █`
	if result != want {
		t.Errorf("rewritegreekfontshift(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
	ttc = `<hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext />$EU)FHMEI=N XRH\ KA)CI/STASQAI TOI=S H(METE/ROISI XOROI=SIN, █`
	result = rewritegreekfontshift(ttc)
	want = `<hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-fs-g-normal>EU)FHMEI=N XRH\ KA)CI/STASQAI TOI=S H(METE/ROISI XOROI=SIN,</hb-fs-g-normal> █`
	if result != want {
		t.Errorf("rewritegreekfontshift(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestApplybetacodeconversion(t *testing.T) {
	ttc := `L. Aelium magistrum suum in <hb-fs-g-normal>E)TUMOLOGI/A|</hb-fs-g-normal><hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	result := applybetacodeconversion(ttc)
	want := `L. Aelium magistrum suum in <hb-fs-g-normal>ἐτυμολογίᾳ</hb-fs-g-normal><hb-fs-l-normal> falsa reprehendit; </hb-fs-l-normal> █`
	if result != want {
		t.Errorf("applybetacodeconversion(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestLatinCleanup(t *testing.T) {
	ttc := `{&7Pa.&} censeo. {&7Ph.&} sed heus tu. {&7Pa.&} quid vis? {&7Ph.&} censen posse me obfirmare et `
	result := LatinCleanup(ttc)
	want := `{<hb-fs-l-smallcapitals>Pa.</hb-fs-l-smallcapitals><hb-fs-l-normal>} censeo. {</hb-fs-l-normal><hb-fs-l-smallcapitals>Ph.</hb-fs-l-smallcapitals>} sed heus tu. {<hb-fs-l-smallcapitals>Pa.</hb-fs-l-smallcapitals><hb-fs-l-normal>} quid vis? {</hb-fs-l-normal><hb-fs-l-smallcapitals>Ph.</hb-fs-l-smallcapitals>} censen posse me obfirmare et `
	if result != want {
		t.Errorf("LatinCleanup(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}
