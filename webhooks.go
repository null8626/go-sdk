package dbl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type rawListener = func(http.ResponseWriter, json.RawMessage, string)

// Webhooks represents a Top.gg webhook manager.
type Webhooks struct {
	Secret    string
	listeners map[string]rawListener
}

// NewWebhooks is a function that creates a new webhook manager instance.
func NewWebhooks(Secret string) *Webhooks {
	return &Webhooks{
		Secret:    Secret,
		listeners: make(map[string]rawListener),
	}
}

func reportPayloadUnmarshalFailure(body []byte, err error) {
	slog.Warn(fmt.Sprintf("Unable to parse Top.gg webhook payload. Please report this bug to the SDK maintainers.\nCause: %s\n--- BEGIN BODY DUMP ---\n%s\n--- END BODY DUMP ---", err.Error(), body))
}

func newRawListener[P Payload](listener func(http.ResponseWriter, *P, string)) rawListener {
	return func(res http.ResponseWriter, rawPayload json.RawMessage, trace string) {
		var payload P

		if err := json.Unmarshal(rawPayload, &payload); err != nil {
			reportPayloadUnmarshalFailure(rawPayload, err)

			res.WriteHeader(http.StatusNoContent)
		} else {
			listener(res, &payload, trace)
		}
	}
}

// OnIntegrationCreate is a method that registers a listener that fires when a user has connected to your webhook integration.
func (webhooks *Webhooks) OnIntegrationCreate(listener func(http.ResponseWriter, *IntegrationCreatePayload, string)) {
	webhooks.listeners["integration.create"] = newRawListener(listener)
}

// OnIntegrationDelete is a method that registers a listener that fires when a user has disconnected from your webhook integration.
func (webhooks *Webhooks) OnIntegrationDelete(listener func(http.ResponseWriter, *IntegrationDeletePayload, string)) {
	webhooks.listeners["integration.delete"] = newRawListener(listener)
}

// OnTest is a method that registers a listener that fires upon sent test from the project dashboard.
func (webhooks *Webhooks) OnTest(listener func(http.ResponseWriter, *TestPayload, string)) {
	webhooks.listeners["webhook.test"] = newRawListener(listener)
}

// OnVoteCreate is a method that registers a listener that fires when a user votes for your project.
func (webhooks *Webhooks) OnVoteCreate(listener func(http.ResponseWriter, *VoteCreatePayload, string)) {
	webhooks.listeners["vote.create"] = newRawListener(listener)
}

type rawPayload struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Handler is the handler function to be passed to http.HandleFunc.
func (webhooks *Webhooks) Handler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	signatureHeader := req.Header.Get("x-topgg-signature")
	trace := req.Header.Get("x-topgg-trace")

	if signatureHeader == "" || trace == "" {
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	var timestamp, signature string

	for pair := range strings.SplitSeq(signatureHeader, ",") {
		parts := strings.Split(pair, "=")

		if len(parts) != 2 {
			res.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		switch parts[0] {
		case "t":
			timestamp = parts[1]
		case APIVersion:
			signature = parts[1]
		}
	}

	if timestamp == "" || signature == "" {
		res.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	mac := hmac.New(sha256.New, []byte(webhooks.Secret))
	mac.Write(fmt.Appendf(nil, "%s.%s", timestamp, body))

	digest := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(digest), []byte(signature)) {
		res.WriteHeader(http.StatusForbidden)

		return
	}

	payload := &rawPayload{}

	if err = json.Unmarshal(body, payload); err != nil {
		reportPayloadUnmarshalFailure(body, err)

		res.WriteHeader(http.StatusNoContent)

		return
	}

	for event, listener := range webhooks.listeners {
		if event == payload.Type {
			listener(res, body, trace)

			return
		}
	}

	res.WriteHeader(http.StatusNoContent)
}
