package errors

import (
	"encoding/json"
)

// DetailedError is a interface to inform the user about the error ocorred.
type DetailedError interface {
	StatusCode() int
	Error() string // to implement the common error interface
}

// WrappedError implements DetailedError, created only for wraps common golang
// errors in this interface.
type WrappedError struct {
	originalError error
	statusCode    int
}

// Wraps return new DetailedError based on a common golang error.
func Wraps(err error, statusCode int) DetailedError {
	return &WrappedError{originalError: err, statusCode: statusCode}
}

func (w *WrappedError) Error() string {
	return w.originalError.Error()
}

// StatusCode return the code to be used in http response.
func (w *WrappedError) StatusCode() int {
	return w.statusCode
}

// MarshalJSON returns the json format or DetailedError
func (w *WrappedError) MarshalJSON() ([]byte, error) {
	errorDescription := map[string]interface{}{
		"_all": []string{
			w.Error(),
		},
	}
	data := map[string]interface{}{
		"errors": []interface{}{
			errorDescription,
		},
	}
	return json.Marshal(&data)
}
