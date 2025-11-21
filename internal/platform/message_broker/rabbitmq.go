package messagebroker

import (
	"fmt"
	"loan_service/configs"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(cfg configs.RabbitMQConfig) (*amqp091.Connection, error) {
	connURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	conn, err := amqp091.Dial(connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	log.Printf("Successfully connected to RabbitMQ")
	return conn, nil
}
