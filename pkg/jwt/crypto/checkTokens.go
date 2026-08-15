package jwt

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Error("BCRYPT MATCH FAILED", 
			"error", err, 
			"hash_received", hash, 
			"password_len", len(password),
		)
		return false
	}
	return true
}