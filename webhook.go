package dbl

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type WebhookListenerFunc func([]byte)

type WebhookListener struct {
	token   string
	path    string
	handler WebhookListenerFunc
	mux     *http.ServeMux
}

type WebhookVotePayload struct {
	// ID of the bot/server that received a vote
	ReceiverId string

	// ID of the user who voted
	VoterId string

	// Whether this vote is just a test done from the page settings
	IsTest bool

	// Whether the weekend multiplier is in effect, meaning users votes count as two
	IsWeekend bool

	// Query string params found on the /bot/:ID/vote page. Example: ?a=1&b=2
	Query url.Values
}

type wVotePayload struct {
	Bot       *string `json:"bot"`
	Server    *string `json:"guild"`
	User      string  `json:"user"`
	Type      string  `json:"type"`
	IsWeekend *bool   `json:"isWeekend"`
	Query     string  `json:"query"`
}

func NewWebhookVotePayload(data []byte) (*WebhookVotePayload, error) {
	p := &wVotePayload{}

	if err := json.Unmarshal(data, p); err != nil {
		return nil, err
	}

	m, err := url.ParseQuery(p.Query)

	if err != nil {
		return nil, err
	}

	receiverId := p.Bot

	if receiverId == nil {
		receiverId = p.Server
	}

	isWeekend := false

	if p.IsWeekend != nil {
		isWeekend = *p.IsWeekend
	}

	return &WebhookVotePayload{
		ReceiverId: *receiverId,
		VoterId:    p.User,
		IsTest:     p.Type == "test",
		IsWeekend:  isWeekend,
		Query:      m,
	}, nil
}

// Create a new webhook listener
func NewWebhookListener(token string, path string, handler func([]byte)) *WebhookListener {
	return &WebhookListener{
		token:   token,
		path:    path,
		handler: WebhookListenerFunc(handler),
	}
}

// Starts listening on specific address. A Blocking call.
// Returns non-nil error from ListenAndServe
func (wl *WebhookListener) Serve(addr string) error {
	wl.mux = http.NewServeMux()

	wl.mux.HandleFunc(wl.path, wl.handleRequest)

	return http.ListenAndServe(addr, wl.mux)
}

func (wl *WebhookListener) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	if r.Header.Get("Authorization") != wl.token {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)

	wl.handler(body)
}
