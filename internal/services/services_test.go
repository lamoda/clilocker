package services

import (
	"fmt"
	config2 "github.com/lamoda/clilocker/internal/config"
	"github.com/lamoda/clilocker/internal/lock"
	"reflect"
	"testing"
)

func TestServicesNew(t *testing.T) {
	cases := []struct {
		config   *config2.Config
		expected *Services
	}{
		{
			config: &config2.Config{
				LockerDsn:      "local://",
			},
			expected: &Services{
				Locker:      lock.NewLocalFileLocker(),
			},
		},
		{
			config: &config2.Config{
				LockerDsn:      "redis://test:1234",
			},
			expected: &Services{
				Locker: mustCreateLocker(func() (lock.Locker, error) {
					return lock.NewRedisLocker("redis://test:1234")
				}),
			},
		},
		{
			config: &config2.Config{
				LockerDsn:      "redis-sentinel://my_master?addrs=sentinel1&addrs=sentine2:1234",
			},
			expected: &Services{
				Locker: mustCreateLocker(func() (lock.Locker, error) {
					return lock.NewRedisSentinelLocker("redis-sentinel://my_master?addrs=sentinel1&addrs=sentine2:1234")
				}),
			},
		},
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("[%d]_%v", i, testCase.config), func(t *testing.T) {
			result, err := New(testCase.config)

			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if reflect.TypeOf(testCase.expected.Locker) != reflect.TypeOf(result.Locker) {
				t.Errorf(
					"expected node lock type %v\ngot %v",
					reflect.TypeOf(testCase.expected.Locker),
					reflect.TypeOf(result.Locker),
				)
			}
		})
	}
}

func mustCreateLocker(factory func() (lock.Locker, error)) lock.Locker {
	locker, err := factory()
	if err != nil {
		panic(err)
	}

	return locker
}
