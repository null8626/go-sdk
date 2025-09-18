package dbl

import (
	"encoding/json"
	"io"
	"net/http"
)

type ratelimitResponse struct {
	RetryAfter int `json:"retry-after"`
}

func (c *Client) readBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	switch res.StatusCode {
	case 400:
		return nil, ErrInvalidRequest
	case 401:
		return nil, ErrUnauthorizedRequest
	case 200:
		break
	default:
		return nil, ErrRequestFailed
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 429 {
		rr := &ratelimitResponse{}

		err = json.Unmarshal(body, rr)

		if err != nil {
			return nil, err
		}

		c.Lock()
		c.RetryAfter = rr.RetryAfter
		c.Unlock()

		return nil, ErrRemoteRatelimit
	}

	return body, nil
}

func (c *Client) createRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	if c.token == "" {
		return nil, ErrRequireAuthentication
	}

	if !c.limiter.Allow() {
		return nil, ErrLocalRatelimit
	}

	req, err := http.NewRequest(method, BaseURL+endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
