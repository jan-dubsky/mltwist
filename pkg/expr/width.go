package expr

// Width is width of an expression result in bytes.
//
// TODO: Doc-comment why uint8.
type Width uint8

const (
	Width8Bit Width = 1 << iota
	Width16Bit
	Width32Bit
	Width64Bit
)
