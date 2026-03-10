package dbl

import "time"

// IntegrationCreatePayload represents an 'integration.create' webhook.payload. Fires when a user has connected to your webhook integration.
type IntegrationCreatePayload struct {
	ConnectionId string         `json:"connection_id"`  // The unique identifier for this connection.
	Secret       string         `json:"webhook_secret"` // The secret used to verify future webhook deliveries.
	Project      PartialProject `json:"project"`        // The project that the integration refers to.
	User         User           `json:"user"`           // The user who triggered this event.
}

// IntegrationDeletePayload represents an 'integration.delete' webhook.payload. Fires when a user has disconnected from your webhook integration.
type IntegrationDeletePayload struct {
	ConnectionId string `json:"connection_id"` // The unique identifier for this connection.
}

// TestPayload represents a 'webhook.test' webhook.payload. Fires upon sent test from the project dashboard.
type TestPayload struct {
	Project PartialProject `json:"project"` // The project that the test refers to.
	User    User           `json:"user"`    // The user who triggered this test.
}

// VoteCreatePayload represents a 'vote.create' webhook.payload. Fires when a user votes for your project.
type VoteCreatePayload struct {
	Id        string         `json:"id"`         // The vote's ID.
	Weight    int            `json:"weight"`     // The number of votes this vote counted for. This is a rounded integer value which determines how many points this individual vote was worth.
	VotedAt   time.Time      `json:"created_at"` // When the vote was cast.
	ExpiresAt time.Time      `json:"expires_at"` // When the vote expires and the user is required to vote again.
	Project   PartialProject `json:"project"`    // The project that received this vote.
	User      User           `json:"user"`       // The user who voted for this project.
}

// Payload represents all possible webhook payloads.
type Payload interface {
	IntegrationCreatePayload | IntegrationDeletePayload | TestPayload | VoteCreatePayload
}
