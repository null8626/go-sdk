package dbl

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const webhooksSecret = "secret"
const webhooksTrace = "trace"

var webhookEvents = []string{"integration.create", "integration.delete", "webhook.test", "vote.create"}

func defaultResponse(res http.ResponseWriter, name, trace string) {
	res.WriteHeader(http.StatusOK)
	res.Write(fmt.Appendf(nil, "%s:%s", name, trace))
}

func mockSignature(body []byte) string {
	timestamp := time.Now().UTC().Unix()

	mac := hmac.New(sha256.New, []byte(webhooksSecret))
	mac.Write(fmt.Appendf(nil, "%d.%s", timestamp, body))

	digest := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("t=%d,"+APIVersion+"=%s", timestamp, digest)
}

func TestWebhooks(t *testing.T) {
	webhooks := NewWebhooks(webhooksSecret)

	webhooks.OnIntegrationCreate(func(res http.ResponseWriter, payload *IntegrationCreatePayload, trace string) {
		defaultResponse(res, "integration.create", trace)
	})

	webhooks.OnIntegrationDelete(func(res http.ResponseWriter, payload *IntegrationDeletePayload, trace string) {
		defaultResponse(res, "integration.delete", trace)
	})

	webhooks.OnTest(func(res http.ResponseWriter, payload *TestPayload, trace string) {
		defaultResponse(res, "webhook.test", trace)
	})

	webhooks.OnVoteCreate(func(res http.ResponseWriter, payload *VoteCreatePayload, trace string) {
		defaultResponse(res, "vote.create", trace)
	})

	for _, event := range webhookEvents {
		rec := httptest.NewRecorder()

		file, err := os.Open(fmt.Sprintf("mocks/%s_payload.json", strings.ReplaceAll(event, ".", "_")))

		assert.Nilf(t, err, "The mock request body JSON file for %s must be available.", event)

		defer file.Close()

		body, err := io.ReadAll(file)

		assert.Nilf(t, err, "The mock request body JSON file for %s must be readable.", event)

		req := httptest.NewRequest(http.MethodPost, "/webhook", io.NopCloser(bytes.NewReader(body)))

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("x-topgg-signature", mockSignature(body))
		req.Header.Add("x-topgg-trace", webhooksTrace)

		webhooks.Handler(rec, req)

		assert.Equalf(t, http.StatusOK, rec.Code, "Sending a %s payload must be responded with 200.", event)
		assert.Equalf(t, fmt.Sprintf("%s:%s", event, webhooksTrace), rec.Body.String(), "Sending a %s payload must be responded with the expected response body.", event)
	}
}
