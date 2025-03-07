package response

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pingencom/pingen2-sdk-go/errors"
)

type JSONResponseHandler struct{}

type DefaultResponse struct {
	Body       string
	StatusCode int
}

type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Self  string `json:"self"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	Total       int `json:"total"`
}

type BaseListResponse struct {
	Included []struct{} `json:"included"`
	Links    Links      `json:"links"`
	Meta     Meta       `json:"meta"`
}

func (r *JSONResponseHandler) InterpretResponse(resp *http.Response, target interface{}) (interface{}, *errors.PingenError) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewPingenError("Failed to read response body", err.Error(), resp.StatusCode, convertHeaders(resp.Header))
	}

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusAccepted {
		return &DefaultResponse{
			Body:       string(bodyBytes),
			StatusCode: resp.StatusCode,
		}, nil
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return nil, errors.NewPingenError("Failed to parse response body", string(bodyBytes), resp.StatusCode, convertHeaders(resp.Header))
		}
		return target, nil
	}

	return nil, errors.NewPingenError("API error", string(bodyBytes), resp.StatusCode, convertHeaders(resp.Header))
}

func convertHeaders(headers http.Header) map[string]string {
	converted := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			converted[key] = values[0]
		}
	}
	return converted
}
