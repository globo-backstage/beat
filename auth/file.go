package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"gopkg.in/yaml.v2"
)


type FileAuthentication struct {
	tokens map[string]interface{}
}

type FileUser struct {
	data map[interface{}]interface{}
}


func NewFileAuthentication(tokensPath string) (Authable, error) {
	data, err := ioutil.ReadFile(tokensPath)

	if (err != nil) {
		return nil, err
	}

	authentication := &FileAuthentication{}
	err = yaml.Unmarshal(data, &authentication.tokens)

	fmt.Println(authentication.tokens)
	return authentication, nil
	
}

func (authenticaton *FileAuthentication) GetUser(header *http.Header) User {
	token := header.Get("Token")
	tokenData := authenticaton.tokens[token]
	
	if tokenData == nil {
		return nil
	}

	data, ok := tokenData.(map[interface{}]interface{})

	if (!ok) {
		return nil
	}
	
	return &FileUser{data: data};
}

func (user *FileUser) Email() string  {
	email, ok := user.data["email"].(string)
	if (ok) {
		return email
	} else {
		return ""
	}
}

