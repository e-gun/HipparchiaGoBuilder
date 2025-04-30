//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package resetdb

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// sendpsqladmincommndstobinary - execute a chain of queries via the psql binary
func sendpsqladmincommndstobinary(pass string, queries []string) {
	const (
		NOTICE = `If you see: "DETAIL:  There are NN other sessions using the database." You need to quit out of HipparchiaGoServer
and/or anything else that is accessing the database. And if you do that and still have problems, then manually drop
the database:

% psql-17 postgres

	psql-17 (17.4 (Homebrew))
	Type "help" for help.
	
	postgres=# DROP DATABASE "hgdb";
	DROP DATABASE
	postgres=# \q`
	)

	binary := getpgbinarypath("psql")
	url := getpostgresadminuri(pass)

	for q := range queries {
		// not everything can run inside a transaction block (esp CREATE DATABASE)

		cmd := exec.Command(binary, "-c", queries[q], url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Could not execute psqladmin command:", queries[q])
			fmt.Println(err)
			if strings.Contains(err.Error(), "other sessions using the database") {
				fmt.Println(NOTICE)
			}
			// os.Exit(1)
		}
	}
}

// sendpsqlusercommandstobinary - execute a chain of queries via the psql binary
func sendpsqlusercommandstobinary(pass string, queries []string) {

	binary := getpgbinarypath("psql")
	url := getpostgresuseruri(pass)

	for q := range queries {
		cmd := exec.Command(binary, "-c", queries[q], url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		fmt.Println(err)
	}
}

// getpgbinarypath - return the path of a psql or pg_restore binary
func getpgbinarypath(command string) string {
	const (
		MACPGAPP = "/Applications/Postgres.app/Contents/Versions/%d/bin/"
		MACPGFD  = "/Applications/Postgres.app"
		MACBREW  = "/opt/homebrew/opt/postgresql@%d/bin/"
		WINPGEXE = `C:\Program Files\PostgreSQL\%d\bin\`
		LNXBIN   = `/usr/bin/`
		LNXLBIN  = `/usr/local/bin/`
		FAIL     = "Cannot find PostgreSQL binaries: aborting"
	)

	bindir := ""
	suffix := ""

	// linux and freebsd need fewer checks
	if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" {
		_, y := os.Stat(LNXBIN + command)
		if y == nil {
			// != nil will trigger a fail later
			return LNXBIN + command
		}
		_, y = os.Stat(LNXLBIN + command)
		if y == nil {
			// != nil will trigger a fail later
			return LNXLBIN + command
		}
	}

	// mac and windows are entangled with versioning issues
	if runtime.GOOS == "darwin" {
		_, y := os.Stat(MACPGFD)
		if y == nil {
			bindir = MACPGAPP
		} else {
			bindir = MACBREW
		}
	} else if runtime.GOOS == "windows" {
		bindir = WINPGEXE
		suffix = ".exe"
	}

	vers := 0

	for i := 21; i > 12; i-- {
		_, y := os.Stat(fmt.Sprintf(bindir, i) + command + suffix)
		if y == nil {
			vers = i
			break
		}
	}

	if vers == 0 {
		fmt.Println(FAIL)
		os.Exit(0)
	}

	bindir = fmt.Sprintf(bindir, vers)
	return bindir + command + suffix
}

// getpostgresadminuri - return a URI to connect to postgres as an administrator; different URI for macOS vs others
func getpostgresadminuri(pgpw string) string {
	const (
		UPWD = `postgresql://%s:%s@%s:%d/%s`
		UBLK = `postgresql://%s:%d/%s`
	)
	var url string
	if runtime.GOOS == "darwin" {
		// macos users have admin access already (on their primary account...) and do not need a pg admin password
		// postgresql://localhost:5432/postgres
		url = fmt.Sprintf(UBLK, pgsq.DEFAULTPSQLHOST, pgsq.DEFAULTPSQLPORT, "postgres")
	} else {
		// postgresql://postgres:password@localhost:5432/postgres
		url = fmt.Sprintf(UPWD, "postgres", pgpw, pgsq.DEFAULTPSQLHOST, pgsq.DEFAULTPSQLPORT, "postgres")
	}
	return url
}

// getpostgresadminuri - return a URI to connect to postgres as an administrator; different URI for macOS vs others
func getpostgresuseruri(userpw string) string {
	const (
		UPWD = `postgresql://%s:%s@%s:%d/%s`
	)
	url := fmt.Sprintf(UPWD, pgsq.DEFAULTPSQLUSER, userpw, pgsq.DEFAULTPSQLHOST, pgsq.DEFAULTPSQLPORT, pgsq.DEFAULTPSQLDB)
	return url
}
