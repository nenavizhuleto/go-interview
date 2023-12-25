package telegram

import (
	"fmt"
	"helpdesk/internals/models/v3/user"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	WelcomeUnathorized = `
*Добро пожаловать в HelpDesk: Уведомления*

Здесь вы сможете получать уведомления об изменениях в системе в реальном времени

/login - Авторизоваться в системе
	`
	WelcomeAuthorized = `
*Добро пожаловать в HelpDesk: Уведомления*

Вы авторизировались как %s
Уведомления будут приходить в реальном времени в этот чат.

	`
)

func (t *TelegramNotificator) StartCommand(update tg.Update) {
	chat := t.chats.ChangeState(update.Message.Chat.ID, ChatNormal)
	if auth, ok := chat.Data["authorized"]; ok {
		if auth.(bool) {
			t.SendMarkdown(update, WelcomeAuthorized, chat.Data["user"].(*user.User).Name)
			return
		}
	}
	t.SendMarkdown(update, WelcomeUnathorized)
}

func (t *TelegramNotificator) HelpCommand(update tg.Update) {
	t.SendMessage(update, "This is help command")

}

func (t *TelegramNotificator) InfoCommand(update tg.Update) {
	chat := t.chats.Get(update.Message.Chat.ID)
	t.SendMessage(update, fmt.Sprintf("Information: %#v", chat))
}

const (
	Login = `
Введите код подтверждения в следующем сообщении
	`
	LoginInvalid = `
Код подтверждения неверный

/start - Вернуться в главное меню
/login - Попробовать снова
	`
)

func (t *TelegramNotificator) LoginCommand(update tg.Update) {
	chat := t.chats.ChangeState(update.Message.Chat.ID, ChatReceive)

	// Here we get received message
	chat.OnReceive = func(u tg.Update) error {
		// Verify user's authorization code

		user, err := user.VerifyTelegramPass(u.Message.Text)
		if err != nil {
			t.SendMarkdown(u, LoginInvalid)
			chat.Data["authorized"] = false
			return nil
		}
		chat.Data["authorized"] = true
		chat.Data["user"] = user.User
		user.ConnectTelegramChat(update.Message.Chat.ID)
		t.SendMarkdown(u, WelcomeAuthorized, update.Message.From.UserName)

		return nil
	}

	t.SendMessage(update, Login)
}

func (t *TelegramNotificator) UnknownCommand(update tg.Update) {
	t.SendMessage(update, "This is unknown command")

}
