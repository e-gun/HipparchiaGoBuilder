//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package structs

import (
	"fmt"
	"log"
	"slices"
)

type ByteStack struct {
	Items []byte
}

// NewByteStack - the factory function; return a ByteStack; not threadsafe
func NewByteStack(items []byte) ByteStack {
	return ByteStack{
		Items: items,
	}
}

// Contents - return the ByteStack contents
func (s *ByteStack) Contents() []byte {
	return s.Items
}

// Flush - empty the ByteStack contents
func (s *ByteStack) Flush() {
	s.Items = []byte{}
}

// IsEmpty - are there ByteStack contents?
func (s *ByteStack) IsEmpty() bool {
	if len(s.Items) == 0 {
		return true
	}
	return false
}

// Pop - get the top item from the stack
func (s *ByteStack) Pop() byte {
	i, e := s.PopWithError()
	if e != nil {
		log.Fatal(fmt.Errorf("ByteStack.Pop(): empty stack"))
	}
	return i
}

// NPop - get the top N item from the stack
func (s *ByteStack) NPop(n int) []byte {
	var out []byte
	var e error
	var p byte
	for i := 0; i < n; i++ {
		p, e = s.PopWithError()
		out = append(out, p)
	}
	if e != nil {
		log.Fatal(fmt.Errorf("ByteStack.NPop(): empty stack"))
	}
	return out
}

// SNPop - get the top N item from the stack and return them as a string
func (s *ByteStack) SNPop(n int) string {
	b := s.NPop(n)
	return string(b)
}

// PopWithError - try to get the top item of the stack
func (s *ByteStack) PopWithError() (byte, error) {
	var i byte
	var e error
	if len(s.Items) == 0 {
		return i, fmt.Errorf("empty stack")
	}
	i = s.Items[len(s.Items)-1]
	s.Items = s.Items[:len(s.Items)-1]
	return i, e
}

func (s *ByteStack) PrintContents() {
	for i, x := range s.Items {
		fmt.Printf("%d\t'%v'\n", i, x)
	}
}

func (s *ByteStack) Push(item byte) {
	s.Items = append(s.Items, item)
}

func (s *ByteStack) Reverse() {
	slices.Reverse(s.Items)
}

func (s *ByteStack) Size() int {
	return len(s.Items)
}
