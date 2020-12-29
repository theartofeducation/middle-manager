package clubhouse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const apiURL = "https://api.clubhouse.io/api/v3"

// Client handles interaction with the Clubhouse API.
type Client struct {
	URL   string
	Token string
}

// NewClient creates and returns a new Clubhouse Client.
func NewClient(token string) Client {
	client := Client{
		Token: token,
	}

	return client
}

// CreateEpic creates an Epic on Clubhouse.
func (c Client) CreateEpic(name, description string) (Epic, error) {
	// TODO: check if epic exists ch246

	epic := Epic{
		Name:        name,
		Description: description,
	}

	body, err := json.Marshal(epic)
	if err != nil {
		return epic, errors.Wrap(err, "Could not create Epic body")
	}

	httpClient := &http.Client{}

	url := apiURL + "/epics"

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	request.Header.Add("Clubhouse-Token", c.Token)
	request.Header.Add("Content-Type", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return epic, errors.Wrap(err, "Could not send request to the Clubhouse API")
	}

	if response.StatusCode != http.StatusCreated {
		return epic, errors.New(fmt.Sprint("Clubhouse returned status", response.StatusCode))
	}

	return epic, nil
}
