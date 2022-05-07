package main

import (
	"fmt"
	"mltwist/internal/console"
	"mltwist/internal/deps"
	"mltwist/internal/elf"
	"mltwist/internal/parser"
	"mltwist/internal/riscv"
	"os"
)

func run() error {
	if l := len(os.Args); l != 2 {
		return fmt.Errorf("unexpected number of arguments: %d", l)
	}
	filename := os.Args[1]

	exec, err := elf.Parse(filename)
	if err != nil {
		return fmt.Errorf("ELF parsing failed: %w", err)
	}

	fmt.Printf("Entrypoint: 0x%x\n", exec.Entrypoint)
	fmt.Printf("Bytes: %b\n", exec.Memory.Address(exec.Entrypoint)[:4])

	riscvParser := riscv.NewParser(riscv.Variant64, riscv.ExtM)
	ins, err := parser.Parse(exec.Memory, riscvParser)
	if err != nil {
		return fmt.Errorf("instruction parsing failed: %w", err)
	}

	program, err := deps.NewProgram(exec.Entrypoint, ins)
	if err != nil {
		return fmt.Errorf("cannot parse model: %w", err)
	}

	console, err := console.New(program)
	if err != nil {
		return fmt.Errorf("cannot create UI: %w", err)
	}

	return console.Run()
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "decomp: %s\n", err.Error())
		os.Exit(1)
	}
}
