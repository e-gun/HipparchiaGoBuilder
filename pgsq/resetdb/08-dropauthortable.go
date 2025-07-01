//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import (
	"context"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"strings"
)

// before (re)building, purge the old authors and works tables
// this will not work properly for PHI...

func DropAuthorMultipleTables(pfx string) {
	const (
		DROPTABLES = `DROP TABLE %s` // comma separated list...
		FINDTABLES = `SELECT universalid FROM authors WHERE universalid ~* '^%s'`
	)

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	q1 := fmt.Sprintf(FINDTABLES, pfx)
	foundrows, e := dbconn.Query(context.Background(), q1)
	if e != nil {
		panic(e)
	}

	var tablestodrop []string
	defer foundrows.Close()
	for foundrows.Next() {
		var thetable string
		err := foundrows.Scan(&thetable)
		if err != nil {
			panic(err)
		}
		tablestodrop = append(tablestodrop, thetable)
	}

	if len(tablestodrop) == 0 {
		// you are already empty...
		return
	}

	q2 := fmt.Sprintf(DROPTABLES, strings.Join(tablestodrop, ","))
	_, err := dbconn.Exec(context.Background(), q2)
	if err != nil {
		fmt.Println(err)
	}
}

func DropOneAuthorTable(oneid string) {
	const (
		DROPTABLE = `DROP TABLE %s` // comma separated list...
	)

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	q2 := fmt.Sprintf(DROPTABLE, oneid)
	_, err := dbconn.Exec(context.Background(), q2)
	if err != nil {
		fmt.Println(err)
	}
}
