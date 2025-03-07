package fileupload

import (
	"net/http"
	"os"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
)

type FileUpload struct {
	APIRequestor *api.APIRequestor
}

type FileResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			URL          string `json:"url"`
			URLSignature string `json:"url_signature"`
			ExpiresAt    string `json:"expires_at"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

func NewFileUpload(apiRequestor *api.APIRequestor) *FileUpload {
	return &FileUpload{
		APIRequestor: apiRequestor,
	}
}

func (f *FileUpload) RequestFileUpload() (FileResponse, *errors.PingenError) {
	var response FileResponse

	_, err := f.APIRequestor.PerformGetRequest("/file-upload", &response, nil, nil)
	if err != nil {
		return FileResponse{}, err
	}

	return response, nil
}

func (f *FileUpload) PutFile(pathToFile, fileURL string) *errors.PingenError {
	file, err := os.Open(pathToFile)
	if err != nil {
		return errors.NewPingenError(
			"Failed to open file",
			"",
			http.StatusInternalServerError,
			nil,
		)
	}
	defer file.Close()

	return f.APIRequestor.PerformPutRequest(fileURL, file)
}
