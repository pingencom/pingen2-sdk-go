package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pingencom/pingen2-sdk-go/config"
)

const DefaultTokenRefreshBufferSeconds int64 = 60

type Token struct {
	AccessToken string
	ExpiresIn   int64
	IssuedAt    time.Time
}

func (t *Token) IsExpired(bufferSeconds int64) bool {
	if t == nil || t.AccessToken == "" {
		return true
	}

	expiresAt := t.IssuedAt.Add(time.Duration(t.ExpiresIn-bufferSeconds) * time.Second)
	return !time.Now().Before(expiresAt)
}

type OAuth struct {
	config        *config.Config
	scope         string
	refreshBuffer int64

	mu    sync.Mutex
	token *Token
}

func NewOAuth(cfg *config.Config, scope string) *OAuth {
	return &OAuth{
		config:        cfg,
		scope:         scope,
		refreshBuffer: DefaultTokenRefreshBufferSeconds,
	}
}

func (o *OAuth) TokenRefreshBufferSeconds() int64 {
	return o.refreshBuffer
}

func (o *OAuth) GetAccessToken() (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.token.IsExpired(o.refreshBuffer) {
		return o.token.AccessToken, nil
	}

	token, err := o.fetchToken()
	if err != nil {
		return "", err
	}

	o.token = token
	return token.AccessToken, nil
}

func (o *OAuth) GetCurrentToken() *Token {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.token
}

func (o *OAuth) RefreshToken() (*Token, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	token, err := o.fetchToken()
	if err != nil {
		return nil, err
	}

	o.token = token
	return token, nil
}

// The caller must hold o.mu.
func (o *OAuth) fetchToken() (*Token, error) {
	params := map[string]string{"grant_type": "client_credentials"}
	if o.scope != "" {
		params["scope"] = o.scope
	}

	issuedAt := time.Now()

	resp, err := GetToken(o.config, params)
	if err != nil {
		return nil, err
	}

	accessToken, ok := resp["access_token"].(string)
	if !ok || accessToken == "" {
		return nil, fmt.Errorf("token response did not contain an access_token")
	}

	token := &Token{
		AccessToken: accessToken,
		IssuedAt:    issuedAt,
	}

	if expiresIn, ok := resp["expires_in"].(float64); ok {
		token.ExpiresIn = int64(expiresIn)
	}

	return token, nil
}

func AuthorizeURL(cfg *config.Config, params map[string]string) (string, error) {
	basePath := cfg.GetAuthBaseURL()

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	values.Set("client_id", cfg.GetClientID())
	if values.Get("response_type") == "" {
		values.Set("response_type", "code")
	}

	authURL, _ := url.Parse(basePath)
	authURL.RawQuery = values.Encode()
	return authURL.String(), nil
}

func GetToken(cfg *config.Config, params map[string]string) (map[string]interface{}, error) {
	apiURL := cfg.GetAPIBaseURL()

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	values.Set("client_id", cfg.GetClientID())
	values.Set("client_secret", cfg.GetClientSecret())

	client := &http.Client{Timeout: cfg.GetRequestTimeout()}
	req, _ := http.NewRequest("POST", apiURL+"/auth/access-tokens", strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", cfg.GetUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

func GetTokenFromImplicit(fragment string) (map[string]string, error) {
	pairs := strings.Split(fragment, "&")
	params := make(map[string]string)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid fragment format: %s", fragment)
		}
		params[kv[0]] = kv[1]
	}

	return map[string]string{
		"access_token": params["access_token"],
		"expires_in":   params["expires_in"],
	}, nil
}
