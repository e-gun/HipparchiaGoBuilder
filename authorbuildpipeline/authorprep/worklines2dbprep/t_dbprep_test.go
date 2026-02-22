//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"fmt"
	"testing"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func TestStripString(t *testing.T) {
	ttc := `ἔντοϲ ἀμώμητον κάλλιπον οὐκ ἐθέλων`
	result := stripstring(ttc)
	want := "εντοϲ αμωμητον καλλιπον ουκ εθελων"
	if result != want {
		t.Errorf("stripstring(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestConsolidatecontiguous(t *testing.T) {
	p := structs.DbWorkline{
		Accented: "nos ambae faciunt in hoc tempore summa gratia et elo-",
		Stripped: "nos ambae faciunt in hoc tempore summa gratia et elo-",
	}
	n := structs.DbWorkline{
		Accented: "quentia quarum alteram C Aquili uereor alteram metuo",
		Stripped: "quentia quarum alteram C Aquili uereor alteram metuo",
	}
	p2, n2 := consolidatecontiguous(p, n)

	// you need
	// p: nos ambae faciunt in hoc tempore, summa gratia et eloquentia
	// t: quarum alteram, C. Aquili, uereor, alteram metuo.

	w1 := `nos ambae faciunt in hoc tempore summa gratia et eloquentia`
	w2 := `quarum alteram C Aquili uereor alteram metuo`
	if p2.Accented != w1 {
		t.Errorf("#1 consolidatecontiguous(%q)\ngot\n%q\nwant\n%q", p.Accented, p2.Accented, w1)
	}
	if n2.Accented != w2 {
		t.Errorf("#2 consolidatecontiguous(%q)\ngot\n%q\nwant\n%q", n.Accented, n2.Accented, w2)
	}

	if p2.Stripped != w1 {
		t.Errorf("#3 consolidatecontiguous(%q)\ngot\n%q\nwant\n%q", p.Stripped, p2.Stripped, w1)
	}
	if n2.Stripped != w2 {
		t.Errorf("#4 consolidatecontiguous(%q)\ngot\n%q\nwant\n%q", n.Stripped, n2.Stripped, w2)
	}
}

func TestFixHyphens(t *testing.T) {
	p := structs.DbWorkline{
		Accented: "nos ambae faciunt in hoc tempore summa gratia et elo-",
		Stripped: "nos ambae faciunt in hoc tempore summa gratia et elo-",
	}
	n := structs.DbWorkline{
		Accented: "quentia quarum alteram C Aquili uereor alteram metuo",
		Stripped: "quentia quarum alteram C Aquili uereor alteram metuo",
	}
	tl := []structs.DbWorkline{p, n}
	rr := FixHyphens(tl)
	p2 := rr[0]
	n2 := rr[1]

	w1 := `nos ambae faciunt in hoc tempore summa gratia et eloquentia`
	w2 := `quarum alteram C Aquili uereor alteram metuo`
	if p2.Accented != w1 {
		t.Errorf("#1 FixHyphens(%q)\ngot\n%q\nwant\n%q", p.Accented, p2.Accented, w1)
	}
	if n2.Accented != w2 {
		t.Errorf("#2 FixHyphens(%q)\ngot\n%q\nwant\n%q", n.Accented, n2.Accented, w2)
	}

	if p2.Stripped != w1 {
		t.Errorf("#3 FixHyphens(%q)\ngot\n%q\nwant\n%q", p.Stripped, p2.Stripped, w1)
	}
	if n2.Stripped != w2 {
		t.Errorf("#4 FixHyphens(%q)\ngot\n%q\nwant\n%q", n.Stripped, n2.Stripped, w2)
	}
}

func TestInsertNBSP(t *testing.T) {
	ttc := `nos ambae <hb_blank_quarter_spaces quantity="15" /> faciunt in <hb-tabbedtext /> hoc`
	n := structs.DbWorkline{
		MarkedUp: ttc,
	}
	tl := []structs.DbWorkline{n}
	result := InsertNBSP(tl)
	want := "nos ambae &nbsp;&nbsp;&nbsp; faciunt in &nbsp;&nbsp;&nbsp; hoc"
	if result[0].MarkedUp != want {
		t.Errorf("InsertNBSP(%q)\ngot\n%q\nwant\n%q", ttc, result[0].MarkedUp, want)
	}
}

func TestExtractNotesAndMarginalText(t *testing.T) {
	n := structs.DbWorkline{
		MarkedUp: `ἄρχ⟨ου⟩ϲι <hb-sp-unconventional_form_written_by_scribe>αρχϲι</hb-sp-unconventional_form_written_by_scribe> Πανοπολιτῶν. τίνα εἰϲ κοινὸ[ν ἐγράφη] ἐ̣μοὶ καὶ ὑμῖν καὶ κονδούκτορϲι <hb-sp-unconventional_form_written_by_scribe>κοντουκτορϲι</hb-sp-unconventional_form_written_by_scribe> περὶ τοῦ ϲυϲτήϲαϲθαι ἁ̣[λιάδαϲ εἰϲ τὴν τῶν]`,
	}
	ll := []structs.DbWorkline{n}
	ll = ExtractNotesAndMarginalText(ll)
	fmt.Println("mu", ll[0].MarkedUp)
	fmt.Println("an", ll[0].Annotations)

}

func TestNoLeadingOrTrailingSpaces(t *testing.T) {
	ttc := ` nos ambae hoc `
	n := structs.DbWorkline{
		MarkedUp: ttc,
		Accented: ttc,
		Stripped: ttc,
	}
	tl := []structs.DbWorkline{n}
	result := NoLeadingOrTrailingSpaces(tl)
	want := "nos ambae hoc"
	if result[0].MarkedUp != want {
		t.Errorf("T: InsertNBSP(%q)\ngot\n%q\nwant\n%q", ttc, result[0].MarkedUp, want)
	}
	if result[0].Accented != want {
		t.Errorf("A: InsertNBSP(%q)\ngot\n%q\nwant\n%q", ttc, result[0].Accented, want)
	}
	if result[0].Stripped != want {
		t.Errorf("S: InsertNBSP(%q)\ngot\n%q\nwant\n%q", ttc, result[0].Stripped, want)
	}
}

func TestAssignCitationValues(t *testing.T) {
	// trying to apply "incr by 1" to "501a"
	lines := `<hb-set_l_1_to_501a />τοὐναντίον (διὸ καὶ μόνον οὐ τὴν αὐτὴν κίνηϲιν ποιεῖ-
<hb-incr_l_0_by_1 />ται τῆϲ πορείαϲ νέοϲ ὢν καὶ τελειωθείϲ, ἀλλὰ τὸ πρῶτον παι-
<hb-incr_l_0_by_1 />δίον ὂν ἕρπει τετραποδίζον), τὰ δ’ ἀνὰ λόγον ἀποδίδωϲι τὴν
<hb-incr_l_0_by_1 />ϲύριγγοϲ καὶ ϲάλπιγγοϲ, ταχὺ δὲ θεῖν οὐχ ἧττον τῶν ἐλά-
<hb-incr_l_1_by_1 />φων, καὶ εἶναι ἄγριον καὶ ἀνθρωποφάγον. 
<hb-set_l_0_to_1 /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext />Ἄνθρωποϲ μὲν οὖν 
<hb-incr_l_0_by_1 />βάλλει τοὺϲ ὀδόνταϲ, βάλλει δὲ καὶ ἄλλα τῶν ζῴων, οἷον 
<hb-incr_l_0_by_1 />ἵπποϲ καὶ ὀρεὺϲ καὶ ὄνοϲ. Βάλλει δ’ ἄνθρωποϲ τοὺϲ προϲθίουϲ, `
	result := AssignCitationValues(lines)
	//for _, r := range result {
	//	fmt.Println(r)
	//}
	if result[5].Lvl1Value != `501b` {
		t.Errorf("T: AssignCitationValues(%q)\ngot\n%q\nwant\n%q", `<hb-incr_l_1_by_1 />`, result[4].Lvl1Value, `501b`)
	}

	lines = `<hb-incr_l_0_by_1 /><hb-tabbedtext />Ναί. 
<hb-incr_l_2_by_1 />
<hb-set_l_1_to_a /><hb-tabbedtext />Ἆρ’ οὖν ἡ τῶνδε τῶν ἀϲκητῶν ἕξιϲ προϲήκουϲ’ ἂν εἴη 
<hb-incr_l_0_by_1 />τούτοιϲ;
<hb-incr_l_0_by_1 />πολλὰϲ μεταβολὰϲ ἐν ταῖϲ ϲτρατείαιϲ μεταβάλλονταϲ ὑδάτων 
<hb-incr_l_1_by_1 />τε καὶ τῶν ἄλλων ϲίτων καὶ εἱλήϲεων καὶ χειμώνων μὴ 
<hb-incr_l_0_by_1 />ἀκροϲφαλεῖϲ εἶναι πρὸϲ ὑγίειαν.
<hb-incr_l_0_by_1 />τιϲ. οἶϲθα γὰρ ὅτι ἐπὶ ϲτρατιᾶϲ ἐν ταῖϲ τῶν ἡρώων 
<hb-incr_l_0_by_1 />ἑϲτιάϲεϲιν οὔτε ἰχθύϲιν αὐτοὺϲ ἑϲτιᾷ, καὶ ταῦτα ἐπὶ 
<hb-incr_l_1_by_1 />θαλάττῃ ἐν Ἑλληϲπόντῳ ὄνταϲ, οὔτε ἑφθοῖϲ κρέαϲιν ἀλλὰ 
<hb-incr_l_0_by_1 />μόνον ὀπτοῖϲ, ἃ δὴ μάλιϲτ’ ἂν εἴη ϲτρατιώταιϲ εὔπορα·
<hb-incr_l_0_by_1 /><hb-tabbedtext />Καὶ ὀρθῶϲ γε, ἔφη, ἴϲαϲί τε καὶ ἀπέχονται. 
<hb-incr_l_1_by_1 /><hb-tabbedtext />Ϲυρακοϲίαν δέ, ὦ φίλε, τράπεζαν καὶ Ϲικελικὴν ποι-
<hb-incr_l_0_by_1 />κιλίαν ὄψου, ὡϲ ἔοικαϲ, οὐκ αἰνεῖϲ, εἴπερ ϲοι ταῦτα δοκεῖ
<hb-incr_l_0_by_1 /><hb-tabbedtext />Ὅλην γὰρ οἶμαι τὴν τοιαύτην ϲίτηϲιν καὶ δίαιταν τῇ 
<hb-incr_l_0_by_1 />μελοποιίᾳ τε καὶ ᾠδῇ τῇ ἐν τῷ παναρμονίῳ καὶ ἐν πᾶϲι 
<hb-incr_l_1_by_1 />ῥυθμοῖϲ πεποιημένῃ ἀπεικάζοντεϲ ὀρθῶϲ ἂν ἀπεικάζοιμεν. 
<hb-incr_l_0_by_1 /><hb-tabbedtext />Πῶϲ γὰρ οὔ; `
	result = AssignCitationValues(lines)
	//for _, r := range result {
	//	fmt.Println(r)
	//}
	if result[16].Lvl1Value != `e` {
		t.Errorf("T: AssignCitationValues(%q)\ngot\n%q\nwant\n%q", `<hb-incr_l_1_by_1 />`, result[16].Lvl1Value, `e`)
	}
}
