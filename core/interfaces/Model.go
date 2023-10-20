package interfaces

type ParserFunc func(out interface{}) error

type User interface {
	Clean() User
	GetId() int
	GetName() string
	GetUsername() string
	UpdateFromAnother(other User)
}
