//this file sets up a Redis-backed queue for job processing

package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	client *redis.Client
}

func New(redisUrl string) *Queue {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}
	return &Queue{
		client: redis.NewClient(opt),
	}
}

func (q *Queue) Enqueue(ctx context.Context, jobId string) error {
	return q.client.LPush(ctx, "job_queue", jobId).Err()
}
