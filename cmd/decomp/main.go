package main

import (
	"decomp/internal/executable"
	"decomp/internal/parser"
	"decomp/internal/riscv"
	"fmt"
	"os"
)

func run() error {
	if l := len(os.Args); l != 2 {
		return fmt.Errorf("unexpected number of arguments: %d", l)
	}
	filename := os.Args[1]

	exec, err := executable.Parse(filename)
	if err != nil {
		return fmt.Errorf("ELF parsing failed: %w", err)
	}

	fmt.Printf("Entrypoint: 0x%x\n", exec.Entrypoint)
	fmt.Printf("Bytes: %b\n", exec.Bytes(exec.Entrypoint, 4))

	_, err = parser.Parse(exec.Entrypoint, &exec, &riscv.ParsingStrategy{})
	if err != nil {
		return fmt.Errorf("instruction parsing failed: %w", err)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "decomp: %s\n", err.Error())
		os.Exit(1)
	}
}
