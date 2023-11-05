package sessions

type Session interface {
	GetToken() string
}
