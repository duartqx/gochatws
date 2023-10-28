package services

import (
	"net/http"

	e "github.com/duartqx/gochatws/core/errors"
	h "github.com/duartqx/gochatws/core/http"
	i "github.com/duartqx/gochatws/core/interfaces"
)

type ChatApiResponse struct {
	Category int    `json:"category"`
	Creator  i.User `json:"creator"`
	Href     string `json:"href"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
}

type ChatRoomService struct {
	chatRoomRepository i.ChatRepository
}

func NewChatRoomService(crr i.ChatRepository) *ChatRoomService {
	return &ChatRoomService{
		chatRoomRepository: crr,
	}
}

func (crs ChatRoomService) All() *h.HttpResponse {
	chatRooms, err := crs.chatRoomRepository.All()
	if err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
	}
	return &h.HttpResponse{Status: http.StatusOK, Body: chatRooms}
}

func (crs ChatRoomService) One(paramId string) *h.HttpResponse {
	chatRoom, err := crs.chatRoomRepository.FindByParamId(paramId)
	if err != nil {
		return &h.HttpResponse{Status: http.StatusNotFound, Body: e.NotFoundError}
	}
	return &h.HttpResponse{Status: http.StatusOK, Body: chatRoom}
}

func (crs ChatRoomService) Create(chatRoom i.ChatRoom) *h.HttpResponse {
	if err := crs.chatRoomRepository.Create(chatRoom); err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
	}
	return &h.HttpResponse{Status: http.StatusCreated, Body: chatRoom}
}
