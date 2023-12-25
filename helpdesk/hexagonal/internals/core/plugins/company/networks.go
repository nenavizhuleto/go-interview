package company

import "net"

type BranchNetworkProperty net.IP

func NewBranchNetworkProperty(ip string) BranchNetworkProperty {
	return BranchNetworkProperty(net.ParseIP(ip))
}

func (e BranchNetworkProperty) GetPropertyName() BranchPropertyName {
	return BranchPropertyName("network")
}
