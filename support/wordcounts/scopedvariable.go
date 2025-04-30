package wordcounts

// HeadwordLookupMap could produce a race condition; but it is read-only after write-once

var (
	HeadwordLookupMap = make(map[string][]string)
)
