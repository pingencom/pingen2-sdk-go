package fileupload_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/fileupload"
	"github.com/stretchr/testify/assert"
)

func setupFileUpload(apiBaseURL string) *fileupload.FileUpload {
    config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	requestor := api.NewAPIRequestor("dummyToken", config)

    return fileupload.NewFileUpload(requestor)
}


func TestRequestFileUpload_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)

		responseJSON := `{
			"data": {
				"id": "filexx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
				"type": "file_uploads",
				"attributes": {
					"url": "https://s3.example/bucket/filename?signer=url",
					"url_signature": "$2y$10$BLOzVbYTXrh4LZbSYNVf7eEDrc58vvQ9PRVZABqV/9WS1eqIcm3M",
					"expires_at": "2020-11-19T09:42:48+0100"
				},
				"links": { "self": "string" }
			}
		}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	fileUploader := setupFileUpload(server.URL)

	response, err := fileUploader.RequestFileUpload()
	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, "https://s3.example/bucket/filename?signer=url", response.Data.Attributes.URL)
	assert.Equal(t, "$2y$10$BLOzVbYTXrh4LZbSYNVf7eEDrc58vvQ9PRVZABqV/9WS1eqIcm3M", response.Data.Attributes.URLSignature)
	assert.Equal(t, "2020-11-19T09:42:48+0100", response.Data.Attributes.ExpiresAt)
}

func TestRequestFileUpload_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	fileUploader := setupFileUpload(server.URL)

	_, err := fileUploader.RequestFileUpload()
	assert.NotNil(t, err)

	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1)"
	assert.Equal(t, expectedMessage, err.Error())
	assert.Equal(t, http.StatusUnauthorized, err.StatusCode)
}

func TestPutFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()

		if string(body) != "This is a test file." {
			t.Errorf("Request body incorrect, got: %s", string(body))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

    fileUploader := setupFileUpload(server.URL)

	tempFile, err := os.CreateTemp("", "testfile*.pdf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("This is a test file."))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	err = fileUploader.PutFile(tempFile.Name(), server.URL)
	assert.Nil(t, err)
}

func TestPutFile_FileNotFound(t *testing.T) {
    fileUploader := setupFileUpload("http://mockserver")

	err := fileUploader.PutFile("/nonexistent/path/to/file.pdf", "https://s3.example/bucket/filename?signer=url")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to open file")
}
