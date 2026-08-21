//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/betacode"
	"github.com/e-gun/HipparchiaGoBuilder/authorbuildpipeline/authorprep/lat"
	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	setter = regexp.MustCompile(`<hb-set_l_(\d)_to_(.*?)\s/>`)
	adder  = regexp.MustCompile(`<hb-incr_l_(\d)_by_1\s`)
	wnv    = regexp.MustCompile(`<hb-metadata_newwork value="(\d{1,3})"`)
	anv    = regexp.MustCompile(`<hb-metadata_newauthor value="(\d{1,4})" />`)
)

// AssignCitationValues - one of the key functions; build []DbWorkline with the right values in each line: `3.4.2`, etc.
func AssignCitationValues(wklinesasstring string) []structs.DbWorkline {
	// fmt.Println("AssignCitationValues")
	// 	will use decoded hex commands to build a citation value for every line in the text file
	//	can produce a formatted line+citation, but really priming us for the move to the dbc

	// inscription level 5 increments are "weird"; will set publication info (vs "document number")
	// <hb-set_l_5_to_3:680Ip1oaLocris, W.obPhyscensesodmid-II bckz/oe />

	levelmapper := map[int]int{
		0: 1,
		1: 1,
		2: 1,
		3: 1,
		4: 1,
		5: 1,
	}

	levelstrmapper := map[int]string{
		// be careful about re returning '1' and not 1
		0: "",
		1: "",
		2: "",
		3: "",
		4: "",
		5: "",
	}

	individuallines := strings.Split(wklinesasstring, "\n")

	dbr := make([]structs.DbWorkline, len(individuallines))
	auth := probeforauthor(individuallines)
	work := probefirstlineforwork(individuallines[0])
	overridelevel := -1
	overridevalue := ""

	for i, line := range individuallines {
		// needed for inscriptions
		pubnotes := ""

		// ToValidUTF8 is required
		// otherwise pgx will say: "ERROR: invalid byte sequence for encoding "UTF8": 0xe2 0x20 0xce (SQLSTATE 22021)"
		line = strings.ToValidUTF8(line, "")

		// Need to strip all 0x00 (null byte) from the lines
		// otherwise pg insert problem: gr0363 ERROR: invalid byte sequence for encoding "UTF8": 0x00 (SQLSTATE 22021)
		line = strings.ReplaceAll(line, string([]byte{0x00}), "")

		if gotwork := wnv.FindStringSubmatch(line); gotwork != nil {
			work = gotwork[1]
			for l := 0; l <= 5; l++ {
				levelmapper[l] = 1
				levelstrmapper[l] = ""
			}
		}

		if gotsetting := setter.FindStringSubmatch(line); gotsetting != nil {
			// Euripides (0006) has <hmu_set_level_0_to_post 961 /> after πῶς οὖν ἔτ’ ἂν θνήισκοιμ’ ἂν ἐνδίκως, πόσι,
			// 'post 961' becomes a problem: you need to add one to 961, but you will fail 'str(int(setting)'
			// slicing at the whitespace will fix this (sort of)
			// but then you get a new problem: UPZ (DDP0155) and its new documents '<hmu_set_level_5_to_2 rp />'
			// the not so pretty solution of the hour is to build a quasi-condition that is seldom met
			// it is almost never true that the split will yield anything other than the original item
			// it also is not clear how many other similar cases are out there: 'after 1001', etc.
			level, _ := strconv.Atoi(gotsetting[1])
			setting := gotsetting[2]

			overridelevel = level
			overridevalue = ""
			val, err := strconv.Atoi(strings.Split(setting, "post ")[len(strings.Split(setting, "post "))-1])
			if err != nil {
				// the error: strconv.Atoi: parsing "t": invalid syntax
				overridevalue = setting
			}

			levelmapper[level] = val
			levelstrmapper[level] = setting

			// you really want to avoid having anything other than "" in the levelstrmapper
			// if a "140" gets lodged in it, you will never be able to increment 140 by 1 (because "140" is not 140...)
			// block reinitialization is where you see the problem
			// <hb-metadata_newauthor value="0012" />
			// <hb-metadata_newwork value="001" /><hb-metadata_workabbrev value="Il" />
			// <hb-set_l_1_to_1 />
			// <hb-set_l_0_to_140 />

			numericalsetting, err2 := strconv.Atoi(setting)
			if err2 == nil {
				// so setting is in fact a number: we should drop the string and just modify levelmapper directly
				levelstrmapper[level] = ""
				levelmapper[level] = numericalsetting
			}

			if level > 0 {
				for l := 0; l < level; l++ {
					levelmapper[l] = 1
					levelstrmapper[l] = ""
				}
			}

			if level == 5 {
				if _, err = strconv.Atoi(getnumberorstring(5, levelmapper, levelstrmapper)); err != nil {
					// so you did not have a number...
					// <hb-set_l_5_to_3:680Ip1oaLocris, W.obPhyscensesodmid-II bckz/oe />
					// pubnotes = "NewDocument: " + lat.ConvertLatinDiacriticals(betacode.BetaCodeCleanup(setting))
					pubnotes = parseinscrl5setvalue(setting)
					levelmapper[level] = 1
					levelstrmapper[level] = ""
					overridevalue = ""
				}
			}
		}

		if gotincrem := adder.FindStringSubmatch(line); gotincrem != nil {
			// fmt.Println(levelmapper)
			// fmt.Println(levelstrmapper)

			level, _ := strconv.Atoi(gotincrem[1])
			setting := 1
			levelmapper[level] = levelmapper[level] + setting

			// trick/trap inside some authors
			// aristotle: <hb-set_l_1_to_501a /> --> ... --> <hb-incr_l_1_by_1 />  (i.e. 501a+1 = 501b)
			levelstrmapper[level] = incrementastring(levelstrmapper[level])

			// if you do not do the next Aristotle will work but Plato will not...
			overridelevel = level
			overridevalue = ""

			// clear/reset the levels below this one
			if level > 0 {
				for l := 0; l < level; l++ {
					levelmapper[l] = 1
					levelstrmapper[l] = ""
				}
			}
		}

		dbr[i] = structs.DbWorkline{
			WkUID:       auth + global.AUTHWORKSEPARATOR + work,
			TbIndex:     i,
			MarkedUp:    line,
			Lvl5Value:   getnumberorstring(5, levelmapper, levelstrmapper),
			Lvl4Value:   getnumberorstring(4, levelmapper, levelstrmapper),
			Lvl3Value:   getnumberorstring(3, levelmapper, levelstrmapper),
			Lvl2Value:   getnumberorstring(2, levelmapper, levelstrmapper),
			Lvl1Value:   getnumberorstring(1, levelmapper, levelstrmapper),
			Lvl0Value:   getnumberorstring(0, levelmapper, levelstrmapper),
			Annotations: pubnotes,
		}
		if overridevalue != "" {
			dbr[i] = applyvalueoverride(dbr[i], overridelevel, overridevalue)
		}
	}

	return dbr
}

// incrementastring - 501a+1 = 501b
func incrementastring(toincrem string) string {
	if toincrem == "" {
		return ""
	}

	_, skipifdigit := strconv.Atoi(toincrem) // will a digit in fact really arrive here?
	if skipifdigit == nil {
		return ""
	}
	rn := []rune(toincrem)
	rn[len(rn)-1] = rn[len(rn)-1] + 1
	return string(rn)
}

func getnumberorstring(level int, levelmapper map[int]int, levelstrmapper map[int]string) string {
	if levelstrmapper[level] != "" {
		return levelstrmapper[level]
	}
	return strconv.Itoa(levelmapper[level])
}

func probeforauthor(individuallines []string) string {
	var author string
	for _, line := range individuallines {
		if gotauth := anv.FindStringSubmatch(line); gotauth != nil {
			author = gotauth[1]
			break
		}
	}
	return author
}

func probefirstlineforwork(line string) string {
	var work string
	if gotwork := wnv.FindStringSubmatch(line); gotwork != nil {
		work = gotwork[1]
	}
	return work
}

func applyvalueoverride(wl structs.DbWorkline, level int, val string) structs.DbWorkline {
	switch level {
	case 0:
		wl.Lvl0Value = val
	case 1:
		wl.Lvl1Value = val
	case 2:
		wl.Lvl2Value = val
	case 3:
		wl.Lvl3Value = val
	case 4:
		wl.Lvl4Value = val
	case 5:
		wl.Lvl5Value = val
	default:
		fmt.Println("applyvalueoverride() sent an invalid value for level:", level)
	}
	return wl
}

var (
	purgedels = regexp.MustCompile(`\x7F.*?`)
)

func parseinscrl5setvalue(s string) string {
	const (
		PUB1   = `<hb-metadata_publicationinfo value="%s" />`
		REGION = `<hb-metadata_region value="%s" />`
		CITY   = `<hb-metadata_city value="%s" />`
		DATE   = `<hb-metadata_date value="%s" />`
		PUB2   = `<hb-metadata_extrapublicationinfo value="%s" />`
		PUB3   = `<hb-metadata_publicationinfo value="(?) %s" />`
	)
	// samples:

	// 98:201,395oaMakedoniaobSveti Nikolaocmdg／232poeSpom 77.57,59／60kz
	// 98:172,359oaMakedoniaobOrehovacocoddate?oecf Spomenik 98.343kzy
	// 71:167,437oaMakedoniaobPeštaniocodIIpoeSEG 2.434[cf Les Villes 289,58]kz
	// 98:180,380p1oaMakedoniaobSenokosocodaet RomoeSpom 71.180,476kz

	// 0: publication
	// oa: region
	// ob: city
	// oc: date
	// od: date
	// oe: additional publication info

	values := strings.Split(s, "\u007F")
	//98:203,406oaMakedoniaobPreotocodaet Romoekz
	// into:
	//98:203,406
	//oaMakedonia
	//obPreot
	//oc
	//odaet Rom
	//oe
	//kz

	if len(values) < 6 || len(values) > 8 {
		// fmt.Println("parseinscrl5setvalue error", s)
		// samples:
		// 16:399p1oaMakedoniaobBudur-C%148iflikocodaet RomoeSEG 13.405hz
		// 857/858
		return fmt.Sprintf(PUB3, s)
	}
	if len(values) == 6 {
		// the next will break the entry where there is a c/d problem, but...
		values = append(values, "")
		values[5] = values[4]
		values[4] = "od"
	}

	if len(values) == 8 {
		// there are 8 in:
		// 98:180,380p1oaMakedoniaobSenokosocodaet RomoeSpom 71.180,476kz
		var newvals []string
		newvals = append(newvals, values[0]+values[1])
		for i, v := range values {
			if i < 2 {
				continue
			}
			newvals = append(newvals, v)
		}
		values = newvals
	}

	for i := range values {
		if i == 0 {
			continue
		}
		trimrunes := 2
		if i == 1 {
			// there is something odd about "oa" that makes it show up as something other than just those two letters
			trimrunes = 3
		}
		rv := []rune(values[i])
		if len(rv) > trimrunes {
			rv = rv[trimrunes:]
		} else {
			rv = []rune{}
		}
		values[i] = lat.ConvertLatinDiacriticals(betacode.ReplacePercentSigns(string(rv)))
	}

	var outvals []string

	outvals = append(outvals, fmt.Sprintf(PUB1, values[0]))
	outvals = append(outvals, fmt.Sprintf(REGION, values[1]))
	outvals = append(outvals, fmt.Sprintf(CITY, values[2]))
	if values[3] != "" {
		outvals = append(outvals, fmt.Sprintf(DATE, values[3]))
	}
	if values[4] != "" {
		outvals = append(outvals, fmt.Sprintf(DATE, values[4]))
	}
	if values[5] != "" {
		outvals = append(outvals, fmt.Sprintf(PUB2, values[5]))
	}

	j := strings.Join(outvals, "")
	j = strings.ReplaceAll(j, "`", "")
	j = purgedels.ReplaceAllString(j, "")
	return j
}
