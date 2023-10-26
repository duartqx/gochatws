package interfaces

import "time"

type ParserFunc func(out interface{}) error

type ChatRoom interface {
	GetCategory() int
	GetCreator() User
	GetCreatorId() int
	GetId() int
	GetName() string
	PopulateCreator(creator User)
	SetCreatorId(id int)
	SetId(id int)
}

type Message interface {
	GetId() int
	GetChatId() int
	GetUserId() int
	GetCreatedAt() *time.Time
	SetId(int) Message
	SetChatId(id int) Message
	SetUserId(id int) Message
	SetCreatedAt() Message
	GetText() string
	GetUser() User
	GetChat() ChatRoom
	PopulateUser(user User) Message
	PopulateChat(chat ChatRoom) Message
}

type Session interface {
	GetToken() string
}

type User interface {
	Clean() User
	GetId() int
	GetName() string
	GetPassword() string
	GetUsername() string
	SetId(id int)
	SetPassword(password string)
	UpdateFromAnother(other User)
}
