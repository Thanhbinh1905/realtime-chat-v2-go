package errors

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrTokenGeneration = errors.New("error generating token")
var ErrTokenInvalid = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")
var ErrInternalServer = errors.New("internal server error")
var ErrBadRequest = errors.New("bad request")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrTooManyRequests = errors.New("too many requests")
var ErrServiceUnavailable = errors.New("service unavailable")
var ErrGatewayTimeout = errors.New("gateway timeout")
var ErrInvalidUserID = errors.New("invalid user ID")
var ErrInvalidEmail = errors.New("invalid email format")
var ErrEmailExists = errors.New("email already exists")
var ErrInternalServerError = errors.New("internal server error")
