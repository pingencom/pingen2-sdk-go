package pingen2sdk

import (
	"fmt"
	"time"
)

type Config struct {
	clientID          string
	clientSecret      string
	environment       string
	requestTimeout    time.Duration
	apiProductionUrl  string
	authProductionUrl string
	apiStagingUrl     string
	authStagingUrl    string
}

func InitSDK(clientID, clientSecret, environment string) (*Config, error) {
	if environment == "" {
		environment = "production"
	}

	config := &Config{
		clientID:          clientID,
		clientSecret:      clientSecret,
		environment:       environment,
		requestTimeout:    20 * time.Second,
		apiProductionUrl:  "https://api.pingen.com",
		authProductionUrl: "https://identity.pingen.com",
		apiStagingUrl:     "https://api-staging.pingen.com",
		authStagingUrl:    "https://identity-staging.pingen.com",
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) SetAPIBaseURL(url string) {
	c.apiProductionUrl = url
}

func (c *Config) GetAPIBaseURL() string {
	if c.environment == "production" {
		return c.apiProductionUrl
	}
	return c.apiStagingUrl
}

func (c *Config) GetAuthBaseURL() string {
	if c.environment == "production" {
		return c.authProductionUrl
	}
	return c.authStagingUrl
}

func (c *Config) GetClientID() string {
	return c.clientID
}

func (c *Config) GetClientSecret() string {
	return c.clientSecret
}

func (c *Config) GetRequestTimeout() time.Duration {
	return c.requestTimeout
}

func (c *Config) GetUserAgent() string {
	return "PINGEN.SDK.GO"
}

func (c *Config) validate() error {
	if c.clientID == "" || c.clientSecret == "" {
		return fmt.Errorf("missing required credentials (ClientID, ClientSecret)")
	}
	return nil
}
