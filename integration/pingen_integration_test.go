//go:build integration
// +build integration

package integration

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pingencom/pingen2-sdk-go/batches"
	"github.com/pingencom/pingen2-sdk-go/ebills"
)

func uploadEbill(t *testing.T, s *session, path string, autoSend bool) ebills.EbillResponse {
	t.Helper()

	response, err := s.Ebills.UploadAndCreate(
		path,
		documentName(path),
		autoSend,
		buildEbillMetaData(),
		nil,
	)
	if err != nil {
		skipIfEbillChannelMissing(t, err)
		t.Fatalf("failed to create e-bill: %v", err)
	}

	return response
}

func TestOrganisations(t *testing.T) {
	s := newSession(t)

	t.Run("1_list_organisations", func(t *testing.T) {
		response, err := s.Organisations.GetCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		require.GreaterOrEqual(t, len(response.Data), 1)
		for _, item := range response.Data {
			assert.NotEmpty(t, item.ID)
		}
	})

	t.Run("2_list_organisations_paginated", func(t *testing.T) {
		response, err := s.Organisations.GetCollection(
			map[string]string{"page[number]": "1", "page[limit]": "5"},
			nil,
		)
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("3_get_organisation_by_id", func(t *testing.T) {
		response, err := s.Organisations.GetDetails(s.OrganisationID(), nil, nil)
		if err != nil {
			t.Fatalf("failed to get organisation details: %v", err)
		}

		assert.Equal(t, s.OrganisationID(), response.Data.ID)
	})
}

func TestLettersHappyCase(t *testing.T) {
	s := newSession(t)

	var letterID string

	t.Run("1_list_letters", func(t *testing.T) {
		response, err := s.Letters.GetCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to list letters: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("2_list_letters_paginated", func(t *testing.T) {
		response, err := s.Letters.GetCollection(
			map[string]string{"page[number]": "1", "page[limit]": "3"},
			nil,
		)
		if err != nil {
			t.Fatalf("failed to list letters: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("3_create_letter", func(t *testing.T) {
		path := documentPath(FileName)

		response, err := s.Letters.UploadAndCreate(
			path,
			documentName(path),
			"left",
			true,
			"cheap",
			"simplex",
			"grayscale",
			"",
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create letter: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)
		assert.Equal(t, "validating", response.Data.Attributes.Status)

		letterID = response.Data.ID

		detail, err := s.Letters.GetDetails(letterID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get letter details: %v", err)
		}
		assert.Equal(t, letterID, detail.Data.ID)

		t.Log("sleep 30 seconds so the collection contains the newly created letter")
		time.Sleep(30 * time.Second)

		collection, err := s.Letters.GetCollection(
			map[string]string{"sort": "-created_at", "page[number]": "1", "page[limit]": "20"},
			nil,
		)
		if err != nil {
			t.Fatalf("failed to list letters: %v", err)
		}

		ids := make([]string, 0, len(collection.Data))
		for _, item := range collection.Data {
			ids = append(ids, item.ID)
		}
		assert.Contains(t, ids, letterID)
	})

	t.Run("4_get_letter_by_id", func(t *testing.T) {
		if letterID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.Letters.GetDetails(letterID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get letter details: %v", err)
		}

		assert.Equal(t, letterID, response.Data.ID)
		t.Logf("Letter status: %s", response.Data.Attributes.Status)
	})

	t.Run("5_get_letter_events", func(t *testing.T) {
		if letterID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.LetterEvents.GetCollection(letterID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get letter events: %v", err)
		}

		t.Logf("Letter events: %d", len(response.Data))
	})

	t.Run("6_get_letter_file", func(t *testing.T) {
		if letterID == "" {
			t.Skip("requires the create subtest to have run")
		}

		file, err := s.Letters.GetFile(letterID)
		if err != nil {
			t.Fatalf("failed to get letter file: %v", err)
		}

		content, readErr := io.ReadAll(io.LimitReader(file, 1024))
		require.NoError(t, readErr)
		assert.NotEmpty(t, content)
		require.NoError(t, file.Close())
	})

	t.Run("7_calculate_letter_price", func(t *testing.T) {
		response, err := s.Letters.CalculatePrice(
			"CH",
			[]string{"normal", "normal"},
			"simplex",
			"grayscale",
			"cheap",
		)
		if err != nil {
			t.Fatalf("failed to calculate letter price: %v", err)
		}

		assert.NotEmpty(t, response.Data.ID)
		t.Logf(
			"Price calculator: currency=%s, price=%v",
			response.Data.Attributes.Currency,
			response.Data.Attributes.Price,
		)
	})

	t.Run("8_get_letter_sent_events", func(t *testing.T) {
		_, err := s.LetterEvents.GetSentCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to get sent letter events: %v", err)
		}
	})

	t.Run("9_get_letter_delivered_events", func(t *testing.T) {
		_, err := s.LetterEvents.GetDeliveredCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to get delivered letter events: %v", err)
		}
	})

	t.Run("10_get_letter_issue_events", func(t *testing.T) {
		_, err := s.LetterEvents.GetIssueCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to get issue letter events: %v", err)
		}
	})

	t.Run("11_get_letter_undeliverable_events", func(t *testing.T) {
		_, err := s.LetterEvents.GetUndeliverableCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to get undeliverable letter events: %v", err)
		}
	})
}

func TestLettersCancelCase(t *testing.T) {
	s := newSession(t)

	var letterID string

	t.Run("1_create_letter", func(t *testing.T) {
		path := documentPath(FileNameCancellable)

		response, err := s.Letters.UploadAndCreate(
			path,
			documentName(path),
			"left",
			true,
			"cheap",
			"simplex",
			"grayscale",
			"",
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create letter: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)
		assert.Equal(t, "validating", response.Data.Attributes.Status)

		letterID = response.Data.ID
	})

	t.Run("2_cancel_letter", func(t *testing.T) {
		if letterID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the letter reaches a cancellable state")
		time.Sleep(10 * time.Second)

		resp, err := s.Letters.Cancel(letterID)
		if err != nil {
			t.Fatalf("failed to cancel letter: %v", err)
		}

		assert.Equal(t, 202, statusCode(t, resp))
	})
}

func TestLettersDeleteCase(t *testing.T) {
	s := newSession(t)

	var letterID string

	t.Run("1_create_letter", func(t *testing.T) {
		path := documentPath(FileName)

		response, err := s.Letters.UploadAndCreate(
			path,
			documentName(path),
			"right",
			false,
			"cheap",
			"simplex",
			"grayscale",
			"",
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create letter: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)
		assert.Equal(t, "validating", response.Data.Attributes.Status)
		time.Sleep(5 * time.Second)

		letterID = response.Data.ID
	})

	t.Run("2_delete_letter", func(t *testing.T) {
		if letterID == "" {
			t.Skip("requires the create subtest to have run")
		}

		resp, err := s.Letters.Delete(letterID)
		if err != nil {
			t.Fatalf("failed to delete letter: %v", err)
		}

		assert.Equal(t, 204, statusCode(t, resp))
		t.Logf("Deleted letter: %s", letterID)
	})
}

func TestBatchesHappyCase(t *testing.T) {
	s := newSession(t)

	var batchID string

	t.Run("1_list_batches", func(t *testing.T) {
		response, err := s.Batches.GetCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to list batches: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("2_list_batches_paginated", func(t *testing.T) {
		response, err := s.Batches.GetCollection(
			map[string]string{"page[number]": "1", "page[limit]": "3"},
			nil,
		)
		if err != nil {
			t.Fatalf("failed to list batches: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("3_create_batch", func(t *testing.T) {
		path := documentPath(FileName)
		splitSize := 2

		response, err := s.Batches.UploadAndCreateBatch(
			path,
			"Integration Test Batch",
			batches.IconDocument,
			batches.ChannelTypePost,
			documentName(path),
			batches.AddressPositionLeft,
			batches.GroupingTypeMerge,
			batches.SplitTypeQRInvoice,
			&splitSize,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create batch: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)

		batchID = response.Data.ID
		t.Logf("Created batch: %s (status: %s)", batchID, response.Data.Attributes.Status)
	})

	t.Run("4_get_batch_by_id", func(t *testing.T) {
		if batchID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.Batches.GetDetails(batchID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get batch details: %v", err)
		}

		assert.Equal(t, batchID, response.Data.ID)
		t.Logf("Batch status: %s", response.Data.Attributes.Status)
	})

	t.Run("5_update_batch", func(t *testing.T) {
		if batchID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the batch reaches an updatable state")
		time.Sleep(10 * time.Second)

		response, err := s.Batches.UpdateBatch(batchID, "Updated Integration Batch", batches.IconRocket)
		if err != nil {
			t.Fatalf("failed to update batch: %v", err)
		}

		assert.Equal(t, batchID, response.Data.ID)
		assert.Equal(t, "Updated Integration Batch", response.Data.Attributes.Name)
		assert.Equal(t, string(batches.IconRocket), response.Data.Attributes.Icon)
	})

	t.Run("6_get_batch_events", func(t *testing.T) {
		if batchID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.BatchEvents.GetCollection(batchID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get batch events: %v", err)
		}

		t.Logf("Batch events: %d", len(response.Data))
	})

	t.Run("7_get_batch_statistics", func(t *testing.T) {
		if batchID == "" {
			t.Skip("requires the create subtest to have run")
		}

		_, err := s.Batches.GetStatistics(batchID)
		if err != nil {
			t.Fatalf("failed to get batch statistics: %v", err)
		}
	})

	t.Run("8_get_batch_channel_and_deliverable_count", func(t *testing.T) {
		if batchID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.Batches.GetDetails(batchID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get batch details: %v", err)
		}

		assert.Equal(t, string(batches.ChannelTypePost), response.Data.Attributes.ChannelType)
		assert.GreaterOrEqual(t, response.Data.Attributes.DeliverableCount, 0)
		t.Logf(
			"Batch channel: %s, deliverables: %d",
			response.Data.Attributes.ChannelType,
			response.Data.Attributes.DeliverableCount,
		)
	})
}

func TestBatchesDeleteCase(t *testing.T) {
	s := newSession(t)

	var batchID string

	t.Run("1_create_batch", func(t *testing.T) {
		path := documentPath(FileNameCancellable)
		splitSize := 2

		response, err := s.Batches.UploadAndCreateBatch(
			path,
			"Integration Test Batch",
			batches.IconDocument,
			batches.ChannelTypePost,
			documentName(path),
			batches.AddressPositionLeft,
			batches.GroupingTypeMerge,
			batches.SplitTypeQRInvoice,
			&splitSize,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create batch: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)

		batchID = response.Data.ID
		t.Logf("Created batch: %s", batchID)
	})

	t.Run("2_delete_batch", func(t *testing.T) {
		if batchID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the batch reaches a deletable state")
		time.Sleep(10 * time.Second)

		resp, err := s.Batches.DeleteBatch(batchID, true)
		if err != nil {
			t.Fatalf("failed to delete batch: %v", err)
		}

		assert.Equal(t, 204, statusCode(t, resp))
		t.Logf("Deleted batch: %s", batchID)
	})
}

func TestWebhooks(t *testing.T) {
	s := newSession(t)

	var webhookID string

	t.Run("1_list_webhooks", func(t *testing.T) {
		response, err := s.Webhooks.GetCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to list webhooks: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("2_create_webhook", func(t *testing.T) {
		response, err := s.Webhooks.Create(
			"issues",
			"https://httpbin.org/post",
			"integration-test-signing-key-32c",
		)
		if err != nil {
			t.Fatalf("failed to create webhook: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)

		webhookID = response.Data.ID
		t.Logf("Created webhook: %s (url: %s)", webhookID, response.Data.Attributes.URL)
	})

	t.Run("3_get_webhook_by_id", func(t *testing.T) {
		if webhookID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.Webhooks.GetDetails(webhookID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get webhook details: %v", err)
		}

		assert.Equal(t, webhookID, response.Data.ID)
		t.Logf("Webhook url: %s", response.Data.Attributes.URL)
	})

	t.Run("4_delete_webhook", func(t *testing.T) {
		if webhookID == "" {
			t.Skip("requires the create subtest to have run")
		}

		resp, err := s.Webhooks.Delete(webhookID)
		if err != nil {
			t.Fatalf("failed to delete webhook: %v", err)
		}

		assert.Equal(t, 204, statusCode(t, resp))
		t.Logf("Deleted webhook: %s", webhookID)
	})
}

func TestEmailsHappyCase(t *testing.T) {
	s := newSession(t)

	var emailID string

	t.Run("1_list_emails", func(t *testing.T) {
		response, err := s.Emails.GetCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to list emails: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("2_create_email", func(t *testing.T) {
		path := documentPath(FileName)

		response, err := s.Emails.UploadAndCreate(
			path,
			documentName(path),
			true,
			buildEmailMetaData(),
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create email: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)

		emailID = response.Data.ID
		t.Logf("Created email: %s (status: %s)", emailID, response.Data.Attributes.Status)
	})

	t.Run("3_get_email_by_id", func(t *testing.T) {
		if emailID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.Emails.GetDetails(emailID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get email details: %v", err)
		}

		assert.Equal(t, emailID, response.Data.ID)
		t.Logf("Email status: %s", response.Data.Attributes.Status)
	})

	t.Run("4_get_email_events", func(t *testing.T) {
		if emailID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.EmailEvents.GetCollection(emailID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get email events: %v", err)
		}

		t.Logf("Email events: %d", len(response.Data))
	})

	t.Run("5_get_email_file", func(t *testing.T) {
		if emailID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 5 seconds so the email reaches a retrievable state")
		time.Sleep(5 * time.Second)

		file, err := s.Emails.GetFile(emailID)
		if err != nil {
			t.Fatalf("failed to get email file: %v", err)
		}

		content, readErr := io.ReadAll(io.LimitReader(file, 1024))
		require.NoError(t, readErr)
		assert.NotEmpty(t, content)
		require.NoError(t, file.Close())
	})
}

func TestEmailsCancelCase(t *testing.T) {
	s := newSession(t)

	var emailID string

	t.Run("1_create_email", func(t *testing.T) {
		path := documentPath(FileNameCancellable)

		response, err := s.Emails.UploadAndCreate(
			path,
			documentName(path),
			true,
			buildEmailMetaData(),
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create email: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)
		emailID = response.Data.ID
	})

	t.Run("2_cancel_email", func(t *testing.T) {
		if emailID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the email reaches a cancellable state")
		time.Sleep(10 * time.Second)

		resp, err := s.Emails.Cancel(emailID)
		if err != nil {
			t.Fatalf("failed to cancel email: %v", err)
		}

		assert.Equal(t, 202, statusCode(t, resp))
	})
}

func TestEmailsDeleteCase(t *testing.T) {
	s := newSession(t)

	var emailID string

	t.Run("1_create_email", func(t *testing.T) {
		path := documentPath(FileName)

		response, err := s.Emails.UploadAndCreate(
			path,
			documentName(path),
			false,
			buildEmailMetaData(),
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create email: %v", err)
		}

		require.NotEmpty(t, response.Data.ID)
		emailID = response.Data.ID
	})

	t.Run("2_delete_email", func(t *testing.T) {
		if emailID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the email reaches a deletable state")
		time.Sleep(10 * time.Second)

		resp, err := s.Emails.Delete(emailID)
		if err != nil {
			t.Fatalf("failed to delete email: %v", err)
		}

		assert.Equal(t, 204, statusCode(t, resp))
		t.Logf("Deleted email: %s", emailID)
	})
}

func TestEbillsHappyCase(t *testing.T) {
	s := newSession(t)

	var ebillID string

	t.Run("1_list_ebills", func(t *testing.T) {
		response, err := s.Ebills.GetCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to list e-bills: %v", err)
		}

		assert.NotNil(t, response.Data)
	})

	t.Run("2_create_ebill", func(t *testing.T) {
		response := uploadEbill(t, s, documentPath(FileName), false)

		require.NotEmpty(t, response.Data.ID)

		ebillID = response.Data.ID
		t.Logf("Created e-bill: %s (status: %s)", ebillID, response.Data.Attributes.Status)
	})

	t.Run("3_get_ebill_by_id", func(t *testing.T) {
		if ebillID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.Ebills.GetDetails(ebillID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get e-bill details: %v", err)
		}

		assert.Equal(t, ebillID, response.Data.ID)
		t.Logf("E-Bill status: %s", response.Data.Attributes.Status)
	})

	t.Run("4_get_ebill_events", func(t *testing.T) {
		if ebillID == "" {
			t.Skip("requires the create subtest to have run")
		}

		response, err := s.EbillEvents.GetCollection(ebillID, nil, nil)
		if err != nil {
			t.Fatalf("failed to get e-bill events: %v", err)
		}

		t.Logf("E-Bill events: %d", len(response.Data))
	})

	t.Run("5_get_ebill_file", func(t *testing.T) {
		if ebillID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the e-bill reaches a retrievable state")
		time.Sleep(10 * time.Second)

		file, err := s.Ebills.GetFile(ebillID)
		if err != nil {
			t.Fatalf("failed to get e-bill file: %v", err)
		}

		content, readErr := io.ReadAll(io.LimitReader(file, 1024))
		require.NoError(t, readErr)
		assert.NotEmpty(t, content)
		require.NoError(t, file.Close())
	})
}

func TestEbillsCancelCase(t *testing.T) {
	s := newSession(t)

	var ebillID string

	t.Run("1_create_ebill", func(t *testing.T) {
		response := uploadEbill(t, s, documentPath(FileNameCancellable), true)

		require.NotEmpty(t, response.Data.ID)
		ebillID = response.Data.ID
	})

	t.Run("2_cancel_ebill", func(t *testing.T) {
		if ebillID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the e-bill reaches a cancellable state")
		time.Sleep(10 * time.Second)

		resp, err := s.Ebills.Cancel(ebillID)
		if err != nil {
			t.Fatalf("failed to cancel e-bill: %v", err)
		}

		assert.Equal(t, 202, statusCode(t, resp))
	})
}

func TestEbillsDeleteCase(t *testing.T) {
	s := newSession(t)

	var ebillID string

	t.Run("1_create_ebill", func(t *testing.T) {
		response := uploadEbill(t, s, documentPath(FileName), false)

		require.NotEmpty(t, response.Data.ID)
		ebillID = response.Data.ID
	})

	t.Run("2_delete_ebill", func(t *testing.T) {
		if ebillID == "" {
			t.Skip("requires the create subtest to have run")
		}

		t.Log("sleep 10 seconds so the e-bill reaches a deletable state")
		time.Sleep(10 * time.Second)

		resp, err := s.Ebills.Delete(ebillID)
		if err != nil {
			t.Fatalf("failed to delete e-bill: %v", err)
		}

		assert.Equal(t, 204, statusCode(t, resp))
		t.Logf("Deleted e-bill: %s", ebillID)
	})
}

func TestUser(t *testing.T) {
	s := newSession(t)

	t.Run("1_get_user", func(t *testing.T) {
		response, err := s.Users.GetDetails(nil, nil)
		if err != nil {
			t.Fatalf("failed to get user details: %v", err)
		}

		assert.NotEmpty(t, response.Data.ID)
		assert.NotEmpty(t, response.Data.Attributes.Email)
		t.Logf(
			"User: %s %s (%s)",
			response.Data.Attributes.FirstName,
			response.Data.Attributes.LastName,
			response.Data.Attributes.Email,
		)
	})

	t.Run("2_get_user_associations", func(t *testing.T) {
		response, err := s.UserAssociations.GetCollection(nil, nil)
		if err != nil {
			t.Fatalf("failed to get user associations: %v", err)
		}

		t.Logf("User associations: %d", len(response.Data))
	})
}
