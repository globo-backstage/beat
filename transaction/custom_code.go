package transaction

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/backstage/beat/errors"
)

type CustomCodeRequestBody struct {
	Hook        string   `json:"hook"`
	CustomCodes []string `json:"customCodes"`
	Req         struct {
		IsUpdate bool `json:"isUpdate"`
		IsCreate bool `json:"isCreate"`
	} `json:"req"`
}

func RunBeforeSave(t *Transaction) errors.Error {
	buf := bytes.NewBuffer([]byte{})
	reqBody := &CustomCodeRequestBody{
		Hook:        "beforeSave",
		CustomCodes: []string{"scheduler"},
	}
	err := json.NewEncoder(buf).Encode(reqBody)
	if err != nil {
		return errors.Wraps(err, http.StatusInternalServerError)
	}
	res, err := http.Post("http://localhost:8100", "application/json", buf)
	if err != nil {
		return errors.Wraps(err, http.StatusInternalServerError)
	}
	println(res)
	return nil
}
