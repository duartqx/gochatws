package core

type Repository interface {
	getModel() *Model
	findById(id int) (*Model, error)
}
