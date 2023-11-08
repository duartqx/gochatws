package ws

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketBroadcaster struct {
	Conns *[]*websocket.Conn
	Send  *chan []byte
	key   int
	mutex *sync.Mutex
}

func NewWebSocketBroadcaster(key int) *WebSocketBroadcaster {
	channel := make(chan []byte)
	return &WebSocketBroadcaster{
		Conns: &[]*websocket.Conn{},
		Send:  &channel,
		mutex: &sync.Mutex{},
		key:   key,
	}
}

func (wb WebSocketBroadcaster) Broadcast(msg *[]byte) error {
	wb.mutex.Lock()
	defer wb.mutex.Unlock()
	for _, conn := range *wb.Conns {
		err := conn.WriteMessage(websocket.TextMessage, *msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wb WebSocketBroadcaster) Close(conn *websocket.Conn) {
	wb.mutex.Lock()
	conn.Close()
	for i, c := range *wb.Conns {
		if c == conn {
			*wb.Conns = append((*wb.Conns)[:i], (*wb.Conns)[i+1:]...)
			break
		}
	}
	wb.mutex.Unlock()
}

func (wb WebSocketBroadcaster) Length() int {
	return len(*wb.Conns)
}

func (wb WebSocketBroadcaster) GetKey() int {
	return wb.key
}
