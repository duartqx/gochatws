package repositories

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"

	i "github.com/duartqx/gochatws/core/interfaces"
	m "github.com/duartqx/gochatws/domains/models"
)

const baseChatJoinQuery string = `
	SELECT c.id, c.creator_id, c.name, c.category, u.id, u.username, u.name
	FROM ChatRoom AS c
	INNER JOIN User AS u
	ON u.id = c.creator_id
`

type ChatRoomRepository struct {
	db             *sqlx.DB
	userRepository i.UserRepository
}

func NewChatRoomRepository(db *sqlx.DB, ur i.UserRepository) *ChatRoomRepository {
	return &ChatRoomRepository{
		db:             db,
		userRepository: ur,
	}
}

func (crr ChatRoomRepository) GetModel() *m.ChatRoomModel {
	return &m.ChatRoomModel{C: &m.UserClean{}}
}

func (crr ChatRoomRepository) FindById(id int) (i.ChatRoom, error) {

	chatRoom := crr.GetModel()

	row := crr.db.QueryRow(baseChatJoinQuery+"WHERE c.id = $1 LIMIT 1", id)
	if err := row.Scan(
		&chatRoom.Id,
		&chatRoom.CreatorId,
		&chatRoom.Name,
		&chatRoom.Category,
		&chatRoom.C.Id,
		&chatRoom.C.Username,
		&chatRoom.C.Name,
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
	rows, err := crr.db.Query(baseChatJoinQuery)
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
			&chatRoom.C.Id,
			&chatRoom.C.Username,
			&chatRoom.C.Name,
		); err != nil {
			return nil, err
		}

		var icr i.ChatRoom = chatRoom

		chatRooms = append(chatRooms, icr)
	}

	return &chatRooms, nil
}

func (crr ChatRoomRepository) Create(cr i.ChatRoom) error {

	if cr.GetCreatorId() == 0 {
		return fmt.Errorf("Invalid Creator Id")
	}

	creator, err := crr.userRepository.FindById(cr.GetCreatorId())
	if err != nil {
		return err
	}

	result, err := crr.db.Exec(
		"INSERT INTO ChatRoom (creator_id, name, category) VALUES ($1, $2, $3)",
		cr.GetCreatorId(),
		cr.GetName(),
		cr.GetCategory(),
	)
	if err != nil {
		return err
	}

	chatId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	cr.SetId(int(chatId))
	cr.PopulateCreator(creator)

	return nil
}
