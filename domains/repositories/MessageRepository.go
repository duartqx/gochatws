package repositories

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"

	i "github.com/duartqx/gochatws/core/interfaces"
	m "github.com/duartqx/gochatws/domains/models"
)

const baseMessageJoinQuery string = `
	SELECT
		-- Message
		m.id,
		m.created_at,
		m.text,

		-- Chat
		c.id,
		c.name,
		c.creator_id,
		c.category,

		-- User
		u.id,
		u.username,
		u.name
	FROM Message AS m

	-- Chat informations
	INNER JOIN ChatRoom AS c
	ON m.chat_id = c.id

	-- User informations
	INNER JOIN User AS u
	ON m.user_id = u.id
`

const messageOrder string = "ORDER BY m.created_at DESC"

type MessageRepository struct {
	db             *sqlx.DB
	userRepository i.UserRepository
	chatRepository i.ChatRepository
}

func NewMessageRepository(
	db *sqlx.DB,
	userRepository i.UserRepository,
	chatRepository i.ChatRepository,
) *MessageRepository {
	return &MessageRepository{
		db:             db,
		userRepository: userRepository,
		chatRepository: chatRepository,
	}
}

func (mr MessageRepository) GetMessageModel() *m.MessageModel {
	return &m.MessageModel{}
}

func (mr MessageRepository) GetChatModel() *m.ChatRoomModel {
	return &m.ChatRoomModel{}
}

func (mr MessageRepository) GetUserModel() *m.UserClean {
	return &m.UserClean{}
}

func (mr MessageRepository) FindById(id int) (i.Message, error) {

	msg := mr.GetMessageModel()
	chat := mr.GetChatModel()
	user := mr.GetUserModel()

	row := mr.db.QueryRow(baseMessageJoinQuery+"WHERE m.id = $1 LIMIT 1", id)
	if err := row.Scan(
		// Message info
		&msg.Id,
		&msg.CreatedAt,
		&msg.Text,
		// Chat info
		&chat.Id,
		&chat.Name,
		&chat.CreatorId,
		&chat.Category,
		// User info
		&user.Id,
		&user.Username,
		&user.Name,
	); err != nil {
		return nil, err
	}
	msg.PopulateChat(chat)
	msg.PopulateUser(user)

	return msg, nil
}

func (mr MessageRepository) FindByParamId(id string) (i.Message, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return mr.FindById(idInt)
}

func (mr MessageRepository) FindByChatId(id int) ([]i.Message, error) {

	messages := []i.Message{}

	query := baseMessageJoinQuery + " WHERE c.id = $1 " + messageOrder

	rows, err := mr.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		msg := mr.GetMessageModel()
		chat := mr.GetChatModel()
		user := mr.GetUserModel()

		if err := rows.Scan(
			// Message info
			&msg.Id,
			&msg.CreatedAt,
			&msg.Text,
			// Chat info
			&chat.Id,
			&chat.Name,
			&chat.CreatorId,
			&chat.Category,
			// User info
			&user.Id,
			&user.Username,
			&user.Name,
		); err != nil {
			return nil, err
		}
		msg.PopulateChat(chat)
		msg.PopulateUser(user)

		var iMessage i.Message = msg

		messages = append(messages, iMessage)
	}

	return messages, nil
}

func (mr MessageRepository) FindByChatParamId(id string) ([]i.Message, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return mr.FindByChatId(idInt)
}

func (mr MessageRepository) Create(m i.Message) error {
	if m.GetChatId() == 0 || m.GetUserId() == 0 {
		return fmt.Errorf("Could not save message, check ChatId and UserId")
	}

	user, err := mr.userRepository.FindById(m.GetUserId())
	if err != nil {
		return fmt.Errorf("Could not find user with id of %d\n", m.GetUserId())
	}

	chat, err := mr.chatRepository.FindById(m.GetChatId())
	if err != nil {
		return fmt.Errorf("Could not find chat with id of %d\n", m.GetChatId())
	}

	m.SetCreatedAt()

	result, err := mr.db.Exec(`
			INSERT INTO Message (chat_id, user_id, text, created_at)
			VALUES ($1, $2, $3, $4)
		`,
		m.GetChatId(),
		m.GetUserId(),
		m.GetText(),
		m.GetCreatedAt(),
	)
	if err != nil {
		return err
	}

	msgId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.SetId(int(msgId)).PopulateChat(chat).PopulateUser(user)

	return nil
}
