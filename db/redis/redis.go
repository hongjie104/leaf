package redis

import (
	"time"

	red "github.com/garyburd/redigo/redis"
)

var (
	redisClient *red.Pool
)

// DialContext DialContext
type DialContext struct {
}

// Dial 连接
func Dial(host string, db int) {
	// 从配置文件获取redis的ip以及db
	// 建立连接池
	redisClient = &red.Pool{
		// 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxIdle: 1,
		// 最大的激活连接数，表示同时最多有N个连接
		MaxActive: 10,
		// 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: 180 * time.Second,
		Wait:        true,
		Dial: func() (red.Conn, error) {
			c, err := red.Dial(
				"tcp",
				host,
				red.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				red.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				red.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				red.DialDatabase(db),
			)
			if err != nil {
				return nil, err
			}
			// // 选择db
			// c.Do("SELECT", db)
			return c, nil
		},
	}
}

// Ref Ref
func (c *DialContext) Ref() red.Conn {
	return redisClient.Get()
}
