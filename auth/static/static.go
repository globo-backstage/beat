package static

import (
	"fmt"
	"github.com/backstage/beat/auth"
	"github.com/spf13/viper"
	"net/http"
)

// StaticAuthentication implements the auth.Authable interface.
//
// This Authable is a simple authentication based on a yaml file
// that contains all tokens allowed to perform a write operation.
//
// An example of yaml file:
// auth:
//   tokens:
//     example1:
//       email: admin@example.net
//
//     example2:
//       email: guest@example.net
//
//
// For each token is allowed to make a request like
// curl -H "Token: example1" http://myserver/api/collection

type StaticAuthentication struct{}

// StaticUser implements the auth.User interface.
type StaticUser struct {
	TokenEmail string `mapstructure:"email"`
}

var (
	TokensConfigPath = "auth.tokens"
	NilStaticUser    = StaticUser{}
)

func init() {
	auth.Register("static", func() (auth.Authable, error) {
		return NewStaticAuthentication(), nil
	})
}

// NewStaticAuthentication return a new StaticAuthentication
func NewStaticAuthentication() *StaticAuthentication {
	return &StaticAuthentication{}
}

// GetUser implements auth.Authable interface.
func (authenticaton *StaticAuthentication) GetUser(header *http.Header) auth.User {
	var user StaticUser
	token := header.Get("Token")

	if token == "" {
		return nil
	}

	viper.UnmarshalKey(fmt.Sprintf("%s.%s", TokensConfigPath, token), &user)

	if user == NilStaticUser {
		return nil
	}

	return user
}

// Email implements auth.User interface.
func (user StaticUser) Email() string {
	return user.TokenEmail
}
