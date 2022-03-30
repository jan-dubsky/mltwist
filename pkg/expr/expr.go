package expr

type Expr interface {
	Width() Width

	// internal is a blocker function which prevents defining custom
	// expressions anywhere else than in this package.
	//
	// Expr package defines a public interface and a list of operations any
	// other part of this module (namely any sort of analysis) must know and
	// be able to handle. For this reason we cannot allow anyone else to
	// introduce its own exprs. This approach as well fits well with the
	// philosophy that this package provides only markers of operations and
	// not the functionality itself as no-one can provide real
	// implementation of any Expr as there is simply no possible
	// implementation, just the marking functionality.
	internal()
}
