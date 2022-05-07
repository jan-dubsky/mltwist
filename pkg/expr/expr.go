// Package expr provides a set of instruction markers.
//
// The instruction set is designed to be as simple as possible, but generic
// enough to represent instruction from different CPU architectures. One way of
// thinking about primitives provided by this package is to reason about it as
// about CPU architecture, or an intermediate code in a compiler. Operations
// exposed by this package are further used in generic analysis of assembly
// code. Introduction of an intermediate code allows most of the further
// analysis to be platform independent.
//
// Our instruction set comprises of two high-level operation types. Those are
// expressions represented by Expr interface and effects which are represented
// by Effect interface. Expressions describe a way how to obtain values. Those
// are register reads, memory reads and constants. On top of that, we provide as
// well conditional expressions and binary operations which allow combining
// expressionsinto more complicated calculations. Effects on the other hand
// describe side-effects of instructions. Those are either memory or register
// writes.
//
// The model we use to represent expressions is basically a syntax tree. In
// other words, we use functional model of operation representation where any
// operation is an expression which can be used as input value of another
// operation. This model is more suitable that imperative model for many
// reasons. First of all, as our goal is to perform some analysis on top of
// instructions, we prefer representing them in as a syntax tree rather than as
// a sequence of operations. Second, our architecture is much simpler than a
// usual CPU architecture, we typically need multiple primitive operations to
// represent a single instruction. In an imperative model, we'd force anyone
// implementing translation from machine code to expression tree to reason where
// to store and load intermediate values. By functional model, we avoid all of
// this. And last but not least, having some intermediate value stored in "dummy
// registers" would complicate further analysis as the analyzer wold see
// non-existing dependencies in between instructions caused just by random
// collisions in intermediate registers
//
// Contrary to expressions, a single instruction is represented by a list of
// side effects which write into registers or memory. This might sound
// inconsistent with the previous paragraph as effects follow imperative
// programming model. But there is good reasoning for use of imperative model in
// between instructions. First of all, we want both the parser and the analysis
// to be able to operate on top of a single instruction. If we used functional
// model to represent a stream of instructions, we'd ultimately have to run an
// analysis on top of the whole program as the whole program would collapse into
// a single gigantic instruction. Second reason is that we'd have quite hard
// times to represent conditional jumps as those are not concept fitting any
// functional world description. The third reason is that we'd suddenly have to
// define the order in which the expression tree would be evaluated. That is
// because memory writes can overwrite other memory writes, so we'd need to say
// which write rewrites which. And that would require defining order in
// evaluation, which is something unusual un functional world. The forth reason
// is that we'd have a hard time to emulate single instruction step as such task
// would be represented by evaluation of part of the expression subtree. And
// last but not least the resulting expression describing the whole program
// would be most likely as big that we'd run out of memory.
//
// To sum up, we have Exprs which serve the purpose of calculation and then we
// have Effects which describe external state change. An instruction can be then
// represented as a set of effects where every effect comprises of a one or
// multiple trees of expressions. This way, description of a single instruction
// is a functional syntex tree, but instructions are not dependent in between
// one another. Moreover this model as well simplifies implementation of
// instruction parsers. It's also worth mentioning that this model is very close
// to the way how people reason and write algorithms - people write expressions
// which are functional representations of how to calculate something and then
// they write and assignment statement and store the value somewhere they can
// use it later.
//
//
//
// The virtual processor operates on values of arbitrary width. The reasoning
// behind this is genericity. We cannot simply pick one width of registers and
// memory accesses and use it to represent width of any CPU architecture.
// Selecting single width of register for a given CPU is also not sufficient as
// even in n-bit architectures, there are typically instructions operating less
// than n bytes. Consequently the only design which seems to be generic enough
// is the one where every operationpackage has a width in bytes and all
// operations truncate of extends values based on their widths.
//
// The width extension logic is following: If width of an input expression
// matches width of an argument, no conversion applies. If width of an operand
// is strictly less than width of the operation, the argument is zero-extended
// to the width of the operation. Analogously, if width of an operand is greater
// than width of an operation, the operand is truncated to the width of the
// operation. It's worth mentioning that width of the operation applies to all
// input values as well as to the output value of any operation. This behaviour
// again corresponds to behaviour of many CPU architectures which define
// operations on a fixed number of bits.
//
//
//
// In order to support as many architectures as possible, our instruction set
// has arbitrary number of registers. It's fully up to the parser how many
// registers it uses. There are architectures which doesn't use just numbers for
// register names. An example of such architecture might x86. It's true that all
// general purpose registers in x86 architecture have numbers, but it's also
// more than usual to refer general purpose registers by their common names
// (i.e. RSP, RBX etc.) rather than by numbers. Very similar statement applies
// to ARM where again, registers have numbers, but are grouped into several
// groups by their purpose. For all those reasons, it's makes sense to refer
// registers by strings rather than by just numbers.
//
// Analogously to registers, our architecture has arbitrary numbers of memory
// address spaces. The reasoning behind is again that there are architectures as
// x86, which have 2 address spaces - normal address space and I/O address
// space. Address spaces are again identified by string keys and address spaces
// with different keys are fully independent on one another.
//
// It's worth mentining that we might represent register accesses as memory
// reads and writes. By using maximal width of a register as stride in between
// individual register records, we could represent register file in a special
// address space. This approach is similar to encoding RAM computer model in a
// Turing machine model, where one tape is used to encode register file. If we
// decided to use this model we'd have to degrade register naming to plain
// numbers as memory addresses are always numbers. On the other hand such
// approach would work and it would be functionally equivalent to having
// separate register operations. The reason we provide register operations is
// just to simplify both the parsing and the code analysis.
//
// Some architectures have registers addresses by other register. A nice examply
// of such registers are MSRs in x86 architecture. There address of MSR register
// is taken from another register. As register operations do not allow
// addressing by another expression, emulation of MSG and similar registers,
// requires usage of registers-in-memory approach described above.
package expr

// Expr represents an arbitrary operation resulting in a value. For more
// detailed explanation, please read doc-comment of this package.
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

// Effect represents an arbitrary operation resulting in an external state
// change. For more detailed explanation, please read doc-comment of this
// package.
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
