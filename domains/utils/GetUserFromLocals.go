package utils

import (
	"encoding/json"
	"fmt"

	i "github.com/duartqx/gochatws/core/interfaces"
	m "github.com/duartqx/gochatws/domains/models"
	"github.com/gofiber/fiber/v2"
)

func GetUserFromLocals(localUser interface{}) (i.User, error) {
	if localUser == nil {
		return nil, fmt.Errorf("User not found on Locals\n")
	}
	userBytes, err := json.Marshal(localUser)
	if err != nil {
		return nil, err
	}

	userStruct := &m.UserModel{}
	err = json.Unmarshal(userBytes, userStruct)
	if err != nil {
		return nil, err
	}
	return userStruct, nil
}

func BuildTemplateContext(c *fiber.Ctx, m *fiber.Map) *fiber.Map {
	user, _ := GetUserFromLocals(c.Locals("user"))

	(*m)["User"] = user

	return m
}
