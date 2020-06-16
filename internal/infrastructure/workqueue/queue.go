package queue

import (
	"fmt"

	"github.com/squadcast_assignment/internal/config"
	"github.com/streadway/amqp"
)

// QueueClient exposes functionalities related to a queue
type QueueClient interface {
	Publish(eName string, qName string, msg []byte) error
	Subscribe(qName string) (<-chan amqp.Delivery, func(), error)
	Close() error
}

type queueClient struct {
	conn *amqp.Connection
}

// GetConnectionToQueue creates and returns a connection to RabbitMQ
func GetConnectionToQueue(address string) (*queueClient, error) {
	queue := queueClient{}
	var err error
	queue.conn, err = amqp.Dial("amqp://" + address)
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (queue *queueClient) Publish(e string, q string, msg []byte) error {
	ch, err := queue.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	payload := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         msg,
	}
	if err := ch.ExchangeDeclare(e, "fanout", false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}
	if err := ch.Publish(e, "", false, false, payload); err != nil {
		return fmt.Errorf("failed to publish to queue: %v", err)
	}
	return nil
}

func (queue *queueClient) Subscribe(qName string) (<-chan amqp.Delivery, func(), error) {
	ch, err := queue.conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	err = ch.ExchangeDeclare(config.App.Queue.Exchange, "fanout", false, false, false, false, nil)
	q, err := ch.QueueDeclare(qName, false, false, false, false, nil)
	err = ch.QueueBind(qName, "", config.App.Queue.Exchange, false, nil)
	c, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	return c, func() { ch.Close() }, err
}

func (queue *queueClient) Close() error {
	if err := queue.conn.Close(); err != nil {
		return err
	}
	return nil
}
