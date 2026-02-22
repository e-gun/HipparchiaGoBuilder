//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

// DeleteOldDataFromAuthorsTable - if you build LAT, then first drop all 'lt' authors from the authors table
func DeleteOldDataFromAuthorsTable(prefix string) {
	const (
		Q = `DELETE FROM public.authors WHERE universalid ~* '^%s'`
	)
	queries := []string{
		fmt.Sprintf(Q, getcorpusabbrev(prefix)),
	}

	dbc.DBCCommandSequence(queries)
}
