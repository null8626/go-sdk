package dbl

type User struct {
	// The id of the user
	ID string `json:"id"`

	// The username of the user
	Username string `json:"username"`

	// [Deprecated since API v0]: The discriminator of the user
	Discriminator string `json:"discriminator"`

	// The avatar hash of the user's avatar (may be empty)
	Avatar string `json:"avatar"`

	// [Deprecated since API v0]: The cdn hash of the user's avatar if the user has none
	DefAvatar string `json:"defAvatar"`

	// [Deprecated since API v0]: The bio of the user
	Biography string `json:"bio"`

	// [Deprecated since API v0]: The banner image url of the user (may be empty)
	Banner string `json:"banner"`

	// [Deprecated since API v0]: The user's socials
	Social *Social `json:"social"`

	// [Deprecated since API v0]: The custom hex color of the user (may be empty)
	Color string `json:"color"`

	// [Deprecated since API v0]: The supporter status of the user
	Supporter bool `json:"supporter"`

	// [Deprecated since API v0]: The certified status of the user
	CertifiedDeveloper bool `json:"certifiedDev"`

	// [Deprecated since API v0]: The mod status of the user
	Moderator bool `json:"mod"`

	// [Deprecated since API v0]: The website moderator status of the user
	WebsiteModerator bool `json:"webMod"`

	// [Deprecated since API v0]: The admin status of the user
	Admin bool `json:"admin"`
}

type Social struct {
	// [Deprecated since API v0]: The youtube channel id of the user (may be empty)
	Youtube string `json:"youtube"`

	// [Deprecated since API v0]: The reddit username of the user (may be empty)
	Reddit string `json:"reddit"`

	// [Deprecated since API v0]:	The twitter username of the user (may be empty)
	Twitter string `json:"twitter"`

	// [Deprecated since API v0]: The instagram username of the user (may be empty)
	Instagram string `json:"instagram"`

	// [Deprecated since API v0]: The github username of the user (may be empty)
	Github string `json:"github"`
}
