package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TelegramNotificator) send(c tg.Chattable) (tg.Message, error) {
	return t.b.Send(c)
}
func (t *TelegramNotificator) SendMessageDirectlyToChat(chat_id int64, text string) (tg.Message, error) {
	msg := tg.NewMessage(chat_id, text)
	return t.b.Send(msg)
}

func (t *TelegramNotificator) SendMessage(update tg.Update, text string) (tg.Message, error) {
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	return t.b.Send(msg)
}

func (t *TelegramNotificator) SendError(update tg.Update, err error) (tg.Message, error) {
	msg := tg.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error occured: %s", err.Error()))
	return t.b.Send(msg)
}

func (t *TelegramNotificator) SendMarkdown(update tg.Update, markdown string, args ...any) (tg.Message, error) {

	text := markdown
	if len(args) > 0 {
		text = fmt.Sprintf(markdown, args...)
	}
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = tg.ModeMarkdown
	return t.b.Send(msg)
}
