//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package dating

import (
	"strings"
)

// ConvertOnePHIStringDate - attempt to convert one of 16k possible strings into a reasonable int...
func ConvertOnePHIStringDate(d string) int {
	if d == "" {
		return 9999
	}
	fp := TakeFingerprint(strings.TrimSpace(d))
	fp = pickandrunparser(fp)

	if fp.Calculated > 9999 {
		// something went terribly wrong in the parser: `unable to encode 317324324326 into binary format for int4 (OID 23)`
		fp.Calculated = 9999
	}
	return fp.Calculated
}

// hgdb=> SELECT COUNT(*) FROM works ;
// count
//--------
// 248061
//(1 row)

// hgdb=> SELECT COUNT(DISTINCT recorded_date) FROM works;
// count
//-------
// 22238
//(1 row)

// hgdb=> SELECT COUNT(*) FROM works where recorded_date != '' and recorded_date != '?' and recorded_date != 'date?' and converted_date = 9999;
// count
//-------
//  2058
//(1 row)

// hgdb=> SELECT distinct recorded_date FROM works where recorded_date not like '' and recorded_date NOT LIKE '?' and converted_date = 9999 order by recorded_date limit 100;
