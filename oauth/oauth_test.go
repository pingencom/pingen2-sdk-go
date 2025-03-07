package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/oauth"
)

func TestAuthorizeURL(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testClientId", "testClientSecret", "production")

	authURL, err := oauth.AuthorizeURL(config, map[string]string{
		"scope":         "letter",
		"state":         "RANDOMGENERATEDSTRING",
		"response_type": "code",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsedURL, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	query := parsedURL.Query()
	if parsedURL.Scheme != "https" {
		t.Fatalf("expected scheme https, got: %s", parsedURL.Scheme)
	}
	if parsedURL.Host != "identity.pingen.com" {
		t.Fatalf("expected host identity.pingen.com, got: %s", parsedURL.Host)
	}
	if query.Get("client_id") != "testClientId" {
		t.Errorf("expected client_id testClientId, got: %s", query.Get("client_id"))
	}
	if query.Get("scope") != "letter" {
		t.Errorf("expected scope letter, got: %s", query.Get("scope"))
	}
	if query.Get("state") != "RANDOMGENERATEDSTRING" {
		t.Errorf("expected state RANDOMGENERATEDSTRING, got: %s", query.Get("state"))
	}
}

func TestAuthorizeURL_Staging(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testClientSecret", "staging")

	authURL, err := oauth.AuthorizeURL(config, map[string]string{
		"scope": "letter",
		"state": "RANDOMGENERATEDSTRING",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsedURL, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	query := parsedURL.Query()
	if parsedURL.Scheme != "https" {
		t.Fatalf("expected scheme https, got: %s", parsedURL.Scheme)
	}
	if parsedURL.Host != "identity-staging.pingen.com" {
		t.Fatalf("expected host identity-staging.pingen.com, got: %s", parsedURL.Host)
	}
	if query.Get("client_id") != "testSetClientId" {
		t.Errorf("expected client_id testSetClientId, got: %s", query.Get("client_id"))
	}
	if query.Get("scope") != "letter" {
		t.Errorf("expected scope letter, got: %s", query.Get("scope"))
	}
	if query.Get("state") != "RANDOMGENERATEDSTRING" {
		t.Errorf("expected state RANDOMGENERATEDSTRING, got: %s", query.Get("state"))
	}
}

func TestMissingClientId(t *testing.T) {
	_, err := pingen2sdk.InitSDK("", "testSetClientSecret", "staging")

	if err == nil {
		t.Fatal("expected an error but got none")
	}

	expectedError := `Missing required credentials (ClientID, ClientSecret)`
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestMissingClientSecret(t *testing.T) {
	_, err := pingen2sdk.InitSDK("testSetClientId", "", "")

	expectedError := `Missing required credentials (ClientID, ClientSecret)`
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetToken(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testClientSecret", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"token_type": "Bearer",
			"expires_in": 43200,
			"access_token": "YOUR_ACCESS_TOKEN"
		}`))
	}))
	defer server.Close()

	config.SetAPIBaseURL(server.URL)
	resp, err := oauth.GetToken(config, map[string]string{
		"grant_type":    "client_credentials",
		"client_secret": "testClientSecret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp["access_token"] != "YOUR_ACCESS_TOKEN" {
		t.Errorf("expected access_token YOUR_ACCESS_TOKEN, got: %v", resp["access_token"])
	}
}

func TestGetToken_InvalidStatus(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	_, err := oauth.GetToken(config, map[string]string{
		"grant_type": "client_credentials",
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	expectedError := `request failed with status code: 400`
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetToken_InvalidJson(t *testing.T) {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Bad Request`))
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	_, err := oauth.GetToken(config, map[string]string{
		"grant_type": "client_credentials",
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	expectedError := `failed to decode response: invalid character 'B' looking for beginning of value`
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetTokenFromImplicit(t *testing.T) {
	resp, err := oauth.GetTokenFromImplicit("access_token=mock_access_token&expires_in=43200")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp["access_token"] != "mock_access_token" {
		t.Errorf("expected access_token mock_access_token, got: %s", resp["access_token"])
	}
	if resp["expires_in"] != "43200" {
		t.Errorf("expected expires_in 43200, got: %s", resp["expires_in"])
	}
}

func TestGetTokenFromImplicit_Invalid(t *testing.T) {
	_, err := oauth.GetTokenFromImplicit("invalid")
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	expectedError := `invalid fragment format: invalid`
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}
}
