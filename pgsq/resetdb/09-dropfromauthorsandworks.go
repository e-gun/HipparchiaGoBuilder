//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import (
	"fmt"

	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

// before (re)building, purge the old authors and works tables
// this will not work properly for PHI...

// hgdb=# select universalid from authors where universalid ~* '^lat11';
// universalid
// -------------
// lat1100
// lat1103
// (2 rows)
//
// hgdb=# delete from authors where universalid ~* '^lat11';
// DELETE 2
// hgdb=# select universalid from authors where universalid ~* '^lat11';
// universalid
// -------------
//(0 rows)

var (
	deleterfromauthors = `DELETE FROM authors WHERE universalid %s '%s'`
	deleterfromworks   = `DELETE FROM works WHERE universalid %s '^%s'`
)

func ResetAuthorsAndWorks(pfx string) {
	queries := []string{
		purgeclassofauthors(pfx),
		purgeclassorworks(pfx),
	}
	dbc.DBCCommandSequence(queries)
}

func ResetOneAuthor(oneid string) {
	fmt.Println("ResetOneAuthor() " + oneid)
	queries := []string{
		purgeoneauthor(oneid),
		purgeallworksofoneauthor(oneid),
	}

	dbc.DBCCommandSequence(queries)
}

func purgeoneauthor(oneid string) string {
	syntax := `=`
	return fmt.Sprintf(deleterfromauthors, syntax, oneid)
}

func purgeclassofauthors(pfx string) string {
	syntax := `~*`
	pfx = `^` + pfx
	return fmt.Sprintf(deleterfromauthors, syntax, pfx)
}

func purgeallworksofoneauthor(oneid string) string {
	syntax := `~*`
	return fmt.Sprintf(deleterfromworks, syntax, oneid)
}

func purgeclassorworks(pfx string) string {
	syntax := `~*`
	return fmt.Sprintf(deleterfromworks, syntax, pfx)
}
