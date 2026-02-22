//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

const (
	WC0 = `DROP TABLE IF EXISTS public.headword_wordcounts;`
	WC1 = `
CREATE TABLE public.headword_wordcounts (
    entry_name character varying(64),
    total_count integer DEFAULT 0,
    gr_count integer DEFAULT 0,
    lt_count integer DEFAULT 0,
    dp_count integer DEFAULT 0,
    in_count integer DEFAULT 0,
    ch_count integer DEFAULT 0,
    frequency_classification character varying(64),
    early_occurrences integer DEFAULT 0,
    middle_occurrences integer DEFAULT 0,
    late_occurrences integer DEFAULT 0,
    acta integer DEFAULT 0,
    agric integer DEFAULT 0,
    alchem integer DEFAULT 0,
    anthol integer DEFAULT 0,
    apocalyp integer DEFAULT 0,
    apocryph integer DEFAULT 0,
    apol integer DEFAULT 0,
    astrol integer DEFAULT 0,
    astron integer DEFAULT 0,
    biogr integer DEFAULT 0,
    bucol integer DEFAULT 0,
    caten integer DEFAULT 0,
    chronogr integer DEFAULT 0,
    comic integer DEFAULT 0,
    comm integer DEFAULT 0,
    concil integer DEFAULT 0,
    coq integer DEFAULT 0,
    dialog integer DEFAULT 0,
    docu integer DEFAULT 0,
    doxogr integer DEFAULT 0,
    eccl integer DEFAULT 0,
    eleg integer DEFAULT 0,
    encom integer DEFAULT 0,
    epic integer DEFAULT 0,
    epigr integer DEFAULT 0,
    epist integer DEFAULT 0,
    evangel integer DEFAULT 0,
    exeget integer DEFAULT 0,
    fab integer DEFAULT 0,
    geogr integer DEFAULT 0,
    gnom integer DEFAULT 0,
    gramm integer DEFAULT 0,
    hagiogr integer DEFAULT 0,
    hexametr integer DEFAULT 0,
    hist integer DEFAULT 0,
    homilet integer DEFAULT 0,
    hymn integer DEFAULT 0,
    hypoth integer DEFAULT 0,
    iamb integer DEFAULT 0,
    ignotum integer DEFAULT 0,
    inscr integer DEFAULT 0,
    invectiv integer DEFAULT 0,
    jurisprud integer DEFAULT 0,
    lexicogr integer DEFAULT 0,
    liturg integer DEFAULT 0,
    lyr integer DEFAULT 0,
    magica integer DEFAULT 0,
    math integer DEFAULT 0,
    mech integer DEFAULT 0,
    med integer DEFAULT 0,
    metrolog integer DEFAULT 0,
    mim integer DEFAULT 0,
    mus integer DEFAULT 0,
    myth integer DEFAULT 0,
    narrfict integer DEFAULT 0,
    nathist integer DEFAULT 0,
    onir integer DEFAULT 0,
    orac integer DEFAULT 0,
    orat integer DEFAULT 0,
    papyrus integer DEFAULT 0,
    paradox integer DEFAULT 0,
    parod integer DEFAULT 0,
    paroem integer DEFAULT 0,
    perieg integer DEFAULT 0,
    phil integer DEFAULT 0,
    physiognom integer DEFAULT 0,
    poem integer DEFAULT 0,
    polyhist integer DEFAULT 0,
    prophet integer DEFAULT 0,
    pseudepigr integer DEFAULT 0,
    rhet integer DEFAULT 0,
    satura integer DEFAULT 0,
    satyr integer DEFAULT 0,
    schol integer DEFAULT 0,
    tact integer DEFAULT 0,
    test integer DEFAULT 0,
    theol integer DEFAULT 0,
    trag integer DEFAULT 0,
    allrhet integer DEFAULT 0,
    allrel integer DEFAULT 0,
    UNIQUE(entry_name)
);

ALTER TABLE public.headword_wordcounts OWNER TO %s;
`

	WC2 = `DROP TABLE IF EXISTS public.unparsed_wordcounts;`
	WC3 = `
CREATE TABLE public.unparsed_wordcounts (
    entry_name character varying(96) NOT NULL,
    total_count integer DEFAULT 0,
    gr_count integer DEFAULT 0,
    lt_count integer DEFAULT 0,
    dp_count integer DEFAULT 0,
    in_count integer DEFAULT 0,
    ch_count integer DEFAULT 0,
    UNIQUE(entry_name)
);`
	WC4 = `ALTER TABLE public.unparsed_wordcounts OWNER TO %s;`
)

func InitializeUnparsedWordCountTable() {
	const (
		DONE = "Initialized the 'unparsed_wordcounts' table"
	)

	queries := []string{
		WC2,
		WC3,
		fmt.Sprintf(WC4, pgsq.DEFAULTPSQLUSER),
	}
	fmt.Println(DONE)
	dbc.DBCCommandSequence(queries)
}

func InitializeHeadwordWordCountTable() {
	const (
		DONE = "Initialized the 'headword_wordcounts' table"
	)

	queries := []string{
		fmt.Sprintf(WC0),
		fmt.Sprintf(WC1, pgsq.DEFAULTPSQLUSER),
	}
	fmt.Println(DONE)
	dbc.DBCCommandSequence(queries)
}
