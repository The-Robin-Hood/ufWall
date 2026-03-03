package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/The-Robin-Hood/ufWall/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	if _, err := exec.LookPath("ufw"); err != nil {
		fmt.Println("Please first install ufw and try again.")
		os.Exit(1)
	}

	if os.Getenv("DEBUG") == "1" {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	p := tea.NewProgram(app.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
