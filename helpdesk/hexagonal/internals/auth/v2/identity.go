package auth

import (
	"fmt"

	"helpdesk/internals/models/v2"
)

func GetIdentity(ip string) (*models.Device, error) {
	return models.DeviceGetByIP(ip)
}

func Register(user *models.User) error {
	if err := user.Create(); err != nil {
		return fmt.Errorf("identity: %w", err)
	}
	return nil
}

func MakeIdentity(ip string) (*models.Device, error) {
	dev, err := models.NewDevice(ip)
	if err != nil {
		return nil, fmt.Errorf("identity: %w", err)
	}

	if err := dev.Create(); err != nil {
		return nil, fmt.Errorf("identity: %w", err)
	}

	return dev, nil
}
