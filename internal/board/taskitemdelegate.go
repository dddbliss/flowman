package board

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dddbliss/flowman/internal/config"
	"github.com/dddbliss/flowman/internal/models"
)

type taskItemDelegate struct {
	model  *Model
	status config.Status
}

func (d taskItemDelegate) Height() int                               { return 2 }
func (d taskItemDelegate) Spacing() int                              { return 1 }
func (d taskItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d taskItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(models.Task)
	if !ok {
		return
	}

	state := (index == m.Index() && d.model.focused == d.status)

	str := fmt.Sprintf("%s\n%s", config.Styles.ItemTitleStyle.Foreground(config.Styles.GetStatefulColor(state)).Render(i.Title()), config.Styles.ItemDescStyle.Render(i.Description()))

	out := config.Styles.ItemStyle.Render(str)

	if index == m.Index() {
		border := d.model.focused == d.status

		if border {
			config.Styles.SelectedItemStyle.Border(lipgloss.NormalBorder(), false, false, false, true)
		} else {
			config.Styles.SelectedItemStyle.Border(lipgloss.HiddenBorder(), false, false, false, true)
		}

		out = config.Styles.SelectedItemStyle.Render(str)
		config.Styles.SelectedItemStyle.BorderLeft(!border)
	}

	fmt.Fprint(w, out)
}
