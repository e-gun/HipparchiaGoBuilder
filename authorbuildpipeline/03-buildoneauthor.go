//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/grk"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/latetidyups"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/worklines1initial"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/worklines2dbprep"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/worklines3cleanup"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/idtandbin"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"os"
	"strings"
	"time"
)

// BuildOneAuthor - the workhorse: this runs the full authorprep sequence on a TXT file: betacode into worklines ... insert into DB
func BuildOneAuthor(datadir string, authorid string) {
	const (
		MSG = "%s\t %.3fs"
	)

	start := time.Now()

	// fmt.Printf("BuildOneAuthor(%s)\n", authorid)

	au, wkk := prepareinitialauthorandworkstructs(datadir, authorid)

	isphi := false
	if authorid[0:3] == "INS" || authorid[0:3] == "CHR" || authorid[0:3] == "DDP" {
		isphi = true
	}

	islat := strings.Contains(au.IDXname, "Latin")

	//L: in0025 Peloponnesos [Latin]
	//L: in0015 Attica [Latin]
	//L: in0085 Black Sea [Chersonesos] [Latin]
	//L: ch0120 Latin [other]
	//L: ch0130 Late Antique Latin
	//L: ch0140 Medieval Latin

	var ttc string
	if authorid[0:3] == "LAT" || islat {
		ttc = initiallatinauthortextfilecleaning(datadir, authorid)
	} else {
		ttc = initialgreekauthortextfilecleaning(datadir, authorid)
	}

	if islat && isphi {
		ttc = latetidyups.PHIAuthorExclamations(ttc)
	}

	if authorid[0:3] == "CHR" && !islat {
		// a problem: "CHR" is mostly latin but some greek; latin cleaning is good until it ruins the greek
		// greek cleaning is the only way to get greek right, but every capital letter will be wrong...
		// a lingering issue will be roman numerals: MCXVI as a series of greek lc letters

		// should consider building an list of files that should be treated as pure latin: 0130, for example

		ttc = latetidyups.PurgeHybrid(ttc)
	}

	if isphi {
		inscriptionorklineprepandcleanupandinsertion(ttc, wkk, authorid)
	} else {
		greekandlatinworklineprepandcleanupandinsertion(ttc, wkk, authorid)
	}

	d := fmt.Sprintf(MSG, authorid, time.Now().Sub(start).Seconds())
	fmt.Println(d)
}

func getcorpusabbrev(authorid string) string {
	pfr := []rune(authorid)
	pfx := string(pfr[0:3])

	var corpusabbrev string

	switch pfx {
	case "TLG":
		corpusabbrev = "gr"
	case "LAT":
		corpusabbrev = "lt"
	case "DDP":
		corpusabbrev = "dp"
	case "CHR":
		corpusabbrev = "ch"
	case "INS":
		corpusabbrev = "in"
	default:
		fmt.Println("Unknown corpus abbrev", pfx)
	}
	return corpusabbrev
}

func prepareinitialauthorandworkstructs(datadir string, authorid string) (structs.DbAuthor, []structs.DbWork) {
	idtdata := idtandbin.BINFileLoad(datadir + authorid + ".IDT")
	abbr := getcorpusabbrev(authorid)
	au, wkk := idtandbin.LoadAuthorIDTData(idtdata)
	idtandbin.CleanDBAuthorNames(&au)
	idtandbin.AssignPrefixToAuthor(abbr, &au)
	for i := range wkk {
		idtandbin.CleanDBWorkNames(&wkk[i])
		idtandbin.AssignPrefixToWork(abbr, &wkk[i])
	}
	global.TheIdtAuMap.Set(au.UID, au)
	return au, wkk
}

func initialgreekauthortextfilecleaning(datadir string, authorid string) string {
	rawtext := idtandbin.HighUnicodeFileLoad(datadir + authorid + ".TXT")
	ttc := betacode.BetaCodeCleanup(rawtext)
	ttc = grk.GreekCleanup(ttc)
	ttc = latetidyups.TidyUp(ttc)
	ttc = worklines1initial.WorklinePrep(ttc)
	return ttc
}

func initiallatinauthortextfilecleaning(datadir string, authorid string) string {
	rawtext := idtandbin.HighUnicodeFileLoad(datadir + authorid + ".TXT")
	ttc := betacode.BetaCodeCleanup(rawtext)
	ttc = lat.LatinCleanup(ttc)
	ttc = latetidyups.TidyUp(ttc)
	ttc = worklines1initial.WorklinePrep(ttc)
	return ttc
}

func greekandlatinworklineprepandcleanupandinsertion(ttc string, wkk []structs.DbWork, authorid string) {
	var dbworklines []structs.DbWorkline
	abbr := getcorpusabbrev(authorid)

	dbworklines = worklines2dbprep.PrepareForDB(ttc)
	dbworklines = worklines3cleanup.FixWorklineMetadata(abbr, dbworklines)
	dbworklines = worklines3cleanup.GRKandLATFixInvalidLevelValues(dbworklines, wkk)
	dbworklines = reindexlines(dbworklines)

	err := insert.InsertWorklinesIntoTable(dbworklines)
	if err != nil {
		fmt.Println(authorid, " warning:", err, "Creating empty table.")
		insert.CreateAuthorTable(authorid)
	}

	if abbr == "gr" {
		a, _ := global.TheCanonAuMap.Get(wkk[0].GetAuthor())
		for i, _ := range wkk {
			wkk[i].RecDate = a.RecDate
			wkk[i].ConvDate = a.ConvDate
		}
	}

	// update the relevant dbwork info while we have it on hand...>
	for i, wk := range wkk {
		wkk[i].FirstLine = findworkstart(wk.UID, dbworklines)
		wkk[i].LastLine = findworkend(wk.UID, dbworklines)
		wkk[i].WdCount = findwordcount(wk.UID, dbworklines)
	}

	// do not insert into dbc; just store this author and work info for later reconciliation by ReconcileWkMaps(), etc
	for _, w := range wkk {
		global.TheIdtWkMap.Set(w.UID, w)
	}

	return
}

func inscriptionorklineprepandcleanupandinsertion(ttc string, wkk []structs.DbWork, authorid string) {
	var dbworklines []structs.DbWorkline

	// authorid looks like "INS0150"
	abbr := getcorpusabbrev(authorid)
	dbworklines = worklines2dbprep.PrepareForDB(ttc)
	dbworklines = worklines3cleanup.FixWorklineMetadata(abbr, dbworklines)
	dbworklines = reindexlines(dbworklines)

	// GRKandLATFixInvalidLevelValues only works for "lt" and "gr"; PhiValueStripper() below is the equivalent
	// dbworklines = worklines2dbprep.GRKandLATFixInvalidLevelValues(dbworklines, wkk)

	oldau, _ := global.TheIdtAuMap.Get(dbworklines[0].WkUID[0:6])
	global.TheDefunctIdtAuMap.Set(oldau.UID, oldau)
	global.TheIdtAuMap.Delete(oldau.UID)

	var auu []structs.DbAuthor
	auu, wkk, dbworklines = worklines3cleanup.RemapInscriptionAuthorsAndWorks(dbworklines, authorid)

	// do this after RemapInscriptionAuthorsAndWorks because you need the bad values until then
	dbworklines = worklines3cleanup.PhiValueStripper(dbworklines)

	// map the lines onto author tables
	newauthorlines := make(map[string][]structs.DbWorkline)
	var brokenworks bool
	for _, l := range dbworklines {
		if len(l.WkUID) == 0 {
			// this is in fact a disaster for the database if you see it
			fmt.Println("No WkUID for", l.TbIndex, l.MarkedUp)
			brokenworks = true
			// do not exit() right away; accumulate the bad data in the console
			continue
		}
		newauthorlines[l.AuID()] = append(newauthorlines[l.WkUID[0:6]], l)
	}

	if brokenworks {
		fmt.Println("Broken works generated by", oldau)
		fmt.Println("Aborting...")
		os.Exit(1)
	}

	for _, ll := range newauthorlines {
		_ = insert.InsertWorklinesIntoTable(ll)
	}

	// inform the global maps of the changes
	for _, a := range auu {
		global.TheIdtAuMap.Set(a.UID, a)
	}
	for _, w := range wkk {
		global.TheIdtWkMap.Set(w.UID, w)
	}

	return
}

func findworkend(id string, lines []structs.DbWorkline) int {
	end := 0
	for i := 0; i < len(lines); i++ {
		if lines[i].WkUID == id && lines[i].TbIndex > end {
			end = lines[i].TbIndex
		}
	}
	// fmt.Println(id, end)
	return end
}

func findworkstart(id string, lines []structs.DbWorkline) int {
	// to find problems:
	// select universalid,title,firstline,wordcount from works where firstline = 9999999 order by universalid asc;

	// InsertNewLinesAtNewLevels() has to be right, for example...

	start := 9999999
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i].WkUID == id && lines[i].TbIndex < start {
			start = lines[i].TbIndex
		}
	}
	// fmt.Println(id, start)
	return start
}

func findwordcount(id string, lines []structs.DbWorkline) int {
	// slow-ish; adds 5-10% to the build time of a corpus
	count := 0
	for i := 0; i < len(lines); i++ {
		if lines[i].WkUID == id {
			sp := strings.Split(lines[i].Stripped, " ")
			for s := range sp {
				if sp[s] != "" {
					count++
				}
			}
		}
	}
	return count
}

// reindexlines - blank line deletions have left a non-contiguous collection of index numbers
func reindexlines(lines []structs.DbWorkline) []structs.DbWorkline {
	for i := 0; i < len(lines); i++ {
		lines[i].TbIndex = i
	}
	return lines
}
