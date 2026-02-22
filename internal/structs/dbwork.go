//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type DbWork struct {
	UID       string
	Title     string
	Language  string
	Pub       string
	LL0       string
	LL1       string
	LL2       string
	LL3       string
	LL4       string
	LL5       string
	Genre     string
	Xmit      string
	Type      string
	Prov      string
	RecDate   string
	ConvDate  int
	WdCount   int
	FirstLine int
	LastLine  int
	Authentic bool
}

// CountLevels - the work structure employs how many levels?
func (dbw *DbWork) CountLevels() int {
	ll := 0
	for _, l := range []string{dbw.LL5, dbw.LL4, dbw.LL3, dbw.LL2, dbw.LL1, dbw.LL0} {
		if len(l) > 0 {
			ll += 1
		}
	}
	return ll
}

// GetCorpus - ex: gr2017w068 --> gr
func (dbw *DbWork) GetCorpus() string {
	return dbw.UID[0:AbbrevLen]
}

// GetAuthor - ex: gr2017w068 --> gr2017
func (dbw *DbWork) GetAuthor() string {
	a := strings.Split(dbw.UID, AuthAndWorkSeparator)
	return a[0]
}

func (dbw *DbWork) PrintOut() {
	const (
		TMPL = `
	UID       {{.UID}}
	Language  {{.Language}}
	Pub       {{.Pub}}
	LL0       {{.LL0}}
	LL1       {{.LL1}}
	LL2       {{.LL2}}
	LL3       {{.LL3}}
	LL4       {{.LL4}}
	LL5       {{.LL5}}
	Genre     {{.Genre}}
	Xmit      {{.Xmit}}
	Type      {{.Type}}
	Prov      {{.Prov}}
	RecDate   {{.RecDate}}
	ConvDate  {{.ConvDate}}
	WdCount   {{.WdCount}}
	FirstLine {{.FirstLine}}
	LastLine  {{.LastLine}}
	Authentic {{.Authentic}}`
	)
	m := map[string]string{
		"UID":       dbw.UID,
		"Language":  dbw.Language,
		"Pub":       dbw.Pub,
		"LL0":       dbw.LL0,
		"LL1":       dbw.LL1,
		"LL2":       dbw.LL2,
		"LL3":       dbw.LL3,
		"LL4":       dbw.LL4,
		"LL5":       dbw.LL5,
		"Genre":     dbw.Genre,
		"Xmit":      dbw.Xmit,
		"Type":      dbw.Type,
		"Prov":      dbw.Prov,
		"RecDate":   dbw.RecDate,
		"ConvDate":  fmt.Sprintf("%d", dbw.ConvDate),
		"WdCount":   fmt.Sprintf("%d", dbw.WdCount),
		"FirstLine": fmt.Sprintf("%d", dbw.FirstLine),
		"LastLine":  fmt.Sprintf("%d", dbw.LastLine),
		"Authentic": fmt.Sprintf("%t", dbw.Authentic),
	}

	t := template.Must(template.New("").Parse(TMPL))

	var b bytes.Buffer
	if ee := t.Execute(&b, m); ee != nil {
		fmt.Println(ee)
	}
	fmt.Println(b.String())
}
