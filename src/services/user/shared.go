package user

import (
	"github.com/ppcamp/go-strings/random"
	"golang.org/x/crypto/bcrypt"
)

// newSecret generates a new string with the fixed lenght of 30.
func newSecret() string {
	return random.String(30)
}

// hashPassword takes a password and apply a bcrypt algorithm with salt 14
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
