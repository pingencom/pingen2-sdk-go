package webhooks_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/webhooks"
	"github.com/stretchr/testify/assert"
)

const mockWebhookResponse = `{
	"data": {
		"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
		"type": "webhooks",
		"attributes": {
			"event_category": "issues",
			"url": "https://example.com",
			"signing_key": "abcdef123456"
		},
		"relationships": {
			"organisation": {
				"links": {
					"related": "https://api.example.com/organisation/1"
				},
				"data": {
					"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					"type": "organisations"
				}
			}
		},
		"links": {
			"self": "https://api.example.com/webhooks/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"
		}
	},
	"included": []
}`

func setupUnauthorizedServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
}

func setupWebhooks(apiBaseURL string) *webhooks.Webhooks {
    config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return webhooks.NewWebhooks("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", apiRequestor)
}

func TestGetWebhookDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/webhooks/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockWebhookResponse))
	}))
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := webhookClient.GetDetails("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", params, headers)

	assert.Nil(t, err)

	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", resp.Data.ID)
	assert.Equal(t, "issues", resp.Data.Attributes.EventCategory)
	assert.Equal(t, "https://example.com", resp.Data.Attributes.URL)
	assert.Equal(t, "abcdef123456", resp.Data.Attributes.SigningKey)
}

func TestGetWebhookDetails_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := webhookClient.GetDetails("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", params, headers)

	assert.NotNil(t, err)
}

func TestGetWebhookCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/webhooks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					"type": "webhooks",
					"attributes": {
						"event_category": "issues",
						"url": "https://valid-url",
						"signing_key": "d09a095a0d1d2ae896f985c0fff1ad51"
					},
					"relationships": {
						"organisation": {
							"links": {
								"related": "string"
							},
							"data": {
								"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
								"type": "organisations"
							}
						}
					},
					"links": {
						"self": "string"
					}
				}
			],
			"included": [],
			"links": {
				"first": "string",
				"last": "string",
				"prev": "string",
				"next": "string",
				"self": "string"
			},
			"meta": {
				"current_page": 0,
				"last_page": 0,
				"per_page": 0,
				"from": 0,
				"to": 0,
				"total": 0
			}
		}`))
	}))
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := webhookClient.GetCollection(params, headers)

	assert.Nil(t, err)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", resp.Data[0].ID)
	assert.Equal(t, "https://valid-url", resp.Data[0].Attributes.URL)
	assert.Equal(t, "d09a095a0d1d2ae896f985c0fff1ad51", resp.Data[0].Attributes.SigningKey)
	assert.Equal(t, "issues", resp.Data[0].Attributes.EventCategory)
}

func TestGetWebhookCollection_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := webhookClient.GetCollection(params, headers)

	assert.NotNil(t, err)
}

func TestCreateWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/webhooks", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), `"event_category":"issues"`)
		assert.Contains(t, string(body), `"url":"https://valid-url"`)
		assert.Contains(t, string(body), `"signing_key":"d09a095a0d1d2ae896f985c0fff1ad51"`)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx2")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"data": {
				"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11",
				"type": "webhooks",
				"attributes": {
					"event_category": "issues",
					"url": "https://valid-url",
					"signing_key": "d09a095a0d1d2ae896f985c0fff1ad51"
				},
				"relationships": {
					"organisation": {
						"links": {
							"related": "string"
						},
						"data": {
							"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
							"type": "organisations"
						}
					}
				}
			}
		}`))
	}))
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	resp, err := webhookClient.Create("issues", "https://valid-url", "d09a095a0d1d2ae896f985c0fff1ad51")

	assert.Nil(t, err)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx11", resp.Data.ID)
	assert.Equal(t, "issues", resp.Data.Attributes.EventCategory)
	assert.Equal(t, "https://valid-url", resp.Data.Attributes.URL)
	assert.Equal(t, "d09a095a0d1d2ae896f985c0fff1ad51", resp.Data.Attributes.SigningKey)
}

func TestCreateWebhook_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	_, err := webhookClient.Create("issues", "https://valid-url", "d09a095a0d1d2ae896f985c0fff1ad51")

	assert.NotNil(t, err)
}

func TestDeleteWebhook(t *testing.T) {
	webhookID := "testdelx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedURL := "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/webhooks/" + webhookID
		assert.Equal(t, expectedURL, r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	resp, err := webhookClient.Delete(webhookID)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteWebhook_Error(t *testing.T) {
	webhookID := "testdelx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"
	server := setupUnauthorizedServer()
	defer server.Close()

	webhookClient := setupWebhooks(server.URL)

	_, err := webhookClient.Delete(webhookID)
	assert.NotNil(t, err)
}
