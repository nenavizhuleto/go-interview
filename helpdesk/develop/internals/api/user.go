package api

import (
	"helpdesk/internals/models"
	"helpdesk/internals/models/comment"
	"helpdesk/internals/models/task"
	"helpdesk/internals/models/user"

	"github.com/gofiber/fiber/v2"
)

func GetUserProfile(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	return c.JSON(Success(fiber.Map{
		"name":  u.Name,
		"phone": u.Phone,
		"company": fiber.Map{
			"name": u.Company.Name,
		},
		"branch": fiber.Map{
			"name":        u.Branch.Name,
			"description": u.Branch.Description,
			"address":     u.Branch.Address,
			"contacts":    u.Branch.Contacts,
		},
	}))
}

var (
	UserTaskFilterBranch  = "branch"
	UserTaskFilterCompany = "company"
)

func GetUserTasks(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	var tasks []task.Task
	var err error
	filter := c.Query("filter")
	switch filter {
	case UserTaskFilterBranch:
		tasks, err = task.ByBranch(u.Branch.ID)
	case UserTaskFilterCompany:
		tasks, err = task.ByCompany(u.Company.ID)
	default:
		tasks, err = task.ByUser(&u)
	}
	if err != nil {
		return err
	}

	res := make([]fiber.Map, 0)
	for _, t := range tasks {
		res = append(res, fiber.Map{
			"id":          t.ID,
			"name":        t.Name,
			"subject":     t.Subject,
			"status":      t.Status,
			"created_at":  t.TimeCreated,
			"activity_at": t.LastActivity,
			"branch": fiber.Map{
				"name":        t.Branch.Name,
				"description": t.Branch.Description,
				"address":     t.Branch.Address,
				"contacts":    t.Branch.Contacts,
			},
			"user": fiber.Map{
				"id":    t.User.ID,
				"name":  t.User.Name,
				"phone": t.User.Phone,
			},
		})
	}

	return c.JSON(Success(res))
}

func NewUserTelegram(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	tg, err := u.CreateTelegram()
	if err != nil {
		return err
	}

	return c.JSON(Success(fiber.Map{
		"pass": tg.Pass,
	}))
}

func GetUserTelegram(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	tg, err := u.GetTelegram()
	if err != nil {
		return err
	}

	return c.JSON(Success(fiber.Map{
		"pass": tg.Pass,
	}))
}

func GetUserTask(c *fiber.Ctx) error {
	task_id := c.Params("id")

	t, err := task.Get(task_id)
	if err != nil {
		return err
	}

	return c.JSON(Success(fiber.Map{
		"id":          t.ID,
		"name":        t.Name,
		"subject":     t.Subject,
		"status":      t.Status,
		"created_at":  t.TimeCreated,
		"activity_at": t.LastActivity,
		"branch": fiber.Map{
			"name":        t.Branch.Name,
			"description": t.Branch.Description,
			"address":     t.Branch.Address,
			"contacts":    t.Branch.Contacts,
		},
		"user": fiber.Map{
			"id":    t.User.ID,
			"name":  t.User.Name,
			"phone": t.User.Phone,
		},
	}))
}

func GetUserTaskComments(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)
	task_id := c.Params("id")

	t, err := task.Get(task_id)
	if err != nil {
		return err
	}

	comments := make([]fiber.Map, 0)
	for _, com := range t.Comments {
		direction := comment.DirectionTo
		if com.User != nil && com.User.ID == u.ID {
			direction = comment.DirectionFrom
		}
		comments = append(comments, fiber.Map{
			"id":      com.ID,
			"content": com.Content,
			"user": fiber.Map{
				"name":  com.User.Name,
				"phone": com.User.Phone,
			},
			"direction":  direction,
			"created_at": com.TimeCreated,
		})
	}

	return c.JSON(Success(comments))
}

func NewUserTaskComment(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)
	forward := c.Query("forward", "")
	if forward == "megaplan" {
		return NewUserTaskCommentMegaplan(c)
	}
	var body struct {
		Content string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("comment_task", err)
	}

	id := c.Params("id")
	t, err := task.Get(id)
	if err != nil {
		return err
	}

	comm := comment.NewUserComment(&u, body.Content)
	t.Comments = append(t.Comments, comm)

	if err := t.Save(); err != nil {
		return err
	}

	return c.JSON(Success(comm.ID))
}

func NewUserTask(c *fiber.Ctx) error {
	forward := c.Query("forward", "")
	if forward == "megaplan" {
		return NewUserTaskMegaplan(c)
	}
	u := c.Locals("user").(user.User)
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
	tk.User = &u
	tk.Branch = u.Branch
	tk.Company = u.Company

	if err := tk.Save(); err != nil {
		return err
	}

	return c.JSON(Success(tk.ID))
}
