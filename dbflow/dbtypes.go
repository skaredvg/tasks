package dbflow

import (
	"time"
)

type User struct {
	id   int
	Name string
}

func (t *User) SetID(id int) {
	t.id = id
}

func (t *User) ID() int {
	return t.id
}

type Label struct {
	id   int
	Name string
}

func (l *Label) SetID(id int) {
	l.id = id
}

func (l Label) ID() int {
	return l.id
}

type Task struct {
	id      int
	opened  time.Duration
	closed  time.Duration
	Author  User
	Assign  User
	Title   string
	Content string
	Labels  []Label
}

func (t *Task) SetID(id int) {
	t.id = id
}

func (t Task) ID() int {
	return t.id
}

func (t *Task) SetOpened(o int) {
	t.opened = time.Duration(o)
}

func (t *Task) SetClosed(c int) {
	t.opened = time.Duration(c)
}

type Tasks []Task
