//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines1initial

import (
	"fmt"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// TLG5035 produces nybbler errors: Scholia In Platonem, Scholia in Platonem (scholia vetera)
// but the problem comes from above... fontshifts; if you turn ReplaceGreekFontMarkup() off,
// the nybble errors disappear; this is a useful lesson for problem that are broadcast here:
// there is a strong chance that they are not really caused by a logic problem here

// SEE BOTTOM OF FILE FOR EXAMPLE/QUIRKS/IRREGULARITIES

const (
	METADATATEMPLATE       = `<hb-metadata_%s value="%s" />`
	LEVELINCREMENTTEMPLATE = "\n<hb-incr_l_%d_by_1 />"
	LEVELSETTEMPLATE       = "<hb-set_l_%d_to_%s />"
	WORKINCREMENTTEMPLATE  = "\n<hb-increment_work_number_by_1 />"
	TMPTTC                 = `. <hb-nb-endofpage /> █ⓕⓔ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █ⓔⓕ █⑧⓪ █ⓑ③ █ⓑ⓪ █ⓑ② █ⓑ⑦ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ③ █ⓕⓕ █ⓐ⑧ █⑧ⓔ █⑨ⓑ █⑧② █ⓔ⑧ █⑧⑨ █⑧⑥ █ⓕ④ <hb-title>ΙΩΑΝΝΟΥ ΤΟΥ ΔΟΞΑΠΑΤΡΗ █⑧⑨ █⑧⑦ █ⓕ④ ΕΙϹ ΤΟ █⑧⑨ █⑧⑧ █ⓕ④ ΠΕΡΙ ΕΥΡΕϹΕΩϹ ΕΡΜΟΓΕΝΟΥϹ ΒΙΒΛΙΟΝ</hb-title> █⑧⑧ █⑧⑨ <hb-tabbedtext />`
)

var (
	findhex     = regexp.MustCompile(`(█[⓪①②③④⑤⑥⑦⑧⑨ⓐⓑⓒⓓⓔⓕ]{1,2}\s+)+`)
	hexreplacer = strings.NewReplacer(
		"⓪", "0",
		"①", "1",
		"②", "2",
		"③", "3",
		"④", "4",
		"⑤", "5",
		"⑥", "6",
		"⑦", "7",
		"⑧", "8",
		"⑨", "9",
		"ⓐ", "a",
		"ⓑ", "b",
		"ⓒ", "c",
		"ⓓ", "d",
		"ⓔ", "e",
		"ⓕ", "f",
	)
	metadatacategories = map[int]string{
		0:   "newauthor",
		1:   "newwork",
		2:   "workabbrev",
		3:   "authabbrev",
		97:  "region",
		98:  "city",
		99:  "notes",
		100: "date",
		101: "publicationinfo",
		102: "additionalpubinfo",
		103: "stillfurtherpubinfo",
		108: "provenance",
		114: "reprints",
		116: "unknownmetadata116",
		122: "documentnumber",
	}
	nybbleactionmapper = map[int]func(structs.HexStringStack) (string, structs.HexStringStack){
		8:  nyb08,
		9:  nyb09,
		10: nyb10,
		11: nyb11,
		12: nyb12,
		13: nyb13,
		14: nyb14,
		15: nyb15,
	}
)

func HexRunner(ttc string) string {
	//  Perhaps the most important single segment: decode the (undocumented) embedded hex values

	// 	First you find the hex runs.
	//	Then you send these to the citation builder to be read/decoded
	//	All the heavy lifting happens there

	// at the end of the process a set of level values might look like "[1 1 1 1 23 543]"
	// and those ones are suspect; they get turned into -1 later when you know the max # of levels in any given text

	// ttc = TMPTTC
	return findhex.ReplaceAllStringFunc(ttc, citationbuilder)
}

func highunicodetohex(highunicode string) string {
	return hexreplacer.Replace(highunicode)
}

func forceregexsafevariants(text string) string {
	// Example implementation, you should replace it with your actual implementation
	return text
}

func avoidregexsafevariants(text string) string {
	// Example implementation, you should replace it with your actual implementation
	return text
}

func citationbuilder(hexsequence string) string {
	// 	parse the sequence of bytes for instructions on how to build a citation
	//	"0xab 0x82 0xff 0x9f 0xe1 0xff" ==> "<set_level_2_to_383>\n<set_level_1_to_a>"
	//
	//	this is very fiddly since each instruction requires different actions upon the data
	//	that follows; furthermore 'level6' is its own world of metadata rather than being a
	//	simple hierarchy like 'level00 ==> verse#' and 'level01 ==> poem#' and 'level02 ==> book#'
	//
	//	everything happens byte-by-byte: that means that if you botch one set of instructions
	//	you will be left with the wrong byte sequence going forward and the next 'instruction'
	//	is likely to be a piece of the previous set of information; this makes debugging rough
	//	since garbage at location B can result from (unrecognized) troubles at location A.
	//
	//	this is perhaps the most important bit of the builder and I never could have coded it if
	//	P. J. Heslin had not released the source code to Diogenes:
	//
	//		https://community.dur.ac.uk/p.j.heslin/Software/Diogenes/
	//
	//	I shudder to think how long it took him to figure out the way to read the bytes and nybbles
	//	properly. Without his efforts Hipparchia would not have been possible.

	thehex := highunicodetohex(hexsequence)
	hexsegments := strings.Split(thehex, "█")
	slices.Reverse(hexsegments)
	hs := structs.NewHexStack(hexsegments)
	hs.Pop() // discard top which should be ""

	fullcitation := ""

	for !hs.IsEmpty() {
		instructions := hs.Pop()
		if instructions == "0 " {
			instructions = ""
			// end of file: ... █⓪ █⓪ █⓪ █⓪ █⓪
			// BUT: runs of zeroes in the middle of TLG3027 so you can't just flush...
			// hs.Flush()
			continue
		}

		textlevel, action := nybbler(instructions)
		// fmt.Println("t, a", textlevel, action)

		if textlevel < 6 {
			// this is very fussy and rather arbitrary: "making it work" in a byte-constrained environments
			// actions 0-6 are "levelincrement" (0) and "levelset" (1-7) (basic stuff)
			// actions 8+ all have different methods of decoding the data that follows and the nybbleactionmapper sends the data to its method
			if action == 0 {
				fullcitation += fmt.Sprintf(LEVELINCREMENTTEMPLATE, textlevel)
			} else if action > 0 && action < 8 {
				lsval := fmt.Sprintf("%d", action)
				fullcitation += fmt.Sprintf(LEVELSETTEMPLATE, textlevel, lsval)
			} else {
				citation, remaining := nybbleactionmapper[action](hs)
				hs.Items = remaining.Items
				fullcitation += fmt.Sprintf(LEVELSETTEMPLATE, textlevel, citation)
			}
		} else if textlevel == 6 {
			cit, leftover := levelsixparsing(action, nybbleactionmapper, fullcitation, hs)
			hs.Items = leftover.Items
			fullcitation = cit
		}
	}

	return fullcitation
}

// levelsixparsing - this is its own world; the papyri metadata is stored here; it is itself betacoded...
func levelsixparsing(action int, actionmapper map[int]func(structs.HexStringStack) (string, structs.HexStringStack), fullcitation string, hs structs.HexStringStack) (string, structs.HexStringStack) {
	// note the recursivity below: level six will itself use the nybbler logic that got you here
	metadata := make(map[string]string)
	category := 0
	if !hs.IsEmpty() {
		category = hs.PopMaskedByte()
	}

	var citation string
	if action == 0 {
		if category == 1 {
			fullcitation += WORKINCREMENTTEMPLATE
		}
		citation = ""
	} else if action > 0 && action < 8 {
		citation, hs = nyb15(hs)
	} else {
		citation, hs = actionmapper[action](hs)
	}

	if category != 100 {
		citation = forceregexsafevariants(citation)
	} else {
		citation = avoidregexsafevariants(citation)
	}

	citation = strings.ReplaceAll(citation, "`", "")
	metadata[metadatacategories[category]] = citation

	for key, value := range metadata {
		if len(value) > 0 {
			// note that you can read this with the naked eye in the hex:
			// 'b3 b0 b2 b7' is '3027' and 'b0 b0 b3' is '003'
			fullcitation += fmt.Sprintf(METADATATEMPLATE, key, cleanembeddedvalue(value))
		}
	}

	return fullcitation, hs
}

func cleanembeddedvalue(ev string) string {
	// you might have "date: III%3IIa"; so you need to pass through betacodecleanup
	// BUT it is not clear that the whole suite is needed (or even works); specifically there are many examples of
	// things like "211%3212p" and here you were supposed to have "211%3`212p" so that the correct %3 is found
	// instead you will look up %3212... Note that there are a smattering of other betacode possibilities in here
	// such as: Welles,RC.28%7 --> Welles,RC.28﹢ ; is it the case that only %3 and %7 are needed? it sure looks like it.

	// nevertheless, work likely remains: see IN0010 and:
	//an	 region: Att. · city: Athens · notes: stoich. · date: 5th-4th bcoeofkz!
	//an	 region: Att. · city: Athens · notes: stoich. · date: 4th-3rd bcoeofkz"
	//an	 region: Att. · city: Athens · notes: stoich. · date: 3rd bc?oeofkz&
	//an	 region: Att. · city: Athens · notes: stoich. · date: 3rd bc?oeofkz)
	//an	 region: Att. · city: Athens · notes: stoich.? · date: 3rd bc?oeofkz+
	// this is a problem in the hexrunner though, right? note that these lines can produce "bell" too

	ev = strings.ReplaceAll(ev, "%3", "／")
	ev = strings.ReplaceAll(ev, "%7", "﹢")

	ev = betacode.SimpleLatinSpanLATE(ev) // Cic. &3Leg.& 1.2 --> <hb-fs-l-italic>Leg.</hb-fs-l-italic> 1.2
	ev = lat.LatinCleanup(ev)

	// re-enable later if needed...
	// ev = betacode.BetaCodeCleanup(ev)
	return ev
}

// nybbler - take a character and split it into two chunks of info: 'textlevel' on left and 'action' on right
func nybbler(hexvalstring string) (int, int) {
	if hexvalstring == "" {
		global.BAD("nybbler(): empty hex string; THIS IS VERY BAD")
		// should probably panic...
		return -1, -1
	}

	hexval := getbase16(hexvalstring)
	textlevel := (hexval & 0x70) >> 4
	action := hexval & 0x0F

	// fmt.Printf("nybbler() lvl: action\t'%s': %d --> %d ; %d.\n", hexvalstring, hexval, int(textlevel), int(action))
	//nybbler() lvl: action	'b0': 176 --> 3 ; 0.
	//nybbler() lvl: action	'ff': 255 --> 7 ; 15.
	//nybbler() lvl: action	'ef': 239 --> 6 ; 15.
	//nybbler() lvl: action	'81': 129 --> 0 ; 1.
	//nybbler() lvl: action	'b0': 176 --> 3 ; 0.
	//nybbler() lvl: action	'b0': 176 --> 3 ; 0.

	return int(textlevel), int(action)
}

//
// NYBs - different ways of parsing the subsequent hex values
//
// NB: if you see errors here there is every chance that they are provoked by a parse problem *before* the
// function that yields garbage yields its garbage; be ready to read whole byte sequences from the top when
// debugging and not just "local" chunks
//
// if you see *any* of the warning messages emerge, then the build of this author is almost certainly badly broken
//

// nyb08 - 8 -> read 7 bits of next number [& int('7f', 16)]
func nyb08(hs structs.HexStringStack) (string, structs.HexStringStack) {
	var citation string
	if !hs.IsEmpty() {
		citation = fmt.Sprintf("%d", hs.PopMaskedByte())
	} else {
		citation = "[unk_nyb_08]"
	}
	return citation, hs
}

// nyb09 - 9 -> read a number and then a character
func nyb09(hs structs.HexStringStack) (string, structs.HexStringStack) {
	citation := ""
	if !hs.IsEmpty() {
		citation = fmt.Sprintf("%d", hs.PopMaskedByte())
	}
	if !hs.IsEmpty() {
		masked := hs.PopMaskedByte()
		if masked != 0xFF {
			citation += string(rune(masked))
		}
	}
	return citation, hs
}

// nyb10 - 10 -> read a number and then an ascii string
func nyb10(hs structs.HexStringStack) (string, structs.HexStringStack) {
	var citation string
	if !hs.IsEmpty() {
		citation = fmt.Sprintf("%d", hs.PopMaskedByte())
		stop := false
		for !hs.IsEmpty() && !stop {
			next := hs.PopMaskedByte()
			if next != 0xFF && next != 0x7F {
				// NB: "next != 0x7F" was added to fix a parsing bug; does that mean that "next != 0xFF" is wrong?
				citation += string(rune(next))
			} else {
				stop = true
			}
		}
	}

	return citation, hs
}

// nyb11 - 11 -> next two bytes are a 14 bit number
func nyb11(hs structs.HexStringStack) (string, structs.HexStringStack) {
	var citation string
	if !hs.IsEmpty() {
		firstbyte := hs.PopMaskedByte()
		secondbyte := hs.PopMaskedByte()
		citation = strconv.Itoa((firstbyte << 7) + secondbyte)
	} else {
		citation = "[unk_nyb_11]"
		// why is the next line there? strongly suspect that you only reach this when there is a parse problem...
		fmt.Println("nyb11 is flushing its HexStringStack: ", hs.Contents())
		hs.Flush() // doing this, but do not really know why...
	}
	return citation, hs
}

// nyb12 - 12 -> a 2-byte number, then a character
func nyb12(hs structs.HexStringStack) (string, structs.HexStringStack) {
	var citation string
	if !hs.IsEmpty() {
		firstbyte := hs.PopMaskedByte()
		if !hs.IsEmpty() {
			secondbyte := hs.PopMaskedByte()
			citation = strconv.Itoa((firstbyte << 7) + secondbyte)
			if !hs.IsEmpty() {
				citation += string(rune(hs.PopMaskedByte()))
			} else {
				citation += "[unk_nyb_12c]"
			}
		} else {
			citation += "[unk_nyb_12b]"
		}
	} else {
		citation += "[unk_nyb_12a]"
	}
	return citation, hs
}

// nyb13 - 13 -> a 2-byte number, then a string
func nyb13(hs structs.HexStringStack) (string, structs.HexStringStack) {
	var citation string
	if !hs.IsEmpty() {
		firstbyte := hs.PopMaskedByte()
		if !hs.IsEmpty() {
			secondbyte := hs.PopMaskedByte()
			citation = strconv.Itoa((firstbyte << 7) + secondbyte)
			if !hs.IsEmpty() {
				stringPart, remaining := nyb15(hs)
				hs = remaining
				citation += stringPart
			} else {
				citation += "[unk_nyb_13b]"
			}
		} else {
			citation += "[unk_nyb_13b]"
		}
	} else {
		citation += "[unk_nyb_13a]"
	}
	return citation, hs
}

// nyb14 - 14 -> append a char to the counter number
func nyb14(hs structs.HexStringStack) (string, structs.HexStringStack) {
	var citation string
	if !hs.IsEmpty() {
		citation += string(rune(hs.PopMaskedByte()))
	} else {
		citation += "[unk_nyb_14]"
	}
	return citation, hs
}

// nyb15 - 15 -> an ascii string follows
func nyb15(hs structs.HexStringStack) (string, structs.HexStringStack) {
	var citation string
	stop := false
	for !hs.IsEmpty() && !stop {
		p16 := hs.PopBase16()
		if p16 != 0xFF {
			citation += string(rune(int(p16) & 0x7F))
		} else {
			stop = true
		}
	}
	return citation, hs
}

func getbase16(s string) int64 {
	i, err := strconv.ParseInt(s[0:2], 16, 64) // "b1 " --> "b1"
	if err != nil {
		global.BAD("getbase16() failed to convert string to int64")
		log.Fatal(err)
	}
	return i
}

// EXAMPLE:
// deep inside gr2102 as of end of GreekCleanup():
// █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ① █ⓑ⑥ █ⓕⓕ █⑨① <hb-tabbedtext />Ἔχ̣[ο]μ̣ε̣ν̣ περὶ θεοῦ διάλημψιν καὶ ἀπὸ τ̣ῆϲ γραφῆϲ καὶ
// █⑧⓪ τῆϲ κ[ο]ινῆϲ ἐννοίαϲ ὅτι ἄτρεπτόϲ ἐϲτιν, │ ὅ[τι ἀ]ναλλοί-
// █⑧② ωτόϲ ἐϲτιν· ὁ γὰρ ὅλωϲ μὴ ὑποκείμενοϲ ποιότη̣[τι οὐ τρέ-
// █⑧⓪ πετ]αι, οὐκ ἀλλοιοῦται· οὐδὲν │ γ̣ὰ̣ρ̣ ἕτερ̣ό̣ν̣ ἐϲτιν ἀλλοί-
// █⑧③ ωϲιϲ ἢ κατὰ ποιὸν μεταβολή. οὐ [πᾶϲα μεταβο]λ̣ὴ ἀλλοίωϲίϲ
// █⑧⓪ ἐϲτιν, │ ἀλλ＇ ἡ κ̣[ατὰ] π̣οιότητα. εἰ[ϲ]ίν γ̣ε̣ κα[ὶ ἄλλαι]
// █⑧④ μεταβολαί, ἐπε[ὶ καὶ κινήϲειϲ εἰϲ]ὶν ἄλλαι. τὸ γ̣ι̣νόμενον │
// █⑧⓪ μεταβάλλει, ἀλλ＇ οὔκ ἐϲτιν ἀ̣λ̣λο[ί]ω̣[ϲ]ι̣ϲ ἡ κί[ν]η̣[ϲ]ιϲ
// █⑧⓪ αὕτη. τὸ αὐξό̣μ̣[ενον μεταβάλλ]ει, ἀλλ＇ οὔκ ἐϲτιν ἀλλοί│ω-
// █⑧⑥ ϲιϲ αὕτη· προϲθήκη γὰρ καὶ αὔξηϲιϲ̣ π̣[ο]ϲ̣[ο]ῦ̣ ἐϲτιν ἡ τοι-
// █⑧⓪ αύτη κ[ί]ν̣η̣[ϲ]ιϲ. ὅ̣[τ]αν δὲ ἐκ φα̣ύ̣λ̣ο̣υ ϲ[πο]υ̣│δα̣ῖ̣οϲ̣ ἢ̣ ἐ̣κ̣
// █⑧⑦ ϲπουδαίου φαῦλοϲ γένητ̣[α]ί τ̣ι̣ϲ̣, ἠ̣λλοίωται κατ̣[ὰ] τ̣[ὴν
// █⑧⓪ ποιότητα], ὡϲ αὖ ὅτε ἐκ νοϲο̣ῦν│τοϲ εἰϲ ὑγείαν ἔλθῃ καὶ
// █⑧⑧ █⑧⑧ ἔνπαλιν.
// █⑧⑧ █⑧⑧ <hb-tabbedtext />ἀ[λ]λ̣ὰ̣ τ[ὰ]ϲ̣ [λ]έ̣[ξ]ειϲ πρὸϲ τὴν ἔν[νοιαν τοῦ πρ]ά̣γ-
// █⑧⓪ ματοϲ, π̣ερὶ οὗ │ λέγονται, δ̣ι̣[ανοού]μεθα. ὁ θεὸϲ οὐκ ἐκ

//<hb-metadata_newwork value="016" /><hb-set_l_1_to_1 /><hb-tabbedtext />Ἔχ̣[ο]μ̣ε̣ν̣ περὶ θεοῦ διάλημψιν καὶ ἀπὸ τ̣ῆϲ γραφῆϲ καὶ
//<hb-incr_l_0_by_1 />τῆϲ κ[ο]ινῆϲ ἐννοίαϲ ὅτι ἄτρεπτόϲ ἐϲτιν, │ ὅ[τι ἀ]ναλλοί-
//<hb-set_l_0_to_2 />ωτόϲ ἐϲτιν· ὁ γὰρ ὅλωϲ μὴ ὑποκείμενοϲ ποιότη̣[τι οὐ τρέ-
//<hb-incr_l_0_by_1 />πετ]αι, οὐκ ἀλλοιοῦται· οὐδὲν │ γ̣ὰ̣ρ̣ ἕτερ̣ό̣ν̣ ἐϲτιν ἀλλοί-
//<hb-set_l_0_to_3 />ωϲιϲ ἢ κατὰ ποιὸν μεταβολή. οὐ [πᾶϲα μεταβο]λ̣ὴ ἀλλοίωϲίϲ
//<hb-incr_l_0_by_1 />ἐϲτιν, │ ἀλλ’ ἡ κ̣[ατὰ] π̣οιότητα. εἰ[ϲ]ίν γ̣ε̣ κα[ὶ ἄλλαι]
//<hb-set_l_0_to_4 />μεταβολαί, ἐπε[ὶ καὶ κινήϲειϲ εἰϲ]ὶν ἄλλαι. τὸ γ̣ι̣νόμενον │
//<hb-incr_l_0_by_1 />μεταβάλλει, ἀλλ’ οὔκ ἐϲτιν ἀ̣λ̣λο[ί]ω̣[ϲ]ι̣ϲ ἡ κί[ν]η̣[ϲ]ιϲ
//<hb-incr_l_0_by_1 />αὕτη. τὸ αὐξό̣μ̣[ενον μεταβάλλ]ει, ἀλλ’ οὔκ ἐϲτιν ἀλλοί│ω-
//<hb-set_l_0_to_6 />ϲιϲ αὕτη· προϲθήκη γὰρ καὶ αὔξηϲιϲ̣ π̣[ο]ϲ̣[ο]ῦ̣ ἐϲτιν ἡ τοι-
//<hb-incr_l_0_by_1 />αύτη κ[ί]ν̣η̣[ϲ]ιϲ. ὅ̣[τ]αν δὲ ἐκ φα̣ύ̣λ̣ο̣υ ϲ[πο]υ̣│δα̣ῖ̣οϲ̣ ἢ̣ ἐ̣κ̣
//<hb-set_l_0_to_7 />ϲπουδαίου φαῦλοϲ γένητ̣[α]ί τ̣ι̣ϲ̣, ἠ̣λλοίωται κατ̣[ὰ] τ̣[ὴν
//<hb-incr_l_0_by_1 />ποιότητα], ὡϲ αὖ ὅτε ἐκ νοϲο̣ῦν│τοϲ εἰϲ ὑγείαν ἔλθῃ καὶ
//<hb-set_l_0_to_8 />ἔνπαλιν.
//<hb-set_l_0_to_8 /><hb-tabbedtext />ἀ[λ]λ̣ὰ̣ τ[ὰ]ϲ̣ [λ]έ̣[ξ]ειϲ πρὸϲ τὴν ἔν[νοιαν τοῦ πρ]ά̣γ-
//<hb-incr_l_0_by_1 />ματοϲ, π̣ερὶ οὗ │ λέγονται, δ̣ι̣[ανοού]μεθα. ὁ θεὸϲ οὐκ ἐκ

// the key thing to note in the above example is `█⑧⑦` which here means `<hb-set_l_0_to_7 />` (as you would expect...)

// Cicero:
// █ⓕⓔ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █⓪ █ⓔⓕ █⑧⓪ █ⓑ⓪ █ⓑ④ █ⓑ⑦ █ⓑ④ █ⓕⓕ █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ② █ⓕⓕ █ⓔⓕ █⑧② █ⓓ③ █ⓓ② █ⓔⓕ █ⓕ③ █ⓔ③ █ⓕⓕ █ⓔⓕ █⑧③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓕⓕ █⑨ⓑ █⑧① █⑧⑦ █⑧④ █ⓔⓕ █ⓔ③ █ⓕⓕ █ⓔ② █ⓕⓐ se neminem putet,
// <hb-metadata_newauthor value="0474" /><hb-metadata_newwork value="002" /><hb-metadata_workabbrev value="SRosc" /><hb-metadata_authabbrev value="Cic" /><hb-set_l_1_to_135 /><hb-set_l_0_to_4 />se neminem putet, ut se solum beatum, solum potentem

// decomposed:
// █ⓔⓕ █⑧⓪ █ⓑ⓪ █ⓑ④ █ⓑ⑦ █ⓑ④ █ⓕⓕ         --> (80) <hb-metadata_newauthor value="0474" />
// █ⓔⓕ █⑧① █ⓑ⓪ █ⓑ⓪ █ⓑ② █ⓕⓕ              --> (81) <hb-metadata_newwork value="002" />
// █ⓔⓕ █⑧② █ⓓ③ █ⓓ② █ⓔⓕ █ⓕ③ █ⓔ③ █ⓕⓕ    --> (82) <hb-metadata_workabbrev value="SRosc" />
// █ⓔⓕ █⑧③ █ⓒ③ █ⓔ⑨ █ⓔ③ █ⓕⓕ               --> (83) <hb-metadata_authabbrev value="Cic" />
// █⑨ⓑ █⑧① █⑧⑦                              -->  "l1 a11" (nyb11: next two bytes are a 14 bit number) <hb-set_l_1_to_135 />
// █⑧④                                         -->  "l0" <hb-set_l_0_to_4 />
// █ⓔⓕ █ⓔ③ █ⓕⓕ                               --> "l6 a15" + █ⓔ③ --> nyb15(hs) --> `c`
// █ⓔ② █ⓕⓐ                                    --> "l6 a2" + █ⓕⓐ --> nyb15(hs) --> `z`
