package jwtutils

import (
	"testing"
)

func TestJWT(t *testing.T) {
	t.Run("TestTokenGenerationAndValidation", func(t *testing.T) {
		jwt := NewJWT()
		userId := 1
		tokenPair, err := jwt.GenerateTokenPair(userId)
		if err != nil {
			t.Errorf("GenerateTokenPair error: %v", err)
		}
		err = jwt.ValidateAccessToken(tokenPair.AccessToken, userId)
		if err != nil {
			t.Errorf("ValidateAccessToken error: %v", err)
		}
		err = jwt.ValidateRefreshToken(tokenPair.RefreshToken, userId)
		if err != nil {
			t.Errorf("ValidateRefreshToken error: %v", err)
		}
	})
}
