package controllers

import (
	"io"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	i "github.com/duartqx/gochatws/core/interfaces"
)

type WebSocketController struct {
	chatRepository i.ChatRepository
}

func NewWebSocketController(chatRepository i.ChatRepository) *WebSocketController {
	return &WebSocketController{chatRepository: chatRepository}
}

func (wsc WebSocketController) Chat() func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				if err == io.EOF {
					break
				}
				continue
			}

			err = c.WriteMessage(mt, []byte("Received: "+string(msg)))
			if err != nil {
				break
			}
		}
	})
}
