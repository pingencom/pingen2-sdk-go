package pingen2sdk

import (
	"github.com/pingencom/pingen2-sdk-go/api"
	"github.com/pingencom/pingen2-sdk-go/batches"
	"github.com/pingencom/pingen2-sdk-go/batchevents"
	"github.com/pingencom/pingen2-sdk-go/config"
	"github.com/pingencom/pingen2-sdk-go/ebillevents"
	"github.com/pingencom/pingen2-sdk-go/ebills"
	"github.com/pingencom/pingen2-sdk-go/emailevents"
	"github.com/pingencom/pingen2-sdk-go/emails"
	"github.com/pingencom/pingen2-sdk-go/letterevents"
	"github.com/pingencom/pingen2-sdk-go/letters"
	"github.com/pingencom/pingen2-sdk-go/oauth"
	"github.com/pingencom/pingen2-sdk-go/organisations"
	"github.com/pingencom/pingen2-sdk-go/userassociations"
	"github.com/pingencom/pingen2-sdk-go/users"
	"github.com/pingencom/pingen2-sdk-go/webhooks"
)

// Config is the SDK configuration. It is an alias, so a *Config obtained from
// InitSDK is the same type the api and oauth packages take.
type Config = config.Config

func InitSDK(clientID, clientSecret, environment string) (*Config, error) {
	return config.InitSDK(clientID, clientSecret, environment)
}

func InitSDKWithoutCredentials(environment string) *Config {
	return config.InitSDKWithoutCredentials(environment)
}

type Options struct {
	ClientID       string
	ClientSecret   string
	Environment    string
	Scope          string
	OrganisationID string
	AccessToken    string
}

type Client struct {
	Config    *Config
	OAuth     *oauth.OAuth
	Requestor *api.APIRequestor

	Organisations    *organisations.Organisations
	Users            *users.Users
	UserAssociations *userassociations.UserAssociations
	Letters          *letters.Letters
	LetterEvents     *letterevents.LetterEvents
	Batches          *batches.Batches
	BatchEvents      *batchevents.BatchEvents
	Webhooks         *webhooks.Webhooks
	Emails           *emails.Emails
	EmailEvents      *emailevents.EmailEvents
	Ebills           *ebills.Ebills
	EbillEvents      *ebillevents.EbillEvents

	organisationID string
}

func New(opts Options) (*Client, error) {
	c := &Client{organisationID: opts.OrganisationID}

	if opts.AccessToken != "" {
		c.Config = config.InitSDKWithoutCredentials(opts.Environment)
		c.Requestor = api.NewAPIRequestor(opts.AccessToken, c.Config)
	} else {
		config, err := config.InitSDK(opts.ClientID, opts.ClientSecret, opts.Environment)
		if err != nil {
			return nil, err
		}

		c.Config = config
		c.OAuth = oauth.NewOAuth(config, opts.Scope)
		c.Requestor = api.NewAPIRequestorWithTokenSource(c.OAuth, config)
	}

	c.bindResources()

	return c, nil
}

func (c *Client) OrganisationID() string {
	return c.organisationID
}

func (c *Client) ForOrganisation(organisationID string) *Client {
	scoped := *c
	scoped.organisationID = organisationID
	scoped.bindResources()

	return &scoped
}

func (c *Client) bindResources() {
	c.Organisations = organisations.NewOrganisations(c.Requestor)
	c.Users = users.NewUsers(c.Requestor)
	c.UserAssociations = userassociations.NewUserAssociations(c.Requestor)
	c.Letters = letters.NewLetters(c.organisationID, c.Requestor)
	c.LetterEvents = letterevents.NewLetterEvents(c.organisationID, c.Requestor)
	c.Batches = batches.NewBatches(c.organisationID, c.Requestor)
	c.BatchEvents = batchevents.NewBatchEvents(c.organisationID, c.Requestor)
	c.Webhooks = webhooks.NewWebhooks(c.organisationID, c.Requestor)
	c.Emails = emails.NewEmails(c.organisationID, c.Requestor)
	c.EmailEvents = emailevents.NewEmailEvents(c.organisationID, c.Requestor)
	c.Ebills = ebills.NewEbills(c.organisationID, c.Requestor)
	c.EbillEvents = ebillevents.NewEbillEvents(c.organisationID, c.Requestor)
}
