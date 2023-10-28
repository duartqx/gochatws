package models

type ChatCategoryModel struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

func GetChatCategories() *[5]ChatCategoryModel {
	return &[5]ChatCategoryModel{
		{Id: 1, Name: "Science"},
		{Id: 2, Name: "Tech"},
		{Id: 3, Name: "Movies"},
		{Id: 4, Name: "Video Games"},
		{Id: 5, Name: "TV"},
	}
}
