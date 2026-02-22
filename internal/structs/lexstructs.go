//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"
)

type DbLexicon struct {
	IdFloat    float32
	EntryName  string
	EntryMetr  string
	IdString   string
	EntryType  string
	Transl     string
	Usedby     string
	PrelimInfo string
	SenseIDs   []string
	Senses     []LexicalSenses
}

func (ent *DbLexicon) PrintOut() {
	const (
		TMPL = `
	EntryName     {{.EntryName}}
	EntryMetr     {{.EntryMetr}}
	IdString         {{.IdString}}
	Transl        {{.Transl}}
	Usedby        {{.Usedby}}
	SenseIDs      {{.SenseIDs}}`
	)
	m := map[string]string{
		"EntryName": ent.EntryName,
		"EntryMetr": ent.EntryMetr,
		"IdString":  ent.IdString,
		"Transl":    ent.Transl,
		"Usedby":    ent.Usedby,
		"SenseIDs":  strings.Join(ent.SenseIDs, ", "),
	}

	t := template.Must(template.New("").Parse(TMPL))

	var b bytes.Buffer
	if ee := t.Execute(&b, m); ee != nil {
		fmt.Println(ee)
	}
	fmt.Println(b.String())
}

func (ent *DbLexicon) SensesJSONOut() []byte {
	pmap := make(map[int]LexicalSenses)
	for i, p := range ent.Senses {
		pmap[i] = p
	}

	js, err := json.Marshal(pmap)
	if err != nil {
		log.Fatal(err)
	}
	return js
}

type LexicalSenses struct {
	ID       string `json:"id"`
	N        string `json:"n"`
	LVL      string `json:"lvl"`
	Contents string `json:"contents"`
}
