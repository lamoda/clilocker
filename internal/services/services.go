package services

import (
	"fmt"
	"github.com/lamoda/clilocker/internal/config"
	"github.com/lamoda/clilocker/internal/lock"
	"net/url"
)

type Services struct {
	Locker      lock.Locker
	ClusterLock lock.Locker
}

func New(conf *config.Config) (*Services, error) {
	services := &Services{}

	locker, err := createLocker(conf.LockerDsn)
	if err != nil {
		return nil, err
	}
	services.Locker = locker

	return services, nil
}

func createLocker(dsnString string) (lock.Locker, error) {
	dsn, err := url.Parse(dsnString)
	if err != nil {
		return nil, err
	}

	switch dsn.Scheme {
	case config.LocalFileScheme:
		return lock.NewLocalFileLocker(), nil
	case config.RedisScheme:
		return lock.NewRedisLocker(dsnString)
	case config.RedisSentinelScheme:
		return lock.NewRedisSentinelLocker(dsnString)
	default:
		return nil, fmt.Errorf("unknown type of locker %s", dsn.Scheme)
	}
}
