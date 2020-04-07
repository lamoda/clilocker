package lock

import (
	"github.com/go-redis/redis/v7"
	"net/url"
	"time"
)

func NewRedisLocker(dsnString string) (Locker, error) {
	dsn, err := url.Parse(dsnString)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Network:     "tcp",
		Addr:        dsn.Host,
		DialTimeout: time.Duration(getQueryIntParam(dsn, "dial_timeout", 1)) * time.Second,
	})

	return NewRedisCommonLocker(client, dsn), nil
}


