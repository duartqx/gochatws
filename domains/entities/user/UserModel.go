package user

type UserModel struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username" validate:"email,required"`
	Name     string `db:"name" json:"name" validate:"required,min=3,max=50"`
	Password string `db:"password" json:"password"`
}

func (u UserModel) Clean() User {
	return &UserDTO{Id: u.Id, Name: u.Name, Username: u.Username}
}

func (u UserModel) GetId() int {
	return u.Id
}

func (u UserModel) GetName() string {
	return u.Name
}

func (u UserModel) GetPassword() string {
	return u.Password
}

func (u UserModel) GetUsername() string {
	return u.Username
}

func (u *UserModel) SetId(id int) User {
	u.Id = id
	return u
}

func (u *UserModel) SetName(name string) User {
	u.Name = name
	return u
}

func (u *UserModel) SetPassword(password string) User {
	u.Password = password
	return u
}

func (u *UserModel) SetUsername(username string) User {
	u.Username = username
	return u
}

func (u *UserModel) UpdateFromAnother(other User) {
	if other.GetName() != "" {
		u.Name = other.GetName()
	}

	if other.GetUsername() != "" {
		u.Username = other.GetUsername()
	}
}
