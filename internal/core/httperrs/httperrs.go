package httperrs

import "errors"

const (
	// JSON
	ErrCodeMalformedJSON = "ERR_MALFORMED_JSON"
	// JWT
	ErrCodeAccessTokenExpired  = "ERR_ACCESS_TOKEN_EXPIRED"
	ErrCodeAccessTokenInvalid  = "ERR_ACCESS_TOKEN_INVALID"
	ErrCodeRefreshTokenExpired = "ERR_REFRESH_TOKEN_EXPIRED"
	ErrCodeRefreshTokenInvalid = "ERR_REFRESH_TOKEN_INVALID"
	// Common errors
	ErrCodeUnautorized = "ERR_UNAUTORIZED"
	ErrCodeInternal    = "ERR_INTERNAL_SERVER_ERROR"
	ErrCodeUnknown     = "ERR_UNKNOWN_ERR"
	// Auth module
	ErrCodeEmailTaken             = "ERR_EMAIL_TAKEN"
	ErrCodeInvalidEmailOrPassword = "ERR_INVALID_EMAIL_OR_PASSWORD"
)

var ErrUnauthorized = errors.New("unauthorized")
