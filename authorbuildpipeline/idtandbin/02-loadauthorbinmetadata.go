//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"slices"
)

// BINDataExtraction - parse bindata and associate categories with author numbers
func BINDataExtraction(bindata []byte) map[string][]string {
	// internal bindata format
	// labels: H-pack
	//	[1 2 3 4]
	//	read #4:
	//		0 --> 'last'
	//		N --> len of python string
	//
	// authors: h-pack
	//	[1 2] --> (& hex 7fff) --> author number in 4-place decimal
	//	[3 4] --> works
	//	read #4
	//		test-for-last, etc

	// aiming for: 'Acta': ['0031', '0300', '0317', '8577', '1025', '1057', '0389', '2038', '9352', '18433', '2948', '9099', '1281'], 'Alchemica': ['1016', '1304', '1306', ...], ... }

	stack := structs.NewByteStack(bindata)
	stack.Reverse()

	quant := 999
	var labels []string
	var datastartpoint [][]byte

	// walk through the file byte by byte

	// (a) first find the categories: labels are things like "Lyrica"; datastartpoint is where to find the data: '[0 0 34 251]'

	stack.NPop(4) // discard
	for quant != 0 {
		head := stack.NPop(4) // where the category data starts
		rv := stack.Pop()     // how many letters in the label for the category
		quant = int(rv)
		labels = append(labels, stack.SNPop(quant)) // pop that many letters and store them as the label
		// example:
		// 0 0 0 5 65 99 116 97 32
		// 0 0 0 5: grab five letters
		// 65 99 116 97 32: the letters are 'Acta '
		datastartpoint = append(datastartpoint, head)
	}

	labels = labels[:len(labels)-1]                         // the last item is ''
	datastartpoint = datastartpoint[:len(datastartpoint)-1] // the last item is '[0 0 0 0]'

	//for i, label := range labels {
	//	fmt.Printf("'%s'\t%v\n", label, datastartpoint[i])
	//}

	//'Acta '	[0 0 0 0]
	//'Alchemica'	[0 0 0 37]
	//'Anthologiae'	[0 0 0 146]
	//'Apocalypses'	[0 0 0 172]
	// ...

	// (b) trim the pile of zeros in the middle of the stack
	stop := false
	for !stop {
		next := stack.Pop()
		if int(next) != 0 {
			stack.Push(next)
			stop = true
		}
	}

	// (c) find where the data relative to each category is located and aggregate this information

	// reverse again
	stack.Reverse()
	remainingbytes := stack.Contents()

	// fmt.Println("Remaining bytes:", len(remainingbytes))

	bytebundles := make(map[string][]byte)

	for i, label := range labels {
		// category data start locations
		// datastartpoint should look like: '[0, 0, 0, 0], [0, 0, 0, 37], [0, 0, 0, 146], [0, 0, 0, 172], [0, 0, 0, 221], [0, 0, 1, 57], ...'

		l1 := decodedatalocation(datastartpoint[i])
		var l2 int
		if i+1 < len(datastartpoint) {
			l2 = decodedatalocation(datastartpoint[i+1])
		} else {
			l2 = stack.Size()
		}

		// fmt.Printf("%s\tl1: %d\tl2: %d\n", label, l1, l2)
		// Acta 	l1: 0	l2: 37
		// Alchemica	l1: 37	l2: 146
		// Anthologiae	l1: 146	l2: 172
		// ...

		bytebundles[label] = remainingbytes[l1:l2]
		// fmt.Printf("remainingbytes[%d:%d]\t%v\n", l1, l2, remainingbytes[l1:l2])
		// remainingbytes[0:37]	[31 5 129 44 1 129 61 1 33 129 132 1 33 4 33 129 133 1 135 246 1 36 136 200 1 139 132 1 35 139 133 1 0 0 0 0 131]
	}

	sortedlabels := StringMapKeysIntoSortedSlice(bytebundles) // for debugging purposes we want repeatable runs

	// (d) unpack the data in the bytebundles and associate those lists of authors with the various headings
	headingsandauthors := make(map[string][]string)
	for _, lab := range sortedlabels {
		// clean the labels while we are at it: "Philosophici%3%19ae" might be in there
		headingsandauthors[betacode.ReplaceSingletons(betacode.ReplacePercentSigns(lab))] = findauthorsinbundles(bytebundles[lab])
	}

	//for _, lab := range labels {
	//	fmt.Printf("%v:\n", lab)
	//	fmt.Printf("\t%v\n", headingsandauthors[lab])
	//}

	// Tactici:
	//	[0058 0546 0556 0648 3075 3181]

	// note that some surprises return like '16416' and '16419' among the Tragica; but this is true in the python too

	return headingsandauthors
}

func decodedatalocation(bb []byte) int {
	// looks like '[0 0 33 75]'
	calc := (256 * int(bb[2])) + int(bb[3])
	return calc
}

func findauthorsinbundles(bb []byte) []string {
	// Acta
	// [0, 0, 0, 0, 1, 133, 139, 35, 1, 132, 139, 1, 200, 136, 36, 1, 246, 135, 1, 133, 129, 33, 4, 33, 1, 132, 129, 33, 1, 61, 129, 1, 44, 129, 5, 31, 128]
	// ['0031', '0300', '0317', '8577', '1025', '1057', '0389', '2038', '9352', '18433', '2948', '9099', '1281']

	// Anthologiae
	// [0, 0, 0, 0, 1, 139, 155, 1, 88, 155, 1, 51, 144, 4, 31, 138, 1, 100, 137, 75, 64, 248, 135, 1, 245, 135]
	// ['2037', '2040', '19337', '0394', '1168', '0411', '0411', '2817']

	bundle := structs.NewByteStack(bb)
	bundle.Reverse()

	var authors []string
	var nstr string
	for bundle.Size() > 1 {
		var auth string
		bynum := bundle.NPop(2)
		aunum := ((int(bynum[0]) * 256) + int(bynum[1])) & 32767 // int('7fff', 16)
		// need 5 to turn into '0005'
		nstr = fmt.Sprintf("%04d", aunum)
		if nstr == "0000" {
			break
		}
		nextbyte := int(bundle.Pop())
		if (nextbyte&128 != 0) || (nextbyte == 0) {
			auth = nstr
			authors = append(authors, auth)
			bundle.Push(byte(nextbyte))
		} else {
			auth = nstr
			authors = append(authors, auth)
			// bundle.Push(byte(nextbyte))

			// but does this sub-information ever really appear?
			// fmt.Printf("need to find works for %s too\n", nstr)
			// bundle = NewByteStack(findworkbundles(bundle.Contents()))
		}
	}
	return authors
}

func findworkbundles(bb []byte) []byte {
	// this is non-functional in the python... no works are ever found
	bundle := structs.NewByteStack(bb)
	nextbyte := int(bundle.Pop())
	for nextbyte&128 == 0 && nextbyte != 0 {
		nextbyte = int(bundle.Pop())
	}
	return bundle.Contents()
}

// StringMapKeysIntoSortedSlice - convert map[string]T to []string
func StringMapKeysIntoSortedSlice[T any](mp map[string]T) []string {
	sl := make([]string, len(mp))
	i := 0
	for k := range mp {
		sl[i] = k
		i += 1
	}
	slices.Sort(sl)
	return sl
}
