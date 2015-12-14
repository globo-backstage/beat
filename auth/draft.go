package auth;

import "net/http"

type DraftAuthentication struct {}
type DraftUser struct {}

func (d *DraftAuthentication) GetUser(*http.Header) User  {
	return &DraftUser{};
}

func (u *DraftUser) Email () string  {
	return "anonymous@example.net"
}
