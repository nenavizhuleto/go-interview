package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/task"
	"helpdesk/internals/models/v3/user"
)

func SetUsersRoutes(path string, router fiber.Router) {
	users := router.Group(path)
	users.Get("/", GetUsers)
	users.Get("/:id", GetUser)
	users.Post("/", CreateUser)
	users.Delete("/:id", DeleteUser)

	users.Get("/:id/tasks", GetUserTasks)
	users.Post("/:id/tasks", CreateUserTask)
	users.Post("/:id/telegram", CreateUserTelegram)
}

func GetUsers(c *fiber.Ctx) error {
	users, err := user.All()
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := user.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {

	var body struct {
		Name  string
		Phone string
		IP    string
	}
	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("createUser: %w", err)
	}
	user, err := user.New(body.Name, body.Phone, body.IP)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func CreateUserTelegram(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := user.Get(id)
	if err != nil {
		return err
	}

	if err := user.CreateTelegram(); err != nil {
		return err
	}

	tg, err := user.GetTelegram()
	if err != nil {
		return err
	}

	return c.JSON(tg)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := user.Get(id)
	if err != nil {
		return err
	}

	if err := user.Delete(); err != nil {
		return err
	}

	return c.JSON(user)
}

func GetUserTasks(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := user.Get(id)
	if err != nil {
		return err
	}

	tasks, err := task.ByUser(user)
	fmt.Println(tasks)
	if err != nil {
		return err
	}

	return c.JSON(tasks)

}

func CreateUserTask(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := user.Get(id)
	if err != nil {
		return err
	}

	var body struct {
		Name    string
		Subject string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("task", err)
	}

	tk, err := task.New()
	if err != nil {
		return err
	}

	tk.Name = body.Name
	tk.Subject = body.Subject
	tk.User = user
	tk.Branch = user.Branch
	tk.Company = user.Company

	if err := tk.Save(); err != nil {
		return err
	}

	return c.JSON(tk)

}
