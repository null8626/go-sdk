package dbl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRoute struct {
	Method   string
	Endpoint *regexp.Regexp
	Name     string
}

const basePath = "/api/" + APIVersion

var mockRoutes = []mockRoute{
	{
		Method:   http.MethodGet,
		Endpoint: regexp.MustCompile(fmt.Sprintf(`^%s\/projects\/@me$`, basePath)),
		Name:     "get_self",
	}, {
		Method:   http.MethodPost,
		Endpoint: regexp.MustCompile(fmt.Sprintf(`^%s\/projects\/@me\/commands$`, basePath)),
		Name:     "",
	}, {
		Method:   http.MethodGet,
		Endpoint: regexp.MustCompile(fmt.Sprintf(`^%s\/projects\/@me\/votes/\d+$`, basePath)),
		Name:     "get_vote",
	}, {
		Method:   http.MethodGet,
		Endpoint: regexp.MustCompile(fmt.Sprintf(`^%s\/projects\/@me\/votes$`, basePath)),
		Name:     "get_votes",
	},
}

type mockHttpClient struct{}

func (mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	res := &http.Response{
		StatusCode: http.StatusNotFound,
		Request:    req,
		Body:       io.NopCloser(nil),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	for _, route := range mockRoutes {
		if req.Method == route.Method && route.Endpoint.MatchString(req.URL.Path) {
			if route.Name == "" {
				res.StatusCode = http.StatusNoContent
			} else {
				file, err := os.Open(fmt.Sprintf("mocks/%s.json", route.Name))

				if err != nil {
					return nil, err
				}

				res.StatusCode = http.StatusOK
				res.Body = file
			}

			break
		}
	}

	return res, nil
}

func TestClient(t *testing.T) {
	client, err := NewClient(
		"token",
		HTTPClientOption(&mockHttpClient{}),
	)

	assert.Nil(t, err, "NewClient() must work.")

	_, err = client.GetSelf()

	assert.Nil(t, err, "Client.GetSelf() must work.")

	err = client.PostCommands([]any{})

	assert.Nil(t, err, "Client.PostCommands() must work.")

	_, err = client.GetVote(UserDiscord, "123456")

	assert.Nil(t, err, "Client.GetVote() must work.")

	_, err = client.GetVote(UserTopgg, "123456")

	assert.Nil(t, err, "Client.GetVote() must work.")

	firstPage, err := client.GetVotes(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))

	assert.Nil(t, err, "Client.GetVotes() must work.")

	_, err = firstPage.Next()

	assert.Nil(t, err, "PaginatedVotes.Next() must work.")
}
