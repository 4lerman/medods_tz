package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func HashToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func ComparePasswords(hashedPassword string, to_compare []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), to_compare)
	return err == nil
}

func CompareTokens(hashedToken string, to_compare []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), to_compare)
	return err == nil
}