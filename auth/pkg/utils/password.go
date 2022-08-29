package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword returns new bcrypt hashed password.
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// ComparePassword compares a bcrypt hashed password with its possible plaintext equivalent.
//
// Returns nil on success, or an error on failure.
func ComparePassword(hashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
