//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"strings"
)

func CleanDBWorkNames(w *structs.DbWork) {
	n := w.Title
	n = simplegreekspan(n)
	n = betacode.BetaCodeCleanup(n)
	n = lat.ConvertLatinDiacriticals(n)
	w.Title = strings.ToValidUTF8(n, "")
}
