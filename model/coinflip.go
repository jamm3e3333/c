package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jamm3e3333/c/game"
)

// Key mapping for coin flip
type coinFlipKeyMap struct {
	Flip key.Binding
	Back key.Binding
	Quit key.Binding
}

// Default key mappings
var defaultCoinFlipKeyMap = coinFlipKeyMap{
	Flip: key.NewBinding(
		key.WithKeys("f", "enter", " "),
		key.WithHelp("f/enter/space", "flip coin"),
	),
	Back: key.NewBinding(
		key.WithKeys("b", "esc"),
		key.WithHelp("b/esc", "back to menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

// Animation frames for coin flip - showing rotation instead of size changes
var coinFrames = []string{
	// Heads up
	`
   ______
  /      \
 |        |
 |  COIN  |
 |        |
  \______/
`,
	// Start rotating - edge view
	`
    ____
   /    \
  |      |
  |______|  
`,
	// Rotating - thin edge view
	`
     __
    |  |
    |  |
    |__|  
`,
	// Almost vertical
	`
      |
      |
      |
      |  
`,
	// Vertical edge
	`
      |
      |
      |
      |  
`,
	// Starting to show other side
	`
     __
    |  |
    |  |
    |__|  
`,
	// More of other side visible
	`
    ____
   /    \
  |      |
  |______|  
`,
	// Full other side
	`
   ______
  /      \
 |        |
 |        |
 |        |
  \______/
`,
	// Starting to rotate back
	`
    ____
   /    \
  |      |
  |______|  
`,
	// Edge view again
	`
     __
    |  |
    |  |
    |__|  
`,
	// Almost vertical again
	`
      |
      |
      |
      |  
`,
	// Vertical edge
	`
      |
      |
      |
      |  
`,
	// Starting to show original side
	`
     __
    |  |
    |  |
    |__|  
`,
	// More of original side visible
	`
    ____
   /    \
  |      |
  |______|  
`,
}

// CoinFlipModel represents the coin flip game view
type CoinFlipModel struct {
	coinFlip      *game.CoinFlip
	keys          coinFlipKeyMap
	flipping      bool
	animationStep int
	result        game.CoinSide
	showResult    bool
	showingResult time.Time // Tracks when the result started showing
	autoFlip      bool      // Whether to automatically flip again after delay
}

// NewCoinFlipModel creates a new coin flip model
func NewCoinFlipModel() CoinFlipModel {
	return CoinFlipModel{
		coinFlip:      game.NewCoinFlip(),
		keys:          defaultCoinFlipKeyMap,
		flipping:      false,
		animationStep: 0,
		showResult:    false,
		showingResult: time.Time{},
		autoFlip:      true, // Enable auto-flipping by default
	}
}

// Init initializes the model
func (m CoinFlipModel) Init() tea.Cmd {
	return nil
}

// Custom message types for the coin flip animation and auto-flip timer
type flipMsg struct{}
type checkAutoFlipMsg struct{}

// Update handles messages and updates the model
func (m CoinFlipModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Back):
			// Return to main menu
			return NewMainModel(), nil

		case key.Matches(msg, m.keys.Flip):
			if !m.flipping {
				m.flipping = true
				m.showResult = false
				return m, m.flipCoin()
			}
		}

	case flipMsg:
		if m.animationStep < len(coinFrames)-1 {
			m.animationStep++
			return m, m.flipCoin()
		} else {
			m.result = m.coinFlip.Flip()
			m.flipping = false
			m.showResult = true
			m.animationStep = 0
			m.showingResult = time.Now() // Record when we started showing the result

			// Start the timer for auto-flip if enabled
			if m.autoFlip {
				return m, m.checkAutoFlip()
			}
			return m, nil
		}

	case checkAutoFlipMsg:
		// Check if it's time to flip again (after 6 seconds)
		if m.showResult && m.autoFlip && !m.flipping {
			elapsed := time.Since(m.showingResult)
			if elapsed >= 6*time.Second {
				// Time to flip again
				m.flipping = true
				m.showResult = false
				return m, m.flipCoin()
			} else {
				// Not time yet, check again in a bit
				return m, m.checkAutoFlip()
			}
		}
	}

	return m, nil
}

// flipCoin is a command that simulates the coin flipping animation
func (m CoinFlipModel) flipCoin() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return flipMsg{}
	})
}

// checkAutoFlip is a command that checks if it's time to automatically flip again
func (m CoinFlipModel) checkAutoFlip() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return checkAutoFlipMsg{}
	})
}

// View renders the coin flip view
func (m CoinFlipModel) View() string {
	var s strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("COIN FLIP")

	s.WriteString(title + "\n\n")

	// Coin animation or result
	coinStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Align(lipgloss.Center).
		Width(40)

	if m.flipping {
		s.WriteString(coinStyle.Render(coinFrames[m.animationStep]) + "\n")
	} else if m.showResult {
		resultStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			Background(lipgloss.Color("#333333")).
			Padding(1, 3).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700"))

		coin := `
   ______
  /      \
 |        |
 |   ` + string(m.result) + `   |
 |        |
  \______/
`
		s.WriteString(coinStyle.Render(coin) + "\n\n")
		s.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(40).Render(
			resultStyle.Render("Result: "+string(m.result)),
		) + "\n\n")

		// Show auto-flip status and countdown
		if m.autoFlip {
			elapsed := time.Since(m.showingResult)
			remaining := 6*time.Second - elapsed
			if remaining < 0 {
				remaining = 0
			}
			countdown := fmt.Sprintf("Next flip in %.1f seconds...", remaining.Seconds())
			s.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(40).Foreground(lipgloss.Color("#888888")).Render(countdown) + "\n")
		}
	} else {
		s.WriteString(coinStyle.Render(coinFrames[0]) + "\n\n")
		s.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(40).Render(
			"Press F or ENTER to flip the coin",
		) + "\n\n")
	}

	// Help
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("f/enter: flip • b/esc: back to menu • q: quit")

	s.WriteString("\n" + helpStyle + "\n")

	return s.String()
}
