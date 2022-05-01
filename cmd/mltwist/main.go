package main

import (
	"fmt"
	"mltwist/internal/console/disassemble"
	"mltwist/internal/console/ui"
	"mltwist/internal/deps"
	"mltwist/internal/executable"
	"mltwist/internal/parser"
	"mltwist/internal/riscv"
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
	fmt.Printf("Bytes: %b\n", exec.Memory.Addr(exec.Entrypoint)[:4])

	riscvParser := riscv.NewParser(riscv.Variant64, riscv.ExtM)
	ins, err := parser.Parse(exec.Entrypoint, exec.Memory, riscvParser)
	if err != nil {
		return fmt.Errorf("instruction parsing failed: %w", err)
	}

	program, err := deps.NewProgram(ins.Entrypoint, ins.Instructions)
	if err != nil {
		return fmt.Errorf("cannot parse model: %w", err)
	}

	ui, err := ui.New(disassemble.New(program))
	if err != nil {
		return fmt.Errorf("cannot create UI: %w", err)
	}

	return ui.Run()
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "decomp: %s\n", err.Error())
		os.Exit(1)
	}
}
