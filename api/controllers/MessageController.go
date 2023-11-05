package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/duartqx/gochatws/api/utils"
	w "github.com/duartqx/gochatws/api/ws"
	e "github.com/duartqx/gochatws/application/errors"
	s "github.com/duartqx/gochatws/application/services"
	m "github.com/duartqx/gochatws/domains/entities/message"
	r "github.com/duartqx/gochatws/domains/repositories"
)

type MessageController struct {
	chatRepository r.ChatRepository
	messageService *s.MessageService
	connStore      *[]*w.WsConnection
	mutex          *sync.Mutex
}

func NewMessageController(
	chatRepository r.ChatRepository,
	messageService *s.MessageService,
	connStore *[]*w.WsConnection,
) *MessageController {
	return &MessageController{
		chatRepository: chatRepository,
		messageService: messageService,
		connStore:      connStore,
		mutex:          &sync.Mutex{},
	}
}

func (mc *MessageController) WebSocketChat() func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		conn := &w.WsConnection{
			Conn:   c,
			Send:   make(chan []byte),
			ChatId: c.Params("chat_id"),
		}
		*mc.connStore = append(*mc.connStore, conn)

		creator, err := utils.GetUserFromLocals(c.Locals("user"))
		if err != nil {
			return
		}

		chat, err := mc.chatRepository.FindByParamId(c.Params("chat_id"))
		if err != nil {
			return
		}

		go func(conn *w.WsConnection) {
			defer func() {
				conn.Conn.Close()
				// Remove the connection from the global list when done
				for i, c := range *mc.connStore {
					if c == conn {
						*mc.connStore = append((*mc.connStore)[:i], (*mc.connStore)[i+1:]...)
						break
					}
				}
			}()

			for {
				_, msg, err := conn.Conn.ReadMessage()
				if err != nil {
					break
				}

				message := &m.MessageModel{Text: string(msg)}

				message.SetChatId(chat.GetId()).SetUserId(creator.GetId())

				response := mc.messageService.Create(message)
				if response.Status != http.StatusCreated {
					break
				}

				msgJson, err := json.Marshal(message)
				if err != nil {
					break
				}

				// Broadcast the message to all connections
				for _, c := range *mc.connStore {
					if c.ChatId == conn.ChatId {
						c.Send <- msgJson
					}
				}
			}
		}(conn)

		// Write to the WebSocket from the broadcast channel
		for {
			select {
			case msg, ok := <-conn.Send:

				mc.mutex.Lock()

				if !ok {
					conn.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				conn.Conn.WriteMessage(websocket.TextMessage, msg)

				mc.mutex.Unlock()
			}
		}
	})
}

func (mc MessageController) GetChatMessages(c *fiber.Ctx) error {
	chatId, err := strconv.Atoi(c.Params("chat_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}
	response := mc.messageService.ChatMessages(chatId)
	return c.Status(response.Status).JSON(response.Body)
}

func (mc MessageController) CreateMessage(c *fiber.Ctx) error {
	creator, err := utils.GetUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(e.UnauthorizedError)
	}

	chatId, err := strconv.Atoi(c.Params("chat_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}

	message := &m.MessageModel{}
	if err := c.BodyParser(message); err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}
	message.SetUserId(creator.GetId()).SetChatId(chatId)

	response := mc.messageService.Create(message)
	return c.Status(response.Status).JSON(response.Body)
}
