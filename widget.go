package dbl

import "fmt"

func asWidgetPath(projectType ProjectType) string {
	switch projectType {
	case ProjectDiscordBot:
		return "discord/bot"
	case ProjectDiscordServer:
		return "discord/server"
	default:
		return ""
	}
}

// LargeWidget is a function that generates a large widget URL.
func LargeWidget(projectType ProjectType, id string) string {
	return fmt.Sprintf(BaseURL+"/widgets/large/%s/%s", asWidgetPath(projectType), id)
}

// VotesWidget is a function that generates a small widget URL for displaying votes.
func VotesWidget(projectType ProjectType, id string) string {
	return fmt.Sprintf(BaseURL+"/widgets/small/votes/%s/%s", asWidgetPath(projectType), id)
}

// OwnerWidget is a function that generates a small widget URL for displaying a project's owner.
func OwnerWidget(projectType ProjectType, id string) string {
	return fmt.Sprintf(BaseURL+"/widgets/small/owner/%s/%s", asWidgetPath(projectType), id)
}

// SocialWidget is a function that generates a small widget URL for displaying social stats.
func SocialWidget(projectType ProjectType, id string) string {
	return fmt.Sprintf(BaseURL+"/widgets/small/social/%s/%s", asWidgetPath(projectType), id)
}
