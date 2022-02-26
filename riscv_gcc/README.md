# RISC-V cross compilation environment

Even though many OS distributions provide a package with RISC-V GCC toolchain,
it shows up that some properties of the RISC-V cross-compiler GCC binary depend
on build options and cannot be change afterwards. The build option we care of
is usege of compressed RISC-V instructions. Those GCC binaries by default use
RISC-V C extension (RVC) which specifies compressed instruction set. Even
though the man page lists an option to specify an architecture and
architectures without the C extension are listed in the man page and accepted
by the binary, it shows up that compressed instructions are generated even for
architecture settings excluding the C extension.

As we in our project don't care as much about the decompilation of RISC-V
machine code, but more about general decompilation, we have decided not to
support the RISC-V C extension in the decompiler. But to be able to do so, we
first need to get GCC compiler which can generate only non-compressed
instructions. And it [shows up][1] that we have to build one.

We could naturally build the RISC-V GCC compiler on our machines. But as
building a GCC is not a super-convenient process, we have decided to automate
and simplify it using `docker`.

## References

[1]: https://stackoverflow.com/questions/43704690/how-to-forbid-the-riscv-compressed-codes
