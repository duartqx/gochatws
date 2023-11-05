package message

import (
	"time"

	c "github.com/duartqx/gochatws/domains/entities/chatroom"
	u "github.com/duartqx/gochatws/domains/entities/user"
)

type Message interface {
	GetChat() c.ChatRoom
	GetChatId() int
	GetCreatedAt() *time.Time
	GetId() int
	GetText() string
	GetUser() u.User
	GetUserId() int
	SetChat(chat c.ChatRoom) Message
	SetChatId(id int) Message
	SetCreatedAt(at time.Time) Message
	SetId(id int) Message
	SetText(text string) Message
	SetUser(user u.User) Message
	SetUserId(id int) Message
}
