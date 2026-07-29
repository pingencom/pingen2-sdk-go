# Pingen Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/pingencom/pingen2-sdk-go)](https://pkg.go.dev/github.com/pingencom/pingen2-sdk-go)

The official [Pingen][pingen] Go client library for letter sending service.

## Requirements

- Go 1.25 or later

## Installation

Make sure your project is using Go Modules (it will have a `go.mod` file in its
root if it already is):

```sh
go mod init
```

Then, reference pingen2-sdk-go in a Go program with `import`:

```go
import (
	"github.com/pingencom/pingen2-sdk-go"
)
```

Run any of the normal `go` commands (`build`/`install`/`test`). The Go
toolchain will resolve and fetch the pingen2-sdk-go module automatically.

Alternatively, you can also explicitly `go get` the package into a project:

```bash
go get -u github.com/pingencom/pingen2-sdk-go
```

# Environments

We have two Environments available: Production and Staging, see [Environments](https://api.pingen.com/documentation#section/Basics/Environments)

This SDK supports staging as well. **When initiating the resource** the optional environment attribute should be set to the 'staging'.

# Usage

The simplest way to integrate is using the client credentials grant, see [Grant type](https://api.pingen.com/documentation#section/Authentication/Which-grant-type-should-i-use)

Usage examples are below. For complete flows against a live API — creating,
sending, cancelling and deleting every delivery type — read the integration
suite in `integration/`.

### Client (recommended)

`pingen2sdk.New` is the single entry point: it builds the configuration, wires an
OAuth token source that refreshes the access token automatically before it
expires, and exposes every resource client ready to use.

Keep your credentials out of the source and read them from the environment:

```sh
export PINGEN2_CLIENT_ID=yourClientId
export PINGEN2_CLIENT_SECRET=yourClientSecret
export PINGEN2_ORGANIZATION_ID=yourOrganisationId
```

The SDK never reads the environment itself — pass the values in explicitly, so
the source of your configuration stays yours to choose.

```go
import (
    "log"
    "os"

    "github.com/pingencom/pingen2-sdk-go"
)

pingen, err := pingen2sdk.New(pingen2sdk.Options{
    ClientID:       os.Getenv("PINGEN2_CLIENT_ID"),
    ClientSecret:   os.Getenv("PINGEN2_CLIENT_SECRET"),
    Environment:    "staging",
    Scope:          "letter batch webhook organisation_read email ebill",
    OrganisationID: os.Getenv("PINGEN2_ORGANIZATION_ID"),
})
if err != nil {
    log.Fatalf("Error creating client: %v", err)
}

params := map[string]string{}
headers := map[string]string{}

letterResp, pErr := pingen.Letters.UploadAndCreate(
    "testdata/test.pdf",
    "sdk.pdf",
    "left",
    true,
    "fast",
    "simplex",
    "color",
    "",
    nil,
    nil,
)
if pErr != nil {
    log.Fatalf("Error creating letter: %v", pErr)
}

letterEvents, _ := pingen.LetterEvents.GetCollection(letterResp.Data.ID, params, headers)
log.Println("LETTER EVENTS:", letterEvents.Data)
```

`pingen2sdk.New` also accepts a static `AccessToken` instead of client credentials,
and `pingen.ForOrganisation("other-organisation-id")` returns a copy of the
client bound to another organisation while sharing the same token source.

### Manual token handling (still supported)

The pre-`pingen2sdk.New` way of assembling the SDK by hand keeps working
unchanged — no existing integration has to be rewritten. You give up automatic
token refresh: the access token below is fetched once and is your job to renew.

Use it when you manage tokens yourself, and `pingen2sdk.New` otherwise.

Configuration lives in
`github.com/pingencom/pingen2-sdk-go/config`; `pingen2sdk.Config` is an alias for
it, so either import works:

```go
  config, _ := pingen2sdk.InitSDK(
    "yourClientId",
    "yourClientSecret",
    "staging",
)

params := map[string]string{
    "grant_type": "client_credentials",
    "scope":      "letter batch webhook organisation_read email ebill",
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

organisationID := "YOUR_ORGANISATION_ID"

letterClient := letters.NewLetters(organisationID, apiRequestor)

fmt.Println("UPLOAD, CREATE AND AUTOSEND LETTER")
letterResp, _ := letterClient.UploadAndCreate(
    "testdata/test.pdf",
    "sdk.pdf",
    "left",
    true,
    "fast",
    "simplex",
    "color",
    "",
    nil,
    nil,
)
fmt.Println("Letter created and sent:", letterResp.Data)

time.Sleep(2 * time.Second)

letterID := letterResp.Data.ID

fmt.Println("LETTER EVENTS")
letterEventsClient := letterevents.NewLetterEvents(organisationID, apiRequestor)
letterEvents, _ := letterEventsClient.GetCollection(letterID, params, headers)
fmt.Println("LETTER EVENTS:", letterEvents.Data)
```

## Documentation

For a comprehensive list of examples, check out the [API
documentation][api-docs].

On the right-hand side of every endpoint you can see request samples for Python and other languages, which you can copy and paste into your application.

## Support

New features and bug fixes are released on the latest major version of the Pingen Go client library. If you are on an older major version, we recommend that you upgrade to the latest in order to use the new features and bug fixes including those for security vulnerabilities. Older major versions of the package will continue to be available for use, but will not be receiving any updates.

## Development

Pull requests from the community are welcome. If you submit one, please keep
the following guidelines in mind:

1. Code must be `go fmt` compliant.
2. All types, structs and funcs should be documented.
3. Ensure that `make test` succeeds.

## Testing

We use makefile for conveniently running development tasks. You can use them directly, or copy the commands out of the `makefile`. To our help docs, run `make`.

Run all tests, lint and formatting:

```sh
  make ci
```

Run all tests with coverage:

```sh
  make test-cov
```

### Integration tests

Next to the unit tests there is an integration suite in `integration/`, hidden
behind the `integration` build tag, that runs against the **real Pingen staging
API**: it uploads the fixtures in `testdata/`, creates letters, batches, emails
and ebills, sends and cancels them and deletes everything again.

To run it, copy the example environment file and fill in your staging
credentials:

```sh
  cp .env.example .env
  # edit .env: PINGEN2_CLIENT_ID / PINGEN2_CLIENT_SECRET (staging)
  make test-integration
```

`PINGEN2_ORGANIZATION_ID` is optional — when it is empty the suite uses the
first organisation your credentials can see. `PINGEN2_ORGANIZATION_NAME` is
optional too and is asserted only when set. `PINGEN2_USE_STAGING` defaults to
`true` and should stay that way.

The suite is **skipped automatically when the credentials are absent**, so
`make test`, `make check` and `make ci` stay offline and green without a `.env`
file — the plain `go test ./...` used by those targets never builds the tagged
files. `make lint-integration` runs staticcheck over the tagged files as well.

The OAuth scope the suite requests is `letter batch webhook organisation_read
email ebill`. Note there is no `user` scope — the `/user` endpoints are covered
by the scopes above, and asking for `user` makes the token request fail with
`invalid_scope`.

### Running the tests in Docker

If you would rather not install a Go toolchain locally, the repository ships a
`Dockerfile` and a `docker-compose.yml` that mount the working tree at `/app`.

Using compose (the container stays up, so you can run several commands in it).
Pass `--build` so a stale image does not shadow a change to the `Dockerfile` —
the official Go images set `GOTOOLCHAIN=local`, so an image older than the `go`
directive in `go.mod` fails outright instead of fetching a newer toolchain:

```sh
  docker compose up -d --build
  docker compose exec go-sdk make test              # unit tests
  docker compose exec go-sdk make test-integration  # integration suite, needs .env
  docker compose down
```

Or as a one-off, without compose:

```sh
  docker run --rm -v "$PWD":/app -w /app golang:1.26 make test
  docker run --rm -v "$PWD":/app -w /app golang:1.26 make test-integration
```

Because the working tree is mounted, the container reads the same `.env` and
`testdata/` fixtures as a local run. Add `-v go-mod-cache:/go/pkg/mod` to reuse
the module cache between runs instead of re-downloading dependencies each time.

For any requests, bug or comments, please [open an issue][issues] or [submit a
pull request][pulls].

[api-docs]: https://api.pingen.com/documentation
[issues]: https://github.com/pingencom/pingen2-sdk-go/issues/new
[pulls]: https://github.com/pingencom/pingen2-sdk-go/pulls
[pingen]: https://pingen.com
