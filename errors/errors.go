package errors

import (
	"encoding/json"
	originalErrors "errors"
	"fmt"
	"strings"
)

// Error is a interface to inform the user about the error ocorred.
type Error interface {
	StatusCode() int
	Error() string // to implement the common error interface
}

// wrappedError implements Error, created only for wraps common golang
// errors in this interface.
type wrappedError struct {
	originalError error
	statusCode    int
}

// Wraps return new Error based on a common golang error.
func Wraps(err error, statusCode int) Error {
	return &wrappedError{originalError: err, statusCode: statusCode}
}

func (w *wrappedError) Error() string {
	return w.originalError.Error()
}

// StatusCode return the code to be used in http response.
func (w *wrappedError) StatusCode() int {
	return w.statusCode
}

// MarshalJSON returns the json format or Error
func (w *wrappedError) MarshalJSON() ([]byte, error) {
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

// New returns new Error based on string.
func New(text string, statusCode int) Error {
	return Wraps(originalErrors.New(text), statusCode)
}

// Newf returns new Error based on string formated by fmt.Errorf.
func Newf(statusCode int, text string, params ...interface{}) Error {
	return Wraps(fmt.Errorf(text, params...), statusCode)
}

// ValidationError is another implementation of Error interface.
// used to group field errors in one.
type ValidationError struct {
	items []validateItemError
}

func (v *ValidationError) StatusCode() int {
	return 422
}

func (v *ValidationError) Error() string {
	parts := []string{}

	for _, item := range v.items {
		parts = append(
			parts,
			fmt.Sprintf("%s: %s", item.field, item.message),
		)
	}

	return strings.Join(parts, ", ")
}

func (v *ValidationError) Length() int {
	return len(v.items)
}

func (v *ValidationError) Put(field, message string) {
	item := validateItemError{field: field, message: message}
	v.items = append(v.items, item)
}

func (v *ValidationError) MarshalJSON() ([]byte, error) {
	errorParts := []interface{}{}
	for _, item := range v.items {
		errorParts = append(errorParts, map[string]interface{}{
			item.field: []string{
				item.message,
			},
		})
	}
	data := map[string]interface{}{
		"errors": errorParts,
	}
	return json.Marshal(&data)
}

type validateItemError struct {
	field   string
	message string
}
