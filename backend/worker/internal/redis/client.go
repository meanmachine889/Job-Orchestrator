package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rds *redis.Client
}

func New(redisUrl string) *Client {
	opt, _ := redis.ParseURL(redisUrl)
	rds := redis.NewClient(opt)
	return &Client{rds: rds}
}

func (c *Client) WaitForJob(ctx context.Context) (string, error) {
	res, err := c.rds.BRPop(ctx, 0, "job_queue").Result()
	if err != nil {
		return "", err
	}
	return res[1], nil
}
