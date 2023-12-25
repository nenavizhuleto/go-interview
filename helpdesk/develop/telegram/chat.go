package telegram

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type State int

const (
	ChatNormal State = iota
	ChatReceive
)



type ChatContext struct {
	State     State
	Data      map[string]interface{}
	OnReceive OnReceiveFunc
}

type OnReceiveFunc func(tg.Update) error

func NewChatContext() *ChatContext {
	return &ChatContext{
		State: ChatNormal,
		Data:  make(map[string]interface{}),
	}
}

type Chats map[int64]*ChatContext

func (c Chats) HandleChat(update tg.Update) error {
	if chat, ok := c[update.Message.Chat.ID]; ok {
		switch chat.State {
		case ChatReceive:
			var err error
			// Check if callback exists
			if chat.OnReceive != nil {
				err = chat.OnReceive(update)
				// Reset OnReceive callback
				chat.OnReceive = nil
			}

			// Reset chat's state
			chat.State = ChatNormal
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c Chats) Get(chatID int64) *ChatContext {
	chat, ok := c[chatID]
	if ok {
		return chat
	} else {
		c[chatID] = NewChatContext()
		return c[chatID]
	}
}

// ChangeState changes state of the chat, if not exists creates and updates it
func (c Chats) ChangeState(chatID int64, state State) *ChatContext {
	chat := c.Get(chatID)
	chat.State = ChatReceive
	return chat
}
