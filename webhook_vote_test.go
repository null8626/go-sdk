package dbl

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testToken = "wblAV@d!Od9uL761Rz23BEQC$#YCJdQ0nDlZUEfnDxY"
)

var (
	testPayload  = []byte(`{"bot":"441751906428256277","user":"105122038586286080","type":"upvote","isWeekend":false,"query":""}`)
	testListener = NewWebhookListener(testToken, "/votes", func(payload []byte) {
		vote, err := NewWebhookVotePayload(payload)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to parse webhook payload: %s", err)
			return
		}

		fmt.Printf("A user with the ID of %s has voted us on Top.gg!", vote.VoterId)
	})
)

func TestWebhookVoteHookMethod(t *testing.T) {
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/votes", bytes.NewBuffer(testPayload))

	testListener.handleRequest(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "GET method should not be allowed")
}

func TestWebhookVoteHookAuthentication(t *testing.T) {
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/votes", bytes.NewBuffer(testPayload))

	testListener.handleRequest(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code, "Unauthorized request should not be processed")
}

func TestWebhookVoteProcessing(t *testing.T) {
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/votes", bytes.NewBuffer(testPayload))
	req.Header.Set("Authorization", testToken)

	testListener.handleRequest(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code, "Request should succeed w/o content")
}
