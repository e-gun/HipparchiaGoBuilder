//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package worklines2dbprep

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/structs"
)

var (
	// needs to be suited to METADATATEMPLATE in hexrunner.go: `<hb-metadata_%s value="%s" />`
	metadatanotationfinder = regexp.MustCompile("<hb-metadata_(.*?) value=\"(.*?)\" />")

	// marginaltextfindera    = regexp.MustCompile("<hb-sp-marginaltext>(.[^<]?)</hb-sp-marginaltext>")
	// marginaltextfinderb will let you find `<hb-sp-marginaltext>⟩ <hb-fs-l-normal>corr.</hb-fs-l-normal></hb-sp-marginaltext>`
	marginaltextfinderb = regexp.MustCompile(`<hb-sp-marginaltext>(.*?)</hb-sp-marginaltext>`)
	// unconventionalfinder - otherwise you see "ἐπίτροπον επιτρον̣ον" in the text; will make it hard to search for "bad Greek" though...
	unconventionalfinder = regexp.MustCompile(`<hb-sp-unconventional_form_written_by_scribe>(.*?)</hb-sp-unconventional_form_written_by_scribe>`)
	alteredformfinder    = regexp.MustCompile(`<hb-sp-form_altered_by_scribe>(.*?)</hb-sp-form_altered_by_scribe>`)
	rectifiedformfinder  = regexp.MustCompile(`<hb-sp-rectified_form>(.*?)</hb-sp-rectified_form>`)
	alternativefinder    = regexp.MustCompile(`<hb-sp-alternative_reading>(.*?)</hb-sp-alternative_reading>`)
	discardedfinder      = regexp.MustCompile(`<hb-sp-discarded_form>(.*?)</hb-sp-discarded_form>`)

	skipcategories = []string{
		"newauthor",
		"newwork",
		// "workabbrev",
	}
)

func ExtractNotesAndMarginalText(lines []structs.DbWorkline) []structs.DbWorkline {
	// further candidates in leftcurlybrackets:
	// 		9: "<hb-sp-alternative_reading>",
	// not sure of how common they are...

	for i, line := range lines {
		lmu := line.MarkedUp
		md := metadatanotationfinder.MatchString(lmu)
		mt := marginaltextfinderb.MatchString(lmu)
		ucf := unconventionalfinder.MatchString(lmu)
		af := alteredformfinder.MatchString(lmu)
		rf := rectifiedformfinder.MatchString(lmu) // rf not in ddp, but in chr and ins
		alt := alternativefinder.MatchString(lmu)
		dis := discardedfinder.MatchString(lmu)

		if !md && !mt && !ucf && !af && !rf && !alt && !dis {
			continue
		}
		var newnotes []string
		if md {
			// remove the notes from the markedup line
			lines[i].MarkedUp = metadatanotationfinder.ReplaceAllString(lmu, "")
			notes := metadatanotationfinder.FindAllStringSubmatch(lmu, -1)
			combinednotes := metadatacollection(notes)
			newnotes = append(newnotes, combinednotes)
		}
		if mt {
			// remove the mt from the markedup line: one of the big places you will see this is Homer ("※", etc.)
			lines[i].MarkedUp = marginaltextfinderb.ReplaceAllString(lmu, "")

			// is there one and only one marginal text relative to any given line? looks like it, but... FindAll
			notes := marginaltextfinderb.FindAllStringSubmatch(lmu, -1)
			combinednotes := metadatacollection(notes)
			newnotes = append(newnotes, combinednotes)
		}
		if ucf {
			lines[i].MarkedUp = unconventionalfinder.ReplaceAllString(lmu, "")
			notes := unconventionalfinder.FindAllStringSubmatch(lmu, -1)
			combinednotes := ucfcollection(notes)
			newnotes = append(newnotes, combinednotes)
		}
		if af {
			lines[i].MarkedUp = alteredformfinder.ReplaceAllString(lmu, "")
			notes := alteredformfinder.FindAllStringSubmatch(lmu, -1)
			combinednotes := afcollection(notes)
			newnotes = append(newnotes, combinednotes)
		}
		if rf {
			lines[i].MarkedUp = rectifiedformfinder.ReplaceAllString(lmu, "")
			notes := rectifiedformfinder.FindAllStringSubmatch(lmu, -1)
			combinednotes := rfcollection(notes)
			newnotes = append(newnotes, combinednotes)
		}
		if alt {
			lines[i].MarkedUp = alternativefinder.ReplaceAllString(lmu, "")
			notes := alternativefinder.FindAllStringSubmatch(lmu, -1)
			combinednotes := altcollection(notes)
			newnotes = append(newnotes, combinednotes)
		}
		if dis {
			lines[i].MarkedUp = discardedfinder.ReplaceAllString(lmu, "")
			notes := discardedfinder.FindAllStringSubmatch(lmu, -1)
			combinednotes := discollection(notes)
			newnotes = append(newnotes, combinednotes)
		}
		// Cicero has `<hb-metadata_date value=\"Non. Dec. 61 (%1411)\" />`, etc.
		// the "%1411" should be "%14`11", i.e., "§11"
		// I wonder how many others we will end up doing by hand...
		nn := strings.Join(newnotes, " ")
		nn = strings.Replace(nn, "%14", "§", -1)
		lines[i].Annotations = nn
	}
	return lines
}

func metadatacollection(notes [][]string) string {
	var allnotes []string
	for _, note := range notes {
		kv := strings.Join(note[1:], ": ")
		skip := false
		for _, skipcat := range skipcategories {
			if strings.Contains(kv, skipcat) {
				skip = true
			}
		}
		if !skip {
			allnotes = append(allnotes, kv)
		}
	}
	combinednotes := strings.Join(allnotes, " · ")
	return combinednotes
}

func ucfcollection(vals [][]string) string {
	const (
		TMPL = `<hb-sp-unconventional_form_written_by_scribe>(%s)</hb-sp-unconventional_form_written_by_scribe>`
	)
	var collector []string
	for _, val := range vals {
		collector = append(collector, val[1])
	}
	combinednotes := strings.Join(collector, " ⁄ ")
	return "corrected: " + fmt.Sprintf(TMPL, combinednotes)
}

func afcollection(vals [][]string) string {
	const (
		TMPL = `<hb-sp-form_altered_by_scribe>[%s]</hb-sp-form_altered_by_scribe>`
	)
	var collector []string
	for _, val := range vals {
		collector = append(collector, val[1])
	}
	combinednotes := strings.Join(collector, " ⁄ ")
	return "altered: " + fmt.Sprintf(TMPL, combinednotes)
}

func rfcollection(vals [][]string) string {
	const (
		TMPL = `<hb-sp-rectified_form>{scil: %s}</hb-sp-rectified_form>`
	)
	var collector []string
	for _, val := range vals {
		collector = append(collector, val[1])
	}
	combinednotes := strings.Join(collector, " ⁄ ")
	return "rectified: " + fmt.Sprintf(TMPL, combinednotes)
}

func altcollection(vals [][]string) string {
	const (
		TMPL = `<hb-sp-alternative_reading>[alt： %s]</hb-sp-alternative_reading>` // note the full width colon
	)
	var collector []string
	for _, val := range vals {
		collector = append(collector, val[1])
	}
	combinednotes := strings.Join(collector, " ⁄ ")
	return "alternate: " + fmt.Sprintf(TMPL, combinednotes)
}

func discollection(vals [][]string) string {
	const (
		TMPL = `<hb-sp-discarded_form>{%s}</hb-sp-discarded_form>` // note the full width colon
	)
	var collector []string
	for _, val := range vals {
		collector = append(collector, val[1])
	}
	combinednotes := strings.Join(collector, " ⁄ ")
	return "discarded: " + fmt.Sprintf(TMPL, combinednotes)
}
