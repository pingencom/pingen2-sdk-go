//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/config"
	"github.com/pingencom/pingen2-sdk-go/oauth"
	"github.com/pingencom/pingen2-sdk-go/organisations"
	"github.com/stretchr/testify/assert"
)

func buildOAuth(t *testing.T) (*oauth.OAuth, *config.Config) {
	t.Helper()

	creds := requireCredentials(t)

	config, err := config.InitSDK(creds.ClientID, creds.ClientSecret, environment(creds))
	if err != nil {
		t.Fatalf("Failed to initialise the SDK: %v", err)
	}

	return oauth.NewOAuth(config, Scope), config
}

func TestOAuthTokenCanBeObtainedAndUsed(t *testing.T) {
	o, config := buildOAuth(t)

	token, err := o.GetAccessToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "Token request must succeed")

	current := o.GetCurrentToken()
	if current == nil {
		t.Fatal("Expected a cached token after a successful token request")
	}
	assert.Greater(t, current.ExpiresIn, int64(0))

	organisationsClient := organisations.NewOrganisations(api.NewAPIRequestorWithTokenSource(o, config))

	collection, pingenErr := organisationsClient.GetCollection(nil, nil)
	assert.Nil(t, pingenErr)
	assert.NotEmpty(t, collection.Data, "Expected at least one organisation")
}

func TestOAuthTokenIsReusedWhileValid(t *testing.T) {
	o, _ := buildOAuth(t)

	first, err := o.GetAccessToken()
	assert.NoError(t, err)

	second, err := o.GetAccessToken()
	assert.NoError(t, err)

	assert.Equal(t, first, second)
}

func TestOAuthExpiredTokenIsRefreshed(t *testing.T) {
	o, config := buildOAuth(t)

	organisationsClient := organisations.NewOrganisations(api.NewAPIRequestorWithTokenSource(o, config))

	_, pingenErr := organisationsClient.GetCollection(nil, nil)
	assert.Nil(t, pingenErr)

	firstToken := o.GetCurrentToken()
	if firstToken == nil {
		t.Fatal("Expected a cached token after a successful request")
	}
	firstValue := firstToken.AccessToken

	// Backdate the token so it is treated as expiring within the refresh
	// buffer, guaranteeing a refresh on the next request.
	firstToken.IssuedAt = firstToken.IssuedAt.Add(-time.Duration(firstToken.ExpiresIn) * time.Second)
	assert.True(t, firstToken.IsExpired(o.TokenRefreshBufferSeconds()))

	_, pingenErr = organisationsClient.GetCollection(nil, nil)
	assert.Nil(t, pingenErr)

	secondToken := o.GetCurrentToken()
	if secondToken == nil {
		t.Fatal("Expected a cached token after the refresh")
	}
	assert.NotEqual(t, firstValue, secondToken.AccessToken, "A new access token must have been fetched")
	assert.False(t, secondToken.IsExpired(o.TokenRefreshBufferSeconds()))
}
