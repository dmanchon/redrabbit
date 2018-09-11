package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/dmanchon/redrabbit"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

func newPool(host string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", host) },
	}
}

func getQueues(app *redrabbit.Server, w http.ResponseWriter) {
	_, queues := app.ListQueues()
	json.NewEncoder(w).Encode(queues)
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	json.NewEncoder(w).Encode(path)

}

func main() {
	host := flag.String("host", "localhost:6379", "Redis host in the format 'host:port'")
	flag.Parse()

	log.Printf("Start, connecting to redis[%s]...", host)
	pool := newPool(*host)
	defer pool.Close()

	app := redrabbit.Start(pool)

	//queue := redrabbit.GetQueue("onna/beats/ds1", 5)

	//for i := 0; i < 1; i++ {
	//	msg := redrabbit.NewMsg(
	//		[]byte("hello"),
	//		nil, redrabbit.GenerateUUID(),
	//		queue,
	//	)
	//	queue.Add(msg)
	//}

	//msg2, _ := queue.Get()
	//log.Println(msg2)

	//time.Sleep(time.Second * 10)
	//msg2, _ = queue.Get()

	router := mux.NewRouter()

	router.HandleFunc("/q", func(w http.ResponseWriter, r *http.Request) {
		getQueues(app, w)
	}).Methods("GET")

	router.HandleFunc("/q/{path:.+?}/@get", getMessage).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}
