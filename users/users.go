package users

import (
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
)

type Users struct {
	apiRequestor *api.APIRequestor
}

type UserResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Email     string   `json:"email"`
			FirstName string   `json:"first_name"`
			LastName  string   `json:"last_name"`
			Status    string   `json:"status"`
			Language  string   `json:"language"`
			Edition   string   `json:"edition"`
			Flags     []string `json:"flags"`
			CreatedAt string   `json:"created_at"`
			UpdatedAt string   `json:"updated_at"`
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
			Notifications struct {
				Links struct {
					Related struct {
						Href string `json:"href"`
						Meta struct {
							Count int `json:"count"`
						} `json:"meta"`
					} `json:"related"`
				} `json:"links"`
			} `json:"notifications"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Meta struct {
			Abilities struct {
				Self struct {
					Reach            string `json:"reach"`
					Act              string `json:"act"`
					ResendActivation string `json:"resend-activation"`
				} `json:"self"`
			} `json:"abilities"`
		} `json:"meta"`
	} `json:"data"`
	Included []struct{} `json:"included"`
}

func NewUsers(apiRequestor *api.APIRequestor) *Users {
	return &Users{
		apiRequestor: apiRequestor,
	}
}

func (u *Users) GetDetails(
	params map[string]string,
	suppliedHeaders map[string]string,
) (UserResponse, *errors.PingenError) {
	var response UserResponse
	_, err := u.apiRequestor.PerformGetRequest("/user", &response, params, suppliedHeaders)
	if err != nil {
		return UserResponse{}, err
	}

	return response, nil
}
