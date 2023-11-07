package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"

	m "github.com/duartqx/gochatws/domains/entities/message"
	r "github.com/duartqx/gochatws/domains/repositories"
)

type WebSocketService struct {
	Conns             *map[string]*WebSocketBroadcaster
	messageRepository r.MessageRepository
}

func GetWebSocketService(messageRepository r.MessageRepository) *WebSocketService {
	return &WebSocketService{
		Conns:             &map[string]*WebSocketBroadcaster{},
		messageRepository: messageRepository,
	}
}

func (wc *WebSocketService) Listen(conn *websocket.Conn, finish chan bool, userId, chatId int) {
	broadcaster := wc.GetBroadcaster(conn.Params("chat_id"))
	*broadcaster.Conns = append(*broadcaster.Conns, conn)

	// Defer removing the connection from the pool
	defer func() {
		broadcaster.Mutex.Lock()
		conn.Close()
		close(finish)
		for i, c := range *broadcaster.Conns {
			if c == conn {
				*broadcaster.Conns = append(
					(*broadcaster.Conns)[:i], (*broadcaster.Conns)[i+1:]...,
				)
				break
			}
		}
		broadcaster.Mutex.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		message := &m.MessageModel{Text: string(msg)}

		message.SetChatId(chatId).SetUserId(userId)

		err = wc.messageRepository.Create(message)
		if err != nil {
			break
		}

		msgJson, err := json.Marshal(message)
		if err != nil {
			break
		}

		broadcaster.Broadcast(&msgJson)
	}
}

func (wc WebSocketService) GetBroadcaster(id string) *WebSocketBroadcaster {
	var broadcaster *WebSocketBroadcaster
	broadcaster, ok := (*wc.Conns)[id]
	if !ok {
		broadcaster = NewWebSocketBroadcaster()
		(*wc.Conns)[id] = broadcaster
	}
	return broadcaster
}
