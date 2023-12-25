package mongo

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/task"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepo struct {
	name string
	c    *mongo.Collection
}

func NewTaskRepo(name string) TaskRepo {
	return TaskRepo{
		name: name,
		c:    GetCollection(name),
	}
}

func (r TaskRepo) Get(id task.TaskID) (task.Task, error) {
	var t task.Task
	if err := r.c.FindOne(nil, bson.D{{Key: "id", Value: id}}).Decode(&t); err != nil {
		return task.Task{}, core.NewDatabaseError(r.name, "get", err)
	}

	return t, nil
}

func (r TaskRepo) All() ([]task.Task, error) {
	cursor, err := r.c.Find(nil, bson.D{})
	if err != nil {
		return nil, core.NewDatabaseError(r.name, "all", err)
	}

	tasks := make([]task.Task, 0)
	if err := cursor.All(nil, &tasks); err != nil {
		return nil, core.NewDatabaseError(r.name, "all", err)
	}

	return tasks, nil
}

func (r TaskRepo) Save(t task.Task) error {
	if _, err := r.Get(t.ID); err != nil {
		if _, err := r.c.InsertOne(nil, t); err != nil {
			return core.NewDatabaseError(r.name, "create", err)
		}
	} else {
		// Exists
		if _, err := r.c.UpdateOne(nil, bson.D{{Key: "id", Value: t.ID}}, bson.D{{Key: "$set", Value: t}}); err != nil {
			return core.NewDatabaseError(r.name, "update")
		}
	}

	return nil
}

func (r TaskRepo) Delete(id task.TaskID) error {
	if err := r.c.FindOneAndDelete(nil, bson.D{{Key: "id", Value: id}}).Err(); err != nil {
		return core.NewDatabaseError(r.name, "delete", err)
	}

	return nil
}
