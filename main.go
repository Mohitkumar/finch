package main

import (
	"fmt"

	"github.com/mohitkumar/finch/rest"
)

func main() {
	data := make(map[string]any)
	inner := make(map[string]any)
	data["a"] = 12
	data["output"] = inner
	inner["status"] = 200
	rdConf := rest.RedisConfig{
		Host:      "localhost",
		Port:      6379,
		Namespace: "finch",
	}
	conf := rest.Config{
		Port:        8080,
		RedisConfig: rdConf,
	}
	server, err := rest.NewServer(conf)
	if err != nil {
		panic(err)
	}

	err = server.Start()

	if err != nil {
		fmt.Println("could not start server")
	}
}
