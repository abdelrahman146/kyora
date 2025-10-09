package utils

import "golang.org/x/crypto/bcrypt"

type hashHelper struct{}

func (hashHelper) Make(val string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (hashHelper) Validate(val, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(val))
	return err == nil
}

var Hash = hashHelper{}
