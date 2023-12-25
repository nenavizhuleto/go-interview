package models

import (
	"errors"
	"fmt"
	"net"

	"helpdesk/internals/data"
)

type Network struct {
	Netmask  string `json:"netmask" db:"netmask"`
	BranchID string `json:"branch_id" db:"branch_id"`
}

var ErrUnsupportedDevice = errors.New("UnsupportedDevice")

func GetNetworkFromDevice(device *Device) (*Network, error) {
	db := data.DB

	var networks []Network
	if err := db.Select(&networks, "SELECT * FROM subnets ORDER by netmask"); err != nil {
		return nil, err
	}
	var network Network
	for _, n := range networks {
		_, ipnet, err := net.ParseCIDR(n.Netmask)
		size, _ := ipnet.Mask.Size()
		ip, _, err := net.ParseCIDR(device.IP + "/" + fmt.Sprint(size))
		if err != nil {
			return nil, err
		}
		if ipnet.Contains(ip) {
			network = n
			return &network, nil
		}
	}
	return nil, ErrUnsupportedDevice
}
