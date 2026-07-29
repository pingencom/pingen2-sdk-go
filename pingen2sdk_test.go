package pingen2sdk_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/stretchr/testify/assert"
)

const mockTokenResponse = `{
	"access_token": "mockAccessToken",
	"token_type": "Bearer",
	"expires_in": 3600
}`

const mockLetterResponse = `{
	"data": {
		"id": "test-letter-id",
		"type": "letters",
		"attributes": {
			"status": "sent",
			"file_original_name": "test.pdf",
			"file_pages": 2,
			"address": "Test Address",
			"address_position": "left",
			"country": "CH",
			"delivery_product": "fast",
			"print_mode": "simplex",
			"print_spectrum": "color",
			"price_currency": "CHF",
			"price_value": 1.5,
			"paper_types": ["normal"],
			"source": "api",
			"tracking_number": "tracking-number",
			"submitted_at": "2021-11-19T09:42:48+0100",
			"created_at": "2020-11-19T09:42:48+0100",
			"updated_at": "2020-11-19T09:42:48+0100"
		}
	}
}`

func TestNewWithClientCredentials(t *testing.T) {
	c, err := pingen2sdk.New(pingen2sdk.Options{
		ClientID:       "testClientId",
		ClientSecret:   "testClientSecret",
		Scope:          "letter batch",
		OrganisationID: "test-organisation-id",
	})

	assert.NoError(t, err)
	assert.NotNil(t, c.Config)
	assert.NotNil(t, c.OAuth)
	assert.NotNil(t, c.Requestor)
	assert.Equal(t, "test-organisation-id", c.OrganisationID())

	assert.NotNil(t, c.Organisations)
	assert.NotNil(t, c.Users)
	assert.NotNil(t, c.UserAssociations)
	assert.NotNil(t, c.Letters)
	assert.NotNil(t, c.LetterEvents)
	assert.NotNil(t, c.Batches)
	assert.NotNil(t, c.BatchEvents)
	assert.NotNil(t, c.Webhooks)
	assert.NotNil(t, c.Emails)
	assert.NotNil(t, c.EmailEvents)
	assert.NotNil(t, c.Ebills)
	assert.NotNil(t, c.EbillEvents)

	assert.Equal(t, "https://api.pingen.com", c.Config.GetAPIBaseURL())
	assert.Equal(t, "https://identity.pingen.com", c.Config.GetAuthBaseURL())
}

func TestNewStagingEnvironment(t *testing.T) {
	c, err := pingen2sdk.New(pingen2sdk.Options{
		ClientID:       "testClientId",
		ClientSecret:   "testClientSecret",
		Environment:    "staging",
		OrganisationID: "test-organisation-id",
	})

	assert.NoError(t, err)
	assert.Equal(t, "https://api-staging.pingen.com", c.Config.GetAPIBaseURL())
	assert.Equal(t, "https://identity-staging.pingen.com", c.Config.GetAuthBaseURL())
}

func TestNewWithStaticAccessToken(t *testing.T) {
	c, err := pingen2sdk.New(pingen2sdk.Options{
		AccessToken:    "staticAccessToken",
		Environment:    "staging",
		OrganisationID: "test-organisation-id",
	})

	assert.NoError(t, err)
	assert.Nil(t, c.OAuth)
	assert.NotNil(t, c.Requestor)
	assert.NotNil(t, c.Letters)
	assert.NotNil(t, c.EbillEvents)
	assert.Equal(t, "https://api-staging.pingen.com", c.Config.GetAPIBaseURL())
}

func TestNewWithoutCredentialsOrToken(t *testing.T) {
	c, err := pingen2sdk.New(pingen2sdk.Options{
		OrganisationID: "test-organisation-id",
	})

	assert.Error(t, err)
	assert.Nil(t, c)
	assert.Contains(t, err.Error(), "missing required credentials")
}

func TestNewWithoutOrganisationID(t *testing.T) {
	c, err := pingen2sdk.New(pingen2sdk.Options{
		AccessToken: "staticAccessToken",
	})

	assert.NoError(t, err)
	assert.Equal(t, "", c.OrganisationID())
	assert.NotNil(t, c.Letters)
	assert.NotNil(t, c.Batches)
	assert.NotNil(t, c.Emails)
}

func TestForOrganisation(t *testing.T) {
	c, err := pingen2sdk.New(pingen2sdk.Options{
		AccessToken:    "staticAccessToken",
		OrganisationID: "first-organisation-id",
	})
	assert.NoError(t, err)

	scoped := c.ForOrganisation("second-organisation-id")

	assert.Equal(t, "second-organisation-id", scoped.OrganisationID())
	assert.Equal(t, "first-organisation-id", c.OrganisationID())

	assert.NotSame(t, c, scoped)
	assert.NotSame(t, c.Letters, scoped.Letters)
	assert.NotSame(t, c.Batches, scoped.Batches)
	assert.NotSame(t, c.EmailEvents, scoped.EmailEvents)

	assert.Same(t, c.Requestor, scoped.Requestor)
	assert.Same(t, c.Config, scoped.Config)
}

func TestForOrganisationRoutesRequests(t *testing.T) {
	var requestedPath, authorizationHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		authorizationHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockLetterResponse))
	}))
	defer server.Close()

	c, err := pingen2sdk.New(pingen2sdk.Options{
		AccessToken:    "staticAccessToken",
		OrganisationID: "first-organisation-id",
	})
	assert.NoError(t, err)
	c.Config.SetAPIBaseURL(server.URL)

	scoped := c.ForOrganisation("second-organisation-id")

	_, pingenErr := scoped.Letters.GetDetails("test-letter-id", nil, nil)

	assert.Nil(t, pingenErr)
	assert.Equal(t, "/organisations/second-organisation-id/deliveries/letters/test-letter-id", requestedPath)
	assert.Equal(t, "Bearer staticAccessToken", authorizationHeader)
}

func TestClientSmokeThroughOAuth(t *testing.T) {
	var authorizationHeader, requestedScope string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/access-tokens" {
			_ = r.ParseForm()
			requestedScope = r.Form.Get("scope")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockTokenResponse))
			return
		}

		authorizationHeader = r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockLetterResponse))
	}))
	defer server.Close()

	c, err := pingen2sdk.New(pingen2sdk.Options{
		ClientID:       "testClientId",
		ClientSecret:   "testClientSecret",
		Scope:          "letter batch",
		OrganisationID: "test-organisation-id",
	})
	assert.NoError(t, err)
	c.Config.SetAPIBaseURL(server.URL)

	response, pingenErr := c.Letters.GetDetails("test-letter-id", nil, nil)

	assert.Nil(t, pingenErr)
	assert.Equal(t, "test-letter-id", response.Data.ID)
	assert.Equal(t, "letters", response.Data.Type)
	assert.Equal(t, "sent", response.Data.Attributes.Status)
	assert.Equal(t, "Bearer mockAccessToken", authorizationHeader)
	assert.NotNil(t, c.OAuth.GetCurrentToken())
	assert.Equal(t, "mockAccessToken", c.OAuth.GetCurrentToken().AccessToken)
	assert.Equal(t, "letter batch", requestedScope)
}

func TestInitSDKDelegatesToConfig(t *testing.T) {
	cfg, err := pingen2sdk.InitSDK("testClientId", "testClientSecret", "staging")

	assert.NoError(t, err)
	assert.Equal(t, "https://api-staging.pingen.com", cfg.GetAPIBaseURL())
	assert.Equal(t, "testClientId", cfg.GetClientID())

	_, err = pingen2sdk.InitSDK("", "", "")
	assert.Error(t, err)
}

func TestInitSDKWithoutCredentialsDelegatesToConfig(t *testing.T) {
	cfg := pingen2sdk.InitSDKWithoutCredentials("staging")

	assert.Equal(t, "https://api-staging.pingen.com", cfg.GetAPIBaseURL())
	assert.Empty(t, cfg.GetClientID())
}
