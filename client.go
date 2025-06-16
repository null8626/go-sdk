package dbl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const defaultTimeout = 3 * time.Second

// HTTPClient is an interface for HTTP client implementations.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// OptionFunc is a function that modifies the the *Client provided.
type OptionFunc func(*Client) error

// Client contains fields and methods for interacting with the Discord Bot List
// API.
type Client struct {
	sync.Mutex

	// bots/* 60/m with 1 hour block if exceeded
	// Indicates how long the timeout period is/when you will be able to send requests again
	// Upon exceeding a rate limit, this will be updated with the retry-after value.
	RetryAfter int

	limiter    *rate.Limiter
	httpClient HTTPClient
	token      string
	id         string
}

type tokenStructure struct {
	Id string `json:"id"`
}

// NewClient returns a new *Client after applying the options provided.
func NewClient(token string, options ...OptionFunc) (*Client, error) {
	tokenSections := strings.Split(token, ".")

	if len(tokenSections) != 3 {
		return nil, ErrRequireAuthentication
	}

	decodedTokenSection, err := base64.RawURLEncoding.DecodeString(tokenSections[1])

	if err != nil {
		return nil, ErrRequireAuthentication
	}

	innerTokenStructure := &tokenStructure{}

	if err = json.Unmarshal(decodedTokenSection, innerTokenStructure); err != nil {
		return nil, ErrRequireAuthentication
	}

	client := &Client{
		limiter:    rate.NewLimiter(1, 60),
		httpClient: &http.Client{Timeout: defaultTimeout},
		token:      token,
		id:         innerTokenStructure.Id,
	}

	for _, optionFunc := range options {
		if optionFunc == nil {
			return nil, fmt.Errorf("invalid nil dbl.Client option func")
		} else if err := optionFunc(client); err != nil {
			return nil, fmt.Errorf("error running dbl.Client option func: %w", err)
		}
	}

	return client, nil
}

// HTTPClientOption allows for customizing the HTTP client used.
func HTTPClientOption(httpClient HTTPClient) OptionFunc {
	return func(client *Client) error {
		client.httpClient = httpClient

		return nil
	}
}

// TimeoutOption allows for customizing the HTTP client timeout.
func TimeoutOption(duration time.Duration) OptionFunc {
	return func(client *Client) error {
		httpClient, ok := client.httpClient.(*http.Client)
		if !ok {
			return fmt.Errorf("unable to type assert Client.httpClient to *http.Client")
		}

		httpClient.Timeout = duration

		return nil
	}
}
