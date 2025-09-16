package dbl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoposter(t *testing.T) {
	client, err := NewClient(os.Getenv("TOPGG_TOKEN"))

	assert.Nil(t, err, "Client should be created w/o error")

	autoposter, err := client.StartAutoposter(900000, func() int {
		return 2
	})

	assert.Nil(t, err, "Autoposter should be created w/o error")

	for i := 0; i < 3; i++ {
		err := <-autoposter.Posted
		assert.Nil(t, err, "Posting should not error")
	}

	autoposter.Stop()
}
