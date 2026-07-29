package oauth_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pingencom/pingen2-sdk-go/config"
	"github.com/pingencom/pingen2-sdk-go/oauth"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeURL(t *testing.T) {
	config, _ := config.InitSDK("testClientId", "testClientSecret", "production")

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
	config, _ := config.InitSDK("testSetClientId", "testClientSecret", "staging")

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
	_, err := config.InitSDK("", "testSetClientSecret", "staging")

	if err == nil {
		t.Fatal("expected an error but got none")
	}

	expectedError := `missing required credentials (ClientID, ClientSecret)`
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestMissingClientSecret(t *testing.T) {
	_, err := config.InitSDK("testSetClientId", "", "")

	expectedError := `missing required credentials (ClientID, ClientSecret)`
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetToken(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testClientSecret", "")

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
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

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
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`Bad Request`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
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

const mockTokenResponseWithoutAccessToken = `{
	"token_type": "Bearer",
	"expires_in": 43200
}`

// tokenServer is a fake /auth/access-tokens endpoint that counts how often it
// was hit and hands out a different access token on every hit.
type tokenServer struct {
	*httptest.Server
	hits     int
	lastForm url.Values
}

func newTokenServer(t *testing.T, expiresIn int64) *tokenServer {
	t.Helper()

	server := &tokenServer{}
	server.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/auth/access-tokens", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := io.ReadAll(r.Body)
		form, parseErr := url.ParseQuery(string(body))
		assert.Nil(t, parseErr)

		server.hits++
		server.lastForm = form

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(
			w,
			`{"token_type":"Bearer","expires_in":%d,"access_token":"accessToken_%d"}`,
			expiresIn,
			server.hits,
		)
	}))

	return server
}

func TestToken_IsExpired(t *testing.T) {
	var nilToken *oauth.Token
	assert.True(t, nilToken.IsExpired(oauth.DefaultTokenRefreshBufferSeconds))

	emptyToken := &oauth.Token{AccessToken: "", ExpiresIn: 43200, IssuedAt: time.Now()}
	assert.True(t, emptyToken.IsExpired(oauth.DefaultTokenRefreshBufferSeconds))

	freshToken := &oauth.Token{AccessToken: "accessToken", ExpiresIn: 43200, IssuedAt: time.Now()}
	assert.False(t, freshToken.IsExpired(oauth.DefaultTokenRefreshBufferSeconds))
	assert.False(t, freshToken.IsExpired(0))

	insideBufferToken := &oauth.Token{
		AccessToken: "accessToken",
		ExpiresIn:   43200,
		IssuedAt:    time.Now().Add(-43170 * time.Second),
	}
	assert.True(t, insideBufferToken.IsExpired(oauth.DefaultTokenRefreshBufferSeconds))
	assert.False(t, insideBufferToken.IsExpired(0))
}

func TestOAuth_Defaults(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	auth := oauth.NewOAuth(config, "")

	assert.Equal(t, oauth.DefaultTokenRefreshBufferSeconds, auth.TokenRefreshBufferSeconds())
	assert.Nil(t, auth.GetCurrentToken())

}

func TestOAuth_GetAccessToken_CachesToken(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := newTokenServer(t, 43200)
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	auth := oauth.NewOAuth(config, "")

	first, err := auth.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_1", first)
	assert.Equal(t, 1, server.hits)

	second, err := auth.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_1", second)
	assert.Equal(t, 1, server.hits)

	current := auth.GetCurrentToken()
	assert.NotNil(t, current)
	assert.Equal(t, "accessToken_1", current.AccessToken)
	assert.Equal(t, int64(43200), current.ExpiresIn)
	assert.False(t, current.IssuedAt.IsZero())
}

func TestOAuth_GetAccessToken_RefreshesExpiredToken(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := newTokenServer(t, 43200)
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	auth := oauth.NewOAuth(config, "")

	first, err := auth.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_1", first)
	assert.Equal(t, 1, server.hits)

	current := auth.GetCurrentToken()
	current.IssuedAt = current.IssuedAt.Add(-time.Duration(current.ExpiresIn) * time.Second)

	second, err := auth.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_2", second)
	assert.NotEqual(t, first, second)
	assert.Equal(t, 2, server.hits)
}

func TestOAuth_RefreshToken_AlwaysFetches(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := newTokenServer(t, 43200)
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	auth := oauth.NewOAuth(config, "")

	first, err := auth.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_1", first)
	assert.Equal(t, 1, server.hits)

	refreshed, err := auth.RefreshToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_2", refreshed.AccessToken)
	assert.Equal(t, 2, server.hits)

	refreshedAgain, err := auth.RefreshToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_3", refreshedAgain.AccessToken)
	assert.Equal(t, 3, server.hits)

	cached, err := auth.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, "accessToken_3", cached)
	assert.Equal(t, 3, server.hits)
}

func TestOAuth_GetAccessToken_MissingAccessToken(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockTokenResponseWithoutAccessToken)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	auth := oauth.NewOAuth(config, "")

	token, err := auth.GetAccessToken()

	assert.Empty(t, token)
	assert.NotNil(t, err)
	assert.Equal(t, "token response did not contain an access_token", err.Error())
	assert.Nil(t, auth.GetCurrentToken())
}

func TestOAuth_GetAccessToken_RequestFailed(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	auth := oauth.NewOAuth(config, "")

	token, err := auth.GetAccessToken()

	assert.Empty(t, token)
	assert.NotNil(t, err)
	assert.Equal(t, "request failed with status code: 401", err.Error())
	assert.Nil(t, auth.GetCurrentToken())
}

func TestOAuth_GetAccessToken_ForwardsScope(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := newTokenServer(t, 43200)
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	auth := oauth.NewOAuth(config, "letter batch webhook")

	_, err := auth.GetAccessToken()

	assert.Nil(t, err)
	assert.Equal(t, "client_credentials", server.lastForm.Get("grant_type"))
	assert.Equal(t, "letter batch webhook", server.lastForm.Get("scope"))
	assert.Equal(t, "testSetClientId", server.lastForm.Get("client_id"))
	assert.Equal(t, "testSetClientSecret", server.lastForm.Get("client_secret"))
}

func TestOAuth_GetAccessToken_OmitsEmptyScope(t *testing.T) {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")

	server := newTokenServer(t, 43200)
	defer server.Close()
	config.SetAPIBaseURL(server.URL)

	auth := oauth.NewOAuth(config, "")

	_, err := auth.GetAccessToken()

	assert.Nil(t, err)
	assert.Equal(t, "client_credentials", server.lastForm.Get("grant_type"))
	_, hasScope := server.lastForm["scope"]
	assert.False(t, hasScope)
}

func TestRefreshToken_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_scope"}`))
	}))
	defer server.Close()

	config, _ := config.InitSDK("testClientId", "testClientSecret", "")
	config.SetAPIBaseURL(server.URL)

	source := oauth.NewOAuth(config, "letter")

	token, err := source.RefreshToken()

	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Nil(t, source.GetCurrentToken())
}

func TestGetToken_TransportError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	unreachable := server.URL
	server.Close()

	config, _ := config.InitSDK("testClientId", "testClientSecret", "")
	config.SetAPIBaseURL(unreachable)

	resp, err := oauth.GetToken(config, map[string]string{"grant_type": "client_credentials"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to send request")
}

func TestOAuth_ConcurrentGetAccessTokenFetchesOnce(t *testing.T) {
	var hits int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"access_token":"tok","expires_in":3600,"token_type":"Bearer"}`)
	}))
	defer server.Close()

	config, _ := config.InitSDK("testClientId", "testClientSecret", "")
	config.SetAPIBaseURL(server.URL)

	source := oauth.NewOAuth(config, "letter")

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			token, err := source.GetAccessToken()
			assert.NoError(t, err)
			assert.Equal(t, "tok", token)
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(1), atomic.LoadInt64(&hits))
}
