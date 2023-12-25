package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
)

type Task struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Subject          string           `json:"subject"`
	Status           string           `json:"status"`
	TimeCreated      time.Time        `json:"created_at"`
	LastActivity     time.Time        `json:"activity_at"`
	Company          Company          `json:"company"`
	Branch           Branch           `json:"branch"`
	User             User             `json:"user"`
	Comments         []Comment        `json:"comments"`
	BeforeCreateHook BeforeCreateHook `json:"-" bson:"-"`
}

const tasks = "tasks"

type BeforeCreateHook func(*Task) error

func NewTask() *Task {
	tm := time.Now()
	return &Task{
		ID:           uuid.NewString(),
		Status:       "created",
		TimeCreated:  tm,
		LastActivity: tm,
		Comments:     make([]Comment, 0),
	}
}

func TaskByID(id string) (*Task, error) {
	coll := data.GetCollection(tasks)

	var task Task
	if err := coll.FindOne(nil, bson.D{{Key: "id", Value: id}}).Decode(&task); err != nil {
		return nil, fmt.Errorf("task: %s id not found. %w", id, err)
	}

	return &task, nil
}

func TasksByUser(user *User) ([]Task, error) {
	coll := data.GetCollection(tasks)

	cursor, err := coll.Find(nil, bson.D{{Key: "user.id", Value: user.ID}})
	if err != nil {
		return nil, fmt.Errorf("tasks: %w", err)
	}

	var tasks []Task
	if err := cursor.All(nil, &tasks); err != nil {
		return nil, fmt.Errorf("tasks: %w", err)
	}

	return tasks, nil
}

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

func (t *Task) Create() error {
	if t.BeforeCreateHook != nil {
		if err := t.BeforeCreateHook(t); err != nil {
			return err
		}
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
