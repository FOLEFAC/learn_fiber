package utils

import "golang.org/x/crypto/bcrypt"

// NormalizePassword func for a returning the users input as a byte slice.
func NormalizePassword(p string) []byte {
	return []byte(p)
}

// GeneratePassword func for a making hash & salt with user password.
func GeneratePassword(p string) string {
	// Normalize password from string to []byte.
	bytePwd := NormalizePassword(p)

	// MinCost is just an integer constant provided by the bcrypt package
	// along with DefaultCost & MaxCost. The cost can be any value
	// you want provided it isn't lower than the MinCost (4).
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return err.Error()
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it.
	return string(hash)
}
