package dbl

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

// A HTTP client implementation.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// A callback that modifies the *Client provided.
type OptionFunc func(*Client) error

// Interact with API v1's endpoints.
type Client struct {
	limiter    *rate.Limiter // The client's ratelimiter.
	httpClient HTTPClient    // The client's HTTP client.
	token      string        // The client's API token.
}

// Creates a new client instance.
func NewClient(token string, options ...OptionFunc) (*Client, error) {
	client := &Client{
		limiter:    rate.NewLimiter(1, 100),
		httpClient: &http.Client{Timeout: defaultTimeout},
		token:      token,
	}

	for _, optionFunc := range options {
		if optionFunc == nil {
			return nil, errors.New("Specified dbl.Client option func must not be null")
		}

		if err := optionFunc(client); err != nil {
			return nil, fmt.Errorf("Unable to run dbl.Client option func: %w", err)
		}
	}

	return client, nil
}

// Creates an option func that customizes the client's HTTP client.
func HTTPClientOption(httpClient HTTPClient) OptionFunc {
	return func(client *Client) error {
		client.httpClient = httpClient

		return nil
	}
}

// Creates an option func that customizes the client's HTTP client timeout.
func TimeoutOption(duration time.Duration) OptionFunc {
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
		return nil, ErrInvalidRequest
	case 401:
		return nil, ErrInvalidToken
	case 404:
		return nil, ErrNotFound
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
		return nil, ErrServerSide
	}
}

func (client *Client) createRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	if client.token == "" {
		return nil, ErrInvalidToken
	}

	if !client.limiter.Allow() {
		return nil, ErrLocalRatelimit
	}

	req, err := http.NewRequest(method, BaseURL+endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+client.token)

	return req, nil
}

// Tries to get your project's information.
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

// Tries to update the application commands list in your Discord bot's Top.gg page.
func (client *Client) PostCommands(commands []any) error {
	body, err := json.Marshal(commands)

	if err != nil {
		return err
	}

	req, err := client.createRequest("POST", "/projects/@me/commands", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	if _, err = client.httpClient.Do(req); err != nil {
		return err
	}

	return nil
}

// Tries to get the latest vote information of a user on your project. Returns nil if the user has not voted.
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
	}

	if res.StatusCode == 404 {
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

// A paginated list of a project's vote information.
type PaginatedVotes struct {
	client *Client
	votes  paginatedVotes
}

// Gets the votes in this page.
func (votes *PaginatedVotes) Votes(client *Client) []Vote {
	return votes.votes.Votes
}

// Tries to advance to the next page.
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

// Tries to get a cursor-based paginated list of votes for your project, ordered by creation date.
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
