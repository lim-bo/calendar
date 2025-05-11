package notifier

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/smtp"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/jordan-wright/email"
	"github.com/lim-bo/calendar/backend/internal/rabbit"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type NotifyService struct {
	inputChan chan *models.Notification
	errorChan chan error
	closeChan chan struct{}

	cfg SMTPConfig

	consumer          *rabbit.Consumer
	consumerCloseFunc rabbit.Finalizer
}

func New() *NotifyService {
	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg := SMTPConfig{
		Host:     viper.GetString("smtp_host"),
		Port:     viper.GetString("smtp_port"),
		Username: viper.GetString("smtp_user"),
		Password: viper.GetString("smtp_pass"),
		From:     viper.GetString("smtp_from"),
	}
	mqcfg := rabbit.RabbitCfg{
		Host:     viper.GetString("rabbit_host"),
		Port:     viper.GetString("rabbit_port"),
		Username: viper.GetString("rabbit_user"),
		Password: viper.GetString("rabbit_pass"),
	}
	c, cancel := rabbit.NewConsumer(mqcfg, "notifications")
	closeChan := make(chan struct{}, 1)
	inputChan := make(chan *models.Notification)
	errorChan := make(chan error)
	return &NotifyService{
		inputChan:         inputChan,
		cfg:               cfg,
		errorChan:         errorChan,
		consumer:          c,
		closeChan:         closeChan,
		consumerCloseFunc: cancel,
	}
}

func (ns *NotifyService) Close() {
	ns.consumerCloseFunc()
	close(ns.closeChan)
	close(ns.inputChan)
	close(ns.errorChan)
}

func (ns *NotifyService) notifyWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ns.closeChan:
			return
		case task := <-ns.inputChan:
			e := email.NewEmail()
			e.To = task.To
			e.From = ns.cfg.From
			e.Subject = task.Subject
			e.Text = []byte(task.Content)
			slog.Debug("emails", slog.String("from", ns.cfg.From), slog.String("to", e.To[0]))
			err := e.Send(ns.cfg.Host+":"+ns.cfg.Port, smtp.PlainAuth("", ns.cfg.Username, ns.cfg.Password, ns.cfg.Host))
			if err != nil {
				ns.errorChan <- errors.New("sending email notification error: " + err.Error())
				continue
			}
			slog.Info("sended email", slog.String("reciever", task.To[0]), slog.String("content", task.Content))
		}
	}
}

func (ns *NotifyService) Run() error {
	go func() {
		for {
			select {
			case err := <-ns.errorChan:
				slog.Error("error occured", slog.String("error_desc", err.Error()))
			case <-ns.closeChan:
				return
			}
		}
	}()
	err := ns.consumer.Run(context.Background(), func(d amqp091.Delivery) error {
		var task models.Notification
		err := sonic.ConfigDefault.Unmarshal(d.Body, &task)
		if err != nil {
			ns.errorChan <- err
			return err
		}
		ns.inputChan <- &task
		return nil
	})
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go ns.notifyWorker(wg)
	}
	log.Println("notifications service started")
	wg.Wait()
	return nil
}
