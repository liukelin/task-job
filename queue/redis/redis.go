package redis

/**
 redis list 消息队列类
 @auth liukelin@qianxin.com
*/

import (
	"log"
	"fmt"
	"time"
	"github.com/garyburd/redigo/redis"
)

/**
 * 必要方法
 */
type QueueInterface interface {
	// 连接
	Connect() (error)
	// push
	Push(string, string) (error)
	// Get
	Get(string) (string, error)
	// 常驻get
	// GetCallback(string, func(string) bool) 
}

type Queue struct {
	Conf 	*Conf
	Pool 	*redis.Pool  // 连接池
	Conn    *redis.Conn  // 冗余
	// Callback 		 // 回调函数, 用于收到事件回调
}

type Conf struct {
	Servers []string	 // ["127.0.0.1:7711", "127.0.0.1:7712"]
	Host 	string
	Qname	string		 // 队列name
}

func (q *Queue) Connect() error {
	_, err := q.initConnection()
	return err
}

// 检测重连
func (q *Queue) initConnection() (redis.Conn, error) {
	if q.Pool == nil {
		pool := &redis.Pool {
			MaxIdle:     5,		// 最大空闲连接数 
			MaxActive:   20,	// 最大连接数
			IdleTimeout: 5 * time.Second, // 空闲超时时间, 不close的回收时间
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				con, err := redis.Dial("tcp", q.Conf.Host,
					// redis.DialPassword(conf["Password"].(string)),
					redis.DialDatabase(0),
					// redis.DialConnectTimeout(timeout*time.Second),
					// redis.DialReadTimeout(timeout*time.Second),
					// redis.DialWriteTimeout(timeout*time.Second)
				)
				if err != nil {
					return nil, err 
				}
				return con, nil 
			},
		}
		q.Pool = pool
	}
	return q.Pool.Get(), nil
}

// 写入元素
func (q *Queue) Push(data string, queueName string) error {

	conn, err := q.initConnection()
	if err != nil {
		return fmt.Errorf("Connect redis Push client error:", err)
	}
	defer conn.Close()

	if queueName == "" {
		queueName = q.Conf.Qname
	}
	
	_, err = conn.Do("rpush", queueName, data)
	return err
}

// 从list获取元素
func (q *Queue) Get(queueName string) (string, error) {

	conn, err := q.initConnection()
	if err != nil {
		log.Println("Connect redis Get client error:", err)
		return "", fmt.Errorf("Connect redis Get client error:", err)
	}

	defer conn.Close()

	if queueName == "" {
		queueName = q.Conf.Qname
	}

	data, _err := conn.Do("lpop", queueName)
	if  _err != nil {
		log.Println("redis Get lpop error:", _err)
		return "", _err
	}

	if data == nil {
		return "", nil
	}

	// 格式判断
	switch vv := data.(type) {
	case []uint8: // 字节切片类型
		return string(vv), nil
    case string:
		return data.(string), nil
    default:
        return "", fmt.Errorf("redis Get value type error:", vv)
	}
}

func dial(addr string) (redis.Conn, error) {
	return redis.Dial("tcp", addr)
}