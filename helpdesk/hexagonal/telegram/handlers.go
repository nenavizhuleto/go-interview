package telegram

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (t *TelegramNotificator) HandleCommand(update tg.Update) {
	if handler, ok := t.commands[update.Message.Command()]; ok {
		handler(update)
	} else {
		t.UnknownCommand(update)
	}
}

func (t *TelegramNotificator) HandleMessage(update tg.Update) {
	err := t.chats.HandleChat(update)
	if err != nil {
		t.SendError(update, err)
		return
	}
}
