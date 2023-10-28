package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	e "github.com/duartqx/gochatws/core/errors"
	i "github.com/duartqx/gochatws/core/interfaces"
	m "github.com/duartqx/gochatws/domains/models"
	s "github.com/duartqx/gochatws/domains/services"
	"github.com/duartqx/gochatws/domains/utils"
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
	chat, ok := response.Body.(i.ChatRoom)
	if !ok {
		return c.Render("404", fiber.Map{"Title": "404 Not Found"})
	}
	return c.Render(
		"chat",
		utils.BuildTemplateContext(c, &fiber.Map{
			"Title":  chat.GetName(),
			"ChatId": chat.GetId(),
			"Host":   "127.0.0.1:8000/",
		}),
	)
}

func (crc ChatRoomController) Create(c *fiber.Ctx) error {

	creator, err := utils.GetUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(e.UnauthorizedError)
	}

	chatRoom := &m.ChatRoomModel{}
	if err := c.BodyParser(chatRoom); err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}
	chatRoom.SetCreatorId(creator.GetId())

	response := crc.chatRoomService.Create(chatRoom)
	return c.Status(response.Status).JSON(response.Body)
}
