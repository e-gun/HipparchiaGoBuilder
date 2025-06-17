package generic

import (
	"fmt"
	"testing"
)

func TestLCGreek(t *testing.T) {
	tg := `a(/bruna nwxelou=s nwxelh/s nwxelw=s nwxelh/s o( o( o(/ o( o(/s o(/,te o(/te,o(/ste o(/te,o(/te`
	g := ConvertLCBetacode(tg)
	fmt.Println(g)
	// ἅβρυνα νωχελοῦϲ νωχελήϲ νωχελῶϲ νωχελήϲ ὁ ὁ ὅ ὁ ὅϲ ὅ,τε ὅτε,ὅϲτε ὅτε,ὅτε
	tg = `o(`
	g = ConvertLCBetacode(tg)
	fmt.Println(g)
	// oops: ο(
}
