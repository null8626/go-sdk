package dbl

type WidgetType string

const (
	DiscordBotWidget    WidgetType = "discord/bot"
	DiscordServerWidget WidgetType = "discord/server"
)

// Generates a large widget URL.
func LargeWidget(ty WidgetType, id string) string {
	return BaseURL + "widgets/large/" + string(ty) + "/" + id
}

// Generates a small widget URL for displaying votes.
func VotesWidget(ty WidgetType, id string) string {
	return BaseURL + "widgets/small/votes/" + string(ty) + "/" + id
}

// Generates a small widget URL for displaying an entity's owner.
func OwnerWidget(ty WidgetType, id string) string {
	return BaseURL + "widgets/small/owner/" + string(ty) + "/" + id
}

// Generates a small widget URL for displaying social stats.
func SocialWidget(ty WidgetType, id string) string {
	return BaseURL + "widgets/small/social/" + string(ty) + "/" + id
}
