package transaction

import (
	"encoding/json"
	"github.com/backstage/beat/errors"
	"github.com/dimfeld/httptreemux"
	"net/http"
)

type TransactionHandler func(*Transaction)
type Transaction struct {
	writer http.ResponseWriter
	Id     string
	Params map[string]string
	Req    *http.Request
}

func (t *Transaction) WriteError(err errors.Error) {
	t.writer.WriteHeader(err.StatusCode())
	json.NewEncoder(t.writer).Encode(err)
}

func (t *Transaction) WriteResultWithStatusCode(statusCode int, result interface{}) {
	t.writer.WriteHeader(statusCode)
	json.NewEncoder(t.writer).Encode(result)
}

func (t *Transaction) WriteResult(result interface{}) {
	json.NewEncoder(t.writer).Encode(result)
}

func (t *Transaction) NoResultWithStatusCode(statusCode int) {
	t.writer.WriteHeader(statusCode)
}

func Handle(handler TransactionHandler) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		t := &Transaction{
			Req:    r,
			writer: w,
			Params: ps,
		}
		handler(t)
	}
}
