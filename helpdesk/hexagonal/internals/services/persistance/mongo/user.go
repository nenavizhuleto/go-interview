package mongo

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	name string
	c    *mongo.Collection
}

func NewUserRepo(name string) UserRepo {

	r := UserRepo{}
	r.name = name
	r.c = GetCollection(r.name)

	return r
}

func (r UserRepo) Get(id user.UserID) (user.User, error) {
	var u user.User
	if err := r.c.FindOne(nil, bson.D{{Key: "id", Value: id}}).Decode(&u); err != nil {
		return user.User{}, core.NewDatabaseError(r.name, "get", err)
	}
	return u, nil
}

func (r UserRepo) All() ([]user.User, error) {
	var users []user.User
	cursor, err := r.c.Find(nil, bson.D{})
	if err != nil {
		return nil, core.NewDatabaseError(r.name, "all", err)
	}

	if err = cursor.All(nil, &users); err != nil {
		return nil, core.NewDatabaseError(r.name, "all", err)
	}

	return users, nil
}

func (r UserRepo) Save(u user.User) error {
	if _, err := r.Get(u.ID); err != nil {
		// Insert if not exists
		if _, err := r.c.InsertOne(nil, u); err != nil {
			return core.NewDatabaseError(r.name, "create", err)
		}
	} else {
		// Update if exists
		if err := r.c.FindOneAndUpdate(nil, bson.D{{Key: "id", Value: u.ID}}, bson.D{{Key: "$set", Value: u}}).Err(); err != nil {
			return core.NewDatabaseError(r.name, "update", err)
		}
	}

	return nil
}

func (r UserRepo) Delete(id user.UserID) error {
	if err := r.c.FindOneAndDelete(nil, bson.D{{Key: "id", Value: id}}).Err(); err != nil {
		return core.NewDatabaseError(r.name, "delete", err)
	}

	return nil
}
