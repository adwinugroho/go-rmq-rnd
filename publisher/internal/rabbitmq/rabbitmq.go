package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-rmq-rnd/publisher/config"
	"github.com/streadway/amqp"
)

var once sync.Once
var conn *amqp.Connection
var ch *amqp.Channel

func GetRmqConnection() (*amqp.Connection, error) {
	amqpServerURL := config.Config.RabbitMQ.Host

	once.Do(func() {
		var err error
		conn, err = amqp.DialConfig(amqpServerURL, amqp.Config{
			Heartbeat: 10 * time.Second,
		})

		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ:", err)
		}
	})

	if conn == nil || conn.IsClosed() {
		var err error
		conn, err = amqp.DialConfig(amqpServerURL, amqp.Config{
			Heartbeat: 10 * time.Second,
		})

		if err != nil {
			fmt.Println("Failed to reconnect to RabbitMQ:", err)
		}
	}

	return conn, nil
}

func GetChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	var err error
	ch, err = conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel:", err)
	}

	return ch, nil
}

func SetupRabbitMQ() (*amqp.Connection, *amqp.Channel, chan *amqp.Error, chan *amqp.Error) {
	conn, err := GetRmqConnection()
	if err != nil {
		log.Fatalf("Error while connecting to RabbitMQ: %v", err)
	}

	ch, err = GetChannel(conn)
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	connErrors := make(chan *amqp.Error)
	conn.NotifyClose(connErrors)

	chErrors := make(chan *amqp.Error)
	ch.NotifyClose(chErrors)

	return conn, ch, connErrors, chErrors
}
