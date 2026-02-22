//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package betacode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

// ReplaceQuotationMarks() - clean up the ", "2, "3, ... markup

const (
	FAILEDQUOTE = `<hb-build-error>＂%s<hb-build-error>`
)

var (
	simplequotesubsa = map[int]string{
		1: "“",
		2: "QUOTE2",
		3: "QUOTE3",
		4: "‘",
		5: "’",
		6: "QUOTE6",
		7: "QUOTE7",
		8: "QUOTE8",
	}
	complexquotesubsa = map[int]string{
		1: "\u201e",
		2: "QUOTE2",
		3: "QUOTE3",
		4: "\u201a",
		5: "\u201b",
		6: "QUOTE6",
		7: "QUOTE7",
		8: "QUOTE8",
	}
	simplequotesubsb = map[int][]string{
		2: {"“", "”"},
		3: {"‵", "′"}, // Reversed prime and prime (for later fixing)
		6: {"“", "”"},
		7: {"‵", "′"}, // Reversed prime and prime (for later fixing)
		8: {"“", "”"},
	}
	complexquotesubsb = map[int][]string{
		2: {"“", "”"},
		3: {"‵", "′"}, // Reversed prime and prime (for later fixing)
		6: {"«", "»"},
		7: {"‹", "›"},
		8: {"“", "„"},
	}
	qreg0 = regexp.MustCompile(`"\d{1,2}`)
	qreq1 = regexp.MustCompile(`"(.*?)"`)
	qreg2 = regexp.MustCompile(`QUOTE\d.*?QUOTE\d`)
	qreg3 = regexp.MustCompile(`QUOTE(\d)(.*?)QUOTE\d`)
)

func ReplaceQuotationMarks(ttc string) string {
	// Format Quotation Marks
	// "50-59 Reserved for Greek documentary papyri
	// "60-69 Reserved for Greek inscriptions

	// brackets work like: text <14 enclosure >14 more text
	// quotes do not use the on/off mechanism

	// instead quote work like: text "7 enclosure "7 more text
	// so you have to catch both sides of the enclosure and cannot just do thing1 and thing2 when you see thingA and thingB

	// a. Purge " markup (first substitution)
	ttc = qreg0.ReplaceAllStringFunc(ttc, func(match string) string {
		return quotesubstitutesa(match)
	})

	// b. Replace text between *remaining* quotes with curly quotes
	ttc = qreq1.ReplaceAllString(ttc, "“$1”")

	// c. Handle QUOTE numbers substitution (second substitution)
	ttc = qreg2.ReplaceAllStringFunc(ttc, func(match string) string {
		return quotesubstitutesb(match)
	})

	// d. fix balance errors...
	// todo...

	return ttc
}

// quotesubstitutesa - Helper function for replacing the first type of quotation pattern
func quotesubstitutesa(match string) string {
	var substitutions map[int]string
	if SIMPLEQUOTES {
		substitutions = simplequotesubsa
	} else {
		substitutions = complexquotesubsa
	}

	match = strings.TrimPrefix(match, "\"")

	val, err := strconv.Atoi(match)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(FAILEDQUOTE, match)
	}

	substitute, exists := substitutions[val]
	if !exists {
		substitute = fmt.Sprintf(FAILEDQUOTE, match)
		global.MSG(fmt.Sprintf("quotesubstitutesa()\t%s", substitute))
	}
	return substitute
}

// quotesubstitutesb - Helper function for replacing the second type of quotation pattern
func quotesubstitutesb(match string) string {
	var substitutions map[int][]string
	if SIMPLEQUOTES {
		substitutions = simplequotesubsb
	} else {
		substitutions = complexquotesubsb
	}

	sm := qreg3.FindStringSubmatch(match)

	val, err := strconv.Atoi(sm[1]) // The first character is the number part (adjust if needed)
	if err != nil {
		// Handle error if the match is not a valid number
		return fmt.Sprintf(FAILEDQUOTE, match)
	}

	core := sm[2]

	substitution, exists := substitutions[val]
	var result string
	if exists {
		result = substitution[0] + core + substitution[1]
	} else {
		result = fmt.Sprintf(FAILEDQUOTE, match)
		global.MSG(fmt.Sprintf("quotesubstitutesa()\t%s", result))
	}

	return result
}
