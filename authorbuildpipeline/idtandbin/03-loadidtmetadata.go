//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package idtandbin

import (
	"bytes"
	"fmt"
	"os"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

// LoadAuthorIDTData - build partial DbAuthor and DbWork entries from IDT data
func LoadAuthorIDTData(idtbytes []byte) (structs.DbAuthor, []structs.DbWork) {
	// painfully walk through the bytes of the idt data
	// Homer (0012):
	// [1 6 133 0 0 239 128 176 176 177 178 255 16 0 16 38 49 72 111 109 101 114 117 115 38 32 69 112 105 99 46 2 3 67 0 0 239 129 176 176 177 255 16 1 5...]

	var thisauth structs.DbAuthor
	var foundworks []structs.DbWork

	leveldict := make(map[int]string)

	offset := -1
	for offset < len(idtbytes)-1 {
		offset++

		if idtbytes[offset] == 0 {
			continue
		} else if idtbytes[offset] == 1 || idtbytes[offset] == 2 {
			offset += 3
			offset++
			offset++

			if idtbytes[offset] == 239 {
				offset++
				// 128 176 176 177 178 255 16 0 16 38 49 72 111 109 101 114 117 115 38 32 69 112 105 99 46 2 3 67 0 0  ...
				level := int(idtbytes[offset] & 0x7f)
				// 128 -> 0
				offset++
				asciistring := getasciistring(idtbytes, offset)
				offset += len(asciistring)

				switch level {
				case 0:
					thisauth.UID = asciistring // but this will be incorrect for now: `0012` and not `gr0012`
					offset++
					if idtbytes[offset] == 16 && idtbytes[offset+1] == 0 {
						offset += 2
						authorname := getpascalstr(idtbytes, offset)
						thisauth.IDXname = authorname
						offset += len(authorname)
					} else {
						fmt.Printf("Author number apparently was not followed by author name: %d = %d\n", offset, idtbytes[offset])
						break
					}
				case 1:
					worknumber := asciistring
					offset++
					if idtbytes[offset] == 16 && idtbytes[offset+1] == 1 {
						offset += 2
						workname := getpascalstr(idtbytes, offset) // this can have special characters and needs cleaning
						offset += len(workname)
						offset++
						for idtbytes[offset] == 17 {
							depth, levellabel := findlabelsforlevels(idtbytes, offset)
							leveldict[depth] = levellabel
							offset += len(levellabel) + 3
						}
						offset--
						work := structs.DbWork{
							UID:   fmt.Sprintf("%sw%s", thisauth.UID, worknumber),
							Title: workname,
						}
						work = insetlevellabels(leveldict, work)
						foundworks = append(foundworks, work)
						//thisauth.Works = append(thisauth.Works, work)
						//thisauth.WorkDict[worknumber] = len(thisauth.Works) - 1
						leveldict = make(map[int]string)
					} else {
						fmt.Printf("LoadAuthorIDTData(): Work number apparently was not followed by work: %d = %d\n", offset, idtbytes[offset])
					}
				case 2:
					fmt.Println("LoadAuthorIDTData() made it to level 2: sub-works. This never happens?")
					os.Exit(1)
				}
			}
		} else if idtbytes[offset] == 3 {
			offset += 4
		} else if idtbytes[offset] == 10 {
			offset++
		} else if idtbytes[offset] == 11 || idtbytes[offset] == 13 {
			offset += 2
		}
	}

	return thisauth, foundworks
}

func getasciistring(filearray []byte, offset int) string {
	var buffer bytes.Buffer
	for _, b := range filearray[offset:] {
		if b == 255 {
			break
		}
		buffer.WriteByte(b & 0x7f)
	}
	return buffer.String()
}

func getpascalstr(filearray []byte, offset int) string {
	strlen := int(filearray[offset])
	offset++
	var buffer bytes.Buffer
	for _, b := range filearray[offset : offset+strlen] {
		buffer.WriteByte(b)
	}
	return buffer.String()
}

func findlabelsforlevels(filearray []byte, offset int) (int, string) {
	newByteCountOffset := 1
	depth := int(filearray[offset+newByteCountOffset])
	newByteCountOffset++
	levelLabel := getpascalstr(filearray, offset+newByteCountOffset)
	return depth, levelLabel
}

func insetlevellabels(ldict map[int]string, wk structs.DbWork) structs.DbWork {
	for level, label := range ldict {
		switch level {
		case 0:
			wk.LL0 = label
		case 1:
			wk.LL1 = label
		case 2:
			wk.LL2 = label
		case 3:
			wk.LL3 = label
		case 4:
			wk.LL4 = label
		case 5:
			wk.LL5 = label
		default:
			fmt.Printf("Unrecognized level label: %d = %s\n", level, label)
			os.Exit(1)
		}
	}
	return wk
}
