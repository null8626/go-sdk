# Top.gg Go SDK

> For more information, see the documentation here: <https://pkg.go.dev/github.com/top-gg/go-dbl>.

The community-maintained Go library for Top.gg.

## Chapters

- [Installation](#installation)
- [Setting up](#setting-up)
- [Usage](#usage)
  - [Getting your project's information](#getting-your-projects-information)
  - [Getting your project's vote information of a user](#getting-your-projects-vote-information-of-a-user)
  - [Getting a cursor-based paginated list of votes for your project](#getting-a-cursor-based-paginated-list-of-votes-for-your-project)
  - [Posting your bot's application commands list](#posting-your-bots-application-commands-list)
  - [Generating widget URLs](#generating-widget-urls)
  - [Webhooks](#webhooks)

## Installation

```sh
$ go get github.com/top-gg/go-dbl
```

## Setting up

```go
import "github.com/top-gg/go-dbl"

client, err := dbl.NewClient(os.Getenv("TOPGG_TOKEN"))
```

## Usage

### Getting your project's information

```go
project, err := client.GetSelf()
```

### Getting your project's vote information of a user

#### Discord ID

```go
vote, err := client.GetVote(dbl.UserDiscord, "661200758510977084")
```

#### Top.gg ID

```go
vote, err := client.GetVote(dbl.UserTopgg, "8226924471638491136")
```

### Getting a cursor-based paginated list of votes for your project

```go
since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

firstPage, err := client.GetVotes(since)

if err != nil {
  return err
}

for _, vote := range firstPage.Votes() {
  // ...
}

secondPage, err := firstPage.Next()

if err != nil {
  return err
}

for _, vote := range secondPage.Votes() {
  // ...
}
```

### Posting your bot's application commands list

#### Disgo

In your bot's ready event listener:

```go
bot.WithEventListenerFunc(func(e *events.Ready) {
	commands, err := e.Client().Rest.GetGlobalCommands(e.Client().ApplicationID, true)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: Unable to get global commands: %s\n", err.Error())

		return
	}

	rawCommands := make([]string, len(commands))

	for i, command := range commands {
		rawCommand, err := command.MarshalJSON()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: Unable to serialize global command #%d: %s\n", i+1, err.Error())

			return
		}

		rawCommands[i] = string(rawCommand)
	}

	if err = client.PostCommands("[" + strings.Join(rawCommands, ",") + "]"); err != nil {
		fmt.Fprintf(os.Stderr, "error: Unable to post commands to Top.gg: %s\n", err.Error())
	}
})
```

#### Raw

```go
// Array of application commands that
// can be serialized to Discord API's raw JSON format.
err := client.PostCommands(`[{
  "id": "1",
  "type": 1,
  "application_id": "1",
  "name": "test",
  "description": "command description",
  "default_member_permissions": "",
  "version": "1"
}]`)
```

### Generating widget URLs

#### Large

```go
widgetUrl := dbl.LargeWidget(dbl.ProjectDiscordBot, "1026525568344264724")
```

#### Votes

```go
widgetUrl := dbl.VotesWidget(dbl.ProjectDiscordBot, "1026525568344264724")
```

#### Owner

```go
widgetUrl := dbl.OwnerWidget(dbl.ProjectDiscordBot, "1026525568344264724")
```

#### Social

```go
widgetUrl := dbl.SocialWidget(dbl.ProjectDiscordBot, "1026525568344264724")
```

### Webhooks

```go
webhooks := dbl.NewWebhooks(os.Getenv("TOPGG_WEBHOOK_SECRET"))

// Optional
webhooks.OnIntegrationCreate(func(res http.ResponseWriter, payload *dbl.IntegrationCreatePayload, trace string) {
	res.WriteHeader(http.StatusNoContent)
})

// Optional
webhooks.OnIntegrationDelete(func(res http.ResponseWriter, payload *dbl.IntegrationDeletePayload, trace string) {
	res.WriteHeader(http.StatusNoContent)
})

// Optional
webhooks.OnTest(func(res http.ResponseWriter, payload *dbl.TestPayload, trace string) {
	res.WriteHeader(http.StatusNoContent)
})

// Optional
webhooks.OnVoteCreate(func(res http.ResponseWriter, payload *dbl.VoteCreatePayload, trace string) {
	res.WriteHeader(http.StatusNoContent)
})
```

Later, in your server's setup:

```go
// POST /webhook
http.HandleFunc("/webhook", webhooks.Handler)
```
