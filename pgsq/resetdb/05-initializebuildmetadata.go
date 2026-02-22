//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

const (
	BMD = `
DROP TABLE IF EXISTS public.buildmetadata;
CREATE TABLE public.buildmetadata (
    category   character varying(64),
    builderver character varying(64),
    gitcommit  character varying(64),
    date       character varying(64),
    notes      character varying(128),
    UNIQUE (category)
);`
)

func InitializeBuildMetadataTable() {
	const (
		DONE = "Initialized the buildmetadata table"
	)

	queries := []string{
		BMD,
	}
	fmt.Println(DONE)
	dbc.DBCCommandSequence(queries)
}
