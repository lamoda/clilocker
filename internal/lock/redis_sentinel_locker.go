package lock

import (
	"github.com/go-redis/redis/v7"
	"net/url"
	"time"
)

func NewRedisSentinelLocker(dsnString string) (Locker, error) {
	dsn, err := url.Parse(dsnString)
	if err != nil {
		return nil, err
	}

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: dsn.Host,
		SentinelAddrs: dsn.Query()["addrs"],
		DialTimeout: time.Duration(getQueryIntParam(dsn, "dial_timeout", 1)) * time.Second,
	})

	return NewRedisCommonLocker(client, dsn), nil
}