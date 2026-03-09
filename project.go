package dbl

// A project's platform.
type Platform string

const (
	Discord Platform = "discord"
)

// A project's type.
type ProjectType string

const (
	DiscordBotProject    ProjectType = "bot"
	DiscordServerProject ProjectType = "server"
)

// A project listed on Top.gg.
type Project struct {
	Id           string      `json:"id"`           // The project's ID.
	Name         string      `json:"name"`         // The project's name sourced from the external platform.
	Platform     Platform    `json:"platform"`     // The project's platform.
	Type         ProjectType `json:"type"`         // The project's type.
	Headline     string      `json:"headline"`     // The project's short description.
	Tags         []string    `json:"tags"`         // The project's tag IDs.
	CurrentVotes int         `json:"votes"`        // The project's current vote count that affects the project's ranking.
	TotalVotes   int         `json:"votes_total"`  // The project's total vote count.
	ReviewScore  float32     `json:"review_score"` // The project's review score out of 5.
	ReviewCount  int         `json:"review_count"` // The project's total review count.
}

// A brief information on project listed on Top.gg.
type PartialProject struct {
	Id         string      `json:"id"`          // The project's ID.
	Type       ProjectType `json:"type"`        // The project's ID.
	Platform   Platform    `json:"platform"`    // The project's platform.
	PlatformId string      `json:"platform_id"` // The project's platform ID.
}
