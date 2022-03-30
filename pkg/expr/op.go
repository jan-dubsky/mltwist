package expr

type UnaryOp uint8

const (
	Negate UnaryOp = iota + 1
	Not
	Abs
	CastUnsigned
)

type BinaryOp uint8

const (
	Add BinaryOp = iota + 1
	Sub
	Mul
	Div
	Mod
	Rsh
	Lsh
	RshL
	And
	Or
	Xor
)
