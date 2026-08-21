//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package authorbuildpipeline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/e-gun/HipparchiaGoBuilder/internal/global"
)

func AssignTLGMetadataToAuthors(tlgmetadata map[string]map[string][]string) {

	// 	allmetadata := map[string]map[string][]string{
	//		"gnn": gnn,
	//		"gnx": gnx,
	//		"dtg": dtg,
	//		"ept": ept,
	//		"loc": loc,
	//	}

	// DbAuthor:
	//    UID       string
	//    Language  string
	//    IDXname   string
	//    Name      string
	//    Shortname string
	//    Cleaname  string
	//    Genres    string
	//    RecDate   string
	//    ConvDate  int
	//    Location  string
	//    WorkList  []string

	// TheCanonAuMap: map[string]structs.DbAuthor

	for k, v := range tlgmetadata["loc"] {
		a, ok := global.TheCanonAuMap.Get(k)
		if ok {
			a.Location = strings.Join(v, "; ")
			global.TheCanonAuMap.Set(k, a)
		}
	}

	for k, v := range tlgmetadata["dtg"] {
		a, ok := global.TheCanonAuMap.Get(k)
		if ok {
			a.RecDate = v[0]
			global.TheCanonAuMap.Set(k, a)
		}
	}

	// use epithets for genres; but genres for genres might be better in the end
	for k, v := range tlgmetadata["ept"] {
		a, ok := global.TheCanonAuMap.Get(k)
		if ok {
			a.Genres = strings.Join(v, "; ")
			global.TheCanonAuMap.Set(k, a)
		}
	}

	// set language while we are here...
	for _, k := range global.TheCanonAuMap.Keys() {
		a, _ := global.TheCanonAuMap.Get(k)
		a.Language = "G"
		global.TheCanonAuMap.Set(k, a)
	}

	// build the shortname while here
	for _, k := range global.TheCanonAuMap.Keys() {
		a, _ := global.TheCanonAuMap.Get(k)
		// because other corpora might be part of this build and are in the mastermap
		a.Cleaname = idxnamecleaner(a.IDXname)
		global.TheCanonAuMap.Set(k, a)
	}
}

func idxnamecleaner(old string) string {
	const (
		TEMPLATE = `%s (%s)`
	)
	// elaborate because...
	// `<hb-fs-l-bold>Eratosthenes </hb-fs-l-bold><hb-fs-l-normal>et </hb-fs-l-normal><hb-fs-l-bold>Eratosthenica</hb-fs-l-bold> Philol.`
	// needs to be `Eratosthenes et Eratosthenica (Philol.)`

	cleaner1 := regexp.MustCompile("<hb-fs-l-bold>(.*?)</hb-fs-l-bold>")
	cleaner2 := regexp.MustCompile(`</hb-fs-l-bold>\s(.*?)$`)

	if !cleaner1.MatchString(old) {
		return old
	}
	sm := cleaner1.FindAllStringSubmatch(old, -1)
	var allnames []string
	for _, v := range sm {
		allnames = append(allnames, v[1])
	}
	clean := strings.Join(allnames, "et ")

	epthm := cleaner2.FindAllStringSubmatch(old, -1)
	if len(epthm) == 0 {
		return clean
	}
	epith := epthm[0][1]
	return fmt.Sprintf(TEMPLATE, clean, epith)
}
