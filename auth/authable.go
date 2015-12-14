package auth;

import (
	"net/http"
)

type User interface {
	Email() string
}

type Authable interface {
	GetUser(*http.Header) User
}
