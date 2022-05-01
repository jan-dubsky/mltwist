package view

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

func Print(e Element) error {
	_, screenLines, err := terminal.GetSize(0)
	if err != nil {
		return fmt.Errorf("cannot get terminal size: %w", err)
	}

	// Clean the screen.
	fmt.Print("\033[H\033[2J")

	minLines := e.MinLines()
	if screenLines < minLines {
		fmt.Printf("screen height is not sufficient: %d < %d", screenLines, minLines)
		return nil
	}

	lines := e.MaxLines()
	if lines < 0 || lines > screenLines {
		lines = screenLines
	}

	err = e.Print(lines)
	if err != nil {
		return fmt.Errorf("element printing failed: %w", err)
	}

	return nil
}
