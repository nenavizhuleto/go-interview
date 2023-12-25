package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
)

type TaskFilter struct {
	ID string `json:"id"`
}

type Task struct {
	ID           string           `json:"id,omitempty" db:"id"`
	HumanNumber  int              `json:"humanNumber,omitempty" db:"human_number"`
	UserID       string           `json:"user_id,omitempty" db:"user_id"`
	Name         string           `json:"name" db:"name"`
	Subject      string           `json:"subject" db:"subject"`
	Status       string           `json:"status,omitempty" db:"status"`
	TimeCreated  *TimeCreated     `json:"timeCreated,omitempty" db:"time_created"`
	Comments     []Comment        `json:"comments,omitempty"`
	Responsible  Employee         `json:"responsible,omitempty"`
	IsUrgent     bool             `json:"isUrgent"`
	IsTemplate   bool             `json:"isTemplate"`
	BeforeCreate BeforeCreateFunc `json:"-" bson:"-"`
}

const tasks = "tasks"

type BeforeCreateFunc func(*Task) error

func TasksAll() ([]Task, error) {
	coll := data.GetCollection(tasks)
	cursor, err := coll.Find(nil, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	var tasks []Task
	if err := cursor.All(nil, &tasks); err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	return tasks, nil
}

func TasksByUserID(id string) ([]Task, error) {
	coll := data.GetCollection(tasks)
	cursor, err := coll.Find(nil, bson.D{{Key: "userid", Value: id}})
	if err != nil {
		return nil, fmt.Errorf("allByUser: %w", err)
	}

	var tasks []Task
	if err := cursor.All(nil, &tasks); err != nil {
		return nil, fmt.Errorf("allByUser: %w", err)
	}

	return tasks, nil
}

func (t *Task) beforeCreateHook() error {
	if t.BeforeCreate != nil {
		return t.BeforeCreate(t)
	} else {
		return nil
	}
}

func (t *Task) Create() error {
	if err := t.beforeCreateHook(); err != nil {
		return fmt.Errorf("beforeCreate: %w", err)
	}

	coll := data.GetCollection(tasks)
	if _, err := coll.InsertOne(nil, t); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	return nil
}

func (t *Task) Update() error {
	collection := data.GetCollection(tasks)
	if _, err := collection.UpdateOne(nil, bson.D{{Key: "id", Value: t.ID}}, bson.D{{Key: "$set", Value: t}}); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (t *Task) Delete() error {
	collection := data.GetCollection(tasks)

	if t.ID == "" {
		return fmt.Errorf("delete: id is empty")
	}

	if res, err := collection.DeleteOne(nil, bson.D{{Key: "id", Value: t.ID}}); err != nil {
		return fmt.Errorf("delete: %w", err)
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("delete: task not found")
	}

	return nil
}

func (t *Task) Save() error {
	db := data.DB
	var exists Task
	err := db.Get(&exists, "SELECT id FROM tasks WHERE id = ?", t.ID)

	if err != nil {
		return t.Create()
	} else {
		return t.Update()
	}
}

func (t *Task) Prettify() error {
	oldTime := t.TimeCreated.Value
	oldTime = strings.Replace(oldTime, "T", " ", 1)
	tm, err := time.Parse(TimeFormatString, oldTime)
	if err != nil {
		return err
	}
	t.TimeCreated.Value = tm.Format(time.DateTime)
	return nil
}

func GetTaskForUser(uid string, id string) (*Task, error) {
	db := data.DB

	var task Task
	if err := db.Get(&task, "SELECT * FROM tasks WHERE user_id = ? AND id = ?", uid, id); err != nil {
		return nil, err
	}

	return &task, nil
}

func GetTasksForUser(uid string) ([]Task, error) {
	db := data.DB

	var tasks []Task
	if err := db.Select(&tasks, "SELECT * FROM tasks WHERE user_id = ?", uid); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return tasks, nil
}

func UpdateTaskForUser(uid string, task *Task) error {
	db := data.DB
	task.UserID = uid
	_, err := db.NamedExec("UPDATE tasks SET name = :name, status = :status WHERE id = :id AND user_id = :user_id", task)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTaskStatus(id string, status string) error {
	db := data.DB
	_, err := db.Exec("UPDATE tasks SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return err
	}
	return nil
}

func SaveTaskForUser(uid string, task *Task) error {
	db := data.DB
	task.UserID = uid
	_, err := db.NamedExec("INSERT INTO tasks VALUES (:id, :human_number, :user_id, :name, :subject, :status, :time_created.value)", task)
	if err != nil {
		return err
	}
	return nil
}
