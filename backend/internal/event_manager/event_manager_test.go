package eventmanager_test

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	eventmanager "github.com/lim-bo/calendar/backend/internal/event_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	err := util.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	m.Run()
}

func TestAddEvent(t *testing.T) {
	cfg := &eventmanager.DBConfig{
		Host:     viper.GetString("events_db_host"),
		Port:     viper.GetString("events_db_port"),
		User:     viper.GetString("events_db_user"),
		Password: viper.GetString("events_db_pass"),
	}
	em := eventmanager.New(cfg)
	uid, err := uuid.Parse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	if err != nil {
		t.Fatal(err)
	}
	event := &models.Event{
		Master:       uid,
		Name:         "test event",
		Description:  "тусняк висняк",
		Type:         "TUSA",
		Start:        time.Now(),
		End:          time.Now().Add(time.Hour * 24),
		Participants: []uuid.UUID{uid},
	}
	err = em.AddEvent(event)
	if err != nil {
		t.Error(err)
	}
}
