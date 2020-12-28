package clickup

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
