//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lexica

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"github.com/e-gun/HipparchiaGoBuilder/pgsq/insert"
)

func ParseAndLoadLSJ(dir string, fn string) {
	// greatscott01.xml is just prefatory material
	// greatscott65.xml contains a single entry: Ϟ

	data := LoadLexFile(dir, fn)

	finder := regexp.MustCompile(`greatscott(\d\d).xml`)
	prefix, _ := strconv.Atoi(finder.ReplaceAllString(fn, "$1"))

	// fmt.Println(fn)

	data = FixSpecificLSJCitations(data)

	// rewrite the hyperlink style right off the bat
	data = generatehyperlinks(data)

	entries, spliterror := FindGreekLexEntries(data)
	if spliterror != nil && fn != "greatscott01.xml" {
		fmt.Println("no entries found for", fn)
	}

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

	entries = removeentrycollisions(validentries)

	if len(entries) == 0 && fn != "greatscott01.xml" {
		// you have a problem...
		fmt.Println("no entries for", fn)
	}

	// HGS wants to know next/previous entry; strings is not helpful, esp since n1111, n1111a, n1112 counting is in place
	// this puts them all in order despite MP disordering of the data: use the fn as the in and the entry # as the float
	// "wc -l *xml" shows that there are no original files longer than about 6300 lines
	for i := range entries {
		entries[i].IdFloat = float32(prefix) + float32(i)/10000
		// fmt.Printf("%d\t%f\t%s\n", i, entries[i].IdFloat, entries[i].EntryName)
		//1961	78.196098	ὑψιτέλεϲτοϲ
		//1962	78.196198	ὑψιτενέω
		//1725	80.172501	φληδάω
		//1726	80.172600	φληναφάω
		//1727	80.172699	φληνάφημα
	}

	err := insert.EntriesIntoGkLexicon(entries)
	if err != nil {
		fmt.Println(err)
	}
}

// removeentrycollisions - no longer an issue because H Dik has edited it away?
func removeentrycollisions(entries []structs.DbLexicon) []structs.DbLexicon {
	collider := make(map[string]bool)
	count := 0
	for i, entry := range entries {
		_, contains := collider[entry.IdString]
		if contains {
			// fmt.Println("removeentrycollisions() collision at", entry.IdString)
			entries[i].IdString = entry.IdString + fmt.Sprintf("_%d", count)
			// fmt.Println(entry.IdString + fmt.Sprintf("_%d", count))
			count++
		} else {
			collider[entry.IdString] = true
		}
	}
	return entries
}

// removesensecollisions - no longer an issue because H Dik has edited it away?
func removesensecollisions(senses []structs.LexicalSenses) []structs.LexicalSenses {
	collider := make(map[string]bool)
	count := 0
	for i, sense := range senses {
		_, contains := collider[sense.ID]
		if contains {
			// fmt.Println("removesensecollisions() collision at", sense.ID)
			senses[i].ID = sense.ID + fmt.Sprintf("_%d", count)
			// fmt.Println(sense.ID + fmt.Sprintf("_%d", count))
			count++
		} else {
			collider[sense.ID] = true
		}
		// supplementary check because someone is bad:
		if len(senses[i].ID) > 32 || len(senses[i].N) > 32 || len(senses[i].LVL) > 32 {
			fmt.Println("senses problem with", senses[i].ID, senses[i].N, senses[i].LVL)
		}
	}
	return senses
}

func generatehyperlinks(data string) string {
	// <bibl n="Perseus:abo:tlg,1342,001:23:6"> --> <bibl id="perseus/gr1342/001/23:6">

	const (
		TMPL = `<bibl id="perseus/%s%s/%s/%s">`
	)
	var (
		hyperlinkeditor = regexp.MustCompile(`<bibl n="Perseus:abo:tlg,(\d\d\d\d),(\d\d\d):(.*?)">`)
	)

	replacer := func(text string) string {
		groups := hyperlinkeditor.FindAllStringSubmatch(text, 1)
		return fmt.Sprintf(TMPL, global.TLGABBREV, groups[0][1], groups[0][2], groups[0][3])
	}

	data = hyperlinkeditor.ReplaceAllStringFunc(data, replacer)
	return data
}
