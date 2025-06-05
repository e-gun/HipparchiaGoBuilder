//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import "github.com/e-gun/HipparchiaGoBuilder/internal/structs"

const (
	LATANALYSISTESTDATA = `discutitur	{22924660 9 discutio	 	pres ind pass 3rd sg}
discutiunt	{22924660 9 discutio	 	pres ind act 3rd pl}
discutiuntque	{22924660 9 discutiunt,discutio	 	pres ind act 3rd pl}
discutiuntur	{22924660 9 discutio	 	pres ind pass 3rd pl}
disdidi	{22932594 0 disdidi_,dis-do	 	perf ind act 1st sg}
disemis	{22932594 9 dise_mi_s,disemus	 	neut abl pl}{22932594 9 dise_mi_s,disemus	 	fem abl pl}{22932594 9 dise_mi_s,disemus	 	masc abl pl}{22932594 9 dise_mi_s,disemus	 	neut dat pl}{22932594 9 dise_mi_s,disemus	 	fem dat pl}{22932594 9 dise_mi_s,disemus	 	masc dat pl}`

	LTD2 = `aduersae	{1699046 9 adversae,adverro	 	perf part pass fem nom/voc pl}{1699046 9 adversae,adverro	 	perf part pass fem dat sg}{1699046 9 adversae,adverro	 	perf part pass fem gen sg}{1699420 0 adversae,adversa	 	fem nom/voc pl}{1699420 0 adversae,adversa	 	fem dat sg}{1699420 0 adversae,adversa	 	fem gen sg}{1747058 9 adversae,adverto	 	perf part pass fem nom/voc pl}{1747058 9 adversae,adverto	 	perf part pass fem dat sg}{1747058 9 adversae,adverto	 	perf part pass fem gen sg}`
)

func ParseLatinAnalyses(entries []string) []structs.GramAnalysis {
	// note that Possibilities []MorphPossib inside each latin GramAnalysis will never contain
	// a value in the Transl field; this means that you cannot build a Latin vocab list off of these
	// in HGS
	return ParseAnalyses("latin", entries)
}
