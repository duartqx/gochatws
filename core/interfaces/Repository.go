package interfaces

type Repository interface {
	FindById(id int) (User, error)
}
