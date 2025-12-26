package response

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type TestData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TestListResponse struct {
	BaseListResponse
	Data []TestData `json:"data"`
}

func TestJSONResponseHandler_InterpretResponse(t *testing.T) {
	handler := &JSONResponseHandler{}

	t.Run("successful JSON response", func(t *testing.T) {
		jsonData := `{"id": 123, "name": "test"}`
		resp := createMockResponse(http.StatusOK, jsonData, map[string]string{
			"Content-Type": "application/json",
		})
		defer resp.Body.Close()

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected result to be non-nil")
		}

		resultData, ok := result.(*TestData)
		if !ok {
			t.Fatal("Expected result to be *TestData")
		}
		if resultData.ID != 123 {
			t.Errorf("Expected ID 123, got %d", resultData.ID)
		}
		if resultData.Name != "test" {
			t.Errorf("Expected name 'test', got '%s'", resultData.Name)
		}
	})

	t.Run("successful list response with pagination", func(t *testing.T) {
		jsonData := `{
			"data": [
				{"id": 1, "name": "first"},
				{"id": 2, "name": "second"}
			],
			"links": {
				"first": "https://api.example.com/data?page=1",
				"last": "https://api.example.com/data?page=5",
				"prev": null,
				"next": "https://api.example.com/data?page=2",
				"self": "https://api.example.com/data?page=1"
			},
			"meta": {
				"current_page": 1,
				"last_page": 5,
				"per_page": 10,
				"from": 1,
				"to": 2,
				"total": 42
			},
			"included": []
		}`
		resp := createMockResponse(http.StatusOK, jsonData, nil)
		defer resp.Body.Close()

		var target TestListResponse
		result, err := handler.InterpretResponse(resp, &target)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		resultData, ok := result.(*TestListResponse)
		if !ok {
			t.Fatal("Expected result to be *TestListResponse")
		}
		if len(resultData.Data) != 2 {
			t.Errorf("Expected 2 data items, got %d", len(resultData.Data))
		}
		if resultData.Meta.Total != 42 {
			t.Errorf("Expected total 42, got %d", resultData.Meta.Total)
		}
		if resultData.Links.Next != "https://api.example.com/data?page=2" {
			t.Errorf("Expected next link, got '%s'", resultData.Links.Next)
		}
	})

	t.Run("no content response", func(t *testing.T) {
		resp := createMockResponse(http.StatusNoContent, "", nil)
		defer resp.Body.Close()

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defaultResp, ok := result.(*DefaultResponse)
		if !ok {
			t.Fatal("Expected result to be *DefaultResponse")
		}
		if defaultResp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, defaultResp.StatusCode)
		}
		if defaultResp.Body != "" {
			t.Errorf("Expected empty body, got '%s'", defaultResp.Body)
		}
	})

	t.Run("accepted response", func(t *testing.T) {
		bodyContent := "Request accepted for processing"
		resp := createMockResponse(http.StatusAccepted, bodyContent, nil)
		defer resp.Body.Close()

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defaultResp, ok := result.(*DefaultResponse)
		if !ok {
			t.Fatal("Expected result to be *DefaultResponse")
		}
		if defaultResp.StatusCode != http.StatusAccepted {
			t.Errorf("Expected status code %d, got %d", http.StatusAccepted, defaultResp.StatusCode)
		}
		if defaultResp.Body != bodyContent {
			t.Errorf("Expected body '%s', got '%s'", bodyContent, defaultResp.Body)
		}
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		invalidJSON := `{"id": 123, "name": "test"` // missing closing brace
		resp := createMockResponse(http.StatusOK, invalidJSON, map[string]string{
			"X-Request-Id": "req-invalid-json",
		})
		defer resp.Body.Close()

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err == nil {
			t.Fatal("Expected error for invalid JSON")
		}
		if result != nil {
			t.Error("Expected result to be nil on error")
		}

		if err.Message != "Failed to parse response body" {
			t.Errorf("Expected parse error message, got '%s'", err.Message)
		}
		if err.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, err.StatusCode)
		}
		if err.RequestID != "req-invalid-json" {
			t.Errorf("Expected request ID 'req-invalid-json', got '%s'", err.RequestID)
		}
	})

	t.Run("client error response", func(t *testing.T) {
		errorBody := `{"error": "Invalid request", "message": "Field 'name' is required"}`
		resp := createMockResponse(http.StatusBadRequest, errorBody, map[string]string{
			"X-Request-Id": "req-bad-request",
			"Content-Type": "application/json",
		})
		defer resp.Body.Close()

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err == nil {
			t.Fatal("Expected error for 4xx status code")
		}
		if result != nil {
			t.Error("Expected result to be nil on error")
		}

		if err.Message != "API error" {
			t.Errorf("Expected API error message, got '%s'", err.Message)
		}
		if err.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, err.StatusCode)
		}

		if err.JSONBody == nil {
			t.Error("Expected JSONBody to be parsed")
		}
	})

	t.Run("server error response", func(t *testing.T) {
		errorBody := `{"error": "Internal server error"}`
		resp := createMockResponse(http.StatusInternalServerError, errorBody, map[string]string{
			"X-Request-Id": "req-server-error",
		})
		defer resp.Body.Close()

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err == nil {
			t.Fatal("Expected error for 5xx status code")
		}
		if result != nil {
			t.Error("Expected result to be nil on error")
		}

		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, err.StatusCode)
		}
	})

	t.Run("response body read error", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       &errorReader{},
		}
		resp.Header.Set("X-Request-Id", "req-read-error")

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err == nil {
			t.Fatal("Expected error for body read failure")
		}
		if result != nil {
			t.Error("Expected result to be nil on error")
		}

		if err.Message != "Failed to read response body" {
			t.Errorf("Expected read error message, got '%s'", err.Message)
		}
	})

	t.Run("empty response body with success status", func(t *testing.T) {
		resp := createMockResponse(http.StatusOK, "", nil)
		defer resp.Body.Close()

		var target TestData
		result, err := handler.InterpretResponse(resp, &target)

		if err == nil {
			t.Fatal("Expected error for empty JSON")
		}
		if result != nil {
			t.Error("Expected result to be nil on error")
		}
	})
}

func TestConvertHeaders(t *testing.T) {
	t.Run("convert http headers to map", func(t *testing.T) {
		httpHeaders := make(http.Header)
		httpHeaders.Add("Content-Type", "application/json")
		httpHeaders.Add("X-Request-Id", "req-123")
		httpHeaders.Add("Authorization", "Bearer token")

		result := convertHeaders(httpHeaders)

		expected := map[string]string{
			"Content-Type":  "application/json",
			"X-Request-Id":  "req-123",
			"Authorization": "Bearer token",
		}

		if len(result) != len(expected) {
			t.Errorf("Expected %d headers, got %d", len(expected), len(result))
		}

		for key, expectedValue := range expected {
			if result[key] != expectedValue {
				t.Errorf("Expected header %s: %s, got %s", key, expectedValue, result[key])
			}
		}
	})

	t.Run("convert headers with multiple values", func(t *testing.T) {
		httpHeaders := make(http.Header)
		httpHeaders.Add("Accept", "application/json")
		httpHeaders.Add("Accept", "text/plain") // Second value should be ignored

		result := convertHeaders(httpHeaders)

		if result["Accept"] != "application/json" {
			t.Errorf("Expected first value 'application/json', got '%s'", result["Accept"])
		}
	})

	t.Run("empty headers", func(t *testing.T) {
		httpHeaders := make(http.Header)
		result := convertHeaders(httpHeaders)

		if len(result) != 0 {
			t.Errorf("Expected empty map, got %v", result)
		}
	})

	t.Run("headers with empty values", func(t *testing.T) {
		httpHeaders := make(http.Header)
		httpHeaders["Empty-Header"] = []string{} // Empty slice

		result := convertHeaders(httpHeaders)

		if _, exists := result["Empty-Header"]; exists {
			t.Error("Expected empty header to be excluded")
		}
	})
}

func TestStructsJSONSerialization(t *testing.T) {
	t.Run("Links JSON serialization", func(t *testing.T) {
		links := Links{
			First: "https://api.example.com/data?page=1",
			Last:  "https://api.example.com/data?page=10",
			Prev:  "https://api.example.com/data?page=1",
			Next:  "https://api.example.com/data?page=3",
			Self:  "https://api.example.com/data?page=2",
		}

		data, err := json.Marshal(links)
		if err != nil {
			t.Fatalf("Failed to marshal Links: %v", err)
		}

		var unmarshaled Links
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal Links: %v", err)
		}

		if unmarshaled.Next != links.Next {
			t.Errorf("Expected Next %s, got %s", links.Next, unmarshaled.Next)
		}
	})

	t.Run("Meta JSON serialization", func(t *testing.T) {
		meta := Meta{
			CurrentPage: 2,
			LastPage:    10,
			PerPage:     25,
			From:        26,
			To:          50,
			Total:       237,
		}

		data, err := json.Marshal(meta)
		if err != nil {
			t.Fatalf("Failed to marshal Meta: %v", err)
		}

		var unmarshaled Meta
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal Meta: %v", err)
		}

		if unmarshaled.Total != meta.Total {
			t.Errorf("Expected Total %d, got %d", meta.Total, unmarshaled.Total)
		}
	})
}

func createMockResponse(statusCode int, body string, headers map[string]string) *http.Response {
	resp := &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	for key, value := range headers {
		resp.Header.Set(key, value)
	}

	return resp
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (e *errorReader) Close() error {
	return nil
}
