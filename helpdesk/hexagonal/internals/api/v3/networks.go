package api

import (
	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/branch"
	"helpdesk/internals/models/v3/network"

	"github.com/gofiber/fiber/v2"
)

func SetNetworksRoutes(path string, router fiber.Router) {
	networks := router.Group(path)

	networks.Get("/", GetNetworks)
	networks.Get("/:id", GetNetwork)
	networks.Post("/", CreateNetwork)
	networks.Delete("/:id", DeleteNetwork)
}

func GetNetworks(c *fiber.Ctx) error {
	networks, err := network.All()
	if err != nil {
		return err
	}

	return c.JSON(networks)
}

func GetNetwork(c *fiber.Ctx) error {
	id := c.Params("id")
	network, err := network.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(network)
}

func CreateNetwork(c *fiber.Ctx) error {
	var body struct {
		BranchID string `json:"branch_id"`
		Netmask  string
	}
	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("create_network", err)
	}

	branch, err := branch.Get(body.BranchID)
	if err != nil {
		return err
	}

	network, err := network.New(branch, body.Netmask)
	if err != nil {
		return err
	}

	if err := network.Save(); err != nil {
		return err
	}

	return c.JSON(network)
}

func DeleteNetwork(c *fiber.Ctx) error {
	id := c.Params("id")
	network, err := network.Get(id)
	if err != nil {
		return err
	}

	if err := network.Delete(); err != nil {
		return err
	}

	return c.JSON(network)
}
