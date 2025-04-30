//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package latetidyups

import "testing"

func TestPurgeHybrid(t *testing.T) {
	ttc := `domno piissimo αugusto λeone anno et ξonstantino s(anctae) ῥomanae λ. Pomponius ξ[apitoli]o or ξ(apitolio)`
	result := PurgeHybrid(ttc)
	want := "domno piissimo Augusto Leone anno et Constantino s(anctae) (Romanae L. Pomponius C[apitoli]o or C(apitolio)"
	if result != want {
		t.Errorf("PurgeHybrid(%q)\ngot\n%q\nwant\n%q", ttc, result, want)
	}
}
