package services

import (
	"net/http"

	e "github.com/duartqx/gochatws/application/errors"
	h "github.com/duartqx/gochatws/application/http"

	c "github.com/duartqx/gochatws/domains/entities/chatroom"
	m "github.com/duartqx/gochatws/domains/entities/message"
	r "github.com/duartqx/gochatws/domains/repositories"
)

type MessageService struct {
	messageRepository  r.MessageRepository
	chatRoomRepository r.ChatRepository
}

func NewMessageService(
	messageRepository r.MessageRepository,
	chatRoomRepository r.ChatRepository,
) *MessageService {
	return &MessageService{
		messageRepository:  messageRepository,
		chatRoomRepository: chatRoomRepository,
	}
}

func (ms MessageService) ChatMessages(id int) *h.HttpResponse {
	messages, err := ms.messageRepository.FindByChatId(id)
	if err != nil {
		return &h.HttpResponse{
			Status: http.StatusNotFound,
			Body:   e.NotFoundError,
		}
	}
	return &h.HttpResponse{Status: http.StatusOK, Body: messages}
}

func (ms MessageService) Create(msg m.Message) *h.HttpResponse {

	if err := ms.messageRepository.Create(msg); err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.CustomMessageError(err.Error()),
		}
	}
	return &h.HttpResponse{Status: http.StatusCreated, Body: msg}
}

func (ms MessageService) FindChatByParamId(id string) (c.ChatRoom, error) {
	return ms.chatRoomRepository.FindByParamId(id)
}
