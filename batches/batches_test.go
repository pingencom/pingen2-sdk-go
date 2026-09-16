package batches_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/batches"
	"github.com/pingencom/pingen2-sdk-go/config"
	"github.com/stretchr/testify/assert"
)

const mockBatchResponse = `{
	"data": {
		"id": "test-batch-id",
		"type": "batches",
		"attributes": {
			"name": "Test Batch",
			"channel_type": "post",
			"icon": "document",
			"status": "draft",
			"file_original_name": "test.pdf",
			"letter_count": 5,
			"deliverable_count": 5,
			"address_position": "left",
			"print_mode": "simplex",
			"print_spectrum": "color",
			"price_currency": "CHF",
			"price_value": 2.50,
			"source": "api",
			"submitted_at": "2021-11-19T09:42:48+0100",
			"created_at": "2020-11-19T09:42:48+0100",
			"updated_at": "2020-11-19T09:42:48+0100"
		},
		"relationships": {
			"organisation": {
				"links": {
					"related": "string"
				},
				"data": {
					"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					"type": "organisations"
				}
			},
			"events": {
				"links": {
					"related": {
						"href": "string",
						"meta": {
							"count": 0
						}
					}
				}
			}
		},
		"links": {
			"self": "string"
		},
		"meta": {
			"abilities": {
				"self": {
					"cancel": "state",
					"delete": "state",
					"submit": "state",
					"edit": "state",
					"change-window-position": "state",
					"add-attachment": "state"
				}
			}
		}
	},
	"included": []
}`

const mockBatchCollectionResponse = `{
	"data": [
		{
			"id": "batch-1",
			"type": "batches",
			"attributes": {
				"name": "Batch 1",
				"channel_type": "post",
				"icon": "campaign",
				"status": "draft",
				"file_original_name": "file1.pdf",
				"letter_count": 3,
				"deliverable_count": 3,
				"address_position": "left",
				"print_mode": "simplex",
				"print_spectrum": "color",
				"price_currency": "CHF",
				"price_value": 1.50,
				"source": "api",
				"submitted_at": null,
				"created_at": "2020-11-19T09:42:48+0100",
				"updated_at": "2020-11-19T09:42:48+0100"
			},
			"relationships": {
				"organisation": {
					"links": {
						"related": "string"
					},
					"data": {
						"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
						"type": "organisations"
					}
				},
				"events": {
					"links": {
						"related": {
							"href": "string",
							"meta": {
								"count": 0
							}
						}
					}
				}
			},
			"links": {
				"self": "string"
			}
		},
		{
			"id": "batch-2",
			"type": "batches",
			"attributes": {
				"name": "Batch 2",
				"channel_type": "email",
				"icon": "document",
				"status": "sent",
				"file_original_name": "file2.pdf",
				"letter_count": 7,
				"deliverable_count": 7,
				"address_position": "right",
				"print_mode": "duplex",
				"print_spectrum": "grayscale",
				"price_currency": "CHF",
				"price_value": 3.50,
				"source": "api",
				"submitted_at": "2021-11-19T09:42:48+0100",
				"created_at": "2020-11-19T09:42:48+0100",
				"updated_at": "2020-11-19T09:42:48+0100"
			},
			"relationships": {
				"organisation": {
					"links": {
						"related": "string"
					},
					"data": {
						"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
						"type": "organisations"
					}
				},
				"events": {
					"links": {
						"related": {
							"href": "string",
							"meta": {
								"count": 0
							}
						}
					}
				}
			},
			"links": {
				"self": "string"
			}
		}
	],
	"included": [],
	"links": {
		"first": "first-link",
		"last": "last-link",
		"prev": null,
		"next": null,
		"self": "self-link"
	},
	"meta": {
		"current_page": 1,
		"last_page": 1,
		"per_page": 20,
		"from": 1,
		"to": 2,
		"total": 2
	}
}`

const mockBatchStatisticsResponse = `{
	"data": {
		"id": "test-batch-id",
		"type": "batch_details_statistics",
		"attributes": {
			"letter_validating": 4,
			"letter_groups": [
				{
					"name": "domestic",
					"count": 10
				},
				{
					"name": "international",
					"count": 2
				}
			],
			"letter_countries": [
				{
					"country": "CH",
					"count": 10
				},
				{
					"country": "DE",
					"count": 2
				}
			],
			"letter_regions": [
				{
					"country": "CH",
					"count": 8
				}
			]
		}
	}
}`

type requestPayload struct {
	Data struct {
		ID            string                 `json:"id"`
		Type          string                 `json:"type"`
		Attributes    map[string]interface{} `json:"attributes"`
		Relationships map[string]interface{} `json:"relationships"`
	} `json:"data"`
}

func decodeRequestPayload(t *testing.T, r *http.Request) requestPayload {
	t.Helper()

	var payload requestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	assert.Nil(t, err)

	return payload
}

func setupUnauthorizedServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2")
		w.WriteHeader(http.StatusUnauthorized)

		responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
		_, _ = w.Write([]byte(responseJSON))
	}))
}

func setupBatch(apiBaseURL string) *batches.Batches {
	config, _ := config.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return batches.NewBatches("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", apiRequestor)
}

func TestGetDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := batchClient.GetDetails("test-batch-id", params, headers)

	assert.Nil(t, err)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
	assert.Equal(t, "Test Batch", resp.Data.Attributes.Name)
	assert.Equal(t, "post", resp.Data.Attributes.ChannelType)
	assert.Equal(t, "document", resp.Data.Attributes.Icon)
	assert.Equal(t, "draft", resp.Data.Attributes.Status)
	assert.Equal(t, "test.pdf", resp.Data.Attributes.FileOriginalName)
	assert.Equal(t, 5, resp.Data.Attributes.LetterCount)
	assert.Equal(t, 5, resp.Data.Attributes.DeliverableCount)
}

func TestGetDetails_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := batchClient.GetDetails("test-batch-id", params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchCollectionResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := batchClient.GetCollection(params, headers)

	assert.Nil(t, err)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, "batch-1", resp.Data[0].ID)
	assert.Equal(t, "Batch 1", resp.Data[0].Attributes.Name)
	assert.Equal(t, "campaign", resp.Data[0].Attributes.Icon)
	assert.Equal(t, "post", resp.Data[0].Attributes.ChannelType)
	assert.Equal(t, 3, resp.Data[0].Attributes.LetterCount)
	assert.Equal(t, 3, resp.Data[0].Attributes.DeliverableCount)
	assert.Equal(t, "batch-2", resp.Data[1].ID)
	assert.Equal(t, "email", resp.Data[1].Attributes.ChannelType)
	assert.Equal(t, 7, resp.Data[1].Attributes.LetterCount)
	assert.Equal(t, 7, resp.Data[1].Attributes.DeliverableCount)
	assert.Equal(t, 1, resp.Meta.CurrentPage)
	assert.Equal(t, 20, resp.Meta.PerPage)
	assert.Equal(t, 2, resp.Meta.Total)
}

func TestGetCollection_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := batchClient.GetCollection(params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestUploadAndCreateBatch(t *testing.T) {
	var server *httptest.Server

	counter := 0
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter == 0 {
			assert.Equal(t, "/file-upload", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			w.WriteHeader(http.StatusOK)
			mockUploadResponse := fmt.Sprintf(`{
                "data": {
                    "attributes": {
                        "url": "%s/upload",
                        "url_signature": "mock-signature"
                    }
                }
            }`, server.URL)

			_, _ = w.Write([]byte(mockUploadResponse))
		}

		if counter == 1 {
			assert.Equal(t, "/upload", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)
			w.WriteHeader(http.StatusNoContent)
		}

		if counter == 2 {
			assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			payload := decodeRequestPayload(t, r)
			assert.Equal(t, "batches", payload.Data.Type)
			assert.Equal(t, "Test Upload Batch", payload.Data.Attributes["name"])
			assert.Equal(t, "rocket", payload.Data.Attributes["icon"])
			assert.Equal(t, "post", payload.Data.Attributes["channel_type"])
			assert.Equal(t, "test.zip", payload.Data.Attributes["file_original_name"])
			assert.Equal(t, "left", payload.Data.Attributes["address_position"])
			assert.Equal(t, "zip", payload.Data.Attributes["grouping_type"])
			assert.Equal(t, "file", payload.Data.Attributes["grouping_options_split_type"])
			assert.Nil(t, payload.Data.Relationships)

			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(mockBatchResponse))
		}

		counter++
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.UploadAndCreateBatch(
		"test.zip",
		"Test Upload Batch",
		batches.IconRocket,
		batches.ChannelTypePost,
		"test.zip",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		nil,
		nil,
		nil,
		nil,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestUploadAndCreateBatch_PutError(t *testing.T) {
	var server *httptest.Server

	counter := 0
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter == 0 {
			assert.Equal(t, "/file-upload", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			w.WriteHeader(http.StatusOK)
			mockUploadResponse := fmt.Sprintf(`{
                "data": {
                    "attributes": {
                        "url": "%s/upload",
                        "url_signature": "mock-signature"
                    }
                }
            }`, server.URL)

			_, _ = w.Write([]byte(mockUploadResponse))
		}

		if counter == 1 {
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.Header().Set("X-Request-Id", "requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2")
			w.WriteHeader(http.StatusUnauthorized)

			responseJSON := `{"error":"invalid_client","error_description":"Client authentication failed","message":"Client authentication failed"}`
			_, _ = w.Write([]byte(responseJSON))
		}
		counter++
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)
	separator := "comma"

	_, err := batchClient.UploadAndCreateBatch(
		"test.zip",
		"Test Upload Batch",
		batches.IconRocket,
		batches.ChannelTypePost,
		"test.zip",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		nil,
		&separator,
		nil,
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: Api error (Status Code: 401, Request ID: )"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestUploadAndCreateBatch_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.UploadAndCreateBatch(
		"test.zip",
		"Test Upload Batch",
		batches.IconRocket,
		batches.ChannelTypePost,
		"test.zip",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		nil,
		nil,
		nil,
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestCreateBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		payload := decodeRequestPayload(t, r)
		assert.Equal(t, "batches", payload.Data.Type)
		assert.Equal(t, "New Test Batch", payload.Data.Attributes["name"])
		assert.Equal(t, "rocket", payload.Data.Attributes["icon"])
		assert.Equal(t, "ebill", payload.Data.Attributes["channel_type"])
		assert.Equal(t, "https://example.com/file.pdf", payload.Data.Attributes["file_url"])
		assert.Equal(t, "signature123", payload.Data.Attributes["file_url_signature"])
		assert.Equal(t, "left", payload.Data.Attributes["address_position"])
		assert.Equal(t, "zip", payload.Data.Attributes["grouping_type"])
		assert.Equal(t, "file", payload.Data.Attributes["grouping_options_split_type"])
		assert.Equal(t, float64(5), payload.Data.Attributes["grouping_options_split_size"])
		assert.Equal(t, "comma", payload.Data.Attributes["grouping_options_split_separator"])
		assert.Equal(t, "first_page", payload.Data.Attributes["grouping_options_split_position"])

		assert.Equal(
			t,
			map[string]interface{}{
				"preset": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   "preset-id",
						"type": "presets",
					},
				},
			},
			payload.Data.Relationships,
		)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	splitSize := 5
	splitPos := batches.SplitPositionFirstPage
	separator := "comma"
	relationships := map[string]interface{}{
		"preset": map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "preset-id",
				"type": "presets",
			},
		},
	}

	resp, err := batchClient.CreateBatch(
		"https://example.com/file.pdf",
		"signature123",
		"New Test Batch",
		batches.IconRocket,
		batches.ChannelTypeEbill,
		"new-test.pdf",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		&splitSize,
		&separator,
		&splitPos,
		relationships,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestCreateBatch_DefaultsChannelTypeAndOmitsNilOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := decodeRequestPayload(t, r)

		assert.Equal(t, "batches", payload.Data.Type)
		assert.Equal(t, "post", payload.Data.Attributes["channel_type"])
		assert.NotContains(t, payload.Data.Attributes, "grouping_options_split_size")
		assert.NotContains(t, payload.Data.Attributes, "grouping_options_split_separator")
		assert.NotContains(t, payload.Data.Attributes, "grouping_options_split_position")
		assert.Nil(t, payload.Data.Relationships)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.CreateBatch(
		"https://example.com/file.pdf",
		"signature123",
		"New Test Batch",
		batches.IconRocket,
		"",
		"new-test.pdf",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		nil,
		nil,
		nil,
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestCreateBatch_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.CreateBatch(
		"https://s3.example.com/file/test",
		"signature-123",
		"Test Batch",
		batches.IconDocument,
		batches.ChannelTypePost,
		"test.pdf",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		nil,
		nil,
		nil,
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestSendBatchPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id/send", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		payload := decodeRequestPayload(t, r)
		assert.Equal(t, "test-batch-id", payload.Data.ID)
		assert.Equal(t, "batches_channel_post_send", payload.Data.Type)
		assert.Equal(
			t,
			map[string]interface{}{
				"delivery_product": "registered",
				"print_mode":       "simplex",
				"print_spectrum":   "color",
			},
			payload.Data.Attributes,
		)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.SendBatchPost(
		"test-batch-id",
		batches.DeliveryProductRegistered,
		batches.PrintModeSimplex,
		batches.PrintSpectrumColor,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestSendBatchPost_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.SendBatchPost(
		"test-batch-id",
		batches.DeliveryProductFast,
		batches.PrintModeDuplex,
		batches.PrintSpectrumGrayscale,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestSendBatchEmail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id/send", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		payload := decodeRequestPayload(t, r)
		assert.Equal(t, "test-batch-id", payload.Data.ID)
		assert.Equal(t, "batches_channel_email_send", payload.Data.Type)
		assert.Equal(
			t,
			map[string]interface{}{"delivery_product": "electronic_email"},
			payload.Data.Attributes,
		)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.SendBatchEmail("test-batch-id")

	assert.Nil(t, err)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestSendBatchEmail_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.SendBatchEmail("test-batch-id")

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestSendBatchEbill(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id/send", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		payload := decodeRequestPayload(t, r)
		assert.Equal(t, "test-batch-id", payload.Data.ID)
		assert.Equal(t, "batches_channel_ebill_send", payload.Data.Type)
		assert.Equal(
			t,
			map[string]interface{}{"delivery_product": "electronic_ebill"},
			payload.Data.Attributes,
		)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.SendBatchEbill("test-batch-id")

	assert.Nil(t, err)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestSendBatchEbill_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.SendBatchEbill("test-batch-id")

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestCancelBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id/cancel", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.CancelBatch("test-batch-id")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		body, readErr := io.ReadAll(r.Body)
		assert.Nil(t, readErr)

		expectedPayload := `{
            "data": {
                "id": "test-batch-id",
                "type": "batches",
                "attributes": {
                    "with_letters": true,
                    "with_deliverables": true
                }
            }
        }`
		assert.JSONEq(t, expectedPayload, string(body))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.DeleteBatch("test-batch-id", true)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteBatch_WithoutDeliverables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		payload := decodeRequestPayload(t, r)
		assert.Equal(t, "test-batch-id", payload.Data.ID)
		assert.Equal(t, "batches", payload.Data.Type)
		assert.Equal(t, false, payload.Data.Attributes["with_letters"])
		assert.Equal(t, false, payload.Data.Attributes["with_deliverables"])

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.DeleteBatch("test-batch-id", false)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteBatch_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.DeleteBatch("test-batch-id", true)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestUpdateBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		payload := decodeRequestPayload(t, r)
		assert.Equal(t, "test-batch-id", payload.Data.ID)
		assert.Equal(t, "batches", payload.Data.Type)
		assert.Equal(
			t,
			map[string]interface{}{
				"name": "Renamed Batch",
				"icon": "rocket",
			},
			payload.Data.Attributes,
		)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.UpdateBatch("test-batch-id", "Renamed Batch", batches.IconRocket)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestUpdateBatch_OmitsEmptyAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := decodeRequestPayload(t, r)

		assert.Equal(t, "batches", payload.Data.Type)
		assert.Equal(t, map[string]interface{}{"name": "Only Name"}, payload.Data.Attributes)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.UpdateBatch("test-batch-id", "Only Name", "")

	assert.Nil(t, err)
}

func TestUpdateBatch_IconOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := decodeRequestPayload(t, r)

		assert.Equal(t, map[string]interface{}{"icon": "crown"}, payload.Data.Attributes)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.UpdateBatch("test-batch-id", "", batches.IconCrown)

	assert.Nil(t, err)
}

func TestUpdateBatch_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.UpdateBatch("test-batch-id", "Renamed Batch", batches.IconRocket)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestGetStatistics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id/statistics", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchStatisticsResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.GetStatistics("test-batch-id")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
	assert.Equal(t, "batch_details_statistics", resp.Data.Type)
	assert.Equal(t, 4, resp.Data.Attributes.LetterValidating)

	assert.Len(t, resp.Data.Attributes.LetterGroups, 2)
	assert.Equal(t, "domestic", resp.Data.Attributes.LetterGroups[0].Name)
	assert.Equal(t, 10, resp.Data.Attributes.LetterGroups[0].Count)
	assert.Equal(t, "international", resp.Data.Attributes.LetterGroups[1].Name)
	assert.Equal(t, 2, resp.Data.Attributes.LetterGroups[1].Count)

	assert.Len(t, resp.Data.Attributes.LetterCountries, 2)
	assert.Equal(t, "CH", resp.Data.Attributes.LetterCountries[0].Country)
	assert.Equal(t, 10, resp.Data.Attributes.LetterCountries[0].Count)
	assert.Equal(t, "DE", resp.Data.Attributes.LetterCountries[1].Country)

	assert.Len(t, resp.Data.Attributes.LetterRegions, 1)
	assert.Equal(t, "CH", resp.Data.Attributes.LetterRegions[0].Country)
	assert.Equal(t, 8, resp.Data.Attributes.LetterRegions[0].Count)
}

func TestGetStatistics_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.GetStatistics("test-batch-id")

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestEnum_Constants(t *testing.T) {
	assert.Equal(t, "document", string(batches.IconDocument))
	assert.Equal(t, "rocket", string(batches.IconRocket))
	assert.Equal(t, "campaign", string(batches.IconCampaign))
	assert.Equal(t, "megaphone", string(batches.IconMegaphone))
	assert.Equal(t, "wave-hand", string(batches.IconWaveHand))
	assert.Equal(t, "flash", string(batches.IconFlash))
	assert.Equal(t, "bell", string(batches.IconBell))
	assert.Equal(t, "percent-tag", string(batches.IconPercentTag))
	assert.Equal(t, "percent-badge", string(batches.IconPercentBadge))
	assert.Equal(t, "present", string(batches.IconPresent))
	assert.Equal(t, "receipt", string(batches.IconReceipt))
	assert.Equal(t, "information", string(batches.IconInformation))
	assert.Equal(t, "calendar", string(batches.IconCalendar))
	assert.Equal(t, "newspaper", string(batches.IconNewspaper))
	assert.Equal(t, "crown", string(batches.IconCrown))
	assert.Equal(t, "virus", string(batches.IconVirus))

	assert.Equal(t, "left", string(batches.AddressPositionLeft))
	assert.Equal(t, "right", string(batches.AddressPositionRight))

	assert.Equal(t, "zip", string(batches.GroupingTypeZip))
	assert.Equal(t, "merge", string(batches.GroupingTypeMerge))

	assert.Equal(t, "file", string(batches.SplitTypeFile))
	assert.Equal(t, "page", string(batches.SplitTypePage))
	assert.Equal(t, "custom", string(batches.SplitTypeCustom))
	assert.Equal(t, "qr_invoice", string(batches.SplitTypeQRInvoice))

	assert.Equal(t, "first_page", string(batches.SplitPositionFirstPage))
	assert.Equal(t, "last_page", string(batches.SplitPositionLastPage))

	assert.Equal(t, "post", string(batches.ChannelTypePost))
	assert.Equal(t, "ebill", string(batches.ChannelTypeEbill))
	assert.Equal(t, "email", string(batches.ChannelTypeEmail))

	assert.Equal(t, "fast", string(batches.DeliveryProductFast))
	assert.Equal(t, "cheap", string(batches.DeliveryProductCheap))
	assert.Equal(t, "bulk", string(batches.DeliveryProductBulk))
	assert.Equal(t, "premium", string(batches.DeliveryProductPremium))
	assert.Equal(t, "registered", string(batches.DeliveryProductRegistered))
	assert.Equal(t, "electronic_email", string(batches.DeliveryProductElectronicEmail))
	assert.Equal(t, "electronic_ebill", string(batches.DeliveryProductElectronicEbill))

	assert.Equal(t, "simplex", string(batches.PrintModeSimplex))
	assert.Equal(t, "duplex", string(batches.PrintModeDuplex))

	assert.Equal(t, "color", string(batches.PrintSpectrumColor))
	assert.Equal(t, "grayscale", string(batches.PrintSpectrumGrayscale))
}
