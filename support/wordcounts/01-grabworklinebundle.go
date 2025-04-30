package wordcounts

import (
	"context"
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/dbc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	WORLINETEMPLATE = `index, level_05_value, level_04_value, level_03_value, level_02_value, level_01_value, level_00_value,
			wkuniversalid, marked_up_line, accented_line, stripped_line, hyphenated_words, annotations`
)

func GetAllLinesFrom(authortable string) *structs.WorkLineBundle {
	const (
		QTMPL = "SELECT %s FROM %s ORDER by index"
	)
	var prq structs.PrerolledQuery
	prq.TempTable = ""
	prq.PsqlQuery = fmt.Sprintf(QTMPL, WORLINETEMPLATE, authortable)
	foundlines := getworklinebundle(prq)

	return foundlines
}

func GetSelectLinesFrom(authortable string, first int, last int) *structs.WorkLineBundle {
	const (
		QTMPL = "SELECT %s FROM %s WHERE index BETWEEN %d and %d ORDER by index"
	)
	var prq structs.PrerolledQuery
	prq.TempTable = ""
	prq.PsqlQuery = fmt.Sprintf(QTMPL, WORLINETEMPLATE, authortable, first, last)
	foundlines := getworklinebundle(prq)

	return foundlines
}

// getworklinebundle - wlbgrabber, but supply a dbconn for it via this function
func getworklinebundle(prq structs.PrerolledQuery) *structs.WorkLineBundle {
	// fanoutsearcher.go `consumeA` does not need to get a connection for every table query
	// it re-uses and sends its connections to wlbgrabber(); nobody else does that
	dbconn := dbc.GetDBConnection()
	defer dbconn.Release()

	return wlbgrabber(prq, dbconn)
}

// wlbgrabber - use a PrerolledQuery to acquire a *WorkLineBundle
func wlbgrabber(prq structs.PrerolledQuery, dbconn *pgxpool.Conn) *structs.WorkLineBundle {
	// NB: you have to use a dbconn.Exec() and can't use SQLPool.Exex() because with the latter the temp table will
	// get separated from the main query:
	// ERROR: relation "{ttname}" does not exist (SQLSTATE 42P01)

	// [a] build a temp table if needed

	if prq.TempTable != "" {
		_, err := dbconn.Exec(context.Background(), prq.TempTable)
		if err != nil {
			fmt.Println(err)
		}
	}

	// [b] execute the main query (nb: query needs to satisfy needs of RowToStructByPos in [c])

	foundrows, err := dbconn.Query(context.Background(), prq.PsqlQuery)
	if err != nil {
		fmt.Println(err)
	}

	// [c] convert the finds into []DbWorkline
	// not possible to convert this to []*DbWorkline; panic if pgx.RowToStructByPos[*structs.DbWorkline])

	thesefinds, err := pgx.CollectRows(foundrows, pgx.RowToStructByPos[structs.DbWorkline])
	if err != nil {
		fmt.Println("wlbgrabber()", err)
	}

	return &structs.WorkLineBundle{Lines: thesefinds}
}
