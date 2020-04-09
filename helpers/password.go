package helpers

import (
	"crypto/sha1"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

var (
	bcryptHashCost = 13
)

// Hash :
func Hash(key string) string {
	byteKey := []byte(key)
	h := sha1.New()

	h.Write([]byte(byteKey))

	return hex.EncodeToString(h.Sum(nil))
}

// HashAndSalt :
func HashAndSalt(pwd string) (string, error) {
	bytePwd := []byte(pwd)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcryptHashCost)
	if nil != err {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// ComparePasswords :
func ComparePasswords(hashedPwd string, plainPwd string) error {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	bytePwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if nil != err {
		return err
	}

	return nil
}
