package topgg

import "fmt"

// LargeWidget is a function that generates a large widget URL.
func LargeWidget(platform Platform, projectType ProjectType, id string) string {
	return fmt.Sprintf("%s/widgets/large/%s/%s/%s", BaseURL, string(platform), string(projectType), id)
}

// VotesWidget is a function that generates a small widget URL for displaying votes.
func VotesWidget(platform Platform, projectType ProjectType, id string) string {
	return fmt.Sprintf("%s/widgets/small/votes/%s/%s/%s", BaseURL, string(platform), string(projectType), id)
}

// OwnerWidget is a function that generates a small widget URL for displaying a project's owner.
func OwnerWidget(platform Platform, projectType ProjectType, id string) string {
	return fmt.Sprintf("%s/widgets/small/owner/%s/%s/%s", BaseURL, string(platform), string(projectType), id)
}

// SocialWidget is a function that generates a small widget URL for displaying social stats.
func SocialWidget(platform Platform, projectType ProjectType, id string) string {
	return fmt.Sprintf("%s/widgets/small/social/%s/%s/%s", BaseURL, string(platform), string(projectType), id)
}
