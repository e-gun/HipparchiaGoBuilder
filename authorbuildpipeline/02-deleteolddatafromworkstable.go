//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

// DeleteOldDataFromWorksTable - if you build LAT, then first drop all 'lt' authors from the works table
func DeleteOldDataFromWorksTable(prefix string) {
	const (
		Q = `DELETE FROM public.works WHERE universalid ~* '^%s'`
	)
	queries := []string{
		fmt.Sprintf(Q, getcorpusabbrev(prefix)),
	}

	dbc.DBCCommandSequence(queries)
}
