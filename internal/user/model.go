package user

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"email"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
}

func (u *User) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("Marshaling new user error: %v", err)
	}
	return bytes, nil
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	ID       string `json:"_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
