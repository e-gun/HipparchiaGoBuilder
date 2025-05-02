package lexica

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
	"strings"
)

var (
	neposstringmap = map[string]string{
		"Ages.":  "Ag",
		"Alcib.": "Alc",
		"Arist.": "Ar",
		// "Att",
		"Cat.":   "Ca",
		"Cato,":  "Ca",
		"Cato":   "Ca",
		"Chab.":  "Cha",
		"Chabr.": "Cha",
		"Chabr,": "Cha",
		// "Cim",
		// "Con",
		"Datam.":  "Dat",
		"Dion.":   "Di",
		"Dion,":   "Di",
		"Epam.":   "Ep",
		"Eun.":    "Eum",
		"Eun,":    "Eum",
		"Hamilc.": "Ham",
		"Hann.":   "Han",
		"Hannib.": "Han",
		"Iphicr.": "Iph",
		"Iphic.":  "Iph",
		"Lyt.":    "Lys",
		"Mil.":    "Milt",
		// "Paus",
		"Pelop.": "Pel",
		"Ph.":    "Phoc",
		"Regg.":  "Reg",
		"Th.":    "Them",
		"Thras.": "Thr",
		"Tim.":   "Timol",
		// "Timoth"
	}
	senecamap = map[string][]string{
		"ad Helv.":          {"1017", "012:12:"},
		"ad Marc. Consol.":  {"1017", "012:6:"},
		"ad Marc.":          {"1017", "012:6:"},
		"adv. Marc.":        {"1017", "012:6:"}, // ?
		"ad Polyb.":         {"1017", "012:11:"},
		"Agam.":             {"1017", "007:"},
		"Agm.":              {"1017", "007:"},
		"Apoc.":             {"1017", "011:"},
		"Apocol.":           {"1017", "011:"},
		"Apocol.p.":         {"1017", "011:"},
		"Ben.":              {"1017", "013:"},
		"Benef.":            {"1017", "013:"},
		"Brev. Vit.":        {"1017", "012:10:"},
		"Clem.":             {"1017", "014:"},
		"Cons ad Helv.":     {"1017", "012:12:"},
		"Cons. ad Helv.":    {"1017", "012:12:"},
		"Cons. ad Marc.":    {"1017", "012:6:"},
		"Cons. ad Pol.":     {"1017", "012:11:"},
		"Cons. ad Polyb.":   {"1017", "012:11:"},
		"Cons. Helv.":       {"1017", "012:12:"},
		"Cons. Marc.":       {"1017", "012:6:"},
		"Cons. Poll.":       {"1017", "012:11:"},
		"Cons. Polyb.":      {"1017", "012:11:"},
		"Cons. Sap.":        {"1017", "012:2:"},
		"Consol. ad Marc.":  {"1017", "012:6:"},
		"Consol. ad Polyb.": {"1017", "012:11:"},
		"Const. Sap.":       {"1017", "012:2:"},
		"Const.":            {"1017", "012:2:"},
		"Contr":             {"1014", "001:"}, // the father...
		"Contr.":            {"1014", "001:"}, // the father...
		"Controv.":          {"1014", "001:"}, // the father...
		"de Ben.":           {"1017", "013:"},
		"de Clem.":          {"1017", "014:"},
		"de Const. Sap.":    {"1017", "012:2:"},
		"de Ira.":           {"1017", "012:3:"},
		"de Prov.":          {"1017", "012:1:"},
		"de Vit. Beat.":     {"1017", "012:7:"},
		"Ep":                {"1017", "015:"},
		"Ep.":               {"1017", "015:"},
		"Exc. Contr.":       {"1014", "002:"}, // the father...
		"Exc. Controv.":     {"1014", "002:"}, // the father...
		"Excerpt. Contr.":   {"1014", "002:"}, // the father...
		"Excerpt. Controv.": {"1014", "002:"}, // the father...
		"Helv.":             {"1017", "012:12:"},
		"Helv. Cons.":       {"1017", "012:12:"},
		"Her. Fur.":         {"1017", "001:"},
		"Herc Fur.":         {"1017", "001:"},
		"Herc Oet.":         {"1017", "009:"},
		"Herc. F.":          {"1017", "001:"},
		"Herc. Fur.":        {"1017", "001:"},
		"Herc. Oet.":        {"1017", "009:"},
		"Here. Fur.":        {"1017", "001:"},
		"Hipp.":             {"1017", "005:"},
		"Hippol.":           {"1017", "005:"},
		"Ira,":              {"1017", "012:3:"},
		"Ira.":              {"1017", "012:3:"}, // note that we just sent I, II, and III to the same place...
		"Lud. Mort. Claud.": {"1017", "011:"},
		"Med.":              {"1017", "004:"},
		"Mort. Claud.":      {"1017", "011:"},
		"N. Q.":             {"1017", "016:"},
		"Oct.":              {"1017", "010:"},
		"Octav.":            {"1017", "010:"},
		"Oed.":              {"1017", "006:"},
		"Oedip":             {"1017", "006:"},
		"Oedip.":            {"1017", "006:"},
		"Oet.":              {"1017", "009:"},
		"Ot. Sap.":          {"1017", "012:8:"},
		"Phaedr.":           {"1017", "005:"},
		"Phoen.":            {"1017", "003:"},
		"Polyb.":            {"1017", "012:11:"},
		"Prov":              {"1017", "012:1:"},
		"Prov.":             {"1017", "012:1:"},
		"Q. N":              {"1017", "016:"},
		"Q. N.":             {"1017", "016:"},
		"Suas.":             {"1014", "003:"}, // the father...
		"Thyest":            {"1017", "008:"},
		"Thyest.":           {"1017", "008:"},
		"Tranq. An.":        {"1017", "012:9:"},
		"Tranq. Vit.":       {"1017", "012:9:"},
		"Tranq.":            {"1017", "012:9:"},
		"Troad.":            {"1017", "002"},
		"Vit. B.":           {"1017", "012:7:"},
		"Vit. Beat.":        {"1017", "012:7:"},
		// "Orest.": {"1017", "006:"},  // ?!
	}
	sallustringmap = map[string]string{
		"C.":    "001",
		"Cat.":  "001",
		"J":     "002",
		"J.":    "002",
		"Jug.":  "002",
		"H.":    "003",
		"Hist.": "003",
	}
	suetostringmap = map[string]string{
		"aug.":  "Aug",
		"cal.":  "Cal",
		"cl.":   "Cl",
		"dom.":  "Dom",
		"gal.":  "Gal",
		"jul.":  "Jul",
		"nero":  "Nero",
		"otho":  "Otho",
		"tib.":  "Tib",
		"tit.":  "Tit",
		"vesp.": "Ves", // the one that breaks the pattern
		"vit.":  "Vit",
	}
	citationcleaner = regexp.MustCompile(`[^0-9:]`)
)

func FixSpecificLatinCitations(lex string) string {
	lex = fixibidem(lex) // do me first...
	lex = purgesq(lex)
	lex = fixciceroverrines(lex) // needs to be first cicero fix
	lex = fixcicerosections(lex)
	lex = fixcicerochapters(lex)
	lex = fixfrontinus(lex)
	lex = fixmartial(lex)
	lex = fixnepos(lex)
	lex = fixpropertius(lex)
	lex = fixsallust(lex)
	lex = fixseneca(lex)
	lex = fixsuetonius(lex)
	lex = fixvarro(lex)
	lex = purgelinklessbibls(lex) // do last
	return lex
}

func purgesq(lex string) string {
	// <bibl n="Perseus:abo:phi,0472,001:61:130 sq" default="NO" valid="yes"><author>Cat.</author> 61, 130 sq.</bibl>
	// get `sq` out of the link
	const (
		TMPL = `"Perseus:abo:phi,%s:%s"`
	)
	sqfinder := regexp.MustCompile(`"Perseus:abo:phi,(\d\d\d\d,\d\d\d):([^"]*) sq"`)
	return sqfinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := sqfinder.FindStringSubmatch(s)
		return fmt.Sprintf(TMPL, groups[1], groups[2])
	})
}

func fixibidem(lex string) string {
	const (
		TMPL = `"Perseus:abo:phi,%s:%s"`
	)
	ibidemfinder := regexp.MustCompile(`"Perseus:abo:phi,(\d\d\d\d,\d\d\d):(ib\. )([^"]*)"`)
	// "Perseus:abo:phi,1014,001:ib. 76:2" --> "Perseus:abo:phi,1014,001:76:2"
	return ibidemfinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := ibidemfinder.FindStringSubmatch(s)
		return fmt.Sprintf(TMPL, groups[1], groups[3])
	})
}

func fixciceroverrines(lex string) string {

	// grep "Perseus:abo:phi,0474,005" * | grep "section="
	// <bibl n="Perseus:abo:phi,0474,005:2:58:section=142" default="NO" valid="yes"><author>id.</author> Verr. 2, 2, 58, § 142</bibl>
	// ( pecunias capere possint: 2.2.142)
	// <bibl n="Perseus:abo:phi,0474,005:4:55:section=123" default="NO" valid="yes"><author>id.</author> Verr. 2, 4, 55, § 123</bibl>
	// <bibl n="Perseus:abo:phi,0474,005:5:65:section=168" default="NO" valid="yes"><author>Cic.</author> Verr. 2, 5, 65, § 168</bibl>
	// <bibl n="Perseus:abo:phi,0474,005:4:28:section=64" default="NO" valid="yes"><author>id.</author> Verr. 2, 4, 28, § 64</bibl>

	// do we never cite Actio I?
	// yes: and it is indisting. from Actio II...
	// <bibl n="Perseus:abo:phi,0474,005:11:32" default="NO" valid="yes"><author>Cic.</author> Verr. 1, 11, 32</bibl>

	const (
		TMPL = `"Perseus:abo:phi,0474,005:%s:%s:%s" default="NO" valid="yes"><author>%s</author> Verr. %s,`
	)
	verrinefinder := regexp.MustCompile(`"Perseus:abo:phi,0474,005:(\d+):\d+:section=(\d+)" default="NO" valid="yes"><author>([^<]*)</author> Verr. (\d),`)
	lex = verrinefinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := verrinefinder.FindStringSubmatch(s)
		actio := "2"
		if groups[4] == "1" {
			actio = "1"
		}
		return fmt.Sprintf(TMPL, actio, groups[1], groups[2], groups[3], groups[4])
	})
	return lex
}

func fixcicerosections(lex string) string {
	// 	RUN THE VERRINES FIRST (because it needs 'section=')
	//
	//	<quote lang="la">omnes de tuā virtute commemorant,</quote> <bibl n="Perseus:abo:phi,0474,058:1:1:13:section=37" default="NO" valid="yes"><author>Cic.</author> Q. Fr. 1, 1, 13, § 37</bibl>
	//
	//	this should just be 1.1.37
	//
	//	the problem is confined almost exclusively to lt0474w058
	const (
		TMPL = `"Perseus:abo:phi,0474,%s:%s"`
	)
	cicsectionfinder := regexp.MustCompile(`"Perseus:abo:phi,0474,(.*):(\d+):section=(\d+)"`)
	lex = cicsectionfinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := cicsectionfinder.FindStringSubmatch(s)
		return fmt.Sprintf(TMPL, groups[1], groups[3])
	})
	return lex
}

func fixcicerochapters(lex string) string {
	// 	RUN THE VERRINES FIRST (because it needs 'section=')
	//
	// <bibl n="Perseus:abo:phi,0474,039:chapter=34" default="NO"><author>Cic.</author> Brut. 34</bibl>: accusatione desistere
	// <bibl n="Perseus:abo:phi,0474,010:chapter=3" default="NO" valid="yes"><author>Cic.</author> Clu. 3</bibl>
	// of course `accusationi respondere` is in Clu 8...

	// so we want the sections and not the chapters; but we do not have them
	// i.e., this is going to produce inaccurate links...
	// purgelinklessbibls() will clear out any `<bibl default="NO">`

	const (
		TMPL = `<bibl default="NO">`
	)

	cicsectionfinder1 := regexp.MustCompile(`<bibl n="Perseus:abo:phi,0474,(.*):chapter=(\d+)" default="NO">`)
	cicsectionfinder2 := regexp.MustCompile(`<bibl n="Perseus:abo:phi,0474,(.*):chapter=(\d+)" default="NO" valid="yes">`)
	//lex = cicsectionfinder.ReplaceAllStringFunc(lex, func(s string) string {
	//	groups := cicsectionfinder.FindStringSubmatch(s)
	//	return fmt.Sprintf(TMPL, groups[1], groups[2])
	//})
	lex = cicsectionfinder1.ReplaceAllString(lex, TMPL)
	return cicsectionfinder2.ReplaceAllString(lex, TMPL)
}

func fixfrontinus(lex string) string {
	// "Perseus:abo:phi,1245,001:Aquaed. 104"
	//	but Aq. is 002
	frontfinder := regexp.MustCompile(`"Perseus:abo:phi,1245,001:Aquaed. ([^"]*)"`)
	return frontfinder.ReplaceAllString(lex, `"Perseus:abo:phi,1245,002:$1"`)
}

func fixmartial(lex string) string {
	// all of martial has been assigned to work 001, but epig is 002 and de spect is 001
	// it is not clear that de spect is ever cited

	return strings.ReplaceAll(lex, "Perseus:abo:phi,1294,001", "Perseus:abo:phi,1294,002")
}

func fixnepos(lex string) string {
	// <bibl n="Perseus:abo:phi,0588,001:Alcib. 11:4" default="NO" valid="yes"><author>Nep.</author> Alcib. 11, 4</bibl></cit>:
	// --> <bibl n="Perseus:abo:phi,0588,001:Alc:11:4" default="NO" valid="yes"><author>Nep.</author> Alcib. 11, 4</bibl>
	const (
		REPL = `"Perseus:abo:phi,0588,001:%s:%s"`
	)
	neposfinder := regexp.MustCompile(`"Perseus:abo:phi,0588,001:([^\s]*\.) ([^"]*)"`)
	lex = neposfinder.ReplaceAllStringFunc(lex, func(s string) string {
		// 0: whole
		// 1: work name
		// 2: remainder of citation
		groups := neposfinder.FindStringSubmatch(s)
		newworkname, ok := neposstringmap[groups[1]]
		if !ok {
			global.MSG(fmt.Sprintf("fixnepos() lookup failed for %s", groups[1]))
			return s
		}
		return fmt.Sprintf(REPL, newworkname, groups[2])
	})

	return lex
}

func fixpropertius(lex string) string {
	// 	<bibl n="Perseus:abo:phi,1224,001:1:8:29" default="NO"><author>Prop.</author> 1, 8, 29</bibl>
	//	but propertius is lt0620
	return strings.ReplaceAll(lex, "Perseus:abo:phi,1224,001", "Perseus:abo:phi,0620,001")
}

func fixsallust(lex string) string {
	const (
		TMPL = `"Perseus:abo:phi,0631,%s:%s"`
	)
	// sallust is only `0631,001` + `workname`: you need to fix the worknumbers
	// <bibl n="Perseus:abo:phi,0631,001:C. 44:4" default="NO" valid="yes"><author>Sall.</author> C. 44, 4</bibl>
	// <bibl n="Perseus:abo:phi,0631,001:J. 57 sq" default="NO" valid="yes"><author>Sall.</author> J. 57 sq.</bibl>
	// 0: whole
	// 1: work #
	// 2: work name
	// 3: rest of citation
	sallustfinder := regexp.MustCompile(`"Perseus:abo:phi,0631,(001):([^\s]*) ([^"]*)"`)
	lex = sallustfinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := sallustfinder.FindStringSubmatch(s)
		newwknum, ok := sallustringmap[groups[2]]
		if !ok {
			// global.MSG(fmt.Sprintf("fixsallust() lookup failed for %s", groups[2]))
			// fmt.Println(groups[1], groups[2], groups[3])

			// MSG gets really spammy: a real issue is the number of Cicero passages that are in here under Sallust
			// but most of the errors look like "58:12" and so are theoretically valid
			return s
		}
		groups[3] = citationcleaner.ReplaceAllString(groups[3], "")
		return fmt.Sprintf(TMPL, newwknum, groups[3])
	})
	return lex
}

func fixseneca(lex string) string {
	// seneca is badly broken: in the database 1014 is Seneca the Elder; 1017 is Seneca the Younger
	// but the dictionary only cites 1014 (and almost always means to cite 1017)

	// next, the citations are always to work `001` + workname:
	// "Perseus:abo:phi,1014,001:Ira. 3:18:3"
	// "Perseus:abo:phi,1014,001:Q. N. 5:16:5"
	const (
		TMPL = `"Perseus:abo:phi,%s,%s%s"`
	)
	senecafinder := regexp.MustCompile(`"Perseus:abo:phi,1014,001:([^\d]*) ([^"]*)"`)
	lex = senecafinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := senecafinder.FindStringSubmatch(s)
		//for i, group := range groups {
		//	fmt.Println(i, group)
		//}
		authorandwork, ok := senecamap[groups[1]]
		if !ok {
			global.MSG(fmt.Sprintf("fixseneca() lookup failed for %s", groups[1]))
			return s
		}

		author := authorandwork[0]
		work := authorandwork[1]
		return fmt.Sprintf(TMPL, author, work, groups[2])
	})
	return lex
}

func fixsuetonius(lex string) string {
	// "Perseus:abo:phi,1348,001:life=vesp.:16"
	// --> "Perseus:abo:phi,1348,001:Ves:16"
	const (
		TMPL = `"Perseus:abo:phi,1348,001:%s:%s"`
	)
	suetoniusfinder := regexp.MustCompile(`"Perseus:abo:phi,1348,001:life=([^:]*):([^"]*)"`)
	lex = suetoniusfinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := suetoniusfinder.FindStringSubmatch(s)
		life, ok := suetostringmap[groups[1]]
		if !ok {
			global.MSG(fmt.Sprintf("fixsuetonius() lookup failed for %s", groups[1]))
			fmt.Println(groups[1], groups[2])
			return s
		}
		groups[2] = citationcleaner.ReplaceAllString(groups[2], "")
		return fmt.Sprintf(TMPL, life, groups[2])
	})
	return lex
}

func fixvarro(lex string) string {
	// 	the DLL citations are wrong:
	//		"Perseus:abo:phi,0684,001:L. L. 5:section=59"
	//
	//	RR too:
	//		"Perseus:abo:phi,0684,001:R. R. 1:32:2"
	//
	//	lingering issue: 'ibidem'
	//		"Perseus:abo:phi,0684,001:ib. 3:15:2"
	//
	//	lingering issue: DNE
	//		"Perseus:abo:phi,0684,001:Sent. Mor. p. 28"
	//		"Perseus:abo:phi,0684,001:Fragm. p. 241"
	//
	//	lingering issue: guaranteed mishit
	//		"Perseus:abo:phi,0684,001:Sat. Men. 95:10"
	//
	//	lingering issue: might-be-right-might-not
	//		"Perseus:abo:phi,0684,001:1:22"
	//		[is wrong since 'vernum' is in 002 but not in 001]
	//		vernum,</quote> <bibl n="Perseus:abo:phi,0684,001:33:3"

	// hgdb=> select universalid,title from works where universalid ~* 'lt0684' order by universalid;;
	// universalid |            title
	//-------------+------------------------------
	// lt0684w001  | De Lingua Latina
	// lt0684w002  | Res Rusticae
	// lt0684w003  | Antiquitates Rerum Humanarum
	// lt0684w004  | Antiquitates Rerum Divinarum
	// lt0684w005  | Annales
	// lt0684w006  | De Gente Populi Romani
	// lt0684w007  | De Vita Populi Romani
	// lt0684w008  | Res Urbanae
	// lt0684w009  | Logistorici
	// lt0684w010  | carmina
	// lt0684w011  | Menippeae
	// lt0684w012  | epistulae
	// lt0684w013  | epistulae Latinae
	// lt0684w014  | fragmenta grammatica
	// lt0684w015  | frr. de historia litterarum
	// lt0684w016  | fragmenta varia
	// lt0684w017  | incertae sedis fragmenta

	const (
		TMPL = `"Perseus:abo:phi,0684,%s:%s"`
	)
	varromap := map[string]string{
		"L. L.": "001",
		"L.L.":  "001",
		"L.":    "001",
		"R. R.": "002",
		"R.":    "002",
		"Sat.":  "011",
	}
	varrofinder := regexp.MustCompile(`"Perseus:abo:phi,0684,001:([^\s"]*) ([^"]*)"`)
	lex = varrofinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := varrofinder.FindStringSubmatch(s)
		work, ok := varromap[groups[1]]
		if !ok {
			global.MSG(fmt.Sprintf("varromap() lookup failed for %s", groups[1]))
			// varromap() lookup failed for 5:32:156
			return s
		}
		groups[2] = citationcleaner.ReplaceAllString(groups[2], "")
		return fmt.Sprintf(TMPL, work, groups[2])
	})

	return lex
}

func purgelinklessbibls(lex string) string {
	nolink := regexp.MustCompile(`<bibl [^n][^>]*>(.*)</bibl>`)
	return nolink.ReplaceAllString(lex, "$1")
}
