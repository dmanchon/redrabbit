package redrabbit

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
	listening map[string]int64

)

func newPool(host string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", host) },
	}
}

func Init() error {
	listening = make(map[string]int64)
	err, _ := ListQueues()
	return err
}


func Run(host string) {
	log.Printf("Start, connecting to redis[%s]...", host)
	pool = newPool(host)
	defer pool.Close()
	Init()

	queue := GetQueue("onna/beats/ds1", 5)

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
