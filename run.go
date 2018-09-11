package redrabbit

import (
	"github.com/gomodule/redigo/redis"
)

var (
	server *Server
)

type Server struct {
	pool      *redis.Pool
	listening map[string]int64
}

func (s Server) ListQueues() (error, []string) {
	conn := s.pool.Get()
	defer conn.Close()

	var next int64
	var data []string
	ret := make([]string, 0)

	for {
		r, err := redis.Values(conn.Do("SSCAN", "Queues", next))
		if err != nil {
			return err, nil
		}
		_, err = redis.Scan(r, &next, &data)
		if err != nil {
			return err, nil
		}
		ret = append(ret, data...)
		if next == 0 {
			break
		}
	}
	return nil, ret
}

func Start(pool *redis.Pool) *Server {
	if server != nil {
		return server
	}

	listening := make(map[string]int64)

	server = &Server{
		pool,
		listening,
	}

	_, queues := server.ListQueues()
	for _, id := range queues {
		GetQueue(id, 5)
	}

	return server
}
