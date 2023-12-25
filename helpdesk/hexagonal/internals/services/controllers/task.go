package controllers

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/task"
	"helpdesk/internals/services/persistance/memory"
)

type TaskRepo core.Repo[task.TaskID, task.Task]

type TaskController struct {
	Repo TaskRepo
}

func NewTaskController(repo TaskRepo) TaskController {
	return TaskController{
		Repo: repo,
	}
}

func NewMemoryTaskController() TaskController {
	repo := memory.NewTaskRepo("tasks")
	return TaskController{
		Repo: &repo,
	}
}

func (cc *TaskController) CreateTask(name string, subject string) (task.TaskID, error) {
	t := task.NewTask(name, subject)

	if err := cc.Repo.Save(t); err != nil {
		return "", err
	}

	return t.ID, nil
}

func (cc *TaskController) SetTaskOwner(task_id task.TaskID, owner core.Owner) error {
	t, err := cc.Repo.Get(task_id)
	if err != nil {
		return err
	}
	if err := t.SetOwner(owner); err != nil {
		return err
	}
	if err := cc.Repo.Save(t); err != nil {
		return err
	}

	return nil
}

func (cc *TaskController) GetTasks() ([]task.Task, error) {
	return cc.Repo.All()
}
