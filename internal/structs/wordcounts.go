//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"fmt"
	"reflect"
	"sort"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type DbHeadwordTimeCounts struct {
	Early  int
	Middle int
	Late   int
}

type DbUnparsedWordCounts struct {
	Word  string
	Total int
	TLG   int
	LAT   int
	DDP   int
	INS   int
	CHR   int
}

func (wc DbUnparsedWordCounts) PrintOut() {
	const (
		TMPL = "%s\tt: %d\tg: %d\tl: %d\ti: %d\td: %d\tc: %d\n"
	)
	fmt.Printf(TMPL, wc.Word, wc.Total, wc.TLG, wc.LAT, wc.INS, wc.DDP, wc.CHR)
}

type DbHeadwordCounts struct {
	Word     string
	Total    int
	TLG      int
	LAT      int
	DDP      int
	INS      int
	CHR      int
	FrqClas  string
	Early    int
	Middle   int
	Late     int
	Acta     int
	Agric    int
	Alchem   int
	Anthol   int
	Apocal   int
	Apocry   int
	Apol     int
	Astrol   int
	Astron   int
	Biogr    int
	Bucol    int
	Caten    int
	Chron    int
	Comic    int
	Comm     int
	Concil   int
	Coq      int
	Dial     int
	Docu     int
	Doxog    int
	Eccl     int
	Eleg     int
	Encom    int
	Epic     int
	Epigr    int
	Epist    int
	Evang    int
	Exeg     int
	Fab      int
	Geog     int
	Gnom     int
	Gram     int
	Hagiog   int
	Hexam    int
	Hist     int
	Homil    int
	Hymn     int
	Hypoth   int
	Iamb     int
	Ignot    int
	Inscr    int
	Invectiv int
	Juris    int
	Lexic    int
	Liturg   int
	Lyr      int
	Magica   int
	Math     int
	Mech     int
	Med      int
	Meteor   int
	Mim      int
	Mus      int
	Myth     int
	NarrFic  int
	NatHis   int
	Onir     int
	Orac     int
	Orat     int
	Paradox  int
	Papyrus  int
	Parod    int
	Paroem   int
	Perig    int
	Phil     int
	Physiog  int
	Poem     int
	Polyhist int
	Proph    int
	Pseud    int
	Rhet     int
	Satura   int
	Satyr    int
	Schol    int
	Tact     int
	Test     int
	Theol    int
	Trag     int
	AllRhet  int
	AllRelig int
}

func (wc DbHeadwordCounts) PrintOut() {
	m := message.NewPrinter(language.English)
	fmt.Println("DbHeadwordCounts:", wc.Word)
	fmt.Println("Epic:", m.Sprintf("%d", wc.Epic))
	fmt.Println("Phil:", m.Sprintf("%d", wc.Phil))
	fmt.Println("Trag:", m.Sprintf("%d", wc.Trag))
	fmt.Println("AllRhet:", m.Sprintf("%d", wc.AllRhet))
	fmt.Println("AllRelig:", m.Sprintf("%d", wc.AllRelig))
	fmt.Println("Total:", m.Sprintf("%d", wc.Total))
}

type FieldValuePair struct {
	Field string
	Value int
}
type WeightedFieldValuePair struct {
	Field string
	Value float32
}

func (wc DbHeadwordCounts) SortedFVPairs(startfield int, stopfield int) []FieldValuePair {
	var pairs []FieldValuePair
	v := reflect.ValueOf(wc)
	for i := startfield; i < stopfield; i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.Field(i).Interface().(int) // Assuming all fields are int
		pairs = append(pairs, FieldValuePair{fieldName, fieldValue})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})
	return pairs
}

func (wc DbHeadwordCounts) SortedEraPairs() []FieldValuePair {
	startindex := 8
	endindex := 11
	return wc.SortedFVPairs(startindex, endindex)
}

func (wc DbHeadwordCounts) SortedCorpusPairs() []FieldValuePair {
	startindex := 2
	endindex := 7
	return wc.SortedFVPairs(startindex, endindex)
}

func (wc DbHeadwordCounts) SortedGenrePairs() []FieldValuePair {
	startindex := 11
	endindex := 91
	return wc.SortedFVPairs(startindex, endindex)
}

func (wc DbHeadwordCounts) SortedWeightedPairs(fvp []FieldValuePair) []WeightedFieldValuePair {
	maxval := float32(fvp[0].Value)
	wfvp := make([]WeightedFieldValuePair, len(fvp))
	for i, pair := range fvp {
		wfvp[i] = WeightedFieldValuePair{
			Field: pair.Field,
			Value: maxval / float32(pair.Value),
		}
	}
	return wfvp
}

func (wc DbHeadwordCounts) SortedWeightedEraPairs() []WeightedFieldValuePair {
	fvp := wc.SortedEraPairs()
	return wc.SortedWeightedPairs(fvp)
}

func (wc DbHeadwordCounts) SortedWeightedGenrePairs() []WeightedFieldValuePair {
	fvp := wc.SortedGenrePairs()
	return wc.SortedWeightedPairs(fvp)
}

func (wc DbHeadwordCounts) SortedWeightedCorpusPairs() []WeightedFieldValuePair {
	fvp := wc.SortedCorpusPairs()
	return wc.SortedWeightedPairs(fvp)
}
