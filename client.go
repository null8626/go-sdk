package dbl

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const defaultTimeout = 3 * time.Second

// A HTTP client implementation.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// A callback modifies the the *Client provided.
type OptionFunc func(*Client) error

// Interact with API v1's endpoints.
type Client struct {
	sync.Mutex
	RetryAfter int           // How long the client should wait (in seconds) until it can make a request to the API again.
	limiter    *rate.Limiter // The client's ratelimiter.
	httpClient HTTPClient    // The client's HTTP client.
	token      string        // The client's API token.
}

// Creates a new client instance.
func NewClient(token string, options ...OptionFunc) (*Client, error) {
	client := &Client{
		limiter:    rate.NewLimiter(1, 60),
		httpClient: &http.Client{Timeout: defaultTimeout},
		token:      token,
	}

	for _, optionFunc := range options {
		if optionFunc == nil {
			return nil, fmt.Errorf("Specified dbl.Client option func must not be null")
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
			return fmt.Errorf("Unable to type assert Client.httpClient to *http.Client")
		}

		httpClient.Timeout = duration

		return nil
	}
}
