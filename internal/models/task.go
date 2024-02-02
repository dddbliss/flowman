package models

import "github.com/dddbliss/flowman/internal/config"

type Task struct {
	id          int
	status      config.Status
	title       string
	description string
}

func (t *Task) Status() config.Status {
	return t.status
}

func (t *Task) Id() int {
	return t.id
}

func (t *Task) Next() {
	if t.status == config.Done {
		t.status = config.Todo
	} else {
		t.status++
	}
}

// implement list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func (t Task) Title() string {
	return firstN(t.title, 20)
}

func (t Task) Description() string {
	return firstN(t.description, 20)
}

func NewTask(id int, status config.Status, title, description string) Task {
	return Task{id: id, status: status, title: title, description: description}
}
