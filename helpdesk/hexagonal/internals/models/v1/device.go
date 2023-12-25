package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"helpdesk/internals/data"
)

type Device struct {
	IP     string  `json:"ip" db:"ip"`
	User   *User   `json:"user_id" db:"user_id"`
	Subnet *Subnet `json:"subnet" db:"subnet_id"`
	Type   string  `json:"type" db:"type"`
}

const devices = "devices"

func DevicesAll() ([]Device, error) {
	coll := data.GetCollection(devices)

	cursor, err := coll.Find(nil, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	var devices []Device
	if err := cursor.All(nil, &devices); err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	return devices, nil
}

func (d *Device) Save() error {
	db := data.DB

	if _, err := db.NamedExec("UPDATE devices SET user_id = :user.id, subnet_id = :subnet.id, type = :type", d); err != nil {
		return err
	}

	return nil
}

func (d *Device) Exists() (bool, error) {
	coll := data.GetCollection(devices)

	count, err := coll.CountDocuments(nil, bson.D{{Key: "ip", Value: d.IP}})
	if err != nil {
		return false, fmt.Errorf("exists: %w", err)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (d *Device) Create() error {
	coll := data.GetCollection(devices)

	if d.IP == "" {
		return fmt.Errorf("create: ip is empty")
	}

	if exists, err := d.Exists(); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("create: already exists")
	}

	if _, err := coll.InsertOne(nil, d); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	return nil
}

func NewDevice(ip string) (*Device, error) {
	db := data.DB

	var device Device
	if err := db.Get(&device, "SELECT * FROM devices WHERE ip = $1", ip); err != nil {
		device = Device{
			IP:   ip,
			Type: "PC",
		}
		_, err := db.NamedExec("INSERT INTO devices(ip, type) VALUES (:ip, :type)", &device)
		if err != nil {
			return &device, err
		}
	}

	return &device, nil
}
