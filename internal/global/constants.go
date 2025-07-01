package global

const (
	NAME    = "HipparchiaGoBuilder"
	VERSION = "0.9.1"
)

// the following *must* be synchronized with HGS's corporaconst.go values

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
