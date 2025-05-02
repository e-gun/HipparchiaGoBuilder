//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package lexica

import (
	"io"
	"log"
	"os"
	"strings"
)

func LoadLexFile(directory string, filename string) string {
	file, err := os.Open(directory + filename)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	o, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(o)
}

func BuildLexDataFileNamesSlice(datadir string, suffix string) []string {
	entries, err := os.ReadDir(datadir)
	if err != nil {
		log.Fatal(err)
	}

	var xmls []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), suffix) {
			xmls = append(xmls, e.Name())
		}
	}
	return xmls
}
