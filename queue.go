package redrabbit

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

type Queue struct {
	conn redis.Conn
	Id   string
}

func (q Queue) Add(msg *Message) error {
	keyWait := addSuffix(q.Id, ".WAIT")

	_, err := q.conn.Do("RPUSH", keyWait, msg.Id)

	if err == nil {
		_, err = q.conn.Do("HMSET", msg.Id, "body", msg.Body, "ttl", msg.ttl)
		log.Printf("[DONE] queue(%s) <= msg(%s)\n", q.Id, msg.Id)
	} else {
		log.Printf("[FAIL] queue(%s) <= msg(%s)\n[FAIL] %s", q.Id, msg.Id, err)
	}

	return err
}

func (q Queue) Get() (Message, error) {
	keyWait := addSuffix(q.Id, ".WAIT")
	keyHold := addSuffix(q.Id, ".HOLD")

	r, err := q.conn.Do("RPOPLPUSH", keyWait, keyHold)
	key, _ := redis.String(r, err)
	r, err = q.conn.Do("HGET", key, "body")

	body, _ := redis.Bytes(r, err)
	msg := Message{Id: key, Body: body,  queue: &q}

	log.Printf("queue(%s) => msg(%s) with body: %s\n", q.Id, msg.Id, body)

	return msg, err
}

func New(Id string, conn redis.Conn) *Queue {
	return &Queue{
		conn,
		Id,
	}
}
