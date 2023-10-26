package services

import (
	"net/http"

	e "github.com/duartqx/gochatws/core/errors"
	h "github.com/duartqx/gochatws/core/http"
	i "github.com/duartqx/gochatws/core/interfaces"
)

type MessageService struct {
	messageRepository i.MessageRepository
}

func NewMessageService(messageRepository i.MessageRepository) *MessageService {
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

func (ms MessageService) Create(msg i.Message) *h.HttpResponse {

	if err := ms.messageRepository.Create(msg); err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.CustomMessageError(err.Error()),
		}
	}
	return &h.HttpResponse{Status: http.StatusCreated, Body: msg}
}
