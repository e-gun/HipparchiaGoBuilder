//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package generic

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

const (
	TERMINATIONS = `(\s|\.|\]|\<|⟩|\)|’|”|\!|,|:|;|\?|⸥|«|·|$)` // circular imports means this is declared 2x... see also "vv.constants.co"
	GLC          = `ΐάέήίΰαβγδεζηθικλμνξοπρτυφχψωϊϋόύώϝϲἀἁἂἃἄἅἆἇἐἑἒἓἔἕἠἡἢἣἤἥἦἧἰἱἲἳἴἵἶἷὀὁὂὃὄὅὐὑὒὓὔὕὖὗὠὡὢὣὤὥὦὧὰὲὴὶὸὺὼᾀᾁᾂᾃᾄᾅᾆᾇᾐᾑᾒᾓᾔᾕᾖᾗᾠᾡᾢᾣᾤᾥᾦᾧᾲᾳᾴᾶᾷῂῃῄῆῇῒῖῢῤῥῦῧῲῳῴῶῷ`
	GUC          = `ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΤΥΦΧΨΩϜϹἈἉἊἋἌἍἎἏἘἙἚἛἜἝἨἩἪἫἬἭἮἯἸἹἺἻἼἽἾἿὈὉὊὋὌὍὙὛὝὟὨὩὪὫὬὭὮὯᾊᾋᾌᾍᾎᾏᾚᾛᾜᾝᾞᾟᾪᾫᾬᾭᾮᾯᾼῌῬῼ⒣`
	QUIKLATIN    = `a-zA-Z`
	ALLGKANDLAT  = GLC + GUC + QUIKLATIN
)

var (
	// runefd etc. to avoid looping this in hot code
	runefd         = getrunefeeder()
	runereduce     = getrunereducer()
	LunateSwap     = regexp.MustCompile("σ" + TERMINATIONS)
	IsGreekLC      = regexp.MustCompile(`[` + GLC + `]`)
	IsGreekUC      = regexp.MustCompile(`[` + GUC + `]`)
	IsGreek        = regexp.MustCompile(`[` + GUC + GLC + `]`)
	IsGreekOrLatin = regexp.MustCompile(`[` + ALLGKANDLAT + `]`)
)

//
// THE EXPORTABLE FUNCTIONS
//

// StripaccentsSTR - ὀκνεῖϲ --> οκνειϲ, etc.
func StripaccentsSTR(u string) string {
	// reducer := getrunereducer()
	ru := []rune(u)
	stripped := make([]rune, len(ru))
	for i, x := range ru {
		stripped[i] = runereduce[x]
	}
	s := string(stripped)
	return s
}

var (
	a4gswap = strings.NewReplacer("ὰ", "ά", "ὲ", "έ", "ὶ", "ί", "ὸ", "ό", "ὺ", "ύ", "ὴ", "ή", "ὼ", "ώ",
		"ἂ", "ἄ", "ἃ", "ἅ", "ᾲ", "ᾴ", "ᾂ", "ᾄ", "ᾃ", "ᾅ", "ἒ", "ἔ", "ἲ", "ἴ", "ὂ", "ὄ", "ὃ", "ὅ", "ὒ", "ὔ", "ὓ", "ὕ",
		"ἢ", "ἤ", "ἣ", "ἥ", "ᾓ", "ᾕ", "ᾒ", "ᾔ", "ὢ", "ὤ", "ὣ", "ὥ", "ᾣ", "ᾥ", "ᾢ", "ᾤ", "á", "a", "é", "e",
		"í", "i", "ó", "o", "ú", "u")
)

// SwapAcuteForGrave - ὰ --> ά
func SwapAcuteForGrave(thetext string) string {
	return a4gswap.Replace(thetext)
}

type PolytonicSorterStruct struct {
	Sortstring     string
	Originalstring string
	Count          int // Count is needed by RtIndexMaker()
}

// PolytonicSort - sort a slice of polytonic words; note that this is going to be relatively costly to execute
func PolytonicSort(pt []string) []string {
	ss := make([]PolytonicSorterStruct, len(pt))
	for i := 0; i < len(pt); i++ {
		ss[i] = PolytonicSorterStruct{
			Sortstring:     strings.Replace(StripaccentsSTR(pt[i]), "ϲ", "σ", -1) + pt[i],
			Originalstring: pt[i],
			Count:          0,
		}
	}

	slices.SortFunc(ss, func(a, b PolytonicSorterStruct) int { return cmp.Compare(a.Sortstring, b.Sortstring) })

	for i := 0; i < len(pt); i++ {
		pt[i] = ss[i].Originalstring
	}

	return pt
}

//
// THE HELPERS/FEEDERS
//

func getrunereducer() map[rune]rune {
	// because we don't have access to python's transtable function
	// runefd := getrunefeeder()
	// runefd now a var at top of file

	reducer := make(map[rune]rune)
	for f, _ := range runefd {
		for _, r := range runefd[f] {
			reducer[r] = f
		}
	}
	return reducer
}

// getrunefeeder - this one will de-capitalize and de-accentuate (needed for various strippers)
func getrunefeeder() map[rune][]rune {
	feeder := make(map[rune][]rune)
	feeder['α'] = []rune("αἀἁἂἃἄἅἆἇᾀᾁᾂᾃᾄᾅᾆᾇᾲᾳᾴᾶᾷᾰᾱὰάᾈᾉᾊᾋᾌᾍᾎᾏἈἉἊἋἌἍἎἏΑ")
	feeder['ε'] = []rune("εἐἑἒἓἔἕὲέἘἙἚἛἜἝΕ")
	feeder['ι'] = []rune("ιἰἱἲἳἴἵἶἷὶίῐῑῒΐῖῗΐἸἹἺἻἼἽἾἿΙ")
	feeder['ο'] = []rune("οὀὁὂὃὄὅόὸὈὉὊὋὌὍΟ")
	feeder['υ'] = []rune("υὐὑὒὓὔὕὖὗϋῠῡῢΰῦῧύὺὙὛὝὟΥ")
	feeder['η'] = []rune("ηᾐᾑᾒᾓᾔᾕᾖᾗῂῃῄῆῇἤἢἥἣὴήἠἡἦἧᾘᾙᾚᾛᾜᾝᾞᾟἨἩἪἫἬἭἮἯΗ")
	feeder['ω'] = []rune("ωὠὡὢὣὤὥὦὧᾠᾡᾢᾣᾤᾥᾦᾧῲῳῴῶῷώὼᾨᾩᾪᾫᾬᾭᾮᾯὨὩὪὫὬὭὮὯ")
	feeder['ρ'] = []rune("ρῤῥῬ")
	feeder['β'] = []rune("βΒ")
	feeder['ψ'] = []rune("ψΨ")
	feeder['δ'] = []rune("δΔ")
	feeder['φ'] = []rune("φΦ")
	feeder['γ'] = []rune("γΓ")
	feeder['ξ'] = []rune("ξΞ")
	feeder['κ'] = []rune("κΚ")
	feeder['λ'] = []rune("λΛ")
	feeder['μ'] = []rune("μΜ")
	feeder['ν'] = []rune("νΝ")
	feeder['π'] = []rune("πΠ")
	feeder['ϙ'] = []rune("ϙϘ")
	feeder['ϲ'] = []rune("ϲσΣςϹ")
	feeder['τ'] = []rune("τΤ")
	feeder['χ'] = []rune("χΧ")
	feeder['θ'] = []rune("θΘ")
	feeder['ζ'] = []rune("ζΖ")
	feeder['a'] = []rune("aAÁÄáäă")
	feeder['b'] = []rune("bB")
	feeder['c'] = []rune("cC")
	feeder['d'] = []rune("dD")
	feeder['e'] = []rune("eEÉËéëāĕē")
	feeder['f'] = []rune("fF")
	feeder['g'] = []rune("gG")
	feeder['h'] = []rune("hH")
	feeder['i'] = []rune("iIÍÏíïJj")
	feeder['k'] = []rune("kK")
	feeder['l'] = []rune("lL")
	feeder['m'] = []rune("mM")
	feeder['n'] = []rune("nN")
	feeder['o'] = []rune("oOÓÖóöŏō")
	feeder['p'] = []rune("pP")
	feeder['q'] = []rune("qQ")
	feeder['r'] = []rune("rR")
	feeder['s'] = []rune("sS")
	feeder['t'] = []rune("tT")
	feeder['u'] = []rune("uUvVÜÚüú")
	feeder['w'] = []rune("wW")
	feeder['x'] = []rune("xX")
	feeder['y'] = []rune("yY")
	feeder['z'] = []rune("zZ")
	return feeder
}

// extendedrunefeeder - this one will do acute for grave (needed for lemma highlighting)
func extendedrunefeeder() map[rune][]rune {
	feeder := getrunefeeder()
	feeder['ά'] = []rune("ὰά")
	feeder['έ'] = []rune("ὲέ")
	feeder['ή'] = []rune("ὴή")
	feeder['ί'] = []rune("ὶί")
	feeder['ό'] = []rune("όὸ")
	feeder['ύ'] = []rune("ύὺ")
	feeder['ώ'] = []rune("ώὼ")
	feeder['ἂ'] = []rune("ἂἄ")
	feeder['ἒ'] = []rune("ἒἔ")
	feeder['ἢ'] = []rune("ἢἤ")
	feeder['ἲ'] = []rune("ἲἴ")
	feeder['ὂ'] = []rune("ὂὄ")
	feeder['ὒ'] = []rune("ὒὔ")
	feeder['ὓ'] = []rune("ὓὕ")
	feeder['ὢ'] = []rune("ὢὤ")
	feeder['ὣ'] = []rune("ὣὥ")
	feeder['ἃ'] = []rune("ἅἃ")
	feeder['ᾲ'] = []rune("ᾲᾴ")
	feeder['ᾂ'] = []rune("ᾂᾄ")
	feeder['ἣ'] = []rune("ἣἥ")
	feeder['ᾒ'] = []rune("ᾒᾔ")
	feeder['ᾓ'] = []rune("ᾓᾕ")
	feeder['ὃ'] = []rune("ὃὅ")
	feeder['ὂ'] = []rune("ὂὄ")
	feeder['ὒ'] = []rune("ὒὔ")
	feeder['ᾂ'] = []rune("ᾂᾄ")
	feeder['ᾃ'] = []rune("ᾃᾅ")
	feeder['ᾢ'] = []rune("ᾢᾤ")
	feeder['ᾣ'] = []rune("ᾣᾥ")
	return feeder
}

// uvσςϲreducer - provide map to UVσςϲ
func uvσςϲreducer() map[rune]rune {
	// map[73:105 74:105 85:117 86:117 105:105 106:105 ...]
	feeder := make(map[rune][]rune)

	feeder['u'] = []rune("uUvVÜÚüú")
	feeder['ϲ'] = []rune("ϲσΣςϹ")
	feeder['i'] = []rune("iIÍÏíïJj")

	reducer := make(map[rune]rune)
	for f, _ := range feeder {
		for _, r := range feeder[f] {
			reducer[r] = f
		}
	}
	return reducer
}

// uvcapsreducer - provide map to UVcaps
func uvcapsreducer() map[rune]rune {
	// map[73:105 74:105 85:117 86:117 105:105 106:105 ...]
	feeder := make(map[rune][]rune)

	feeder['u'] = []rune("uv")
	feeder['i'] = []rune("ij")
	feeder['U'] = []rune("ÜÚ")
	feeder['V'] = []rune("U")
	feeder['I'] = []rune("IJ")

	reducer := make(map[rune]rune)
	for f, _ := range feeder {
		for _, r := range feeder[f] {
			reducer[r] = f
		}
	}
	return reducer
}
