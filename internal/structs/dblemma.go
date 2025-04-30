//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

type DbLemma struct {
	// dictionary_entry | xref_number |    derivative_forms
	Entry string
	Xref  int
	Deriv []string
}

func (dbl DbLemma) EntryRune() []rune {
	return []rune(dbl.Entry)
}
