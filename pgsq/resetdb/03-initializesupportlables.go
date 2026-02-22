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
	S1 = `
DROP TABLE IF EXISTS public.%s_dictionary;
CREATE TABLE public.%s_dictionary (
    idval float,
    entry_name character varying(256),
    entry_metr character varying(256),
    idname character varying(16),
    entry_type character varying(16),
    translations text,
    usedby character varying(1500),
    prelim_info text,
    sense_ids text,
    senses jsonb,
    UNIQUE (idval)
);

ALTER TABLE public.%s_dictionary OWNER TO %s;
`
	S3 = `
DROP TABLE IF EXISTS public.%s_lemmata;
CREATE TABLE public.%s_lemmata (
    dictionary_entry character varying(64),
    xref_number integer,
    derivative_forms text[]
);

ALTER TABLE public.%s_lemmata OWNER TO %s;
`
	S4 = `
DROP TABLE IF EXISTS public.%s_morphology;
CREATE TABLE public.%s_morphology (
    observed_form character varying(64),
    xrefs character varying(512),
    prefixrefs character varying(256),
    possible_dictionary_forms jsonb,
    related_headwords character varying(512)
);

ALTER TABLE public.%s_morphology OWNER TO %s;
`
)

func InitializeAllSupportTables(lang string) {
	const (
		DONE = "Initialized the %s support tables\n"
	)
	queries := []string{
		fmt.Sprintf(S1, lang, lang, lang, pgsq.DEFAULTPSQLUSER),
		fmt.Sprintf(S3, lang, lang, lang, pgsq.DEFAULTPSQLUSER),
		fmt.Sprintf(S4, lang, lang, lang, pgsq.DEFAULTPSQLUSER),
	}

	dbc.DBCCommandSequence(queries)

	fmt.Printf(DONE, lang)
}

func InitializeLexicalSupportTables(lang string) {
	const (
		DONE = "Initialized the %s lexical support tables\n"
	)
	queries := []string{
		fmt.Sprintf(S1, lang, lang, lang, pgsq.DEFAULTPSQLUSER),
		// fmt.Sprintf(S2, lang, lang, lang, pgsq.DEFAULTPSQLUSER),
	}

	dbc.DBCCommandSequence(queries)

	fmt.Printf(DONE, lang)
}

func InitializeGrammarSupportTables(lang string) {
	const (
		DONE = "Initialized the %s grammar support tables\n"
	)
	queries := []string{
		fmt.Sprintf(S3, lang, lang, lang, pgsq.DEFAULTPSQLUSER),
		fmt.Sprintf(S4, lang, lang, lang, pgsq.DEFAULTPSQLUSER),
	}

	// sendpsqlusercommandstobinary(userpass, queries)
	dbc.DBCCommandSequence(queries)

	fmt.Printf(DONE, lang)
}

func InitializeAllGreekSupportTables() {
	InitializeAllSupportTables("greek")
}

func InitializeGreekLexicalSupportTables() {
	InitializeLexicalSupportTables("greek")
}

func InitializeGreekGrammarSupportTables() {
	InitializeGrammarSupportTables("greek")
}

func InitializeAllLatinSupportTables() {
	InitializeAllSupportTables("latin")
}

func InitializeLatinLexicalSupportTables() {
	InitializeLexicalSupportTables("latin")
}

func InitializeLatinGrammarSupportTables() {
	InitializeGrammarSupportTables("latin")
}
