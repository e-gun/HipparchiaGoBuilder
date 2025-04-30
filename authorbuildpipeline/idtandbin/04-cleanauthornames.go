//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/grk"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"regexp"
	"slices"
	"strings"
)

var (
	findgreek     = regexp.MustCompile(`\$\d?([^&]*?)&\d?`)
	findlatin     = regexp.MustCompile(`&1([^&]*?)&`)
	namemaker1    = regexp.MustCompile(`(Pseudo-|)<hb-fs-l-bold>([^<]*?)</hb-fs-l-bold> ([^‹]*?)$`) // note ">" and not "›"
	namemaker2    = regexp.MustCompile(`(Pseudo-|)<hb-fs-l-bold>([^<]*?)</hb-fs-l-bold>`)
	shtrnamemaker = regexp.MustCompile(`(.*?) \(.*?\)$`)
)

func CleanDBAuthorNames(au *structs.DbAuthor) {
	// `&1Menippus& Phil.` --> `<hb-fs-l-bold>Menippus</hb-fs-l-bold> Phil.`
	// `‹hb-fs-l-bold›Plutarchus‹/hb-fs-l-bold› Biogr., Phil.` [gr0007]`
	n := au.IDXname
	n = strings.Replace(n, "›", ">", -1)
	n = strings.Replace(n, "‹", "<", -1)
	n = betacode.BetaCodeCleanup(n)
	n = lat.ConvertLatinDiacriticals(n)
	n = simplegreekspan(n)
	n = simplelatinspan(n)
	n = strings.ReplaceAll(n, "&", "")
	n = strings.ToValidUTF8(n, "")

	n = strings.ReplaceAll(n, "<hb-fs-l-bold>", "")
	n = strings.ReplaceAll(n, "</hb-fs-l-bold>", "")
	n = strings.ReplaceAll(n, "<hb-fs-l-normal>", "")
	n = strings.ReplaceAll(n, "</hb-fs-l-normal>", "")

	if namemaker1.MatchString(n) {
		au.Cleaname = namemaker1.ReplaceAllString(n, "$1$2 ($3)")
		// two-passes to fix: `<hb-fs-l-bold>Eratosthenes </hb-fs-l-bold>et <hb-fs-l-bold>Eratosthenica</hb-fs-l-bold> Philol.`
		au.Cleaname = namemaker2.ReplaceAllString(au.Cleaname, "$1$2")
	} else {
		au.Cleaname = namemaker2.ReplaceAllString(n, "$1$2")
	}
	au.Shortname = shtrnamemaker.ReplaceAllString(au.Cleaname, "$1")
	// au.Shortname = removegenres(au.Shortname)

	// this should be gotten rid of, but the server uses both
	au.Name = au.Shortname
	au.IDXname = n
}

func removegenres(a string) string {
	var accum []string
	a = strings.ReplaceAll(a, ",", "")
	nameparts := strings.Split(a, " ")
	gg := append(global.KnownGenres, "scr") // to clean "Scr. Eccl."
	for i := len(nameparts) - 1; i >= 0; i-- {
		shoulddrop := false
		for _, g := range gg {
			g = strings.Title(g)
			if nameparts[i] == g {
				shoulddrop = true
			}
			g = g + "."
			if nameparts[i] == g {
				shoulddrop = true
			}
		}
		if !shoulddrop {
			accum = append(accum, nameparts[i])
		}
	}
	slices.Reverse(accum)
	return strings.Join(accum, " ")
}

func simplegreekspan(ttc string) string {
	return findgreek.ReplaceAllStringFunc(ttc, simplegreekshifter)
}

func simplegreekshifter(match string) string {
	// 0: whole match
	// 1: the text span
	groups := findgreek.FindStringSubmatch(match)
	return grk.ConvertGreekLowers(grk.ConvertGreekCapitals(groups[1]))
}

func simplelatinspan(ttc string) string {
	const (
		FSB = `<hb-fs-l-bold>$1</hb-fs-l-bold>`
	)
	return findlatin.ReplaceAllString(ttc, FSB)
}
