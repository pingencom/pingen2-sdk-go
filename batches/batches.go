package batches

import (
	"encoding/json"
	"fmt"

	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/errors"
	"github.com/pingencom/pingen2-sdk-go/fileupload"
)

type Icon string

const (
	IconCampaign     Icon = "campaign"
	IconMegaphone    Icon = "megaphone"
	IconWaveHand     Icon = "wave-hand"
	IconFlash        Icon = "flash"
	IconRocket       Icon = "rocket"
	IconBell         Icon = "bell"
	IconPercentTag   Icon = "percent-tag"
	IconPercentBadge Icon = "percent-badge"
	IconPresent      Icon = "present"
	IconReceipt      Icon = "receipt"
	IconDocument     Icon = "document"
	IconInformation  Icon = "information"
	IconCalendar     Icon = "calendar"
	IconNewspaper    Icon = "newspaper"
	IconCrown        Icon = "crown"
	IconVirus        Icon = "virus"
)

type ChannelType string

const (
	ChannelTypePost  ChannelType = "post"
	ChannelTypeEbill ChannelType = "ebill"
	ChannelTypeEmail ChannelType = "email"
)

type DeliveryProduct string

const (
	DeliveryProductFast       DeliveryProduct = "fast"
	DeliveryProductCheap      DeliveryProduct = "cheap"
	DeliveryProductBulk       DeliveryProduct = "bulk"
	DeliveryProductPremium    DeliveryProduct = "premium"
	DeliveryProductRegistered DeliveryProduct = "registered"

	DeliveryProductElectronicEmail DeliveryProduct = "electronic_email"
	DeliveryProductElectronicEbill DeliveryProduct = "electronic_ebill"
)

type PrintMode string

const (
	PrintModeSimplex PrintMode = "simplex"
	PrintModeDuplex  PrintMode = "duplex"
)

type PrintSpectrum string

const (
	PrintSpectrumColor     PrintSpectrum = "color"
	PrintSpectrumGrayscale PrintSpectrum = "grayscale"
)

type AddressPosition string

const (
	AddressPositionLeft  AddressPosition = "left"
	AddressPositionRight AddressPosition = "right"
)

type GroupingType string

const (
	GroupingTypeZip   GroupingType = "zip"
	GroupingTypeMerge GroupingType = "merge"
)

type SplitType string

const (
	SplitTypeFile      SplitType = "file"
	SplitTypePage      SplitType = "page"
	SplitTypeCustom    SplitType = "custom"
	SplitTypeQRInvoice SplitType = "qr_invoice"
)

type SplitPosition string

const (
	SplitPositionFirstPage SplitPosition = "first_page"
	SplitPositionLastPage  SplitPosition = "last_page"
)

type Batches struct {
	organisationID string
	apiRequestor   *api.APIRequestor
}

type BatchResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name        string `json:"name"`
			ChannelType string `json:"channel_type"`
			Icon        string `json:"icon"`
			Status      string `json:"status"`

			FileOriginalName string `json:"file_original_name"`
			// LetterCount is deprecated, use DeliverableCount instead.
			LetterCount      int     `json:"letter_count"`
			DeliverableCount int     `json:"deliverable_count"`
			AddressPosition  string  `json:"address_position"`
			PrintMode        string  `json:"print_mode"`
			PrintSpectrum    string  `json:"print_spectrum"`
			PriceCurrency    string  `json:"price_currency"`
			PriceValue       float64 `json:"price_value"`
			Source           string  `json:"source"`
			SubmittedAt      string  `json:"submitted_at"`
			CreatedAt        string  `json:"created_at"`
			UpdatedAt        string  `json:"updated_at"`
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
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Meta struct {
			Abilities struct {
				Self struct {
					Cancel               string `json:"cancel"`
					Delete               string `json:"delete"`
					Submit               string `json:"submit"`
					Edit                 string `json:"edit"`
					ChangeWindowPosition string `json:"change-window-position"`
					AddAttachment        string `json:"add-attachment"`
				} `json:"self"`
			} `json:"abilities"`
		} `json:"meta"`
	} `json:"data"`
	Included []struct{} `json:"included"`
}

type BatchCollectionResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name        string `json:"name"`
			ChannelType string `json:"channel_type"`
			Icon        string `json:"icon"`
			Status      string `json:"status"`

			FileOriginalName string `json:"file_original_name"`
			// LetterCount is deprecated, use DeliverableCount instead.
			LetterCount      int     `json:"letter_count"`
			DeliverableCount int     `json:"deliverable_count"`
			AddressPosition  string  `json:"address_position"`
			PrintMode        string  `json:"print_mode"`
			PrintSpectrum    string  `json:"print_spectrum"`
			PriceCurrency    string  `json:"price_currency"`
			PriceValue       float64 `json:"price_value"`
			Source           string  `json:"source"`
			SubmittedAt      string  `json:"submitted_at"`
			CreatedAt        string  `json:"created_at"`
			UpdatedAt        string  `json:"updated_at"`
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

type BatchStatisticsResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			LetterValidating int `json:"letter_validating"`
			LetterGroups     []struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"letter_groups"`
			LetterCountries []struct {
				Country string `json:"country"`
				Count   int    `json:"count"`
			} `json:"letter_countries"`
			LetterRegions []struct {
				Country string `json:"country"`
				Count   int    `json:"count"`
			} `json:"letter_regions"`
		} `json:"attributes"`
	} `json:"data"`
}

func NewBatches(organisationID string, apiRequestor *api.APIRequestor) *Batches {
	return &Batches{
		organisationID: organisationID,
		apiRequestor:   apiRequestor,
	}
}

func (b *Batches) GetDetails(batchID string, params map[string]string, suppliedHeaders map[string]string) (BatchResponse, *errors.PingenError) {
	var response BatchResponse
	url := fmt.Sprintf("/organisations/%s/batches/%s", b.organisationID, batchID)
	_, err := b.apiRequestor.PerformGetRequest(url, &response, params, suppliedHeaders)
	if err != nil {
		return BatchResponse{}, err
	}

	return response, nil
}

func (b *Batches) GetCollection(params map[string]string, suppliedHeaders map[string]string) (BatchCollectionResponse, *errors.PingenError) {
	var response BatchCollectionResponse
	url := fmt.Sprintf("/organisations/%s/batches", b.organisationID)

	_, err := b.apiRequestor.PerformGetRequest(url, &response, params, suppliedHeaders)
	if err != nil {
		return BatchCollectionResponse{}, err
	}

	return response, nil
}

func (b *Batches) UploadAndCreateBatch(
	pathToFile, name string, icon Icon, channelType ChannelType, fileOriginalName string, addressPosition AddressPosition,
	groupingType GroupingType, groupingOptionsSplitType SplitType,
	groupingOptionsSplitSize *int,
	groupingOptionsSplitSeparator *string, groupingOptionsSplitPosition *SplitPosition,
	relationships map[string]interface{},
) (BatchResponse, *errors.PingenError) {
	fileUpload := fileupload.NewFileUpload(b.apiRequestor)

	fileResponse, err := fileUpload.RequestFileUpload()
	if err != nil {
		return BatchResponse{}, err
	}

	err = fileUpload.PutFile(pathToFile, fileResponse.Data.Attributes.URL)
	if err != nil {
		return BatchResponse{}, err
	}

	return b.CreateBatch(
		fileResponse.Data.Attributes.URL,
		fileResponse.Data.Attributes.URLSignature,
		name,
		icon,
		channelType,
		fileOriginalName,
		addressPosition,
		groupingType,
		groupingOptionsSplitType,
		groupingOptionsSplitSize,
		groupingOptionsSplitSeparator,
		groupingOptionsSplitPosition,
		relationships,
	)
}

func (b *Batches) CreateBatch(
	fileURL, fileURLSignature, name string, icon Icon, channelType ChannelType, fileOriginalName string, addressPosition AddressPosition,
	groupingType GroupingType, groupingOptionsSplitType SplitType,
	groupingOptionsSplitSize *int,
	groupingOptionsSplitSeparator *string, groupingOptionsSplitPosition *SplitPosition,
	relationships map[string]interface{},
) (BatchResponse, *errors.PingenError) {
	attributes := map[string]interface{}{
		"file_url":                    fileURL,
		"file_url_signature":          fileURLSignature,
		"name":                        name,
		"icon":                        string(icon),
		"file_original_name":          fileOriginalName,
		"address_position":            string(addressPosition),
		"grouping_type":               string(groupingType),
		"grouping_options_split_type": string(groupingOptionsSplitType),
	}

	if channelType == "" {
		channelType = ChannelTypePost
	}
	attributes["channel_type"] = string(channelType)

	if groupingOptionsSplitSize != nil {
		attributes["grouping_options_split_size"] = *groupingOptionsSplitSize
	}

	if groupingOptionsSplitSeparator != nil {
		attributes["grouping_options_split_separator"] = *groupingOptionsSplitSeparator
	}

	if groupingOptionsSplitPosition != nil {
		attributes["grouping_options_split_position"] = string(*groupingOptionsSplitPosition)
	}

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       "batches",
			"attributes": attributes,
		},
	}

	if relationships != nil {
		dataMap := payload["data"].(map[string]interface{})
		dataMap["relationships"] = relationships
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/batches", b.organisationID)

	var response BatchResponse

	_, err := b.apiRequestor.PerformPostRequest(url, &response, data, nil)
	if err != nil {
		return BatchResponse{}, err
	}

	return response, nil
}

func (b *Batches) SendBatchPost(
	batchID string, deliveryProduct DeliveryProduct, printMode PrintMode, printSpectrum PrintSpectrum,
) (BatchResponse, *errors.PingenError) {
	return b.sendBatch(batchID, "batches_channel_post_send", map[string]interface{}{
		"delivery_product": string(deliveryProduct),
		"print_mode":       string(printMode),
		"print_spectrum":   string(printSpectrum),
	})
}

func (b *Batches) SendBatchEmail(batchID string) (BatchResponse, *errors.PingenError) {
	return b.sendBatch(batchID, "batches_channel_email_send", map[string]interface{}{
		"delivery_product": string(DeliveryProductElectronicEmail),
	})
}

func (b *Batches) SendBatchEbill(batchID string) (BatchResponse, *errors.PingenError) {
	return b.sendBatch(batchID, "batches_channel_ebill_send", map[string]interface{}{
		"delivery_product": string(DeliveryProductElectronicEbill),
	})
}

func (b *Batches) sendBatch(batchID, dataType string, attributes map[string]interface{}) (BatchResponse, *errors.PingenError) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         batchID,
			"type":       dataType,
			"attributes": attributes,
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/batches/%s/send", b.organisationID, batchID)

	var response BatchResponse

	_, err := b.apiRequestor.PerformPatchRequest(url, &response, data, nil)
	if err != nil {
		return BatchResponse{}, err
	}

	return response, nil
}

func (b *Batches) CancelBatch(batchID string) (interface{}, *errors.PingenError) {
	url := fmt.Sprintf("/organisations/%s/batches/%s/cancel", b.organisationID, batchID)
	return b.apiRequestor.PerformCancelRequest(url)
}

func (b *Batches) DeleteBatch(batchID string, withDeliverables bool) (interface{}, *errors.PingenError) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":   batchID,
			"type": "batches",
			"attributes": map[string]interface{}{
				// with_letters is deprecated but still required, it mirrors with_deliverables.
				"with_letters":      withDeliverables,
				"with_deliverables": withDeliverables,
			},
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/batches/%s", b.organisationID, batchID)

	return b.apiRequestor.PerformDeleteRequestWithPayload(url, data)
}

func (b *Batches) UpdateBatch(batchID, name string, icon Icon) (BatchResponse, *errors.PingenError) {
	attributes := map[string]interface{}{}

	if name != "" {
		attributes["name"] = name
	}

	if icon != "" {
		attributes["icon"] = string(icon)
	}

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         batchID,
			"type":       "batches",
			"attributes": attributes,
		},
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("/organisations/%s/batches/%s", b.organisationID, batchID)

	var response BatchResponse

	_, err := b.apiRequestor.PerformPatchRequest(url, &response, data, nil)
	if err != nil {
		return BatchResponse{}, err
	}

	return response, nil
}

func (b *Batches) GetStatistics(batchID string) (BatchStatisticsResponse, *errors.PingenError) {
	var response BatchStatisticsResponse
	url := fmt.Sprintf("/organisations/%s/batches/%s/statistics", b.organisationID, batchID)
	_, err := b.apiRequestor.PerformGetRequest(url, &response, nil, nil)
	if err != nil {
		return BatchStatisticsResponse{}, err
	}

	return response, nil
}
