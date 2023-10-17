package users

type UserModel struct {
	Id       int    `db:"id" json:"user_id"`
	Username string `db:"username" json:"username" validate:"email,required"`
	Name     string `db:"name" json:"name" validate:"required"`
	Password string `db:"password" json:"-"`
}

func (u *UserModel) UpdateFromAnother(another *UserModel) {
	if another.Name != "" {
		u.Name = another.Name
	}

	if another.Username != "" {
		u.Username = another.Username
	}
}
