package interfaces

type ParserFunc func(out interface{}) error

type User interface {
	Clean() User
	GetId() int
	GetName() string
	GetPassword() string
	GetUsername() string
	SetId(id int)
	SetPassword(password string)
	UpdateFromAnother(other User)
}

type ChatRoom interface {
	GetCategory() int
	GetCreator() User
	GetCreatorId() int
	GetId() int
	GetName() string
	PopulateCreator(creator User)
	SetCreatorId(id int)
	SetId(id int)
}

type Session interface {
	GetToken() string
}
