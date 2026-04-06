package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hackIDLE/fedramp-browser/internal/tui"
)

func main() {
	refresh := flag.Bool("refresh", false, "Force fresh fetch, ignoring cache")
	flag.Parse()

	var opts []tui.ModelOption
	if *refresh {
		opts = append(opts, tui.WithRefresh(true))
	}

	p := tea.NewProgram(
		tui.NewModel(opts...),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
