//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lexica

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

const (
	LTLEXTESTDATA = `<entryFree id="n9" type="hapax" key="a^bactus2" n="2" opt="n"><orth extent="full" lang="la" opt="n">ăbactus</orth>, <itype opt="n">ūs</itype>, <gen opt="n">m.</gen> <etym opt="n">abigo</etym>, <sense id="n9.0" n="I" level="1" opt="n"><hi rend="ital">a driving away, robbing</hi> (of cattle, vessels, etc.), <bibl n="Perseus:abo:phi,1318,002:20:4" default="NO"><author>Plin.</author> Pan. 20, 4</bibl>.</sense></entryFree>
<entryFree id="n10" type="hapax" key="a^ba^cu^lus" opt="n"><orth extent="full" lang="la" opt="n">ăbăcŭlus</orth>, <itype opt="n">i</itype>, <gen opt="n">m.</gen> <lbl opt="n">dim.</lbl> <etym opt="n">abacus</etym>, <sense id="n10.0" n="I" level="1" opt="n"><hi rend="ital">a small cube or tile of colored glass for making ornamental pavements</hi>, the Gr. <foreign lang="greek">u\buki/skos</foreign>, <bibl n="Perseus:abo:phi,0978,001:36:199" default="NO" valid="yes"><author>Plin.</author> 36, 26, 67, § 199</bibl>. <cb n="ABAL" /></sense></entryFree>
<entryFree id="n11" type="greek" key="a^ba^cus" opt="n"><orth extent="full" lang="la" opt="n">ăbăcus</orth>, <itype opt="n">i</itype> (according to <bibl default="NO"><author>Prisc.</author> 752</bibl> P. also <orth extent="full" lang="la" opt="n">ăbax</orth>, ăcis; cf. <bibl default="NO"><author>id.</author> p. 688</bibl>), <gen opt="n">m.</gen>,=<foreign lang="greek">a)/bac, a^kos</foreign>, prop. <sense id="n11.0" n="I" level="1" opt="n"><hi rend="ital">a square tublet;</hi> hence, in partic., </sense><sense id="n11.1" n="I" level="1" opt="n"> <hi rend="ital">A sideboard, the top of which was made of marble, sometimes of silver, gold, or other precious material, chiefly used for the display of gold and silver vessels</hi>, <bibl n="Perseus:abo:phi,0474,005:4:16:section=35" default="NO" valid="yes"><author>Cic.</author> Verr. 2, 4, 16, § 35</bibl>; <bibl default="NO">2, 4, 25, § 57</bibl>; <bibl n="Perseus:abo:phi,0474,049:5:21:61" default="NO"><author>id.</author> Tusc. 5, 21, 61</bibl>; <bibl n="Perseus:abo:phi,0684,001:L. L. 9:section=46" default="NO"><author>Varr.</author> L. L. 9, § 46</bibl> Mūll.; <bibl n="Perseus:abo:phi,0978,001:37:14" default="NO" valid="yes"><author>Plin.</author> 37, 2, 6, § 14</bibl>; <bibl n="Perseus:abo:phi,1276,001:3:2" default="NO"><author>Juv.</author> 3, 2</bibl>-0-4: <cit><quote lang="la">perh. also called mensae Delphicae,</quote> <bibl n="Perseus:abo:phi,0474,005:4:59" default="NO" valid="yes"><author>Cic.</author> Verr. 2, 4, 59 <hi rend="ital">init.</hi></bibl></cit> Zumpt; <bibl n="Perseus:abo:phi,1294,001:12:67" default="NO"><author>Mart.</author> 12, 67</bibl>. Accord. to <bibl n="Perseus:abo:phi,0914,001:39:6:7" default="NO" valid="yes"><author>Liv.</author> 39, 6, 7</bibl>, and <bibl n="Perseus:abo:phi,0978,001:34:14" default="NO" valid="yes"><author>Plin.</author> 34, 3, 8, § 14</bibl>, Cn. Manlius Vulso flrst brought them from Asia to Rome, B.C. 187, in his triumph over the Galatae; cf. Becker, Gall. 2, p. 258 (2d edit.).—</sense><sense id="n11.2" n="II" level="1" opt="n"> <hi rend="ital">A gaming-board, divided into compurtments</hi>, for playing with dice or counters, <bibl n="Perseus:abo:phi,1348,001:life=nero:22" default="NO"><author>Suet.</author> Ner. 22</bibl>; <bibl default="NO"><author>Macr.</author> S. 1, 5</bibl>.— </sense><sense id="n11.3" n="III" level="1" opt="n"> <hi rend="ital">A counting-table,</hi> covered with sand or dust, and used for arithmetical computation, <bibl n="Perseus:abo:phi,0969,001:1:131" default="NO"><author>Pers.</author> 1, 131</bibl>; <bibl default="NO"><author>App.</author> Mag. p. 284</bibl>; cf. Becker, Gall. 2, p. 65. —</sense><sense id="n11.4" n="IV" level="1" opt="n"> <hi rend="ital">A wooden tray</hi>, <bibl default="NO"><author>Cato</author>, R. R. 10, 4</bibl>.—</sense><sense id="n11.5" n="V" level="1" opt="n"> <hi rend="ital">A painted panel or square compariment in the wall or ceiling of a chamber</hi>, <bibl n="Perseus:abo:phi,1056,001:7:3:10" default="NO" valid="yes"><author>Vitr.</author> 7, 3, 10</bibl>; <bibl n="Perseus:abo:phi,0978,001:33:159" default="NO" valid="yes"><author>Plin.</author> 33, 12, 56, § 159</bibl>; <bibl n="Perseus:abo:phi,0978,001:35:3" default="NO" valid="yes">35, 1, 1, § 3</bibl>, and 35, 6, 13, § 32.—</sense><sense id="n11.6" n="VI" level="1" opt="n"> In architecture, <hi rend="ital">a fiat, square stone on the top of a column</hi>, immediately under the architrare, <bibl n="Perseus:abo:phi,1056,001:3:5:5 sq" default="NO" valid="yes"><author>Vitr.</author> 3, 5, 5 sq.</bibl>; <bibl n="Perseus:abo:phi,1056,001:4:1:11" default="NO" valid="yes">4, 1, 11</bibl> sq.</sense></entryFree>`
)

var (
	entryfinder   = regexp.MustCompile(`<entryFree.*?</entryFree>`)
	headparser1   = regexp.MustCompile(`<entryFree id="(\S*?)" type="(\S*?)" key="(\S*?)">`)
	headparser2   = regexp.MustCompile(`<entryFree id="(\S*?)" type="(\S*?)" key="(\S*?)" n="\S*?">`)
	stripdigits   = regexp.MustCompile(`[0-9]`)
	prelimcleaner = regexp.MustCompile(`<entryFree id=".*?" type=".*?" key=.*?">(.*?)</entryFree>`)
)

func FindLatinLexEntries(lexdata string) []structs.DbLexicon {
	lexdata = delunate.Replace(lexdata)
	lexdata = strings.ReplaceAll(lexdata, ` opt="n"`, ``)
	lexdata = strings.ReplaceAll(lexdata, `extent="full" `, ``)
	lexdata = strings.ReplaceAll(lexdata, ` default="NO"`, ``)

	entries := entryfinder.FindAllString(lexdata, -1)

	latee := make([]structs.DbLexicon, len(entries))
	for i, entry := range entries {
		// translations now handled by extracttranslations()
		// var translations []string
		var rem string
		var de structs.DbLexicon
		entry = reformatltxml(entry)
		usedby := collectauthors(entry)
		// entry, translations = gatherlattransl(entry)
		de, rem = extractlatinheadmaterial(entry)
		latee[i] = extractletinsenses(de, rem)
		latee[i].Usedby = usedby
		// latee[i].Transl = strings.Join(translations, SEPARATOR)
		latee[i].PrelimInfo = prelimcleaner.ReplaceAllString(latee[i].PrelimInfo, "$1")
	}
	return latee
}

func extractlatinheadmaterial(fullentry string) (structs.DbLexicon, string) {
	var entry structs.DbLexicon
	headgroups := headparser1.FindAllStringSubmatch(fullentry, 1)
	remainder := headparser1.ReplaceAllString(fullentry, "")
	if len(headgroups) == 1 {
		entry.IdString = headgroups[0][1]
		entry.EntryType = headgroups[0][2]
		entry.EntryName = headgroups[0][3]
	}
	headgroups = headparser2.FindAllStringSubmatch(fullentry, 1)
	remainder = headparser2.ReplaceAllString(fullentry, "")
	if len(headgroups) == 1 {
		entry.IdString = headgroups[0][1]
		entry.EntryType = headgroups[0][2]
		entry.EntryName = headgroups[0][3]
	}
	entry.EntryMetr = stripdigits.ReplaceAllString(generic.HandleVowelLengths(entry.EntryName), "")
	entry.EntryMetr = generic.LunatesAndUV(entry.EntryMetr)
	entry.EntryName = generic.LunatesAndUV(generic.SuperScriptNumbers(generic.StripVowelLengths(entry.EntryName)))
	return entry, remainder
}

var (
	fmtcit        = regexp.MustCompile("<cit>([^<]*?)</cit>")
	fmtquote      = regexp.MustCompile(`<quote lang="la">([^<]*?)</quote>`)
	fmttrans      = regexp.MustCompile(`<trans><tr>([^<]*?)</tr>(,|)</trans>`)
	fmtgreek      = regexp.MustCompile(`<foreign lang="greek">([^<]*?)</foreign>`)
	ltsensefinder = regexp.MustCompile(`<sense id="(.*?)" n="(.*?)" level="(.*?)">(.*?)</sense>`)
	fmttrans2     = regexp.MustCompile(`<hi rend="ital">([^<]*?)</hi>`)
	lttransfinder = regexp.MustCompile(`<hb-lx-hri>([^<]*?)</hb-lx-hri>`)
	// notatransl - to drop things that are not part of a definition; any "." is tempting, but you definitely lose material that way
	notatransl = []string{
		"a.",
		"A.",
		"abl.",
		"absol.",
		"Absol.",
		"acc.",
		"Acc.",
		"act.",
		"Act.",
		"adj.",
		"Adj.",
		"adjj.",
		"adv.",
		"Adv.",
		"advv.",
		"Aes.",
		"Al.",
		"Am.",
		"Ar.",
		"Ba.",
		"bis",
		"Ca.",
		"Ch.",
		"Comp",
		"comp.",
		"Comp.",
		"conj.",
		"Conj.",
		"Da.",
		"dat.",
		"Dat.",
		"De.",
		"dep.",
		"demonstr.",
		"Di.",
		"Do.",
		"Ep.",
		"Eut.",
		"fem.",
		"Fem.",
		"fin.",
		"fin.—Sup.",
		"fut.",
		"Fut.",
		"Ge.",
		"gen.",
		"Gen.",
		"gerund.",
		"Gn.",
		"h.",
		"imp.",
		"imper.",
		"imperf.",
		"impers.",
		"Impers.",
		"ind.",
		"indef.",
		"Indef.",
		"indic.",
		"inf.",
		"init.",
		"insepar.",
		"interrog.",
		"Interrog.",
		"l.",
		"Li.",
		"Ly.",
		"M.",
		"masc.",
		"Masc.",
		"med.",
		"met",
		"Mi.",
		"Mil.",
		"n.",
		"N.",
		"neutr.",
		"Neutr.",
		"no.",
		"nom.",
		"num.",
		"object-acc.",
		"obj.clause",
		"obj.-clause",
		"P.",
		"Pa.",
		"parag.",
		"part.",
		"Part.",
		"partt.",
		"pass.",
		"Pass.",
		"Pe",
		"per",
		"perf.",
		"Perf.",
		"pers.",
		"person.",
		"Ph.",
		"Pi.",
		"Pl.",
		"pluperf.",
		"Pluperf.",
		"plur.",
		"Plur.",
		"posit.",
		"Posit.",
		"praes.",
		"prep.",
		"Prep.",
		"pres.",
		"pron.",
		"pronn.",
		"prepp.",
		"Ps.",
		"Py.",
		"rel.",
		"Rel.",
		"rel.-clause",
		"relat.",
		"reflex.",
		"S.",
		"Sc.",
		"Si.",
		"sing.",
		"Sing.",
		"So.",
		"subj.",
		"Subj.",
		"subst.",
		"Subst.",
		"substt.",
		"Substt.",
		"sup.",
		"Sup.",
		"Sy.",
		"tempp.",
		"ter",
		"Th.",
		"Thr.",
		"Tr.",
		"ut",
		"ut.",
		"v.",
		"v.a.",
		"verb.",
		"voc.",
	}
)

func reformatltxml(xml string) string {
	const (
		AU = `<hb-lx-au>$1</hb-lx-au>`
		CT = `<hb-lx-cit>$1</hb-lx-cit>`
		QT = `<hb-lx-lt-quote>$1</hb-lx-lt-quote>`
		TR = `<hb-lx-tr>$1</hb-lx-tr>$2`
		GK = "<hb-lx-grk>%s</hb-lx-grk>"
		HR = `<hb-lx-hri>$1</hb-lx-hri>`
	)
	xml = fmtauthor.ReplaceAllStringFunc(xml, lookupltauthor)
	xml = fmtcit.ReplaceAllString(xml, CT)
	xml = fmtquote.ReplaceAllString(xml, QT)
	xml = fmttrans.ReplaceAllString(xml, TR)
	xml = fmttrans2.ReplaceAllString(xml, HR)
	xml = fmtgreek.ReplaceAllStringFunc(xml, formatgreekinlatinlex)

	return xml
}

func lookupltauthor(match string) string {
	const (
		AU = "<hb-lx-au>%s</hb-lx-au>"
	)
	groups := fmtauthor.FindAllStringSubmatch(match, 1)
	if len(groups) > 0 {
		au, ok := LatinAuthors[groups[0][1]]
		if !ok {
			au = groups[0][1]
		}
		return fmt.Sprintf(AU, au)
	}
	return fmt.Sprintf(AU, match)
}

func formatgreekinlatinlex(match string) string {
	const (
		GK = "<hb-lx-grk>%s</hb-lx-grk>"
	)
	groups := fmtgreek.FindAllStringSubmatch(match, 1)
	if len(groups) > 0 {
		unicode := generic.ConvertLCBetacode(groups[0][1])
		match = fmt.Sprintf(GK, unicode)
	}
	return match
}

func extractletinsenses(latlex structs.DbLexicon, prunedentry string) structs.DbLexicon {
	if ltsensefinder.MatchString(prunedentry) {
		latlex, prunedentry = els1(latlex, prunedentry)
	}

	var trr []string

	for _, s := range latlex.Senses {
		latlex.SenseIDs = append(latlex.SenseIDs, s.ID)
		nt := extracttranslations(s.Contents)
		if len(nt) > 0 {
			trr = append(trr, nt)
		}
	}

	var deduptrr []string
	deduper := make(map[string]bool)
	for _, t := range trr {
		if _, present := deduper[t]; present {
			continue
		} else {
			deduper[t] = true
			deduptrr = append(deduptrr, t)
		}
	}

	latlex.PrelimInfo = prunedentry
	latlex.Transl = strings.Join(deduptrr, SEPARATOR)

	return latlex
}

// els1 - helper for extractletinsenses; depends on the order of the xml fields
func els1(latentry structs.DbLexicon, prunedentry string) (structs.DbLexicon, string) {
	headgroups := ltsensefinder.FindAllStringSubmatch(prunedentry, -1)
	remainder := ltsensefinder.ReplaceAllString(prunedentry, "")

	count := 100
	for _, headgroup := range headgroups {
		var newsense structs.LexicalSenses
		newsense.ID = headgroup[1]
		newsense.N = headgroup[2]
		newsense.LVL = headgroup[3]
		newsense.Contents = headgroup[4]
		if newsense.ID == "" {
			// never hitting this...
			newsense.ID = latentry.IdString + fmt.Sprintf(".%d", count)
			count++
			fmt.Println("no sense id for", latentry.IdString)
		}
		latentry.Senses = append(latentry.Senses, newsense)
	}
	return latentry, remainder
}

// extracttranslations - extract translations from a "sense" so you can set the Transl field of a structs.DbLexicon
func extracttranslations(s string) string {
	// this is tricky since the markup is the same for grammar and translations; see "perolesco", for example,
	// where "hri" sets off both "it's a verb" and "it means x":
	// `<hb-lx-hri>v. inch. n.</hb-lx-hri>, <hb-lx-hri>to grow up</hb-lx-hri>, Lucil. ap. <notbibl><hb-lx-au>Prisc.</hb-lx-au> p. 872</notbibl> P.`

	groups := lttransfinder.FindAllStringSubmatch(s, -1)
	var senses []string
	if len(groups) > 0 {
		for _, group := range groups {
			elem := strings.Split(group[1], " ")
			if !generic.SliceOverlap(notatransl, elem) {
				g1 := strings.TrimSpace(group[1])
				g1 = strings.TrimPrefix(g1, ", ")
				g1 = strings.TrimSuffix(g1, ".") // note that this will make debugging hard may need to be toggled
				if len(group[1]) > 1 {
					senses = append(senses, g1)
				}
			}
		}
	}
	return strings.Join(senses, ", ")
}
