package repositories

import (
	c "github.com/duartqx/gochatws/domains/entities/chatroom"
)

type ChatRepository interface {
	All() (*[]c.ChatRoom, error)
	Create(cr c.ChatRoom) error
	FindById(id int) (c.ChatRoom, error)
	FindByParamId(id string) (c.ChatRoom, error)
}
