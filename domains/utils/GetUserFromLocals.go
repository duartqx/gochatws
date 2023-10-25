package utils

import (
	"encoding/json"
	"fmt"

	i "github.com/duartqx/gochatws/core/interfaces"
	m "github.com/duartqx/gochatws/domains/models"
)

func GetUserFromLocals(localUser interface{}) (i.User, error) {
	if localUser == nil {
		return nil, fmt.Errorf("User not found on Locals\n")
	}
	userBytes, err := json.Marshal(localUser)
	if err != nil {
		return nil, err
	}

	userStruct := &m.UserModel{}
	err = json.Unmarshal(userBytes, userStruct)
	if err != nil {
		return nil, err
	}
	return userStruct, nil
}
