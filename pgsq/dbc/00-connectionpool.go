//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dbc

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	SQLPool *pgxpool.Pool // filled by `main.go` and a call to FillDBConnectionPool
)

type PostgresLogin struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBName string
}

// FillDBConnectionPool - build the pgxpool that the whole program will Acquire() from
func FillDBConnectionPool(pl PostgresLogin) *pgxpool.Pool {
	// nb: macos users can send ANYTHING as a password for the dbc user: admin access already (on their primary account...)

	const (
		UTPL    = "postgres://%s:%s@%s:%d/%s?pool_min_conns=%d&pool_max_conns=%d"
		FAIL1   = "FillDBConnectionPool() Configuration error. Could not execute ParseConfig(url) via '%s'"
		FAIL2   = "FillDBConnectionPool() could not connect to PostgreSQL"
		ERRRUN  = `dial error`
		FAILRUN = `'%s': the PostgreSQL server cannot be found; check that it is running and serving on port %d`
		ERRSRV  = `server error`
		FAILSRV = `'%s': there is configuration problem; see the following response from PostgreSQL:`
	)

	mn := runtime.NumCPU() / 2
	mx := runtime.NumCPU()

	url := fmt.Sprintf(UTPL, pl.User, pl.Pass, pl.Host, pl.Port, pl.DBName, mn, mx)

	config, e := pgxpool.ParseConfig(url)
	if e != nil {
		fmt.Println(fmt.Sprintf(FAIL1, url))
		os.Exit(0)
	}

	thepool, e := pgxpool.NewWithConfig(context.Background(), config)
	if e != nil {
		fmt.Println(fmt.Sprintf(FAIL2))
		if strings.Contains(e.Error(), ERRRUN) {
			fmt.Println(fmt.Sprintf(FAILRUN, ERRRUN, pl.Port))
		}
		if strings.Contains(e.Error(), ERRSRV) {
			fmt.Println(fmt.Sprintf(FAILSRV, ERRSRV))
			parts := strings.Split(e.Error(), ERRSRV)
			fmt.Println(parts[1])
		}
		os.Exit(0)
	}
	return thepool
}
