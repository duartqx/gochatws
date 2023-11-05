package chatroom

import (
	u "github.com/duartqx/gochatws/domains/entities/user"
)

type ChatRoom interface {
	GetCategory() int
	GetCreator() u.User
	GetCreatorId() int
	GetId() int
	GetName() string
	SetCategory(category int) ChatRoom
	SetCreator(user u.User) ChatRoom
	SetCreatorId(id int) ChatRoom
	SetId(id int) ChatRoom
	SetName(name string) ChatRoom
}
