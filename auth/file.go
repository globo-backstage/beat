package auth

import (
	"fmt"
	_ "github.com/backstage/beat/config"
	"github.com/spf13/viper"
	"net/http"
)

// FileAuthentication implements the auth.Authable interface.
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

type FileAuthentication struct{}

// FileUser implements the auth.User interface.
type FileUser struct {
	TokenEmail string `mapstructure:"email"`
}

var (
	TokensConfigPath = "auth.tokens"
	NilFileUser      = FileUser{}
)

// NewFileAuthentication return a new FileAuthentication
func NewFileAuthentication() *FileAuthentication {
	return &FileAuthentication{}
}

// GetUser implements auth.Authable interface.
func (authenticaton *FileAuthentication) GetUser(header *http.Header) User {
	var user FileUser
	token := header.Get("Token")

	if token == "" {
		return nil
	}

	viper.UnmarshalKey(fmt.Sprintf("%s.%s", TokensConfigPath, token), &user)

	if user == NilFileUser {
		return nil
	}

	return user
}

// Email implements auth.User interface.
func (user FileUser) Email() string {
	return user.TokenEmail
}
