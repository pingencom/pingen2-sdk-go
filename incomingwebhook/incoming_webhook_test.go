package incomingwebhook_test

import (
	"encoding/json"
	"testing"

	"github.com/pingencom/pingen2-sdk-go/errors"
	"github.com/pingencom/pingen2-sdk-go/incomingwebhook"
)

func TestValidSignature(t *testing.T) {
	payload := `{"data":{"type":"webhook_issues","id":"a3233e48-5e70-4138-95b2-a72d4875016b","attributes":{"reason":"Page limit exceeded","url":"https:\/\/test\/receiver","created_at":"2023-08-03T11:24:39+0200"},"relationships":{"organisation":{"links":{"related":"http:\/\/api-test.v2.pingen.com\/organisations\/2017973a-6403-444d-af05-eb4b2b7f5e2f"},"data":{"type":"organisations","id":"2017973a-6403-444d-af05-eb4b2b7f5e2f"}},"letter":{"links":{"related":"http:\/\/api-test.v2.pingen.com\/organisations\/2017973a-6403-444d-af05-eb4b2b7f5e2f\/letters\/4f31cdb2-bc0d-4db5-a13d-3336958dba02"},"data":{"type":"letters","id":"4f31cdb2-bc0d-4db5-a13d-3336958dba02"}},"event":{"data":{"type":"letters_events","id":"ba08eb5f-413c-4dd1-8ed6-aac2b96124d0"}}}},"included":[{"type":"organisations","id":"2017973a-6403-444d-af05-eb4b2b7f5e2f","attributes":{"name":"Prof. Leopoldo Hahn","status":"active","plan":"free","billing_mode":"postpaid","billing_currency":"CHF","billing_balance":0,"default_country":"CH","edition":"pingen","default_address_position":"left","data_retention_addresses":12,"data_retention_pdf":12,"color":"#0758FF","created_at":"2023-08-03T11:24:39+0200","updated_at":"2023-08-03T11:24:39+0200"},"links":{"self":"http:\/\/api-test.v2.pingen.com\/organisations\/2017973a-6403-444d-af05-eb4b2b7f5e2f"}},{"type":"letters","id":"4f31cdb2-bc0d-4db5-a13d-3336958dba02","attributes":{"status":"validating","file_original_name":"ullam.pdf","file_pages":null,"address":null,"address_position":"left","country":null,"delivery_product":"fast","print_mode":"simplex","print_spectrum":"color","price_currency":null,"price_value":null,"paper_types":null,"fonts":null,"source":"app","tracking_number":null,"submitted_at":null,"created_at":"2023-08-03T11:24:39+0200","updated_at":"2023-08-03T11:24:39+0200"},"links":{"self":"http:\/\/api-test.v2.pingen.com\/organisations\/2017973a-6403-444d-af05-eb4b2b7f5e2f\/letters\/4f31cdb2-bc0d-4db5-a13d-3336958dba02"}},{"type":"letters_events","id":"ba08eb5f-413c-4dd1-8ed6-aac2b96124d0","attributes":{"code":"file_too_many_pages","name":"Page limit exceeded","producer":"Pingen","location":"","has_image":false,"data":[],"emitted_at":"2023-08-03T11:24:39+0200","created_at":"2023-08-03T11:24:39+0200","updated_at":"2023-08-03T11:24:39+0200"}}]}`
	headers := map[string]string{
		"Content-Type": "application/vnd.api+json",
		"Signature":    "812ac7c9776458ce47f1796faec3eca4b15f47b3b9bcfeca3ad90cc190ba0c27",
	}
	secret := "webhook_test_secret"

	webhook := incomingwebhook.IncomingWebhook{}
	event, err := webhook.ConstructEvent(payload, headers, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var eventData struct {
		Data struct {
			Type string `json:"type"`
		} `json:"data"`
	}
	err = json.Unmarshal([]byte(event.Payload), &eventData)
	if err != nil {
		t.Fatalf("failed to unmarshal event payload: %v", err)
	}

	if eventData.Data.Type != "webhook_issues" {
		t.Errorf("expected type 'webhook_issues', got '%s'", eventData.Data.Type)
	}
}

func TestMissingSignatureHeader(t *testing.T) {
	payload := `{"data":{"type":"webhook_issues","id":"309a31e0-1abe-4034-8e7e-1fd473a802fd","attributes":{"reason":"Page limit exceeded","url":"https://5f2e-5-173-206-46.ngrok-free.app/webhook","created_at":"2024-04-26T09:54:38+0200"},"relationships":{"organisation":{"links":{"related":"https://api-integration.pingen.com/organisations/b85c7b52-debb-4b15-b5db-d86b0a4e9bbf"},"data":{"type":"organisations","id":"b85c7b52-debb-4b15-b5db-d86b0a4e9bbf"}},"letter":{"links":{"related":"https://api-integration.pingen.com/organisations/b85c7b52-debb-4b15-b5db-d86b0a4e9bbf/letters/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"},"data":{"type":"letters","id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}},"event":{"data":{"type":"letters_events","id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}}},"included":[]}`
	headers := map[string]string{
		"Content-Type": "application/vnd.api+json",
	}
	secret := "webhook_test_secret"

	webhook := incomingwebhook.IncomingWebhook{}
	_, err := webhook.ConstructEvent(payload, headers, secret)

	if err == nil {
		t.Fatalf("expected an error, but got nil")
	}

	if _, ok := err.(*errors.WebhookSignatureException); !ok {
		t.Fatalf("expected WebhookSignatureException, got %v", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	payload := `{"data":{"type":"webhook_issues","id":"309a31e0-1abe-4034-8e7e-1fd473a802fd","attributes":{"reason":"Page limit exceeded","url":"https://5f2e-5-173-206-46.ngrok-free.app/webhook","created_at":"2024-04-26T09:54:38+0200"},"relationships":{"organisation":{"links":{"related":"https://api-integration.pingen.com/organisations/b85c7b52-debb-4b15-b5db-d86b0a4e9bbf"},"data":{"type":"organisations","id":"b85c7b52-debb-4b15-b5db-d86b0a4e9bbf"}},"letter":{"links":{"related":"https://api-integration.pingen.com/organisations/b85c7b52-debb-4b15-b5db-d86b0a4e9bbf/letters/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"},"data":{"type":"letters","id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}},"event":{"data":{"type":"letters_events","id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}}}},"included":[]}`
	headers := map[string]string{
		"Content-Type": "application/vnd.api+json",
		"Signature":    "wrong99999999999999999999999999999999999999999999999999signature",
	}
	secret := "webhook_test_secret"

	webhook := incomingwebhook.IncomingWebhook{}
	_, err := webhook.ConstructEvent(payload, headers, secret)

	if err == nil {
		t.Fatalf("expected an error, but got nil")
	}

	if _, ok := err.(*errors.WebhookSignatureException); !ok {
		t.Fatalf("expected WebhookSignatureException, got %v", err)
	}
}
