//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"fmt"
	"reflect"
	"sort"
)

type FingerPrint struct {
	OrigDateString      string
	Calculated          int
	ManualFudgeFactor   int
	Has1dArabic         bool
	Has2dArabic         bool
	Has3dArabic         bool
	Has4Arabic          bool
	HasMultiDigitArabic bool
	HasAnte             bool
	HasBeginning        bool
	HasCirca            bool
	HasEarly            bool // includes "in"
	HasEnd              bool
	HasGerman           bool
	HasHasEra           bool
	HasBCE              bool // "a", "BC", "B.C.",
	HasCE               bool
	HasLate             bool // includes "ex"
	HasMid              bool
	HasOneArabic        bool
	HasOneRoman         bool
	HasTH               bool
	HasOR               bool
	HasPaulo            bool
	HasPost             bool
	HasQuestion         bool
	HasTwoArabic        bool
	HasTwoRoman         bool
	HasThreeArabic      bool
	HasOneSpan          bool
	HasTwoSpans         bool
	HasThreeSpans       bool
	HasFourSpans        bool
	HasMixedSpans       bool
	HasSlashDate        bool
	HasMultiSlashDate   bool
	HasBracket          bool
	HasAnd              bool
	ContainsNoDigits    bool
	ParserFailed        bool
	LacksParser         bool
}

type FieldBoolValuePair struct {
	Field string
	Value bool
}

func (fp *FingerPrint) Rationalize() {
	// number of blocks
	if fp.HasThreeArabic {
		fp.HasTwoArabic = false
		fp.HasOneArabic = false
	}
	if fp.HasTwoArabic {
		fp.HasOneArabic = false
	}
	if fp.HasTwoRoman {
		fp.HasOneRoman = false
	}

	// digit counter
	if fp.Has2dArabic || fp.Has3dArabic || fp.Has4Arabic {
		fp.HasMultiDigitArabic = true
	}
	if fp.Has4Arabic {
		fp.Has3dArabic = false
		fp.Has2dArabic = false
		fp.Has1dArabic = false
	}
	if fp.Has3dArabic {
		fp.Has2dArabic = false
		fp.Has1dArabic = false
	}
	if fp.Has2dArabic {
		fp.Has1dArabic = false
	}

	// span counter
	if fp.HasTwoSpans {
		fp.HasOneSpan = false
	}
}

func (fp *FingerPrint) ApplyBCE() {
	if fp.HasBCE && !fp.HasCE {
		fp.Calculated = fp.Calculated * -1
	}
}

func (fp *FingerPrint) ApplySimpleFudges() {
	if fp.HasBeginning && !fp.HasEnd {
		fp.Calculated = fp.Calculated - 75
	} else if fp.HasEnd && !fp.HasBeginning {
		fp.Calculated = fp.Calculated + 25
	}
	if fp.HasAnte && !fp.HasPost {
		fp.Calculated = fp.Calculated - 15
	} else if !fp.HasAnte && fp.HasPost {
		fp.Calculated = fp.Calculated + 15
	}
	fp.Calculated = fp.Calculated + fp.ManualFudgeFactor
}

func (fp *FingerPrint) PrepPrint() []FieldBoolValuePair {
	var pairs []FieldBoolValuePair
	v := reflect.ValueOf(fp)
	startfield := 3
	stopfield := v.NumField()

	for i := startfield; i < stopfield; i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.Field(i).Interface().(bool) // Assuming all fields are bool
		pairs = append(pairs, FieldBoolValuePair{fieldName, fieldValue})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Field > pairs[j].Field
	})
	return pairs
}

func (fp *FingerPrint) Print() {
	fvp := fp.PrepPrint()
	fmt.Println(fp.OrigDateString)
	for _, pair := range fvp {
		if pair.Value {
			fmt.Printf("\t%s\n", pair.Field)
		}
	}
}

func (fp *FingerPrint) PrintWithCalc() {
	if fp.Calculated > 5000 {
		fp.Print()
		if fp.ParserFailed {
			fmt.Println("Parser Failed")
		}
	} else {
		fmt.Printf("`%s`", fp.OrigDateString)
		fmt.Printf(" --> %d\n", fp.Calculated)
	}
}
