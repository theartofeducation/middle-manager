package clickup

// Client handles interaction with the ClickUp API.
type Client struct {
	URL                     string
	Key                     string
	TaskStatusUpdatedSecret string
}
