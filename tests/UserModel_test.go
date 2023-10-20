package tests

import (
	"encoding/json"
	"testing"

	i "github.com/duartqx/gochatws/core/interfaces"
	u "github.com/duartqx/gochatws/domains/auth/users"
	"github.com/go-playground/validator/v10"
)

func TestUserModel(t *testing.T) {

	v := validator.New()

	var user i.User = &u.UserModel{}

	if _, ok := user.(*u.UserModel); !ok {
		t.Errorf("UserModel does not implements Model!")
	}

	parser := func(out interface{}) error {
		jsonBody := []byte(`{
			"name": "Other",
			"username": "other@teste.com",
			"password": "password"
		}`)
		return json.Unmarshal(jsonBody, &out)
	}

	parsedOther, err, validationErrs := u.UserModel{}.ParseAndValidate(parser, v)
	if validationErrs != nil || err != nil {
		t.Errorf("Could not Validate UserModel: %s", err)
	}

	if parsedOther.Username != "other@teste.com" ||
		parsedOther.Name != "Other" ||
		parsedOther.Password != "password" {

		t.Errorf(
			"Did not parsed correctly UserModel: %s, %s, %s",
			parsedOther.Name,
			parsedOther.Username,
			parsedOther.Password,
		)
	}

	userFromOther := u.UserModel{}

	userFromOther.UpdateFromAnother(parsedOther)

	if userFromOther.Username != parsedOther.Username ||
		userFromOther.Name != parsedOther.Name {
		t.Errorf("UpdateFromAnother did not update a User")
	}
}
