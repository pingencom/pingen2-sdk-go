package ebillevents_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/config"
	"github.com/pingencom/pingen2-sdk-go/ebillevents"
	"github.com/stretchr/testify/assert"
)

const mockValidJSONResponse = `{
  "data": [
    {
      "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "type": "deliverables_events",
      "attributes": {
        "code": "delivered",
        "name": "Delivered",
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
        "ebill": {
          "links": {
            "related": "string"
          },
          "data": {
            "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
            "type": "ebills"
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

func setupEbillEvents(apiBaseURL string) *ebillevents.EbillEvents {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return ebillevents.NewEbillEvents("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", apiRequestor)
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/deliveries/ebills/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/events", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(mockValidJSONResponse))
	}))
	defer server.Close()

	ebillEvents := setupEbillEvents(server.URL)

	response, err := ebillEvents.GetCollection("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "deliverables_events", response.Data[0].Type)
	assert.Equal(t, "delivered", response.Data[0].Attributes.Code)
	assert.Equal(t, "Delivered", response.Data[0].Attributes.Name)
	assert.Equal(t, "Pingen", response.Data[0].Attributes.Producer)
	assert.Equal(t, "8051 Zürich, CH", response.Data[0].Attributes.Location)
	assert.False(t, response.Data[0].Attributes.HasImage)
	assert.Equal(t, []string{"string"}, response.Data[0].Attributes.Data)
	assert.Equal(t, "ebills", response.Data[0].Relationships.Ebill.Data.Type)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", response.Data[0].Relationships.Ebill.Data.ID)
	assert.Equal(t, 1, response.Meta.CurrentPage)
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

	ebillEvents := setupEbillEvents(server.URL)

	_, err := ebillEvents.GetCollection("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", nil, nil)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}
