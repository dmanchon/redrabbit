package redrabbit

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Queue struct {
	Id      string `json:"id"`
	keyWait string
	keyHold string
	keyExps string
	ttl     int64 `json:"ttl"`
	pool    *redis.Pool
}

func (q Queue) Add(msg *Message) error {
	conn := q.pool.Get()
	defer conn.Close()

	_, err := conn.Do("RPUSH", q.keyWait, msg.Id)

	if err == nil {
		keyInfo := fmt.Sprintf("%s/__messages__/%s", q.Id, msg.Id)

		_, err = conn.Do("HMSET", keyInfo, "body", msg.Body)
		log.Printf("[DONE] queue(%s) <= msg(%s)\n", q.Id, msg.Id)
	} else {
		log.Printf("[FAIL] queue(%s) <= msg(%s)\n[FAIL] %s", q.Id, msg.Id, err)
	}

	return err
}

func (q Queue) waitForTTL() {
	for {
		conn := q.pool.Get()
		defer conn.Close()
		r, err := conn.Do("BZPOPMIN", q.keyExps, 0)
		values, err := redis.Strings(r, err)
		expiration, err := strconv.Atoi(values[2])
		diff := int64(expiration) - time.Now().Unix()
		if diff <= 0 {
			log.Println("Expired ", values)
			r, err = conn.Do("LREM", q.keyHold, values[1])
			r, err = conn.Do("RPUSH", q.keyWait, values[1])
		} else {
			log.Println("Not expired wait for ", diff, " seconds")
			r, err = conn.Do("ZADD", q.keyExps, "NX", "CH", expiration, values[1])
			time.Sleep(time.Second * time.Duration(diff))
		}
	}
}

func (q Queue) Get() (Message, error) {
	conn := q.pool.Get()
	defer conn.Close()
	r, err := conn.Do("RPOPLPUSH", q.keyWait, q.keyHold)
	//check error
	key, _ := redis.String(r, err)
	keyInfo := fmt.Sprintf("%s/__messages__/%s", q.Id, key)

	r, err = conn.Do("HGET", keyInfo, "body")
	//check error

	body, _ := redis.Bytes(r, err)
	msg := Message{Id: key, Body: body, queue: &q}

	//add the TTL
	now := time.Now().Unix()
	expiration := now + q.ttl

	r, err = conn.Do("ZADD", q.keyExps, "NX", "CH", expiration, msg.Id)
	count, _ := redis.Int64(r, err)

	if err != nil || count != 1 {

	}

	log.Printf("queue(%s) => msg(%s) with body: %s\n", q.Id, msg.Id, body)

	return msg, err
}

func (q Queue) register() error {
	conn := q.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", "Queues", q.Id)
	go q.waitForTTL()

	server.listening[q.Id] = q.ttl
	return err
}

func GetQueue(Id string, ttl int64) *Queue {
	keyWait := fmt.Sprintf("%s.WAIT", Id)
	keyHold := fmt.Sprintf("%s.HOLD", Id)
	keyExps := fmt.Sprintf("%s.EXPS", Id)

	q := &Queue{
		Id,
		keyWait,
		keyHold,
		keyExps,
		ttl,
		server.pool,
	}
	_, ok := server.listening[Id]
	if !ok {
		q.register()
	}
	return q
}
