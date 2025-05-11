package main

import (
	// "fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type OutputFunc func() string

// RunTUI starts a Bubble Tea TUI that displays the output of the given function
func RunTUI(outputFn OutputFunc, refreshRate time.Duration) error {
	initialModel := tuiModel{
		outputFn:    outputFn,
		refreshRate: refreshRate,
		// lastRefresh: time.Now(),
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// tuiModel represents the state of the TUI
type tuiModel struct {
	output      string
	outputFn    OutputFunc
	refreshRate time.Duration
	// lastRefresh time.Time
}

type tickMsg time.Time
type outputMsg string

func (m tuiModel) Init() tea.Cmd {
	// Run the function immediately and set up the ticker
	return tea.Batch(
		executeOutputFunc(m.outputFn),
		tickEvery(m.refreshRate),
	)
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			// m.lastRefresh = time.Now()
			// m.refreshCount++
			return m, executeOutputFunc(m.outputFn)
		}

	case outputMsg:
		m.output = string(msg)
		return m, nil

	case tickMsg:
		// m.lastRefresh = time.Now()
		// m.refreshCount++
		return m, tea.Batch(
			executeOutputFunc(m.outputFn),
			tickEvery(m.refreshRate),
		)
	}

	return m, nil
}

func (m tuiModel) View() string {
	return m.output
}

// executeOutputFunc runs the output function and returns its result as a message
func executeOutputFunc(fn OutputFunc) tea.Cmd {
	return func() tea.Msg {
		result := fn()
		return outputMsg(result)
	}
}

func tickEvery(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
