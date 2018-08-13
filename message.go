package redrabbit

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

type Message struct {
	Id      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
	queue   *Queue
}


func (msg Message) Ack() error {
	conn := pool.Get()
	defer conn.Close()
	r, err := conn.Do("LREM", msg.queue.keyHold, 0, msg.Id)
	count, _ := redis.Int64(r, err)
	log.Println(count)

	return err
}

func (msg Message) Ping() error {
	return nil
}

func (msg Message) Nack() error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("RPOPLPUSH", msg.queue.keyHold, msg.queue.keyWait)

	return err
}

func NewMsg(body []byte, headers map[string]string, Id string, queue *Queue) *Message {
	return &Message{
		Id,
		headers,
		body,
		queue,
	}
}
