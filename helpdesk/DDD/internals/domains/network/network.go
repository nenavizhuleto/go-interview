package network

import (
	"errors"
	"fmt"
	"helpdesk/internals/data"
	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/branch"
	"net"
)

var (
	ErrUnknownIP = errors.New("unknown ip address")
)

type Network struct {
	BranchID string `json:"branch_id" db:"branch_id"`
	Netmask  string `json:"netmask" db:"netmask"`
}

func New(branch *branch.Branch, netmask string) (*Network, error) {
	branch_id := branch.ID
	validNetmask, err := newNetmask(netmask)
	if err != nil {
		return nil, models.NewValidationError("network", "netmask")
	}

	return &Network{
		BranchID: branch_id,
		Netmask:  validNetmask,
	}, nil
}

func Get(netmask string) (*Network, error) {
	db := data.DB

	var network Network
	if err := db.Get(&network, "SELECT * FROM networks WHERE netmask = ?", netmask); err != nil {
		return nil, models.NewDatabaseError("network", "get", err)
	}

	return &network, nil
}

func GetByIP(ip string) (*Network, error) {
	networks, err := All()
	if err != nil {
		return nil, err
	}

	var network Network
	for _, n := range networks {
		_, ipnet, err := net.ParseCIDR(n.Netmask)
		size, _ := ipnet.Mask.Size()
		ip, _, err := net.ParseCIDR(ip + "/" + fmt.Sprint(size))
		if err != nil {
			return nil, models.NewDatabaseError("network", "get_by_ip", err)
		}
		if ipnet.Contains(ip) {
			network = n
		}
	}

	if network == (Network{}) {
		return nil, models.NewDatabaseError("network", "get_by_ip", ErrUnknownIP)
	}

	return &network, nil
}

func All() ([]Network, error) {
	db := data.DB

	var networks []Network
	if err := db.Select(&networks, "SELECT * FROM networks"); err != nil {
		return nil, models.NewDatabaseError("network", "all", err)
	}

	return networks, nil
}

func (n *Network) Save() error {
	db := data.DB

	if _, err := Get(n.Netmask); err != nil {
		// Not exists
		if _, err := db.NamedExec("INSERT INTO networks VALUES (:netmask, :branch_id)", n); err != nil {
			return models.NewDatabaseError("network", "create", err)
		}
	} else {
		// Exists
		if _, err := db.Exec("UPDATE networks SET netmask = :netmask WHERE branch_id = ? AND netmask = ?", n.BranchID, n.Netmask); err != nil {
			return models.NewDatabaseError("network", "update", err)
		}
	}

	return nil
}

func (n *Network) Delete() error {
	db := data.DB

	if _, err := db.Exec("DELETE FROM networks WHERE branch_id = ? AND netmask = ?", n.BranchID, n.Netmask); err != nil {
		return models.NewDatabaseError("network", "delete", err)
	}

	return nil
}

// Private functions

func newNetmask(netmask string) (string, error) {
	// Validate netmask
	if _, _, err := net.ParseCIDR(netmask); err != nil {
		return "", models.NewValidationError("network", "netmask", err.Error())
	}

	return netmask, nil
}
