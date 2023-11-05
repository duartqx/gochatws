package user

type User interface {
	Clean() User
	GetId() int
	GetName() string
	GetPassword() string
	GetUsername() string
	SetId(id int) User
	SetName(name string) User
	SetPassword(password string) User
	SetUsername(username string) User
	UpdateFromAnother(other User)
}
