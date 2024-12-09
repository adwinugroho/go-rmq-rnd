package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/apsyadira-jubelio/go-amqp-reconnect/rabbitmq"
	"github.com/davecgh/go-spew/spew"
	"github.com/gammazero/workerpool"
	"github.com/go-rmq-rnd/consumer/config"
	"github.com/go-rmq-rnd/consumer/internal/server"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("RUNNING CONSUMER RabbitMQ")

	godotenv.Load()

	spew.Dump(config.Config)
	ctx, closeCtx := context.WithCancel(context.Background())

	queueName := "consumer-message"

	rabbitmq.Debug = true
	rmqClient := rabbitmq.NewRmqClient(&rabbitmq.ChannelOptions{
		ConnectionString: config.Config.RabbitMQ.Host,
	})

	err := rmqClient.OpenChannel(&rabbitmq.ChannelOptions{
		QueueName:        queueName,
		QueueTTL:         3600000,
		QueueDurable:     true,
		QueueExchange:    fmt.Sprintf("%s-exchange", queueName),
		DeadQueue:        fmt.Sprintf("%s-dead-queue", queueName),
		DeadExchange:     fmt.Sprintf("%s-dead-exchange", queueName),
		DeadQueueDurable: false,
		DeadQueueTTL:     3600000,
		QueueRoutingKey:  queueName,
	})

	if err != nil {
		log.Fatal("error while openning RabbitMQ channel", err)
	}

	log.Println("successfully connecting to RabbitMQ with queue: ", queueName)
	go server.HTTPServer(ctx, queueName)

	// Create a context that is cancelled on system interrupt or SIGTERM signal
	defer closeCtx()

	consumerDone := make(chan bool, 1)
	wp := workerpool.New(1)
	defer wp.Stop() // Ensure all workers are stopped before exiting

	go func() {
		d, err := rmqClient.Consume()
		if err != nil {
			log.Fatal("Error cause:", err)
		}

		for {
			stop := false

			select {
			case <-ctx.Done():
				stop = true
			case msg := <-d:
				message := msg
				// count := 1
				arrMessage := make([]amqp.Delivery, 0)
				arrMessage = append(arrMessage, message)
				wp.Submit(func() {
					if len(arrMessage) > 2 {
						log.Println("start to processing 100 message:", len(arrMessage))
						for _, eachMessage := range arrMessage {
							var messageBody map[string]interface{}
							if err := json.Unmarshal(eachMessage.Body, &messageBody); err != nil {
								message.Ack(false)
								log.Println("error unmarshalling payload", err)
								return
							}
							log.Println("message number 1:", eachMessage)
							// var payload map[string]interface{}
							// data, _ := json.Marshal(messageBody.Data)
							// if err := json.Unmarshal(data, &payload); err != nil {
							// 	message.Ack(false)
							// 	log.Println("error extracting payload data", err)
							// 	return
							// }

							eachMessage.Ack(false)
						}

					}
				})
			}

			if stop {
				break
			}
		}

		consumerDone <- true
	}()

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-signals
		log.Println("received signal", map[string]interface{}{"signal": s})
		closeCtx()

		<-consumerDone
		rmqClient.Close()
		log.Println("connection and channel closed")

		done <- true
	}()

	log.Println("waiting...")
	<-done

	wp.StopWait()
	log.Println("exiting")
}
