package dbl

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	const clientTimeout = 5 * time.Second

	httpClient := &http.Client{}

	client, err := NewClient(
		os.Getenv("apikey"),
		HTTPClientOption(httpClient), // Setting a custom HTTP client. Default is *http.Client with default timeout.
		TimeoutOption(clientTimeout), // Setting timeout option. Default is 3 seconds
	)

	if err != nil {
		log.Fatalf("Error creating new Discord Bot List client: %s", err)
	}

	_, err = client.GetBotStats()

	assert.Nil(t, err, "GetBotStats should succeed")

	err = client.PostBotStats(&BotStats{
		ServerCount: 2,
	})

	assert.Nil(t, err, "PostBotStats should succeed")

	time.Sleep(1 * time.Second)
	_, err = client.GetBot("264811613708746752")

	assert.Nil(t, err, "GetBot should succeed")

	getBotsPayload := GetBotsPayload{
		Limit:  50,
		Offset: 0,
		Sort:   "id",
	}

	time.Sleep(1 * time.Second)
	_, err = client.GetBots(&getBotsPayload)

	assert.Nil(t, err, "GetBots should succeed")

	time.Sleep(1 * time.Second)
	_, err = client.GetVotes(1)

	assert.Nil(t, err, "GetVotes should succeed")

	time.Sleep(1 * time.Second)
	_, err = client.HasUserVoted("661200758510977084")

	assert.Nil(t, err, "HasUserVoted should succeed")

	time.Sleep(1 * time.Second)
	_, err = client.IsMultiplierActive()

	assert.Nil(t, err, "IsMultiplierActive should succeed")
}
