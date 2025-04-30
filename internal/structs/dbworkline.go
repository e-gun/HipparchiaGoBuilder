//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"regexp"
	"strings"
)

var (
	findtrailinghypen = regexp.MustCompile(`-$`)
	findgreek         = regexp.MustCompile(`[α-ωϲ]`)
)

const (
	EMPTYLEVELINFO         = ""
	NUMBEROFCITATIONLEVELS = 6
	LENGTHOFAUTHORID       = 6
)

type DbWorkline struct {
	// WkUID       string
	TbIndex     int
	Lvl5Value   string
	Lvl4Value   string
	Lvl3Value   string
	Lvl2Value   string
	Lvl1Value   string
	Lvl0Value   string
	WkUID       string
	MarkedUp    string
	Accented    string
	Stripped    string
	Hyphenated  string
	Annotations string
}

func (dbl *DbWorkline) IsEmply() bool {
	// use Accented because MarkedUp might only contain markup/metadata
	if dbl.Accented == "" {
		return true
	} else {
		return false
	}
}

func (dbl *DbWorkline) HasHyphen() bool {
	if findtrailinghypen.MatchString(dbl.Accented) {
		return true
	} else {
		return false
	}
}

func (dbl *DbWorkline) HasGreek() bool {
	if findgreek.MatchString(dbl.Stripped) {
		return true
	} else {
		return false
	}
}

func (dbl *DbWorkline) FirstAccentedWord() string {
	accwords := strings.Split(dbl.Accented, " ")
	return accwords[0]
}

func (dbl *DbWorkline) LastAccentedWord() string {
	accwords := strings.Split(dbl.Accented, " ")
	return accwords[len(accwords)-1]
}

func (dbl *DbWorkline) AllButFirstAccentedWord() []string {
	accwords := strings.Split(dbl.Accented, " ")
	if len(accwords) > 2 {
		return accwords[1:]
	} else {
		return []string{}
	}
}

func (dbl *DbWorkline) AllButLastAccentedWord() []string {
	accwords := strings.Split(dbl.Accented, " ")
	if len(accwords) > 2 {
		return accwords[0 : len(accwords)-1]
	} else {
		return accwords
	}
}

func (dbl *DbWorkline) FirstStrippedWord() string {
	stwds := strings.Split(dbl.Stripped, " ")
	return stwds[0]
}

func (dbl *DbWorkline) LastStrippedWord() string {
	stwds := strings.Split(dbl.Stripped, " ")
	return stwds[len(stwds)-1]
}

func (dbl *DbWorkline) AllButFirstStrippedWord() []string {
	stwds := strings.Split(dbl.Stripped, " ")
	if len(stwds) > 2 {
		return stwds[1:]
	} else {
		return []string{}
	}
}

func (dbl *DbWorkline) AllButLastStrippedWord() []string {
	stwds := strings.Split(dbl.Stripped, " ")
	if len(stwds) > 2 {
		return stwds[0 : len(stwds)-1]
	} else {
		return stwds
	}
}

func (dbw *DbWorkline) FindLocus() []string {
	loc := [NUMBEROFCITATIONLEVELS]string{
		dbw.Lvl5Value,
		dbw.Lvl4Value,
		dbw.Lvl3Value,
		dbw.Lvl2Value,
		dbw.Lvl1Value,
		dbw.Lvl0Value,
	}

	var trim []string
	for _, l := range loc {
		if l != EMPTYLEVELINFO {
			trim = append(trim, l)
		}
	}
	return trim
}

// AuID - gr0001w001 --> gr0001
func (dbw *DbWorkline) AuID() string {
	return dbw.WkUID[:LENGTHOFAUTHORID]
}

// GetAccentedWordSlice - split up the accented words
func (dbw *DbWorkline) GetAccentedWordSlice() []string {
	return strings.Split(dbw.Accented, " ")
}

func (dbw *DbWorkline) GatherMetadata() map[string]string {
	md := make(map[string]string)

	// Plautus: "notes: Prisc. &3GL& 2.575K · #8"
	distinctnotes := strings.Split(dbw.Annotations, " · ")
	for _, dn := range distinctnotes {
		kv := strings.Split(dn, ":")
		if len(kv) != 2 {
			continue
		}
		md[kv[0]] = strings.TrimSpace(kv[1])
	}

	return md
}

func (dbw *DbWorkline) GetWordcount() int {
	return len(strings.Split(dbw.Stripped, " "))
}
