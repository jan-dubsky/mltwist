package lines

// MaxMarkLen is maximal allowed length of mark string for a line. Any attempt
// to set longer mark will result in panic.
const MaxMarkLen = 3

type Mark string

const (
	MarkMovedFrom Mark = "<"
	MarkMovedTo   Mark = ">"

	MarkLowerBound Mark = "vvv"
	MarkUpperBound Mark = "^^^"

	MarkFound Mark = "-->"
)
