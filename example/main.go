package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/letters"
	"github.com/pingencom/pingen2-sdk-go/oauth"
	"github.com/pingencom/pingen2-sdk-go/organisations"
	"github.com/pingencom/pingen2-sdk-go/userassociations"
	"github.com/pingencom/pingen2-sdk-go/users"
	"github.com/pingencom/pingen2-sdk-go/webhooks"
)

func main() {
	pingen2sdk.DefaultConfig.ClientID = "I5WGC4NAW0MF40S569RP"
	pingen2sdk.DefaultConfig.ClientSecret = "7gXoaESAB3/CE9enImXN32DmJ7byEsaswSQaDarFdTQSNrFF1nzb+NDHovqRqU1onPUiFdlyK9mvb0AL"

	useStaging := true
	params := map[string]string{
		"grant_type": "client_credentials",
		"scope":      "letter batch webhook organisation_read",
	}

	tokenResp, err := oauth.GetToken(useStaging, params)
	if err != nil {
		log.Fatalf("Error obtaining token: %v", err)
	}
	accessToken := tokenResp["access_token"].(string)
	fmt.Println("Access token obtained")

	apiRequestor := api.NewAPIRequestor(accessToken, true)

	params = map[string]string{}
	headers := map[string]string{}

	userClient := users.NewUsers(apiRequestor)
	userResp, err := userClient.GetDetails(params, headers)
	fmt.Println("USER DETAILS:", userResp.Data)

	orgClient := organisations.NewOrganisations(apiRequestor)
	orgCollection, err := orgClient.GetCollection(params, headers)
	fmt.Println("ORGANISATIONS LIST:", orgCollection.Data)

	organisationID := orgCollection.Data[0].ID

	orgDetails, err := orgClient.GetDetails(organisationID, params, headers)
	fmt.Println("ORGANISATION DETAILS:", orgDetails.Data)

	userAssocClient := userassociations.New(accessToken, true)
	assocCollection, err := userAssocClient.GetCollection()
	if err != nil {
		log.Fatalf("Error fetching user associations: %v", err)
	}
	fmt.Println("COLLECTION OF ASSOCIATIONS:", assocCollection.Data)

	letterClient := letters.NewLetters(organisationID, apiRequestor)

	fmt.Println("UPLOAD AND CREATE LETTER")
	letterResp, err := letterClient.UploadAndCreate(
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
	if err != nil {
		log.Fatalf("Error uploading and creating letter: %v", err)
	}
	fmt.Println("Letter created:", letterResp.Data)

	letterID := "e17eccd3-f8a6-4e97-9f76-f9eb56b11208"

	fmt.Println("Get File")
	data, err := letterClient.GetFile(letterID)
	defer data.Close()

	content, _ := io.ReadAll(data)
	fmt.Printf("File content: %s\n", string(content))

	fmt.Println("LETTER EVENTS")
	letterEvents, err := letterClient.GetEvents(letterID)
	if err != nil {
		log.Fatalf("Error fetching letter events: %v", err)
	}
	fmt.Println("LETTER EVENTS:", letterEvents.Data)

	fmt.Println("SEND LETTER")
	sendResp, err := letterClient.Send(letterID, "fast", "simplex", "color")
	if err != nil {
		log.Fatalf("Error sending letter: %v", err)
	}
	fmt.Println("SEND LETTER RESPONSE:", sendResp.Data)

	time.Sleep(2 * time.Second)

	// Cancel letter
	fmt.Println("CANCEL LETTER")
	cancelResp, err := letterClient.Cancel(letterID)
	if err != nil {
		log.Printf("Error canceling letter: %v", err)
	} else {
		fmt.Println("CANCEL LETTER RESPONSE:", cancelResp)
	}

	// Delete letter
	fmt.Println("DELETE LETTER")
	delResp, err := letterClient.Delete(letterID)
	if err != nil {
		log.Printf("Error deleting letter: %v", err)
	} else {
		fmt.Println("DELETE LETTER RESPONSE:", delResp.Data)
	}

	webhookClient := webhooks.New(organisationID, accessToken, true)

	fmt.Println("CREATE WEBHOOK")
	webhookResp, err := webhookClient.Create(
		"issues",
		"https://valid-url",
		"d09a095a0d1d2ae896f985c0fff1ad51",
	)
	if err != nil {
		log.Fatalf("Error creating webhook: %v", err)
	}
	fmt.Println("Webhook created:", webhookResp.Data)

	webhookID := webhookResp.Data.ID

	fmt.Println("GET DETAILS OF WEBHOOK")
	webhookDetails, err := webhookClient.GetDetails(webhookID)
	if err != nil {
		log.Fatalf("Error fetching webhook details: %v", err)
	}
	fmt.Println("Webhook details:", webhookDetails.Data)

	fmt.Println("GET LIST OF WEBHOOKS")
	webhookList, err := webhookClient.GetCollection()
	if err != nil {
		log.Fatalf("Error fetching list of webhooks: %v", err)
	}
	fmt.Println("Webhook list:", webhookList.Data)

	fmt.Println("DELETE WEBHOOK")
	delWebhookResp, err := webhookClient.Delete(webhookID)
	if err != nil {
		log.Printf("Error deleting webhook: %v", err)
	} else {
		fmt.Println("DELETE WEBHOOK RESPONSE:", delWebhookResp.Data)
	}
}
