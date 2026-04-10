package client

import (
	"errors"
	"fmt"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Type       string `json:"type"`
	Code       string `json:"code"`
	RequestID  string `json:"requestId"`
	Message    string `json:"message,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d (%s): %s [request_id=%s]", e.StatusCode, e.Code, e.Message, e.RequestID)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 404
}

func IsConflict(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 409
}
