package models

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
)

type User struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Phone string `json:"phone" db:"phone"`
}

func (u *User) Value() string {
	return u.ID
}

func UsersAll() ([]User, error) {
	collection := data.GetCollection("users")
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
	collection := data.GetCollection("users")
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
	collection := data.GetCollection("users")
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
	collection := data.GetCollection("users")
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

func (u *User) Create() error {
	collection := data.GetCollection("users")

	if exists, err := u.Exists(); err != nil {
		return fmt.Errorf("create: %w", err)
	} else if exists {
		return fmt.Errorf("create: user with specified phone %s is found", u.Phone)
	}

	if err := u.Validate(); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if _, err := collection.InsertOne(nil, u); err != nil {
		return fmt.Errorf("create: %w", err)
	}
	return nil
}

func (u *User) Delete() error {
	collection := data.GetCollection("users")

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

func GetUserFromDevice(device *Device) (*User, error) {
	db := data.DB
	var user User
	if device.User != nil {
		if err := db.Get(&user, "SELECT * FROM users WHERE id = $1", device.User.ID); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, ErrUserNotFound
}

func NewUser(ip, name, phone string) (*User, error) {
	db := data.DB

	user := User{
		ID:    uuid.NewString(),
		Name:  name,
		Phone: phone,
	}

	_, err := db.NamedExec("INSERT INTO users VALUES (:id, :name, :phone)", &user)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("UPDATE devices SET user_id = $1 WHERE ip = $2", user.ID, ip)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
