package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pingencom/pingen2-sdk-go"
)

type OAuth struct{}

func AuthorizeURL(config *pingen2sdk.Config, params map[string]string) (string, error) {
	basePath := config.GetAuthBaseURL()

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	values.Set("client_id", config.GetClientID())
	if values.Get("response_type") == "" {
		values.Set("response_type", "code")
	}

	authURL, _ := url.Parse(basePath)
	authURL.RawQuery = values.Encode()
	return authURL.String(), nil
}

func GetToken(config *pingen2sdk.Config, params map[string]string) (map[string]interface{}, error) {
	apiURL := config.GetAPIBaseURL()

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	values.Set("client_id", config.GetClientID())
	values.Set("client_secret", config.GetClientSecret())

	client := &http.Client{Timeout: config.GetRequestTimeout()}
	req, _ := http.NewRequest("POST", apiURL+"/auth/access-tokens", strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", config.GetUserAgent())

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
