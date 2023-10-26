package models

import (
	"time"

	i "github.com/duartqx/gochatws/core/interfaces"
)

type MessageModel struct {
	Id        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Text      string    `db:"text" json:"text"`

	// ForeignKeys
	ChatId int `db:"chat_id"`
	UserId int `db:"user_id"`

	// Structs not part of the db
	C i.ChatRoom `json:"chat"`
	U i.User     `json:"user"`
}

func (mm MessageModel) GetId() int {
	return mm.Id
}

func (mm *MessageModel) SetId(id int) i.Message {
	mm.Id = id
	return mm
}

func (mm MessageModel) GetChatId() int {
	return mm.ChatId
}

func (mm MessageModel) GetUserId() int {
	return mm.UserId
}

func (mm MessageModel) GetText() string {
	return mm.Text
}

func (mm MessageModel) GetUser() i.User {
	return mm.U
}

func (mm MessageModel) GetChat() i.ChatRoom {
	return mm.C
}

func (mm *MessageModel) PopulateUser(user i.User) i.Message {
	mm.U = &UserClean{
		Id:       user.GetId(),
		Username: user.GetUsername(),
		Name:     user.GetName(),
	}
	return mm
}

func (mm *MessageModel) PopulateChat(chat i.ChatRoom) i.Message {
	mm.C = &ChatRoomModel{
		Id:       chat.GetId(),
		Category: chat.GetCategory(),
		Name:     chat.GetName(),
	}
	return mm
}
