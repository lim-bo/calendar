package rabbit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/lim-bo/calendar/backend/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Finalizer func()

type Producer struct {
	mqComponents struct {
		conn *amqp.Connection
		ch   *amqp.Channel
		q    *amqp.Queue
	}
}

type RabbitCfg struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewProducer(cfg RabbitCfg, queueName string) (*Producer, Finalizer) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		log.Fatal("error connecting to rabbit: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("error getting channel: " + err.Error())
	}
	err = ch.ExchangeDeclare(
		"delayed_notifications",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		log.Fatal("error declaring delayed exchange: " + err.Error())
	}
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("error declaring queue: " + err.Error())
	}
	err = ch.QueueBind(q.Name, "delayed_routing", "delayed_notifications", false, nil)
	if err != nil {
		log.Fatal("error binding queue: " + err.Error())
	}
	return &Producer{
			mqComponents: struct {
				conn *amqp.Connection
				ch   *amqp.Channel
				q    *amqp.Queue
			}{conn: conn, ch: ch, q: &q},
		}, func() {
			ch.Close()
			conn.Close()
		}
}

func (p *Producer) ProduceWithJSON(jsonMsg []byte) error {

	var msg models.Notification
	var delayed bool
	err := sonic.ConfigDefault.Unmarshal(jsonMsg, &msg)
	if err == nil {
		delayed = msg.Delayed
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         jsonMsg,
		DeliveryMode: amqp.Persistent,
	}
	if delayed {
		delay := time.Until(msg.Deadline)
		publishing.Headers = amqp.Table{
			"x-delay": delay.Milliseconds(),
		}
	}
	err = p.mqComponents.ch.PublishWithContext(ctx, "", p.mqComponents.q.Name, false, false, publishing)
	if err != nil {
		return errors.New("error producing json msg: " + err.Error())
	}
	return nil
}
