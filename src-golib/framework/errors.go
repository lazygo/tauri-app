package framework

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lazygo/lazygo/utils"
)

// Errors
var (
	ErrUnsupportedMediaType        = NewError(http.StatusUnsupportedMediaType)
	ErrNotFound                    = NewError(http.StatusNotFound)
	ErrUnauthorized                = NewError(http.StatusUnauthorized)
	ErrForbidden                   = NewError(http.StatusForbidden)
	ErrMethodNotAllowed            = NewError(http.StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewError(http.StatusRequestEntityTooLarge)
	ErrTooManyRequests             = NewError(http.StatusTooManyRequests)
	ErrBadRequest                  = NewError(http.StatusBadRequest)
	ErrBadGateway                  = NewError(http.StatusBadGateway)
	ErrInternalServerError         = NewError(http.StatusInternalServerError)
	ErrRequestTimeout              = NewError(http.StatusRequestTimeout)
	ErrServiceUnavailable          = NewError(http.StatusServiceUnavailable)
	ErrValidatorNotRegistered      = errors.New("validator not registered")
	ErrActionNotExists             = errors.New("action not exists")
	ErrInvalidRedirectCode         = errors.New("invalid redirect status code")
	ErrCookieNotFound              = errors.New("cookie not found")
	ErrInvalidListenerNetwork      = errors.New("invalid listener network")
)

// Error handlers
var (
	NotFoundHandler = func(c Context) error {
		return ErrNotFound
	}

	MethodNotAllowedHandler = func(c Context) error {
		return ErrMethodNotAllowed
	}
)

// Error represents an error that occurred while handling a request.
type Error struct {
	Code     int         `json:"code"`
	Errno    int         `json:"errno"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}

// NewError creates a new HTTPError instance.
func NewError(code int, message ...interface{}) *Error {
	he := &Error{Code: code, Errno: code, Message: http.StatusText(code)}
	switch len(message) {
	case 0:
	case 1:
		he.Message = message[0]
	case 2:
		he.Errno = utils.ToInt(message[0], -1)
		he.Message = message[1]
	default:
		he.Message = message[0]
	}
	return he
}

// Error makes it compatible with `error` interface.
func (he *Error) Error() string {
	if he.Internal == nil {
		return fmt.Sprintf("code=%d, errno=%d, message=%v", he.Code, he.Errno, he.Message)
	}
	return fmt.Sprintf("code=%d, errno=%d, message=%v, internal=%v", he.Code, he.Errno, he.Message, he.Internal)
}

// SetInternal sets error to HTTPError.Internal
func (he *Error) SetInternal(err error) *Error {
	he.Internal = err
	return he
}

// Unwrap satisfies the Go 1.13 error wrapper interface.
func (he *Error) Unwrap() error {
	return he.Internal
}
