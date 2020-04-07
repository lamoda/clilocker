package lock

import "testing"

func TestSentinelDsnParsing(t *testing.T) {
	locker, err := NewRedisSentinelLocker(
		"redis-sentinel://my_master?addrs=test:1234&addrs=test:4567&key_prefix=test",
	)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	redisLocker, ok := locker.(*RedisCommonLocker)
	if !ok {
		t.Fatal("Expected locker to be redis based")
	}

	if redisLocker.keyPrefix != "test" {
		t.Errorf("keyPrefix expected to be %s, got %s", "test", redisLocker.keyPrefix)
	}
}
