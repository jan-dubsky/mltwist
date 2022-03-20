package main

import (
	"decomp/internal/console"
	"decomp/internal/deps"
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
	fmt.Printf("Bytes: %b\n", exec.Memory.Addr(exec.Entrypoint)[:4])

	riscvParser := riscv.NewParser(riscv.Variant64, riscv.ExtM)
	ins, err := parser.Parse(exec.Entrypoint, exec.Memory, riscvParser)
	if err != nil {
		return fmt.Errorf("instruction parsing failed: %w", err)
	}

	model, err := deps.NewModel(ins.Instructions)
	if err != nil {
		return fmt.Errorf("cannot parse model: %w", err)
	}

	/*
		for i, b := range model.Blocks() {
			if i > 0 {
				fmt.Printf("\n")
			}
			fmt.Printf("Basic block %d:\n", i)

			for j, in := range b.Instructions() {
				fmt.Printf("\t%00d (0x%x): %s\n", j, in.Address(), in.String())
			}
		}
	*/

	ui := console.NewUI(model)
	return ui.Run()
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "decomp: %s\n", err.Error())
		os.Exit(1)
	}
}
