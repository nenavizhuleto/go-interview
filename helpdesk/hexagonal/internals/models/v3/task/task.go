package task

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/branch"
	"helpdesk/internals/models/v3/company"
	"helpdesk/internals/models/v3/user"
	"helpdesk/internals/models/v3/comment"
)

type updatable int
const (
	StatusUpdate updatable = iota // type string
	CommentUpdate // type models.Comment
	ActivityUpdate // type time.Time
)

type UpdateEvent []updatable

type Task struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Subject          string            `json:"subject"`
	Status           string            `json:"status"`
	TimeCreated      time.Time         `json:"created_at"`
	LastActivity     time.Time         `json:"activity_at"`
	Company          *company.Company  `json:"company"`
	Branch           *branch.Branch    `json:"branch"`
	User             *user.User        `json:"user"`
	Comments         []comment.Comment `json:"comments"`
	BeforeCreateHook HookFunc          `json:"-" bson:"-"`
}

const tasks = "tasks"

type HookFunc func(*Task) error

func New() (*Task, error) {
	tm := time.Now()
	return &Task{
		ID:           uuid.NewString(),
		Status:       "created",
		TimeCreated:  tm,
		LastActivity: tm,
		Comments:     make([]comment.Comment, 0),
	}, nil
}

func Get(id string) (*Task, error) {
	coll := data.GetCollection(tasks)
	var task Task
	if err := coll.FindOne(nil, bson.D{{Key: "id", Value: id}}).Decode(&task); err != nil {
		return nil, models.NewDatabaseError("task", "get", err)
	}

	return &task, nil
}

func ByUser(user *user.User) ([]Task, error) {
	coll := data.GetCollection(tasks)
	cursor, err := coll.Find(nil, bson.D{{Key: "user.id", Value: user.ID}})
	if err != nil {
		return nil, models.NewDatabaseError("task", "by_user", err)
	}

	tasks := make([]Task, 0)
	if err := cursor.All(nil, &tasks); err != nil {
		return nil, models.NewDatabaseError("task", "by_user", err)
	}

	return tasks, nil
}

func All() ([]Task, error) {
	coll := data.GetCollection(tasks)
	cursor, err := coll.Find(nil, bson.D{})
	if err != nil {
		return nil, models.NewDatabaseError("task", "all", err)
	}

	tasks := make([]Task, 0)
	if err := cursor.All(nil, &tasks); err != nil {
		return nil, models.NewDatabaseError("task", "all", err)
	}

	return tasks, nil
}

func (t *Task) Save() error {
	coll := data.GetCollection(tasks)

	if _, err := Get(t.ID); err != nil {
		// Not exists
		if t.BeforeCreateHook != nil {
			if err := t.BeforeCreateHook(t); err != nil {
				return models.NewDatabaseError("task", "before_create", err)
			}
		}

		if _, err := coll.InsertOne(nil, t); err != nil {
			return models.NewDatabaseError("task", "create", err)
		}
	} else {
		// Exists
		if _, err := coll.UpdateOne(nil, bson.D{{Key: "id", Value: t.ID}}, bson.D{{Key: "$set", Value: t}}); err != nil {
			return models.NewDatabaseError("task", "update")
		}
	}

	return nil
}

func (t *Task) Delete() error {
	coll := data.GetCollection(tasks)
	if err := coll.FindOneAndDelete(nil, bson.D{{Key: "id", Value: t.ID}}).Err(); err != nil {
		return models.NewDatabaseError("task", "delete", err)
	}
	return nil
}
