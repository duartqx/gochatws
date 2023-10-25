package interfaces

type ChatRepository interface {
	FindById(id int) (ChatRoom, error)
	FindByParamId(id string) (ChatRoom, error)
	Create(cr ChatRoom) error
	All() (*[]ChatRoom, error)
}

type UserRepository interface {
	All() (*[]User, error)
	FindById(id int) (User, error)
	FindByUsername(username string) (User, error)
	ExistsByUsername(username string) bool
	Create(user User) error
	Delete(user User) error
}
