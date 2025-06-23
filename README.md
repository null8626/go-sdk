# Top.gg Go SDK

The community-maintained Go library for Top.gg.

## Installation

```sh
$ go get -u github.com/top-gg/go-dbl
```

## Setting up

### With defaults

```go
client, err := dbl.NewClient(os.Getenv("TOPGG_TOKEN"))

if err != nil {
  log.Fatalf("Error creating new Top.gg client: %s", err)
}
```

### With explicit options

```go
clientTimeout := 5 * time.Second
httpClient := &http.Client{}

client, err := dbl.NewClient(
  os.Getenv("TOPGG_TOKEN"),
  dbl.HTTPClientOption(httpClient),
  dbl.TimeoutOption(clientTimeout),
)

if err != nil {
  log.Fatalf("Error creating new Top.gg client: %s", err)
}
```

## Usage

### Getting a bot

```go
bot, err := client.GetBot("574652751745777665")

if err != nil {
  log.Fatalf("Unable to get a bot: %s", err)
}
```

### Getting several bots

```go
bots, err := client.GetBots(&GetBotsPayload{
  Limit: 20,
})

if err != nil {
  log.Fatalf("Unable to get several bots: %s", err)
}
```

### Getting your bot's voters

#### First page

```go
firstPageVoters, err := client.GetVoters(1)

if err != nil {
  log.Fatalf("Unable to get voters: %s", err)
}
```

#### Subsequent pages

```go
secondPageVoters, err := client.GetVoters(2)

if err != nil {
  log.Fatalf("Unable to get voters: %s", err)
}
```

### Check if a user has voted for your bot

```go
has_voted, err := client.HasUserVoted("661200758510977084")

if err != nil {
  log.Fatalf("Unable to check if a user has voted: %s", err)
}
```

### Getting your bot's server count

```go
serverCount, err := client.GetServerCount()

if err != nil {
  log.Fatalf("Unable to get server count: %s", err)
}
```

### Posting your bot's server count

```go
err := client.PostServerCount(bot.GetServerCount())

if err != nil {
  log.Fatalf("Unable to post server count: %s", err)
}
```

### Automatically posting your bot's server count every few minutes

```go
// Posts once every 30 minutes
autoposter, err := client.StartAutoposter(1800000, func() int {
  return bot.GetServerCount()
})

if err != nil {
  log.Fatalf("Unable to start autoposter: %s", err)
}

go func() {
  for {
    post_err := <-autoposter.Posted

    if post_err != nil {
      log.Fatalf("Unable to post server count: %s", post_err)
    }
  }
}()

// ...

autoposter.Stop() // Optional
```

### Checking if the weekend vote multiplier is active

```go
multiplierActive, err := client.IsMultiplierActive()

if err != nil {
  log.Fatalf("Unable to check weekend vote multiplier: %s", err)
}
```

### Generating widget URLs

#### Large

```go
widgetUrl := dbl.LargeWidget(dbl.DiscordBotWidget, "574652751745777665")
```

#### Votes

```go
widgetUrl := dbl.VotesWidget(dbl.DiscordBotWidget, "574652751745777665")
```

#### Owner

```go
widgetUrl := dbl.OwnerWidget(dbl.DiscordBotWidget, "574652751745777665")
```

#### Social

```go
widgetUrl := dbl.SocialWidget(dbl.DiscordBotWidget, "574652751745777665")
```

### Webhooks

#### Being notified whenever someone voted for your bot

```go
package main

import (
  "errors"
  "fmt"
  "log"
  "os"
  "net/http"

  "github.com/top-gg/go-dbl"
)

func main() {
  listener := dbl.NewWebhookListener(os.Getenv("MY_TOPGG_WEBHOOK_SECRET"), "/votes", handleVote)

  // Serve is a blocking call
  err := listener.Serve(":8080")

  if !errors.Is(err, http.ErrServerClosed) {
    log.Fatalf("HTTP server error: %s", err)
  }
}

func handleVote(payload []byte) {
  vote, err := dbl.NewWebhookVotePayload(payload)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: Unable to parse webhook payload: %s", err)
    return
  }

  fmt.Printf("A user with the ID of %s has voted us on Top.gg!", vote.VoterId)
}
```