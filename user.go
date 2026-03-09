package dbl

import "time"

// A project's type.
type UserSource string

const (
	UserDiscord UserSource = "discord"
	UserTopgg   UserSource = "topgg"
)

// A Top.gg user.
type User struct {
	Id         string `json:"id"`          // The user's ID.
	Name       string `json:"name"`        // The user's name.
	Avatar     string `json:"avatar_url"`  // The user's avatar URL.
	PlatformId string `json:"platform_id"` // The user's platform ID.
}

// A project's vote information.
type Vote struct {
	PartialVote
	VoterId    string `json:"user_id"`     // The voter's ID.
	PlatformId string `json:"platform_id"` // The voter's ID on the project's platform.
}

// A brief information of a project's vote.
type PartialVote struct {
	VotedAt   time.Time `json:"created_at"` // When the vote was cast.
	ExpiresAt time.Time `json:"expires_at"` // When the vote expires and the user is required to vote again.
	Weight    int       `json:"weight"`     // The number of votes this vote counted for. This is a rounded integer value which determines how many points this individual vote was worth.
}
