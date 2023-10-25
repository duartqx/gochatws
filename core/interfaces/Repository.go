package interfaces

type ChatRepository interface {
	FindById(id int) (ChatRoom, error)
	FindByParamId(id string) (ChatRoom, error)
	ParseAndValidate(parser ParserFunc) (ChatRoom, error)
	Create(cr ChatRoom) error
	All() (*[]ChatRoom, error)
}

type UserRepository interface {
	FindById(id int) (User, error)
	FindByUsername(username string) (User, error)
}
