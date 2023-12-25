package user

import (
	"net"
)

type DeviceType string

const (
	TypePC      = DeviceType("PC")
	TypeUnknown = DeviceType("unknown")
)

type Device struct {
	IP   net.IPAddr `json:"ip"`
	Type DeviceType `json:"type"`
}

func (d Device) GetPropertyName() UserPropertyName {
	return UserPropertyName("device")
}
