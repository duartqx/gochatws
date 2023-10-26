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
	chatId int `db:"chat_id"`
	userId int `db:"user_id"`

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
	return mm.chatId
}

func (mm MessageModel) GetUserId() int {
	return mm.userId
}

func (mm MessageModel) GetCreatedAt() *time.Time {
	return &mm.CreatedAt
}

func (mm *MessageModel) SetUserId(id int) i.Message {
	mm.userId = id
	return mm
}

func (mm *MessageModel) SetChatId(id int) i.Message {
	mm.chatId = id
	return mm
}

func (mm *MessageModel) SetCreatedAt() i.Message {
	mm.CreatedAt = time.Now()
	return mm
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
		Id:        chat.GetId(),
		Category:  chat.GetCategory(),
		CreatorId: chat.GetCreatorId(),
		Name:      chat.GetName(),
	}
	return mm
}
