package utils

import (
	m "github.com/duartqx/gochatws/domains/models"
)

func GetChatCategories() *[5]m.ChatCategoryModel {
	return m.GetChatCategories()
}
