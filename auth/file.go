package auth

import (
	"io/ioutil"
	"net/http"
	"gopkg.in/yaml.v2"
)

type FileAuthentication struct {
	Tokens map[string]*FileUser `yaml:"tokens"`
}

type FileUser struct {
	TokenEmail string `yaml:"email"`
}

func NewFileAuthentication(tokensPath string) (Authable, error) {
	data, err := ioutil.ReadFile(tokensPath)

	if (err != nil) {
		return nil, err
	}

	authentication := &FileAuthentication{}
	err = yaml.Unmarshal(data, authentication)

	if (err != nil) {
		return nil, err
	}

	return authentication, nil

}

func (authenticaton *FileAuthentication) GetUser(header *http.Header) User {
	token := header.Get("Token")
	return authenticaton.Tokens[token]
}

func (user FileUser) Email() string  {
	return user.TokenEmail
}
