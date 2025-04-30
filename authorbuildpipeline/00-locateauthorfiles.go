//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"log"
	"os"
	"strings"
)

var (
	exclusions = map[string]bool{
		"INS8000": true, // Delphi Bibliography
		"INS9900": true, // Bibliography [Epigr., general]
		"INS9910": true, // Bibliography [Epigr., Mysia and Troas [Munich]]
		"INS9920": true, // Bibliography [Epigr., Ionia]
		"INS9930": true,
		"CHR0021": true, // Judaica [Hebrew/Aramaic]
		"CHR9900": true, // Bibliography [Epigr., late ant./med.]
		"CHR9910": true, // Bibliography [Epigr., other Latin]
		"DDP9999": true, // Papyrological Guide
	}
	// FYI: COP0001 and COP0002 seem to contain only biblical texts
	// CIV0001 is the Herbrew bible; CIV0002 contains other bible versions; CIV003 is Greek New Testament (NT UBS3 edition)
	// CIV0004 is the Latin vulgate; CIV0005 English Bible (KJV or AV); CIV006 English Bible (RSV); CIV007 John Milton (English and Latin)
)

func BuildAuthorsSlice(datadir string, prefix string) []string {
	const (
		SUFFIX = ".TXT"
	)

	entries, err := os.ReadDir(datadir)
	if err != nil {
		log.Fatal(err)
	}

	var authors []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), SUFFIX) {
			authors = append(authors, strings.TrimSuffix(e.Name(), SUFFIX))
		}
	}

	var pruned []string
	for _, author := range authors {
		_, isexcluded := exclusions[author]
		if !isexcluded {
			pruned = append(pruned, author)
		}
	}

	return pruned
}
