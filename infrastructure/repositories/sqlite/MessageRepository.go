package repositories

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	c "github.com/duartqx/gochatws/domains/entities/chatroom"
	m "github.com/duartqx/gochatws/domains/entities/message"
	u "github.com/duartqx/gochatws/domains/entities/user"
	r "github.com/duartqx/gochatws/domains/repositories"
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

const messageOrder string = "ORDER BY m.created_at ASC"

type MessageRepository struct {
	db             *sqlx.DB
	userRepository r.UserRepository
	chatRepository r.ChatRepository
}

func NewMessageRepository(
	db *sqlx.DB,
	userRepository r.UserRepository,
	chatRepository r.ChatRepository,
) r.MessageRepository {
	return &MessageRepository{
		db:             db,
		userRepository: userRepository,
		chatRepository: chatRepository,
	}
}

func (mr MessageRepository) GetMessageModel() *m.MessageModel {
	return &m.MessageModel{}
}

func (mr MessageRepository) GetChatModel() *c.ChatRoomModel {
	return &c.ChatRoomModel{}
}

func (mr MessageRepository) GetUserModel() *u.UserDTO {
	return &u.UserDTO{}
}

func (mr MessageRepository) FindById(id int) (m.Message, error) {

	msg := mr.GetMessageModel()
	chat := mr.GetChatModel()
	user := mr.GetUserModel()

	var (
		// MsgInfo
		msgCreatedAt time.Time
		msgText      string

		// ChatInfo
		chatId        int
		chatName      string
		chatCreatorId int
		chatCategory  int

		// MsgUserInfo
		userId       int
		userUsername string
		userName     string
	)

	row := mr.db.QueryRow(baseMessageJoinQuery+"WHERE m.id = $1 LIMIT 1", id)
	if err := row.Scan(
		// Message info
		&id,
		&msgCreatedAt,
		&msgText,
		// Chat info
		&chatId,
		&chatName,
		&chatCreatorId,
		&chatCategory,
		// User info
		&userId,
		&userUsername,
		&userName,
	); err != nil {
		return nil, err
	}

	msg.
		SetId(id).
		SetCreatedAt(msgCreatedAt).
		SetText(msgText).
		SetChat(
			chat.
				SetId(chatId).
				SetName(chatName).
				SetCreatorId(chatCreatorId).
				SetCategory(chatCategory),
		).
		SetUser(
			user.
				SetId(userId).
				SetUsername(userUsername).
				SetName(userName),
		)

	return msg, nil
}

func (mr MessageRepository) FindByParamId(id string) (m.Message, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return mr.FindById(idInt)
}

func (mr MessageRepository) FindByChatId(id int) ([]m.Message, error) {

	messages := []m.Message{}

	query := baseMessageJoinQuery + " WHERE c.id = $1 " + messageOrder

	rows, err := mr.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		msg := mr.GetMessageModel()
		chat := mr.GetChatModel()
		user := mr.GetUserModel()

		var (
			// MsgInfo
			msgCreatedAt time.Time
			msgText      string

			// ChatInfo
			chatId        int
			chatName      string
			chatCreatorId int
			chatCategory  int

			// MsgUserInfo
			userId       int
			userUsername string
			userName     string
		)

		if err := rows.Scan(
			// Message info
			&id,
			&msgCreatedAt,
			&msgText,
			// Chat info
			&chatId,
			&chatName,
			&chatCreatorId,
			&chatCategory,
			// User info
			&userId,
			&userUsername,
			&userName,
		); err != nil {
			return nil, err
		}

		msg.
			SetId(id).
			SetCreatedAt(msgCreatedAt).
			SetText(msgText).
			SetChat(
				chat.
					SetId(chatId).
					SetName(chatName).
					SetCreatorId(chatCreatorId).
					SetCategory(chatCategory),
			).
			SetUser(
				user.
					SetId(userId).
					SetUsername(userUsername).
					SetName(userName),
			)

		var iMessage m.Message = msg

		messages = append(messages, iMessage)
	}

	return messages, nil
}

func (mr MessageRepository) FindByChatParamId(id string) ([]m.Message, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return mr.FindByChatId(idInt)
}

func (mr MessageRepository) Create(msg m.Message) error {
	if msg.GetChatId() == 0 || msg.GetUserId() == 0 {
		return fmt.Errorf("Could not save message, check ChatId and UserId")
	}

	user, err := mr.userRepository.FindById(msg.GetUserId())
	if err != nil {
		return fmt.Errorf("Could not find user with id of %d\n", msg.GetUserId())
	}

	chat, err := mr.chatRepository.FindById(msg.GetChatId())
	if err != nil {
		return fmt.Errorf("Could not find chat with id of %d\n", msg.GetChatId())
	}

	msg.SetCreatedAt(time.Now())

	result, err := mr.db.Exec(`
			INSERT INTO Message (chat_id, user_id, text, created_at)
			VALUES ($1, $2, $3, $4)
		`,
		msg.GetChatId(),
		msg.GetUserId(),
		msg.GetText(),
		msg.GetCreatedAt(),
	)
	if err != nil {
		return err
	}

	msgId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	msg.SetId(int(msgId)).SetChat(chat).SetUser(user)

	return nil
}
