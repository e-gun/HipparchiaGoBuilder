//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	// needs to match the inverse function in hexrunner.go: highunicodetohex()
	// regex: "(.)", "(.)", --> "$2", "$1",
	hexreplacer = strings.NewReplacer(
		"0", "⓪",
		"1", "①",
		"2", "②",
		"3", "③",
		"4", "④",
		"5", "⑤",
		"6", "⑥",
		"7", "⑦",
		"8", "⑧",
		"9", "⑨",
		"a", "ⓐ",
		"b", "ⓑ",
		"c", "ⓒ",
		"d", "ⓓ",
		"e", "ⓔ",
		"f", "ⓕ",
	)
)

// IdtFileLoad opens an IDT file and prepares it for parsing
func IdtFileLoad(filepath string) string {
	file, err := os.Open(filepath)
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

	txt := make([]string, 0, len(o))
	for _, c := range o {
		if c >= 128 || c <= 31 {
			x := fmt.Sprintf("%x", c)
			x = "█" + hextohighunicode(x)
			txt = append(txt, x+" ")
		} else {
			txt = append(txt, string(c))
		}
	}

	return strings.Join(txt, "")
	// return string(o)
}

// BINFileLoad opens a BIN file and prepares it for parsing (super-vanilla at the moment)
func BINFileLoad(filepath string) []byte {
	// fmt.Println("BINFileLoad", filepath)
	file, err := os.Open(filepath)
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

	return o
}

// hextohighunicode - Returns the hex sequence corresponding to the given high unicode string
func hextohighunicode(highunicode string) string {
	return hexreplacer.Replace(highunicode)
}

// HighUnicodeFileLoad opens a file and prepares it for parsing
// Returns a collection of characters with the unprintable chars swapped out for their hex representation
func HighUnicodeFileLoad(filepath string) string {
	file, err := os.Open(filepath)
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

	txt := make([]string, 0, len(o))
	for _, c := range o {
		if c >= 128 || c <= 31 {
			x := fmt.Sprintf("%x", c)
			x = "█" + hextohighunicode(x)
			txt = append(txt, x+" ")
		} else {
			txt = append(txt, string(c))
		}
	}

	// fmt.Println(filepath, len(txt), "characters")
	return strings.Join(txt, "")
}
