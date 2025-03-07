package letters

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
	"github.com/pingencom/pingen2-sdk-go/fileupload"
)

type Letters struct {
	organisationID string
	apiRequestor   *api.APIRequestor
}

type LetterResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Status           string   `json:"status"`
			FileOriginalName string   `json:"file_original_name"`
			FilePages        int      `json:"file_pages"`
			Address          string   `json:"address"`
			AddressPosition  string   `json:"address_position"`
			Country          string   `json:"country"`
			DeliveryProduct  string   `json:"delivery_product"`
			PrintMode        string   `json:"print_mode"`
			PrintSpectrum    string   `json:"print_spectrum"`
			PriceCurrency    string   `json:"price_currency"`
			PriceValue       float64  `json:"price_value"`
			PaperTypes       []string `json:"paper_types"`
			Fonts            []struct {
				Name       string `json:"name"`
				IsEmbedded bool   `json:"is_embedded"`
			} `json:"fonts"`
			Source         string `json:"source"`
			TrackingNumber string `json:"tracking_number"`
			SubmittedAt    string `json:"submitted_at"`
			CreatedAt      string `json:"created_at"`
			UpdatedAt      string `json:"updated_at"`
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
			Events struct {
				Links struct {
					Related struct {
						Href string `json:"href"`
						Meta struct {
							Count int `json:"count"`
						} `json:"meta"`
					} `json:"related"`
				} `json:"links"`
			} `json:"events"`
			Batch struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"batch"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Meta struct {
			Abilities struct {
				Self struct {
					Cancel                      string `json:"cancel"`
					Delete                      string `json:"delete"`
					Submit                      string `json:"submit"`
					SendSimplex                 string `json:"send-simplex"`
					Edit                        string `json:"edit"`
					GetPdfRaw                   string `json:"get-pdf-raw"`
					GetPdfValidation            string `json:"get-pdf-validation"`
					RestoreOriginal             string `json:"restore-original"`
					ChangePaperType             string `json:"change-paper-type"`
					ChangeWindowPosition        string `json:"change-window-position"`
					CreateCoverpage             string `json:"create-coverpage"`
					AddAttachment               string `json:"add-attachment"`
					FixOverwriteRestrictedAreas string `json:"fix-overwrite-restricted-areas"`
					FixCoverPage                string `json:"fix-coverpage"`
					FixCountry                  string `json:"fix-country"`
					FixRegularPaper             string `json:"fix-regular-paper"`
					FixAddress                  string `json:"fix-address"`
					FixInteractiveContent       string `json:"fix-interactive-content"`
					FixFormat                   string `json:"fix-format"`
					ApplyPreset                 string `json:"apply-preset"`
					CreatePreset                string `json:"create-preset"`
				} `json:"self"`
			} `json:"abilities"`
		} `json:"meta"`
	} `json:"data"`
	Included []struct{} `json:"included"`
}

type LetterCollectionResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Status           string   `json:"status"`
			FileOriginalName string   `json:"file_original_name"`
			FilePages        int      `json:"file_pages"`
			Address          string   `json:"address"`
			AddressPosition  string   `json:"address_position"`
			Country          string   `json:"country"`
			DeliveryProduct  string   `json:"delivery_product"`
			PrintMode        string   `json:"print_mode"`
			PrintSpectrum    string   `json:"print_spectrum"`
			PriceCurrency    string   `json:"price_currency"`
			PriceValue       float64  `json:"price_value"`
			PaperTypes       []string `json:"paper_types"`
			Fonts            []struct {
				Name       string `json:"name"`
				IsEmbedded bool   `json:"is_embedded"`
			} `json:"fonts"`
			Source         string `json:"source"`
			TrackingNumber string `json:"tracking_number"`
			SubmittedAt    string `json:"submitted_at"`
			CreatedAt      string `json:"created_at"`
			UpdatedAt      string `json:"updated_at"`
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
			Events struct {
				Links struct {
					Related struct {
						Href string `json:"href"`
						Meta struct {
							Count int `json:"count"`
						} `json:"meta"`
					} `json:"related"`
				} `json:"links"`
			} `json:"events"`
			Batch struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"batch"`
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

type PriceCalculationResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Currency string  `json:"currency"`
			Price    float64 `json:"price"`
		} `json:"attributes"`
	} `json:"data"`
}

func NewLetters(organisationID string, apiRequestor *api.APIRequestor) *Letters {
	return &Letters{
		organisationID: organisationID,
		apiRequestor:   apiRequestor,
	}
}

func (l *Letters) GetDetails(letterID string, params map[string]string, suppliedHeaders map[string]string) (LetterResponse, *errors.PingenError) {
	var response LetterResponse
	url := fmt.Sprintf("/organisations/%s/letters/%s", l.organisationID, letterID)
	_, err := l.apiRequestor.PerformGetRequest(url, &response, params, suppliedHeaders)
	if err != nil {
		return LetterResponse{}, err
	}

	return response, nil
}

func (l *Letters) GetCollection(params map[string]string, suppliedHeaders map[string]string) (LetterCollectionResponse, *errors.PingenError) {
	var response LetterCollectionResponse
	url := fmt.Sprintf("/organisations/%s/letters", l.organisationID)

	_, err := l.apiRequestor.PerformGetRequest(url, &response, params, suppliedHeaders)
	if err != nil {
		return LetterCollectionResponse{}, err
	}

	return response, nil
}

func (l *Letters) UploadAndCreate(
	pathToFile, fileOriginalName, addressPosition string,
	autoSend bool,
	deliveryProduct, printMode, printSpectrum, senderAddress string,
	metaData map[string]interface{},
) (LetterResponse, *errors.PingenError) {
	fileUpload := fileupload.NewFileUpload(l.apiRequestor)

	fileResponse, err := fileUpload.RequestFileUpload()
	if err != nil {
		return LetterResponse{}, err
	}

	err = fileUpload.PutFile(pathToFile, fileResponse.Data.Attributes.URL)
	if err != nil {
		return LetterResponse{}, err
	}

	return l.Create(
		fileResponse.Data.Attributes.URL,
		fileResponse.Data.Attributes.URLSignature,
		fileOriginalName,
		addressPosition,
		autoSend,
		deliveryProduct,
		printMode,
		printSpectrum,
		senderAddress,
		metaData,
	)
}

func (l *Letters) Create(
	fileURL, fileSignature, fileOriginalName, addressPosition string,
	autoSend bool,
	deliveryProduct, printMode, printSpectrum, senderAddress string,
	metaData map[string]interface{},
) (LetterResponse, *errors.PingenError) {
	attributes := map[string]interface{}{
		"file_original_name": fileOriginalName,
		"file_url":           fileURL,
		"file_url_signature": fileSignature,
		"address_position":   addressPosition,
		"auto_send":          autoSend,
	}

	if deliveryProduct != "" {
		attributes["delivery_product"] = deliveryProduct
	}

	if printMode != "" {
		attributes["print_mode"] = printMode
	}

	if printSpectrum != "" {
		attributes["print_spectrum"] = printSpectrum
	}

	if senderAddress != "" {
		attributes["sender_address"] = senderAddress
	}

	if metaData != nil {
		attributes["meta_data"] = metaData
	}

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       "letters",
			"attributes": attributes,
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/letters", l.organisationID)

	var response LetterResponse

	_, err := l.apiRequestor.PerformPostRequest(url, &response, data, nil)
	if err != nil {
		return LetterResponse{}, err
	}

	return response, nil
}

func (l *Letters) Send(letterID, deliveryProduct, printMode, printSpectrum string) (LetterResponse, *errors.PingenError) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":   letterID,
			"type": "letters",
			"attributes": map[string]string{
				"delivery_product": deliveryProduct,
				"print_mode":       printMode,
				"print_spectrum":   printSpectrum,
			},
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/letters/%s/send", l.organisationID, letterID)

	var response LetterResponse

	_, err := l.apiRequestor.PerformPatchRequest(url, &response, data, nil)
	if err != nil {
		return LetterResponse{}, err
	}

	return response, nil

}

func (l *Letters) Cancel(letterID string) (interface{}, *errors.PingenError) {
	url := fmt.Sprintf("/organisations/%s/letters/%s/cancel", l.organisationID, letterID)
	return l.apiRequestor.PerformCancelRequest(url)
}

func (l *Letters) Delete(letterID string) (interface{}, *errors.PingenError) {
	url := fmt.Sprintf("/organisations/%s/letters/%s", l.organisationID, letterID)
	return l.apiRequestor.PerformDeleteRequest(url)
}

func (l *Letters) Edit(letterID string, paperTypes []string) (LetterResponse, *errors.PingenError) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":   letterID,
			"type": "letters",
			"attributes": map[string]interface{}{
				"paper_types": paperTypes,
			},
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/letters/%s", l.organisationID, letterID)

	var response LetterResponse

	_, apiErr := l.apiRequestor.PerformPatchRequest(url, &response, data, nil)
	if apiErr != nil {
		return LetterResponse{}, apiErr
	}

	return response, nil
}

func (l *Letters) GetFile(letterID string) (io.ReadCloser, *errors.PingenError) {
	url := fmt.Sprintf("/organisations/%s/letters/%s/file", l.organisationID, letterID)
	return l.apiRequestor.PerformStreamRequest(url)
}

func (l *Letters) CalculatePrice(
	country string,
	paperTypes []string,
	printMode, printSpectrum, deliveryProduct string,
) (PriceCalculationResponse, *errors.PingenError) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "letter_price_calculator",
			"attributes": map[string]interface{}{
				"country":          country,
				"paper_types":      paperTypes,
				"print_mode":       printMode,
				"print_spectrum":   printSpectrum,
				"delivery_product": deliveryProduct,
			},
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/letters/price-calculator", l.organisationID)

	var response PriceCalculationResponse

	_, apiErr := l.apiRequestor.PerformPostRequest(url, &response, data, nil)
	if apiErr != nil {
		return PriceCalculationResponse{}, apiErr
	}

	return response, nil
}
