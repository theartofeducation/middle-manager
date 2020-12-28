package clickup

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
