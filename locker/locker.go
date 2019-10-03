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
			sts = append(sts, pool)
		} else {
			log.Printf("Redis server error: %v", err)
		}

	}

	if len(sts) == 0 {
		return nil, errors.New("No redis server available")
	}
	locker := &Locker{
		sync: redsync.New(sts),
	}

	return locker, nil
}

// LockForKey lock for special key
func (l *Locker) LockForKey(key string, expiry time.Duration) (*redsync.Mutex, error) {
	// opt := redsync.SetTries(5)
	// redsync.SetGenValueFunc()
	// mx := l.sync.NewMutex(key, redsync.SetExpiry(expiry), redsync.SetTries(retry), redsync.SetRetryDelay(time.Second*5))
	return l.LockForKeyWithRetry(key, expiry, 5)
}

// LockForKeyWithRetry lock for special key
func (l *Locker) LockForKeyWithRetry(key string, expiry time.Duration, retry int) (*redsync.Mutex, error) {
	// opt := redsync.SetTries(5)
	// redsync.SetGenValueFunc()
	// mx := l.sync.NewMutex(key, redsync.SetExpiry(expiry), redsync.SetTries(retry), redsync.SetRetryDelay(time.Second*5))
	return l.LockForKeyWithRetryDelay(key, expiry, retry, time.Second*2)
}

// LockForKeyWithRetryDelay lock for special key
func (l *Locker) LockForKeyWithRetryDelay(key string, expiry time.Duration, retry int, retryDelay time.Duration) (*redsync.Mutex, error) {
	// opt := redsync.SetTries(5)
	// redsync.SetGenValueFunc()
	mx := l.sync.NewMutex(key, redsync.SetExpiry(expiry), redsync.SetTries(retry), redsync.SetRetryDelay(retryDelay))
	return mx, mx.Lock()
}

// LockForKeyWithNoRetry lock for special key
func (l *Locker) LockForKeyWithNoRetry(key string, expiry time.Duration) (*redsync.Mutex, error) {
	opt := redsync.SetTries(1)
	// redsync.SetGenValueFunc()
	mx := l.sync.NewMutex(key, redsync.SetExpiry(expiry), opt)
	return mx, mx.Lock()
}
