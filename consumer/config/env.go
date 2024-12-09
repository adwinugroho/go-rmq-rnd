package config

var Config config

type config struct {
	Environment string   `mapstructure:"environment"`
	RabbitMQ    rabbitMQ `mapstructure:"rabbitmq"`
}

type rabbitMQ struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
