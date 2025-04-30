package wordcounts

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/generic"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

func CountThisBundle(priorcounts map[string]int, bundle *structs.WorkLineBundle) map[string]int {
	lines := bundle.Yield()
	for line := range lines {
		words := line.GetAccentedWordSlice()
		for _, word := range words {
			word = generic.SwapAcuteForGrave(word)
			if count, ok := priorcounts[word]; !ok {
				priorcounts[word] = 1
			} else {
				priorcounts[word] = count + 1
			}
		}
	}
	return priorcounts
}

func ParseAndCountThisBundle(priorcounts map[string]int, bundle *structs.WorkLineBundle) map[string]int {
	lines := bundle.Yield()
	for line := range lines {
		words := line.GetAccentedWordSlice()
		for _, unparsedword := range words {
			unparsedword = generic.SwapAcuteForGrave(unparsedword)
			possibilities, ok := HeadwordLookupMap[unparsedword]
			if !ok {
				continue
			}
			for _, possibility := range possibilities {
				if count, already := priorcounts[possibility]; !already {
					priorcounts[possibility] = 1
				} else {
					priorcounts[possibility] = count + 1
				}
			}
		}
	}
	return priorcounts
}
