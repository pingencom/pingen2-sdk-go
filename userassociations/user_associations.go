package userassociations

import (
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
	"github.com/pingencom/pingen2-sdk-go/response"
)

type UserAssociations struct {
	apiRequestor *api.APIRequestor
}

type AssociationCollectionResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Role      string `json:"role"`
			Status    string `json:"status"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
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
		Meta struct {
			Abilities struct {
				Self struct {
					Join  string `json:"join"`
					Leave string `json:"leave"`
					Block string `json:"block"`
				} `json:"self"`
				Organisation struct {
					Manage string `json:"manage"`
				} `json:"organisation"`
			} `json:"abilities"`
		} `json:"meta"`
	} `json:"data"`
	response.BaseListResponse
}

func NewUserAssociations(apiRequestor *api.APIRequestor) *UserAssociations {
	return &UserAssociations{
		apiRequestor: apiRequestor,
	}
}

func (ua *UserAssociations) GetCollection(
	params map[string]string,
	suppliedHeaders map[string]string,
) (AssociationCollectionResponse, *errors.PingenError) {
	var response AssociationCollectionResponse
	_, err := ua.apiRequestor.PerformGetRequest("/user/associations", &response, params, suppliedHeaders)
	if err != nil {
		return AssociationCollectionResponse{}, err
	}

	return response, nil
}
