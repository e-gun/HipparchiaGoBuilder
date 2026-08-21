//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lexica

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	d2headparser   = regexp.MustCompile(`<div2 id="(.*?)" orig_id="(.*?)" key="(.*?)" type="(.*?)">`)
	gksensefinder1 = regexp.MustCompile(`<sense id="([^"]*?)" n="([^"]*?)" level="([^"]*?)">(.*?)</sense>`)
	gksensefinder2 = regexp.MustCompile(`<sense n="([^"]*?)" id="([^"]*?)" level="([^"]*?)">(.*?)</sense>`)
	gksensefinder3 = regexp.MustCompile(`<sense n="" level="([^"]*?)">(.*?)</sense>`)
	lsjheadfinder  = regexp.MustCompile(`<head.*?orth_orig="(.*?)">(.*?)</head>`)
	delunate       = strings.NewReplacer("ς", "ϲ", "σ", "ϲ", "Σ", "Ϲ")
	fmtauthor      = regexp.MustCompile("<author>(.*?)</author>")
)

const (
	TESTENTRY = `<div2 id="crossu(pno/w" orig_id="n108086" key="u(pno/w" type="main" opt="n"><head extent="full" lang="greek" opt="n" orth_orig="ὑπν-όω">ὑπνόω</head>, <tns opt="n">fut.</tns> <cit><quote lang="greek">-ώσω</quote> <bibl n="Perseus:abo:tlg,4080,001:18:14:3" default="NO"><title>Gp.</title> 18.14.3</bibl></cit>: <tns opt="n">aor.</tns> <cit><quote lang="greek">ὕπνωσα</quote> <bibl n="Perseus:abo:tlg,0627,006:3:1:3" default="NO"><author>Hp.</author> <title>Epid.</title> 3.1.γʹ</bibl></cit>, <bibl n="Perseus:abo:tlg,0543,001:3:81:5"><author>Plb.</author> 3.81.5</bibl>, <bibl n="Perseus:abo:tlg,0527,034:46:20"><author>LXX</author> <title>Si.</title> 46.20</bibl>, <bibl n="Perseus:abo:tlg,0526,001:1:12:1"><author>J.</author> <title>AJ</title> 1.12.1</bibl>, <bibl n="Perseus:abo:tlg,0007,047:76" default="NO"><author>Plu.</author> <title>Alex.</title> 76</bibl>, etc.: <tns opt="n">pf.</tns> <cit><quote lang="greek">ὕπνωκα</quote> <bibl n="Perseus:abo:tlg,0007,082:236b"><author>Id.</author> <title>Apophth. Lac.</title> 2.236b</bibl></cit>, (<etym lang="greek" opt="n">καθ-</etym>) <bibl n="Perseus:abo:tlg,0526,001:5:9:3"><author>J.</author> <title>AJ</title> 5.9.3</bibl>:—<gramGrp opt="n"><gram type="voice" opt="n">Med.</gram></gramGrp>, <tns opt="n">fut.</tns> <foreign lang="greek">ὑπνώσομαι</foreign> ibid.:—<gramGrp opt="n"><gram type="voice" opt="n">Pass.</gram></gramGrp>, <tns opt="n">pf.</tns> <mood opt="n">part.</mood> <cit><quote lang="greek">ὑπνωμένος</quote> <bibl n="Perseus:abo:tlg,0016,001:1:11"><author>Hdt.</author> 1.11</bibl></cit>, <bibl n="Perseus:abo:tlg,0016,001:3:69">3.69</bibl>: <tns opt="n">aor.</tns> <cit><quote lang="greek">ὑπνώθην</quote> <bibl n="Perseus:abo:tlg,0007,085:313a"><author>Plu.</author> <title>Par.min.</title> 2.313a</bibl></cit>:—<sense id="n108086.0" n="" level="1" opt="n"><i>put to sleep,</i> only in <author>Dsc.</author> 4.63:—<gramGrp opt="n"><gram type="voice" opt="n">Pass.</gram></gramGrp>, <i>fall asleep, sleep,</i> <author>Hdt.</author> ll.cc.:—so in <gramGrp opt="n"><gram type="voice" opt="n">Med.</gram></gramGrp>, <author>J.</author> l.c. </sense><sense n="II" id="n108086.1" level="2" opt="n"> intr., like <gramGrp opt="n"><gram type="voice" opt="n">Pass.</gram></gramGrp>, <bibl n="Perseus:abo:tlg,0627,006:3:1:3" default="NO"><author>Hp.</author> <title>Epid.</title> 3.1.γʹ</bibl>, <bibl n="Perseus:abo:tlg,0627,006:7:11" default="NO">7.11</bibl> (<foreign lang="greek">ὑπνώσσουσα</foreign> Littré, with cod. C), <bibl n="Perseus:abo:tlg,0086,042:454a:2" default="NO"><author>Arist.</author> <title>Somn.Vig.</title> 454a2</bibl>, <bibl n="Perseus:abo:tlg,0086,051:10" default="NO"><title>Fr.</title> 10</bibl>, <bibl n="Perseus:abo:tlg,0526,001:1:12:1"><author>J.</author> <title>AJ</title> 1.12.1</bibl>; <gramGrp opt="n"><gram type="dialect" opt="n">Lacon.</gram></gramGrp> <mood opt="n">inf.</mood> <foreign lang="greek">ὑπνῶν,</foreign> for <foreign lang="greek">-οῦν,</foreign> <bibl n="Perseus:abo:tlg,0019,007:143"><author>Ar.</author> <title>Lys.</title> 143</bibl>. </sense><sense n="III" id="n108086.2" level="2" opt="n"> <i>die,</i> <author>LXX</author> l.c. (Cf. <foreign lang="greek">ὑπνώω.</foreign>) </sense></div2>`
	SEPARATOR = ` ‖ `
)

// FindGreekLexEntries - take a whole xml file and turn it into []DbLexicon (where each entry contains []LexicalSenses)
func FindGreekLexEntries(lexdata string) ([]structs.DbLexicon, error) {
	lexdata = delunate.Replace(lexdata)

	// be careful that this does not have knock-on effects AND recognize that copy/paste of the original XML just got
	// complicated; but reformatgkxml() will not work until this has been done

	lexdata = strings.ReplaceAll(lexdata, ` opt="n"`, ``)
	lexdata = strings.ReplaceAll(lexdata, ` default="NO">`, `>`)

	entries, spliterror := gklexsplitter(lexdata)
	// sample entry:
	// <div2 id="crossu(ywth/s" orig_id="n110025" key="u(ywth/s" type="main" opt="n"><head extent="suff" lang="greek" opt="n" orth_orig="ὑψ-ωτής">ὑψωτής</head>, <itype lang="greek" opt="n">οῦ</itype>, <gen lang="greek" opt="n">ὁ</gen>, <sense id="n110025.0" n="" level="1" opt="n"><i>one who exalts,</i> <title>PMag.Leid.V.</title> 7.11 (pl.).</sense>

	lsjee := make([]structs.DbLexicon, len(entries))
	for i, entry := range entries {
		var translations []string
		entry = reformatgkxml(entry)
		usedby := collectauthors(entry)
		entry, translations = reformatandextracttranslations(entry)

		var deduptrr []string
		deduper := make(map[string]bool)
		for _, t := range translations {
			if _, present := deduper[t]; present {
				continue
			} else {
				deduper[t] = true
				deduptrr = append(deduptrr, t)
			}
		}

		en, tx := extractdiv2material(entry)
		en, tx = extractheadmaterial(en, tx)
		lsjee[i] = extractlsjsenses(en, tx)
		lsjee[i].Transl = strings.Join(deduptrr, SEPARATOR)
		lsjee[i].Usedby = usedby
	}
	return lsjee, spliterror
}

// gklexsplitter - (1) split the xml data into []string with each string being an entry
func gklexsplitter(lexdata string) ([]string, error) {
	// sample entry:
	// <div2 id="crossu(ywth/s" orig_id="n110025" key="u(ywth/s" type="main" opt="n"><head extent="suff" lang="greek" opt="n" orth_orig="ὑψ-ωτής">ὑψωτής</head>, <itype lang="greek" opt="n">οῦ</itype>, <gen lang="greek" opt="n">ὁ</gen>, <sense id="n110025.0" n="" level="1" opt="n"><i>one who exalts,</i> <title>PMag.Leid.V.</title> 7.11 (pl.).</sense>
	// header ends:
	//  </item></change></revisionDesc></teiHeader>
	//<text>
	//<div1>

	// then `<div2 id=... >.*</div2>` is an entry
	gle := strings.Split(lexdata, "<div1>")
	if len(gle) != 2 {
		return []string{}, fmt.Errorf("gklexsplitter lexdata error")
	}
	div2s := gle[1]
	gle = strings.Split(div2s, "</div2>")
	return gle, nil
}

// extractdiv2material - (2) take an individual entry string and start building an DbLexicon by extracting ID, etc
func extractdiv2material(fullentry string) (structs.DbLexicon, string) {
	// sample div2:
	// <div2 id="crossu(ywth/s" orig_id="n110025" key="u(ywth/s" type="main" opt="n">

	var entry structs.DbLexicon

	headgroups := d2headparser.FindAllStringSubmatch(fullentry, 1)
	remainder := d2headparser.ReplaceAllString(fullentry, "")

	if len(headgroups) > 0 {
		// entry.(nullfield) = headgroups[0][1] // not useful...
		entry.IdString = headgroups[0][2]
		// "head" has the same info as headgroups[0][3]?
		entry.EntryName = generic.ConvertLCBetacode(headgroups[0][3]) // "key" is the word... IDStr might have "cross" in front of it
		entry.EntryType = headgroups[0][4]
		if entry.IdString == "" {
			// never hitting this...
			fmt.Println("no entry ID value found for", entry.EntryName)
			entry.IdString = entry.EntryName
		}
	}

	return entry, remainder
}

// extractheadmaterial - (3) extract "orth_orig"
func extractheadmaterial(lsj structs.DbLexicon, prunedentry string) (structs.DbLexicon, string) {
	headgroups := lsjheadfinder.FindAllStringSubmatch(prunedentry, 1)
	remainder := lsjheadfinder.ReplaceAllString(prunedentry, "")
	if len(headgroups) > 0 {
		lsj.EntryMetr = headgroups[0][1]
		lsj.EntryName = headgroups[0][2]
	}
	return lsj, remainder
}

// extractlsjsenses - (4) take an individual entry string and flesh out its DbLexicon with senses; what remains of the string is PrelimInfo
func extractlsjsenses(lsj structs.DbLexicon, prunedentry string) structs.DbLexicon {
	// sample:
	// <sense id="n108066.0" n="" level="1" opt="n"><i>take upon oneself,</i> <abbr>i.e.</abbr> ... <foreign lang="greek">τὰ στύφειν ὑπισχνούμενα</foreign> ib. <bibl n="Perseus:abo:tlg,0565,001:120" default="NO">120</bibl>. </sense>

	// done above; but can uncomment to test just this function
	// prunedentry = reformatgkxml(prunedentry)

	if gksensefinder1.MatchString(prunedentry) {
		lsj, prunedentry = egs1(lsj, prunedentry)
	}

	if gksensefinder2.MatchString(prunedentry) {
		lsj, prunedentry = egs2(lsj, prunedentry)
	}

	if gksensefinder3.MatchString(prunedentry) {
		lsj, prunedentry = egs3(lsj, prunedentry)
	}

	for _, s := range lsj.Senses {
		lsj.SenseIDs = append(lsj.SenseIDs, s.ID)
	}

	lsj.PrelimInfo = prunedentry

	return lsj
}

// egs1 - helper for extractlsjsenses; depends on the order of the xml fields
func egs1(lsj structs.DbLexicon, prunedentry string) (structs.DbLexicon, string) {
	headgroups := gksensefinder1.FindAllStringSubmatch(prunedentry, -1)
	remainder := gksensefinder1.ReplaceAllString(prunedentry, "")

	count := 100
	for _, headgroup := range headgroups {
		var newsense structs.LexicalSenses
		newsense.ID = headgroup[1]
		newsense.N = headgroup[2]
		newsense.LVL = headgroup[3]
		newsense.Contents = headgroup[4]
		if newsense.ID == "" {
			// never hitting this...
			newsense.ID = lsj.IdString + fmt.Sprintf(".%d", count)
			count++
			fmt.Println("no sense id for", lsj.IdString)
		}
		lsj.Senses = append(lsj.Senses, newsense)
	}
	return lsj, remainder
}

// egs2 - helper for extractlsjsenses; depends on the order of the xml fields
func egs2(lsj structs.DbLexicon, prunedentry string) (structs.DbLexicon, string) {
	headgroups := gksensefinder2.FindAllStringSubmatch(prunedentry, -1)
	remainder := gksensefinder2.ReplaceAllString(prunedentry, "")

	for _, headgroup := range headgroups {
		var newsense structs.LexicalSenses
		newsense.N = headgroup[1]
		newsense.ID = headgroup[2]
		newsense.LVL = headgroup[3]
		newsense.Contents = headgroup[4]
		lsj.Senses = append(lsj.Senses, newsense)
	}
	return lsj, remainder
}

func egs3(lsj structs.DbLexicon, prunedentry string) (structs.DbLexicon, string) {
	// this time you do not have enough info so you have to patch it together
	// only one match?!
	// egs3 n29171

	// fmt.Println("egs3", lsj.IdString)
	headgroups := gksensefinder3.FindAllStringSubmatch(prunedentry, -1)
	remainder := gksensefinder3.ReplaceAllString(prunedentry, "")

	count := 100
	var priorsense structs.LexicalSenses
	if len(lsj.Senses) == 0 {
		priorsense = structs.LexicalSenses{
			ID: lsj.IdString + fmt.Sprintf(".%d", count),
			N:  fmt.Sprintf("%d", count),
		}
		count++
	} else {
		priorsense = lsj.Senses[len(lsj.Senses)-1]
	}

	for _, headgroup := range headgroups {
		var newsense structs.LexicalSenses
		newsense.N = fmt.Sprintf("%d", count)
		newsense.ID = priorsense.ID + fmt.Sprintf(".%d", count)
		newsense.LVL = headgroup[1]
		newsense.Contents = headgroup[2]
		lsj.Senses = append(lsj.Senses, newsense)
		count++
	}
	return lsj, remainder
}

func reformatgkxml(xml string) string {
	const (
		AU = `<hb-lx-au>$1</hb-lx-au>`
		TI = `<hb-lx-tit>$1</hb-lx-tit>`
		TN = `<hb-lx-tns>$1</hb-lx-tns>`
		MD = `<hb-lx-md>$1</hb-lx-md>`
		LG = `<hb-lx-lg-$1>$2</hb-lx-lg-$1>`
		GR = `<hb-lx-gr-$1>$2</hb-lx-gr-$1>`
	)
	var (
		// fmtauthor = regexp.MustCompile("<author>(.*?)</author>")
		fmttit   = regexp.MustCompile("<title>(.*?)</title>")
		fmttns   = regexp.MustCompile("<tns>(.*?)</tns>")
		fmtmood  = regexp.MustCompile("<mood>(.*?)</mood>")
		fmtflang = regexp.MustCompile(`<foreign lang="(.*?)">(.*?)</foreign>`)
		fmtgram  = regexp.MustCompile(`<gram type="(.*?)">(.*?)</gram>`)
	)

	xml = fmtauthor.ReplaceAllStringFunc(xml, lookupgkauthor)
	xml = fmttit.ReplaceAllString(xml, TI)
	xml = fmttns.ReplaceAllString(xml, TN)
	xml = fmtmood.ReplaceAllString(xml, MD)
	xml = fmtflang.ReplaceAllString(xml, LG)
	xml = fmtgram.ReplaceAllString(xml, GR)

	return xml
}

func lookupgkauthor(match string) string {
	const (
		AU = "<hb-lx-au>%s</hb-lx-au>"
	)
	groups := fmtauthor.FindAllStringSubmatch(match, 1)
	if len(groups) > 0 {
		au, ok := GreekAuthors[groups[0][1]]
		if !ok {
			au = groups[0][1]
		}
		return fmt.Sprintf(AU, au)
	}
	return fmt.Sprintf(AU, match)
}

func reformatandextracttranslations(sensebody string) (string, []string) {
	const (
		TR = `<hb-lx-tr>$1</hb-lx-tr>`
	)
	var (
		findtransl = regexp.MustCompile("<i>(.*?)</i>")
	)

	var meanings []string
	alltrans := findtransl.FindAllString(sensebody, -1)
	for _, tr := range alltrans {
		cln := findtransl.ReplaceAllString(tr, "$1")
		cln = strings.TrimSuffix(cln, ",")
		meanings = append(meanings, cln)
	}
	// fmt.Printf("alltrans: %v\n", meanings)
	sensebody = findtransl.ReplaceAllString(sensebody, TR)

	return sensebody, meanings
}

func collectauthors(xml string) string {
	var (
		findrewrittenauth = regexp.MustCompile("<hb-lx-au>(.*?)</hb-lx-au>")
	)
	// must run after reformatgkxml()
	var authors []string
	groups := findrewrittenauth.FindAllStringSubmatch(xml, -1)
	for _, group := range groups {
		authors = append(authors, group[1])
	}
	authors = generic.Unique(authors)
	slices.Sort(authors)
	return strings.Join(authors, SEPARATOR)
}
