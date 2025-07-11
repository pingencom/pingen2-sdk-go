package letterevents_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/letterevents"
	"github.com/stretchr/testify/assert"
)

const mockValidJSONResponse = `{
  "data": [
    {
      "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "type": "letters_events",
      "attributes": {
        "code": "undeliverable",
        "name": "Content failed inspection",
        "producer": "Pingen",
        "location": "8051 Zürich, CH",
        "has_image": false,
        "data": [
          "string"
        ],
        "emitted_at": "2020-11-19T09:42:48+0100",
        "created_at": "2020-11-19T09:42:48+0100",
        "updated_at": "2020-11-19T09:42:48+0100"
      },
      "relationships": {
        "letter": {
          "links": {
            "related": "string"
          },
          "data": {
            "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
            "type": "letters"
          }
        }
      },
      "links": {
        "self": "string"
      }
    }
  ],
  "included": [
    {}
  ],
  "links": {
    "first": "string",
    "last": "string",
    "prev": "string",
    "next": "string",
    "self": "string"
  },
  "meta": {
    "current_page": 1,
    "last_page": 1,
    "per_page": 10,
    "from": 1,
    "to": 10,
    "total": 0
  }
}`

func setupLetterEvents(apiBaseURL string) *letterevents.LetterEvents {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return letterevents.NewLetterEvents("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", apiRequestor)
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/events", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(mockValidJSONResponse))
	}))
	defer server.Close()

	letterEvents := setupLetterEvents(server.URL)

	response, err := letterEvents.GetCollection("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, 1, response.Meta.CurrentPage)
	assert.Equal(t, 200, http.StatusOK)
}

func TestGetCollection_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	letterEvents := setupLetterEvents(server.URL)

	_, err := letterEvents.GetCollection("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", nil, nil)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestGetIssueCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/events/issues", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(mockValidJSONResponse))
	}))
	defer server.Close()

	letterEvents := setupLetterEvents(server.URL)

	response, err := letterEvents.GetIssueCollection(nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, 1, response.Meta.CurrentPage)
}

func TestGetUndeliverableCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/events/undeliverable", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(mockValidJSONResponse))
	}))
	defer server.Close()

	letterEvents := setupLetterEvents(server.URL)

	response, err := letterEvents.GetUndeliverableCollection(nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, 1, response.Meta.CurrentPage)
}

func TestGetSentCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/events/sent", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(mockValidJSONResponse))
	}))
	defer server.Close()

	letterEvents := setupLetterEvents(server.URL)

	response, err := letterEvents.GetSentCollection(nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, 1, response.Meta.CurrentPage)
}
