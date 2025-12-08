package emails_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/emails"
	"github.com/stretchr/testify/assert"
)

const mockResponse = `{
	"data": {
		"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11",
		"type": "emails",
		"attributes": {
			"status": "send",
			"file_original_name": "lorem.pdf",
			"file_pages": 2,
			"recipient_identifier": "info@test.com",
			"price_currency": "CHF",
			"price_value": 1.25,
			"source": "api",
			"submitted_at": "2025-11-29T09:42:48+0100",
			"created_at": "2025-11-19T09:42:48+0100",
			"updated_at": "2025-11-19T09:42:48+0100"
		},
		"relationships": {
			"organisation": {
				"links": {
					"related": "string"
				},
				"data": {
					"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11",
					"type": "organisations"
				}
			},
			"events": {
				"links": {
					"related": {
						"href": "string",
						"meta": {
							"count": 0
						}
					}
				}
			}
		},
        "links": {
			"self": "string"
		},
		"meta": {
			"abilities": {
				"self": {
					"delete": "state"
				}
			}
		}
	},
	"included": [{}]
}`

func setupUnauthorizedServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyy12")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
}

func setupEmail(apiBaseURL string) *emails.Emails {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return emails.NewEmails("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", apiRequestor)
}

func TestGetDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11/deliveries/emails/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	emailClient := setupEmail(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := emailClient.GetDetails("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", params, headers)

	assert.Nil(t, err)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", resp.Data.ID)
	assert.Equal(t, "lorem.pdf", resp.Data.Attributes.FileOriginalName)
	assert.Equal(t, "2025-11-29T09:42:48+0100", resp.Data.Attributes.SubmittedAt)
}

func TestGetDetails_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	emailClient := setupEmail(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := emailClient.GetDetails("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyy12)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestGetCollection(t *testing.T) {
	mockCollectionResponse := `{
		"data": [
			{
				"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11",
				"type": "emails",
				"attributes": {
                    "status": "send",
                    "file_original_name": "lorem.pdf",
                    "file_pages": 2,
                    "recipient_identifier": "info@test.com",
                    "price_currency": "CHF",
                    "price_value": 1.25,
                    "source": "api",
                    "submitted_at": "2025-11-29T09:42:48+0100",
                    "created_at": "2025-11-19T09:42:48+0100",
                    "updated_at": "2025-11-19T09:42:48+0100"
                },
                "relationships": {
                    "organisation": {
                        "links": {
                            "related": "string"
                        },
                        "data": {
                            "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11",
                            "type": "organisations"
                        }
                    },
                    "events": {
                        "links": {
                            "related": {
                                "href": "string",
                                "meta": {
                                    "count": 0
                                }
                            }
                        }
                    }
                },
				"links": { "self": "string" }
			}
		],
		"included": [{}],
		"links": { "first": "first-link", "last": "last-link" },
		"meta": { "current_page": 1, "last_page": 1, "per_page": 10, "total": 1 }
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11/deliveries/emails", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockCollectionResponse))
	}))
	defer server.Close()

	emailClient := setupEmail(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := emailClient.GetCollection(params, headers)

	assert.Nil(t, err)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", resp.Data[0].ID)
	assert.Equal(t, "lorem.pdf", resp.Data[0].Attributes.FileOriginalName)
	assert.Equal(t, 2, resp.Data[0].Attributes.FilePages)
	assert.Equal(t, 1, resp.Meta.CurrentPage)
	assert.Equal(t, 10, resp.Meta.PerPage)
	assert.Equal(t, 1, resp.Meta.Total)
}

func TestGetCollection_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()
	emailClient := setupEmail(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := emailClient.GetCollection(params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyy12)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestUploadAndCreate(t *testing.T) {
	metaData := map[string]interface{}{
		"sender_name":     "Test Example",
		"recipient_email": "info@test.com",
		"recipient_name":  "R_Example",
		"reply_email":     "info_reply@test.com",
		"reply_name":      "Reply Example",
		"subject":         "Your new invoice Number xyz",
		"content":         "Dear recipient\\n\\nAttached is your invoice",
	}

	relationships := map[string]interface{}{
		"preset": map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11",
				"type": "presets",
			},
		},
	}

	var server *httptest.Server

	counter := 0
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter == 0 {
			assert.Equal(t, "/file-upload", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			mockUploadResponse := fmt.Sprintf(`{
                "data": {
                    "attributes": {
                        "url": "%s/upload",
                        "url_signature": "mock-signature"
                    }
                }
            }`, server.URL)

			_, _ = w.Write([]byte(mockUploadResponse))
		}

		if counter == 1 {
			assert.Equal(t, "/upload", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)
			w.WriteHeader(http.StatusNoContent)
		}

		if counter == 2 {
			assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11/deliveries/emails", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			body, _ := io.ReadAll(r.Body)
			assert.Contains(t, string(body), `"file_original_name":"uploaded-test.pdf"`)

			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(mockResponse))
		}

		counter++
	}))
	defer server.Close()

	emailClient := setupEmail(server.URL)

	filePath := "testFile.pdf"

	resp, err := emailClient.UploadAndCreate(
		filePath,
		"uploaded-test.pdf",
		true,
		metaData,
		relationships,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", resp.Data.ID)
}

func TestUploadAndCreate_Error(t *testing.T) {
	metaData := map[string]interface{}{
		"sender_name":     "Test Example",
		"recipient_email": "info@test.com",
		"recipient_name":  "R_Example",
		"reply_email":     "info_reply@test.com",
		"reply_name":      "Reply Example",
		"subject":         "Your new invoice Number xyz",
		"content":         "Dear recipient\\n\\nAttached is your invoice",
	}

	server := setupUnauthorizedServer()
	defer server.Close()

	emailClient := setupEmail(server.URL)
	filePath := "testFile.pdf"

	_, err := emailClient.UploadAndCreate(
		filePath,
		"uploaded-test.pdf",
		true,
		metaData,
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyy12)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestUploadAndCreate_ErrorInPut(t *testing.T) {
	var server *httptest.Server
	metaData := map[string]interface{}{
		"sender_name":     "Test Example",
		"recipient_email": "info@test.com",
		"recipient_name":  "R_Example",
		"reply_email":     "info_reply@test.com",
		"reply_name":      "Reply Example",
		"subject":         "Your new invoice Number xyz",
		"content":         "Dear recipient\\n\\nAttached is your invoice",
	}

	counter := 0
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter == 0 {
			assert.Equal(t, "/file-upload", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			mockUploadResponse := fmt.Sprintf(`{
                "data": {
                    "attributes": {
                        "url": "%s/upload",
                        "url_signature": "mock-signature"
                    }
                }
            }`, server.URL)

			_, _ = w.Write([]byte(mockUploadResponse))
		}

		if counter == 1 {
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyy12")
			w.WriteHeader(http.StatusUnauthorized)

			responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
			_, _ = w.Write([]byte(responseJSON))
		}
		counter++
	}))
	defer server.Close()

	emailClient := setupEmail(server.URL)
	filePath := "testFile.pdf"

	_, err := emailClient.UploadAndCreate(
		filePath,
		"uploaded-test.pdf",
		true,
		metaData,
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: Api error (Status Code: 401, Request ID: )"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11/deliveries/emails", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), `"file_original_name":"uploaded-test.pdf"`)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	metaData := map[string]interface{}{
		"sender_name":     "Test Example",
		"recipient_email": "info@test.com",
		"recipient_name":  "R_Example",
		"reply_email":     "info_reply@test.com",
		"reply_name":      "Reply Example",
		"subject":         "Your new invoice Number xyz",
		"content":         "Dear recipient\\n\\nAttached is your invoice",
	}
	emailClient := setupEmail(server.URL)

	filePath := "testFile.pdf"

	resp, err := emailClient.Create(
		filePath,
		"test-signature",
		"uploaded-test.pdf",
		true,
		metaData,
		nil,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", resp.Data.ID)
}

func TestCreate_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	emailClient := setupEmail(server.URL)

	_, err := emailClient.Create(
		"https://s3.example.com/file/test",
		"signature-123",
		"test.pdf",
		true,
		nil,
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyy12)"
	assert.Equal(t, expectedMessage, err.Error())
}
