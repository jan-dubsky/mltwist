// Package expr provides a model of simple but generic CPU.
//
//
//
// Our virtual processor operates with values or arbitrary byte width. All
// operations understand values as unsigned unless stated otherwise.
//
// Operation and value width represents width of input and output operands. If
// width of input operand is less then width of operation it's used in, the
// value is always zero extended. If width of operation is less than width of
// input operand width, the value is cropped to width of operation.
package expr

type Expr interface {
	// Width returns width of an expression.
	Width() Width

	// internal is a blocker function which prevents defining custom
	// expressions anywhere else than in this package.
	//
	// Expr package defines a public interface and a list of operations any
	// other part of this module (namely any sort of analysis) must know and
	// be able to handle. For this reason we cannot allow anyone else to
	// introduce its own exprs. This approach as well fits well with the
	// philosophy that this package provides only markers of operations and
	// not the functionality itself. As no-one can provide real
	// implementation of any Expr as there is simply no possible
	// implementation, just the marking functionality.
	internalExpr()
}

type Effect interface {
	// Width returns width of an effect.
	Width() Width
	// internalEffect is a blocker function which prevents defining custom
	// effects anywhere else than in this package.
	//
	// Expr package defines a public interface and a list of operations any
	// other part of this module (namely any sort of analysis) must know and
	// be able to handle. For this reason we cannot allow anyone else to
	// introduce its own effects. This approach as well fits well with the
	// philosophy that this package provides only markers of operations and
	// not the functionality itself. As no-one can provide real
	// implementation of any Effect as there is simply no possible
	// implementation, just the marking functionality.
	internalEffect()
}
