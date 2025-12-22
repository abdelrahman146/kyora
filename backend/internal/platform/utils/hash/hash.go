package hash

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Password(val string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func ValidatePassword(val, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(val))
	return err == nil
}

func Make(signature string) string {
	hash := sha256.Sum256([]byte(signature))
	return fmt.Sprintf("%x", hash)
}

func Validate(signature, hash string) bool {
	return Make(signature) == hash
}
