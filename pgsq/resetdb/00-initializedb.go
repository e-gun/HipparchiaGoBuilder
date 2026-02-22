//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import (
	"fmt"
	"runtime"

	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
)

// DropDatabase - simply drop the database
func DropDatabase(adminpass string) {
	const (
		C1    = `DROP DATABASE IF EXISTS %s;`
		DONE1 = "Dropped the database "
	)

	queries := []string{
		fmt.Sprintf(C1, pgsq.DEFAULTPSQLDB),
	}

	sendpsqladmincommndstobinary(adminpass, queries)

	fmt.Println(DONE1)
}

// InitializeHDB - insert the hipparchiaDB table and its user into postgres
func InitializeHDB(adminpass string, hdbpw string) {
	fmt.Println("InitializeHDB()")
	const (
		C0   = `DROP DATABASE IF EXISTS %s;`
		C1   = `DROP ROLE IF EXISTS %s;`
		C2   = `CREATE ROLE %s LOGIN ENCRYPTED PASSWORD '%s' NOSUPERUSER INHERIT CREATEDB NOCREATEROLE NOREPLICATION;`
		C3   = `CREATE DATABASE "%s" WITH OWNER = %s ENCODING = 'UTF8';`
		C4   = `GRANT CREATE ON DATABASE %s TO %s;`
		C5   = `CREATE EXTENSION IF NOT EXISTS pg_trgm;` // this one is not executing properly?
		DONE = "Initialized the database framework"
	)

	queries := []string{
		fmt.Sprintf(C0, pgsq.DEFAULTPSQLDB),
		fmt.Sprintf(C1, pgsq.DEFAULTPSQLUSER),
		fmt.Sprintf(C2, pgsq.DEFAULTPSQLUSER, hdbpw),
		fmt.Sprintf(C3, pgsq.DEFAULTPSQLDB, pgsq.DEFAULTPSQLUSER),
		fmt.Sprintf(C4, pgsq.DEFAULTPSQLDB, pgsq.DEFAULTPSQLUSER),
	}

	sendpsqladmincommndstobinary(adminpass, queries)

	// need to create the extension on the db, not just "in the abstract"
	queries = []string{
		fmt.Sprintf(C5),
	}

	dbc.DBCCommandSequence(queries)
	fmt.Println(DONE)
}

// RequestPostgresAdminPW - ask for the password for the postgres admin user
func RequestPostgresAdminPW() string {
	const (
		PWD2 = "I also need the database password for the postgres administrator -> "
	)
	var pgpw string
	if runtime.GOOS != "darwin" {
		// macos users have admin access already (on their primary account...) and do not need a pg admin password
		fmt.Printf(fmt.Sprintf(PWD2))
		_, ee := fmt.Scan(&pgpw)
		if ee != nil {
			panic(ee)
		}
	}
	return pgpw
}
