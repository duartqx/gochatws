package repositories

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"

	c "github.com/duartqx/gochatws/domains/entities/chatroom"
	u "github.com/duartqx/gochatws/domains/entities/user"
	r "github.com/duartqx/gochatws/domains/repositories"
)

const baseChatJoinQuery string = `
	SELECT c.id, c.name, c.category, u.id, u.username, u.name
	FROM ChatRoom AS c
	INNER JOIN User AS u
	ON u.id = c.creator_id
`

type ChatRoomRepository struct {
	db             *sqlx.DB
	userRepository r.UserRepository
}

func NewChatRoomRepository(db *sqlx.DB, ur r.UserRepository) r.ChatRepository {
	return &ChatRoomRepository{
		db:             db,
		userRepository: ur,
	}
}

func (crr ChatRoomRepository) GetModel() *c.ChatRoomModel {
	return &c.ChatRoomModel{C: &u.UserDTO{}}
}

func (crr ChatRoomRepository) FindById(id int) (c.ChatRoom, error) {

	chatRoom := crr.GetModel()

	var (
		name            string
		category        int
		creatorId       int
		creatorUsername string
		creatorName     string
	)

	row := crr.db.QueryRow(baseChatJoinQuery+"WHERE c.id = $1 LIMIT 1", id)
	if err := row.Scan(
		&id,
		&name,
		&category,
		&creatorId,
		&creatorUsername,
		&creatorName,
	); err != nil {
		return nil, err
	}

	chatRoom.
		SetId(id).
		SetName(name).
		SetCategory(category).
		SetCreator(
			&u.UserDTO{Id: creatorId, Username: creatorUsername, Name: creatorName},
		)

	return chatRoom, nil
}

func (crr ChatRoomRepository) FindByParamId(id string) (c.ChatRoom, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return crr.FindById(idInt)
}

func (crr ChatRoomRepository) All() (*[]c.ChatRoom, error) {
	chatRooms := []c.ChatRoom{}
	rows, err := crr.db.Query(baseChatJoinQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		chatRoom := crr.GetModel()

		var (
			id              int
			name            string
			category        int
			creatorId       int
			creatorUsername string
			creatorName     string
		)

		if err := rows.Scan(
			&id,
			&name,
			&category,
			&creatorId,
			&creatorUsername,
			&creatorName,
		); err != nil {
			return nil, err
		}

		var icr c.ChatRoom = chatRoom.
			SetId(id).
			SetName(name).
			SetCategory(category).
			SetCreator(
				&u.UserDTO{
					Id:       creatorId,
					Username: creatorUsername,
					Name:     creatorName,
				},
			)

		chatRooms = append(chatRooms, icr)
	}

	return &chatRooms, nil
}

func (crr ChatRoomRepository) Create(cr c.ChatRoom) error {

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

	cr.SetId(int(chatId)).SetCreator(creator)

	return nil
}
