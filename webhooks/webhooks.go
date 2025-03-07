package webhooks

import (
	"encoding/json"
	"fmt"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
)

type Webhooks struct {
	organisationID string
	apiRequestor   *api.APIRequestor
}

type WebhookResponseData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		EventCategory string `json:"event_category"`
		URL           string `json:"url"`
		SigningKey    string `json:"signing_key"`
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
	} `json:"relationships"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type WebhookResponse struct {
	Data     WebhookResponseData `json:"data"`
	Included []struct{}          `json:"included"`
}

type WebhookCollectionResponse struct {
	Data     []WebhookResponseData `json:"data"`
	Included []struct{}            `json:"included"`
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

func NewWebhooks(organisationID string, apiRequestor *api.APIRequestor) *Webhooks {
	return &Webhooks{
		organisationID: organisationID,
		apiRequestor:   apiRequestor,
	}
}

func (w *Webhooks) GetDetails(webhookID string, params map[string]string, headers map[string]string) (WebhookResponse, *errors.PingenError) {
	requestUrl := fmt.Sprintf("/organisations/%s/webhooks/%s", w.organisationID, webhookID)
	var response WebhookResponse

	_, err := w.apiRequestor.PerformGetRequest(requestUrl, &response, params, headers)
	if err != nil {
		return WebhookResponse{}, err
	}

	return response, nil
}

func (w *Webhooks) GetCollection(params map[string]string, headers map[string]string) (WebhookCollectionResponse, *errors.PingenError) {
	requestUrl := fmt.Sprintf("/organisations/%s/webhooks", w.organisationID)
	var response WebhookCollectionResponse

	_, err := w.apiRequestor.PerformGetRequest(requestUrl, &response, params, headers)
	if err != nil {
		return WebhookCollectionResponse{}, err
	}

	return response, nil
}

func (w *Webhooks) Create(eventCategory string, url string, signingKey string) (WebhookResponse, *errors.PingenError) {
	requestUrl := fmt.Sprintf("/organisations/%s/webhooks", w.organisationID)

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "webhooks",
			"attributes": map[string]string{
				"event_category": eventCategory,
				"url":            url,
				"signing_key":    signingKey,
			},
		},
	}

	data, _ := json.Marshal(payload)
	var response WebhookResponse

	_, err := w.apiRequestor.PerformPostRequest(requestUrl, &response, data, nil)
	if err != nil {
		return WebhookResponse{}, err
	}

	return response, nil
}

func (w *Webhooks) Delete(webhookID string) (interface{}, *errors.PingenError) {
	requestUrl := fmt.Sprintf("/organisations/%s/webhooks/%s", w.organisationID, webhookID)
	return w.apiRequestor.PerformDeleteRequest(requestUrl)
}
