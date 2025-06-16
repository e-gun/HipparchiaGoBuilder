//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

// longest recorded_date is gr2023: `A.D. 3–4 vid Scholia: Cf. &3SCHOLIA IN IAMBLICHUM PHILOSOPHUM& (5027)`
const (
	CA = `
DROP TABLE IF EXISTS public.authors;
CREATE TABLE public.authors (
    universalid character(7),
    language character varying(11),
    idxname character varying(256),
    akaname character varying(256),
    shortname character varying(256),
    cleanname character varying(256),
    genres character varying(256),
    recorded_date character varying(256),
    converted_date integer,
    location character varying(256),
    UNIQUE(universalid)
);


ALTER TABLE public.authors OWNER TO %s;`
)

func InitializeAuthorTable() {
	const (
		DONE = "Initialized the authors table"
	)
	queries := []string{
		fmt.Sprintf(CA, pgsq.DEFAULTPSQLUSER),
	}

	dbc.DBCCommandSequence(queries)

	fmt.Println(DONE)
}
