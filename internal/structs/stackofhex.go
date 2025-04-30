//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"fmt"
	"log"
	"strconv"
)

type HexStringStack struct {
	Items []string
}

// NewHexStack - the factory function; return a HexStringStack; not threadsafe
func NewHexStack(items []string) HexStringStack {
	return HexStringStack{
		Items: items,
	}
}

// Contents - return the HexStringStack contents
func (s *HexStringStack) Contents() []string {
	return s.Items
}

// Flush - empty the HexStringStack contents
func (s *HexStringStack) Flush() {
	s.Items = []string{}
}

// IsEmpty - are there HexStringStack contents?
func (s *HexStringStack) IsEmpty() bool {
	if len(s.Items) == 0 {
		return true
	} else {
		return false
	}
}

// Pop - get the top item from the stack
func (s *HexStringStack) Pop() string {
	i, e := s.PopWithError()
	if e != nil {
		// log.Fatal(fmt.Errorf("HexStringStack.Pop(): empty stack"))
		fmt.Println("HexStringStack.Pop(): empty stack - this is bad...")
	}
	return i
}

// NPop - get the top N item from the stack
func (s *HexStringStack) NPop(n int) string {
	var out string
	var e error
	var p string
	for i := 0; i < n; i++ {
		p, e = s.PopWithError()
		out += p
	}
	if e != nil {
		log.Fatal(fmt.Errorf("HexStringStack.NPop(): empty stack"))
	}
	return out
}

// PopBase16 - get the top item as an int64: 'ef' --> 239
func (s *HexStringStack) PopBase16() int64 {
	v, e := s.PopWithError()
	if e != nil {
		fmt.Println("HexStringStack.PopBase16(): empty stack - this is bad...")
		return int64(0)
	}
	i, err := strconv.ParseInt(v[0:2], 16, 64) // "b1 " --> "b1"
	if err != nil {
		fmt.Println("HexStringStack.PopBase16() failed to convert string to int64:", v)
		log.Fatal(err)
	}
	return i
}

// PopMaskedByte - get the top item and 'unmask' its (ascii-ready) value with an AND on its bytes
func (s *HexStringStack) PopMaskedByte() int {
	return int(s.PopBase16() & 0x7F)
}

// PopWithError - try to get the top item of the stack
func (s *HexStringStack) PopWithError() (string, error) {
	var i string
	var e error
	if len(s.Items) == 0 {
		return i, fmt.Errorf("empty stack")
	} else {
		i = s.Items[len(s.Items)-1]
		s.Items = s.Items[:len(s.Items)-1]
		return i, e
	}
}

func (s *HexStringStack) PrintContents() {
	for i, x := range s.Items {
		fmt.Printf("%d\t'%s'\n", i, x)
	}
}

func (s *HexStringStack) Push(item string) {
	s.Items = append(s.Items, item)
}

func (s *HexStringStack) Size() int {
	return len(s.Items)
}
