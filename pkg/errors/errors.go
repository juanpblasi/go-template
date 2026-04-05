package errors

import (
	"errors"
	"fmt"
)

// Definimos errores que representen códigos HTTP comunes sin acoplar a HTTP framework
var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidRequest = errors.New("invalid request format")
	ErrInternalError  = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized access")
)

type CustomError struct {
	Type    error
	Message string
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("%v: %s", c.Type, c.Message)
}

func New(errType error, msg string) error {
	return &CustomError{
		Type:    errType,
		Message: msg,
	}
}

// IsType checks if error is of a specific CustomError Type
func IsType(err error, targetType error) bool {
	var customErr *CustomError
	if errors.As(err, &customErr) {
		return errors.Is(customErr.Type, targetType)
	}
	return errors.Is(err, targetType)
}
