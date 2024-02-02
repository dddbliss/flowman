package tasks

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dddbliss/flowman/internal/config"
	"github.com/dddbliss/flowman/internal/data"
	"github.com/dddbliss/flowman/internal/models"
)

type NewTaskModel struct {
	focused     config.Status
	title       textinput.Model
	description textarea.Model
}

func NewTaskForm(focused config.Status) *NewTaskModel {
	form := &NewTaskModel{focused: focused}
	form.title = textinput.New()
	form.title.Placeholder = "take out the trash"
	form.title.Focus()
	form.description = textarea.New()
	form.description.SetHeight(12)
	form.description.ShowLineNumbers = false
	form.description.Placeholder = "make sure to separate your recycle from the rest of the trash"
	return form
}

type TaskAddedMsg struct {
	Task models.Task
}

func (m NewTaskModel) CreateTask() tea.Msg {
	task := models.NewTask(100, m.focused, m.title.Value(), m.description.Value())
	data.CreateTask(task)
	return TaskAddedMsg{Task: task}
}

func (m NewTaskModel) Init() tea.Cmd {
	return nil
}

func (m NewTaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		config.ScreenWidth = msg.Width
		config.ScreenHeight = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return config.Models[config.Board], nil
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				config.Models[config.NewTask] = m
				return config.Models[config.Board], m.CreateTask
			}
		}
	}

	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}

func (m NewTaskModel) View() string {
	ui := lipgloss.JoinVertical(lipgloss.Left,
		config.Styles.TitleStyle.Render("Create new task"),
		config.Styles.FormFieldStyle.Render("task title:"),
		config.Styles.FormFieldStyle.Render(m.title.View()),
		config.Styles.FormFieldStyle.Render("task description:"),
		config.Styles.FormFieldStyle.Render(m.description.View()))

	dialog := lipgloss.Place(config.ScreenWidth, config.ScreenHeight,
		lipgloss.Center, lipgloss.Center,
		config.Styles.DialogBoxStyle.Width(config.ScreenWidth/2).Render(ui),
		lipgloss.WithWhitespaceChars("·"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("238")),
	)
	return dialog

}
