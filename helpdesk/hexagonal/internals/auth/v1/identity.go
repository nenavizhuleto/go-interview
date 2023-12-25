package auth

import (
	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/models/v1"
)

type Identity struct {
	User struct {
		ID    string
		Name  string
		Phone string
	}
	Device struct {
		IP   string
		Type string
	}
	Subnet struct {
		Network string
	}
	Branch struct {
		Name        string
		Address     string
		Contacts    string
		Description string
	}
	Company struct {
		ID   string
		Name string
		Slug string
	}
}

func GetIdentity(c *fiber.Ctx) *Identity {
	i := c.Locals("Identity")
	if i.(*Identity) == nil {
		return nil
	}
	return i.(*Identity)
}

func MakeIdentity(ip string) (*Identity, error) {
	i := &Identity{}

	d, err := models.NewDevice(ip)
	if err != nil {
		return nil, err
	}

	i.Device.IP = d.IP
	i.Device.Type = d.Type

	s, err := models.GetSubnetFromDevice(d)
	if err != nil {
		return nil, err
	}

	i.Subnet.Network = s.Netmask

	u, err := models.GetUserFromDevice(d)
	if err != nil {
		return nil, err
	}

	i.User.ID = u.ID
	i.User.Name = u.Name
	i.User.Phone = u.Phone

	b, err := models.GetBranchFromSubnet(s)
	if err != nil {
		return nil, err
	}

	i.Branch.Name = b.Name
	i.Branch.Address = b.Address
	i.Branch.Contacts = b.Contacts
	i.Branch.Description = b.Description

	c, err := models.GetCompanyFromBranch(b)
	if err != nil {
		return nil, err
	}

	i.Company.ID = c.ID
	i.Company.Name = c.Name
	i.Company.Slug = c.Slug

	return i, nil
}
