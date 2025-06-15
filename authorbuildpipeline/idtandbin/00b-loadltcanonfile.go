package idtandbin

import (
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/latetidyups"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/worklines1initial"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/worklines2dbprep"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
	"strings"
)

const (
	LTCANONFILE = `LAT9999.TXT`
)

func LoadLatinCanon(dir string) map[string]string {
	// LAT9999.TXT works a lot like a standard author TXT...
	can := IdtFileLoad(dir + LTCANONFILE)
	ttc := betacode.BetaCodeCleanup(can)
	ttc = lat.LatinCleanup(ttc)
	ttc = latetidyups.TidyUp(ttc)
	ttc = worklines1initial.WorklinePrep(ttc)
	dbworklines := worklines2dbprep.PrepareForDB(ttc)

	for i, _ := range dbworklines {
		dbworklines[i].MarkedUp = strings.ReplaceAll(dbworklines[i].MarkedUp, "&nbsp;", "")
	}

	//for i, _ := range dbworklines {
	//	fmt.Println(i, dbworklines[i].MarkedUp)
	//}

	// now you have something that looks like:
	//2896 <hb-fs-l-italic>Vita Iuvenalis</hb-fs-l-italic><hb-fs-l-normal> (</hb-fs-l-normal><hb-fs-l-italic>A. Persi Flacci et D. Iuni Iuvenalis Saturae</hb-fs-l-italic><hb-fs-l-normal>, ed. </hb-fs-l-normal>
	//2897 W. V. Clausen, 1959). <hb-speaker>9969.001</hb-speaker>
	//2898
	//2899 <hb-fs-l-bold>Vitruvius Pollio</hb-fs-l-bold><hb-fs-l-normal>. </hb-fs-l-normal><hb-fs-l-italic>De Architectura</hb-fs-l-italic><hb-fs-l-normal> (</hb-fs-l-normal><hb-fs-l-italic>Vitruvii de Architectura</hb-fs-l-italic><hb-fs-l-normal>, ed. F. </hb-fs-l-normal>
	//2900 Krohn, 1912). <hb-speaker>1056.001</hb-speaker>
	//2901
	//2902 <hb-fs-l-bold>Volcacius Sedigitus</hb-fs-l-bold><hb-fs-l-normal>. carmina (</hb-fs-l-normal><hb-fs-l-italic>Fragmenta Poetarum Latinorum Epicorum</hb-fs-l-italic><hb-fs-l-normal> </hb-fs-l-normal>
	//2903 <hb-fs-l-italic>et Lyricorum praeter Ennium et Lucilium</hb-fs-l-italic><hb-fs-l-normal>, ed. W. Morel, 1927). </hb-fs-l-normal>
	//2904 <hb-speaker>0330.001</hb-speaker>

	// the pattern is Pubinfo ... <hb-speaker>AUID.WKID</hb-speaker>

	textblock := make([]string, len(dbworklines))
	for i, _ := range dbworklines {
		textblock[i] = dbworklines[i].MarkedUp
	}

	joined := strings.Join(textblock, "")
	split := strings.Split(joined, "</hb-speaker>")

	parser1 := regexp.MustCompile(`(.*?)<hb-speaker>(\d\d\d\d)\.(\d\d\d)`)

	uidpublicationmap := make(map[string]string)
	for _, s := range split {
		groups := parser1.FindStringSubmatch(s)
		if len(groups) == 4 {
			uid := "lt" + groups[2] + global.AUTHWORKSEPARATOR + groups[3]
			uidpublicationmap[uid] = groups[1]
		}
	}

	// now:
	//lt0684w015  </hb-fs-l-normal><hb-fs-l-bold>M. Terentius Varro</hb-fs-l-bold><hb-fs-l-normal>. fragmenta de historia litterarum (</hb-fs-l-normal><hb-fs-l-italic>Grammaticae</hb-fs-l-italic><hb-fs-l-normal> </hb-fs-l-normal><hb-fs-l-italic>Romanae Fragmenta</hb-fs-l-italic><hb-fs-l-normal>, ed. G. Funaioli, 1907).
	//lt0692w004 <hb-fs-l-bold>Appendix Vergiliana</hb-fs-l-bold><hb-fs-l-normal>. </hb-fs-l-normal><hb-fs-l-italic>Aetna</hb-fs-l-italic><hb-fs-l-normal> (</hb-fs-l-normal><hb-fs-l-italic>Appendix Vergiliana</hb-fs-l-italic><hb-fs-l-normal>, ed. F. R. D. </hb-fs-l-normal>Goodyear, 1966).
	//lt0692w005 <hb-fs-l-bold>Appendix Vergiliana</hb-fs-l-bold><hb-fs-l-normal>. </hb-fs-l-normal><hb-fs-l-italic>Copa</hb-fs-l-italic><hb-fs-l-normal> (</hb-fs-l-normal><hb-fs-l-italic>Appendix Vergiliana</hb-fs-l-italic><hb-fs-l-normal>, ed. E. J. Kenney, </hb-fs-l-normal>1966).

	// parser1 will return something like when you hit HGS ProlixBrowswerCitations():
	// Q. Tullius Cicero. carmina (Fragmenta Poetarum Latinorum Epicorum et Lyricorum praeter Ennium et Lucilium, ed. W. Morel, 1927).

	// but, in practice, you already know "Author and Title" and really only need "Text and Edition"

	parser2 := regexp.MustCompile(`([^(]+)\((.+)\)`)
	for k, v := range uidpublicationmap {
		groups := parser2.FindStringSubmatch(v)
		if len(groups) == 3 {
			newv := groups[2]
			uidpublicationmap[k] = newv
		}
	}

	// down to just: </hb-fs-l-normal><hb-fs-l-italic>Grammatici Latini ex</hb-fs-l-italic><hb-fs-l-normal> </hb-fs-l-normal><hb-fs-l-italic>Recensione Henrici Keilii</hb-fs-l-italic><hb-fs-l-normal>. Vol. 6, ed. H. Keil, 1874

	return uidpublicationmap
}
