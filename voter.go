package dbl

type Voter struct {
	// The id of the voter
	ID string `json:"id"`

	// The username of the voter
	Username string `json:"username"`

	// The voter's avatar url
	Avatar string `json:"avatar"`
}
