package memory

import (
	"mltwist/internal/state/interval"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// Memory represents a single linear address space where expressions of
// different width are stored to and loaded from.
//
// Unlike standard program memory which stores individual bytes, our memory has
// to store generic expressions instead as during static analysis, our
// expressions might not be fully evaluated yet. Storing expression instead of
// bytes brings some challenges in mapping the expression model to a
// byte-oriented model. Namely a new write can rewrite just part of an
// expression written before. For this reason, we have to store an information
// which bytes of an expression are still valid and which of them were
// rewritten. On read, we are able to compose any expression as combination of
// written expression using by bit shifts, ANDs and OSs.
//
// The expression partial rewrite expression described above might result in
// duplication of expression. This happens if a new write rewrites just some
// bytes in the middle of a previously written expression. Such write results in
// 2 places in memory being calculated by the same expression which differ only
// by bytes taken. Given that every write can result in duplication of at most
// one expression, the ultimate growth of expression complexity (and memory
// consumption) is only linear in number of writes. Consequently the overall
// number of expression used to represent the state of memory after n writes is
// always O(n) independently on the fact whether expression splitting happens or
// not. So expression splitting is not much of an issue as it doesn't increase
// number of expressions to evaluate significantly.
type Memory interface {
	// Load loads w bytes from memory address addr.
	//
	// If any byte in range [addr, addr+w) is missing, this method returns
	// (nil, false). Caller can use Missing method of this interface to get
	// more detailed information which bytes are missing. If returned value
	// is true, the returned expression is always non-nil.
	Load(addr model.Addr, w expr.Width) (expr.Expr, bool)

	// Store stores expression ex of width w to memory address range [addr,
	// addr+w).
	Store(addr model.Addr, ex expr.Expr, w expr.Width)

	// Missing returns all intervals missing in the memory in range [addr,
	// addr+w). This method can be used to identify bytes which are missing
	// to a Load operation.
	Missing(addr model.Addr, w expr.Width) interval.Map[model.Addr]

	// Blocks returns all continuous memory address intervals the memory
	// stores.
	Blocks() interval.Map[model.Addr]
}
