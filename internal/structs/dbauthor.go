//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"bytes"
	"fmt"
	"text/template"
)

type DbAuthor struct {
	UID       string
	Language  string
	IDXname   string
	Name      string
	Shortname string
	Cleaname  string
	Genres    string
	RecDate   string
	ConvDate  int
	Location  string
	// beyond the DB starts here
	WorkList []string
}

func (dba *DbAuthor) PrintOut() {
	const (
		TMPL = `
	UID       {{.UID}}
	Language  {{.Language}}
	IDXname   {{.IDXname}}
	Name      {{.Name}}
	Shortname {{.Shortname}}
	Cleaname  {{.Cleaname}}
	Genres    {{.Genres}}
	RecDate   {{.RecDate}}
	ConvDate  {{.ConvDate}}
	Location  {{.Location}}`
	)
	m := map[string]string{
		"UID":       dba.UID,
		"Language":  dba.Language,
		"IDXname":   dba.IDXname,
		"Name":      dba.Name,
		"Shortname": dba.Shortname,
		"Cleaname":  dba.Cleaname,
		"Genres":    dba.Genres,
		"RecDate":   dba.RecDate,
		"ConvDat":   fmt.Sprintf("%d", dba.ConvDate),
		"Location":  dba.Location,
	}

	t := template.Must(template.New("").Parse(TMPL))

	var b bytes.Buffer
	if ee := t.Execute(&b, m); ee != nil {
		fmt.Println(ee)
	}
	fmt.Println(b.String())
}
