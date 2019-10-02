package locker

import (
	"time"

	"errors"

	"log"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

// Locker gloable locker
type Locker struct {
	sync *redsync.Redsync
}

// CreateNew new instance
func CreateNew(config *Config) (*Locker, error) {
	sts := []redsync.Pool{}

	if config.RedisAddress == nil || len(config.RedisAddress) == 0 {
		return nil, errors.New("Redis config is empty")
	}

	for _, ops := range config.RedisAddress {
		pool := redis.NewPool(func() (redis.Conn, error) {
			redi, err := redis.DialURL(ops)
			if err != nil {
				log.Printf("Cannot dial redis: %v", err)
			}

			return redi, err
		}, 3)
		connt := pool.Get()

		_, err := connt.Do("INFO")
		connt.Close()
		if err == nil {
			// log.Printf("Redis server response: %v", resp)
			sts = append(sts, pool)
		} else {
			log.Printf("Redis server error: %v", err)
		}

	}

	if len(sts) == 0 {
		return nil, errors.New("No redis server available")
	}

	// sts.
	/*
		pool, err :=  &redis.Pool{
			MaxIdle: 3,
			IdleTimeout: 240 * time.Second,
			// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
			Dial: func () (redis.Conn, error) { return redis.Dial("tcp", config.RedisAddress) },
		  }
	*/
	// sts = redis.NewPool()
	locker := &Locker{
		sync: redsync.New(sts),
	}
	// redsync.N

	return locker, nil
}

// LockForKeyWithRetry lock for special key
func (l *Locker) LockForKeyWithRetry(key string, expiry time.Duration, retry int) (*redsync.Mutex, error) {
	// opt := redsync.SetTries(5)
	// redsync.SetGenValueFunc()
	mx := l.sync.NewMutex(key, redsync.SetExpiry(expiry), redsync.SetTries(retry), redsync.SetRetryDelay(time.Second*5))
	return mx, mx.Lock()
}

// LockForKeyWithNoRetry lock for special key
func (l *Locker) LockForKeyWithNoRetry(key string, expiry time.Duration) (*redsync.Mutex, error) {
	opt := redsync.SetTries(1)
	// redsync.SetGenValueFunc()
	mx := l.sync.NewMutex(key, redsync.SetExpiry(expiry), opt)
	return mx, mx.Lock()
}
