package utils

import (
	"encoding/json"
	"fmt"

	u "github.com/duartqx/gochatws/domains/entities/user"
)

func GetUserFromLocals(localUser interface{}) (u.User, error) {
	if localUser == nil {
		return nil, fmt.Errorf("User not found on Locals\n")
	}
	userBytes, err := json.Marshal(localUser)
	if err != nil {
		return nil, err
	}

	userStruct := &u.UserModel{}
	err = json.Unmarshal(userBytes, userStruct)
	if err != nil {
		return nil, err
	}
	return userStruct, nil
}
