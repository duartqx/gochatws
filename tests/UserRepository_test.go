package tests

import (
	"log"
	"testing"

	i "github.com/duartqx/gochatws/core/interfaces"
	u "github.com/duartqx/gochatws/domains/auth/users"
)

func TestUserRepositoryImplementsRepository(t *testing.T) {

	var userRepository i.Repository = u.UserRepository{}

	if _, ok := userRepository.(i.Repository); !ok {
		t.Errorf("UserRepository does not Implements Repository!")
	}
	log.Println("UserRepository does Implements Repository!")

}
