//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"testing"
)

func TestSimplegreekspan(t *testing.T) {
	ttc := `$*PERI\ TOU= XRHSTE/ON STUPTHRI/A| STROGGU/LH| A)NTI/LOGOS &(fort. pars operis $*XEIRO/MHKTA&) (e cod. Venet. Mar`
	result := simplegreekspan(ttc)
	want := "Περὶ τοῦ χρηϲτέον ϲτυπτηρίᾳ ϲτρογγύλῃ ἀντίλογοϲ (fort. pars operis Χειρόμηκτα) (e cod. Venet. Mar"
	if result != want {
		t.Errorf("simplegreekspan(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}

func TestGatherGreekMetadata(t *testing.T) {
	GatherTLGAuthorMetadata()
}

func TestRemoveGenres(t *testing.T) {
	aa := []string{
		"Evangelium Secundum Hebraeos",
		"Alexander Scr. Eccl.",
		"Hegesianax Astron., Epic.",
		"Anonymus Discipulus Isidori Milesii Mech.",
		"[Agathodaemon] Alchem.",
		"Xenophilus Mus., Phil.",
	}
	for _, a := range aa {
		fmt.Println(removegenres(a))
	}
}

func TestCleanDBAuthorNames(t *testing.T) {
	a := structs.DbAuthor{
		IDXname: "<hb-fs-l-bold>Diogenes Laertius</hb-fs-l-bold> Biogr.",
	}
	CleanDBAuthorNames(&a)
	a.PrintOut()
	a = structs.DbAuthor{
		IDXname: "Joannes <hb-fs-l-bold>Tzetzes</hb-fs-l-bold> Gramm., Poeta",
	}
	CleanDBAuthorNames(&a)
	a.PrintOut()
}

func TestLoadLatinCanon(t *testing.T) {
	_, _ = LoadLatinCanon("/Users/erik/Development/go/src/github.com/e-gun/HipparchiaGoBuilder/data/LAT/")
}
