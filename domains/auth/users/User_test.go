package users

import (
	"encoding/json"
	"testing"

	i "github.com/duartqx/gochatws/core/interfaces"
	"github.com/go-playground/validator/v10"
)

func TestUserModel(t *testing.T) {

	v := validator.New()

	var user i.User = &UserModel{}

	if _, ok := user.(*UserModel); !ok {
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

	parsedOther, err, validationErrs := UserModel{}.ParseAndValidate(parser, v)
	if validationErrs != nil || err != nil {
		t.Errorf("FAILED: Could not Validate UserModel: %s\n", err)
	}
	t.Logf("PASSED: UserModel validated\n")

	if parsedOther.Username != "other@teste.com" ||
		parsedOther.Name != "Other" ||
		parsedOther.Password != "password" {

		t.Errorf(
			"FAILED: Did not parsed correctly UserModel: %s, %s, %s",
			parsedOther.Name,
			parsedOther.Username,
			parsedOther.Password,
		)
	}
	t.Logf("PASSED: UserModel parsed correctly\n")

	userFromOther := UserModel{}

	userFromOther.UpdateFromAnother(parsedOther)

	if userFromOther.Username != parsedOther.Username ||
		userFromOther.Name != parsedOther.Name {
		t.Errorf("FAILED: UpdateFromAnother did not update a User")
	}
	t.Logf("PASSED: UpdateFromAnother updates UserModel correctly\n")
}

func TestUserRepositoryImplementsRepository(t *testing.T) {

	var userRepository i.Repository = UserRepository{}

	if _, ok := userRepository.(i.Repository); !ok {
		t.Errorf("FAILED: UserRepository does not Implements Repository!\n")
	}
	t.Logf("PASSED: UserRepository implements Repository interface.\n")
}
