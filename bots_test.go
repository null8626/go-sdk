package dbl

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBotID  = "574652751745777665"
	fetchLimit = 20
)

func TestBots(t *testing.T) {
	client, err := NewClient(os.Getenv("TOPGG_TOKEN"))

	assert.Nil(t, err, "Client should be created w/o error")

	bots, err := client.GetBots(&GetBotsPayload{
		Limit: fetchLimit,
	})

	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Request should be successful (API depended)")

	assert.Equal(t, fetchLimit, len(bots.Results), "Results array size should match request limit")
	assert.Equal(t, fetchLimit, bots.Count, "Count should match request limit")
	assert.Equal(t, 0, bots.Offset, "Offset should be zero or non-specified")
}

func TestBot(t *testing.T) {
	client, err := NewClient(os.Getenv("TOPGG_TOKEN"))

	assert.Nil(t, err, "Client should be created w/o error")

	bot, err := client.GetBot(testBotID)

	assert.Nil(t, err, "Unable to get user data")

	assert.Equal(t, testBotID, bot.ID, "Request & result bot ID should match")
}
