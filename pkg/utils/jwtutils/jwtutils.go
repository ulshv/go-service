package jwtutils

import (
	"fmt"
	"log/slog"
	"time"

	jwt_mod "github.com/golang-jwt/jwt/v5"
	"github.com/ulshv/go-service/pkg/logs"
	"github.com/ulshv/go-service/pkg/utils/envutils"
)

type TokenType string

type Claims struct {
	UserId    int       `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	jwt_mod.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Jwt struct {
	accessTokenSecret    []byte
	refreshTokenSecret   []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	logger               *slog.Logger
}

func NewJWT() *Jwt {
	envutils.LoadEnvFiles()

	return &Jwt{
		accessTokenSecret:    []byte(envutils.RequireEnv("JWT_ACCESS_TOKEN_SECRET")),
		refreshTokenSecret:   []byte(envutils.RequireEnv("JWT_REFRESH_TOKEN_SECRET")),
		accessTokenDuration:  15 * time.Minute,
		refreshTokenDuration: 7 * 24 * time.Hour,
		logger:               logs.NewLogger("JWTUtils"),
	}
}

var jwtLogger = logs.NewLogger("JWTUtils")

func (j *Jwt) GenerateTokenPair(userId int) (TokenPair, error) {
	// Generate access token
	accessToken, err := j.generateToken(userId, AccessToken, j.accessTokenSecret, j.accessTokenDuration)
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to generate access token: %v", err)
	}

	// Generate refresh token
	refreshToken, err := j.generateToken(userId, RefreshToken, j.refreshTokenSecret, j.refreshTokenDuration)
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *Jwt) ValidateAccessToken(tokenString string) (*Claims, error) {
	return j.validateToken(tokenString, AccessToken, j.accessTokenSecret)
}

func (j *Jwt) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return j.validateToken(tokenString, RefreshToken, j.refreshTokenSecret)
}

func (j *Jwt) RefreshTokenPair(refreshToken string) (TokenPair, error) {
	// Parse the refresh token
	token, err := jwt_mod.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt_mod.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt_mod.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.refreshTokenSecret, nil
	})

	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to parse refresh token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return TokenPair{}, fmt.Errorf("invalid refresh token claims")
	}

	// Validate token type
	if claims.TokenType != RefreshToken {
		return TokenPair{}, fmt.Errorf("invalid token type")
	}

	// Generate new token pair
	return j.GenerateTokenPair(claims.UserId)
}

func (j *Jwt) generateToken(userId int, tokenType TokenType, secret []byte, duration time.Duration) (string, error) {
	claims := Claims{
		UserId:    userId,
		TokenType: tokenType,
		RegisteredClaims: jwt_mod.RegisteredClaims{
			ExpiresAt: jwt_mod.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt_mod.NewNumericDate(time.Now()),
		},
	}

	token := jwt_mod.NewWithClaims(jwt_mod.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		jwtLogger.Error("failed to generate token", "err", err)
		return "", err
	}

	return signedToken, nil
}

func (j *Jwt) validateToken(tokenString string, tokenType TokenType, secret []byte) (*Claims, error) {
	token, err := jwt_mod.ParseWithClaims(tokenString, &Claims{}, func(token *jwt_mod.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt_mod.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate token type
	if claims.TokenType != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	// // Validate user ID
	// if claims.UserId != userId {
	// 	return fmt.Errorf("token doesn't belong to user")
	// }

	return claims, nil
}
