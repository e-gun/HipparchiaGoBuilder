//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func FixHyphens(lines []structs.DbWorkline) []structs.DbWorkline {
	// note that BuildStrippedAndAccentedlines() can be a problem if you are not careful
	// 	sample in:
	//		'ἀγῶϲι πρόνοιαν· ὃϲ καὶ τότε περαιούμενοϲ ναυϲὶν ἐϲ Ἰτα-', 'ἀγῶϲι πρόνοιαν ὃϲ καὶ τότε περαιούμενοϲ ναυϲὶν ἐϲ ἰτα-'
	//	sample out:
	//		'ἀγῶϲι πρόνοιαν· ὃϲ καὶ τότε περαιούμενοϲ ναυϲὶν ἐϲ Ἰτα-', 'ἀγῶϲι πρόνοιαν ὃϲ καὶ τότε περαιούμενοϲ ναυϲὶν ἐϲ ἰταλίαν'
	//		+ add the hyphenated word to the hyphens column
	//		+ modify the next line...

	previous := structs.DbWorkline{}

	for i, thisline := range lines {
		if previous.HasHyphen() {
			previous, thisline = consolidatecontiguous(previous, thisline)
		}

		lines[i] = thisline
		if i > 0 {
			lines[i-1] = previous
		}
		previous = thisline
	}
	for i := range lines {
		// strippedliner can not do the following; so we do it here
		lines[i].Accented = strings.ReplaceAll(lines[i].Accented, "-", "")
		lines[i].Stripped = strings.ReplaceAll(lines[i].Stripped, "-", "")
	}

	return lines
}

func consolidatecontiguous(previous structs.DbWorkline, this structs.DbWorkline) (structs.DbWorkline, structs.DbWorkline) {
	// previous has a hyphen, so you have something like:
	// p: nos ambae faciunt in hoc tempore, summa gratia et elo-
	// t: quentia; quarum alteram, C. Aquili, uereor, alteram metuo.

	// you need
	// p: nos ambae faciunt in hoc tempore, summa gratia et eloquentia
	// t: quarum alteram, C. Aquili, uereor, alteram metuo.

	// fix the accented line
	hw := strings.TrimSuffix(previous.LastAccentedWord(), "-") + this.FirstAccentedWord()
	previous.Hyphenated = hw
	newwords := append(previous.AllButLastAccentedWord(), hw)
	previous.Accented = strings.Join(newwords, " ")
	this.Accented = strings.Join(this.AllButFirstAccentedWord(), " ")

	// now the stripped line
	hw = strings.TrimSuffix(previous.LastStrippedWord(), "-") + this.FirstStrippedWord()
	newwords = append(previous.AllButLastStrippedWord(), hw)
	previous.Stripped = strings.Join(newwords, " ")
	this.Stripped = strings.Join(this.AllButFirstStrippedWord(), " ")

	return previous, this
}
