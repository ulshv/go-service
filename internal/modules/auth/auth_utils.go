package auth

import (
	"fmt"
	"strconv"
	"strings"
)

func hashPassword(password string) string {
	// todo, mock implementation
	return password
}

func validatePassword(password, hash string) bool {
	// todo, mock implementation
	return password == hash
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
