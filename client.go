package topgg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

const defaultTimeout = 3 * time.Second

// HTTPClient is an interface for HTTP client implementations.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// OptionFunc is a callback type that modifies the *Client provided.
type OptionFunc func(*Client) error

// Client is a struct that lets you interact with Top.gg API v1's endpoints.
type Client struct {
	limiter    *rate.Limiter // The client's ratelimiter.
	httpClient HTTPClient    // The client's HTTP client.
	token      string        // The client's API token.
}

// NewClient is a function that tries to create a new client instance.
func NewClient(token string, options ...OptionFunc) (*Client, error) {
	client := &Client{
		limiter:    rate.NewLimiter(1, 100),
		httpClient: &http.Client{Timeout: defaultTimeout},
		token:      token,
	}

	for _, optionFunc := range options {
		if optionFunc == nil {
			return nil, errors.New("Specified topgg.Client option func must not be null")
		} else if err := optionFunc(client); err != nil {
			return nil, fmt.Errorf("Unable to run topgg.Client option func: %w", err)
		}
	}

	return client, nil
}

// WithHTTPClient is an option func factory that customizes the client's HTTP client.
func WithHTTPClient(httpClient HTTPClient) OptionFunc {
	return func(client *Client) error {
		client.httpClient = httpClient

		return nil
	}
}

// WithTimeout is an option func factory that customizes the client's HTTP client timeout.
func WithTimeout(duration time.Duration) OptionFunc {
	return func(client *Client) error {
		httpClient, ok := client.httpClient.(*http.Client)

		if !ok {
			return errors.New("Unable to type assert Client.httpClient to *http.Client")
		}

		httpClient.Timeout = duration

		return nil
	}
}

func (client *Client) readBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	switch res.StatusCode {
	case 400:
		return nil, errors.New("Attempted to send an invalid request to the API.")
	case 401, 403:
		return nil, errors.New("Invalid Top.gg API token.")
	case 404:
		return nil, errors.New("Such query does not exist.")
	case 429:
		retryAfter, err := strconv.ParseFloat(res.Header.Get("Retry-After"), 32)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("The client is blocked by the API. Please try again in %f seconds.", retryAfter)
	case 200:
		return io.ReadAll(res.Body)
	case 204:
		return []byte{}, nil
	default:
		return nil, errors.New("Received an unexpected error from Top.gg's end.")
	}
}

func (client *Client) createRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	if client.token == "" {
		return nil, errors.New("Missing API token.")
	}

	if !client.limiter.Allow() {
		return nil, errors.New("Temporarily prevented from sending requests by local ratelimiter.")
	}

	req, err := http.NewRequest(method, BaseURL+endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+client.token)

	return req, nil
}

// GetSelf is a method that tries to get your project's information.
func (client *Client) GetSelf() (*Project, error) {
	req, err := client.createRequest("GET", "/projects/@me", nil)

	if err != nil {
		return nil, err
	}

	res, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := client.readBody(res)

	if err != nil {
		return nil, err
	}

	project := &Project{}

	if err = json.Unmarshal(body, project); err != nil {
		return nil, err
	}

	return project, nil
}

// PostCommands is a method that tries to update the application commands list in your Discord bot's Top.gg page.
func (client *Client) PostCommands(commands any) error {
	var body []byte

	switch c := commands.(type) {
	case string:
		body = []byte(c)
	case []byte:
		body = c
	default:
		b, err := json.Marshal(commands)

		if err != nil {
			return err
		}

		body = b
	}

	req, err := client.createRequest("POST", "/projects/@me/commands", bytes.NewBuffer(body))

	if err != nil {
		return err
	} else if _, err = client.httpClient.Do(req); err != nil {
		return err
	}

	return nil
}

// GetVote is a method that tries to get the latest vote information of a user on your project. Returns nil if the user has not voted.
func (client *Client) GetVote(userSource UserSource, id string) (*PartialVote, error) {
	req, err := client.createRequest("GET", "/projects/@me/votes/"+id, nil)

	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("source", string(userSource))

	req.URL.RawQuery = query.Encode()

	res, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	} else if res.StatusCode == 404 {
		return nil, nil
	}

	body, err := client.readBody(res)

	if err != nil {
		return nil, err
	}

	vote := &PartialVote{}

	if err = json.Unmarshal(body, vote); err != nil {
		return nil, err
	}

	return vote, nil
}

type paginatedVotes struct {
	Votes  []Vote `json:"data"`
	Cursor string `json:"cursor"`
}

// PaginatedVotes represents a paginated list of a project's vote information.
type PaginatedVotes struct {
	client *Client
	votes  paginatedVotes
}

// Votes is a method that gets the votes in this page.
func (votes *PaginatedVotes) Votes(client *Client) []Vote {
	return votes.votes.Votes
}

// Next is a method that tries to advance to the next page.
func (votes *PaginatedVotes) Next() (*PaginatedVotes, error) {
	req, err := votes.client.createRequest("GET", "/projects/@me/votes", nil)

	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("cursor", votes.votes.Cursor)

	req.URL.RawQuery = query.Encode()

	return votes.client.getVotes(req)
}

func (client *Client) getVotes(req *http.Request) (*PaginatedVotes, error) {
	res, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := client.readBody(res)

	if err != nil {
		return nil, err
	}

	votes := &PaginatedVotes{
		client: client,
	}

	err = json.Unmarshal(body, &votes.votes)

	if err != nil {
		return nil, err
	}

	return votes, nil
}

// GetVotes is a method that tries to get a cursor-based paginated list of votes for your project, ordered by creation date.
func (client *Client) GetVotes(since time.Time) (*PaginatedVotes, error) {
	req, err := client.createRequest("GET", "/projects/@me/votes", nil)

	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("startDate", since.Format("2006-01-02T15:04:05.000Z07:00"))

	req.URL.RawQuery = query.Encode()

	return client.getVotes(req)
}
