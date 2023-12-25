package memory

import (
	"errors"
	"helpdesk/internals/core/task"
)

type TaskRepo struct {
	name string
	c    map[task.TaskID]task.Task
}

var (
	ErrTaskNotFound = errors.New("task not found")
)

func NewTaskRepo(name string) TaskRepo {
	return TaskRepo{
		name: name,
		c:    make(map[task.TaskID]task.Task),
	}
}

func (r TaskRepo) Get(id task.TaskID) (task.Task, error) {
	if c, ok := r.c[id]; ok {
		return c, nil
	} else {
		return task.Task{}, ErrTaskNotFound
	}
}

func (r TaskRepo) All() ([]task.Task, error) {
	tasks := make([]task.Task, 0)
	for _, c := range r.c {
		tasks = append(tasks, c)
	}

	return tasks, nil
}

func (r *TaskRepo) Save(c task.Task) error {
	r.c[c.ID] = c
	return nil
}

func (r *TaskRepo) Delete(id task.TaskID) error {
	if _, ok := r.c[id]; ok {
		delete(r.c, id)
		return nil
	} else {
		return ErrTaskNotFound
	}
}
