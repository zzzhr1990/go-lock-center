package locker

// "github.com/go-redsync/redsync"
// "time"

// Config for redis config locker
type Config struct {
	// RedisAddress address for redis
	RedisAddress []string `yaml:"redis-address" json:"redisAddress,omitempty"`
	// RedisPassword passwd
	// RedisPassword string
}
