package lexica

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
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
		"ad Marc.":          {"1017", "012:6:"},
		"ad Polyb.":         {"1017", "012:11:"},
		"Agm.":              {"1017", "007:"},
		"Agam.":             {"1017", "007:"},
		"Apoc.":             {"1017", "011:"},
		"Apocol.":           {"1017", "011:"},
		"Apocol.p.":         {"1017", "011:"},
		"Ben.":              {"1017", "013:"},
		"Benef.":            {"1017", "013:"},
		"Brev. Vit.":        {"1017", "012:10:"},
		"Clem.":             {"1017", "014:"},
		"de Clem.":          {"1017", "014:"},
		"de Prov.":          {"1017", "012:1:"},
		"Ep":                {"1017", "015:"},
		"Ep.":               {"1017", "015:"},
		"Consol. ad Marc.":  {"1017", "012:6:"},
		"Cons. ad Marc.":    {"1017", "012:6:"},
		"Cons. ad Helv.":    {"1017", "012:12:"},
		"Cons. ad Polyb.":   {"1017", "012:11:"},
		"Cons. Helv.":       {"1017", "012:12:"},
		"Cons. Marc.":       {"1017", "012:6:"},
		"Cons. Polyb.":      {"1017", "012:11:"},
		"Const.":            {"1017", "012:2:"},
		"Const. Sap.":       {"1017", "012:2:"},
		"Contr":             {"1014", "001:"}, // the father...
		"Contr.":            {"1014", "001:"}, // the father...
		"Controv.":          {"1014", "001:"}, // the father...
		"Exc. Contr.":       {"1014", "002:"}, // the father...
		"Exc. Controv.":     {"1014", "002:"}, // the father...
		"Excerpt. Contr.":   {"1014", "002:"}, // the father...
		"Excerpt. Controv.": {"1014", "002:"}, // the father...
		"Helv.":             {"1017", "012:12:"},
		"Herc. Fur.":        {"1017", "001:"},
		"Herc Oet.":         {"1017", "009:"},
		"Herc. Oet.":        {"1017", "009:"},
		"Hipp.":             {"1017", "005:"},
		"Hippol.":           {"1017", "005:"},
		"Ira,":              {"1017", "012:3"}, // note that we just sent I, II, and III to the same place...
		"Lud. Mort. Claud.": {"1017", "011:"},
		"Med.":              {"1017", "004:"},
		"Mort. Claud.":      {"1017", "011:"},
		"N. Q.":             {"1017", "016:"},
		"Oct.":              {"1017", "010:"},
		"Octav.":            {"1017", "010:"},
		"Ot. Sap.":          {"1017", "012:8:"},
		"Oed.":              {"1017", "006:"},
		"Oedip":             {"1017", "006:"},
		"Oedip.":            {"1017", "006:"},
		"Oet.":              {"1017", "009:"},
		// "Orest.": {"1017", "006:"},  // ?!
		"Phaedr.":     {"1017", "005:"},
		"Phoen.":      {"1017", "003:"},
		"Polyb.":      {"1017", "012:11:"},
		"Prov":        {"1017", "012:1:"},
		"Prov.":       {"1017", "012:1:"},
		"Q. N.":       {"1017", "016:"},
		"Suas.":       {"1014", "003:"}, // the father...
		"Thyest.":     {"1017", "008:"},
		"Tranq.":      {"1017", "012:9:"},
		"Tranq. An.":  {"1017", "012:9:"},
		"Tranq. Vit.": {"1017", "012:9:"},
		"Troad.":      {"1017", "002"},
		"Vit. B.":     {"1017", "012:7:"},
		"Vit. Beat.":  {"1017", "012:7:"},
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
	ibidemfinder    = regexp.MustCompile(`"Perseus:abo:phi,(\d\d\d\d,\d\d\d):(ib\. )([^"]*)"`)
	neposfinder     = regexp.MustCompile(`"Perseus:abo:phi,0588,001:([^\s]*\.) ([^"]*)"`)
	sallustfinder   = regexp.MustCompile(`"Perseus:abo:phi,0631,(001):([^\s]*) ([^"]*)"`)
	citationcleaner = regexp.MustCompile(`[^0-9:]`)
	suetoniusfinder = regexp.MustCompile(`"Perseus:abo:phi,1348,001:life=([^:]*):([^"]*)"`)
)

func FixSpecificLatinCitations(lex string) string {
	lex = fixibidem(lex)
	lex = fixnepos(lex)
	lex = fixsallust(lex)
	lex = fixsuetonius(lex)
	return lex
}

func fixibidem(lex string) string {
	const (
		TMPL = `"Perseus:abo:phi,%s:%s"`
	)
	// "Perseus:abo:phi,1014,001:ib. 76:2" --> "Perseus:abo:phi,1014,001:76:2"
	return ibidemfinder.ReplaceAllStringFunc(lex, func(s string) string {
		groups := ibidemfinder.FindStringSubmatch(s)
		return fmt.Sprintf(TMPL, groups[1], groups[3])
	})
}

func fixnepos(lex string) string {
	// <bibl n="Perseus:abo:phi,0588,001:Alcib. 11:4" default="NO" valid="yes"><author>Nep.</author> Alcib. 11, 4</bibl></cit>:
	// --> <bibl n="Perseus:abo:phi,0588,001:Alc:11:4" default="NO" valid="yes"><author>Nep.</author> Alcib. 11, 4</bibl>
	const (
		REPL = `"Perseus:abo:phi,0588,001:%s:%s"`
	)

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

func fixsuetonius(lex string) string {
	// "Perseus:abo:phi,1348,001:life=vesp.:16"
	// --> "Perseus:abo:phi,1348,001:Ves:16"
	const (
		TMPL = `"Perseus:abo:phi,1348,001:%s:%s"`
	)
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
