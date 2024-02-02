package board

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dddbliss/flowman/internal/config"
	"github.com/dddbliss/flowman/internal/data"
	"github.com/dddbliss/flowman/internal/models"
	"github.com/dddbliss/flowman/internal/tasks"
)

type Model struct {
	keys     keyMap
	help     help.Model
	focused  config.Status
	lists    []list.Model
	loaded   bool
	quitting bool
}

func New() *Model {
	return &Model{keys: keys,
		help: help.New()}
}

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(models.Task)
	m.lists[selectedTask.Status()].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	data.UpdateStatus(selectedTask.Id(), selectedTask.Status())
	m.lists[selectedTask.Status()].InsertItem(len(m.lists[selectedTask.Status()].Items()), list.Item(selectedTask))
	return nil
}

func (m *Model) Next() {
	if m.focused == config.Done {
		m.focused = config.Todo
	} else {
		m.focused++
	}
}

func (m *Model) Previous() {
	if m.focused == config.Todo {
		m.focused = config.Done
	} else {
		m.focused--
	}
}

// TODO: Call on tea.WindowSizeMsg
func (m *Model) initLists(width, height int) {
	d := taskItemDelegate{}

	defaultList := list.New([]list.Item{}, d, width, height)
	defaultList.SetShowHelp(false)
	defaultList.SetFilteringEnabled(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// To Do
	m.lists[config.Todo].Title = "To Do"
	m.lists[config.Todo].SetDelegate(taskItemDelegate{model: m, status: config.Todo})
	todoTasks, _ := data.GetTasksByStatus(config.Todo)
	var todoItems []list.Item

	for _, task := range todoTasks {
		todoItems = append(todoItems, task)
	}

	m.lists[config.Todo].SetItems(todoItems)

	// In Progress
	m.lists[config.InProgress].Title = "In Progress"
	m.lists[config.InProgress].SetDelegate(taskItemDelegate{model: m, status: config.InProgress})
	inProgTasks, _ := data.GetTasksByStatus(config.InProgress)
	var inProgItems []list.Item

	for _, task := range inProgTasks {
		inProgItems = append(inProgItems, task)
	}

	m.lists[config.InProgress].SetItems(inProgItems)

	// Done

	m.lists[config.Done].SetDelegate(taskItemDelegate{model: m, status: config.Done})
	m.lists[config.Done].Title = "Done"
	doneTasks, _ := data.GetTasksByStatus(config.Done)
	var doneItems []list.Item

	for _, task := range doneTasks {
		doneItems = append(doneItems, task)
	}

	m.lists[config.Done].SetItems(doneItems)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		config.ScreenWidth = msg.Width
		config.ScreenHeight = msg.Height

		x, y := config.Styles.ColumnStyle.GetFrameSize()

		if !m.loaded {
			config.Styles.ColumnStyle.Width((msg.Width / config.Divisor) - x)
			config.Styles.FocusedStyle.Width((msg.Width / config.Divisor) - x)
			config.Styles.FocusedStyle.Height(msg.Height - y - 1)
			config.Styles.ColumnStyle.Height(msg.Height - y - 1)
			m.initLists(msg.Width, msg.Height-y-2)
			m.help.Width = msg.Width
			m.loaded = true
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Previous()
		case "right", "l":
			m.Next()
		case "enter":
			m.MoveToNext()
		case "n":
			config.Models[config.Board] = m
			config.Models[config.NewTask] = tasks.NewTaskForm(m.focused)
			return config.Models[config.NewTask].Update(nil)
		}

	case tasks.TaskAddedMsg:
		task := msg.Task
		return m, m.lists[task.Status()].InsertItem(len(m.lists[task.Status()].Items()), task)
	}

	var cmd tea.Cmd
	for i := 0; i < len(m.lists); i++ {
		m.lists[i].SetDelegate(taskItemDelegate{model: &m, status: config.Status(i)})
		m.lists[i].Update(msg)
	}

	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.loaded {
		todoView := m.lists[config.Todo].View()
		inProgressView := m.lists[config.InProgress].View()
		doneView := m.lists[config.Done].View()

		out := ""
		switch m.focused {
		case config.InProgress:
			out = lipgloss.JoinHorizontal(lipgloss.Left, config.Styles.ColumnStyle.Render(todoView), config.Styles.FocusedStyle.Render(inProgressView), config.Styles.ColumnStyle.Render(doneView))
		case config.Done:
			out = lipgloss.JoinHorizontal(lipgloss.Left, config.Styles.ColumnStyle.Render(todoView), config.Styles.ColumnStyle.Render(inProgressView), config.Styles.FocusedStyle.Render(doneView))
		default:
			out = lipgloss.JoinHorizontal(lipgloss.Left, config.Styles.FocusedStyle.Render(todoView), config.Styles.ColumnStyle.Render(inProgressView), config.Styles.ColumnStyle.Render(doneView))
		}

		helpView := m.help.View(m.keys)
		out = out + "\n" + helpView
		return out
	} else {
		return "loading..."
	}
}
