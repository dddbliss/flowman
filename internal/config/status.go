package config

import tea "github.com/charmbracelet/bubbletea"

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

var Models []tea.Model

const (
	Board Status = iota
	NewTask
)
