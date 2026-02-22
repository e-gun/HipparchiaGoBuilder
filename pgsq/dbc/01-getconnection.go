//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dbc

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetDBConnection - Acquire() a connection from the main pgxpool
func GetDBConnection() *pgxpool.Conn {
	dbc, e := SQLPool.Acquire(context.Background())
	if e != nil {
		os.Exit(0)
	}
	return dbc
}
