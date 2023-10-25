package interfaces

type ParserFunc func(out interface{}) error

type User interface {
	Clean() User
	GetId() int
	SetId(id int)
	SetPassword(password string)
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
	SetCreatorId(id int)
	GetCreator() User
	PopulateCreator(creator User)
}

type Session interface {
	GetToken() string
}
