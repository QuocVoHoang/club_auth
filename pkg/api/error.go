package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Resp *http.Response
	Body []byte
}

var errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) StatusCode() int {
	return e.Resp.StatusCode
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http error: %s %s", e.Resp.Status, e.Body)
}

// HandleHTTPError converts HTTPError to appropriate domain errors
func HandleHTTPError(err error, serviceName string) error {
	var httpErr HTTPError
	if !errors.As(err, &httpErr) {
		return fmt.Errorf("%s call failed: %w", serviceName, err)
	}

	if json.Unmarshal(httpErr.Body, &errorResponse) == nil && errorResponse.Message != "" {
		// Use the message from the service response
		return fmt.Errorf("%s service error: %s", serviceName, errorResponse.Message)
	}

	// Fallback to HTTPError's built-in error message
	return fmt.Errorf("%s service error (HTTP %d): %s", serviceName, httpErr.StatusCode(), httpErr.Error())
}
