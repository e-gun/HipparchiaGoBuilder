//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lexica

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
	"regexp"
	"strings"
)

func ParseAndLoadLatinLex(id int, splitdata []string) {
	// about to split all over again soon, though...
	data := strings.Join(splitdata, "\n")

	data = FixSpecificLatinCitations(data)

	// rewrite the hyperlink style right off the bat; FixSpecificLatinCitations() must execute before this
	data = generatelatinhyperlinks(data)

	entries := FindLatinLexEntries(data)

	// remove blanks
	var validentries []structs.DbLexicon
	count := 0
	for _, entry := range entries {
		if entry.EntryName != "" {
			validentries = append(validentries, entry)
		} else {
			count++
		}
	}

	// HGS wants to know next/previous entry; strings is not helpful, esp since n1111, n1111a, n1112 counting is in place
	// this puts them all in order despite MP disordering of the data: use the fn as the in and the entry # as the float
	// wc -l *xml
	//   51857 latin-lexicon_1999.04.0059.xml
	// that means 6 cpus only need to divide by 10k
	// fewer cores is 100k; but you might be doing an "-rp" run: so 1 core and 100K...
	// the corner case comes if you look for the next word at one of the boundaries of the workpiles...
	for i, _ := range entries {
		entries[i].IdFloat = float32(id) + float32(i)/100000
		// fmt.Printf("%d\t%f\t%s\n", i, entries[i].IdFloat, entries[i].EntryName)
		//1313	14.013130	spuo
		//1314	14.013140	spurcalia
		//2548	13.025480	Zoilus
		//2549	13.025490	zomoteganite
		//2550	13.025500	zona
		//2551	13.025510	zonalis
		//2552	13.025520	zonarius
	}

	err := insert.InsertEntriesIntoLatinLexicon(entries)
	if err != nil {
		fmt.Println(err)
	}

}

func generatelatinhyperlinks(data string) string {
	// <bibl n="Perseus:abo:tlg,1342,001:23:6"> --> <bibl id="perseus/gr1342/001/23:6">
	// <bibl n="Perseus:abo:phi,0474,037:2:60 al" default="NO" valid="yes"> --> <bibl id="perseus/lt0474/037/2:60">
	const (
		TMPL = `<bibl id="perseus/lt%s/%s/%s">`
	)
	var (
		hyperlinkeditor = regexp.MustCompile(`<bibl n="Perseus:abo:phi,(\d\d\d\d),(\d\d\d):(.+?)">`)
	)

	replacer := func(text string) string {
		groups := hyperlinkeditor.FindAllStringSubmatch(text, 1)
		return fmt.Sprintf(TMPL, groups[0][1], groups[0][2], groups[0][3])
	}

	data = hyperlinkeditor.ReplaceAllStringFunc(data, replacer)
	return data
}
