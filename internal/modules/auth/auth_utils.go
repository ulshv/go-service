package auth

import (
	"fmt"

	"github.com/ulshv/go-service/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

var utilsLogger = logger.NewLogger("AuthUtils")

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedBytes), nil
}

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		utilsLogger.Debug("validatePassword error", "err", err)
	}
	return err == nil
}
