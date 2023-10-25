package interfaces

type ChatRepository interface {
	All() (*[]ChatRoom, error)
	Create(cr ChatRoom) error
	FindById(id int) (ChatRoom, error)
	FindByParamId(id string) (ChatRoom, error)
}

type UserRepository interface {
	All() (*[]User, error)
	Create(user User) error
	Delete(user User) error
	ExistsByUsername(username string) bool
	FindById(id int) (User, error)
	FindByUsername(username string) (User, error)
}
