package helper

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	DefaultPasswordSaltSize int = 32
)

// HashPassword hashes the password based on a salt
func HashPassword(appsalt, usersalt, password string) string {
	h := hmac.New(sha256.New, []byte(appsalt))
	h.Write([]byte(fmt.Sprintf("%s%s", usersalt, password)))
	return hex.EncodeToString(h.Sum(nil))
}

// SaltAndHashPassword creates a salt and returns the salt + hashed password
func SaltAndHashPassword(appsalt, password string) (string, string) {
	// create a salt
	salt := make([]byte, DefaultPasswordSaltSize)
	rand.Read(salt)
	saltHex := hex.EncodeToString(salt)

	hash := HashPassword(appsalt, saltHex, password)
	return saltHex, hash
}
