package interfaces

type ParserFunc func(out interface{}) error

type User interface {
	Clean() User
	GetId() int
	GetName() string
	GetUsername() string
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
