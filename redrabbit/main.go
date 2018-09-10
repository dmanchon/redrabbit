package main

import "github.com/dmanchon/redrabbit"
import "flag"

func main() {
	host := flag.String("host", "localhost:6379", "Redis host in the format 'host:port'")
	flag.Parse()

	redrabbit.Run(*host)
}
