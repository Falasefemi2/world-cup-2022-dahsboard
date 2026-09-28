package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/Falasefemi2/worldcupdashboard/data/local"
	"github.com/Falasefemi2/worldcupdashboard/ui"
)

func main() {
	dashboard := ui.NewDashboard(&local.Client{})
	p := tea.NewProgram(dashboard)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oh no, there's been an error: %v", err)
		os.Exit(1)
	}
}
