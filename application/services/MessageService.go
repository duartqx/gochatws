package services

import (
	"net/http"

	e "github.com/duartqx/gochatws/application/errors"
	h "github.com/duartqx/gochatws/application/http"

	m "github.com/duartqx/gochatws/domains/entities/message"
	r "github.com/duartqx/gochatws/domains/repositories"
)

type MessageService struct {
	messageRepository r.MessageRepository
}

func NewMessageService(messageRepository r.MessageRepository) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
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
