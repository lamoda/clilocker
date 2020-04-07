package lock

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	defaultRedisPrefix  = "locker_redis"
	lockTimeout         = 1000 * time.Millisecond
	lockRefreshInterval = 500 * time.Millisecond
)

type RedisCommonLocker struct {
	locker    *redislock.Client
	keyPrefix string
	locks     map[string]*innerLock
}

func NewRedisCommonLocker(client *redis.Client, dsn *url.URL) Locker {
	keyPrefix := dsn.Query().Get("key_prefix")
	if keyPrefix == "" {
		keyPrefix = defaultRedisPrefix
	}

	return &RedisCommonLocker{
		locker:    redislock.New(client),
		keyPrefix: keyPrefix,
		locks:     make(map[string]*innerLock, 0),
	}
}

func (l *RedisCommonLocker) Lock(name string, quantity int) (*Lock, error) {
	for i := 0; i < quantity; i++ {
		lockId := createLockId(name, i)

		innerLock, err := newInnerLock(l.locker, lockId)
		if err != nil {
			return nil, err
		}
		if innerLock == nil {
			continue
		}

		l.locks[lockId] = innerLock

		lock := NewTakenLock(lockId, l)
		return lock, nil
	}

	return NewUntakenLock(), nil
}

func (l *RedisCommonLocker) Release(lock Lock) error {
	innerLock, ok := l.locks[lock.Id()]
	if !ok {
		return fmt.Errorf("can not find lock with id %s", lock.Id())
	}

	return innerLock.release()
}

func (l *RedisCommonLocker) createKeyNameFromId(lockId string) string {
	return l.keyPrefix + hex.EncodeToString(md5.New().Sum([]byte(lockId)))
}

type innerLock struct {
	ticker  *time.Ticker
	lock    *redislock.Lock
	mux     sync.Mutex
	running bool
}

func newInnerLock(locker *redislock.Client, key string) (*innerLock, error) {
	lock, err := locker.Obtain(key, lockTimeout, nil)
	if err == redislock.ErrNotObtained {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	value := &innerLock{
		ticker:  time.NewTicker(lockRefreshInterval),
		lock:    lock,
		running: true,
	}
	value.refresh()
	return value, nil
}

func (p innerLock) refresh() {
	go func() {
		for range p.ticker.C {
			p.mux.Lock()
			if p.running {
				err := p.lock.Refresh(lockTimeout, nil)
				if err != nil {
					p.ticker.Stop()
				}
			}
			p.mux.Unlock()
		}
	}()
}

func (p innerLock) release() error {
	p.mux.Lock()
	p.running = false
	p.ticker.Stop()

	err := p.lock.Release()
	p.mux.Unlock()
	return err
}

func getQueryIntParam(dsn *url.URL, param string, defaultValue int) int {
	value, err := strconv.Atoi(dsn.Query().Get(param))
	if err != nil {
		return defaultValue
	}

	return value
}
