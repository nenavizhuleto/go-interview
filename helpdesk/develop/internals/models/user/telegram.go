package user

import (
	"helpdesk/internals/data"
	"helpdesk/internals/models"

	"go.mongodb.org/mongo-driver/bson"
)

const TelegramPassLength = 7

type TelegramUser struct {
	User   User   `json:"user"`
	Pass   string `json:"pass"`
	ChatID int64  `json:"chat_id"`
}

func VerifyTelegramPass(pass string) (*TelegramUser, error) {
	coll := data.GetCollection(telegram)

	var tg TelegramUser
	if err := coll.FindOne(nil, bson.M{"pass": pass}).Decode(&tg); err != nil {
		return nil, models.NewDatabaseError("user", "verify_telegram", err)
	}

	return &tg, nil
}

func (tu *TelegramUser) ConnectTelegramChat(chat_id int64) error {
	coll := data.GetCollection(telegram)

	tu.ChatID = chat_id

	if err := coll.FindOneAndUpdate(nil, bson.M{"user.id": tu.User.ID}, bson.M{"$set": tu}).Err(); err != nil {
		return models.NewDatabaseError("user", "verify_telegram", err)
	}

	return nil
}
