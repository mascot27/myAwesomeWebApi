package models

import "errors"

type tPassword string

type User struct {
	Name     string    `json:"username"`
	Password tPassword `json:"password"`
	Uuid     string    `json:"uuid"`
}

// Marshaler ignores the field value completely.
func (tPassword) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

/**
Mock data:
TODO: implement db
TODO: hash passwords
 */
var Users = []User{
	User{
		Name:     "user1",
		Password: "password",
	},
	User{
		Name:     "user2",
		Password: "password",
	},
}

func IsValidCredentials(username string, password tPassword) (string, error) {
	for _, user := range Users {
		if user.Name == user.Name && user.Password == password {
			return user.Uuid, nil
		}
	}
	return "", errors.New("invalid credentials")
}

func GetUuidForName(name string) (string, error) {

	for _, elem := range Users {
		if elem.Name == name {
			return elem.Uuid, nil
		}
	}

	return "", errors.New("invalid name")
}

func HashPassword(plaintext string) tPassword {
	//TODO: implement me
	return tPassword(plaintext)
}
