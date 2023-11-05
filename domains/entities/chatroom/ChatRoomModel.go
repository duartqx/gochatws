package chatroom

import (
	u "github.com/duartqx/gochatws/domains/entities/user"
)

type ChatRoomModel struct {
	Id        int    `db:"id" json:"id"`
	CreatorId int    `db:"creator_id" json:"creator_id"`
	Name      string `db:"name" json:"name"`
	Category  int    `db:"category" json:"category"`
	C         u.User `json:"creator"`
}

func (crm ChatRoomModel) GetCategory() int {
	return crm.Category
}

func (crm ChatRoomModel) GetCreator() u.User {
	return crm.C
}

func (crm ChatRoomModel) GetCreatorId() int {
	return crm.CreatorId
}

func (crm ChatRoomModel) GetId() int {
	return crm.Id
}

func (crm ChatRoomModel) GetName() string {
	return crm.Name
}

func (crm *ChatRoomModel) SetCategory(category int) ChatRoom {
	crm.Category = category
	return crm
}

func (crm *ChatRoomModel) SetCreator(user u.User) ChatRoom {
	crm.CreatorId = user.GetId()
	crm.C = &u.UserDTO{
		Id:       user.GetId(),
		Username: user.GetUsername(),
		Name:     user.GetName(),
	}
	return crm
}

func (crm *ChatRoomModel) SetCreatorId(id int) ChatRoom {
	crm.CreatorId = id
	return crm
}

func (crm *ChatRoomModel) SetId(id int) ChatRoom {
	crm.Id = id
	return crm
}

func (crm *ChatRoomModel) SetName(name string) ChatRoom {
	crm.Name = name
	return crm
}
