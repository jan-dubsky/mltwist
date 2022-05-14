package main

import (
	"fmt"
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/disassemble"
	"mltwist/internal/consoleui/emulate"
	"mltwist/internal/deps"
	"mltwist/internal/elf"
	"mltwist/internal/parser"
	"mltwist/internal/riscv"
	"mltwist/internal/state"
	"mltwist/internal/state/memory"
	"mltwist/pkg/model"
	"os"
)

func parseElf(filename string) (*elf.Memory, model.Addr, *elf.Memory, error) {
	p, err := elf.NewParser(filename)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("cannot create elf parser: %w", err)
	}
	defer p.Close()

	code, err := p.MachineCode()
	if err != nil {
		err = fmt.Errorf("machine code cannot be extracted from ELF: %w", err)
		return nil, 0, nil, err
	}

	mem, err := p.Memory()
	if err != nil {
		err = fmt.Errorf("cannot extract program memory from ELF: %w", err)
		return nil, 0, nil, err
	}

	return code, p.Entrypoint(), mem, nil
}

func runIU(p *deps.Program, mem *elf.Memory) error {
	memBlocks := make([]memory.ByteBlock, len(mem.Blocks))
	for i, b := range mem.Blocks {
		memBlocks[i] = b
	}

	byteMem, err := memory.NewBytes(memBlocks)
	if err != nil {
		return fmt.Errorf("cannot create byte memory of a program: %w", err)
	}

	emulF := func(p *deps.Program, ip model.Addr) (consoleui.Mode, error) {
		m := memory.NewOverlay(byteMem, memory.NewSparse())

		stat := &state.State{
			Regs: state.NewRegMap(),
			Mems: memory.MemMap{
				riscv.MemoryKey: m,
			},
		}

		emul, err := emulate.New(p, ip, stat)
		if err != nil {
			return nil, fmt.Errorf("cannot create emulation mode: %w", err)
		}

		return emul, nil
	}

	disass := disassemble.New(p, emulF)
	ui, err := consoleui.New(disass)
	if err != nil {
		return fmt.Errorf("cannot create console UI: %w", err)
	}

	return ui.Run()
}

func run() error {
	if l := len(os.Args); l != 2 {
		return fmt.Errorf("unexpected number of arguments: %d", l)
	}
	filename := os.Args[1]

	code, entrypoint, memory, err := parseElf(filename)
	if err != nil {
		return fmt.Errorf("ELF parsing failed: %w", err)
	}

	riscvParser := riscv.NewParser(riscv.Variant64, riscv.ExtM)
	ins, err := parser.Parse(code, riscvParser)
	if err != nil {
		return fmt.Errorf("instruction parsing failed: %w", err)
	}

	program, err := deps.NewProgram(entrypoint, ins)
	if err != nil {
		return fmt.Errorf("cannot parse model: %w", err)
	}

	return runIU(program, memory)
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "mltwist: %s\n", err.Error())
		os.Exit(1)
	}
}
