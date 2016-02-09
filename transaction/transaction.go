package transaction

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	"github.com/dimfeld/httptreemux"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

const (
	TransactionHeader    = "Backstage-Transaction"
	MaxTransactionHeader = 22
	SlowTransactionWarn  = time.Millisecond * 100
)

type TransactionHandler func(*Transaction)
type Transaction struct {
	writer         http.ResponseWriter
	statusCode     int
	Id             string
	CollectionName string
	ItemSchema     *schemas.ItemSchema
	Params         map[string]string
	Req            *http.Request
	Log            *log.Entry
}

func (t *Transaction) WriteError(err errors.Error) {
	t.statusCode = err.StatusCode()
	t.writer.WriteHeader(err.StatusCode())
	json.NewEncoder(t.writer).Encode(err)
}

func (t *Transaction) WriteResultWithStatusCode(statusCode int, result interface{}) {
	t.statusCode = statusCode
	t.writer.WriteHeader(statusCode)
	json.NewEncoder(t.writer).Encode(result)
}

func (t *Transaction) WriteResult(result interface{}) {
	t.WriteResultWithStatusCode(http.StatusOK, result)
}

func (t *Transaction) NoResultWithStatusCode(statusCode int) {
	t.statusCode = statusCode
	t.writer.WriteHeader(statusCode)
}

func (t *Transaction) BaseUrl() string {
	host := t.Req.URL.Host
	scheme := t.Req.URL.Scheme

	if host == "" {
		host = t.Req.Host
	}

	if scheme == "" {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/api", scheme, host)
}

func Handle(handler TransactionHandler) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		start := time.Now()
		id := IdFromRequest(r)

		t := &Transaction{
			Id:     id,
			Req:    r,
			writer: w,
			Params: ps,
			Log: log.WithFields(log.Fields{
				"transaction": id,
			}),
		}

		handler(t)
		logTransaction(t, time.Since(start))
	}
}

func CollectionHandle(handler TransactionHandler) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		start := time.Now()
		collectionName := ps["collectionName"]
		id := IdFromRequest(r)

		t := &Transaction{
			Id:             id,
			CollectionName: collectionName,
			Req:            r,
			writer:         w,
			Params:         ps,
			Log: log.WithFields(log.Fields{
				"transaction":    id,
				"collectionName": collectionName,
			}),
		}

		handler(t)
		logTransaction(t, time.Since(start))
	}
}

func IdFromRequest(r *http.Request) string {
	header := r.Header.Get(TransactionHeader)
	if header == "" || len(header) > MaxTransactionHeader {
		header = base64.RawStdEncoding.EncodeToString(uuid.NewV4().Bytes())
	}
	return header
}

func logTransaction(t *Transaction, latency time.Duration) {
	msg := fmt.Sprintf(
		"%s %s %d %s", t.Req.Method, t.Req.URL.RequestURI(), t.statusCode,
		latency.String(),
	)

	switch {
	case t.statusCode >= http.StatusInternalServerError:
		t.Log.Error(msg)
	case t.statusCode >= http.StatusBadRequest || latency > SlowTransactionWarn:
		t.Log.Warn(msg)
	default:
		t.Log.Info(msg)
	}
}
