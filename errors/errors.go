package errors

import (
	"encoding/json"
	"fmt"
)

type PingenError struct {
	Message    string            `json:"message"`
	JSONBody   interface{}       `json:"json_body"`
	StatusCode int               `json:"-"`
	Headers    map[string]string `json:"-"`
	RequestID  string            `json:"-"`
}

func NewPingenError(message string, body string, statusCode int, headers map[string]string) *PingenError {
	pErr := &PingenError{
		Message:    message,
		StatusCode: statusCode,
		Headers:    headers,
	}

	if body != "" {
		var jsonBody interface{}
		if err := json.Unmarshal([]byte(body), &jsonBody); err == nil {
			pErr.JSONBody = jsonBody
		}
	}

	if headers != nil {
		if requestID, ok := headers["X-Request-Id"]; ok {
			pErr.RequestID = requestID
		}
	}

	return pErr
}

func (e *PingenError) Error() string {
	return fmt.Sprintf("PingenError: %s (Status Code: %d, Request ID: %s)", e.Message, e.StatusCode, e.RequestID)
}

type AuthenticationError struct {
	PingenError
}

func NewAuthenticationError(message string, body string, statusCode int, headers map[string]string) *AuthenticationError {
	baseError := NewPingenError(message, body, statusCode, headers)
	return &AuthenticationError{*baseError}
}

type WebhookSignatureException struct {
	Message string
}

func NewWebhookSignatureException(message string) *WebhookSignatureException {
	return &WebhookSignatureException{Message: message}
}

func (e *WebhookSignatureException) Error() string {
	return fmt.Sprintf("WebhookSignatureException: %s", e.Message)
}
