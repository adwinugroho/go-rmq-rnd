package config

var Config config

type config struct {
	RabbitMQ rabbitMQ `mapstructure:"rabbitmq"`
}

type rabbitMQ struct {
	Host      string `mapstructure:"host"`
	QueueName string `mapstructure:"queue_name"`
}
