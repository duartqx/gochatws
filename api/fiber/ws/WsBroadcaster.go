package ws

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketBroadcaster struct {
	Conns *[]*websocket.Conn
	Send  *chan []byte
	Mutex *sync.Mutex
}

func (wb WebSocketBroadcaster) Broadcast(msg *[]byte) error {
	wb.Mutex.Lock()
	defer wb.Mutex.Unlock()
	for _, conn := range *wb.Conns {
		err := conn.WriteMessage(websocket.TextMessage, *msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewWebSocketBroadcaster() *WebSocketBroadcaster {
	channel := make(chan []byte)
	return &WebSocketBroadcaster{
		Conns: &[]*websocket.Conn{},
		Send:  &channel,
		Mutex: &sync.Mutex{},
	}
}
