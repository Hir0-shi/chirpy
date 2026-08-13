package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	// argon2id will choose sane defaults for this library version.
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
