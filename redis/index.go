package lredis

import (
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
)

var lxpool *redigo.Pool
var lxresconfig MapRedisConfig //默认配置

type MapRedisConfig struct {
	Host   string
	Port   int
	Pwd    string
	DB     string
	Prefix string
}

type LXredis struct {
	Name string
	Val  string
	DB   string
}

func GetLxRedis() redigo.Conn {
	return lxpool.Get()
}

func InitRedis(config MapRedisConfig) (err error) {

	lxresconfig = config

	pool_size := 20

	lxpool = redigo.NewPool(func() (redigo.Conn, error) {
		c, err := redigo.Dial("tcp", fmt.Sprintf("%s:%d", lxresconfig.Host, lxresconfig.Port))
		if err != nil {
			return nil, err
		}
		if _, err := c.Do("AUTH", lxresconfig.Pwd); err != nil {
			c.Close()
			return nil, err
		}
		if _, err := c.Do("SELECT", lxresconfig.DB); err != nil {
			c.Close()
			return nil, err
		}
		return c, nil
	}, pool_size)

	return
}

//判断key是否存在
func (con *LXredis) IsExist() (isexit bool, err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	isexit, err = redigo.Bool(redis.Do("EXISTS", nname))
	return
}

//设置redis字符串
func (con *LXredis) Set() (err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	_, err = redis.Do("set", nname, con.Val)

	return
}

//设置redis字符串 加过期时间
func (con *LXredis) SetTime(times int) (err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	_, err = redis.Do("set", nname, con.Val, "EX", times)

	return
}

//获取redis
func (con *LXredis) GET() (gstr string, err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	gstr, err = redigo.String(redis.Do("GET", nname))

	return
}

//自增
func (con *LXredis) INCR() (nums int64, err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	nums, err = redigo.Int64(redis.Do("INCR", nname))

	return
}

//删除
func (con *LXredis) DELETE() (err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	is_key_exit, _ := redigo.Bool(redis.Do("EXISTS", nname))
	if is_key_exit {
		_, err = redis.Do("DEL", nname)
		if err != nil {
			return
		}
	}

	return
}

//设置过期时间
func (con *LXredis) EXPIRE(times int) (err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	_, err = redigo.Int(redis.Do("EXPIRE", nname, times))

	return
}

//查看过期时间
func (con *LXredis) TTL() (ttl int, err error) {

	nname := lxresconfig.Prefix + con.Name

	redis := GetLxRedis()

	defer redis.Close()

	if con.DB != "" {
		if _, err = redis.Do("SELECT", con.DB); err != nil {
			return
		}
	}

	ttl, err = redigo.Int(redis.Do("ttl", nname))

	return
}
