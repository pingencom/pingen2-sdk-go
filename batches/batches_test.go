
package batches_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/batches"
	"github.com/stretchr/testify/assert"
)

const mockBatchResponse = `{
	"data": {
		"id": "test-batch-id",
		"type": "batches",
		"attributes": {
			"name": "Test Batch",
			"icon": "document",
			"status": "draft",
			"file_original_name": "test.pdf",
			"letter_count": 5,
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
				"icon": "campaign",
				"status": "draft",
				"file_original_name": "file1.pdf",
				"letter_count": 3,
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
				"icon": "document",
				"status": "sent",
				"file_original_name": "file2.pdf",
				"letter_count": 7,
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
		"id": "stats-id",
		"type": "batch_statistics",
		"attributes": {
			"total_letters": 100,
			"processed_letters": 95,
			"sent_letters": 90,
			"cancelled_letters": 5,
			"error_letters": 0
		}
	}
}`

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
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
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
	assert.Equal(t, "document", resp.Data.Attributes.Icon)
	assert.Equal(t, "draft", resp.Data.Attributes.Status)
	assert.Equal(t, "test.pdf", resp.Data.Attributes.FileOriginalName)
	assert.Equal(t, 5, resp.Data.Attributes.LetterCount)
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
	assert.Equal(t, 3, resp.Data[0].Attributes.LetterCount)
	assert.Equal(t, "batch-2", resp.Data[1].ID)
	assert.Equal(t, 7, resp.Data[1].Attributes.LetterCount)
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

			body, _ := io.ReadAll(r.Body)
			assert.Contains(t, string(body), `"name":"Test Upload Batch"`)
			assert.Contains(t, string(body), `"icon":"rocket"`)
			assert.Contains(t, string(body), `"file_original_name":"test.zip"`)
			assert.Contains(t, string(body), `"address_position":"left"`)
			assert.Contains(t, string(body), `"grouping_type":"zip"`)

			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(mockBatchResponse))
		}

		counter++
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	splitSize := 10
	splitPos := batches.SplitPositionFirstPage

	resp, err := batchClient.UploadAndCreateBatch(
		"test.zip",
		"Test Upload Batch",
		batches.IconRocket,
		"test.zip",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		&splitSize,
		nil,
		&splitPos,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestUploadAndCreateBatch_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	_, err := batchClient.UploadAndCreateBatch(
		"test.zip",
		"Test Upload Batch",
		batches.IconRocket,
		"test.zip",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
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

		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), `"name":"New Test Batch"`)
		assert.Contains(t, string(body), `"icon":"rocket"`)
		assert.Contains(t, string(body), `"address_position":"left"`)
		assert.Contains(t, string(body), `"grouping_type":"zip"`)
		assert.Contains(t, string(body), `"grouping_options_split_type":"file"`)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	splitSize := 5
	splitPos := batches.SplitPositionFirstPage

	resp, err := batchClient.CreateBatch(
		"https://example.com/file.pdf",
		"signature123",
		"New Test Batch",
		batches.IconRocket,
		"new-test.pdf",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		&splitSize,
		nil,
		&splitPos,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
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
		"test.pdf",
		batches.AddressPositionLeft,
		batches.GroupingTypeZip,
		batches.SplitTypeFile,
		nil,
		nil,
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestSendBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id/send", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), `"print_mode":"simplex"`)
		assert.Contains(t, string(body), `"print_spectrum":"color"`)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	deliveryProducts := map[string]string{
		"standard": "economy",
	}

	resp, err := batchClient.SendBatch("test-batch-id", deliveryProducts, "simplex", "color")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestSendBatch_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	deliveryProducts := map[string]string{
		"standard": "economy",
	}

	_, err := batchClient.SendBatch("test-batch-id", deliveryProducts, "simplex", "color")

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

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	resp, err := batchClient.DeleteBatch("test-batch-id")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestEditBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/batches/test-batch-id", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		expectedPayload := `{
            "data": {
                "id": "test-batch-id",
                "type": "batches",
                "attributes": {
                    "paper_types": ["normal", "qr"]
                }
            }
        }`
		assert.JSONEq(t, expectedPayload, string(body))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockBatchResponse))
	}))
	defer server.Close()

	batchClient := setupBatch(server.URL)

	paperTypes := []string{"normal", "qr"}
	resp, err := batchClient.EditBatch("test-batch-id", paperTypes)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-batch-id", resp.Data.ID)
}

func TestEditBatch_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	batchClient := setupBatch(server.URL)

	paperTypes := []string{"normal", "qr"}
	_, err := batchClient.EditBatch("test-batch-id", paperTypes)

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
	assert.Equal(t, "stats-id", resp.Data.ID)
	assert.Equal(t, 100, resp.Data.Attributes.TotalLetters)
	assert.Equal(t, 95, resp.Data.Attributes.ProcessedLetters)
	assert.Equal(t, 90, resp.Data.Attributes.SentLetters)
	assert.Equal(t, 5, resp.Data.Attributes.CancelledLetters)
	assert.Equal(t, 0, resp.Data.Attributes.ErrorLetters)
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
}