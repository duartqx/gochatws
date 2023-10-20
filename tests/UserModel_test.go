package tests

import (
	"log"
	"testing"

	i "github.com/duartqx/gochatws/core/interfaces"
	u "github.com/duartqx/gochatws/domains/auth/users"
)

func TestUserModel(t *testing.T) {

	var user i.User = &u.UserModel{}
	t.Log(user)

	if _, ok := user.(*u.UserModel); !ok {
		t.Errorf("UserModel does not implements Model!")
	}

	log.Println("UserModel implements Model!")
}
