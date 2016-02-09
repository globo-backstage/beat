package auth

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	auths = map[string]RegisterFunc{}
)

// User is the basic interface to identify the logged user provided by each
// Authable implementation.
type User interface {
	Email() string
}

type RegisterFunc func() (Authable, error)

// Authable is the interface that provides all capacity to handle autenticated
// and authorized transations.
//
// GetUser returns the current user based on http header for a transaction.
type Authable interface {
	GetUser(*http.Header) User
}

// Register inserts a implementation of `Authable` in the register, is useful
// to auto discover implementations and change it without changing the code.
func Register(name string, fn RegisterFunc) {
	auths[name] = fn
}

// New returns a implementation of `Authable` found in the register, if not found
// return an error.
func New(name string) (Authable, error) {
	fn := auths[name]
	if fn == nil {
		return nil, ErrNotFound{name: name}
	}
	db, err := fn()

	if err != nil {
		return nil, authError{name: name, originalErr: err}
	}

	return db, nil
}

type ErrNotFound struct {
	name string
}

func (a ErrNotFound) Error() string {
	availableAuths := make([]string, 0, len(auths))
	for auth := range auths {
		availableAuths = append(availableAuths, auth)
	}

	return fmt.Sprintf(`Authentication "%s" not found, are available: %s.`, a.name, strings.Join(availableAuths, ", "))
}

type authError struct {
	name        string
	originalErr error
}

func (a authError) Error() string {
	return fmt.Sprintf(`[authentication][%s] %s`, a.name, a.originalErr.Error())
}
