//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

const (
	LATLEMTESTDATA = `Aaron	19126	Aaron (masc nom sg)	Aaronem (masc acc sg)	Aaroni (masc dat sg)	Aaronis (masc gen sg)
Aba	90033	Abas (masc acc pl)	Abin (masc abl pl) (masc dat pl)	Abis (masc abl pl) (masc dat pl)	Abisue (masc abl pl) (masc dat pl)
Abanteus	106721	Abanteis (neut abl pl) (fem abl pl) (masc abl pl) (neut dat pl) (fem dat pl) (masc dat pl)
Abantiades	106721	Abantiadas (masc acc pl)	Abantiades (masc nom sg)
Abantius	106721	Abanti (masc/neut gen sg)	Abantia (neut nom/voc/acc pl) (fem abl sg) (fem nom/voc sg)	Abantias (fem acc pl)`
)

func ParseLatinLemmata(entries []string) []structs.HeadwordAndForms {
	return ParseLemmata("latin", entries)
}
