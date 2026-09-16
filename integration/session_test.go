//go:build integration
// +build integration

package integration

import (
	"fmt"
	"sync"
	"testing"

	"github.com/pingencom/pingen2-sdk-go"
	"github.com/pingencom/pingen2-sdk-go/organisations"
)

type session struct {
	*pingen2sdk.Client

	creds credentials
}

var (
	sessionOnce   sync.Once
	sharedSession *session
	sessionErr    error
)

func newSession(t *testing.T) *session {
	t.Helper()

	creds := requireCredentials(t)

	sessionOnce.Do(func() {
		sharedSession, sessionErr = buildSession(creds)
	})

	if sessionErr != nil {
		t.Fatalf("Integration session setup failed: %v", sessionErr)
	}

	return sharedSession
}

func buildSession(creds credentials) (*session, error) {
	pingen, err := pingen2sdk.New(pingen2sdk.Options{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Environment:  environment,
		Scope:        Scope,
	})
	if err != nil {
		return nil, err
	}

	organisationID, err := resolveOrganisationID(creds, pingen.Organisations)
	if err != nil {
		return nil, err
	}

	return &session{
		Client: pingen.ForOrganisation(organisationID),
		creds:  creds,
	}, nil
}

func resolveOrganisationID(creds credentials, client *organisations.Organisations) (string, error) {
	if creds.OrganisationID != "" {
		return creds.OrganisationID, nil
	}

	collection, err := client.GetCollection(nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list organisations: %w", err)
	}

	if len(collection.Data) == 0 {
		return "", fmt.Errorf("no organisations returned - check the staging credentials")
	}

	return collection.Data[0].ID, nil
}
