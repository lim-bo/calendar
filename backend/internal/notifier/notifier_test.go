package notifier_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/lim-bo/calendar/backend/internal/notifier"
	"github.com/lim-bo/calendar/backend/internal/rabbit"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	err := notifier.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}

func TestSendEmails(t *testing.T) {
	cfg := rabbit.RabbitCfg{
		Host:     viper.GetString("rabbit_host"),
		Port:     viper.GetString("rabbit_port"),
		Username: viper.GetString("rabbit_user"),
		Password: viper.GetString("rabbit_pass"),
	}
	p, cancel := rabbit.NewProducer(cfg, "notifications")
	defer cancel()
	for i := range 1 {
		task := models.Notification{
			To:      []string{"lim-bo@yandex.ru"},
			Content: fmt.Sprintf("message №%d", i),
			Subject: "Test message",
		}
		raw, err := sonic.ConfigDefault.Marshal(task)
		if err != nil {
			t.Fatal(err)
		}
		err = p.ProduceWithJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
	}
}
