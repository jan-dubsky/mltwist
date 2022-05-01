package ui

import "fmt"

type commandPrompt struct{}

func (commandPrompt) MinLines() int { return 2 }
func (commandPrompt) MaxLines() int { return 2 }
func (commandPrompt) Print(_ int) error {
	fmt.Printf("Enter command: ")
	return nil
}
