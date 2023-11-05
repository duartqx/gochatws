package message

import (
	"time"

	c "github.com/duartqx/gochatws/domains/entities/chatroom"
	u "github.com/duartqx/gochatws/domains/entities/user"
)

type MessageModel struct {
	Id        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Text      string    `db:"text" json:"text"`

	// ForeignKeys
	chatId int `db:"chat_id"`
	userId int `db:"user_id"`

	// Structs not part of the db
	C c.ChatRoom `json:"chat"`
	U u.User     `json:"user"`
}

func (mm MessageModel) GetChat() c.ChatRoom {
	return mm.C
}

func (mm MessageModel) GetChatId() int {
	return mm.chatId
}

func (mm MessageModel) GetCreatedAt() *time.Time {
	return &mm.CreatedAt
}

func (mm MessageModel) GetId() int {
	return mm.Id
}

func (mm MessageModel) GetText() string {
	return mm.Text
}

func (mm MessageModel) GetUser() u.User {
	return mm.U
}

func (mm MessageModel) GetUserId() int {
	return mm.userId
}

func (mm *MessageModel) SetChat(chat c.ChatRoom) Message {
	mm.chatId = chat.GetId()
	mm.C = &c.ChatRoomModel{
		Id:        chat.GetId(),
		Category:  chat.GetCategory(),
		CreatorId: chat.GetCreatorId(),
		Name:      chat.GetName(),
	}
	return mm
}

func (mm *MessageModel) SetChatId(id int) Message {
	mm.chatId = id
	return mm
}

func (mm *MessageModel) SetCreatedAt(at time.Time) Message {
	mm.CreatedAt = at
	return mm
}

func (mm *MessageModel) SetId(id int) Message {
	mm.Id = id
	return mm
}

func (mm *MessageModel) SetText(text string) Message {
	mm.Text = text
	return mm
}

func (mm *MessageModel) SetUser(user u.User) Message {
	mm.userId = user.GetId()
	mm.U = &u.UserDTO{
		Id:       user.GetId(),
		Username: user.GetUsername(),
		Name:     user.GetName(),
	}
	return mm
}

func (mm *MessageModel) SetUserId(id int) Message {
	mm.userId = id
	return mm
}
