package main

import (
	"encoding/json"
	"flag"
	"github.com/dmanchon/redrabbit"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

	q := redrabbit.GetQueue(path, 5)
	msg, _ := q.Get()

	json.NewEncoder(w).Encode(msg)
}

func createQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	q := redrabbit.GetQueue(path, 5)

	json.NewEncoder(w).Encode(q)

}

func putMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	q := redrabbit.GetQueue(path, 5)
	body, _ := ioutil.ReadAll(r.Body)

	msg := redrabbit.NewMsg(
		body,
		nil, redrabbit.GenerateUUID(),
		q,
	)
	q.Add(msg)
	json.NewEncoder(w).Encode(msg)
}

func main() {
	host := flag.String("host", "localhost:6379", "Redis host in the format 'host:port'")
	flag.Parse()

	log.Printf("Start, connecting to redis[%s]...", *host)
	pool := newPool(*host)
	defer pool.Close()

	app := redrabbit.Start(pool)

	router := mux.NewRouter()

	// The order of the definition of the routes is important
	router.HandleFunc("/q/{path:.+?}/@get", getMessage).Methods("GET")
	router.HandleFunc("/q/{path:.+?}/@put", putMessage).Methods("POST")

	// Get all queues
	router.HandleFunc("/q", func(w http.ResponseWriter, r *http.Request) {
		getQueues(app, w)
	}).Methods("GET")

	router.HandleFunc("/q/{path:.+?}", createQueue).Methods("POST")
	//router.HandleFunc("/q/{path:.+?}", deleteQueue).Methods("DELETE")

	//router.HandleFunc("/q/{path:.+?}/@info", queueInfo).Methods("GET")

	//router.HandleFunc("/q/{path:.+?}/{msg}/@ack", ackMessage).Methods("POST")
	//router.HandleFunc("/q/{path:.+?}/{msg}/@nack", nackMessage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
