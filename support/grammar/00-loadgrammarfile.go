//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grammar

import (
	"io"
	"os"
	"strings"
)

func LoadGrammarFile(directory string, filename string) []string {
	file, err := os.Open(directory + filename)

	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		ee := file.Close()
		if ee != nil {

		}
	}(file)

	o, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(o), "\n")
}
