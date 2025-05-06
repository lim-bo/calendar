package rabbit_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/lim-bo/calendar/backend/internal/notifier"
	"github.com/lim-bo/calendar/backend/internal/rabbit"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	err := notifier.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}

func TestRabbit(t *testing.T) {
	cfg := rabbit.RabbitCfg{
		Host:     viper.GetString("rabbit_host"),
		Port:     viper.GetString("rabbit_port"),
		Username: viper.GetString("rabbit_user"),
		Password: viper.GetString("rabbit_pass"),
	}
	p, cancel := rabbit.NewProducer(cfg, "test")
	defer cancel()
	for i := range 10 {
		msg := fmt.Appendf(make([]byte, 0), `{"message_num": %d}`, i)
		err := p.ProduceWithJSON(msg)
		if err != nil {
			t.Fatal(err)
		}
	}
	cancel()
	c, cancel := rabbit.NewConsumer(cfg, "test")
	defer cancel()
	ctx, ctxcancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(10)
	c.Run(ctx, func(d amqp091.Delivery) error {
		t.Log(string(d.Body))
		wg.Done()
		return nil
	})
	wg.Wait()
	ctxcancel()
}
