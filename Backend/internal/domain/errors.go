// for error showing from user side

package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidInput         = errors.New("invalid input")
	ErrInternalServer       = errors.New("internal server error")
	ErrEmailNotVerified     = errors.New("email not verified")
	ErrInvoiceLimitExceeded = errors.New("monthly invoice limit exceeded")
)
