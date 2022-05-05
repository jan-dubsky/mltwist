package emulate

import (
	"fmt"
	"math/big"
	"mltwist/internal/console/internal/linereader"
	"mltwist/internal/emulator"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

var _ emulator.StateProvider = &stateProvider{}

type stateProvider struct{}

func (*stateProvider) Register(key expr.Key, w expr.Width) expr.Const {
	return readRegister(key, w)
}

func (*stateProvider) Memory(key expr.Key, addr model.Addr, w expr.Width) expr.Const {
	p := fmt.Sprintf(
		"Please enter value of memory %s at address 0x%x (%d) [%d bytes]: ",
		key, addr, addr, w)
	return readValueNoErr(p, w)
}

func readRegister(key expr.Key, w expr.Width) expr.Const {
	p := fmt.Sprintf("Please enter value of register %s [%d bytes]: ", key, w)
	return readValueNoErr(p, w)
}

func readValueNoErr(prompt string, w expr.Width) expr.Const {
	for {
		fmt.Print(prompt)

		c, err := readValue(w)
		if err == nil {
			return c
		}

		fmt.Printf("error: %s\n", err.Error())
		_, _ = linereader.ReadLine()
	}
}

func readValue(w expr.Width) (expr.Const, error) {
	line, err := linereader.ReadLine()
	if err != nil {
		return expr.Const{}, fmt.Errorf("cannot read line: %w", err)
	}

	if len(line) == 0 {
		return expr.Const{}, fmt.Errorf("value entered cannot be empty")
	}

	// For some reason big.Int.SetString allows underscores at some
	// positions. To avoid ambiguity and not to leak implementation details
	// we prohibit them at all.
	for i, c := range line {
		if c == '_' {
			return expr.Const{}, fmt.Errorf("underscore at index %d", i)
		}
	}

	num, ok := (&big.Int{}).SetString(line, 0)
	if !ok {
		return expr.Const{}, fmt.Errorf("cannot parse number: %s", line)
	}

	// Convert to little endian.
	bs := num.Bytes()
	revertBytes(bs)

	abs := expr.NewConst(bs, w)
	if num.Sign() >= 0 {
		return abs, nil
	}

	c := expr.NewBinary(expr.Sub, expr.Zero, abs, w)
	return exprtransform.ConstFold(c).(expr.Const), nil
}
