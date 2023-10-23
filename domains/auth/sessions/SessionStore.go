package sessions

import (
	"fmt"

	i "github.com/duartqx/gochatws/core/interfaces"
)

type SessionStore map[string]i.Session

func NewSessionStore() i.SessionStore {
	return &SessionStore{}
}

func (ss *SessionStore) Set(key string, session i.Session) error {
	(*ss)[key] = session
	return nil
}

func (ss SessionStore) Get(key string) (i.Session, error) {
	session, ok := ss[key]
	if !ok {
		return nil, fmt.Errorf("Key not found")
	}
	return session, nil
}

func (ss *SessionStore) Delete(key string) error {
	delete(*ss, key)
	return nil
}
