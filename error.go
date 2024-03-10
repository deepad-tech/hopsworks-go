package hopsworks

import (
	"fmt"
)

// APIError provides error information returned by the Hopsworks API.
type APIError struct {
	Code           any    `errorCode:"code,omitempty"`
	UserMessage    string `json:"usrMsg,omitempty"`
	ErrMessage     string `json:"errorMsg,omitempty"`
	HTTPStatusCode int    `json:"-"`
}

// RequestError provides information about generic request errors.
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

type ErrorResponse struct {
	Error *APIError
}

func (e *APIError) Error() string {
	if e.HTTPStatusCode > 0 {
		return fmt.Sprintf("error, status code: %d, message: %s", e.HTTPStatusCode, e.ErrMessage)
	}

	return e.ErrMessage
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("error, status code: %d, message: %s", e.HTTPStatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}
