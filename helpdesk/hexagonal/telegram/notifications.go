package telegram

import (
	"fmt"
	"helpdesk/internals/models/v3/task"
	"helpdesk/internals/models/v3/user"
)

func (t *TelegramNotificator) NotifyUser(user *user.TelegramUser, tk *task.Task, update task.UpdateEvent) {
	chat := t.chats.Get(user.ChatID)
	authorized, ok := chat.Data["authorized"]
	if !ok {
		return
	}

	if authorized.(bool) {
		msg := "Вам новое уведомление от Helpdesk\n"
		msg += fmt.Sprintf("Номер задачи: %s\n", tk.ID)
		for _, updatable := range update {
			switch updatable {
			case task.StatusUpdate:
				msg += fmt.Sprintf("Статус изменен на %s\n", tk.Status)
			case task.CommentUpdate:
				msg += fmt.Sprintf("Новый комментарий: \n%s\n", tk.Comments[len(tk.Comments)-1].Content)
			}
		}

		t.SendMessageDirectlyToChat(user.ChatID, msg)
	}

}
