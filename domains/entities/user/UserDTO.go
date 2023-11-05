package user

// UserDTO is returned from UserModel.Clean method to make sure the Password
// field is not leaked even if it's hashed.
type UserDTO struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Name     string `db:"name" json:"name"`
}

func (u UserDTO) Clean() User {
	return &u
}

func (u UserDTO) GetId() int {
	return u.Id
}

func (u UserDTO) GetName() string {
	return u.Name
}

func (u UserDTO) GetPassword() string {
	return ""
}

func (u UserDTO) GetUsername() string {
	return u.Username
}

func (u *UserDTO) SetId(id int) User {
	u.Id = id
	return u
}

func (u *UserDTO) SetName(name string) User {
	u.Name = name
	return u
}

func (u *UserDTO) SetPassword(password string) User {
	return u
}

func (u *UserDTO) SetUsername(username string) User {
	u.Username = username
	return u
}

func (u *UserDTO) UpdateFromAnother(other User) {}
