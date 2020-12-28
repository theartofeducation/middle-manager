package clubhouse

// Client handles interaction with the Clubhouse API.
type Client struct {
	URL   string
	Token string
}

// NewClient creates and returns a new Clubhouse Client.
func NewClient(url, token string) Client {
	client := Client{
		URL:   url,
		Token: token,
	}

	return client
}
