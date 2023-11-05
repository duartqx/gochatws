package repositories

import (
	m "github.com/duartqx/gochatws/domains/entities/message"
)

type MessageRepository interface {
	FindById(id int) (m.Message, error)
	FindByParamId(id string) (m.Message, error)
	FindByChatId(id int) ([]m.Message, error)
	FindByChatParamId(id string) ([]m.Message, error)
	Create(m m.Message) error
}
