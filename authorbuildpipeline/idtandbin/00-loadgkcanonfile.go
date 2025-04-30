//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"bytes"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/worklines1initial"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/dating"
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

const (
	GKCANONFILE = `DOCCAN2.TXT`
)

var (
	findnewauthor = regexp.MustCompile(`<hb-set_l_1_to_a />key (\d\d\d\d)`)
	findnewwork   = regexp.MustCompile(`<hb-set_l_1_to_\d />key (\d\d\d\d)\s(...)`)
	findkvpairs   = regexp.MustCompile(`^([a-z][a-z][a-z])\s(.*?)$`)
	findtitle     = regexp.MustCompile(`&3([^&]*?)&`)
)

func LoadGreekCanon(dir string) (map[string]structs.DbAuthor, map[string]structs.DbWork) {
	// DOCCAN2.TXT works a lot like a standard author TXT...
	aumap := make(map[string]structs.DbAuthor)
	wkmap := make(map[string]structs.DbWork)

	can := IdtFileLoad(dir + GKCANONFILE)
	can = initialcanoncleaner(can)

	auu, wkk := findauthorsandworks(can)

	for k, v := range auu {
		id := strings.TrimSpace(k)
		aumap[id] = parseauthdata(id, v)
	}

	for k, v := range wkk {
		id := strings.TrimSpace(k)
		wkmap[id] = parseworkdata(id, v)
	}

	return aumap, wkmap
}

func StoreCanonInSharedMaps(aumap map[string]structs.DbAuthor, wkmap map[string]structs.DbWork) {
	for k, v := range aumap {
		global.TheCanonAuMap.Set(k, v)
	}

	for k, v := range wkmap {
		global.TheCanonWkMap.Set(k, v)
	}
}

func initialcanoncleaner(ttc string) string {
	can := betacode.BetaCodeCleanup(ttc)
	can = worklines1initial.WorklinePrep(can) // WorklinePrep will make some betacode reappear; clean it later
	can = strings.ReplaceAll(can, `<hb-incr_l_0_by_1 />`, ``)
	can = strings.ReplaceAll(can, `<hb-incr_l_2_by_1 />`, ``)
	can = strings.ReplaceAll(can, `<hb-metadata_newauthor value="9998" /><hb-metadata_newwork value="001" />`, ``)
	// now you have authors-and-works that look like:
	//<hb-set_l_1_to_a />key 9022
	//nam Joannes TZETZES
	//srt Tzetzes, Joannes
	//epi Gramm.
	//epi Poeta
	//geo Constantinopolitanus
	//dat A.D. 12
	//<hb-set_l_1_to_3 />key 9022 003
	//wrk &1Prolegomena de comoedia Aristophanis&
	//cla Gramm.
	//xmt Cod
	//typ Book in series
	//wct 3,126
	//cit Section／line
	//tit &3Prolegomena de comoedia. Scholia in Acharnenses, Equites, Nubes&
	//pub Bouma
	//pla Groningen
	//pyr 1975
	//pag 22–38
	//ser &3Scholia in Aristophanem &1.1A
	//brk *)ALE/CANDROS O( *AI)TWLO\S KAI\ *LUKO/FRWN O( *XALKIDEU\S
	//MEGALODWRI/AIS& (prooemium i): pp. 22–31
	//brk *)ALE/CANDROS O( *AI)TWLO\S KAI\ *LUKO/FRWN O( *XALKIDEU/S, A)LLA\
	//&(prooemium ii): pp. 31–38
	//edr Koster, W.J.W.

	return can
}

func findauthorsandworks(can string) (map[string][]string, map[string][]string) {
	authors := make(map[string][]string)
	works := make(map[string][]string)
	var appender []string
	var previous string
	previousiswork := false

	lines := strings.Split(can, "\n")
	for _, line := range lines {
		if findnewauthor.MatchString(line) {
			if previousiswork {
				works[previous] = appender
			} else {
				authors[previous] = appender
			}
			na := findnewauthor.ReplaceAllString(line, "gr$1")
			previous = na
			previousiswork = false
			appender = []string{}
		}
		if findnewwork.MatchString(line) {
			if previousiswork {
				works[previous] = appender
			} else {
				authors[previous] = appender
			}

			nw := findnewwork.ReplaceAllStringFunc(line, buildworkname)
			previous = nw
			previousiswork = true
			appender = []string{}
		}
		appender = append(appender, line)
	}

	if previousiswork {
		works[previous] = appender
	} else {
		authors[previous] = appender
	}

	// now you will have works that look like
	// <hb-set_l_1_to_1 />key 9021 001
	//wrk &1De magna et sacra arte &(sub nomine Stephani Alexandrini philosophi)
	//cla Alchem.
	//xmt Cod
	//typ Book
	//wct 17,336
	//cit Volume／page／line
	//tit &3Physici et medici Graeci minores, &vol. 2
	//pub Reimer
	//pla Berlin
	//pyr 1842
	//rpu Hakkert
	//rpl Amsterdam
	//ryr 1963
	//pag 199–253
	//brk *PRA/CEIS& 1–2: pp. 199–208
	//brk Epistula ad Theodorum: p. 208
	//brk *PRA/CEIS& 3–8: pp. 209–242
	//brk *PRA=CIS& 9 (Didascalia ad Heraclium): pp. 243–253
	//edr Ideler, J.L.

	// authors look like
	//<hb-set_l_1_to_a />key 9021
	//nam STEPHANUS
	//syn fiq Stephanus Phil.
	//epi Alchem.
	//geo Alexandrinus
	//dat A.D. 7
	//vid Cf. et STEPHANUS Phil. (9019)

	return authors, works
}

func buildworkname(match string) string {
	groups := findnewwork.FindStringSubmatch(match)
	return "gr" + groups[1] + "w" + groups[2]
}

func parseauthdata(id string, data []string) structs.DbAuthor {
	// tags seen:
	// [brk, cit, cla, crf, dat, edr, epi, gen, geo, nam, pag, pla, pub, pyr, rpl, rpu, ryr, ser, srt, syn, tit, typ, vid, wct, wrk, xmt]

	var theauthor structs.DbAuthor
	theauthor.UID = id
	theauthor.Language = "G"
	var epislice []string

	for _, line := range data {
		kvp := findkvpairs.FindStringSubmatch(line)
		if kvp != nil {
			kvp[2] = strings.TrimSpace(kvp[2])
			switch kvp[1] {
			case "nam":
				theauthor.Name = kvp[2]
			case "srt":
				theauthor.Shortname = kvp[2]
			case "syn":
				theauthor.Shortname = kvp[2]
			case "geo":
				theauthor.Location = kvp[2]
			case "dat":
				theauthor.RecDate = kvp[2]
				theauthor.ConvDate = dating.ParseTLGDate(theauthor.RecDate)
			case "epi":
				// people can be both "Phil." and something else...
				epislice = append(epislice, kvp[2])
			}
		}
	}
	theauthor.Genres = strings.Join(epislice, "; ")

	return theauthor
}

func parseworkdata(id string, data []string) structs.DbWork {
	// work tags seen
	// [brk, cit, cla, crf, dat, edr, epi, geo, nam, pag, pla, pub, pyr, ref, rpl, rpu, ryr, ser, sur, tit, typ, vel,
	// vid, wct, wrk, xmt]

	var thework structs.DbWork
	var strayinfo []string
	var genreslice []string
	var tit string
	var pub string
	var pla string
	var pyr string
	var edr string
	var cit string

	thework.UID = id
	thework.Language = "G"

	for _, line := range data {
		kvp := findkvpairs.FindStringSubmatch(line)
		if kvp != nil {
			kvp[2] = strings.TrimSpace(kvp[2])
			switch kvp[1] {
			case "wrk":
				thework.Title = kvp[2]
			case "cla":
				genreslice = append(genreslice, kvp[2])
			case "xmit":
				thework.Xmit = kvp[2]
			case "wct":
				thework.WdCount, _ = strconv.Atoi(strings.ReplaceAll(kvp[2], ",", ""))
				if thework.WdCount == 0 {
					thework.WdCount = -1
				}
			case "cit":
				cit = kvp[2]
			case "tit":
				tit = kvp[2]
			case "pub":
				pub = kvp[2]
			case "pla":
				pla = kvp[2]
			case "pyr":
				pyr = kvp[2]
			case "edr":
				edr = kvp[2]
			default:
				// fmt.Println(kvp[0])
				strayinfo = append(strayinfo, kvp[2])
			}
		}
	}

	setcitationinformation(&thework, cit)
	setpublicationinformation(&thework, tit, pub, pla, pyr, edr)

	genreslice = generic.DropEmptyStrings(generic.Unique(genreslice))
	// todo: there is a genres parser error below
	// gr0007w009: Rhet. Epist. Dialog. Biogr. xmt Cod Med. Theol. Gnom. Hist. Phil. cla Myth. Paroem.  Polyhist. Narr. Fict. Nat. Astron.
	thework.Genre = strings.Join(genreslice, "; ")

	// fmt.Println(thework.UID, thework.Title, thework.Genre)
	// gr9019w001 <hb-fs-l-bold>In Aristotelis librum de interpretatione commentarium<hb-fs-l-normal> </hb-fs-l-normal> Comm.; Phil.

	return thework
}

func setpublicationinformation(thework *structs.DbWork, tit string, pub string, pla string, pyr string, edr string) {
	const (
		// PUBTEMPL - lots of corner cases cut...
		PUBTEMPL = `{{.tit}}, {{.pub}} {{.pla}} {{.pyr}} ({{.edr}})`
	)
	m := map[string]string{
		"tit": tit,
		"pub": pub,
		"pla": pla,
		"pyr": pyr,
		"edr": edr,
	}

	t := template.Must(template.New("").Parse(PUBTEMPL))

	var b bytes.Buffer
	if ee := t.Execute(&b, m); ee != nil {
		fmt.Println(ee)
	}
	fpub := lat.ConvertLatinDiacriticals(b.String())
	fpub = simpletitlespan(fpub)
	thework.Pub = fpub
}

func simpletitlespan(ttc string) string {
	return findtitle.ReplaceAllString(ttc, `<hb-fs-l-italic>$1</hb-fs-l-italic>`)
}

func setcitationinformation(w *structs.DbWork, cit string) {
	cc := strings.Split(cit, "／")
	slices.Reverse(cc)

	assigntolevel := func(i int, c string, w *structs.DbWork) {
		switch i {
		case 0:
			w.LL0 = c
		case 1:
			w.LL1 = c
		case 2:
			w.LL2 = c
		case 3:
			w.LL3 = c
		case 4:
			w.LL4 = c
		case 5:
			w.LL5 = c
		default:
			fmt.Println("setcitationinformation() invalid citation level")
		}
	}

	for i, c := range cc {
		assigntolevel(i, c, w)
	}
}

const (
	CTESTDAT = `<hb-set_l_1_to_a />key 9021
nam STEPHANUS
syn fiq Stephanus Phil.
epi Alchem.
geo Alexandrinus
dat A.D. 7
vid Cf. et STEPHANUS Phil. (9019)
<hb-set_l_1_to_1 />key 9021 001
wrk &1De magna et sacra arte &(sub nomine Stephani Alexandrini philosophi)
cla Alchem.
xmt Cod
typ Book
wct 17,336
cit Volume／page／line
tit &3Physici et medici Graeci minores, &vol. 2
pub Reimer
pla Berlin
pyr 1842
rpu Hakkert
rpl Amsterdam
ryr 1963
pag 199–253
brk *PRA/CEIS& 1–2: pp. 199–208
brk Epistula ad Theodorum: p. 208
brk *PRA/CEIS& 3–8: pp. 209–242
brk *PRA=CIS& 9 (Didascalia ad Heraclium): pp. 243–253
edr Ideler, J.L.

<hb-set_l_1_to_a />key 9022
nam Joannes TZETZES
srt Tzetzes, Joannes
epi Gramm.
epi Poeta
geo Constantinopolitanus
dat A.D. 12
<hb-set_l_1_to_3 />key 9022 003
wrk &1Prolegomena de comoedia Aristophanis&
cla Gramm.
xmt Cod
typ Book in series
wct 3,126
cit Section／line
tit &3Prolegomena de comoedia. Scholia in Acharnenses, Equites, Nubes&
pub Bouma
pla Groningen
pyr 1975
pag 22–38
ser &3Scholia in Aristophanem &1.1A
brk *)ALE/CANDROS O( *AI)TWLO\S KAI\ *LUKO/FRWN O( *XALKIDEU\S
MEGALODWRI/AIS& (prooemium i): pp. 22–31
brk *)ALE/CANDROS O( *AI)TWLO\S KAI\ *LUKO/FRWN O( *XALKIDEU/S, A)LLA\
&(prooemium ii): pp. 31–38
edr Koster, W.J.W.
<hb-set_l_1_to_4 />key 9022 004
wrk &1Versus de poematum generibus&
cla Gramm.
cla Iamb.
xmt Cod
typ Book in series
wct 2,813
cit Section／line
tit &3Prolegomena de comoedia. Scholia in Acharnenses, Equites, Nubes&
pub Bouma
pla Groningen
pyr 1975
pag 84–109
ser &3Scholia in Aristophanem &1.1A
brk Introductio: pp. 84–94
brk De comoedia: pp. 94–98
brk De tragoedia: pp. 99–109
edr Koster, W.J.W.
<hb-set_l_1_to_6 />key 9022 006
wrk &1De Pleiadibus &(excerptum)
cla Astron.
cla Comm.
xmt Cod
typ Book
wct 949
cit Page／line
tit &3Scholia in Aratum vetera&
pub Teubner
pla Stuttgart
pyr 1974
pag 547–551
edr Martin, J.
<hb-set_l_2_to_9022x />
<hb-set_l_1_to_1 />key 9022 X01
wrk &1Exegesis in scutum Hesiodi&
cla Schol.
xmt Cod
ref Gaisford, pp. 609–654
crf Cf. Joannes PEDIASIMUS Gramm. (2592 003)
<hb-set_l_1_to_3 />key 9022 X03
wrk &1De poeticae generibus &(fort. auctore Joanne Tzetza)
cla Gramm.
xmt Cod
ref Koster, vol. 1.1A, p. 50
crf Cf. &3PROLEGOMENA DE COMOEDIA& (3002 015)
<hb-set_l_1_to_5 />key 9022 X05
wrk &1Vita Oppiani&
cla Biogr.
xmt Cod
ref Colonna, p. 40
crf Cf. &3VITAE OPPIANI& (4172 006)
<hb-set_l_1_to_6 />key 9022 X06
wrk &1Scholia et argumenta in Aristophanem&
cla Schol.
xmt Cod
crf Cf. &3SCHOLIA IN ARISTOPHANEM& (5014 015–21, 23)
<hb-set_l_1_to_7 />key 9022 X07
wrk &1Epigrammata in Oppianum&
cla Epigr.
cla Schol.
xmt Cod
ref Bussemaker, pp. 260, 276
crf Cf. &3SCHOLIA IN OPPIANUM& (5032 002)
<hb-set_l_1_to_8 />key 9022 X08
wrk &1Interpretatio et scholia in Hesiodi opera et dies&
cla Schol.
xmt Cod
ref Gaisford, pp. 10–22, 23–447 passim
crf Cf. &3SCHOLIA IN HESIODUM& (5025 002)
<hb-set_l_1_to_9 />key 9022 X09
wrk &1Introductio et scholia in Lycophronem&
cla Schol.
xmt Cod
brk Introductio: olim sub auctore Isaac Tzetza
ref Scheer, pp. 1–4, 8–398 passim
crf Cf. &3SCHOLIA IN LYCOPHRONEM& (5030 001)
<hb-set_l_2_to_9023 />
<hb-set_l_1_to_a />key 9023
nam THOMAS MAGISTER
syn vel Theodulus
epi Philol.
geo Thessalonicensis
geo Constantinopolitanus
dat A.D. 13–14
<hb-set_l_1_to_1 />key 9023 001
wrk &1Ecloga nominum et verborum Atticorum&
cla Lexicogr.
xmt Cod
typ Book
wct 47,537
cit Alphabetic letter／page／line
tit &3Thomae Magistri sive Theoduli monachi ecloga vocum Atticarum&
pub Orphantropheus
pla Halle
pyr 1832
rpu Olms
rpl Hildesheim
ryr 1970
pag 1–411
edr Ritschl, F.
<hb-set_l_1_to_2 />key 9023 002
wrk &1Poemata de Arato &[Dub.]
cla Iamb.
xmt Cod
typ Book
wct 285
cit Poem／line
tit &3Scholia in Aratum vetera&
pub Teubner
pla Stuttgart
pyr 1974
pag 558–559
edr Martin, J.
<hb-set_l_2_to_9023x />
<hb-set_l_1_to_2 />key 9023 X02
wrk &1Scholia et argumenta in Aristophanem&
cla Schol.
xmt Cod
crf Cf. &3SCHOLIA IN ARISTOPHANEM& (5014 005, 012)
<hb-set_l_1_to_3 />key 9023 X03
wrk &1Scholia in Pindarum&
cla Schol.
xmt Cod
crf Cf. &3SCHOLIA IN PINDARUM& (5034 003–005, 007)
<hb-set_l_1_to_4 />key 9023 X04
wrk &1Scholia in Aeschylum&
cla Schol.
xmt Cod
crf Cf. &3SCHOLIA IN AESCHYLUM& (5010 006, 007, 010)
<hb-set_l_1_to_5 />key 9023 X05
wrk &1Scholia in Sophoclis Oedipum tyrannum&
cla Schol.
xmt Cod
ref Longo, pp. 167–265
crf Cf. &3SCHOLIA IN SOPHOCLEM& (5037 005)`
)
