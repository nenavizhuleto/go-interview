package user

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
	"helpdesk/internals/models"
)

type User struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Phone         string   `json:"phone"`
	Devices       []string `json:"devices"`    // Will be used in the future, to manage multiple device access by user
	OnAfterCreate HookFunc `json:"-" bson:"-"` // Execute after user have been successfully inserted into database
	OnAfterUpdate HookFunc `json:"-" bson:"-"` // Execute after user have been successfully inserted into database
}

type HookFunc func(*User) error

const users = "users"

func New(username string, phone string) (*User, error) {
	validName, err := newName(username)
	if err != nil {
		return nil, err
	}
	validPhone, err := newPhone(phone)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:      newID(),
		Name:    validName,
		Phone:   validPhone,
		Devices: make([]string, 0),
	}, nil
}

func Get(id string) (*User, error) {
	coll := data.GetCollection(users)

	var user User
	if err := coll.FindOne(nil, bson.D{{Key: "id", Value: id}}).Decode(&user); err != nil {
		return nil, models.NewDatabaseError("user", "get", err)
	}

	return &user, nil
}

func All() ([]User, error) {
	coll := data.GetCollection(users)
	var users []User
	cursor, err := coll.Find(nil, bson.D{})
	if err != nil {
		return nil, models.NewDatabaseError("user", "all", err)
	}

	if err = cursor.All(nil, &users); err != nil {
		return nil, models.NewDatabaseError("user", "all", err)
	}

	return users, nil
}

func (u *User) Save() error {
	coll := data.GetCollection(users)
	if _, err := Get(u.ID); err != nil {
		// Insert if not exists
		if _, err := coll.InsertOne(nil, u); err != nil {
			return models.NewDatabaseError("user", "create", err)
		}

		if u.OnAfterCreate != nil {
			if err := u.OnAfterCreate(u); err != nil {
				return models.NewDatabaseError("user", "create_hook", err)
			}
		}
	} else {
		// Update if exists
		if err := coll.FindOneAndUpdate(nil, bson.D{{Key: "id", Value: u.ID}}, bson.D{{Key: "$set", Value: u}}).Err(); err != nil {
			return models.NewDatabaseError("user", "update", err)
		}

		if u.OnAfterUpdate != nil {
			if err := u.OnAfterUpdate(u); err != nil {
				return models.NewDatabaseError("user", "update_hook", err)
			}
		}
	}

	return nil
}

func (u *User) Delete() error {
	coll := data.GetCollection(users)
	if err := coll.FindOneAndDelete(nil, bson.D{{Key: "id", Value: u.ID}}).Err(); err != nil {
		return models.NewDatabaseError("user", "delete", err)
	}

	return nil
}

// Private functions

func newID() string {
	return uuid.NewString()
}

func newName(username string) (string, error) {
	// Domain logic goes here...
	// Validation etc...
	if len(username) == 0 {
		return "", models.NewValidationError("user", "username")
	}
	return username, nil
}

func newPhone(phone string) (string, error) {
	// Validation etc...
	if len(phone) == 0 {
		return "", models.NewValidationError("user", "phone")
	}
	return phone, nil
}

func newDevice(ip string) string {
	return ip
}
