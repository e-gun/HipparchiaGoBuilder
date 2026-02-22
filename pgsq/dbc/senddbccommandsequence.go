//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dbc

import (
	"context"
	"fmt"
)

// DBCCommandSequence - execute a chain of queries via a psql connection
func DBCCommandSequence(queries []string) {
	dbconn := GetDBConnection()
	defer dbconn.Release()

	for _, query := range queries {
		_, err := dbconn.Exec(context.Background(), query)
		if err != nil {
			fmt.Println(err)
		}
	}
}
