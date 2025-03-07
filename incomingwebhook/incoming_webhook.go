package incomingwebhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/pingencom/pingen2-sdk-go/errors"
)

type WebhookEvent struct {
	Payload string
}

type IncomingWebhook struct{}

func (iw *IncomingWebhook) ConstructEvent(payload string, headers map[string]string, secret string) (*WebhookEvent, error) {
	if err := VerifyHeader(payload, headers, secret); err != nil {
		return nil, err
	}

	return &WebhookEvent{Payload: payload}, nil
}

func VerifyHeader(payload string, headers map[string]string, secret string) error {
	signature, ok := headers["Signature"]
	if !ok {
		return errors.NewWebhookSignatureException("signature missing")
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if signature != expectedSig {
		return errors.NewWebhookSignatureException("webhook signature matching failed")
	}

	return nil
}
