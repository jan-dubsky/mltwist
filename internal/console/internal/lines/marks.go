package lines

// MaxMarkLen is maximal allowed length of mark string for a line. Any attempt
// to set longer mark will result in panic.
const MaxMarkLen = 3

const (
	MarkMovedFrom = "<"
	MarkMovedTo   = ">"
)

const (
	MarkLowerBound = "vvv"
	MarkUpperBound = "^^^"
)
