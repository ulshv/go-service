package auth

import (
	"fmt"
	"strconv"
	"strings"

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

func generateToken(userId int) string {
	// todo, mock implementation
	return fmt.Sprintf("token-%v", userId)
}

func validateToken(token string, userId int) (error, bool) {
	// todo, mock implementation
	idStr := strings.TrimPrefix(token, "token-")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err, false
	}
	return nil, id == userId
}
