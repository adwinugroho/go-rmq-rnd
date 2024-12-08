package main

import (
	"log"

	"github.com/go-rmq-rnd/publisher/config"
)

func init() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("RUNNING PUBLISHER RabbitMQ")
}
