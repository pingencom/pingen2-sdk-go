package organisations

import (
	"fmt"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
)

type Organisations struct {
	apiRequestor *api.APIRequestor
}

type OrganisationResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name                      string   `json:"name"`
			Status                    string   `json:"status"`
			Plan                      string   `json:"plan"`
			BillingMode               string   `json:"billing_mode"`
			BillingCurrency           string   `json:"billing_currency"`
			BillingBalance            float64  `json:"billing_balance"`
			MissingCredits            int      `json:"missing_credits"`
			Edition                   string   `json:"edition"`
			DefaultCountry            string   `json:"default_country"`
			DefaultAddressPosition    string   `json:"default_address_position"`
			DataRetentionAddresses    int      `json:"data_retention_addresses"`
			DataRetentionPDF          int      `json:"data_retention_pdf"`
			LimitsMonthlyLettersCount int      `json:"limits_monthly_letters_count"`
			Color                     string   `json:"color"`
			Flags                     []string `json:"flags"`
			CreatedAt                 string   `json:"created_at"`
			UpdatedAt                 string   `json:"updated_at"`
		} `json:"attributes"`
		Relationships struct {
			Associations struct {
				Links struct {
					Related struct {
						Href string `json:"href"`
						Meta struct {
							Count int `json:"count"`
						} `json:"meta"`
					} `json:"related"`
				} `json:"links"`
			} `json:"associations"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Meta struct {
			Abilities struct {
				Self struct {
					Manage string `json:"manage"`
				} `json:"self"`
			} `json:"abilities"`
		} `json:"meta"`
	} `json:"data"`
	Included []struct{} `json:"included"`
}

type OrganisationCollectionResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name                      string   `json:"name"`
			Status                    string   `json:"status"`
			Plan                      string   `json:"plan"`
			BillingMode               string   `json:"billing_mode"`
			BillingCurrency           string   `json:"billing_currency"`
			BillingBalance            float64  `json:"billing_balance"`
			MissingCredits            int      `json:"missing_credits"`
			Edition                   string   `json:"edition"`
			DefaultCountry            string   `json:"default_country"`
			DefaultAddressPosition    string   `json:"default_address_position"`
			DataRetentionAddresses    int      `json:"data_retention_addresses"`
			DataRetentionPDF          int      `json:"data_retention_pdf"`
			LimitsMonthlyLettersCount int      `json:"limits_monthly_letters_count"`
			Color                     string   `json:"color"`
			Flags                     []string `json:"flags"`
			CreatedAt                 string   `json:"created_at"`
			UpdatedAt                 string   `json:"updated_at"`
		} `json:"attributes"`
		Relationships struct {
			Associations struct {
				Links struct {
					Related struct {
						Href string `json:"href"`
						Meta struct {
							Count int `json:"count"`
						} `json:"meta"`
					} `json:"related"`
				} `json:"links"`
			} `json:"associations"`
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

func NewOrganisations(apiRequestor *api.APIRequestor) *Organisations {
	return &Organisations{
		apiRequestor: apiRequestor,
	}
}

func (o *Organisations) GetDetails(
	organisationID string,
	params map[string]string,
	suppliedHeaders map[string]string,
) (OrganisationResponse, *errors.PingenError) {
	var response OrganisationResponse
	endpoint := fmt.Sprintf("/organisations/%s", organisationID)
	_, err := o.apiRequestor.PerformGetRequest(endpoint, &response, params, suppliedHeaders)
	if err != nil {
		return OrganisationResponse{}, err
	}

	return response, nil
}

func (o *Organisations) GetCollection(
	params map[string]string,
	suppliedHeaders map[string]string,
) (OrganisationCollectionResponse, *errors.PingenError) {
	var response OrganisationCollectionResponse
	_, err := o.apiRequestor.PerformGetRequest("/organisations", &response, params, suppliedHeaders)
	if err != nil {
		return OrganisationCollectionResponse{}, err
	}

	return response, nil
}
