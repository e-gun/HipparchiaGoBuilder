//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"slices"
	"strings"
	"testing"
)

func TestNybbler(t *testing.T) {
	ttc := `01`
	result := getbase16(ttc)
	textlevel := (result & 0x70) >> 4
	action := result & 0x0F
	fmt.Println(textlevel, action)
}

func TestFixIrrationalTags(t *testing.T) {
	invalidline := `<hb-sp-expanded_text><hb-fs-g-smallerthannormal>τίϲ ἡ τάραξιϲ</hb-sp-expanded_text> τοῦ βίου; τί βάρβιτοϲ</hb-fs-g-smallerthannormal>`
	result := FixIrrationalTags(invalidline)
	want := `<hb-sp-expanded_text><hb-fs-g-smallerthannormal>τίϲ ἡ τάραξιϲ</hb-fs-g-smallerthannormal></hb-sp-expanded_text><hb-fs-g-smallerthannormal> τοῦ βίου; τί βάρβιτοϲ</hb-fs-g-smallerthannormal>`
	if result != want {
		t.Errorf("fixhmuirrationaloragnization(%q)\ngot\n%q\nwant\n%q", invalidline, result, want)
	}
	invalidline = `[]κακ<hb-sp-superscript>η</hb-sp-superscript> β<hb-sp-x>ο<hb-sp-y>υab</hb-sp-x>c<hb-sp-superscript>λ</hb-sp-y></hb-sp-superscript>`
	result = FixIrrationalTags(invalidline)
	want = `[]κακ<hb-sp-superscript>η</hb-sp-superscript> β<hb-sp-x>ο<hb-sp-y>υab</hb-sp-y></hb-sp-x><hb-sp-y>c<hb-sp-superscript>λ</hb-sp-superscript></hb-sp-y><hb-sp-superscript></hb-sp-superscript>`
	if result != want {
		t.Errorf("fixhmuirrationaloragnization(%q)\ngot\n%q\nwant\n%q", invalidline, result, want)
	}
}

func TestBracketSimplification(t *testing.T) {
	ttc := `❨❩❴❵⟦⟧⟪⟫《`
	result := findbrackets.ReplaceAllStringFunc(ttc, bracketsimplifier)
	want := `(){}[]⟪⟫⟨`
	if result != want {
		t.Errorf("swapitemordersuite(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestSwapitemordersuite(t *testing.T) {
	ttc := `<span class="latin smallerthannormal">Gnom. Vatic. 743 [</span>zzz`
	result := swapitemordersuite(ttc)
	want := `<span class="latin smallerthannormal">Gnom. Vatic. 743 </span>[zzz`
	if result != want {
		t.Errorf("swapitemordersuite(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestSmartersinglequotes(t *testing.T) {
	ttc := ` ’open’ for close ’zzz’`
	result := smartersinglequotes(ttc)
	want := ` ‘open’ for close ‘zzz’`
	if result != want {
		t.Errorf("smartersinglequotes(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestBalanceQuotes(t *testing.T) {
	ttc := `	<hb-incr_l_0_by_1 />Ϲμινθεῦ εἴ ποτέ τοι χαρίεντ＇ ἐπὶ νηὸν ἔρεψα,
	<hb-incr_l_0_by_1 />ἢ εἰ δή ποτέ τοι κατὰ πίονα μηρ ’ίἔκηα
	<hb-incr_l_0_by_1 />ταύρων ἠδ＇ αἰγῶν, τὸ δέ μοι κρήηνον ἐέλδωρ·`
	result := BalanceQuotes(ttc)
	fmt.Println(result)
}

func TestNyb15(t *testing.T) {
	thehex := highunicodetohex(`ⓕⓐ`)
	hexsegments := strings.Split(thehex, "█")
	slices.Reverse(hexsegments)
	hs := structs.NewHexStack(hexsegments)
	s, hs := nyb15(hs)
	fmt.Println("s", s)

}

func TestHexRunner(t *testing.T) {
	// ttc := `. <hb-nb-endofpage /> █ⓕⓔ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █ⓔⓕ █⑧⓪ █ⓑ③ █ⓑ⓪ █ⓑ② █ⓑ⑦ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ③ █ⓕⓕ █ⓐ⑧ █⑧ⓔ █⑨ⓑ █⑧② █ⓔ⑧ █⑧⑨ █⑧⑥ █ⓕ④ <hb-title>ΙΩΑΝΝΟΥ ΤΟΥ ΔΟΞΑΠΑΤΡΗ █⑧⑨ █⑧⑦ █ⓕ④ ΕΙϹ ΤΟ █⑧⑨ █⑧⑧ █ⓕ④ ΠΕΡΙ ΕΥΡΕϹΕΩϹ ΕΡΜΟΓΕΝΟΥϹ ΒΙΒΛΙΟΝ</hb-title> █⑧⑧ █⑧⑨ <hb-tabbedtext />`
	//hexsequence := `█ⓑ⓪ █ⓑ⓪ █ⓑ① █ⓑ⑤ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ① █ⓑ② █ⓕⓕ █ⓔⓕ █⑧② █ⓒⓓ █ⓔ① █ⓔ④ █ⓔ① █ⓔ② █ⓔ① █ⓕⓕ █ⓔⓕ █⑧③ █ⓓ③ █ⓕ⑨ █ⓕ② █ⓐⓒ █ⓐ⓪ █ⓒ① █ⓕ② █ⓔ① █ⓔ② █ⓐⓒ █ⓐ⓪ █ⓓ⓪ █ⓔ① █ⓔⓒ █ⓕⓕ █⑨ⓐ █ⓕ② █ⓐ⑧ █ⓑ① █ⓑ① █ⓑ⓪ █ⓐ⑨ █ⓕⓕ █ⓔⓕ █ⓔ① █ⓒ① █ⓕ② █ⓔ① █ⓔ② █ⓔ⑨ █ⓔ① █ⓕⓕ █ⓔⓕ █ⓔ② █ⓒⓓ █ⓔ① █ⓔ④ █ⓔ① █ⓔ② █ⓔ① █ⓕⓕ █ⓔⓕ █ⓔ③ █ⓕⓕ █ⓔⓕ █ⓔ④ █ⓔ③ █ⓐ⓪ █ⓑ⑤ █ⓑ⑥ █ⓑ⓪ █ⓐⓓ █ⓑ⑤ █ⓑ⑥ █ⓑ⑤ █ⓐ⓪ █ⓔ① █ⓔ③ █ⓑⓕ █ⓕⓕ █ⓔⓕ █ⓔ⑤ █ⓒ⑨ █ⓒ⑦ █ⓒⓒ █ⓓ③ █ⓕ⑨ █ⓕ② █ⓐ⓪ █ⓑ② █ⓑ① █ⓑⓐ █ⓑ② █ⓐⓔ █ⓑ① █ⓑ⑤ █ⓑ③ █ⓕⓕ █ⓔ① █ⓕⓐ`
	// hexsequence := ` █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ① █ⓑ⑥ █ⓕⓕ █⑨① <hb-tabbedtext />Ἔχ̣[ο]μ̣ε̣ν̣ περὶ θεοῦ διάλημψιν καὶ ἀπὸ τ̣ῆϲ γραφῆϲ καὶ █⑧⓪ τῆϲ κ[ο]ινῆϲ ἐννοίαϲ ὅτι ἄτρεπτόϲ ἐϲτιν, │ ὅ[τι ἀ]ναλλοί-█⑧② ωτόϲ ἐϲτιν· ὁ γὰρ ὅλωϲ μὴ ὑποκείμενοϲ ποιότη̣[τι οὐ τρέ-█⑧⓪ πετ]αι, οὐκ ἀλλοιοῦται· οὐδὲν │ γ̣ὰ̣ρ̣ ἕτερ̣ό̣ν̣ ἐϲτιν ἀλλοί-█⑧③ ωϲιϲ ἢ κατὰ ποιὸν μεταβολή. οὐ [πᾶϲα μεταβο]λ̣ὴ ἀλλοίωϲίϲ █⑧⓪ ἐϲτιν, │ ἀλλ＇ ἡ κ̣[ατὰ] π̣οιότητα. εἰ[ϲ]ίν γ̣ε̣ κα[ὶ ἄλλαι] █⑧④ μεταβολαί, ἐπε[ὶ καὶ κινήϲειϲ εἰϲ]ὶν ἄλλαι. τὸ γ̣ι̣νόμενον │ █⑧⓪ μεταβάλλει, ἀλλ＇ οὔκ ἐϲτιν ἀ̣λ̣λο[ί]ω̣[ϲ]ι̣ϲ ἡ κί[ν]η̣[ϲ]ιϲ █⑧⓪ αὕτη. τὸ αὐξό̣μ̣[ενον μεταβάλλ]ει, ἀλλ＇ οὔκ ἐϲτιν ἀλλοί│ω-█⑧⑥ ϲιϲ αὕτη· προϲθήκη γὰρ καὶ αὔξηϲιϲ̣ π̣[ο]ϲ̣[ο]ῦ̣ ἐϲτιν ἡ τοι-█⑧⓪ αύτη κ[ί]ν̣η̣[ϲ]ιϲ. ὅ̣[τ]αν δὲ ἐκ φα̣ύ̣λ̣ο̣υ ϲ[πο]υ̣│δα̣ῖ̣οϲ̣ ἢ̣ ἐ̣κ̣ █⑧⑦ ϲπουδαίου φαῦλοϲ γένητ̣[α]ί τ̣ι̣ϲ̣, ἠ̣λλοίωται κατ̣[ὰ] τ̣[ὴν █⑧⓪ ποιότητα], ὡϲ αὖ ὅτε ἐκ νοϲο̣ῦν│τοϲ εἰϲ ὑγείαν ἔλθῃ καὶ █⑧⑧ █⑧⑧ ἔνπαλιν. █⑧⑧ █⑧⑧ <hb-tabbedtext />ἀ[λ]λ̣ὰ̣ τ[ὰ]ϲ̣ [λ]έ̣[ξ]ειϲ πρὸϲ τὴν ἔν[νοιαν τοῦ πρ]ά̣γ-█⑧⓪ ματοϲ, π̣ερὶ οὗ │ λέγονται, δ̣ι̣[ανοού]μεθα. ὁ θεὸϲ οὐκ ἐκ `
	// hexsequence := ` █⑧⓪ Domus referta vasis Corinthiis et Deliacis, in quibus est █⑧⓪ authepsa illa quam tanto pretio nuper mercatus est ut qui █⑧⓪ praetereuntes quid praeco enumeraret audiebant fundum █⑧⓪ venire arbitrarentur. Quid praeterea caelati argenti, quid █⑧⓪ stragulae vestis, quid pictarum tabularum, quid signorum, █⑧⓪ quid marmoris apud illum putatis esse? Tantum scilicet █⑧⓪ quantum e multis splendidisque familiis in turba et rapinis █⑧⓪ coacervari una in domo potuit. Familiam vero quantam et █⑨⓪ quam variis cum artificiis habeat quid ego dicam? Mitto █⑧⓪ hasce artis volgaris, coquos, pistores, lecticarios; animi et █⑧⓪ aurium causa tot homines habet ut cotidiano cantu vocum █⑧⓪ et nervorum et tibiarum nocturnisque conviviis tota vicinitas █⑧⓪ personet. In hac vita, iudices, quos sumptus cotidianos, █⑧⓪ quas effusiones fieri putatis, quae vero convivia? honesta, █⑧⓪ credo, in eius modi domo, si domus haec habenda est potius █⑧⓪ ⟨quam⟩ officina nequitiae ac deversorium flagitiorum omnium. █⑨⓪ Ipse vero quem ad modum composito et dilibuto capillo █⑧⓪ passim per forum volitet cum magna caterva togatorum vi-█⑧⓪ detis, iudices; videtis ut omnis despiciat, ut hominem prae █ⓕⓔ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █ⓔⓕ █⑧⓪ █ⓑ⓪ █ⓑ④ █ⓑ⑦ █ⓑ④ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ② █ⓕⓕ █ⓔⓕ █⑧② █ⓓ③ █ⓓ② █ⓔⓕ █ⓕ③ █ⓔ③ █ⓕⓕ █ⓔⓕ █⑧③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓕⓕ █⑨ⓑ █⑧① █⑧⑦ █⑧④ █ⓔⓕ █ⓔ③ █ⓕⓕ █ⓔ② █ⓕⓐ se neminem putet, ut se solum beatum, solum potentem █⑧⓪ putet. Quae vero efficiat et quae conetur si velim comme-█⑧⓪ morare, vereor, iudices, ne quis imperitior existimet me █⑧⓪ causam nobilitatis victoriamque voluisse laedere. Tametsi █⑧⓪ meo iure possum, si quid in hac parte mihi non placeat, █⑧⓪ vituperare; non enim vereor ne quis alienum me animum █⑧⓪ habuisse a causa nobilitatis existimet. █⑨⓪ <hb-tabbedtext /><hb-tabbedtext />Sciunt ei qui me norunt me pro mea tenui infirmaque █⑧⓪ parte, postea quam id quod maxime volui fieri non potuit, ut █⑧⓪ componeretur, id maxime defendisse ut ei vincerent qui █⑧⓪ vicerunt. Quis enim erat qui non videret humilitatem cum █⑧⓪ dignitate de amplitudine contendere? quo in certamine █⑧⓪ perditi civis erat non se ad eos iungere quibus incolumibus █⑧⓪ et domi dignitas et foris auctoritas retineretur. Quae per-█⑧⓪ fecta esse et suum cuique honorem et gradum redditum █⑧⓪ gaudeo, iudices, vehementerque laetor eaque omnia deorum █⑧⓪ voluntate, studio populi Romani, consilio et imperio et █⑨⓪ felicitate L. Sullae gesta esse intellego. Quod animadversum █⑧⓪ est in eos qui contra omni ratione pugnarunt, non debeo █⑧⓪ reprehendere; quod viris fortibus quorum opera eximia in █⑧⓪ rebus gerendis exstitit honos habitus est, laudo. Quae ut █⑧⓪ fierent idcirco pugnatum esse arbitror meque in eo studio █⑧⓪ partium fuisse confiteor. Sin autem id actum est et idcirco █⑧⓪ arma sumpta sunt ut homines postremi pecuniis alienis █⑧⓪ locupletarentur et in fortunas unius cuiusque impetum █⑧⓪ facerent, et id non modo re prohibere non licet sed ne █⑧⓪ verbis quidem vituperare, tum vero in isto bello non re-█⑧⓪ creatus neque restitutus sed subactus oppressusque populus █⑨⓪ Romanus est. Verum longe aliter est; nil horum est, █⑧⓪ iudices. Non modo non laedetur causa nobilitatis, si istis █⑧⓪ hominibus resistetis, verum etiam ornabitur. Etenim qui █⑧⓪ haec vituperare volunt Chrysogonum tantum posse querun-█⑧⓪ tur; qui laudare volunt concessum ei non esse commemo-█⑧⓪ rant. Ac iam nihil est quod quisquam aut tam stultus aut █⑧⓪ tam improbus sit qui dicat: ＇Vellem quidem liceret; █⑧⓪ hoc dixissem.＇ Dicas licet. ＇Hoc fecissem.＇ Facias licet; █⑧⓪ nemo prohibet. ＇Hoc decrevissem.＇ Decerne, modo █⑧⓪ recte; omnes approbabunt. ＇Hoc iudicassem.＇ Laudabunt █⑨⓪ omnes, si recte et ordine iudicaris. Dum necesse erat █⑧⓪ resque ipsa cogebat, unus omnia poterat; qui postea quam `
	// hexsequence := `█⑨ⓑ █⑧① █⑧⑦ se neminem putet,`
	// hexsequence := `ut ita █ⓕⓔ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █ⓔⓕ █⑧⓪ █ⓑ⓪ █ⓑ④ █ⓑ⑦ █ⓑ④ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⑤ █ⓑ⑦ █ⓕⓕ █ⓔⓕ █⑧② █ⓒ① █ⓕ④ █ⓕ④ █ⓕⓕ █ⓔⓕ █⑧③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓕⓕ █ⓑ① █ⓐ⑧ █⑨① █⑨④ █⑧⑦ █ⓔ⑧ █ⓕⓐ █⑨① █ⓔⓕ █ⓔⓒ █ⓓ② █ⓔⓕ █ⓔⓓ █ⓔ① █ⓔ⑤ █ⓕⓕ █ⓔⓕ █ⓔ④ █ⓒⓔ █ⓔⓕ █ⓔⓔ █ⓐⓔ █ⓐ⓪ █ⓒ④ █ⓔ⑤ █ⓔ③ █ⓐⓔ █ⓐ⓪ █ⓑ⑥ █ⓑ① █ⓐ⓪ █ⓐ⑧ █ⓐ⑤ █ⓑ① █ⓑ④ █ⓔ⓪ █ⓑ① █ⓑ① █ⓐ⑨ █ⓕⓕ dicam, mollitiamque naturae `
	hexsequence := ` hos retine atque auge famam laudesque bonorum. █⑨⓪ █⑧ⓕ █ⓕ④ █ⓕⓕ █ⓔⓕ █ⓔ③ █ⓕⓕ <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-title>EX INCERTIS LIBRIS</hb-title> █⑧① █ⓔⓕ █ⓔ③ █ⓒⓔ █ⓔⓕ █ⓔⓔ █ⓐⓔ █ⓐ⓪ █ⓑ② █ⓑ⓪ █ⓑ② █ⓒⓓ █ⓕⓕ █ⓔ⑧ █ⓕⓐ █⑧⑧ quorum luxuries fortunas, censa peredit █⑨⓪ █ⓔⓕ █ⓔ③ █ⓓⓑ █ⓓ⓪ █ⓕ② █ⓔⓕ █ⓔ② █ⓐⓔ █ⓓⓓ █ⓐ⓪ █ⓐ⑥ █ⓑ③ █ⓕ⑤ █ⓔⓒ █ⓕ④ █ⓐⓔ █ⓐ⓪ █ⓕ③ █ⓕ⑨ █ⓔⓒ █ⓔⓒ █ⓐⓔ █ⓐ⓪ █ⓒ⑦ █ⓒⓒ █ⓐ⑥ █ⓐ⓪ █ⓑ④ █ⓐⓔ █ⓑ② █ⓑ④ █ⓑ⑧ █ⓒⓑ █ⓕⓕ █ⓔ⑧ █ⓕⓐ █⑧⑨ nam quasi vos sibi dedecori genuere parentes █⑨⓪ █ⓔⓕ █ⓔ③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓐⓔ █ⓐ⓪ █ⓐ⑥ █ⓑ③ █ⓒⓕ █ⓔ⑥ █ⓔ⑥ █ⓐⓔ █ⓐ⑥ █ⓐ⓪ █ⓑ① █ⓐⓔ █ⓑ⑦ █ⓑ⑦ █ⓕⓕ █ⓔ⑧ █ⓕⓐ █⑧ⓐ cedant arma togae, concedat laurea laudi █⑨⓪ █ⓔⓕ █ⓔ③ █ⓓⓑ █ⓓ③ █ⓔ① █ⓔⓒ █ⓐⓔ █ⓓⓓ █ⓐ⓪ █ⓐ⑥ █ⓑ③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓐⓔ █ⓐ⑥ █ⓐ⓪ █ⓑ⑤ █ⓕⓕ █ⓔ⑧ █ⓕⓐ █⑧ⓑ O fortunatam natam me consule Romam! <hb-nb-endofpage /> █⑨⓪ █ⓔⓕ █ⓔ③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓐⓔ █ⓐ⓪ █ⓐ⑥ █ⓑ③ █ⓒ① █ⓕ④ █ⓕ④ █ⓐ⑥ █ⓐⓔ █ⓐ⓪ █ⓑ② █ⓐⓔ █ⓑ① █ⓑ⑤ █ⓐⓔ █ⓑ③ █ⓕⓕ █ⓔ⑧ █ⓕⓐ █⑧ⓒ in montes patrios et ad incunabula nostra <hb-nb-endofpage /> █⑨⑧ █⑨③ █⑧ⓕ █ⓕ④ █ⓕⓕ █ⓔⓕ █ⓔ③ █ⓕⓕ <hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-tabbedtext /><hb-title>MARIVS</hb-title> █⑧① █ⓔⓕ █ⓔ③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓐⓔ █ⓐ⓪ █ⓐ⑥ █ⓑ③ █ⓒⓒ █ⓔ⑤ █ⓔ⑦ █ⓐⓔ █ⓐ⑥ █ⓐ⓪ █ⓑ① █ⓐⓔ █ⓑ② █ⓕⓕ █ⓔ⑧ █ⓕⓐ █⑧ⓓ nuntia fulva Iovis, miranda visa figura █⑨⓪ █ⓔⓕ `

	result := HexRunner(hexsequence)
	//want = strings.ReplaceAll(want, `\`, ``)
	//result = strings.ReplaceAll(result, `\`, "")
	// fails because of the representation of newlines...
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		fmt.Printf("[%d]\t%s\n", i, line)
	}
}
