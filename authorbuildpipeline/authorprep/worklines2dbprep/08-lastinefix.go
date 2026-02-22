//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func LastLinefix(lines []structs.DbWorkline) []structs.DbWorkline {
	// this might no longer be true...

	// findhyphens() vs addcdlabels()

	// '-█ⓕⓔ' lines are not read as ending with '-' by the time they get to findhyphens().
	//
	// it's an ugly issue

	// the following lines of isaeus refuse to match any conditional to test for a hyphen you throw at them and so
	// will not enter into the 'if...' clause; but if you cut and paste the text '-' == 'True'

	//   "ἔφη τήν τε ἡλικίαν ὑφορᾶϲθαι τὴν ἑαυτοῦ καὶ τὴν ἀπαι-"
	//   █⑧⓪ E)/FH TH/N TE H(LIKI/AN U(FORA=SQAI TH\N E(AUTOU= KAI\ TH\N A)PAI-█ⓕⓔ
	//   "Εἶτα αὐτὸϲ μὲν εἰ ἦν ἄπαιϲ, ἐποιήϲατ’ ἄν· τὸν δὲ Με-"
	//   █⑧⓪ *EI)=TA AU)TO\S ME\N EI) H)=N A)/PAIS, E)POIH/SAT' A)/N: TO\N DE\ *ME-█ⓕⓔ
	//
	// >>> x = 'ἀπαι-'
	// >>> if '-' in x: print('yes')
	// ...
	// yes
	// >>>

	//  the fix is to remove the trailing space in regexsubs addcdlabels(). but then that kills your ability to get
	//  the last line of a work into the dbc
	//  replace = '\n<hmu_end_of_cd_block_re-initialize_key_variables />'
	//  vs replace = '\n<hmu_end_of_cd_block_re-initialize_key_variables /> '

	// so the simple solution to a complex problem is to slap a single whitespace at the end of the file
	lines[len(lines)-1].MarkedUp += " "
	return lines
}
