package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type MyRedis struct {
	db *redis.Client
}

const (
	Step = 1000
)

var GRedis *MyRedis

func InitRedis() error {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	GRedis = &MyRedis{
		db: db,
	}

	return nil
}

func (myRedis *MyRedis) NextStep(business string) (int64, int64, error) {
	// 总耗时小于2秒
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(20000)*time.Millisecond)
	defer cancelFunc()

	for counter := 0; counter < NewSegmentRetryTimes; counter++ {
		currentId, err := myRedis.db.IncrBy(ctx, business, Step).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				currentId = Step
			} else {
				println(fmt.Sprintf("err:+%v", err))
				continue
			}
		}

		// 执行成功
		return currentId, Step, nil
	}

	println("failed")
	return 0, 0, errors.New("new segment error")
}
