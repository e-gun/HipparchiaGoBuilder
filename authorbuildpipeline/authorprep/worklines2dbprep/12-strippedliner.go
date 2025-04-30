//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"regexp"
	"strings"
)

const (
	combininglowerdot = "\u0323"
	straypunct        = `\<\>\{\}\[\]\(\)⟨⟩₍₎\'\.\?\!†⌉⎜͙＃✳※¶§͜﹖→𐄂𝕔;:,‚·‧∣‖⸏Ͼ⊗⸎×⁝≌﹕∙͵﹡*&﹢+"／◦⎫⎬⎭⎩⌈⌊⌋⸐⸑῀﹦~±⎪⌐⁄│□/` // can not purge hyphens until findhyphens happens
	zappedpunct       = `❨❩⟨⟩⟪⟫⦅⦆❴❵`
	homericeditorial  = `※⸖⟩—Ͽ`
	straydigits       = `\d`
	quotish           = `ˈ＇‛‘’“”„ʼ´ʹ«»`
	metrics           = `⏑–⏔⏕⏒⏓⏖̆̅̄̂ʃ`
	buggyhexvals      = ``// actually worth keeping in so you can find problem spots in the parser ( is especially important)
)

var (
	// findmarkup will kill even editorial insertions 'do <ut> des': you have to be sure that no 'pure' angled brackets made it this far
	findmarkup = regexp.MustCompile(`<[^>]*?>`)
	findpunct  = regexp.MustCompile("[" + combininglowerdot + straydigits + straypunct + zappedpunct + homericeditorial + quotish + metrics + "]")

	// deletetcontents: a list of tags whose contents will be removed from the search column: removes things like '<speaker>Th.</speaker>'
	//               ' κρ ' will no longer find every line spoken by Creon in Antigone if you make 'speaker' unsearchable
	//               it is risky to add items to 'unsearchable' without reading/modifying the source for tag rules
	deletetcontents     = buildunsearchable()
	forcelunate         = strings.NewReplacer("σ", "ϲ", "ς", "ϲ")
	findextrawhitespace = regexp.MustCompile(`\s{2,}`)
	strippermap         = buildstrippermap()
)

func buildunsearchable() []*regexp.Regexp {
	const (
		PAT = `<%s>.*?</%s>\s*`
	)
	unsearchablelist := []string{
		"speaker",
	}

	uslice := make([]*regexp.Regexp, len(unsearchablelist))
	for i, item := range unsearchablelist {
		uslice[i] = regexp.MustCompile(fmt.Sprintf(PAT, item, item))
	}
	return uslice
}

func buildstrippermap() map[rune]rune {

	ivals := []rune("ἀἁἂἃἄἅἆἇᾀᾁᾂᾃᾄᾅᾆᾇᾲᾳᾴᾶᾷᾰᾱὰάἐἑἒἓἔἕὲέἰἱἲἳἴἵἶἷὶίῐῑῒΐῖῗΐὀὁὂὃὄὅόὸὐὑὒὓὔὕὖὗϋῠῡῢΰῦῧύὺᾐᾑᾒᾓᾔᾕᾖᾗῂῃῄῆῇἤἢἥἣὴήἠἡἦἧὠὡὢὣὤὥὦὧᾠᾡᾢᾣᾤᾥᾦᾧῲῳῴῶῷώὼ")
	ovals := []rune("αααααααααααααααααααααααααεεεεεεεειιιιιιιιιιιιιιιιιοοοοοοοουυυυυυυυυυυυυυυυυηηηηηηηηηηηηηηηηηηηηηηηωωωωωωωωωωωωωωωωωωωωωωω")

	ivals = append(ivals, []rune("ᾈᾉᾊᾋᾌᾍᾎᾏἈἉἊἋἌἍἎἏΑἘἙἚἛἜἝΕἸἹἺἻἼἽἾἿΙὈὉὊὋὌὍΟὙὛὝὟΥᾘᾙᾚᾛᾜᾝᾞᾟἨἩἪἫἬἭἮἯΗᾨᾩᾪᾫᾬᾭᾮᾯὨὩὪὫὬὭὮὯῤῥῬΒΨΔΦΓΞΚΛΜΝΠϘΡσΣςϹΤΧΘΖ")...)
	ovals = append(ovals, []rune("αααααααααααααααααεεεεεεειιιιιιιιιοοοοοοουυυυυηηηηηηηηηηηηηηηηηωωωωωωωωωωωωωωωωρρρβψδφγξκλμνπϙρϲϲϲϲτχθζ")...)

	ivals = append(ivals, []rune("vUjÁÄáäÉËéëÍÏíïÓÖóöÜÚüú")...)
	ovals = append(ovals, []rune("uViaaaaeeeeiiiioooouuuu")...)

	stripmap := make(map[rune]rune)
	for i, iv := range ivals {
		stripmap[iv] = ovals[i]
	}

	return stripmap
}

func stripstring(in string) string {
	//	sample in:
	//		ἔντοϲ ἀμώμητον κάλλιπον οὐκ ἐθέλων
	//
	//	sample out:
	//		εντοϲ αμωμητον καλλιπον ουκ εθελων
	rr := []rune(in)
	for i, r := range rr {
		found, ok := strippermap[r]
		if !ok {
			found = r
		}
		rr[i] = found
	}
	return string(rr)
}

func BuildStrippedAndAccentedlines(lines []structs.DbWorkline) []structs.DbWorkline {
	// 	generate the easy to search stripped column
	//
	// findpunct / straypunct is a big deal: it defines what a clean line will look like and so what you can search for
	//   sadly can't nuke :punct: as a class because we need hyphens
	//   if you want to find »αʹ« you need ʹ
	//   if you want to find »͵α« you need ͵
	//   if you want to search for undocumented/idiosyncratic chars you need ◦⊚
	//   misc other things that one might want to exclude but are currently included: ☩ͻ
	//   the following are supposed to be killed off by bracketsimplifier(): ❨❩⟨⟩⟪⟫⦅⦆❴❵
	//   no longer relevant?: ⸨⸩｟｠《
	//	 homeric editorial marks: ※⸖⟩—Ͽ

	for i, line := range lines {
		toclean := line.MarkedUp
		for _, fp := range deletetcontents {
			toclean = fp.ReplaceAllString(toclean, "")
		}
		toclean = strings.ToLower(toclean)
		toclean = findmarkup.ReplaceAllString(toclean, "")
		toclean = findpunct.ReplaceAllString(toclean, "")
		toclean = forcelunate.Replace(toclean)
		toclean = findextrawhitespace.ReplaceAllString(toclean, " ")
		lines[i].Accented = toclean
		lines[i].Stripped = stripstring(toclean)
	}
	return lines
}

// zBuildStrippedAndAccentedlines - leave in the capitalization (for Helma)
func zBuildStrippedAndAccentedlines(lines []structs.DbWorkline) []structs.DbWorkline {
	// 	generate the easy to search stripped column
	//
	// findpunct / straypunct is a big deal: it defines what a clean line will look like and so what you can search for
	//   sadly can't nuke :punct: as a class because we need hyphens
	//   if you want to find »αʹ« you need ʹ
	//   if you want to find »͵α« you need ͵
	//   if you want to search for undocumented/idiosyncratic chars you need ◦⊚
	//   misc other things that one might want to exclude but are currently included: ☩ͻ
	//   the following are supposed to be killed off by bracketsimplifier(): ❨❩⟨⟩⟪⟫⦅⦆❴❵
	//   no longer relevant?: ⸨⸩｟｠《
	//	 homeric editorial marks: ※⸖⟩—Ͽ

	for i, line := range lines {
		toclean := line.MarkedUp
		for _, fp := range deletetcontents {
			toclean = fp.ReplaceAllString(toclean, "")
		}
		// toclean = strings.ToLower(toclean)
		toclean = findmarkup.ReplaceAllString(toclean, "")
		toclean = findpunct.ReplaceAllString(toclean, "")
		// toclean = forcelunate.Replace(toclean)
		toclean = findextrawhitespace.ReplaceAllString(toclean, " ")
		lines[i].Accented = toclean
		lines[i].Stripped = stripstring(toclean)
	}
	return lines
}
