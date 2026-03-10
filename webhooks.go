package dbl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type rawListener = func(http.ResponseWriter, json.RawMessage, string)

// A Top.gg webhook manager.
type Webhooks struct {
	Secret    string
	listeners map[string]rawListener
}

// Creates a new webhook manager instance.
func NewWebhooks(Secret string) *Webhooks {
	return &Webhooks{
		Secret:    Secret,
		listeners: make(map[string]rawListener),
	}
}

func newRawListener[P Payload](listener func(http.ResponseWriter, *P, string)) rawListener {
	return func(res http.ResponseWriter, rawPayload json.RawMessage, trace string) {
		var payload P

		if err := json.Unmarshal(rawPayload, &payload); err != nil {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			listener(res, &payload, trace)
		}
	}
}

// Registers a listener that fires when a user has connected to your webhook integration.
func (webhooks *Webhooks) OnIntegrationCreate(listener func(http.ResponseWriter, *IntegrationCreatePayload, string)) {
	webhooks.listeners["integration.create"] = newRawListener(listener)
}

// Registers a listener that fires when a user has disconnected from your webhook integration.
func (webhooks *Webhooks) OnIntegrationDelete(listener func(http.ResponseWriter, *IntegrationDeletePayload, string)) {
	webhooks.listeners["integration.delete"] = newRawListener(listener)
}

// Registers a listener that fires upon sent test from the project dashboard.
func (webhooks *Webhooks) OnTest(listener func(http.ResponseWriter, *TestPayload, string)) {
	webhooks.listeners["webhook.test"] = newRawListener(listener)
}

// Registers a listener that fires when a user votes for your project.
func (webhooks *Webhooks) OnVoteCreate(listener func(http.ResponseWriter, *VoteCreatePayload, string)) {
	webhooks.listeners["vote.create"] = newRawListener(listener)
}

type rawPayload struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// The handler function to be passed to HandleFunc.
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
		res.WriteHeader(http.StatusBadRequest)

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
