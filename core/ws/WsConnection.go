package ws

import (
	"github.com/gofiber/contrib/websocket"
)

type WsConnection struct {
	Conn   *websocket.Conn
	Send   chan []byte
	ChatId string
}

func GetConnectionStore() *[]*WsConnection {
	return &[]*WsConnection{}
}

func (wc WsConnection) GetConn() *websocket.Conn {
	return wc.Conn
}
