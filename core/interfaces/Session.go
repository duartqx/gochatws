package interfaces

type SessionStore interface {
	Set(key string, session Session) error
	Get(key string) (Session, error)
	Delete(key string) error
}
