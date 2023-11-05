package utils

import "github.com/gofiber/fiber/v2"

func BuildTemplateContext(c *fiber.Ctx, m *fiber.Map) *fiber.Map {
	user, _ := GetUserFromLocals(c.Locals("user"))

	(*m)["User"] = user

	return m
}
