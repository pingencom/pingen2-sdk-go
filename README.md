# Pingen Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/pingencom/pingen2-sdk-go)](https://pkg.go.dev/github.com/pingencom/pingen2-sdk-go)

The official [Pingen][pingen] Go client library for letter sending service.

## Requirements

- Go 1.24 or later

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

Examples of using you can find in Example folder.  

```sh
  config, _ := pingen2sdk.InitSDK(
    "yourClientId",
    "yourClientSecret",
    "staging",
)

params := map[string]string{
    "grant_type": "client_credentials",
    "scope":      "letter batch webhook organisation_read user",
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

organisationID := 'YOUR_ORGANISATION_ID'

letterClient := letters.NewLetters(organisationID, apiRequestor)

fmt.Println("UPLOAD, CREATE AND AUTOSEND LETTER")
letterResp, _ := letterClient.UploadAndCreate(
    "/app/example/testFile.pdf",
    "sdk.pdf",
    "left",
    true,
    "fast",
    "simplex",
    "color",
    "",
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

For any requests, bug or comments, please [open an issue][issues] or [submit a
pull request][pulls].

[api-docs]: https://api.pingen.com/documentation
[issues]: https://github.com/pingencom/pingen2-sdk-go/issues/new
[pulls]: https://github.com/pingencom/pingen2-sdk-go/pulls
[pingen]: https://pingen.com
