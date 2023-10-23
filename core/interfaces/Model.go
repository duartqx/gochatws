package interfaces

type ParserFunc func(out interface{}) error

type User interface {
	Clean() User
	GetId() int
	SetId(id int)
	GetName() string
	GetUsername() string
	GetPassword() string
	UpdateFromAnother(other User)
}

type ChatRoom interface {
	GetId() int
	SetId(id int)
	GetName() string
	GetCategory() int
	GetCreatorId() int
	GetCreator() User
}
