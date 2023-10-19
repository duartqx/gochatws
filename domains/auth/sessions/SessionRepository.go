package sessions

import "github.com/jmoiron/sqlx"

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (sr SessionRepository) GetModel() *SessionModel {
	return &SessionModel{}
}

func (sr SessionRepository) FindByToken(token string) (*SessionModel, error) {
	s := sr.GetModel()
	err := sr.db.Get(s, "SELECT * FROM Session WHERE ID = $1 LIMIT 1", token)
	return s, err
}

func (sr SessionRepository) FindByUserId(id int) (*[]SessionModel, error) {
	sessions := &[]SessionModel{}
	err := sr.db.Get(sessions, "SELECT * FROM Session WHERE UserId = $1", id)
	return sessions, err
}
