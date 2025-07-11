package errors

import (
	"encoding/json"
	"testing"
)

func TestNewPingenError(t *testing.T) {
	t.Run("basic error creation", func(t *testing.T) {
		message := "Test error"
		statusCode := 400
		headers := map[string]string{
			"Content-Type": "application/json",
			"X-Request-Id": "req-123",
		}

		err := NewPingenError(message, "", statusCode, headers)

		if err.Message != message {
			t.Errorf("Expected message %s, got %s", message, err.Message)
		}
		if err.StatusCode != statusCode {
			t.Errorf("Expected status code %d, got %d", statusCode, err.StatusCode)
		}
		if err.RequestID != "req-123" {
			t.Errorf("Expected request ID 'req-123', got '%s'", err.RequestID)
		}
		if err.Headers["Content-Type"] != "application/json" {
			t.Errorf("Expected Content-Type header, got %v", err.Headers)
		}
	})

	t.Run("with valid JSON body", func(t *testing.T) {
		jsonBody := `{"error": "validation failed", "field": "email"}`
		err := NewPingenError("Test error", jsonBody, 422, nil)

		if err.JSONBody == nil {
			t.Error("Expected JSONBody to be parsed")
		}

		bodyMap, ok := err.JSONBody.(map[string]interface{})
		if !ok {
			t.Error("Expected JSONBody to be a map")
		}

		if bodyMap["error"] != "validation failed" {
			t.Errorf("Expected error field, got %v", bodyMap["error"])
		}
	})

	t.Run("with invalid JSON body", func(t *testing.T) {
		invalidJSON := `{"invalid": json}`
		err := NewPingenError("Test error", invalidJSON, 400, nil)

		if err.JSONBody != nil {
			t.Error("Expected JSONBody to be nil for invalid JSON")
		}
	})

	t.Run("with empty body", func(t *testing.T) {
		err := NewPingenError("Test error", "", 400, nil)

		if err.JSONBody != nil {
			t.Error("Expected JSONBody to be nil for empty body")
		}
	})

	t.Run("without X-Request-Id header", func(t *testing.T) {
		headers := map[string]string{"Content-Type": "application/json"}
		err := NewPingenError("Test error", "", 400, headers)

		if err.RequestID != "" {
			t.Errorf("Expected empty request ID, got '%s'", err.RequestID)
		}
	})

	t.Run("with nil headers", func(t *testing.T) {
		err := NewPingenError("Test error", "", 400, nil)

		if err.Headers != nil {
			t.Error("Expected headers to be nil")
		}
		if err.RequestID != "" {
			t.Error("Expected request ID to be empty")
		}
	})
}

func TestPingenError_Error(t *testing.T) {
	t.Run("error string format", func(t *testing.T) {
		headers := map[string]string{"X-Request-Id": "req-456"}
		err := NewPingenError("Something went wrong", "", 500, headers)

		expected := "PingenError: Something went wrong (Status Code: 500, Request ID: req-456)"
		if err.Error() != expected {
			t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("error string without request ID", func(t *testing.T) {
		err := NewPingenError("Something went wrong", "", 500, nil)

		expected := "PingenError: Something went wrong (Status Code: 500, Request ID: )"
		if err.Error() != expected {
			t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
		}
	})
}

func TestPingenError_JSONSerialization(t *testing.T) {
	t.Run("JSON marshaling excludes status code and headers", func(t *testing.T) {
		headers := map[string]string{
			"X-Request-Id": "req-789",
			"Content-Type": "application/json",
		}
		jsonBody := `{"error": "test"}`

		err := NewPingenError("Test error", jsonBody, 400, headers)

		data, marshalErr := json.Marshal(err)
		if marshalErr != nil {
			t.Fatalf("Failed to marshal PingenError: %v", marshalErr)
		}

		var result map[string]interface{}
		if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", unmarshalErr)
		}

		if result["message"] != "Test error" {
			t.Errorf("Expected message in JSON, got %v", result["message"])
		}
		if result["json_body"] == nil {
			t.Error("Expected json_body in JSON")
		}

		if _, exists := result["status_code"]; exists {
			t.Error("Expected status_code to be excluded from JSON")
		}
		if _, exists := result["headers"]; exists {
			t.Error("Expected headers to be excluded from JSON")
		}
		if _, exists := result["request_id"]; exists {
			t.Error("Expected request_id to be excluded from JSON")
		}
	})
}

func TestNewAuthenticationError(t *testing.T) {
	t.Run("authentication error creation", func(t *testing.T) {
		message := "Invalid API key"
		statusCode := 401
		headers := map[string]string{"X-Request-Id": "auth-123"}

		authErr := NewAuthenticationError(message, "", statusCode, headers)

		if authErr.Message != message {
			t.Errorf("Expected message %s, got %s", message, authErr.Message)
		}
		if authErr.StatusCode != statusCode {
			t.Errorf("Expected status code %d, got %d", statusCode, authErr.StatusCode)
		}
		if authErr.RequestID != "auth-123" {
			t.Errorf("Expected request ID 'auth-123', got '%s'", authErr.RequestID)
		}
	})

	t.Run("authentication error implements error interface", func(t *testing.T) {
		authErr := NewAuthenticationError("Invalid API key", "", 401, nil)

		var err error = authErr
		expected := "PingenError: Invalid API key (Status Code: 401, Request ID: )"
		if err.Error() != expected {
			t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
		}
	})
}

func TestNewWebhookSignatureException(t *testing.T) {
	t.Run("webhook signature exception creation", func(t *testing.T) {
		message := "Invalid webhook signature"
		err := NewWebhookSignatureException(message)

		if err.Message != message {
			t.Errorf("Expected message %s, got %s", message, err.Message)
		}
	})
}

func TestWebhookSignatureException_Error(t *testing.T) {
	t.Run("error string format", func(t *testing.T) {
		message := "Signature verification failed"
		err := NewWebhookSignatureException(message)

		expected := "WebhookSignatureException: Signature verification failed"
		if err.Error() != expected {
			t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
		}
	})
}
