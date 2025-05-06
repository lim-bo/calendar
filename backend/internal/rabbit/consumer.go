package rabbit

import (
	"context"
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	mqComponents struct {
		conn *amqp.Connection
		ch   *amqp.Channel
		q    *amqp.Queue
	}
	closeChan chan struct{}
}

func NewConsumer(cfg RabbitCfg, queueName string) (*Consumer, Finalizer) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		log.Fatal("error connecting to rabbit: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("error getting channel: " + err.Error())
	}
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("error declaring queue: " + err.Error())
	}
	closeChan := make(chan struct{}, 1)
	return &Consumer{
			mqComponents: struct {
				conn *amqp.Connection
				ch   *amqp.Channel
				q    *amqp.Queue
			}{conn: conn, ch: ch, q: &q},
			closeChan: closeChan,
		}, func() {
			closeChan <- struct{}{}
			close(closeChan)
			ch.Close()
			conn.Close()

		}
}

func (c *Consumer) Run(ctx context.Context, handler func(d amqp.Delivery) error) error {
	consumeChan, err := c.mqComponents.ch.Consume(c.mqComponents.q.Name, "", false, false, false, false, nil)
	if err != nil {
		return errors.New("setting handler error: " + err.Error())
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.closeChan:
				return
			case d, ok := <-consumeChan:
				if !ok {
					return
				}
				if err := handler(d); err != nil {
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}
			}
		}
	}()
	return nil
}
