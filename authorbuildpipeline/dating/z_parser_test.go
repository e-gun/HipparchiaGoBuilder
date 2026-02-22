//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dating

import "testing"

func TestTwoarabicsimple(t *testing.T) {
	d := `1011／1012 ac`
	fp := TakeFingerprint(d)
	fp = twoarabicsimple(fp)
	fp.PrintWithCalc()
}

func TestOnearabiccentury(t *testing.T) {
	d := `5th ad`
	fp := TakeFingerprint(d)
	fp = onearabiccentury(fp)
	fp.PrintWithCalc()
}

func TestOneromancentury(t *testing.T) {
	d := `mid-IV ac?`
	fp := TakeFingerprint(d)
	fp = oneromancentury(fp)
	fp.PrintWithCalc()
}

func TestSlashdated(t *testing.T) {
	d := `c.63／2-51／0 bc`
	fp := TakeFingerprint(d)
	fp = slashdated(fp)
	fp.PrintWithCalc()
}

func TestEitherordate(t *testing.T) {
	d := `618 or 633 ac`
	fp := TakeFingerprint(d)
	fp = eitherordate(fp)
	fp.PrintWithCalc()
}

func TestMultislashdate(t *testing.T) {
	d := `109／108／106／105 BC`
	fp := TakeFingerprint(d)
	fp = multislashdate(fp)
	fp.PrintWithCalc()
	d = `AD 652／3／667／8'`
	fp = TakeFingerprint(d)
	fp = multislashdate(fp)
	fp.PrintWithCalc()
}

func TestMulticommadate(t *testing.T) {
	d := `114, 116, ﹠ 156 ac`
	fp := TakeFingerprint(d)
	fp = multicommadate(fp)
	fp.PrintWithCalc()
}

func TestAndsigndate(t *testing.T) {
	d := `1299 ﹠ 1344 ac`
	fp := TakeFingerprint(d)
	fp = andsigndate(fp)
	fp.PrintWithCalc()
}

func TestTwoarabiccomplex(t *testing.T) {
	d := `30 BC-AD 68`
	fp := TakeFingerprint(d)
	fp = twoarabiccomplex(fp)
	fp.PrintWithCalc()
	d = `100 bc-100 ac`
	fp = TakeFingerprint(d)
	fp = twoarabiccomplex(fp)
	fp.PrintWithCalc()
}

func TestPickandrunparser(t *testing.T) {
	d := `end 11th-beg. 12th ac`
	// d := `c II／IIIp`
	fp := TakeFingerprint(d)
	fp = pickandrunparser(fp)
	fp.Print()
	fp.PrintWithCalc()
}
