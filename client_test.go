package dbl

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	const clientTimeout = 5 * time.Second

	httpClient := &http.Client{}

	_, err := NewClient(
		"token",
		HTTPClientOption(httpClient), // Sets an custom HTTP client (optional.)
		TimeoutOption(clientTimeout), // Sets an custom HTTP client timeout (optional.)
	)

	if err != nil {
		log.Fatalf("Unable to create new Top.gg client: %s", err)
	}

	// ...
}
