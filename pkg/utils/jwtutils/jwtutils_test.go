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
		_, err = jwt.ValidateAccessToken(tokenPair.AccessToken)
		if err != nil {
			t.Errorf("ValidateAccessToken error: %v", err)
		}
		_, err = jwt.ValidateRefreshToken(tokenPair.RefreshToken)
		if err != nil {
			t.Errorf("ValidateRefreshToken error: %v", err)
		}
	})
}
