//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package insert

import (
	"context"
	"fmt"
	"time"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

func InsertBuildMetadata(cat string, notes string) {
	const (
		INS = `INSERT INTO buildmetadata (category, builderver, gitcommit, date, notes)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (category)
			    DO UPDATE SET
			                  category = EXCLUDED.category,
			                  builderver = EXCLUDED.builderver,
			                  gitcommit = EXCLUDED.gitcommit,
			                  date = EXCLUDED.date,
			                  notes = EXCLUDED.notes`
	)
	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	now := time.Now().Format("15:04:05 02 Jan 2006")

	// UPSERT: update if the author already exists
	_, err := dbconn.Exec(context.Background(), INS, cat, global.VERSION, global.GitHash, now, notes)
	if err != nil {
		fmt.Println(err)
	}
}
