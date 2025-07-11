package userassociations_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/userassociations"
	"github.com/stretchr/testify/assert"
)

func setupUserAssociations(apiBaseURL string) *userassociations.UserAssociations {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return userassociations.NewUserAssociations(apiRequestor)
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1")
		w.WriteHeader(http.StatusOK)

		responseJSON := `{
            "data": [
                {
                    "id": "userxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
                    "type": "users",
                    "attributes": {
                        "role": "owner",
                        "status": "pending",
                        "created_at": "2020-11-19T09:42:48+0100",
                        "updated_at": "2020-11-19T09:42:48+0100"
                    },
                    "relationships": {
                        "organisation": {
                            "links": {"related": "string"},
                            "data": {
                                "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
                                "type": "organisations"
                            }
                        }
                    },
                    "links": {
                        "self": "string"
                    },
                    "meta": {
                        "abilities": {
                            "self": {"join": "ok", "leave": "ok", "block": "ok"},
                            "organisation": {"manage": "ok"}
                        }
                    }
                }
            ],
            "included": [{}],
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
                "to": 1,
                "total": 1
            }
        }`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	userAssociations := setupUserAssociations(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	response, err := userAssociations.GetCollection(params, headers)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, "userxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", response.Data[0].ID)
	assert.Equal(t, "users", response.Data[0].Type)
	assert.Equal(t, "owner", response.Data[0].Attributes.Role)
}

func TestGetCollection_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	userAssociations := setupUserAssociations(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := userAssociations.GetCollection(params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1)"
	assert.Equal(t, expectedMessage, err.Error())
	assert.Equal(t, http.StatusUnauthorized, err.StatusCode, "unexpected status code")
}
