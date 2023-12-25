package api

import (
	"encoding/json"
	"fmt"
	"helpdesk/internals/megaplan"
	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/comment"
	"helpdesk/internals/models/v3/task"
	tk "helpdesk/internals/models/v3/task"
	"helpdesk/internals/models/v3/user"
	"helpdesk/telegram"
	"log"

	"github.com/gofiber/fiber/v2"
)

func HandleMegaplanEvent(c *fiber.Ctx) error {
	var event megaplan.TaskEvent
	if err := c.BodyParser(&event); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	if event.Model != "Task" {
		return c.SendStatus(200)
	}

	str, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("Event: %s", string(str))

	dto := event.Data
	task, err := tk.Get(dto.ID)
	if err != nil {
		log.Printf("event: %s", err.Error())
		return c.SendStatus(200)
	}

	needUpdate := false
	update := make(tk.UpdateEvent, 0)
	newStatus := dto.GetStatus()
	if task.Status != newStatus {
		task.Status = newStatus
		update = append(update, tk.StatusUpdate)
		needUpdate = true
	}
	newComment := dto.GetLastComment()
	if newComment != nil {
		task.Comments = append(task.Comments, *newComment)
		update = append(update, tk.CommentUpdate)
		needUpdate = true
	}
	if dto.Activity != nil && needUpdate {
		task.LastActivity = dto.Activity.Value
		update = append(update, tk.ActivityUpdate)
	}

	if needUpdate {
		if err := task.Save(); err != nil {
			log.Println(err)
			return c.SendStatus(500)
		}

		if tg, _ := task.User.GetTelegram(); tg != nil {
			telegram.Bot.NotifyUser(tg, task, update)
		}
	}

	return c.SendStatus(200)
}

func CommentTaskMegaplan(c *fiber.Ctx) error {
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

	content := fmt.Sprintf("#[FROMUSER]: %s", body.Content)
	com, err := megaplan.MP.CommentTask(t.ID, content)
	if err != nil {
		return err
	}

	var comm comment.Comment
	comm.ID = com.ID
	comm.Content = com.Content
	comm.Direction = comment.DirectionFrom

	t.Comments = append(t.Comments, comm)

	return c.JSON(comm)
}

func CreateUserTaskMegaplan(c *fiber.Ctx) error {
	id := c.Params("id")

	log.Printf("Body: %s", string(c.BodyRaw()))

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

	tk.BeforeCreateHook = func(t *task.Task) error {
		var TaskSubjectFormat = `
			<h2>от %s:</h2>
			<h3>Суть обращения:</h3>
			<p>%s</p>
			<hr/>
			<h3>Дополнительная информания:</h3>
			<ul>
			<li>Контакты: %s</li>
			<li>Устройство: %s</li>
			<li>Отдел: <br/>Название: %s <br/>Описание: %s <br/>Адрес: %s <br/>Контакты: %s</li>
			</ul>
		`
		task_name := fmt.Sprintf("%s: %s", t.Company.Name, t.Name)
		task_subject := fmt.Sprintf(TaskSubjectFormat,
			t.User.Name,
			t.Subject,
			t.User.Phone,
			t.User.Devices[0],
			t.Branch.Name,
			t.Branch.Description,
			t.Branch.Address,
			t.Branch.Contacts,
		)

		dto, err := megaplan.MP.CreateTask(task_name, task_subject)
		if err != nil {
			return fmt.Errorf("before_create_hook: %w", err)
		}

		if dto.TimeCreated != nil {
			t.ID = dto.ID
			t.TimeCreated = dto.TimeCreated.Value
			t.LastActivity = dto.TimeCreated.Value
		}

		return nil
	}

	if err := tk.Save(); err != nil {
		return err
	}

	return c.JSON(tk)

}
