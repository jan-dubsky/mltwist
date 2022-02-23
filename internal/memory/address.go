package memory

// Address represents an arbitrary memory address in a program (user-space)
// address space.
//
// We might use plain old uint64 to represent any memory address and we would be
// most likely fine for following 10 years. On the other hand given that RISC-V
// already has 128bit instruction set described, it makes sense to introduce
// this alias and to make the code more variable in the future.
type Address = uint64
