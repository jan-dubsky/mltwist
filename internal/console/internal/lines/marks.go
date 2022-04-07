package lines

// MaxMarkLen is maximal allowed length of Mark string for a line. Any attempt
// to set longer mark will result in panic.
const MaxMarkLen = 3

// Mark is a printable and human understandable string which symbolically
// describes some additional properties of a line.
//
// All marks has to be shorter or equal length as MaxMarkLen. An attempt to set
// a longer line Mark will panic.
type Mark string

const (
	MarkNone Mark = ""

	MarkMovedFrom Mark = "<"
	MarkMovedTo   Mark = ">"

	MarkLowerBound Mark = "vvv"
	MarkUpperBound Mark = "^^^"

	MarkErrMovedFrom Mark = "!<"
	MarkErrMovedTo   Mark = "!>"

	MarkErr Mark = "!"
)
