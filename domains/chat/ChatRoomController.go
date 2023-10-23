package chat

import (
	"net/http"

	e "github.com/duartqx/gochatws/core/errors"
	i "github.com/duartqx/gochatws/core/interfaces"
	"github.com/gofiber/fiber/v2"
)

type ChatRoomController struct {
	crr i.ChatRepository
}

func NewChatRoomController(crr i.ChatRepository) *ChatRoomController {
	return &ChatRoomController{
		crr: crr,
	}
}

func (crc ChatRoomController) All(c *fiber.Ctx) error {
	chatRooms, err := crc.crr.All()
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.InternalError)
	}
	return c.
		Status(http.StatusOK).
		JSON(chatRooms)
}

func (crc ChatRoomController) One(c *fiber.Ctx) error {
	chatRoom, err := crc.crr.FindByParamId(c.Params("id"))
	if err != nil {
		return c.
			Status(http.StatusNotFound).
			JSON(e.NotFoundError)
	}
	return c.
		Status(http.StatusOK).
		JSON(chatRoom)
}

func (crc ChatRoomController) Create(c *fiber.Ctx) error {
	parsedChatRoom, err := crc.crr.ParseAndValidate(c.BodyParser)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(e.BadRequestError)
	}
	if err := crc.crr.Create(parsedChatRoom); err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.InternalError)
	}
	return c.
		Status(http.StatusCreated).
		JSON(parsedChatRoom)
}
