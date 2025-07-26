package auth

import "crypto/sha256"

func EncodePassword(pass string) []byte{
	hash := sha256.New()

	hash.Write([]byte(pass))
	hashed := hash.Sum(nil)
	return hashed
}

func CheckPassword(attempt string, secret string) bool {
	correct := false
	hash := sha256.New()
	hash.Write([]byte(attempt))
	test := hash.Sum(nil)
	if string(test) == secret {
		correct = true
	}

	return correct
}