package model

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainModel represents the main menu
type MainModel struct {
	choices  []string
	cursor   int
	selected int
}

// NewMainModel creates a new main menu model
func NewMainModel() MainModel {
	return MainModel{
		choices:  []string{"Coin Flip", "Dice Roll", "Blackjack", "Quit"},
		cursor:   0,
		selected: -1,
	}
}

// Init initializes the model
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Quit the program
		case "ctrl+c", "q":
			return m, tea.Quit

		// Move cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Move cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// Select current option
		case "enter", " ":
			m.selected = m.cursor

			// Handle the selected option
			switch m.cursor {
			case 0: // Coin Flip
				return NewCoinFlipModel(), nil
			case 1: // Dice Roll
				// Placeholder for dice roll game
				// Will be implemented later
				return m, nil
			case 2: // Blackjack
				// Placeholder for blackjack game
				// Will be implemented later
				return m, nil
			case 3: // Quit
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	itemStyle  = lipgloss.NewStyle().PaddingLeft(4)
	selected   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#7D56F4")).SetString("▸ ")
	unselected = lipgloss.NewStyle().PaddingLeft(2).SetString("  ")
)

// View renders the main menu
func (m MainModel) View() string {
	s := titleStyle.Render("CLI GAMES") + "\n\n"

	for i, choice := range m.choices {
		cursor := unselected.String()
		if m.cursor == i {
			cursor = selected.String()
		}

		s += fmt.Sprintf("%s%s\n", cursor, choice)
	}

	s += "\n" + itemStyle.Render("Press q to quit, up/k and down/j to navigate, enter to select") + "\n"

	return s
}
