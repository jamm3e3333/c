package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"

	"github.com/jamm3e3333/c/model"
)

func main() {
	// Initialize with the main menu model
	p := tea.NewProgram(model.NewMainModel())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
