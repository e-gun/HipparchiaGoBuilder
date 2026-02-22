//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package grk

import (
	"fmt"
	"regexp"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

var (
	insetcoptic     = regexp.MustCompile(`(<hb-fs-g-coptic>)(.*?)(&|</hb-fs-g-coptic>)`)
	copticcapfinder = regexp.MustCompile(`\*([a-zA-Z])`)
	copticlcfinder  = regexp.MustCompile(`[a-zA-Z]`)
	copticuppermap  = map[string]string{
		"A": "\u2c80",
		"B": "\u2c82",
		"C": "\u2c9c",
		"D": "\u2c86",
		"E": "\u2c86",
		"F": "\u2caa",
		"f": "\u03e4",
		"G": "\u2c84",
		"g": "\u03ec",
		"H": "\u2c8e",
		"h": "\u03e8",
		"I": "\u2c92",
		"j": "\u03ea",
		"K": "\u2c94",
		"L": "\u2c96",
		"M": "\u2c98",
		"N": "\u2c9a",
		"O": "\u2c9e",
		"P": "\u2ca0",
		"Q": "\u2c90",
		"R": "\u2ca2",
		"S": "\u2ca4",
		"s": "\u03e2",
		"T": "\u2ca6",
		"t": "\u03ee",
		"U": "\u2ca8",
		"V": "\u2c8a",
		"W": "\u2cb0",
		"X": "\u2cac",
		"Y": "\u2cae",
		"Z": "\u2c8c",
	}
	copticlowermap = map[string]string{
		"A": "\u2c81",
		"B": "\u2c83",
		"C": "\u2c9d",
		"D": "\u2c87",
		"E": "\u2c89",
		"F": "\u2cab",
		"f": "\u03e5",
		"G": "\u2c85",
		"g": "\u03ed",
		"H": "\u2c8f",
		"h": "\u03e9",
		"I": "\u2c93",
		// "J": "\u03eb",  // ?! DDP0088 has a "J"
		"j": "\u03eb",
		"K": "\u2c95",
		"L": "\u2c97",
		"M": "\u2c99",
		"N": "\u2c9b",
		"O": "\u2c9f",
		"P": "\u2ca1",
		"Q": "\u2c91",
		"R": "\u2ca3",
		"S": "\u2ca5",
		"s": "\u03e3",
		"T": "\u2ca7",
		"t": "\u03ef",
		"U": "\u2ca9",
		"V": "\u2c8b",
		"W": "\u2cb1",
		"X": "\u2cad",
		"Y": "\u2caf",
		"Z": "\u2c8d",
		"a": "\u2c81", // inscriptions were producing lots of copticlowercases() ; the following almost certainly breaks something
		"b": "\u2c83", // we are just going to guess that all the INS0180 lowers were really uppers
		"c": "\u2c9d",
		"d": "\u2c87",
		"e": "\u2c89",
		//"F": "\u2cab",
		//"f": "\u03e5",
		//"G": "\u2c85",
		//"g": "\u03ed",
		//"H": "\u2c8f",
		//"h": "\u03e9",
		"i": "\u2c93",
		"J": "\u03eb", // ?! DDP0088 has a "J"
		// "j": "\u03eb",
		"k": "\u2c95",
		"l": "\u2c97",
		"m": "\u2c99",
		"n": "\u2c9b",
		"o": "\u2c9f",
		"p": "\u2ca1",
		"q": "\u2c91",
		"r": "\u2ca3",
		//"S": "\u2ca5",
		//"s": "\u03e3",
		//"T": "\u2ca7",
		//"t": "\u03ef",
		"u": "\u2ca9",
		"v": "\u2c8b",
		"w": "\u2cb1",
		"x": "\u2cad",
		"y": "\u2caf",
		"z": "\u2c8d",
	}
)

func InstertCoptic(ttc string) string {
	// Replace the matches using the 'copticprobe' function
	cleaned := insetcoptic.ReplaceAllStringFunc(ttc, copticprobe)
	return cleaned
}

func copticprobe(match string) string {
	// Define the regex to capture the parts of the match
	// We use a regex to simulate what Python's match.group() does.
	matches := insetcoptic.FindStringSubmatch(match)

	if len(matches) == 4 {
		opentag := matches[1]
		body := matches[2]
		closetag := matches[3]

		// Find capital letters and process them
		body = copticcapfinder.ReplaceAllStringFunc(body, func(s string) string {
			return copticuppercases(string(s[1]))
		})

		// Replace all alphabetic characters with lowercase processed version
		body = copticlcfinder.ReplaceAllStringFunc(body, func(s string) string {
			return copticlowercases(s)
		})

		// Combine the open tag, processed body, and close tag
		return fmt.Sprintf("%s%s%s", opentag, body, closetag)
	}

	// If the match does not follow the expected format, return the original match
	return match
}

func copticuppercases(swap string) string {
	// Look up the substitution for the given toReplace character
	substitute, found := copticuppermap[swap]
	if !found {
		// If no substitution is found, return the original character
		global.MSG(fmt.Sprintf("copticuppercases() confusion: %s", swap))
		return swap
	}

	return substitute
}

func copticlowercases(swap string) string {
	// Look up the substitution for the given toReplace character
	substitute, found := copticlowermap[swap]
	if !found {
		// If no substitution is found, return the original character
		global.MSG(fmt.Sprintf("copticlowercases() confusion: %s", swap))
		return swap
	}
	return substitute
}
