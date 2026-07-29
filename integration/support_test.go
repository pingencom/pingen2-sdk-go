//go:build integration
// +build integration

package integration

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pingencom/pingen2-sdk-go/errors"
	"github.com/pingencom/pingen2-sdk-go/response"
)

const Scope = "letter batch webhook organisation_read email ebill"

// Document names sent to the API. The staging environment recognises the magic
// `simulate_cancellable` suffix and keeps such deliveries in a state that can
// be cancelled, which lets the suite exercise the cancel flow deterministically.
const (
	FileName            = "test.pdf"
	FileNameCancellable = "test_simulate_cancellable.pdf"
)

const (
	envClientID         = "PINGEN2_CLIENT_ID"
	envClientSecret     = "PINGEN2_CLIENT_SECRET"
	envOrganisationID   = "PINGEN2_ORGANIZATION_ID"
	envOrganisationName = "PINGEN2_ORGANIZATION_NAME"
	envUseStaging       = "PINGEN2_USE_STAGING"
)

const missingCredentialsMessage = "Integration credentials not configured. " +
	"Copy .env.example to .env and fill in PINGEN2_CLIENT_ID / PINGEN2_CLIENT_SECRET."

const ebillChannelMissingCode = "conflict_missing_configuration"

type credentials struct {
	ClientID         string
	ClientSecret     string
	OrganisationID   string
	OrganisationName string
	UseStaging       bool
}

func loadCredentials() credentials {
	dotenv := parseDotenv(filepath.Join(repoRoot(), ".env"))

	value := func(key string) string {
		if fromEnv := os.Getenv(key); fromEnv != "" {
			return fromEnv
		}
		return dotenv[key]
	}

	return credentials{
		ClientID:         value(envClientID),
		ClientSecret:     value(envClientSecret),
		OrganisationID:   value(envOrganisationID),
		OrganisationName: value(envOrganisationName),
		UseStaging:       useStaging(value(envUseStaging)),
	}
}

func requireCredentials(t *testing.T) credentials {
	t.Helper()

	creds := loadCredentials()
	if creds.ClientID == "" || creds.ClientSecret == "" {
		t.Skip(missingCredentialsMessage)
	}

	return creds
}

// Defaults to staging: the suite must never run against production.
func environment(c credentials) string {
	if c.UseStaging {
		return "staging"
	}

	return "production"
}

func useStaging(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "0", "false", "no", "off":
		return false
	default:
		return true
	}
}

func documentPath(fileName string) string {
	return filepath.Join(repoRoot(), "testdata", fileName)
}

func documentName(path string) string {
	return filepath.Base(path)
}

func buildEmailMetaData() map[string]interface{} {
	return map[string]interface{}{
		"sender_name":     "Pingen Test",
		"recipient_email": "grzegorz.morgas@pingen.com",
		"recipient_name":  "Test Recipient",
		"reply_email":     "noreply@example.com",
		"reply_name":      "Reply Test",
		"subject":         "Integration Test Email",
		"content":         "Dear Recipient\\n\\nThis is an integration test.\\n\\nBest regards",
	}
}

func buildEbillMetaData() map[string]interface{} {
	today := time.Now()

	return map[string]interface{}{
		"invoice_number":       fmt.Sprintf("INV-%s", uniqueSuffix()),
		"invoice_date":         today.Format("2006-01-02"),
		"invoice_due_date":     today.AddDate(0, 0, 30).Format("2006-01-02"),
		"recipient_identifier": "41100000014283293",
	}
}

func statusCode(t *testing.T, resp interface{}) int {
	t.Helper()

	defaultResponse, ok := resp.(*response.DefaultResponse)
	if !ok {
		t.Fatalf("expected *response.DefaultResponse, got %T", resp)
		return 0
	}

	return defaultResponse.StatusCode
}

func skipIfEbillChannelMissing(t *testing.T, err *errors.PingenError) {
	t.Helper()

	if err == nil {
		return
	}

	if strings.Contains(errorBody(err), ebillChannelMissingCode) {
		t.Skipf("Organisation has no ebill channel configured: %s", errorBody(err))
	}
}

func errorBody(err *errors.PingenError) string {
	if err == nil {
		return ""
	}

	body, marshalErr := json.Marshal(err.JSONBody)
	if marshalErr != nil {
		return fmt.Sprintf("%s: %v", err.Message, err.JSONBody)
	}

	return fmt.Sprintf("%s: %s", err.Message, string(body))
}

func uniqueSuffix() string {
	buffer := make([]byte, 6)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%012d", time.Now().UnixNano()%1e12)
	}

	return hex.EncodeToString(buffer)
}

var (
	repoRootOnce  sync.Once
	repoRootValue string
)

func repoRoot() string {
	repoRootOnce.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			repoRootValue = "."
			return
		}

		for {
			if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
				repoRootValue = dir
				return
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				repoRootValue = "."
				return
			}

			dir = parent
		}
	})

	return repoRootValue
}

func parseDotenv(path string) map[string]string {
	values := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		return values
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		key, value, _ := strings.Cut(line, "=")
		value = strings.Trim(strings.Trim(strings.TrimSpace(value), `"`), `'`)
		values[strings.TrimSpace(key)] = value
	}

	return values
}
