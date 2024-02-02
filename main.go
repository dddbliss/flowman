package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dddbliss/flowman/internal/board"
	"github.com/dddbliss/flowman/internal/config"
	"github.com/dddbliss/flowman/internal/data"
	"github.com/dddbliss/flowman/internal/tasks"
)

func main() {
	data.Open()
	config.Models = []tea.Model{board.New(), tasks.NewTaskForm(config.Todo)}
	m := config.Models[config.Board]
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
