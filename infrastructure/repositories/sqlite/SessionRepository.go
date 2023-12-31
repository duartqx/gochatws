package repositories

import (
	s "github.com/duartqx/gochatws/infrastructure/sessions"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (sr SessionRepository) GetModel() *s.SessionModel {
	return &s.SessionModel{}
}

func (sr SessionRepository) FindByToken(token string) (s.Session, error) {
	s := sr.GetModel()
	err := sr.db.Get(s, "SELECT * FROM Session WHERE ID = $1 LIMIT 1", token)
	return s, err
}

func (sr SessionRepository) FindByUserId(id int) (*[]s.Session, error) {
	sessions := []s.Session{}
	rows, err := sr.db.Query(
		"SELECT token, user_id, creation_at FROM Session WHERE UserId = $1",
		id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		session := sr.GetModel()

		err := rows.Scan(&session.Token, &session.UserId, &session.CreationAt)
		if err != nil {
			return nil, err
		}

		var iSession s.Session = session

		sessions = append(sessions, iSession)

	}
	return &sessions, err
}
