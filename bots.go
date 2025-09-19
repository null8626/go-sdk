package dbl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Reviews struct {
	// The project's average review score out of 5
	Score float64 `json:"averageScore"`

	// The project's review count
	Count int `json:"count"`
}

type Bot struct {
	// The Top.gg id of the bot
	ID string `json:"id"`

	// The Discord id of the bot
	ClientID string `json:"clientid"`

	// The username of the bot
	Username string `json:"username"`

	// The bot's avatar url
	Avatar string `json:"avatar"`

	// The prefix of the bot
	Prefix string `json:"prefix"`

	// The short description of the bot
	ShortDescription string `json:"shortdesc"`

	// The long description of the bot. Can contain HTML and/or Markdown (may be empty)
	LongDescription string `json:"longdesc"`

	// The tags of the bot
	Tags []string `json:"tags"`

	// The website url of the bot (may be empty)
	Website string `json:"website"`

	// The support server invite code of the bot (may be empty)
	Support string `json:"support"`

	// The link to the github repo of the bot (may be empty)
	Github string `json:"github"`

	// The owners of the bot. First one in the array is the main owner
	Owners []string `json:"owners"`

	// The custom bot invite url of the bot (may be empty)
	Invite string `json:"invite"`

	// The date when the bot was submitted
	Date time.Time `json:"date"`

	// The vanity url of the bot (may be empty)
	Vanity string `json:"vanity"`

	// The monthly amount of votes the bot has
	MonthlyPoints int `json:"monthlyPoints"`

	// The amount of votes the bot has
	Points int `json:"points"`

	// The amount of servers the bot is in
	ServerCount int `json:"server_count"`

	// The bot's reviews
	Review *Reviews `json:"reviews"`
}

type GetBotsPayload struct {
	// The amount of bots to return. Max. 500
	// Default 50
	Limit int

	// Amount of bots to skip
	// Default 0
	Offset int

	// [Deprecated since API v0] Field search filter
	Search map[string]string

	// The field to sort by descending, valid field names are "id", "date", and "monthlyPoints".
	Sort string

	// A list of fields to show
	Fields []string
}

type GetBotsResult struct {
	// Slice of Bot pointers of matching bots
	Results []*Bot `json:"results"`

	// The limit used
	Limit int `json:"limit"`

	// The offset used
	Offset int `json:"offset"`

	// The length of the results array
	Count int `json:"count"`

	// The total number of bots matching your search
	// Not limited by "limit" field
	Total int `json:"total"`
}

type BotStats struct {
	// The amount of servers the bot is in (may be empty)
	ServerCount int `json:"server_count"`

	// [Deprecated since API v0] The amount of servers the bot is in per shard. Always present but can be empty
	Shards []int `json:"shards"`

	// [Deprecated since API v0] The amount of shards a bot has (may be empty)
	ShardCount int `json:"shard_count"`
}

type checkResponse struct {
	Voted int `json:"voted"`
}

type BotStatsPayload struct {
	// The amount of servers the bot is in, must not be zero.
	ServerCount int `json:"server_count"`

	// [Deprecated since API v0] The amount of servers the bot is in per shard.
	Shards []int `json:"shards"`

	// [Deprecated since API v0] The zero-indexed id of the shard posting. Makes server_count set the shard specific server count (optional)
	ShardID int `json:"shard_id"`

	// [Deprecated since API v0] The amount of shards the bot has (optional)
	ShardCount int `json:"shard_count"`
}

// Information about different bots with an optional filter parameter
//
// Use nil if no option is passed
func (c *Client) GetBots(filter *GetBotsPayload) (*GetBotsResult, error) {
	req, err := c.createRequest("GET", "bots", nil)

	if err != nil {
		return nil, err
	} else if filter != nil {
		q := req.URL.Query()

		if filter.Limit > 500 {
			filter.Limit = 500
		}

		if filter.Limit > 0 {
			q.Add("limit", strconv.Itoa(filter.Limit))
		} else {
			q.Add("limit", "50")
		}

		if filter.Offset >= 0 {
			q.Add("offset", strconv.Itoa(filter.Offset))
		}

		if filter.Sort != "" {
			switch filter.Sort {
			case "id", "date", "monthlyPoints":
				q.Add("sort", filter.Sort)
			default:
				return nil, ErrInvalidRequest
			}
		}

		if len(filter.Fields) != 0 {
			q.Add("fields", strings.Join(filter.Fields, ","))
		}

		req.URL.RawQuery = q.Encode()
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := c.readBody(res)

	if err != nil {
		return nil, err
	}

	bots := &GetBotsResult{}

	err = json.Unmarshal(body, bots)

	if err != nil {
		return nil, err
	}

	return bots, nil
}

// Information about a specific bot
func (c *Client) GetBot(botID string) (*Bot, error) {
	req, err := c.createRequest("GET", "bots/"+botID, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := c.readBody(res)

	if err != nil {
		return nil, err
	}

	bot := &Bot{}

	err = json.Unmarshal(body, bot)

	if err != nil {
		return nil, err
	}

	return bot, nil
}

// Fetches your project's recent 100 unique voters
//
// # Requires authentication
func (c *Client) GetVotes(page int) ([]*User, error) {
	if page <= 0 {
		return nil, ErrInvalidRequest
	}

	req, err := c.createRequest("GET", fmt.Sprintf("bots/%s/votes", c.id), nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	q.Add("page", strconv.Itoa(page))

	req.URL.RawQuery = q.Encode()

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := c.readBody(res)

	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)

	err = json.Unmarshal(body, &users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// Use this endpoint to see who have upvoted for your project in the past 12 hours. It is safe to use this even if you have over 1k votes.
//
// # Requires authentication
//
// [Deprecated since API v0]: The _botID argument is no longer used.
func (c *Client) HasUserVoted(_botID, userID string) (bool, error) {
	req, err := c.createRequest("GET", "bots/check", nil)

	if err != nil {
		return false, err
	}

	q := req.URL.Query()

	q.Add("userId", userID)

	req.URL.RawQuery = q.Encode()

	res, err := c.httpClient.Do(req)

	if err != nil {
		return false, err
	}

	body, err := c.readBody(res)

	if err != nil {
		return false, err
	}

	cr := &checkResponse{}

	err = json.Unmarshal(body, cr)

	if err != nil {
		return false, err
	}

	return cr.Voted == 1, nil
}

// Information about a specific bot's stats
//
// [Deprecated since API v0]: The _botID argument is no longer used.
func (c *Client) GetBotStats(_botID string) (*BotStats, error) {
	req, err := c.createRequest("GET", "bots/stats", nil)

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := c.readBody(res)

	if err != nil {
		return nil, err
	}

	botStats := &BotStats{}

	err = json.Unmarshal(body, botStats)

	if err != nil {
		return nil, err
	}

	return botStats, nil
}

// Post your bot's server count
//
// # Requires authentication
//
// [Deprecated since API v0]: The _botID argument is no longer used.
func (c *Client) PostBotStats(_botID string, payload *BotStatsPayload) error {
	if payload.ServerCount <= 0 {
		return ErrInvalidRequest
	}

	encoded, err := json.Marshal(&BotStats{
		ServerCount: payload.ServerCount,
		Shards:      []int{},
	})

	if err != nil {
		return err
	}

	req, err := c.createRequest("POST", "bots/stats", bytes.NewBuffer(encoded))

	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return ErrRequestFailed
	}

	return nil
}
