package domain

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/squadcast_assignment/internal/serializer"
)

type Error struct {
	err error
}

func NewError(msg string) *Error {
	e := Error{}
	e.err = errors.New(msg)
	return &e
}

func (e *Error) Error() string {
	return e.err.Error()
}

func ErrToJSON(e error, code int) []byte {
	errResponse := serializer.Error{
		Status: http.StatusText(code),
		Error:  e.Error(),
	}
	errJSON, marshalErr := json.Marshal(errResponse)
	if marshalErr != nil {
		panic("code error: json marshalling failed for internally constructed struct")
	}
	return errJSON
}
