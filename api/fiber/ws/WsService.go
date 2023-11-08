package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"

	m "github.com/duartqx/gochatws/domains/entities/message"
	r "github.com/duartqx/gochatws/domains/repositories"
)

type WebSocketService struct {
	ConnsMap          *map[int]*WebSocketBroadcaster
	messageRepository r.MessageRepository
}

func GetWebSocketService(messageRepository r.MessageRepository) *WebSocketService {
	return &WebSocketService{
		ConnsMap:          &map[int]*WebSocketBroadcaster{},
		messageRepository: messageRepository,
	}
}

func (wc *WebSocketService) Listen(conn *websocket.Conn, finish chan bool, userId, chatId int) {

	broadcaster := wc.getBroadcaster(chatId)
	*broadcaster.Conns = append(*broadcaster.Conns, conn)

	// Defer removing the connection from the pool
	defer wc.close(broadcaster, conn, finish)

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

func (wc *WebSocketService) getBroadcaster(id int) *WebSocketBroadcaster {
	var broadcaster *WebSocketBroadcaster
	broadcaster, ok := (*wc.ConnsMap)[id]
	if !ok {
		broadcaster = NewWebSocketBroadcaster(id)
		(*wc.ConnsMap)[id] = broadcaster
	}
	return broadcaster
}

func (wc *WebSocketService) close(
	broadcaster *WebSocketBroadcaster,
	conn *websocket.Conn,
	finish chan bool,
) {
	defer close(finish)
	broadcaster.Close(conn)
	if broadcaster.Length() == 0 {
		delete((*wc.ConnsMap), broadcaster.GetKey())
	}
}
