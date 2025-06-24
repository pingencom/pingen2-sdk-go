package pingen2sdk

import (
	"testing"
	"time"
)

func TestInitSDK(t *testing.T) {
	t.Run("successful initialization with production environment", func(t *testing.T) {
		clientID := "test-client-id"
		clientSecret := "test-client-secret"
		environment := "production"

		config, err := InitSDK(clientID, clientSecret, environment)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if config == nil {
			t.Fatal("Expected config to be non-nil")
		}
		if config.clientID != clientID {
			t.Errorf("Expected clientID %s, got %s", clientID, config.clientID)
		}
		if config.clientSecret != clientSecret {
			t.Errorf("Expected clientSecret %s, got %s", clientSecret, config.clientSecret)
		}
		if config.environment != environment {
			t.Errorf("Expected environment %s, got %s", environment, config.environment)
		}
		if config.requestTimeout != 20*time.Second {
			t.Errorf("Expected timeout 20s, got %v", config.requestTimeout)
		}
	})

	t.Run("successful initialization with staging environment", func(t *testing.T) {
		config, err := InitSDK("client", "secret", "staging")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if config.environment != "staging" {
			t.Errorf("Expected environment staging, got %s", config.environment)
		}
	})

	t.Run("default environment when empty", func(t *testing.T) {
		config, err := InitSDK("client", "secret", "")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if config.environment != "production" {
			t.Errorf("Expected default environment production, got %s", config.environment)
		}
	})

	t.Run("error with empty clientID", func(t *testing.T) {
		config, err := InitSDK("", "secret", "production")

		if err == nil {
			t.Error("Expected error for empty clientID")
		}
		if config != nil {
			t.Error("Expected config to be nil on error")
		}
		expectedError := "missing required credentials (ClientID, ClientSecret)"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("error with empty clientSecret", func(t *testing.T) {
		config, err := InitSDK("client", "", "production")

		if err == nil {
			t.Error("Expected error for empty clientSecret")
		}
		if config != nil {
			t.Error("Expected config to be nil on error")
		}
	})

	t.Run("error with both credentials empty", func(t *testing.T) {
		config, err := InitSDK("", "", "production")

		if err == nil {
			t.Error("Expected error for empty credentials")
		}
		if config != nil {
			t.Error("Expected config to be nil on error")
		}
	})

	t.Run("default URLs are set correctly", func(t *testing.T) {
		config, err := InitSDK("client", "secret", "production")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		expectedAPIUrl := "https://api.pingen.com"
		expectedAuthUrl := "https://identity.pingen.com"
		expectedAPIStagingUrl := "https://api-staging.pingen.com"
		expectedAuthStagingUrl := "https://identity-staging.pingen.com"

		if config.apiProductionUrl != expectedAPIUrl {
			t.Errorf("Expected API production URL %s, got %s", expectedAPIUrl, config.apiProductionUrl)
		}
		if config.authProductionUrl != expectedAuthUrl {
			t.Errorf("Expected Auth production URL %s, got %s", expectedAuthUrl, config.authProductionUrl)
		}
		if config.apiStagingUrl != expectedAPIStagingUrl {
			t.Errorf("Expected API staging URL %s, got %s", expectedAPIStagingUrl, config.apiStagingUrl)
		}
		if config.authStagingUrl != expectedAuthStagingUrl {
			t.Errorf("Expected Auth staging URL %s, got %s", expectedAuthStagingUrl, config.authStagingUrl)
		}
	})
}

func TestConfig_SetAPIBaseURL(t *testing.T) {
	t.Run("set custom API base URL", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "production")
		customURL := "https://custom-api.example.com"

		config.SetAPIBaseURL(customURL)

		if config.apiProductionUrl != customURL {
			t.Errorf("Expected API URL %s, got %s", customURL, config.apiProductionUrl)
		}
	})
}

func TestConfig_GetAPIBaseURL(t *testing.T) {
	t.Run("production environment returns production URL", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "production")

		url := config.GetAPIBaseURL()
		expected := "https://api.pingen.com"

		if url != expected {
			t.Errorf("Expected URL %s, got %s", expected, url)
		}
	})

	t.Run("staging environment returns staging URL", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "staging")

		url := config.GetAPIBaseURL()
		expected := "https://api-staging.pingen.com"

		if url != expected {
			t.Errorf("Expected URL %s, got %s", expected, url)
		}
	})

	t.Run("custom environment returns staging URL", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "development")

		url := config.GetAPIBaseURL()
		expected := "https://api-staging.pingen.com"

		if url != expected {
			t.Errorf("Expected URL %s, got %s", expected, url)
		}
	})

	t.Run("custom API URL in production", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "production")
		customURL := "https://custom-api.example.com"
		config.SetAPIBaseURL(customURL)

		url := config.GetAPIBaseURL()

		if url != customURL {
			t.Errorf("Expected custom URL %s, got %s", customURL, url)
		}
	})
}

func TestConfig_GetAuthBaseURL(t *testing.T) {
	t.Run("production environment returns production auth URL", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "production")

		url := config.GetAuthBaseURL()
		expected := "https://identity.pingen.com"

		if url != expected {
			t.Errorf("Expected auth URL %s, got %s", expected, url)
		}
	})

	t.Run("staging environment returns staging auth URL", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "staging")

		url := config.GetAuthBaseURL()
		expected := "https://identity-staging.pingen.com"

		if url != expected {
			t.Errorf("Expected auth URL %s, got %s", expected, url)
		}
	})

	t.Run("custom environment returns staging auth URL", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "test")

		url := config.GetAuthBaseURL()
		expected := "https://identity-staging.pingen.com"

		if url != expected {
			t.Errorf("Expected auth URL %s, got %s", expected, url)
		}
	})
}

func TestConfig_GetClientID(t *testing.T) {
	t.Run("returns correct client ID", func(t *testing.T) {
		clientID := "test-client-12345"
		config, _ := InitSDK(clientID, "secret", "production")

		result := config.GetClientID()

		if result != clientID {
			t.Errorf("Expected client ID %s, got %s", clientID, result)
		}
	})
}

func TestConfig_GetClientSecret(t *testing.T) {
	t.Run("returns correct client secret", func(t *testing.T) {
		clientSecret := "super-secret-key-67890"
		config, _ := InitSDK("client", clientSecret, "production")

		result := config.GetClientSecret()

		if result != clientSecret {
			t.Errorf("Expected client secret %s, got %s", clientSecret, result)
		}
	})
}

func TestConfig_GetRequestTimeout(t *testing.T) {
	t.Run("returns default timeout", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "production")

		timeout := config.GetRequestTimeout()
		expected := 20 * time.Second

		if timeout != expected {
			t.Errorf("Expected timeout %v, got %v", expected, timeout)
		}
	})
}

func TestConfig_GetUserAgent(t *testing.T) {
	t.Run("returns correct user agent", func(t *testing.T) {
		config, _ := InitSDK("client", "secret", "production")

		userAgent := config.GetUserAgent()
		expected := "PINGEN.SDK.GO"

		if userAgent != expected {
			t.Errorf("Expected user agent %s, got %s", expected, userAgent)
		}
	})
}

func TestConfig_validate(t *testing.T) {
	t.Run("valid config passes validation", func(t *testing.T) {
		config := &Config{
			clientID:     "valid-client",
			clientSecret: "valid-secret",
		}

		err := config.validate()

		if err != nil {
			t.Errorf("Expected no error for valid config, got %v", err)
		}
	})

	t.Run("empty clientID fails validation", func(t *testing.T) {
		config := &Config{
			clientID:     "",
			clientSecret: "valid-secret",
		}

		err := config.validate()

		if err == nil {
			t.Error("Expected error for empty clientID")
		}
	})

	t.Run("empty clientSecret fails validation", func(t *testing.T) {
		config := &Config{
			clientID:     "valid-client",
			clientSecret: "",
		}

		err := config.validate()

		if err == nil {
			t.Error("Expected error for empty clientSecret")
		}
	})

	t.Run("validation error message", func(t *testing.T) {
		config := &Config{
			clientID:     "",
			clientSecret: "",
		}

		err := config.validate()
		expected := "missing required credentials (ClientID, ClientSecret)"

		if err.Error() != expected {
			t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
		}
	})
}
