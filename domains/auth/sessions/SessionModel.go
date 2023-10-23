package sessions

import "time"

type SessionModel struct {
	Token      string    `db:"token"`
	UserId     int       `db:"user_id"`
	CreationAt time.Time `db:"creation_at"`
}

func (sm SessionModel) GetToken() string {
	return sm.Token
}
