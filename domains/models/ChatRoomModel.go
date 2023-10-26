package models

import (
	i "github.com/duartqx/gochatws/core/interfaces"
)

type ChatRoomModel struct {
	Id        int        `db:"id" json:"id"`
	CreatorId int        `db:"creator_id" json:"creator_id"`
	Name      string     `db:"name" json:"name"`
	Category  int        `db:"category" json:"category"`
	C         *UserClean `json:"creator"`
}

func (crm ChatRoomModel) GetId() int {
	return crm.Id
}

func (crm *ChatRoomModel) SetId(id int) {
	crm.Id = id
}

func (crm ChatRoomModel) GetName() string {
	return crm.Name
}

func (crm ChatRoomModel) GetCategory() int {
	return crm.Category
}

func (crm ChatRoomModel) GetCreatorId() int {
	return crm.CreatorId
}

func (crm *ChatRoomModel) SetCreatorId(id int) {
	crm.CreatorId = id
}

func (crm ChatRoomModel) GetCreator() i.User {
	return crm.C
}

func (crm *ChatRoomModel) PopulateCreator(user i.User) {
	crm.C = &UserClean{
		Id:       user.GetId(),
		Username: user.GetUsername(),
		Name:     user.GetName(),
	}
}