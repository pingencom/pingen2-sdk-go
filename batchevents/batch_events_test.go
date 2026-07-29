package batchevents_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/batchevents"
	"github.com/pingencom/pingen2-sdk-go/config"
	"github.com/stretchr/testify/assert"
)

const mockValidJSONResponse = `{
  "data": [
    {
      "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "type": "batches_events",
      "attributes": {
        "code": "created",
        "name": "Content failed inspection",
        "producer": "Pingen",
        "location": "8051 Zürich, CH",
        "data": [
          "string"
        ],
        "emitted_at": "2020-11-19T09:42:48+0100",
        "created_at": "2020-11-19T09:42:48+0100",
        "updated_at": "2020-11-19T09:42:48+0100"
      },
      "relationships": {
        "batch": {
          "links": {
            "related": "string"
          },
          "data": {
            "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
            "type": "batches"
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

func setupBatchEvents(apiBaseURL string) *batchevents.BatchEvents {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return batchevents.NewBatchEvents("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", apiRequestor)
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/events", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(mockValidJSONResponse))
	}))
	defer server.Close()

	batchEvents := setupBatchEvents(server.URL)

	response, err := batchEvents.GetCollection("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, 1, response.Meta.CurrentPage)
	assert.Equal(t, 200, http.StatusOK)
}

func TestGetCollection_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"errors":[{"code":"500","title":"Server error"}]}`))
	}))
	defer server.Close()

	config, _ := config.InitSDK("testClientId", "testClientSecret", "")
	config.SetAPIBaseURL(server.URL)
	requestor := api.NewAPIRequestor("dummyToken", config)

	batchEvents := batchevents.NewBatchEvents("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", requestor)

	resp, err := batchEvents.GetCollection("test-batch-id", nil, nil)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	assert.Empty(t, resp.Data)
}
