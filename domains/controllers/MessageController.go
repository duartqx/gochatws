package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	e "github.com/duartqx/gochatws/core/errors"
	i "github.com/duartqx/gochatws/core/interfaces"
	m "github.com/duartqx/gochatws/domains/models"
	s "github.com/duartqx/gochatws/domains/services"
	"github.com/duartqx/gochatws/domains/utils"
)

type MessageController struct {
	chatRepository i.ChatRepository
	messageService *s.MessageService
}

func NewMessageController(
	chatRepository i.ChatRepository,
	messageService *s.MessageService,
) *MessageController {
	return &MessageController{
		chatRepository: chatRepository,
		messageService: messageService,
	}
}

func (mc MessageController) WebSocketChat() func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)

		creator, err := utils.GetUserFromLocals(c.Locals("user"))
		if err != nil {
			return
		}

		chat, err := mc.chatRepository.FindByParamId(c.Params("chat_id"))
		if err != nil {
			return
		}

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				if err == io.EOF {
					break
				}
				continue
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

			if err = c.WriteMessage(mt, msgJson); err != nil {
				break
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
