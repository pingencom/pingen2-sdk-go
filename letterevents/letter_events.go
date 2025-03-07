package letterevents

import (
	"fmt"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
)

type LetterEvents struct {
	organisationID string
	apiRequestor   *api.APIRequestor
}

type LetterEventsCollectionResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Code      string   `json:"code"`
			Name      string   `json:"name"`
			Producer  string   `json:"producer"`
			Location  string   `json:"location"`
			HasImage  bool     `json:"has_image"`
			Data      []string `json:"data"`
			EmittedAt string   `json:"emitted_at"`
			CreatedAt string   `json:"created_at"`
			UpdatedAt string   `json:"updated_at"`
		} `json:"attributes"`
		Relationships struct {
			Letter struct {
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

func NewLetterEvents(organisationID string, apiRequestor *api.APIRequestor) *LetterEvents {
	return &LetterEvents{
		organisationID: organisationID,
		apiRequestor:   apiRequestor,
	}
}

func (le *LetterEvents) fetchCollection(
	url string,
	params map[string]string,
	headers map[string]string,
) (LetterEventsCollectionResponse, *errors.PingenError) {
	var response LetterEventsCollectionResponse

	_, err := le.apiRequestor.PerformGetRequest(url, &response, params, headers)

	if err != nil {
		return LetterEventsCollectionResponse{}, err
	}

	return response, nil
}

func (le *LetterEvents) GetCollection(
	letterID string,
	params map[string]string,
	headers map[string]string,
) (LetterEventsCollectionResponse, *errors.PingenError) {
	requestURL := fmt.Sprintf("/organisations/%s/letters/%s/events", le.organisationID, letterID)

	return le.fetchCollection(requestURL, params, headers)
}

func (le *LetterEvents) GetIssueCollection(
	params map[string]string,
	headers map[string]string,
) (LetterEventsCollectionResponse, *errors.PingenError) {
	requestURL := fmt.Sprintf("/organisations/%s/letters/events/issues", le.organisationID)

	return le.fetchCollection(requestURL, params, headers)
}

func (le *LetterEvents) GetUndeliverableCollection(
	params map[string]string,
	headers map[string]string,
) (LetterEventsCollectionResponse, *errors.PingenError) {
	requestURL := fmt.Sprintf("/organisations/%s/letters/events/undeliverable", le.organisationID)

	return le.fetchCollection(requestURL, params, headers)
}

func (le *LetterEvents) GetSentCollection(
	params map[string]string,
	headers map[string]string,
) (LetterEventsCollectionResponse, *errors.PingenError) {
	requestURL := fmt.Sprintf("/organisations/%s/letters/events/sent", le.organisationID)

	return le.fetchCollection(requestURL, params, headers)
}
