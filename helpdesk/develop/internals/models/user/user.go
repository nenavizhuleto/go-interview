package user

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
	"helpdesk/internals/models"
	"helpdesk/internals/models/branch"
	"helpdesk/internals/models/company"
	"helpdesk/internals/models/device"
	"helpdesk/internals/models/network"
	"helpdesk/internals/util"
)

type User struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Phone         string           `json:"phone"`
	Network       *network.Network `json:"network"`
	Branch        *branch.Branch   `json:"branch"`
	Company       *company.Company `json:"company"`
	Devices       []*device.Device `json:"devices"` // Will be used in the future, to manage multiple device access by user
	Identified    bool             `json:"identified"`
	OnAfterCreate HookFunc         `json:"-" bson:"-"` // Execute after user have been successfully inserted into database
	OnAfterUpdate HookFunc         `json:"-" bson:"-"` // Execute after user have been successfully inserted into database
}

type HookFunc func(*User) error

const users = "users"
const telegram = "telegram"

func New(username string, phone string, ip string) (*User, error) {
	network, err := network.GetByIP(ip)
	if err != nil {
		return nil, err
	}

	branch, err := branch.Get(network.BranchID)
	if err != nil {
		return nil, err
	}

	company, err := company.Get(branch.CompanyID)
	if err != nil {
		return nil, err
	}

	validName, err := newName(username)
	if err != nil {
		return nil, err
	}

	validPhone, err := newPhone(phone, company)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:         newID(),
		Name:       validName,
		Phone:      validPhone,
		Network:    network,
		Branch:     branch,
		Company:    company,
		Devices:    make([]*device.Device, 0),
		Identified: true,
	}, nil
}

func (u User) CreateTelegram() (*TelegramUser, error) {
	coll := data.GetCollection(telegram)

	tg := TelegramUser{
		User: u,
		Pass: util.RandStringBytes(TelegramPassLength),
	}

	if _, err := coll.InsertOne(nil, tg); err != nil {
		return nil, models.NewDatabaseError("user", "telegram", err)
	}

	return &tg, nil
}

func (u User) GetTelegram() (*TelegramUser, error) {
	coll := data.GetCollection(telegram)

	var tg TelegramUser
	if err := coll.FindOne(nil, bson.M{"user.id": u.ID}).Decode(&tg); err != nil {
		return nil, models.NewDatabaseError("user", "get_telegram", err)
	}

	return &tg, nil
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

func (u *User) Identify(ip string) error {
	network, err := network.GetByIP(ip)
	if err != nil {
		return err
	}

	branch, err := branch.Get(network.BranchID)
	if err != nil {
		return err
	}

	company, err := company.Get(branch.CompanyID)
	if err != nil {
		return err
	}

	u.Network = network
	u.Branch = branch
	u.Company = company

	u.Identified = true

	return nil
}

func (u *User) AddDevice(dev *device.Device) error {
	if !u.Identified {
		return models.NewIdentificationError("user", nil)
	}

	dev.OwnerID = u.ID
	if err := dev.Save(); err != nil {
		return err
	}

	u.Devices = append(u.Devices, dev)
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

func newPhone(phone string, c *company.Company) (string, error) {
	// Validation etc...
	if len(phone) == 0 {
		return "", models.NewValidationError("user", "phone")
	}

	coll := data.GetCollection(users)
	users, err := coll.CountDocuments(nil, bson.M{"company.id": c.ID, "phone": phone})
	if err != nil {
		return "", models.NewValidationError("user", "phone")
	}

	if users != 0 {
		return "", models.NewValidationError("user", "phone", "user exists")
	}

	return phone, nil
}

func newDevice(ip string) string {
	return ip
}
