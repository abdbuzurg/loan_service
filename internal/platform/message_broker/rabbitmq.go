package messagebroker

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func NewRabbitMQConnection() (*amqp091.Connection, error) {
	connURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)

	conn, err := amqp091.Dial(connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	log.Printf("Successfully connected to RabbitMQ")
	return conn, nil
}
