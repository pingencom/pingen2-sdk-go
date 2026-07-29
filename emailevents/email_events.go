package emailevents

import (
	"fmt"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
)

type EmailEvents struct {
	organisationID string
	apiRequestor   *api.APIRequestor
}

type EmailEventsCollectionResponse struct {
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
			Email struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"email"`
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

func NewEmailEvents(organisationID string, apiRequestor *api.APIRequestor) *EmailEvents {
	return &EmailEvents{
		organisationID: organisationID,
		apiRequestor:   apiRequestor,
	}
}

func (ee *EmailEvents) GetCollection(
	emailID string,
	params map[string]string,
	headers map[string]string,
) (EmailEventsCollectionResponse, *errors.PingenError) {
	requestURL := fmt.Sprintf("/organisations/%s/deliveries/emails/%s/events", ee.organisationID, emailID)

	var response EmailEventsCollectionResponse

	_, err := ee.apiRequestor.PerformGetRequest(requestURL, &response, params, headers)

	if err != nil {
		return EmailEventsCollectionResponse{}, err
	}

	return response, nil
}
