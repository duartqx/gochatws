package controllers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/duartqx/gochatws/api/fiber/utils"
	e "github.com/duartqx/gochatws/application/errors"
	s "github.com/duartqx/gochatws/application/services"
	cr "github.com/duartqx/gochatws/domains/entities/chatroom"
)

type ChatRoomController struct {
	chatRoomService *s.ChatRoomService
}

func NewChatRoomController(crs *s.ChatRoomService) *ChatRoomController {
	return &ChatRoomController{
		chatRoomService: crs,
	}
}

func (crc ChatRoomController) All(c *fiber.Ctx) error {
	response := crc.chatRoomService.All()
	return c.Status(response.Status).JSON(response.Body)
}

func (crc ChatRoomController) One(c *fiber.Ctx) error {
	response := crc.chatRoomService.One(c.Params("id"))
	return c.Status(response.Status).JSON(response.Body)
}

func (crc ChatRoomController) ChatView(c *fiber.Ctx) error {
	response := crc.chatRoomService.One(c.Params("id"))
	if response.Status != http.StatusOK {
		return c.Render("404", fiber.Map{"Title": "404 Not Found"})
	}
	chat, ok := response.Body.(cr.ChatRoom)
	if !ok {
		return c.Render("404", fiber.Map{"Title": "404 Not Found"})
	}

	endpoint := fmt.Sprintf(
		"ws://%s/api/chat/%d/ws/connect", c.Hostname(), chat.GetId(),
	)

	return c.Render(
		"chat",
		utils.BuildTemplateContext(c, &fiber.Map{
			"Title":      chat.GetName(),
			"ChatId":     chat.GetId(),
			"WsEndpoint": endpoint,
		}),
	)
}

func (crc ChatRoomController) Create(c *fiber.Ctx) error {

	creator, err := utils.GetUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(e.UnauthorizedError)
	}

	chatRoom := &cr.ChatRoomModel{}
	if err := c.BodyParser(chatRoom); err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}
	chatRoom.SetCreatorId(creator.GetId())

	response := crc.chatRoomService.Create(chatRoom)
	return c.Status(response.Status).JSON(response.Body)
}
