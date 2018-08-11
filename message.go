package redrabbit

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

type Message struct {
	Id      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
	ttl     int64
	queue   *Queue
}


func (msg Message) Ack() error {
	keyHold := addSuffix(msg.queue.Id, ".HOLD")
	r, err := msg.queue.conn.Do("LREM", keyHold, 0, msg.Id)
	count, _ := redis.Int64(r, err)
	log.Println(count)

	return err
}

func (msg Message) Ping() error {
	return nil
}

func (msg Message) Nack() error {
	keyWait := addSuffix(msg.queue.Id, ".WAIT")
	keyHold := addSuffix(msg.queue.Id, ".HOLD")

	_, err := msg.queue.conn.Do("RPOPLPUSH", keyHold, keyWait)

	return err
}

func NewMsg(body []byte, headers map[string]string, Id string, ttl int64, queue *Queue) *Message {
	return &Message{
		Id,
		headers,
		body,
		ttl,
		queue,
	}
}
