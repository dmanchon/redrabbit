package redrabbit

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

func Run() {
	log.Println("Start...")
	conn, err := redis.Dial("tcp", ":6379")
	defer conn.Close()
	if err != nil {
		// handle error
	}

	queue := New("onna/beats/ds1", conn)
	msg := NewMsg(
		[]byte("hello"),
		nil, generateUUID(),
		100.00, queue,
	)

	queue.Add(msg)
	queue.Add(msg)

	msg2, _ := queue.Get()
	log.Println(msg2)

	msg2.Nack()
}
