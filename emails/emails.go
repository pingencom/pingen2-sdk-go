package emails

import (
	"encoding/json"
	"fmt"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
	"github.com/pingencom/pingen2-sdk-go/fileupload"
)

type Emails struct {
	organisationID string
	apiRequestor   *api.APIRequestor
}

type EmailResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Status              string  `json:"status"`
			FileOriginalName    string  `json:"file_original_name"`
			FilePages           int     `json:"file_pages"`
			RecipientIdentifier string  `json:"recipient_identifier"`
			PriceCurrency       string  `json:"price_currency"`
			PriceValue          float64 `json:"price_value"`
			Source              string  `json:"source"`
			SubmittedAt         string  `json:"submitted_at"`
			CreatedAt           string  `json:"created_at"`
			UpdatedAt           string  `json:"updated_at"`
		} `json:"attributes"`
		Relationships struct {
			Organisation struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"organisation"`
			Events struct {
				Links struct {
					Related struct {
						Href string `json:"href"`
						Meta struct {
							Count int `json:"count"`
						} `json:"meta"`
					} `json:"related"`
				} `json:"links"`
			} `json:"events"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Meta struct {
			Abilities struct {
				Self struct {
					Delete string `json:"delete"`
				} `json:"self"`
			} `json:"abilities"`
		} `json:"meta"`
	} `json:"data"`
	Included []struct{} `json:"included"`
}

type EmailCollectionResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Status              string  `json:"status"`
			FileOriginalName    string  `json:"file_original_name"`
			FilePages           int     `json:"file_pages"`
			RecipientIdentifier string  `json:"recipient_identifier"`
			PriceCurrency       string  `json:"price_currency"`
			PriceValue          float64 `json:"price_value"`
			Source              string  `json:"source"`
			SubmittedAt         string  `json:"submitted_at"`
			CreatedAt           string  `json:"created_at"`
			UpdatedAt           string  `json:"updated_at"`
		} `json:"attributes"`
		Relationships struct {
			Organisation struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"organisation"`
			Events struct {
				Links struct {
					Related struct {
						Href string `json:"href"`
						Meta struct {
							Count int `json:"count"`
						} `json:"meta"`
					} `json:"related"`
				} `json:"links"`
			} `json:"events"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
	Included []struct{} `json:"included"`
	Links    struct {
		First string `json:"first"`
		Last  string `json:"last"`
		Prev  string `json:"prev"`
		Next  string `json:"next"`
		Self  string `json:"self"`
	} `json:"links"`
	Meta struct {
		CurrentPage int `json:"current_page"`
		LastPage    int `json:"last_page"`
		PerPage     int `json:"per_page"`
		From        int `json:"from"`
		To          int `json:"to"`
		Total       int `json:"total"`
	} `json:"meta"`
}

func NewEmails(organisationID string, apiRequestor *api.APIRequestor) *Emails {
	return &Emails{
		organisationID: organisationID,
		apiRequestor:   apiRequestor,
	}
}

func (e *Emails) GetDetails(emailID string, params map[string]string, suppliedHeaders map[string]string) (EmailResponse, *errors.PingenError) {
	var response EmailResponse
	url := fmt.Sprintf("/organisations/%s/deliveries/emails/%s", e.organisationID, emailID)
	_, err := e.apiRequestor.PerformGetRequest(url, &response, params, suppliedHeaders)
	if err != nil {
		return EmailResponse{}, err
	}

	return response, nil
}

func (e *Emails) GetCollection(params map[string]string, suppliedHeaders map[string]string) (EmailCollectionResponse, *errors.PingenError) {
	var response EmailCollectionResponse
	url := fmt.Sprintf("/organisations/%s/deliveries/emails", e.organisationID)

	_, err := e.apiRequestor.PerformGetRequest(url, &response, params, suppliedHeaders)
	if err != nil {
		return EmailCollectionResponse{}, err
	}

	return response, nil
}

func (e *Emails) UploadAndCreate(
	pathToFile, fileOriginalName string,
	autoSend bool,
	metaData, relationships map[string]interface{},
) (EmailResponse, *errors.PingenError) {
	fileUpload := fileupload.NewFileUpload(e.apiRequestor)

	fileResponse, err := fileUpload.RequestFileUpload()
	if err != nil {
		return EmailResponse{}, err
	}

	err = fileUpload.PutFile(pathToFile, fileResponse.Data.Attributes.URL)
	if err != nil {
		return EmailResponse{}, err
	}

	return e.Create(
		fileResponse.Data.Attributes.URL,
		fileResponse.Data.Attributes.URLSignature,
		fileOriginalName,
		autoSend,
		metaData,
		relationships,
	)
}

func (e *Emails) Create(
	fileURL, fileSignature, fileOriginalName string,
	autoSend bool,
	metaData, relationships map[string]interface{},
) (EmailResponse, *errors.PingenError) {
	attributes := map[string]interface{}{
		"file_original_name": fileOriginalName,
		"file_url":           fileURL,
		"file_url_signature": fileSignature,
		"auto_send":          autoSend,
		"meta_data":          metaData,
	}

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       "emails",
			"attributes": attributes,
		},
	}

	if relationships != nil {
		dataMap := payload["data"].(map[string]interface{})
		dataMap["relationships"] = relationships
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/deliveries/emails", e.organisationID)

	var response EmailResponse

	_, err := e.apiRequestor.PerformPostRequest(url, &response, data, nil)
	if err != nil {
		return EmailResponse{}, err
	}

	return response, nil
}
