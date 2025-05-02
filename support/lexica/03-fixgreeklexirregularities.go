//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lexica

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"regexp"
)

var (
	eurconverter = map[string]string{
		// the LSJ Eur is 001-018 + 040, 041, 050, 052 (and the last are `correct`)
		/// <bibl n="Perseus:abo:tlg,0006,040:1152"><author>E.</author> <title>Hec.</title>
		// <bibl n="Perseus:abo:tlg,0006,041:111"><author>E.</author> <title>Supp.</title>
		// <bibl n="Perseus:abo:tlg,0006,050:853"><author>Id.</author> <title>Ba.</title>
		// <bibl n="Perseus:abo:tlg,0006,052:97"><author>E.</author> <title>Rh.</title>
		"001": "034", // Cyc.
		"002": "035", // Alc.
		"003": "036", // Med.
		"004": "037", // Heracl.
		"005": "038", // Hipp.
		"006": "039", // Andr.
		"007": "040", // Hec.
		"008": "041", // Supp.
		"009": "043", // HF
		"010": "046", // Ion
		"011": "044", // Tr.
		"012": "042", // El.
		"013": "045", // IT
		"014": "047", // Hel.
		"015": "048", // Ph.
		"016": "049", // Or.
		"017": "050", // Ba.
		"018": "051", // IA
		"038": "038", // Hipp.
		"040": "040", // Hec.
		"041": "041", // Supp.
		"044": "044", // Tr.
		"047": "047", // Hel.
		"050": "050", // Ba.
		"052": "052", // Rh
	}
	eurfinder = regexp.MustCompile(`<bibl n="Perseus:abo:tlg,0006,(\d\d\d):(\d*)">`)
)

func FixSpecificLSJCitations(lex string) string {
	return fixeuripides(lex)
}

// hgdb=> select universalid,title from works where universalid ~* 'gr0006' order by universalid;
// universalid |                title
//-------------+--------------------------------------
// gr0006w020  | Fragmenta
// gr0006w021  | Fragmenta papyracea
// gr0006w022  | Epinicium in Alcibiadem (fragmenta)
// gr0006w023  | Fragmenta Phaethontis
// gr0006w024  | Fragmenta Antiopes
// gr0006w025  | Fragmenta Alexandri
// gr0006w026  | Fragmenta Hypsipyles
// gr0006w027  | Fragmenta Phrixei (P. Oxy. 34.2685)
// gr0006w028  | Fragmenta fabulae incertae
// gr0006w029  | Fragmenta
// gr0006w030  | Fragmenta Oenei
// gr0006w031  | Epigrammata
// gr0006w032  | Fragmenta Phaethontis incertae sedis
// gr0006w033  | Fragmenta
// gr0006w034  | Cyclops
// gr0006w035  | Alcestis
// gr0006w036  | Medea
// gr0006w037  | Heraclidae
// gr0006w038  | Hippolytus
// gr0006w039  | Andromacha
// gr0006w040  | Hecuba
// gr0006w041  | Supplices
// gr0006w042  | Electra
// gr0006w043  | Hercules
// gr0006w044  | Troiades
// gr0006w045  | Iphigenia Taurica
// gr0006w046  | Ion
// gr0006w047  | Helena
// gr0006w048  | Phoenisae
// gr0006w049  | Orestes
// gr0006w050  | Bacchae
// gr0006w051  | Iphigenia Aulidensis
// gr0006w052  | Rhesus

func fixeuripides(lex string) string {
	// <bibl n="Perseus:abo:tlg,0006,018:601"><author>E.</author> <title>IA</title> 601</bibl>
	// <bibl n="Perseus:abo:tlg,0006,015:309"><author>E.</author> <title>Ph.</title> 309</bibl>

	const (
		REPL = `<bibl n="Perseus:abo:tlg,0006,%s:%s">`
	)

	lex = eurfinder.ReplaceAllStringFunc(lex, func(s string) string {
		// 0: whole
		// 1: work#
		// 2: line#
		groups := eurfinder.FindStringSubmatch(s)
		newworknum, ok := eurconverter[groups[1]]
		if !ok {
			global.MSG(fmt.Sprintf("fixeuripides() lookup failed for %s", groups[1]))
			return s
		}
		return fmt.Sprintf(REPL, newworknum, groups[2])
	})

	return lex
}
