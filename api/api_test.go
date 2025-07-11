package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestPreparePath(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	requestor := NewAPIRequestor("dummyToken", config)
	urlPath := "/documents"
	params := map[string]string{"key": "value", "anotherKey": "anotherValue"}

	expected := "https://api.pingen.com/documents?anotherKey=anotherValue&key=value"
	actual := requestor.preparePath(urlPath, params)

	assert.Equal(t, expected, actual)
}

func TestRequestHeaders(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	requestor := NewAPIRequestor("dummyToken", config)

	extraHeaders := map[string]string{
		"Custom-Header": "CustomValue",
	}

	headers := requestor.requestHeaders(extraHeaders)

	assert.Equal(t, "PINGEN.SDK.GO", headers.Get("User-Agent"))
	assert.Equal(t, "Bearer dummyToken", headers.Get("Authorization"))
	assert.Equal(t, "application/vnd.api+json", headers.Get("Content-Type"))
	assert.Equal(t, "application/vnd.api+json", headers.Get("Accept"))
	assert.Equal(t, "CustomValue", headers.Get("Custom-Header"))
}

func TestPerformGetRequest_Success(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer dummyToken", r.Header.Get("Authorization"))
		assert.Equal(t, "/api/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"result":"success"}`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	params := map[string]string{}
	headers := map[string]string{}

	var result map[string]interface{}
	resp, err := requestor.PerformGetRequest("/api/test", &result, params, headers)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", result["result"])
}

func TestPerformPostRequest_Success(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/vnd.api+json", r.Header.Get("Content-Type"))
		assert.Equal(t, "/api/test", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"key":"value"}`, string(body))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"result":"created"}`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	payload := []byte(`{"key":"value"}`)
	headers := map[string]string{}

	var result map[string]interface{}
	resp, err := requestor.PerformPostRequest("/api/test", &result, payload, headers)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "created", result["result"])
}

func TestPerformPutRequest_Success(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/test", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "test content", string(body))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	tempFile := strings.NewReader("test content")
	err := requestor.PerformPutRequest(server.URL+"/api/test", tempFile)

	assert.Nil(t, err)
}

func TestPerformPutRequest_NetworkError(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL("http://invalid-url")
	requestor := NewAPIRequestor("dummyToken", config)

	tempFile := strings.NewReader("test content")
	err := requestor.PerformPutRequest("/api/test", tempFile)

	assert.NotNil(t, err)
	assert.Equal(t, "Internal error", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
}

func TestPerformPutRequest_ApiError(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/test", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "test content", string(body))

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	tempFile := strings.NewReader("test content")
	err := requestor.PerformPutRequest(server.URL+"/api/test", tempFile)

	assert.NotNil(t, err)
	assert.Equal(t, "Api error", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestPerformPatchRequest_Success(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/test", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"key":"value"}`, string(body))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"result":"updated"}`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	payload := []byte(`{"key":"value"}`)
	headers := map[string]string{}

	var result map[string]interface{}
	resp, err := requestor.PerformPatchRequest("/api/test", &result, payload, headers)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "updated", result["result"])
}

func TestPerformCancelRequest_Success(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/test", r.URL.Path)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	resp, err := requestor.PerformCancelRequest("/api/test")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestPerformDeleteRequest_Success(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/test", r.URL.Path)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	resp, err := requestor.PerformDeleteRequest("/api/test")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestPerformGetRequest_Error(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error":"internal server error"}`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	params := map[string]string{}
	headers := map[string]string{}

	var result map[string]interface{}
	_, err := requestor.PerformGetRequest("/api/fail", &result, params, headers)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 500, Request ID: requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestPerformStreamRequest(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mock-url", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "mock file content from redirect")
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	stream, err := requestor.PerformStreamRequest("/mock-url")

	assert.Nil(t, err)
	assert.NotNil(t, stream)

	defer stream.Close()
	responseData, readErr := io.ReadAll(stream)
	assert.Nil(t, readErr)
	assert.Equal(t, "mock file content from redirect", string(responseData))
}

func TestPerformStreamRequest_WrongStatus(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mock-url", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusMultipleChoices)
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	requestor := NewAPIRequestor("dummyToken", config)

	stream, err := requestor.PerformStreamRequest("/mock-url")

	assert.NotNil(t, err)
	assert.Nil(t, stream)
	expectedMessage := "PingenError: Invalid HTTP response (Status Code: 300, Request ID: )"
	assert.Equal(t, expectedMessage, err.Error())
}
