package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/duartqx/gochatws/api/fiber/utils"
	w "github.com/duartqx/gochatws/api/fiber/ws"
	e "github.com/duartqx/gochatws/application/errors"
	s "github.com/duartqx/gochatws/application/services"
	m "github.com/duartqx/gochatws/domains/entities/message"
)

type MessageController struct {
	messageService *s.MessageService
	service        *w.WebSocketService
}

func NewMessageController(
	messageService *s.MessageService,
	service *w.WebSocketService,
) *MessageController {
	return &MessageController{
		messageService: messageService,
		service:        service,
	}
}

func (mc *MessageController) WebSocketChatController() func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		user, err := utils.GetUserFromLocals(c.Locals("user"))
		if err != nil {
			return
		}

		chat, err := mc.messageService.FindChatByParamId(c.Params("chat_id"))
		if err != nil {
			return
		}

		finish := make(chan bool, 1)

		go mc.service.Listen(c, finish, user.GetId(), chat.GetId())

		for {
			select {
			case <-finish:
				return
			default:
				// This case will run if the 'finish' channel is not ready to send
				// You can add a sleep here to prevent busy waiting
				time.Sleep(500 * time.Millisecond)
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
