package dbl

import (
	"flag"
	"time"
)

type BotAutoposter struct {
	stopChannel chan bool
	// The channel on which errors are delivered after every attempted post API request.
	Posted chan error
}

type BotAutoposterCallback func() int

// Automates your bot's server count posting
//
// # Requires authentication
func (c *Client) StartBotAutoposter(delay int, callback BotAutoposterCallback) (*BotAutoposter, error) {
	if c.token == "" {
		return nil, ErrRequireAuthentication
	}

	if flag.Lookup("test.v") == nil && delay < 900000 {
		delay = 900000
	}

	stopChannel := make(chan bool)
	postedChannel := make(chan error)

	go func() {
		ticker := time.NewTicker(time.Duration(delay) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopChannel:
				close(postedChannel)
				return
			case <-ticker.C:
				postedChannel <- c.PostBotServerCount(callback())
			}
		}
	}()

	return &BotAutoposter{
		stopChannel: stopChannel,
		Posted:      postedChannel,
	}, nil
}

// Stops the bot autoposter
func (a *BotAutoposter) Stop() {
	a.stopChannel <- true
}
