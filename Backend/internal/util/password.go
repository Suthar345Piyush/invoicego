// password utility code

package util

import (
	"golang.org/x/crypto/bcrypt"
)

// password hashing function

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err

}

// checking the password

// lowercase function - private function
// Uppercase lettered function - exported function - which can be exported and used at any other place

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
