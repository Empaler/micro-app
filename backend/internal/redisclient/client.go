package redisclient

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type PopularItem struct {
	ID    int64   `json:"id"`
	Score float64 `json:"score"`
}

type Client struct {
	client *redis.Client
}

func NewClient(addr string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func NewClientFromConfig(host, port, db string) (*Client, error) {
	parsedDB, err := strconv.Atoi(db)
	if err != nil {
		return nil, fmt.Errorf("invalid redis db: %w", err)
	}

	addr := host
	if port != "" {
		addr = fmt.Sprintf("%s:%s", host, port)
	}

	return NewClient(addr, parsedDB)
}

func (c *Client) IncrementLookups(ctx context.Context, collection string, itemID int64) error {
	if c == nil || c.client == nil {
		return nil
	}

	key := fmt.Sprintf("%s:popular", strings.TrimSpace(collection))
	return c.client.ZIncrBy(ctx, key, 1, strconv.FormatInt(itemID, 10)).Err()
}

func (c *Client) GetTopLookedUp(ctx context.Context, collection string, limit int) ([]PopularItem, error) {
	if c == nil || c.client == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 10
	}

	key := fmt.Sprintf("%s:popular", strings.TrimSpace(collection))
	result, err := c.client.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	items := make([]PopularItem, 0, len(result))
	for _, member := range result {
		id, err := strconv.ParseInt(member.Member.(string), 10, 64)
		if err != nil {
			continue
		}
		items = append(items, PopularItem{ID: id, Score: member.Score})
	}

	return items, nil
}
