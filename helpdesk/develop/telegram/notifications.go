package telegram

import (
	"fmt"
	"helpdesk/internals/models/task"
	"helpdesk/internals/models/user"
)

const CommentUpdateMessage = `
Новый комментарий по обращению: #%s %s 
Комментарий:
%s
От: HelpDesk (Технический специалист)
Для ответа на комментарий перейдите в приложение helpdesk с вашего рабочего устройства
`

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
				msg += fmt.Sprintf(CommentUpdateMessage, tk.ID, tk.Name, tk.Comments[len(tk.Comments)-1].Content)
			}
		}

		t.SendMessageDirectlyToChat(user.ChatID, msg)
	}

}
