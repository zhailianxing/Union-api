package share

import (
	"fmt"
	log "github.com/alecthomas/log4go"
	"tantanbei.com/redis"
	"time"
)

var Redis *redis.Redis
var err error

func init() {
	log.Debug("create redis")
	Redis, err = redis.NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	if err != nil {
		panic(fmt.Sprint("create redis error:", err))
	}
}
