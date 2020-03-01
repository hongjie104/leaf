package redis

import (
	"time"

	red "github.com/garyburd/redigo/redis"
)

// DialContext DialContext
type DialContext struct {
	redisClient *red.Pool
}

// Dial 连接
func Dial(host string, db, maxIdle, maxActive int) *DialContext {
	c := &DialContext{}
	// 建立连接池
	c.redisClient = &red.Pool{
		// 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxIdle: maxIdle,
		// 最大的激活连接数，表示同时最多有N个连接
		MaxActive: maxActive,
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
	return c
}

// Exec Exec
func (c *DialContext) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := c.redisClient.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}
