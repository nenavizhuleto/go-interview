package company

import "github.com/google/uuid"

type NetworkID string

type Network struct {
	ID  NetworkID `json:"id"`
	Net string    `json:"netmask"`
}

func (n Network) PropertyName() string {
	return "network"
}

func NewNetwork(net string) Network {
	return Network{
		ID:  NetworkID(uuid.NewString()),
		Net: net,
	}
}

func (n *Network) SetNet(net string) {
	n.Net = net
}

func (n Network) Value() string {
	return n.Net
}
