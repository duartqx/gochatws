package sessions

import "fmt"

type SessionStore map[string]Session

func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

func (ss *SessionStore) Set(key string, session Session) error {
	(*ss)[key] = session
	return nil
}

func (ss SessionStore) Get(key string) (Session, error) {
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
