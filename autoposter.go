package dbl

import (
	"flag"
	"time"
)

type Autoposter struct {
	stopChannel chan bool
	// The channel on which errors are delivered after every attempted post API request.
	Posted chan error
}

type AutoposterCallback func() *BotStatsPayload

// Automates your bot's server count posting
//
// # Requires authentication
func (c *Client) StartAutoposter(delay int, callback AutoposterCallback) (*Autoposter, error) {
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
				postedChannel <- c.PostBotStats("", callback())
			}
		}
	}()

	return &Autoposter{
		stopChannel: stopChannel,
		Posted:      postedChannel,
	}, nil
}

// Stops the autoposter
func (a *Autoposter) Stop() {
	a.stopChannel <- true
}
