package redrabbit

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") },
	}
}

func Run() {
	log.Println("Start...")
	pool = newPool()
	defer pool.Close()

	queue := NewQueue("onna/beats/ds1", 5)

	for i := 0; i < 1; i++ {
		msg := NewMsg(
			[]byte("hello"),
			nil, generateUUID(),
			queue,
		)
		queue.Add(msg)
	}

	msg2, _ := queue.Get()
	log.Println(msg2)

	time.Sleep(time.Second * 10)
	msg2, _ = queue.Get()

}
