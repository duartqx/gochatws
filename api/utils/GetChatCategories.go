package utils

import (
	c "github.com/duartqx/gochatws/domains/entities/chatroom"
)

func GetChatCategories() *[5]c.ChatCategoryModel {
	return c.GetChatCategories()
}
