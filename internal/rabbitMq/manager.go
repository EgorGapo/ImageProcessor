package rabbitmq

import (
	"encoding/json"
	"test/internal/domain"

	"github.com/streadway/amqp"
)

type RabbitMQConnectionManager struct {
	channel    *amqp.Channel
	connection *amqp.Connection
	queueName  string
}

func NewRabbitMQConnection(amqpUrl, queueName string) (*RabbitMQConnectionManager, error) {
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(
		queueName,
		true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQConnectionManager{
		channel:    ch,
		connection: conn,
		queueName:  queueName,
	}, nil
}

func (s *RabbitMQConnectionManager) Send(task domain.Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = s.channel.Publish(
		"",
		s.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}

func (s *RabbitMQConnectionManager) Receive() (<-chan amqp.Delivery, error) {
	msgs, err := s.channel.Consume(
		s.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, err
}

func (s *RabbitMQConnectionManager) Close() {
	s.channel.Close()
	s.connection.Close()
}
