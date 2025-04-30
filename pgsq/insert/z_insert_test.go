//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"context"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"testing"
)

func ReadyDBConnection() {
	pl := dbc.PostgresLogin{
		Host:   pgsq.DEFAULTPSQLHOST,
		Port:   pgsq.DEFAULTPSQLPORT,
		User:   pgsq.DEFAULTPSQLUSER,
		Pass:   pgsq.DEFAULTPSQLPASS,
		DBName: pgsq.DEFAULTPSQLDB,
	}

	dbc.SQLPool = dbc.FillDBConnectionPool(pl)
}

func TestInsertHeadwordWordcountsIntoTable(t *testing.T) {
	const (
		EN1 = `TestInsertHeadwordWordcountsIntoTable`
		Q1  = `DELETE FROM headword_wordcounts WHERE entry_name = $1;`
	)

	ReadyDBConnection()

	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	// reset the old data
	_, err := dbconn.Exec(context.Background(), Q1, EN1)
	if err != nil {
		fmt.Println(err)
	}

	wct := structs.DbHeadwordCounts{Word: EN1}
	towrite := map[string]structs.DbHeadwordCounts{EN1: wct}
	InsertHeadwordWordcountsIntoTable(towrite)
}
