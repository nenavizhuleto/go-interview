package api

import (
	"encoding/json"
	"fmt"
	"helpdesk/internals/megaplan"
	"helpdesk/internals/models"
	"helpdesk/internals/models/comment"
	"helpdesk/internals/models/task"
	"helpdesk/internals/models/user"
	"helpdesk/telegram"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ProcessTaskEvent(dto megaplan.TaskDTO) error {
	_task, err := task.Get(dto.ID)
	if err != nil {
		return err
	}

	updateTask := func(t *task.Task, u task.UpdateEvent) error {
		if dto.Activity != nil {
			_task.LastActivity = dto.Activity.Value
			u = append(u, task.ActivityUpdate)
		}
		log.Printf("Updating task %s", t.ID)
		if err := t.Save(); err != nil {
			return err
		}

		if tg, _ := t.User.GetTelegram(); tg != nil {
			telegram.Bot.NotifyUser(tg, t, u)
		}

		return nil
	}

	update := make(task.UpdateEvent, 0)

	// --- Did status changed? ---
	newStatus := dto.GetStatus()
	if _task.Status != newStatus {
		_task.Status = newStatus
		update = append(update, task.StatusUpdate)
		return updateTask(_task, update)
	}

	// --- Did new comment appeared? ---
	newComment := dto.GetLastComment()
	if newComment != nil && newComment.Direction == comment.DirectionTo {
		for _, c := range _task.Comments {
			if c.ID == newComment.ID {
				return nil
			}
		}
		_task.Comments = append(_task.Comments, *newComment)
		update = append(update, task.CommentUpdate)
		return updateTask(_task, update)
	}

	return nil
}

func HandleMegaplanEvent(c *fiber.Ctx) error {
	var event megaplan.TaskEvent
	if err := c.BodyParser(&event); err != nil {
		return c.SendStatus(500)
	}

	if event.Model != "Task" {
		return c.SendStatus(200)
	}

	str, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("Event: %s", string(str))

	dto := event.Data
	if err := ProcessTaskEvent(dto); err != nil {
		return c.SendStatus(200)
	}

	return c.SendStatus(200)
}

func NewUserTaskCommentMegaplan(c *fiber.Ctx) error {
	_user := c.Locals("user").(user.User)
	var body struct {
		Content   string
		Direction string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("comment_task", err)
	}

	id := c.Params("id")
	t, err := task.Get(id)
	if err != nil {
		return err
	}

	var content string
	switch body.Direction {
	case comment.DirectionTo:
		content = fmt.Sprintf("%s %s", megaplan.CommentTagTo, body.Content)
	case comment.DirectionFrom:
		content = fmt.Sprintf("%s %s", megaplan.CommentTagFrom, body.Content)
	default:
		content = fmt.Sprintf("%s %s", megaplan.CommentTagFrom, body.Content)
	}
	mp_comment, err := megaplan.MP.CommentTask(t.ID, content)
	if err != nil {
		return err
	}

	var _comment comment.Comment
	_comment.ID = mp_comment.ID
	_comment.Content = body.Content
	_comment.Direction = comment.DirectionFrom
	_comment.User = &_user

	t.Comments = append(t.Comments, _comment)

	if err := t.Save(); err != nil {
		return err
	}

	return c.JSON(Success(_comment))
}

func NewUserTaskMegaplan(c *fiber.Ctx) error {
	_user := c.Locals("user").(user.User)

	log.Printf("Body: %s", string(c.BodyRaw()))

	var body struct {
		Name    string
		Subject string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("task", err)
	}

	_task, err := task.New()
	if err != nil {
		return err
	}

	_task.Name = body.Name
	_task.Subject = body.Subject
	_task.User = &_user
	_task.Branch = _user.Branch
	_task.Company = _user.Company
	_task.BeforeSaveHook = PrepareTaskForMegaplan

	if err := _task.Save(); err != nil {
		return err
	}

	return c.JSON(Success(_task))

}

func PrepareTaskForMegaplan(_task *task.Task) error {
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
	task_name := fmt.Sprintf("%s: %s", _task.Company.Name, _task.Name)
	task_subject := fmt.Sprintf(TaskSubjectFormat,
		_task.User.Name,
		_task.Subject,
		_task.User.Phone,
		_task.User.Devices[0],
		_task.Branch.Name,
		_task.Branch.Description,
		_task.Branch.Address,
		_task.Branch.Contacts,
	)

	dto, err := megaplan.MP.CreateTask(task_name, task_subject)
	if err != nil {
		return fmt.Errorf("before_create_hook: %w", err)
	}

	if dto.TimeCreated != nil {
		_task.ID = dto.ID
		_task.TimeCreated = dto.TimeCreated.Value
		_task.LastActivity = dto.TimeCreated.Value
	}

	return nil

}
