//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/letterevents"
	"github.com/pingencom/pingen2-sdk-go/letters"
	"github.com/pingencom/pingen2-sdk-go/oauth"
	"github.com/pingencom/pingen2-sdk-go/organisations"
	"github.com/pingencom/pingen2-sdk-go/userassociations"
	"github.com/pingencom/pingen2-sdk-go/users"
	"github.com/pingencom/pingen2-sdk-go/webhooks"
	"github.com/pingencom/pingen2-sdk-go/batches"
	"github.com/pingencom/pingen2-sdk-go/batchevents"
	"github.com/pingencom/pingen2-sdk-go/ebills"
	"github.com/pingencom/pingen2-sdk-go/emails"
)

func prettyPrint(label string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("%s: %+v\n", label, data)
		return
	}
	fmt.Printf("%s:\n%s\n\n", label, string(jsonData))
}

func main() {
	config, _ := pingen2sdk.InitSDK(
		"clientId",
		"clientSecret",
		"staging",
	)

	params := map[string]string{
		"grant_type": "client_credentials",
		"scope":      "letter batch webhook organisation_read user email ebill",
	}

	tokenResp, err := oauth.GetToken(config, params)
	if err != nil {
		log.Fatalf("Error obtaining token: %v", err)
	}
	accessToken := tokenResp["access_token"].(string)
	fmt.Println("Access token obtained")

	apiRequestor := api.NewAPIRequestor(accessToken, config)

	params = map[string]string{}
	headers := map[string]string{}

	userClient := users.NewUsers(apiRequestor)
	userResp, _ := userClient.GetDetails(params, headers)
	prettyPrint("USER DETAILS:", userResp.Data)

	orgClient := organisations.NewOrganisations(apiRequestor)
	orgCollection, _ := orgClient.GetCollection(params, headers)
	prettyPrint("ORGANISATIONS LIST:", orgCollection.Data)

	organisationID := orgCollection.Data[0].ID

	orgDetails, _ := orgClient.GetDetails(organisationID, params, headers)
	prettyPrint("ORGANISATION DETAILS:", orgDetails.Data)

	userAssocClient := userassociations.NewUserAssociations(apiRequestor)
	assocCollection, _ := userAssocClient.GetCollection(params, headers)
	prettyPrint("COLLECTION OF ASSOCIATIONS:", assocCollection.Data)

	letterClient := letters.NewLetters(organisationID, apiRequestor)

	fmt.Println("UPLOAD AND CREATE LETTER")
	letterResp, _ := letterClient.UploadAndCreate(
		"/app/example/testFile.pdf",
		"sdk.pdf",
		"left",
		false,
		"fast",
		"simplex",
		"color",
		"",
		nil,
	)
	prettyPrint("Letter created:", letterResp.Data)

    time.Sleep(2 * time.Second)

	letterID := letterResp.Data.ID

	fmt.Println("LETTER EVENTS")
	letterEventsClient := letterevents.NewLetterEvents(organisationID, apiRequestor)
	letterEvents, _ := letterEventsClient.GetCollection(letterID, params, headers)
	prettyPrint("LETTER EVENTS:", letterEvents.Data)

	fmt.Println("SEND LETTER")
	sendResp, _ := letterClient.Send(letterID, "fast", "simplex", "color")
	prettyPrint("SEND LETTER RESPONSE:", sendResp.Data)

	time.Sleep(2 * time.Second)

	// Cancel letter
	fmt.Println("CANCEL LETTER")
	cancelResp, _ := letterClient.Cancel(letterID)
	prettyPrint("CANCEL LETTER RESPONSE:", cancelResp)

	// Delete letter
	fmt.Println("DELETE LETTER")
	delResp, _ := letterClient.Delete(letterID)
	prettyPrint("DELETE LETTER RESPONSE:", delResp)

	webhookClient := webhooks.NewWebhooks(organisationID, apiRequestor)

	fmt.Println("CREATE WEBHOOK")
	webhookResp, _ := webhookClient.Create(
		"issues",
		"https://valid-url",
		"d09a095a0d1d2ae896f985c0fff1ad51",
	)
	prettyPrint("Webhook created:", webhookResp.Data)

	webhookID := webhookResp.Data.ID

	fmt.Println("GET DETAILS OF WEBHOOK")
	webhookDetails, _ := webhookClient.GetDetails(webhookID, params, headers)
	prettyPrint("Webhook details:", webhookDetails.Data)

	fmt.Println("GET LIST OF WEBHOOKS")
	webhookList, _ := webhookClient.GetCollection(params, headers)
	prettyPrint("Webhook list:", webhookList.Data)

	fmt.Println("DELETE WEBHOOK")
	delWebhookResp, _ := webhookClient.Delete(webhookID)
	prettyPrint("DELETE WEBHOOK RESPONSE:", delWebhookResp)

	emailClient := emails.NewEmails(organisationID, apiRequestor)

	fmt.Println("UPLOAD AND CREATE EMAIL")
	emailMetaData := map[string]interface{}{
		"sender_name":    "Test Example",
        "recipient_email":  "info@test.com",
		"recipient_name":    "R_Example",
        "reply_email":  "info_reply@test.com",
        "reply_name":     "Reply Example",
        "subject":    "Your new invoice Number xyz",
        "content": "Dear recipient\\n\\nAttached is your invoice",
	}
	emailResp, _ := emailClient.UploadAndCreate(
		"/app/example/testFile.pdf",
		"sdk.pdf",
		false,
		emailMetaData,
		nil,
	)
	prettyPrint("Email created:", emailResp.Data)

    time.Sleep(2 * time.Second)

    ebillClient := ebills.NewEbills(organisationID, apiRequestor)

	fmt.Println("UPLOAD AND CREATE EMAIL")
	ebillMetaData := map[string]interface{}{
		"invoice_number": "Invoice 8051",
        "invoice_date": "2025-01-01",
        "invoice_due_date": "2025-01-01",
        "recipient_identifier": "41010560425610173",
	}

    relationships := map[string]interface{}{
        "preset": map[string]interface{}{
            "data": map[string]interface{}{
                "id": "7e500ff8-3407-45d1-9d1e-b7de227d1081",
                "type": "presets",
            },
        },
    }

	ebillResp, _ := ebillClient.UploadAndCreate(
		"/app/example/testFile.pdf",
		"sdk.pdf",
		false,
		ebillMetaData,
		relationships,
	)
	prettyPrint("Ebill created:", ebillResp.Data)

    time.Sleep(2 * time.Second)
}
