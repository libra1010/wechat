package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	//"github.com/gomodule/redigo/redis"
)

//Redis redis cache
type Redis struct {
	conn *redis.Client
}

//RedisOpts redis 连接属性
type RedisOpts struct {
	Host        string `yml:"host" json:"host"`
	Password    string `yml:"password" json:"password"`
	Database    int    `yml:"database" json:"database"`
	MaxIdle     int    `yml:"max_idle" json:"max_idle"`
	MaxActive   int    `yml:"max_active" json:"max_active"`
	IdleTimeout int32  `yml:"idle_timeout" json:"idle_timeout"` //second
}

type SetinelRedisOpts struct {
	MasterName    string
	Password      string
	Database      int
	SentinelNodes []string
}

//NewRedis 实例化
func NewRedis(opts *RedisOpts) *Redis {

	client := redis.NewClient(&redis.Options{
		Addr:     opts.Host,
		Password: opts.Password, // no password set
		DB:       opts.Database, // use default DB
	})

	return &Redis{conn: client}
}

func NewSetinelRedis(opts *SetinelRedisOpts) *Redis {
	sf := &redis.FailoverOptions{
		// The master name.
		MasterName: opts.MasterName,
		// A seed list of host:port addresses of sentinel nodes.
		SentinelAddrs: opts.SentinelNodes,

		// Following options are copied from Options struct.
		Password: opts.Password,
		DB:       opts.Database,
	}
	client := redis.NewFailoverClient(sf)
	return &Redis{conn: client}
}

//SetConn 设置conn
func (r *Redis) SetConn(conn *redis.Client) {
	r.conn = conn
}

//Get 获取一个值
func (r *Redis) Get(key string) interface{} {
	res, err := r.conn.Get(key).Result()
	if err != nil {
		return nil
	}

	var reply interface{}
	if err = json.Unmarshal([]byte(res), &reply); err != nil {
		return nil
	}

	return reply
}

//Set 设置一个值
func (r *Redis) Set(key string, val interface{}, timeout time.Duration) (err error) {

	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}

	err = r.conn.Set(key, data, timeout).Err()
	return

}

//IsExist 判断key是否存在
func (r *Redis) IsExist(key string) bool {
	exists := r.conn.Exists(key)

	i, _ := exists.Result()

	if i > 0 {
		return true
	}
	return false
}

//Delete 删除
func (r *Redis) Delete(key string) error {
	status := r.conn.Del(key)

	if _, err := status.Result(); err != nil {
		return err
	}

	return nil
}
