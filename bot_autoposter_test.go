package dbl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBotAutoposter(t *testing.T) {
	client, err := NewClient(os.Getenv("TOPGG_TOKEN"))

	assert.Nil(t, err, "Client should be created w/o error")

	botAutoposter, err := client.StartBotAutoposter(3000, func() int {
		return 2
	})

	assert.Nil(t, err, "BotAutoposter should be created w/o error")

	for i := 0; i < 3; i++ {
		err := <-botAutoposter.Posted
		assert.Nil(t, err, "Posting should not error")
	}

	botAutoposter.Stop()
}
