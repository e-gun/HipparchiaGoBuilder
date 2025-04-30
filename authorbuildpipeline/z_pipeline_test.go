//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"testing"
)

func TestFindwordcount(t *testing.T) {
	l1 := structs.DbWorkline{
		WkUID:    "1",
		Stripped: "one two three four five",
	}
	l2 := structs.DbWorkline{
		WkUID:    "1",
		Stripped: "six seven   eight   nine",
	}
	ll := []structs.DbWorkline{l1, l2}
	result := findwordcount("1", ll)
	want := 9
	if result != want {
		t.Errorf("findwordcount)\ngot\n%d\nwant\n%d", result, want)
	}
}
