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
	TLGABBREV         = "tlg"
	LATABBREV         = "lat"
	INSABBREV         = "inx"
	DDPABBREV         = "dpx"
	CHRABBREV         = "chx"
	DDPFIRSTPASS      = "ddp"
	INSFIRSTPASS      = "ins"
	CHRFIRSTPASS      = "chr"
)
