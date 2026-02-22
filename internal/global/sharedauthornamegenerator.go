//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package global

import (
	"regexp"
	"strings"
	"sync"
)

var (
	UniqueAuthorNames = NewUniqueNamer(3)
	trailingdigit     = regexp.MustCompile("[0-9]$")
)

type UniqueNamer struct {
	Last int
	Max  int
	mu   sync.Mutex
}

func (n *UniqueNamer) GetName(pad int) string {
	n.mu.Lock()
	defer n.mu.Unlock()
	b36 := n.tobase36(n.Last)
	n.Last = n.Last + 1

	if n.Last > n.Max {
		n.Last = 0
	}

	needpad := pad - len(b36)
	for i := 0; i < needpad; i++ {
		b36 = "0" + b36
	}
	return b36
}

// GetNameWithTrailingDigit - because the "w" between AUID and WKID is easier to read this way
func (n *UniqueNamer) GetNameWithTrailingDigit(pad int) string {
	nm := n.GetName(pad)
	if trailingdigit.MatchString(nm) {
		return nm
	} else {
		return n.GetNameWithTrailingDigit(pad)
	}
}

func (n *UniqueNamer) NumberOfDigits(digits int) {
	// 1,296 available in a "2-digit" name
	// 46,655 available in a "3-digit" name
	// 1,679,615 available with 4-digits

	// `select count(*) from works;` --> 248016
	// `select count(*) from works where universalid ~* '^in';` --> 144419
	// `select count(*) from authors where universalid ~* '^in';` --> 463

	switch digits {
	case 0:
		n.Max = 0
	case 1:
		n.Max = 35
	case 2:
		n.Max = 1296
	case 3:
		n.Max = 46655
	case 4:
		n.Max = 1679615
	default:
		n.Max = 0
	}
}

func (n *UniqueNamer) tobase36(num int) string {
	if num == 0 {
		return "0"
	}

	// nb: cannot use upper case unless you are willing to work with postgres re case-sensitive table names
	const base36chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	var result strings.Builder

	// While the number is greater than 0, divide it by 36 and append the character
	for num > 0 {
		remainder := num % 36
		result.WriteByte(base36chars[remainder])
		num /= 36
	}

	// The result will be in reverse order, so reverse it before returning
	return n.reverse(result.String())
}

// Helper function to reverse a string
func (n *UniqueNamer) reverse(s string) string {
	var reversed strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		reversed.WriteByte(s[i])
	}
	return reversed.String()
}

func NewUniqueNamer(digits int) *UniqueNamer {
	nmr := &UniqueNamer{
		Last: 10, // there can be some NewSafeAuMap key deletion issues if you do not start at '0a', in impossible value
	}
	nmr.NumberOfDigits(digits)
	return nmr
}
