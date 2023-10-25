package chat

import (
	"strconv"

	i "github.com/duartqx/gochatws/core/interfaces"
	u "github.com/duartqx/gochatws/domains/users"
	"github.com/jmoiron/sqlx"
)

const baseJoinQuery string = `
	SELECT c.id, c.creator_id, c.name, c.category, u.id, u.username, u.name
	FROM ChatRoom AS c
	INNER JOIN User AS u
	ON u.id = c.creator_id
`

type ChatRoomRepository struct {
	db *sqlx.DB
	ur i.UserRepository
}

func NewChatRoomRepository(db *sqlx.DB, ur i.UserRepository) *ChatRoomRepository {
	return &ChatRoomRepository{
		db: db,
		ur: ur,
	}
}

func (crr ChatRoomRepository) GetModel() *ChatRoomModel {
	return &ChatRoomModel{U: &u.UserClean{}}
}

func (crr ChatRoomRepository) FindById(id int) (i.ChatRoom, error) {

	chatRoom := crr.GetModel()

	row := crr.db.QueryRow(baseJoinQuery+"WHERE c.id = $1 LIMIT 1", id)
	if err := row.Scan(
		&chatRoom.Id,
		&chatRoom.CreatorId,
		&chatRoom.Name,
		&chatRoom.Category,
		&chatRoom.U.Id,
		&chatRoom.U.Username,
		&chatRoom.U.Name,
	); err != nil {
		return nil, err
	}
	return chatRoom, nil
}

func (crr ChatRoomRepository) FindByParamId(id string) (i.ChatRoom, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return crr.FindById(idInt)
}

func (crr ChatRoomRepository) All() (*[]i.ChatRoom, error) {
	chatRooms := []i.ChatRoom{}
	rows, err := crr.db.Query(baseJoinQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		chatRoom := crr.GetModel()

		if err := rows.Scan(
			&chatRoom.Id,
			&chatRoom.CreatorId,
			&chatRoom.Name,
			&chatRoom.Category,
			&chatRoom.U.Id,
			&chatRoom.U.Username,
			&chatRoom.U.Name,
		); err != nil {
			return nil, err
		}

		var icr i.ChatRoom = chatRoom

		chatRooms = append(chatRooms, icr)
	}

	return &chatRooms, nil
}

func (crr ChatRoomRepository) Create(cr i.ChatRoom) error {
	_, err := crr.db.Exec(`
		INSERT INTO ChatRoom (creator_id, name, category)
		VALUES ($1, $2, $3)
	`, cr.GetCreatorId(), cr.GetName(), cr.GetCategory())
	if err != nil {
		return err
	}
	var chatId int
	err = crr.db.QueryRow("SELECT last_insert_rowid()").Scan(&chatId)
	if err != nil {
		return err
	}

	cr.SetId(chatId)

	return nil
}

func (crr ChatRoomRepository) populateCreator(cr *ChatRoomModel) error {
	creator, err := crr.ur.FindById(cr.CreatorId)
	if err != nil {
		return err
	}

	cr.U.Id = creator.GetId()
	cr.U.Username = creator.GetUsername()
	cr.U.Name = creator.GetName()

	return nil
}

func (crr ChatRoomRepository) ParseAndValidate(parser i.ParserFunc) (i.ChatRoom, error) {
	parsedChatRoom := crr.GetModel()
	if err := parser(parsedChatRoom); err != nil {
		return nil, err
	}
	if err := crr.populateCreator(parsedChatRoom); err != nil {
		return nil, err
	}
	return parsedChatRoom, nil
}
