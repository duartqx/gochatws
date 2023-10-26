package tests

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"

	i "github.com/duartqx/gochatws/core/interfaces"
	m "github.com/duartqx/gochatws/domains/models"
	r "github.com/duartqx/gochatws/domains/repositories"
)

func TestUserModel(t *testing.T) {

	v := validator.New()

	var user i.User = &m.UserModel{}

	if _, ok := user.(*m.UserModel); !ok {
		t.Errorf("FAILED: UserModel does not implements the User interface!\n")
	}
	t.Logf("PASSED: UserModel implements User interface\n")

	parser := func(out interface{}) error {
		jsonBody := []byte(`{
			"name": "Other",
			"username": "other@teste.com",
			"password": "password"
		}`)
		return json.Unmarshal(jsonBody, &out)
	}

	parsedOther := &m.UserModel{}
	if err := parser(parsedOther); err != nil {
		t.Errorf("FAILED: User was not able to parse!\n")
	}
	t.Logf("PASSED: User parsed correctly!\n")

	if err := v.Struct(parsedOther); err != nil {
		t.Errorf("FAILED: User did not validate!\n")
	}
	t.Logf("PASSED: User validated correctly!\n")

	if parsedOther.GetUsername() != "other@teste.com" ||
		parsedOther.GetName() != "Other" ||
		parsedOther.GetPassword() != "password" {

		t.Errorf(
			"FAILED: Did not parsed correctly UserModel: %s, %s, %s",
			parsedOther.GetPassword(),
			parsedOther.GetUsername(),
			parsedOther.GetName(),
		)
	}
	t.Logf("PASSED: UserModel parsed correctly\n")

	userFromOther := m.UserModel{}

	userFromOther.UpdateFromAnother(parsedOther)

	if userFromOther.GetUsername() != parsedOther.GetUsername() ||
		userFromOther.GetName() != parsedOther.GetName() {
		t.Errorf("FAILED: UpdateFromAnother did not update a User")
	}
	t.Logf("PASSED: UpdateFromAnother updates UserModel correctly\n")
}

func TestUserRepositoryImplementsRepository(t *testing.T) {

	var userRepository i.UserRepository = r.UserRepository{}

	if _, ok := userRepository.(i.UserRepository); !ok {
		t.Errorf("FAILED: UserRepository does not Implements Repository!\n")
	}
	t.Logf("PASSED: UserRepository implements Repository interface.\n")
}