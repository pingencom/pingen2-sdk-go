package organisations_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/organisations"
	"github.com/stretchr/testify/assert"
)

func setupUnauthorizedServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
}

func setupOrganisations(apiBaseURL string) *organisations.Organisations {
    config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return organisations.NewOrganisations(apiRequestor)
}

func TestGetDetails(t *testing.T) {
	organisationID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy1")
		w.WriteHeader(http.StatusOK)

		responseJSON := `{
			"data": {
				"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				"type": "organisations",
				"attributes": {
					"name": "ACME GmbH",
					"status": "active",
					"plan": "free",
					"billing_mode": "prepaid",
					"billing_currency": "CHF",
					"billing_balance": 11.23,
					"missing_credits": 0,
					"edition": "string",
					"default_country": "CH",
					"default_address_position": "left",
					"data_retention_addresses": 18,
					"data_retention_pdf": 12,
					"limits_monthly_letters_count": 5000,
					"color": "#0758FF",
					"flags": [
						"string"
					],
					"created_at": "2020-11-19T09:42:48+0100",
					"updated_at": "2020-11-19T09:42:48+0100"
				},
				"relationships": {
					"associations": {
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
							"manage": "ok"
						}
					}
				}
			},
			"included": [{}]
		}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	orgClient := setupOrganisations(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	response, err := orgClient.GetDetails(organisationID, params, headers)

	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, organisationID, response.Data.ID)
	assert.Equal(t, "organisations", response.Data.Type)
	assert.Equal(t, "ACME GmbH", response.Data.Attributes.Name)
	assert.Equal(t, 11.23, response.Data.Attributes.BillingBalance)
	assert.Equal(t, "CH", response.Data.Attributes.DefaultCountry)
	assert.Equal(t, "left", response.Data.Attributes.DefaultAddressPosition)
}

func TestGetDetails_Unauthorized(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	orgClient := setupOrganisations(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := orgClient.GetDetails("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
	assert.Equal(t, http.StatusUnauthorized, err.StatusCode, "unexpected status code")
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy3")
		w.WriteHeader(http.StatusOK)

		responseJSON := `{
			"data": [
				{
					"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					"type": "organisations",
					"attributes": {
						"name": "ACME GmbH",
						"status": "active",
						"plan": "free",
						"billing_mode": "prepaid",
						"billing_currency": "CHF",
						"billing_balance": 11.23,
						"default_country": "CH",
						"created_at": "2020-11-19T09:42:48+0100",
						"updated_at": "2020-11-19T09:42:48+0100"
					},
					"links": {
						"self": "string"
					}
				}
			],
			"links": {
				"self": "string",
				"first": "string",
				"last": "string"
			},
			"meta": {
				"current_page": 1,
				"last_page": 1,
				"total": 1
			}
		}`
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	orgClient := setupOrganisations(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	response, err := orgClient.GetCollection(params, headers)

	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, "ACME GmbH", response.Data[0].Attributes.Name)
	assert.Equal(t, 11.23, response.Data[0].Attributes.BillingBalance)
	assert.Equal(t, "CH", response.Data[0].Attributes.DefaultCountry)
	assert.Equal(t, "string", response.Links.Self)
	assert.Equal(t, 1, response.Meta.CurrentPage)
	assert.Equal(t, 1, response.Meta.Total)
}

func TestGetCollection_Unauthorized(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	orgClient := setupOrganisations(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := orgClient.GetCollection(params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
	assert.Equal(t, http.StatusUnauthorized, err.StatusCode, "unexpected status code")
}
