//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines3cleanup

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/dating"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/idtandbin"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

// the goal is to take an ins, ddp, or chr dbc as built by the standard parser and to break it down into a new set of databases
//
// the original files just heap up various documents inside of larger works, but this means you lose useful access to the
// location and date information for each individual document, information that could be used as the basis for a search
//
// for example 'INS0080' ==> 'Black Sea and Scythia Minor'
//	w001 is 'IosPE I(2)'
//	w002 is 'CIRB'
//
// document01 of w002 has the following associated with it:
//	<hmu_metadata_region value="N. Black Sea" />
//	<hmu_metadata_city value="Pantikapaion" />
//	<hmu_metadata_date value="344-310a" />
//	<hmu_metadata_publicationinfo value="IosPE II 1" />
//
// document180 of w002 has the following:
//	<hmu_metadata_region value="N. Black Sea" />
//	<hmu_metadata_city value="Myrmekion" />
//	<hmu_metadata_date value="c 400-350a" />
//	<hmu_metadata_publicationinfo value="IosPE IV 294[cf CIRB p.479]" />
//
// What if you wanted to search only documents from Pantikapaion?
// What if you wanted to restrict your search to 375BCE?
// What if you want to search *all* Greek authors and inscriptions older than 400BCE?
//
// what needs to happen is that the works need to become authors and the documents need to become works
// then as works the metadata can be assigned to the relevant fields of the workdb
//
// the results of the first pass db generation are the following:
//	there are 24 inscription files that will turn into 24 authors
//	463 works will be distributed among these 24 authors
//
//	there are 213 papyrus files that will turn into 213 authors
//	516 works will be distributed among these 213 authors
//
// things have to happen in several passes
// first you need to build the new hybrid authors from authors+works
// but you can't fill out 'floruit' yet: you need line 1 of the newwork dbc
// but the new work dbs need to be extracted from the old work dbc
//
// level05 = document number (sorta: not always the same as what is asserted at level06)
// level01 = face (recto, verso, etc; usually recto: this will result in spammy output unless you use workarounds in HipparchiaServer)
// level00 = line
// AND a random use of the other levels which sometimes contain information

// what you would see if you did not subdivide...

// hgdb=> select index,wkuniversalid,annotations from in0060 where annotations ~* '^region' order by index asc;

// hgdb=> select * from authors where universalid = 'in0060';
// universalid | language |  idxname  |  akaname  | shortname | cleanname | genres | recorded_date | converted_date | location
//-------------+----------+-----------+-----------+-----------+-----------+--------+---------------+----------------+----------
// in0060      |          | Macedonia | Macedonia | Macedonia | Macedonia |        |               |              0 |

// hgdb=> select universalid,title from works where universalid ~* 'in0060';
// universalid |               title
//-------------+-----------------------------------
// in0060w002  | Epigraphes Ano Makedonias I [EAM]
// in0060w003  | SEG 1–41
// in0060w004  | Vulic/, Spomenik 71–98 [NW Mac.]
// in0060w001  | IG X 2,1 [Thessalonike]

var (
	ddpnewworkspan  = regexp.MustCompile(`^<hb-fs-l-normal>(.*?﹦.*?)</hb-fs-l-normal>$`)
	nomarkup        = regexp.MustCompile(`<[^>]*>`)
	remapcorpusname = map[string]string{
		global.INSFIRSTPASS: global.INSABBREV,
		global.DDPFIRSTPASS: global.DDPABBREV,
		global.CHRFIRSTPASS: global.CHRABBREV,
	}
)

func RemapInscriptionAuthorsAndWorks(lines []structs.DbWorkline, authorid string) ([]structs.DbAuthor, []structs.DbWork, []structs.DbWorkline) {
	// note that in a multiprocessing environment you will have problems with the new names so use UniqueAuthorNames
	UniqueWorkNamer := global.NewUniqueNamer(3)

	globalauthornameprefix, specificauthornamesuffix := buildremapper(authorid)

	newworkmap := make(map[string]structs.DbWork)

	// blank values
	previouswork := structs.DbWork{
		FirstLine: 1,
	}

	previousworkannotations := ""

	prevlline := structs.DbWorkline{
		TbIndex: 1,
	}

	workingwithauthorid := "-null-workingwithauthorid-value-" // the first line had best set this...
	oldwkuid := "-null-oldwkuid-value-"

	theprefix := remapcorpusname[lines[0].WkUID[0:global.ABBREVLEN]] // "chr" --> "chx"

	wordcount := 0

	authorsderivedfromworks := make(map[string]structs.DbAuthor)

	for i, l := range lines {
		//fmt.Println(l)
		//if i > 40 {
		//	os.Exit(1)
		//}

		// the first lines at this point will look like...
		// 0 {0 1 -1 -1 -1 -1 a in0150w001 [τῷ  —  ἱερ]- τῷ ἱερ ἱερεῖ τω ιερ ιερει ἱερεῖ publicationinfo: AD 1915, 76 · workabbrev: Chios }
		// 1 {1 1 -1 -1 -1 -1 a in0150w001 εῖ̣, [ὅταν  —  θύῃ δίδοϲ]- ὅταν θύῃ δίδοϲθαι οταν θυη διδοϲθαι δίδοϲθαι }
		// 2 {2 1 -1 -1 -1 -1 1 in0150w001 θα[ι  —  ἡμίεκτον χρ]- ἡμίεκτον χρ χρυϲοῦ ημιεκτον χρ χρυϲου χρυϲοῦ }

		if i > 0 {
			prevlline = lines[i-1]
		}

		if lines[i].WkUID != oldwkuid {
			// in0150w001 --> in0150w002
			// you changed ORIGNIALWORK so you just triggered a NEWAUTHOR

			// assignnametoauthor, _ := lookforauthornameinmetadata(l)
			// newidxname := fmt.Sprintf("%s · %s", oldau.IDXname, assignnametoauthor)

			newidxname := ""
			if specificauthornamesuffix[lines[i].WkUID[global.ABBREVLEN:]] != " " {
				// note that `empty` is ` ` and not ``
				newidxname = globalauthornameprefix + " · " + specificauthornamesuffix[lines[i].WkUID[global.ABBREVLEN:]]
			} else {
				newidxname = globalauthornameprefix
			}
			newidxname = betacode.ReplacePercentSigns(newidxname)
			newidxname = lat.ConvertLatinDiacriticals(newidxname)
			newidxname = betacode.SimpleLatinSpanLATE(newidxname)
			newidxname = strings.ReplaceAll(newidxname, "`", "")
			newidxname = strings.ReplaceAll(newidxname, "&", "")

			// need to ensure that old filename collisions impossible; otherwise 'in0020' can happen in the new data
			// but this should now be guaranteed by the gap between DDPFIRSTPASS and DDPABBREV
			workingwithauthorid = theprefix + global.UniqueAuthorNames.GetNameWithTrailingDigit(4)
			strippedname := nomarkup.ReplaceAllString(newidxname, "")

			a := structs.DbAuthor{
				UID:       workingwithauthorid,
				IDXname:   newidxname,
				Cleaname:  strippedname, // HGS: aunamehint()
				Shortname: newidxname,
				Name:      newidxname, // HGS: ProlixBrowswerCitations()
			}
			authorsderivedfromworks[a.UID] = a
			oldwkuid = l.WkUID

			// fmt.Println(a)
			// {in000a  (Ionia) Did  (Ionia) Did (Ionia) Did   0  []}
			// {in000b  (Ionia) Eph  (Ionia) Eph (Ionia) Eph   0  []}
			// {in000c  (Ionia) Sup Eph  (Ionia) Sup Eph (Ionia) Sup Eph   0  []}
			// {in000d  (Ionia) Eryth  (Ionia) Eryth (Ionia) Eryth   0  []}
			// {in000e  (Ionia) Klaz  (Ionia) Klaz (Ionia) Klaz   0  []}
		}

		if foundanewwork(l, prevlline) {
			//fmt.Println("foundanewwork")
			//fmt.Println(l)
			// {75857 280 -1 -1 -1 -1 1 in0150w020 [ —  κ]οινῇ ἄρ[χοντεϲ﹖  — ] κοινῇ ἄρχοντεϲ κοινη αρχοντεϲ  publicationinfo: BCH 1922, 343, no. 34 · documentnumber: 286}
			// {75858 281 -1 -1 -1 -1 1 in0150w020 [ — ]α̣[ — ] α α  publicationinfo: BCH 1925, 310, no. 7 · documentnumber: 287}
			// {75872 282 -1 -1 -1 -1 1 in0150w020 ∙    publicationinfo: BCH 1925, 312, no. 12 · documentnumber: 288}
			// {75876 284 -1 -1 -1 -1 1 in0150w020 [ — ]α̣υτα παρα αυτα παρα αυτα παρα  publicationinfo: Anadolu 1965, 29-157, III · documentnumber: 289}

			newwork := gathernewworkinfo(l)

			// Res Gestae just reasserts itself periodically and tries to fool you:
			// `region: Galatia · city: Ankyra · date: 19 ac · workabbrev: RG` (8x or so)

			if strings.Contains(l.WkUID, global.CHRFIRSTPASS+"0120") && previousworkannotations == l.Annotations {
				global.MSG("'new' PHI work is not new: \t" + l.Annotations)
				continue
			} else {
				previousworkannotations = l.Annotations
			}

			wn := UniqueWorkNamer.GetName(global.WKIDLEN)
			newwork.UID = workingwithauthorid + global.AUTHWORKSEPARATOR + wn

			newwork.FirstLine = l.TbIndex

			// todo: this is truly suspect because so many chr lines are broken right now and need ungreeking
			if l.HasGreek() {
				newwork.Language = "G"
			} else {
				newwork.Language = "L"
			}

			newworkmap[newwork.UID] = newwork

			previouswork.WdCount = wordcount
			wordcount = 0

			previouswork.LastLine = l.TbIndex - 1
			newworkmap[previouswork.UID] = previouswork
			previouswork = newwork
		}

		wordcount += l.GetWordcount()
		lines[i].WkUID = previouswork.UID
	}
	// need to catch the last work
	previouswork.LastLine = lines[len(lines)-1].TbIndex - 1
	previouswork.WdCount = wordcount
	newworkmap[previouswork.UID] = previouswork

	var works []structs.DbWork
	for _, v := range newworkmap {
		// there is an empty work from when you start
		if v.UID != "" {
			works = append(works, v)
		}
	}

	// a sample new work:
	// {in001yw1f1 region: C. Crete · city: Gortyna · date: IIa · documentnumber: 361 · workabbrev: IC 4  #361IC 4 line face -1 -1 -1 -1 inscription inscribed  C. Crete; Gortyna IIa 2500 0 14729 14730 true}

	// the authors are hollow: effectively a slice of IDs
	// in000a, in000b, in000c, ...

	// hgdb=> select * from authors where universalid ~* 'in006o';
	// universalid | language | idxname | akaname | shortname | cleanname | genres | recorded_date | converted_date | location
	//-------------+----------+---------+---------+-----------+-----------+--------+---------------+----------------+----------
	// in006o      |          |         |         |           |           |        |               |              0 |
	//(1 row)

	var authors []structs.DbAuthor
	for _, newauth := range authorsderivedfromworks {
		authors = append(authors, newauth)
	}

	// you have to have an author for every work or HGS will panic when it tries to do calculatewholeauthorsearches()
	checkwkauthors := make(map[string]bool)
	for _, a := range authors {
		checkwkauthors[a.UID] = true
	}
	for _, w := range works {
		if _, ok := checkwkauthors[w.UID[0:global.AUIDLEN]]; !ok {
			fmt.Println("fatal build error no author found", w)
			fmt.Println(checkwkauthors)
			fmt.Println(w.UID[0:global.AUIDLEN])
			os.Exit(0)
		}
	}

	return authors, works, lines
}

func gathernewworkinfo(l structs.DbWorkline) structs.DbWork {
	corpus := l.WkUID[0:global.ABBREVLEN]
	switch corpus {
	case global.INSFIRSTPASS:
		return gathernewinsworkinfo(l)
	case global.CHRFIRSTPASS:
		return gathernewchrworkinfo(l)
	case global.DDPFIRSTPASS:
		return gathernewddpworkinfo(l)
	default:
		return gathernewinsworkinfo(l)
	}
}

func gathernewinsworkinfo(l structs.DbWorkline) structs.DbWork {
	const (
		TEMPL1 = `%s%s (%s)`
		TEMPL2 = `%s%s`
	)
	// 	metadatacategories = map[int]string{
	//		0:   "newauthor",
	//		1:   "newwork",
	//		2:   "workabbrev",
	//		3:   "authabbrev",
	//		97:  "region",
	//		98:  "city",
	//		99:  "notes",
	//		100: "date",
	//		101: "publicationinfo",
	//		102: "additionalpubinfo",
	//		103: "stillfurtherpubinfo",
	//		108: "provenance",
	//		114: "reprints",
	//		116: "unknownmetadata116",
	//		122: "documentnumber",

	// workabbrev contains "author"?

	kvp := l.GatherMetadata()
	for k, v := range kvp {
		kvp[k] = betacode.ReplacePercentSigns(lat.ConvertLatinDiacriticals(betacode.ReplacePoundSigns(v)))
	}

	date, ok := kvp["date"]
	if !ok {
		date = ""
	}

	var allprov []string
	region, ok := kvp["region"]
	if ok {
		allprov = append(allprov, region)
	}

	city, ok := kvp["city"]
	if ok {
		allprov = append(allprov, city)
	}

	prov := strings.Join(allprov, ", ")

	var allpi []string
	publicationinfo, ok := kvp["publicationinfo"]
	if ok {
		allpi = append(allpi, publicationinfo)
	}
	additionalpubinfo, ok := kvp["additionalpubinfo"]
	if ok {
		allpi = append(allpi, additionalpubinfo)
	}
	stillfurtherpubinfo, ok := kvp["stillfurtherpubinfo"]
	if ok {
		allpi = append(allpi, stillfurtherpubinfo)
	}

	pi := strings.Join(allpi, "; ")

	doc, ok := kvp["documentnumber"]
	if ok {
		doc = "#" + doc + ": "
	}

	// many/most inscriptions have only lines and so l0 alone matters; but some are more complex and use other levels
	// see Rhodes and Rhodian Peraia · IG XII,1 [Rhodes], (Rhodos, Rhodos): `γʹ Οὐηράνιοϲ Εὔφη[μοϲ...` (in0100)
	// here l2 has `fr a` and `fr b` and l1 has `col I` and `col II`
	// so you should test for blanks at higher levels; here `fr` and `col` are the proper labels
	// l3 never matters?

	// this turns into a problem for HGS: FindValidLevelValues() will do `= 'fr'` at l2, but only `= 'fr a'` finds
	// you can't store either `fr a` or `fr b` as the level label, though because both are present
	// in effect you have a fourth level (`fr`) with two values: `a` and `b`; the same goes for `col I`, etc. too...
	// it does not look like the builder can really handle this; HGS needs to try to accommodate it...

	levelnameifnotempty := func(s string) string {
		spl := strings.Split(s, " ")
		// "1" is the standard blank value; "0" can be found too
		if spl[0] == "1" || spl[0] == "0" {
			return ""
		}
		return spl[0]
	}

	var tit string
	if prov != "" {
		tit = fmt.Sprintf(TEMPL1, doc, publicationinfo, prov)
	} else {
		tit = fmt.Sprintf(TEMPL2, doc, publicationinfo)
	}

	tit = strings.ReplaceAll(tit, "  ", " ") // `#1312:  (Phryg., Apameia (Dinar))` --> `#1312: (Phryg., Apameia (Dinar))`

	gen := "inscr"
	if l.WkUID[0:global.ABBREVLEN] == global.DDPFIRSTPASS {
		gen = "docu"
	}
	nw := structs.DbWork{
		UID:       "",
		Title:     betacode.SimpleLatinSpanLATE(lat.ConvertLatinDiacriticals(tit)),
		Language:  "",
		Pub:       pi,
		LL0:       "line",
		LL1:       levelnameifnotempty(l.Lvl1Value),
		LL2:       levelnameifnotempty(l.Lvl2Value),
		LL3:       "",
		LL4:       "",
		LL5:       "",
		Genre:     gen,
		Xmit:      fmt.Sprintf("direct (%s)", l.WkUID[0:global.AUIDLEN]), // so you can go back and bug hunt...
		Type:      "",
		Prov:      remapprov(prov),
		RecDate:   date,
		ConvDate:  dating.ConvertOnePHIStringDate(date),
		WdCount:   0,
		FirstLine: l.TbIndex,
		LastLine:  0,
		Authentic: true,
	}
	return nw
}

func gathernewchrworkinfo(l structs.DbWorkline) structs.DbWork {
	// 15 194 {15 194 -1 -1 -1 -1 1 ch0017w001 ἐνθάδε κῖτ[ε  — <hb-fs-l-normal>name</hb-fs-l-normal> — ] ἐνθάδε κῖτε name ενθαδε κιτε name  region: Sicilia · city: Syracusae · date: date?}
	// this marks the arrival of "Kaibel, IG XIV, 194"; note that the "194" is stored as the l5 value

	// {43 237 -1 -1 -1 -1 1 ch0017w001 ἐνθάδε κῖτε πε͂ϲ̣ <hb-sp-rectified_form>παῖϲ</hb-sp-rectified_form> ἐνθάδε κῖτε πε͂ϲ παῖϲ ενθαδε κιτε πε͂ϲ παιϲ  region: Sicilia · city: Acrae · date: date? · documentnumber: 9}
	// this marks the arrival of "Kaibel, IG XIV, 237"; note that the "237" is stored as the l5 value

	wk := gathernewinsworkinfo(l)
	//wk.Title = l.WkUID // ch0017w001 which we will decode later with an appeal to the IDT (that is: leverage the OLDWORK title)
	//wk.Title = wk.Title + " · " + l.Lvl5Value

	if fnwsimple(l) {
		// "date", etc. triggered this and the wk looks ok
		return wk
	}

	// a mere change in lvl5 value triggered this and the wk title makes no sense
	// ideally you can say that the title is l.Lvl5Value, BUT note that chr0130 can do +1 changes at l5:
	// so you can see 417, 1, 2, 419, ...
	// it turns out that this is 417, 417a, 417b, 419, ...
	// this also means you need fussy code which should have already eliminated this via fixchrlvl5values()

	wk.Title = "#" + l.Lvl5Value

	return wk
}

func gathernewddpworkinfo(l structs.DbWorkline) structs.DbWork {
	wk := gathernewinsworkinfo(l)
	if fnwsimple(l) {
		// "date", etc. triggered this and the wk looks ok
		return wk
	}

	// "<hb-fs-l-normal>PPar 5﹦UPZ 2,180a </hb-fs-l-normal>" was all you were told...
	if ddpnewworkspan.MatchString(l.MarkedUp) {
		tit := ddpnewworkspan.ReplaceAllString(l.MarkedUp, "$1")
		tit = betacode.SimpleLatinSpanLATE(tit)
		wk.Title = tit
	}

	return wk
}

func foundanewwork(l structs.DbWorkline, prevl structs.DbWorkline) bool {
	// should see something like:

	// there are lots of signals for a new work...

	// DDP has provenance, not region...

	// INS0150 will fail simple regionfinder and provenancefinder tests because it marks new documents differently: publicationinfo:

	// INS0140 ...
	// {0 1 -1 -1 -1 -1 1 in0140w001 βαϲιλεὺϲ Βιθυνῶν Ζιαήλαϲ βαϲιλεὺϲ βιθυνῶν ζιαήλαϲ βαϲιλευϲ βιθυνων ζιαηλαϲ  region: Aegean Isl. · city: Kos: Asklepieion · date: 246-242 bc · workabbrev: TAM IV,1 }

	// note that ch0130 changes works *exclusively* by changing the lvl 5 values...
	// here you can see the end of #2, the whole of #3, and the start of #4
	// 38 2 {38 2 -1 -1 -1 -1 1 ch0130w001 sic χ, sic χχ. sic χ sic χχ sic χ sic χχ  }
	// 39 2 {39 2 -1 -1 -1 -1 1 ch0130w001 votis χ, votis χχ. votis χ votis χχ uotis χ uotis χχ  }
	// 40 3 {40 3 -1 -1 -1 -1 39 ch0130w001 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;... quibus omnibus quas[i] quibus omnibus quasi quibus omnibus quasi  }
	// 41 3 {41 3 -1 -1 -1 -1 40 ch0130w001 quidam cumulus accedit, quod omnes quidam cumulus accedit quod omnes quidam cumulus accedit quod omnes  }
	// 42 3 {42 3 -1 -1 -1 -1 41 ch0130w001 [i]bidem sectatores sanctissimae religi- ibidem sectatores sanctissimae religionis ibidem sectatores sanctissimae religionis religionis }
	// 43 3 {43 3 -1 -1 -1 -1 42 ch0130w001 onis habitare dicantur ... habitare dicantur habitare dicantur  }
	// 44 4 {44 4 -1 -1 -1 -1 1 ch0130w001 ddd. ☧ nnn. ddd ☧ nnn ddd ☧ nnn  }

	// contrast ch0017 which both changes those values and has Annotations
	// here you can see the end of #194, the whole of #197 and #198, and the start of #199
	// {20 194 -1 -1 -1 -1 6 ch0017w001 τοῦ ο[ —   —   — ]. τοῦ ο του ο  }
	// {21 197 -1 -1 -1 -1 1 ch0017w001 [ἐνθά]δε κῖτε Ζ[ — <hb-fs-l-normal>name</hb-fs-l-normal> — (﹖)  —  τῆϲ μακαρίαϲ] ἐνθάδε κῖτε ζ name τῆϲ μακαρίαϲ ενθαδε κιτε ζ name τηϲ μακαριαϲ  region: Sicilia · city: Syracusae · date: date?}
	// {22 197 -1 -1 -1 -1 2 ch0017w001 μνήμη[ϲ  —   —   — ] μνήμηϲ μνημηϲ  }
	// {23 197 -1 -1 -1 -1 3 ch0017w001 —   —   — ΗΓ —   —   —  ηγ ηγ  }
	// {24 198 -1 -1 -1 -1 1 ch0017w001 [ἐ]νθάδε [κῖτε  —   —   —  ζή]- ἐνθάδε κῖτε ζήϲαϲ ενθαδε κιτε ζηϲαϲ ζήϲαϲ region: Sicilia · city: Syracusae · date: date?}
	// {25 198 -1 -1 -1 -1 2 ch0017w001 [ϲ]α⟨ϲ⟩ ἔτη [ — ʹ· ἐτελεύτηϲεν(﹖) τῇ πρὸ  — ʹ] ἔτη ἐτελεύτηϲεν τῇ πρὸ ετη ετελευτηϲεν τη προ  }
	// {26 198 -1 -1 -1 -1 3 ch0017w001 [κ]αλ(ανδῶν) [ —   — ων]. καλανδῶν ων καλανδων ων  }
	// {27 199 -1 -1 -1 -1 1 ch0017w001 [ —   —   — ο]υ μακαρ[ι —   —   — ] ου μακαρι ου μακαρι  region: Sicilia · city: Syracusae · date: date?}

	// one more set of irregularities in DDP:
	// No WkUID for 0 <hb-fs-l-normal>PPetr 1,intro,pg 43﹦WChr 55 </hb-fs-l-normal>
	//No WkUID for 1 <hb-fs-l-normal>PPetr 1,11,ctr﹦PPetrWills 18</hb-fs-l-normal>
	//No WkUID for 2 <hb-fs-l-normal>PPetr 1,11,Fr,pg 34﹦PPetr 3,54,A2</hb-fs-l-normal>
	//No WkUID for 3 <hb-fs-l-normal>PPetr 1,11,Fr1,pg 35﹦PPetrWills 16</hb-fs-l-normal>
	//No WkUID for 4 <hb-fs-l-normal>PPetr 1,11,Fr2,pg 35﹦PPetrWills 18</hb-fs-l-normal>
	//No WkUID for 5 <hb-fs-l-normal>PPetr 1,12﹦PPetrWills 13 </hb-fs-l-normal>
	//No WkUID for 6 <hb-fs-l-normal>PPetr 1,13,1﹢2﹢3﹦PPetrWills 6</hb-fs-l-normal>
	//No WkUID for 7 <hb-fs-l-normal>PPetr 1,13,Fr,pg 42﹦PPetrWills 16</hb-fs-l-normal>
	//No WkUID for 8 <hb-fs-l-normal>PPetr 1,14﹦PPetrWills 3 </hb-fs-l-normal>
	//No WkUID for 9 <hb-fs-l-normal>PPetr 1,15﹦PPetrWills 3 </hb-fs-l-normal>
	//No WkUID for 10 <hb-fs-l-normal>PPetr 1,16,Fr1﹦PPetrWills 3</hb-fs-l-normal>
	//No WkUID for 0 <hb-fs-l-normal>OWilck 1﹦OLeid 19 </hb-fs-l-normal>
	//No WkUID for 1 <hb-fs-l-normal>OWilck 2﹦OLeid 175 </hb-fs-l-normal>
	//No WkUID for 2 <hb-fs-l-normal>OWilck 3﹦OLeid 176 </hb-fs-l-normal>

	kvp := l.GatherMetadata()
	var anything string
	region, _ := kvp["region"]
	city, _ := kvp["city"]
	publicationinfo, _ := kvp["publicationinfo"]
	prov, _ := kvp["provenance"]
	doc, _ := kvp["documentnumber"]
	date, _ := kvp["date"]
	workabbrev, _ := kvp["workabbrev"]

	anything = region + city + publicationinfo + prov + doc + date + workabbrev

	// stupid carve-out for CHR0130; which also increments work numbers irregularly via l5...
	if l.WkUID[0:global.ABBREVLEN] == global.CHRFIRSTPASS && l.Lvl5Value != prevl.Lvl5Value {
		return true
	}

	if l.WkUID[0:global.ABBREVLEN] == global.DDPFIRSTPASS && ddpnewworkspan.MatchString(l.MarkedUp) {
		return true
	}

	// NOTE: there seems to be ONE (and only one?) 'impossible' new author with 0 associated works
	// the word counter will bring this to light: `wlbgrabber() ERROR: relation "dpx02d3" does not exist (SQLSTATE 42P01)`
	// if you dig around via `hgdb=# select * from authors where universalid = 'dpx02d3';` you will see that this is
	// `PEdfou · Vol 3` and it comes from `DDP0044.TXT`
	// debugging that file via `authorprep.WriteWorklineProgress(lines)` will produce a file that is 114 lines long
	// the last line reads: `[140] 0044w003 [8 1 1 1 1 1] 	 <hb-fs-l-normal>PEdfou 3,8﹦SB 6,9302</hb-fs-l-normal>  	 Notes: `
	// and examining `DDP0044.TXT` itself shows at the end `&PEdfou `3,`8%6SB `6,`9302$ ðþ(null)(null)(null)(null)...`
	// so that really does seem to be where the file ends, i.e, with the announcement of a work (that will be mapped
	// to a new author) and then... nothing.

	// so we have this incredibly stupid check...
	if l.WkUID == "0044w003" && global.WorkingOnCorpus == "DDP" {
		return false
	}

	if len(anything) == 0 {
		return false
	}
	return true
}

// fnwsimple - without the chr kludge in it
func fnwsimple(l structs.DbWorkline) bool {
	kvp := l.GatherMetadata()
	var anything string
	region, _ := kvp["region"]
	city, _ := kvp["city"]
	publicationinfo, _ := kvp["publicationinfo"]
	prov, _ := kvp["provenance"]
	doc, _ := kvp["documentnumber"]
	date, _ := kvp["date"]

	anything = region + city + publicationinfo + prov + doc + date

	if len(anything) == 0 {
		return false
	}
	// fmt.Println(l)
	return true
}

func buildremapper(idxfilename string) (string, map[string]string) {
	var where string
	switch idxfilename[0:3] {
	case "INS":
		where = global.Config.InsDir
	case "CHR":
		where = global.Config.ChrDir
	case "DDP":
		where = global.Config.PapDir
	default:
		fmt.Println("Unkown indexfilename prefix:", idxfilename[0:3])
		where = global.Config.InsDir
	}

	idtdata := idtandbin.BINFileLoad(where + idxfilename + ".IDT")
	au, wkk := idtandbin.LoadAuthorIDTData(idtdata)

	// au {0150  Ionia      0  []}
	// wkk [{0150w001 Chios   document           0 0 0 0 false} {0150w002 Didyma   document           0 0 0 0 false} ...]

	mapper := make(map[string]string)

	for _, w := range wkk {
		mapper[w.UID] = w.Title
	}
	return au.IDXname, mapper
}
