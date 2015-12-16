package auth

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

// FileAuthentication implements the auth.Authable interface.
//
// This Authable is a simple authentication based on a yaml file
// that contains all tokens allowed to perform a write operation.
//
// An example of yaml file:
// tokens:
//   example1:
//     email: admin@example.net
//
//   example2:
//     email: guest@example.net
//
//
// For each token is allowed to make a request like
// curl -H "Token: example1" http://myserver/api/collection
type FileAuthentication struct {
	Tokens map[string]*FileUser `yaml:"tokens"`
}

// FileUser implements the auth.User interface.
type FileUser struct {
	TokenEmail string `yaml:"email"`
}

// NewFileAuthentication return a new FileAuthentication by path that contains yaml file.
func NewFileAuthentication(tokensPath string) (Authable, error) {
	data, err := ioutil.ReadFile(tokensPath)

	if err != nil {
		return nil, err
	}

	authentication := &FileAuthentication{}
	err = yaml.Unmarshal(data, authentication)

	if err != nil {
		return nil, err
	}

	return authentication, nil
}

// GetUser implements auth.Authable interface.
func (authenticaton *FileAuthentication) GetUser(header *http.Header) User {
	token := header.Get("Token")
	return authenticaton.Tokens[token]
}

// Email implements auth.User interface.
func (user FileUser) Email() string {
	return user.TokenEmail
}
