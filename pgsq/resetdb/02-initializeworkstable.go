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

const (
	CW = `
DROP TABLE IF EXISTS public.works;
CREATE TABLE public.works (
    universalid character(11),
    title character varying(512),
    language character varying(1),
    publication_info text,
    levellabels_00 character varying(64),
    levellabels_01 character varying(64),
    levellabels_02 character varying(64),
    levellabels_03 character varying(64),
    levellabels_04 character varying(64),
    levellabels_05 character varying(64),
    workgenre character varying(256),
    transmission character varying(64),
    worktype character varying(64),
    provenance character varying(196),
    recorded_date character varying(196),
    converted_date integer,
    wordcount integer,
    firstline integer,
    lastline integer,
    authentic boolean,
    UNIQUE(universalid)
);

ALTER TABLE public.works OWNER TO %s;`
)

func InitializeWorksTable() {
	const (
		DONE = "Initialized the works table"
	)
	queries := []string{
		fmt.Sprintf(CW, pgsq.DEFAULTPSQLUSER),
	}

	// sendpsqlusercommandstobinary(userpass, queries)
	dbc.DBCCommandSequence(queries)
	fmt.Println(DONE)
}
