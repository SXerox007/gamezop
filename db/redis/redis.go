package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	MAX_IDLE   = 80
	MAX_ACTIVE = 12000
	PORT       = 6379 // default port
)

var redisDBPool *redis.Pool

func Init() {
	redisDBPool = &redis.Pool{
		MaxIdle:   MAX_IDLE,
		MaxActive: MAX_ACTIVE,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			// ping redis test the connectivity
			// ping(c)
			// fmt.Println("Redis connected with success.")
			return c, err
		},
	}
}

// ping tests connectivity for redis
func ping(c redis.Conn) error {
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}
	fmt.Printf("PING Response = %s\n", s)
	return nil
}

func GetClient() redis.Conn {
	return redisDBPool.Get()
}
