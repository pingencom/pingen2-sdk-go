package letters_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/letters"
	"github.com/stretchr/testify/assert"
)

const mockResponse = `{
	"data": {
		"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
		"type": "letters",
		"attributes": {
			"status": "send",
			"file_original_name": "lorem.pdf",
			"file_pages": 2,
			"address": "Hans Meier\nExample street 4\n8000 Zürich\nSwitzerland",
			"address_position": "left",
			"country": "CH",
			"delivery_product": "fast",
			"print_mode": "simplex",
			"print_spectrum": "color",
			"price_currency": "CHF",
			"price_value": 1.25,
			"paper_types": ["normal", "qr"],
			"fonts": [
				{
					"name": "Helvetica",
					"is_embedded": true
				},
				{
					"name": "Helvetica-Bold",
					"is_embedded": false
				}
			],
			"source": "api",
			"tracking_number": "98.1234.11",
			"sender_address": "ACME GmbH | Strasse 3 | 8000 Zürich",
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
			},
			"batch": {
				"links": {
					"related": "string"
				},
				"data": {
					"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					"type": "batches"
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
					"send-simplex": "state",
					"edit": "state",
					"get-pdf-raw": "ok",
					"get-pdf-validation": "ok",
					"restore-original": "state",
					"change-paper-type": "state",
					"change-window-position": "state",
					"create-coverpage": "state",
					"add-attachment": "state",
					"fix-overwrite-restricted-areas": "state",
					"fix-coverpage": "state",
					"fix-country": "state",
					"fix-regular-paper": "state",
					"fix-address": "state",
					"fix-interactive-content": "state",
					"fix-format": "state",
					"apply-preset": "state",
					"create-preset": "state"
				}
			}
		}
	},
	"included": [{}]
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

func setupLetter(apiBaseURL string) *letters.Letters {
	config, _ := pingen2sdk.InitSDK("testSetClientId", "testSetClientSecret", "")
	config.SetAPIBaseURL(apiBaseURL)
	apiRequestor := api.NewAPIRequestor("dummyToken", config)

	return letters.NewLetters("testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", apiRequestor)
}

func TestGetDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := letterClient.GetDetails("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", params, headers)

	assert.Nil(t, err)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", resp.Data.ID)
	assert.Equal(t, "lorem.pdf", resp.Data.Attributes.FileOriginalName)
	assert.Equal(t, "CH", resp.Data.Attributes.Country)
	assert.Equal(t, "2021-11-19T09:42:48+0100", resp.Data.Attributes.SubmittedAt)
}

func TestGetDetails_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	letterClient := setupLetter(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := letterClient.GetDetails("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestGetCollection(t *testing.T) {
	mockCollectionResponse := `{
		"data": [
			{
				"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				"type": "letters",
				"attributes": {
					"status": "string",
					"file_original_name": "lorem.pdf",
					"file_pages": 2,
					"address": "Hans Meier\nExample street 4\n8000 Zürich\nSwitzerland",
					"address_position": "left",
					"country": "CH",
					"delivery_product": "fast",
					"print_mode": "simplex",
					"print_spectrum": "color",
					"price_currency": "CHF",
					"price_value": 1.25,
					"paper_types": ["normal", "qr"],
					"fonts": [
						{"name": "Helvetica", "is_embedded": true},
						{"name": "Helvetica-Bold", "is_embedded": false}
					],
					"source": "api",
					"tracking_number": "98.1234.11",
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
                    },
                    "batch": {
                        "links": {
                            "related": "string"
                        },
                        "data": {
                            "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
                            "type": "batches"
                        }
                    }
                },
				"links": { "self": "string" }
			}
		],
		"included": [{}],
		"links": { "first": "first-link", "last": "last-link" },
		"meta": { "current_page": 1, "last_page": 1, "per_page": 10, "total": 1 }
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockCollectionResponse))
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	resp, err := letterClient.GetCollection(params, headers)

	assert.Nil(t, err)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", resp.Data[0].ID)
	assert.Equal(t, "lorem.pdf", resp.Data[0].Attributes.FileOriginalName)
	assert.Equal(t, 2, resp.Data[0].Attributes.FilePages)
	assert.Equal(t, 1, resp.Meta.CurrentPage)
	assert.Equal(t, 10, resp.Meta.PerPage)
	assert.Equal(t, 1, resp.Meta.Total)
}

func TestGetCollection_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()
	letterClient := setupLetter(server.URL)

	params := map[string]string{}
	headers := map[string]string{}
	_, err := letterClient.GetCollection(params, headers)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestUploadAndCreate(t *testing.T) {
	metaData := map[string]interface{}{
		"recipient": map[string]string{
			"name":    "R_Example",
			"street":  "R_Street",
			"number":  "R_12",
			"zip":     "R_12",
			"city":    "R_Warsaw",
			"country": "PL",
		},
		"sender": map[string]string{
			"name":    "S_Example",
			"street":  "S_Street",
			"number":  "S_12",
			"zip":     "S_12",
			"city":    "S_Warsaw",
			"country": "PL",
		},
	}

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
			assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			body, _ := io.ReadAll(r.Body)
			assert.Contains(t, string(body), `"file_original_name":"uploaded-test.pdf"`)

			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(mockResponse))
		}

		counter++
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	filePath := "testFile.pdf"

	resp, err := letterClient.UploadAndCreate(
		filePath,
		"uploaded-test.pdf",
		"left",
		true,
		"fast",
		"simplex",
		"color",
		"Example Sender Address",
		metaData,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", resp.Data.ID)
}

func TestUploadAndCreate_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	letterClient := setupLetter(server.URL)
	filePath := "testFile.pdf"

	_, err := letterClient.UploadAndCreate(
		filePath,
		"uploaded-test.pdf",
		"left",
		true,
		"fast",
		"simplex",
		"color",
		"Example Sender Address",
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestUploadAndCreate_ErrorInPut(t *testing.T) {
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

	letterClient := setupLetter(server.URL)
	filePath := "testFile.pdf"

	_, err := letterClient.UploadAndCreate(
		filePath,
		"uploaded-test.pdf",
		"left",
		true,
		"fast",
		"simplex",
		"color",
		"Example Sender Address",
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: Api error (Status Code: 401, Request ID: )"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), `"file_original_name":"uploaded-test.pdf"`)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	filePath := "testFile.pdf"

	resp, err := letterClient.Create(
		filePath,
		"test-signature",
		"uploaded-test.pdf",
		"left",
		false,
		"",
		"",
		"",
		"",
		nil,
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", resp.Data.ID)
}

func TestCreate_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	letterClient := setupLetter(server.URL)

	_, err := letterClient.Create(
		"https://s3.example.com/file/test",
		"signature-123",
		"test.pdf",
		"left",
		true,
		"fast",
		"simplex",
		"color",
		"Example Street",
		nil,
	)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestSendLetter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/send", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), `"delivery_product":"fast"`)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	resp, err := letterClient.Send("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", "fast", "simplex", "color")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", resp.Data.ID)
}

func TestSendLetter_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	letterClient := setupLetter(server.URL)

	_, err := letterClient.Send("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", "fast", "simplex", "color")

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestCancelLetter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/test-letter-id/cancel", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	resp, err := letterClient.Cancel("test-letter-id")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteLetter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/test-letter-id", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	resp, err := letterClient.Delete("test-letter-id")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestEdit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		expectedPayload := `{
            "data": {
                "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1",
                "type": "letters",
                "attributes": {
                    "paper_types": ["normal", "qr"]
                }
            }
        }`
		assert.JSONEq(t, expectedPayload, string(body))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	paperTypes := []string{"normal", "qr"}
	resp, err := letterClient.Edit("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", paperTypes)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", resp.Data.ID)
}

func TestEdit_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	letterClient := setupLetter(server.URL)

	paperTypes := []string{"normal", "qr"}
	_, err := letterClient.Edit("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1", paperTypes)

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}

func TestGetFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/letterxx-xxxx-xxxx-xxxx-xxxxxxxxxx21/file", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "mock file content")
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	stream, err := letterClient.GetFile("letterxx-xxxx-xxxx-xxxx-xxxxxxxxxx21")

	assert.Nil(t, err)
	assert.NotNil(t, stream)

	defer stream.Close()
	responseData, readErr := io.ReadAll(stream)
	assert.Nil(t, readErr)
	assert.Equal(t, "mock file content", string(responseData))
}

func TestCalculatePrice(t *testing.T) {
	mockPriceResponse := `{
		"data": {
			"id": "price-id",
			"type": "letter_price_calculator",
			"attributes": {
				"currency": "CHF",
				"price": 12.34
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/organisations/testxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx1/letters/price-calculator", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), `"country":"CH"`)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockPriceResponse))
	}))
	defer server.Close()

	letterClient := setupLetter(server.URL)

	resp, err := letterClient.CalculatePrice("CH", []string{"normal", "qr"}, "simplex", "color", "fast")

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "CHF", resp.Data.Attributes.Currency)
	assert.Equal(t, 12.34, resp.Data.Attributes.Price)
}

func TestCalculatePrice_Error(t *testing.T) {
	server := setupUnauthorizedServer()
	defer server.Close()

	letterClient := setupLetter(server.URL)

	_, err := letterClient.CalculatePrice("CH", []string{"normal", "qr"}, "simplex", "color", "fast")

	assert.NotNil(t, err)
	expectedMessage := "PingenError: API error (Status Code: 401, Request ID: requestx-yyyy-yyyy-yyyy-yyyyyyyyyyy2)"
	assert.Equal(t, expectedMessage, err.Error())
}
