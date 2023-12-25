package models

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
)

type User struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Phone         string       `json:"phone"`
	Devices       []string     `json:"devices"`    // Will be used in the future, to manage multiple device access by user
	OnAfterCreate UserHookFunc `json:"-" bson:"-"` // Execute after user have been successfully inserted into database
}

type UserHookFunc func(*User) error

const users = "users"

func NewUser() *User {
	return &User{
		ID:      uuid.NewString(),
		Devices: make([]string, 0),
	}
}

func UserByID(id string) (*User, error) {
	coll := data.GetCollection(users)

	var user User
	if err := coll.FindOne(nil, bson.D{{Key: "id", Value: id}}).Decode(&user); err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (u *User) Create() error {

	coll := data.GetCollection(users)
	if _, err := coll.InsertOne(nil, u); err != nil {
		return fmt.Errorf("user: %w", err)
	}

	if u.OnAfterCreate != nil {
		if err := u.OnAfterCreate(u); err != nil {
			return fmt.Errorf("user_after_create_hook: %w", err)
		}
	}

	return nil
}

func UsersAll() ([]User, error) {
	collection := data.GetCollection(users)
	var users []User
	cursor, err := collection.Find(nil, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	return users, nil
}

// Validate returns nil if validation succeded
func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("validate: user's name is empty")
	}

	if u.ID == "" {
		return fmt.Errorf("validate: user's id is empty")
	}

	if u.Phone == "" {
		return fmt.Errorf("validate: user's phone is empty")
	}

	return nil
}

func (u *User) Exists() (bool, error) {
	collection := data.GetCollection(users)
	count, err := collection.CountDocuments(nil, bson.D{{Key: "phone", Value: u.Phone}})
	if err != nil {
		return false, fmt.Errorf("exists: %w", err)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (u *User) Fetch() error {
	collection := data.GetCollection(users)
	if u.ID == "" {
		return fmt.Errorf("fetch: id is empty")
	}

	res := collection.FindOne(nil, bson.D{{Key: "id", Value: u.ID}})
	if err := res.Decode(u); err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	return nil
}

func (u *User) Update() error {
	collection := data.GetCollection(users)
	if exists, err := u.Exists(); err != nil {
		return fmt.Errorf("update: %w", err)
	} else if !exists {
		return fmt.Errorf("update: not found")
	}
	if err := u.Validate(); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if _, err := collection.UpdateOne(nil, bson.D{{Key: "id", Value: u.ID}}, bson.D{{Key: "$set", Value: u}}); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (u *User) Delete() error {
	collection := data.GetCollection(users)

	if u.ID == "" {
		return fmt.Errorf("delete: user's id is empty")
	}

	if res, err := collection.DeleteOne(nil, bson.D{{Key: "id", Value: u.ID}}); err != nil {
		return fmt.Errorf("delete: %w", err)
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("delete: user not found")
	}

	return nil
}
