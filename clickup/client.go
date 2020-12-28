package clickup

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// Client handles interaction with the ClickUp API.
type Client struct {
	URL                     string
	Key                     string
	TaskStatusUpdatedSecret string
}

// NewClient creates and returns a new ClickUp Client.
func NewClient(url, key, taskStatusUpdatedSecret string) Client {
	client := Client{
		URL:                     url,
		Key:                     key,
		TaskStatusUpdatedSecret: taskStatusUpdatedSecret,
	}

	return client
}

// VerifySignature validates a Webhook's signature.
func (c Client) VerifySignature(signature string, body []byte) error {
	secret := []byte(c.TaskStatusUpdatedSecret)

	hash := hmac.New(sha256.New, secret)
	hash.Write(body)
	generatedSignature := hex.EncodeToString(hash.Sum(nil))

	if signature == generatedSignature {
		return ErrSignatureMismatch
	}

	return nil
}

// GetWebhook parses a Webhook's body and returns a Webhook struct.
func (c Client) GetWebhook(body io.ReadCloser) (Webhook, error) {
	defer body.Close()

	var webhook Webhook

	if err := json.NewDecoder(body).Decode(&webhook); err != nil {
		return webhook, errors.Wrap(err, "Could not parse Webhook body")
	}

	return webhook, nil
}
