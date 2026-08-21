//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package global

const (
	NAME    = "HipparchiaGoBuilder"
	VERSION = "1.0.0"
)

// the following *must* be synchronized with HGS's corporaconst.go values
// if you change these values HGS will be unable to find any corpus until you change those values there and rebuild it

const (
	AUNAMELEN         = 4
	WKIDLEN           = 3
	ABBREVLEN         = 3
	AUIDLEN           = ABBREVLEN + AUNAMELEN
	AUTHWORKSEPARATOR = "w"
	TLGABBREV         = "tlg" // this and next should *not* be config-able because they *must* match HGS values
	LATABBREV         = "lat" // any edits to this code require simultaneous edits to that code
	INSABBREV         = "inx"
	DDPABBREV         = "dpx"
	CHRABBREV         = "chx"
	DDPFIRSTPASS      = "ddp"
	INSFIRSTPASS      = "ins"
	CHRFIRSTPASS      = "chr"
)
