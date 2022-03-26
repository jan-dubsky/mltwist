package model

// Register represents an arbitrary (but single) cpu register.
//
// As different architectures have different registry set, there is no more
// specific way how to represent them then just by numbers. Most likely there
// will be some way how to number registers on every platform. That is at least
// true for most common architectures (x86, x86_64, MIPS, ARM, RISC V).
//
// Note that some architectures can have some special meaning associated with
// some register numbers which cannot be depicted by this respresentation.
type Register uint64

// Registers represents a set of Register values.
type Registers map[Register]struct{}
