//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package generic

import "strings"

// ToSet - returns a blank map of a slice
func ToSet[T comparable](sl []T) map[T]struct{} {
	m := make(map[T]struct{})
	for i := 0; i < len(sl); i++ {
		m[sl[i]] = struct{}{}
	}
	return m
}

// Unique - return only the Unique items from a slice
func Unique[T comparable](s []T) []T {
	// can't use slices.Compact because that only looks as consecutive repeats: [a, a, b, a] -> [a, b, a]
	set := ToSet(s)

	result := make([]T, len(set))
	count := 0
	for k := range set {
		result[count] = k
		count += 1
	}

	return result
}

// DropEmptyStrings - get rid of empty slice values
func DropEmptyStrings(slc []string) []string {
	var result []string
	for _, item := range slc {
		if len(item) != 0 {
			result = append(result, item)
		}
	}
	return result
}

func TrimAllStrings(slc []string) []string {
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}
	return slc
}

// StringMapKeysIntoSlice - convert map[string]T to []string
func StringMapKeysIntoSlice[T any](mp map[string]T) []string {
	sl := make([]string, len(mp))
	i := 0
	for k := range mp {
		sl[i] = k
		i += 1
	}
	return sl
}

// SliceContains - is X in slice A?
func SliceContains[T comparable](sl []T, seek T) bool {
	for _, v := range sl {
		if v == seek {
			return true
		}
	}
	return false
}

// SliceOverlap - is any X in Y?
func SliceOverlap[T comparable](x []T, y []T) bool {
	for _, v := range x {
		if SliceContains(y, v) {
			return true
		}
	}
	return false
}
