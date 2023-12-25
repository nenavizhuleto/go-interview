package models

import (
	"fmt"
	"net"

	"helpdesk/internals/data"
)

type Subnet struct {
	Netmask  string `json:"netmask" db:"netmask"`
	BranchID string `json:"branch_id" db:"branch_id"`
}

func GetSubnetFromDevice(device *Device) (*Subnet, error) {
	db := data.DB

	var subnets []Subnet
	if err := db.Select(&subnets, "SELECT * FROM subnets ORDER by netmask"); err != nil {
		return nil, err
	}
	var subnet Subnet
	for _, s := range subnets {
		_, network, err := net.ParseCIDR(s.Netmask)
		size, _ := network.Mask.Size()
		ip, _, err := net.ParseCIDR(device.IP + "/" + fmt.Sprint(size))
		if err != nil {
			return nil, err
		}
		if network.Contains(ip) {
			subnet = s
			return &subnet, nil
		}
	}
	return nil, ErrUnsupportedDevice
}
