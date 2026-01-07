package batchevents

import (
	"fmt"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
)

type BatchEvents struct {
	organisationID string
	apiRequestor   *api.APIRequestor
}

type BatchEventsCollectionResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Code      string   `json:"code"`
			Name      string   `json:"name"`
			Producer  string   `json:"producer"`
			Location  string   `json:"location"`
			Data      []string `json:"data"`
			EmittedAt string   `json:"emitted_at"`
			CreatedAt string   `json:"created_at"`
			UpdatedAt string   `json:"updated_at"`
		} `json:"attributes"`
		Relationships struct {
			Batch struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"letter"`
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

func NewBatchEvents(organisationID string, apiRequestor *api.APIRequestor) *BatchEvents {
	return &BatchEvents{
		organisationID: organisationID,
		apiRequestor:   apiRequestor,
	}
}

func (be *BatchEvents) fetchCollection(
	url string,
	params map[string]string,
	headers map[string]string,
) (BatchEventsCollectionResponse, *errors.PingenError) {
	var response BatchEventsCollectionResponse

	_, err := be.apiRequestor.PerformGetRequest(url, &response, params, headers)

	if err != nil {
		return BatchEventsCollectionResponse{}, err
	}

	return response, nil
}

func (be *BatchEvents) GetCollection(
	batchID string,
	params map[string]string,
	headers map[string]string,
) (BatchEventsCollectionResponse, *errors.PingenError) {
	requestURL := fmt.Sprintf("/organisations/%s/batches/%s/events", be.organisationID, batchID)

	return be.fetchCollection(requestURL, params, headers)
}
