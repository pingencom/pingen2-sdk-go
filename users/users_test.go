package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/users"
	"github.com/stretchr/testify/assert"
)

func setupUsers(apiBaseURL string) *users.Users {
    config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return users.NewUsers(apiRequestor)
}

func TestGetDetails(t *testing.T) {
	userID := "userxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1")
		w.WriteHeader(http.StatusOK)

		responseJSON := `{
			"data": {
				"id": "userxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
				"type": "users",
				"attributes": {
					"email": "email",
					"first_name": "John",
					"last_name": "Snow",
					"status": "active",
					"language": "en-GB",
					"edition": "string",
					"created_at": "2020-11-19T09:42:48+0100",
					"updated_at": "2020-11-19T09:42:48+0100"
				},
				"relationships": {
					"associations": {
						"links": {
							"related": { "href": "string", "meta": { "count": 0 } }
						}
					},
					"notifications": {
						"links": {
							"related": { "href": "string", "meta": { "count": 0 } }
						}
					}
				},
				"links": { "self": "string" },
				"meta": {
					"abilities": {
						"self": {
							"reach": "ok",
							"act": "ok",
							"resend-activation": "ok"
						}
					}
				}
			},
			"included": [{}]
		}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	user := setupUsers(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	response, err := user.GetDetails(params, headers)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, userID, response.Data.ID)
	assert.Equal(t, "users", response.Data.Type)
	assert.Equal(t, "email", response.Data.Attributes.Email)
}

func TestGetDetails_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	user := setupUsers(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := user.GetDetails(params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-xxxx-xxxx-xxxx-xxxxxxxxxxx1)"
	assert.Equal(t, expectedMessage, err.Error())
	assert.Equal(t, http.StatusUnauthorized, err.StatusCode, "unexpected status code")
}
