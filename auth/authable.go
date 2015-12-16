package auth

import (
	"net/http"
)

// User is the basic interface to identify the logged user provided by each
// Authable implementation.
type User interface {
	Email() string
}

// Authable is the interface that provides all capacity to handle autenticated
// and authorized transations.
//
// GetUser returns the current user based on http header for a transaction.
type Authable interface {
	GetUser(*http.Header) User
}
