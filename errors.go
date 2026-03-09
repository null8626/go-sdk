package dbl

import "errors"

var (
	ErrServerSide     = errors.New("Unable to make request due to Top.gg's end.")
	ErrNotFound       = errors.New("Such query does not exist.")
	ErrInvalidRequest = errors.New("Attempted to send an invalid request to the API.")
	ErrLocalRatelimit = errors.New("Temporarily prevented from sending requests by local ratelimiter.")
	ErrInvalidToken   = errors.New("Invalid API token.")
)
