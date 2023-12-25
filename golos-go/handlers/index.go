package handlers

import "github.com/gofiber/fiber/v2"

func HandleMain(c *fiber.Ctx) error {
	return c.Render("screens/main", nil, "layout/system")
}

func HandleGetInfo(c *fiber.Ctx) error {
	return c.Render("screens/info", nil)
}

func HandleGetDevices(c *fiber.Ctx) error {
	return c.Render("screens/devices", nil)
}

func HandleGetRooms(c *fiber.Ctx) error {
	return c.Render("screens/rooms", nil)
}

func HandleGetPorches(c *fiber.Ctx) error {
	return c.Render("screens/porches", nil)
}

func HandleGetBuildings(c *fiber.Ctx) error {
	return c.Render("screens/buildings", nil)
}

func HandleGetBlocks(c *fiber.Ctx) error {
	return c.Render("screens/blocks", nil)
}

func HandleGetStreets(c *fiber.Ctx) error {
	return c.Render("screens/streets", nil)
}

func HandleGetConfig(c *fiber.Ctx) error {
	return c.Render("screens/config", nil)
}

func HandleIndex(c *fiber.Ctx) error {
	return c.Redirect("/system")
}
