package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/errors"
	"github.com/pingencom/pingen2-sdk-go/response"
)

type APIRequestor struct {
	accessToken     string
	config          *pingen2sdk.Config
	responseHandler *response.JSONResponseHandler
}

func NewAPIRequestor(accessToken string, config *pingen2sdk.Config) *APIRequestor {
	return &APIRequestor{
		accessToken:     accessToken,
		config:          config,
		responseHandler: &response.JSONResponseHandler{},
	}
}

func (r *APIRequestor) PerformGetRequest(
	url string,
	target interface{},
	params map[string]string,
	extraHeaders map[string]string,
) (interface{}, *errors.PingenError) {
	return r.performHTTPRequest(http.MethodGet, url, nil, extraHeaders, params, target)
}

func (r *APIRequestor) PerformPutRequest(
	url string,
	file io.Reader,
) *errors.PingenError {
	req, _ := http.NewRequest(http.MethodPut, url, file)

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.NewPingenError(
			"Internal error",
			fmt.Sprintf("Failed to send PUT request: %v", err.Error()),
			http.StatusInternalServerError,
			nil,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.NewPingenError(
			"Api error",
			fmt.Sprintf("PUT request failed with status %d", resp.StatusCode),
			resp.StatusCode,
			nil,
		)
	}

	return nil
}

func (r *APIRequestor) PerformPostRequest(
	url string,
	target interface{},
	payload []byte,
	extraHeaders map[string]string,
) (interface{}, *errors.PingenError) {
	body := bytes.NewBuffer(payload)
	return r.performHTTPRequest(http.MethodPost, url, body, extraHeaders, nil, target)
}

func (r *APIRequestor) PerformPatchRequest(
	url string,
	target interface{},
	payload []byte,
	extraHeaders map[string]string,
) (interface{}, *errors.PingenError) {
	body := bytes.NewBuffer(payload)
	return r.performHTTPRequest(http.MethodPatch, url, body, extraHeaders, nil, target)
}

func (r *APIRequestor) PerformCancelRequest(
	urlPath string,
) (interface{}, *errors.PingenError) {
	return r.performHTTPRequest(http.MethodPatch, urlPath, nil, nil, nil, nil)
}

func (r *APIRequestor) PerformDeleteRequest(
	urlPath string,
) (interface{}, *errors.PingenError) {
	return r.performHTTPRequest(http.MethodDelete, urlPath, nil, nil, nil, nil)
}

func (r *APIRequestor) PerformStreamRequest(url string) (io.ReadCloser, *errors.PingenError) {
	reqURL := r.preparePath(url, nil)
	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	req.Header = r.requestHeaders(nil)

	client := &http.Client{Timeout: r.config.GetRequestTimeout()}
	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, errors.NewPingenError(
			"Invalid HTTP response",
			fmt.Sprintf("Stream request failed with status %d", resp.StatusCode),
			resp.StatusCode,
			nil,
		)
	}

	return resp.Body, nil
}

func (r *APIRequestor) performHTTPRequest(
	method string,
	urlPath string,
	body io.Reader,
	headers map[string]string,
	params map[string]string,
	target interface{},
) (interface{}, *errors.PingenError) {
	reqURL := r.preparePath(urlPath, params)

	req, _ := http.NewRequest(method, reqURL, body)
	req.Header = r.requestHeaders(headers)

	client := &http.Client{Timeout: r.config.GetRequestTimeout()}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	return r.responseHandler.InterpretResponse(resp, target)
}

func (r *APIRequestor) preparePath(
	urlPath string,
	params map[string]string,
) string {
	reqURL, _ := url.Parse(r.config.GetAPIBaseURL() + urlPath)
	query := reqURL.Query()

	for key, value := range params {
		query.Add(key, value)
	}

	reqURL.RawQuery = query.Encode()
	return reqURL.String()
}

func (r *APIRequestor) requestHeaders(extraHeaders map[string]string) http.Header {
	headers := http.Header{}

	headers.Add("User-Agent", r.config.GetUserAgent())
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", r.accessToken))
	headers.Add("Content-Type", "application/vnd.api+json")
	headers.Add("Accept", "application/vnd.api+json")

	for key, value := range extraHeaders {
		headers.Add(key, value)
	}

	return headers
}
