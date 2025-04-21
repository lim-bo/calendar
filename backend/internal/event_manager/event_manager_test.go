package eventmanager_test

import (
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	eventmanager "github.com/lim-bo/calendar/backend/internal/event_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	err := util.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	m.Run()
}

func TestAddEvent(t *testing.T) {
	cfg := eventmanager.DBConfig{
		Host:     viper.GetString("events_db_host"),
		Port:     viper.GetString("events_db_port"),
		User:     viper.GetString("events_db_user"),
		Password: viper.GetString("events_db_pass"),
	}
	em := eventmanager.New(cfg)
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	event := &models.Event{
		Master:       uid,
		Name:         "test event 3",
		Description:  "тусняк висняк",
		Type:         "TUSA",
		Start:        time.Now(),
		End:          time.Now().Add(time.Hour * 24),
		Participants: []uuid.UUID{uid},
	}
	err := em.AddEvent(event)
	if err != nil {
		t.Error(err)
	}

}

func TestGetEvents(t *testing.T) {
	cfg := eventmanager.DBConfig{
		Host:     viper.GetString("events_db_host"),
		Port:     viper.GetString("events_db_port"),
		User:     viper.GetString("events_db_user"),
		Password: viper.GetString("events_db_pass"),
	}
	em := eventmanager.New(cfg)
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	events, err := em.GetEvents(uid)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range events {
		t.Log(*e)
	}

}

func TestDeleteEvent(t *testing.T) {
	cfg := eventmanager.DBConfig{
		Host:     viper.GetString("events_db_host"),
		Port:     viper.GetString("events_db_port"),
		User:     viper.GetString("events_db_user"),
		Password: viper.GetString("events_db_pass"),
	}
	em := eventmanager.New(cfg)
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	id, err := primitive.ObjectIDFromHex("67fd5a426b30852b941ef893")
	if err != nil {
		t.Fatal(err)
	}
	err = em.DeleteEvent(uid, id)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEventsByMonth(t *testing.T) {
	cfg := eventmanager.DBConfig{
		Host:     viper.GetString("events_db_host"),
		Port:     viper.GetString("events_db_port"),
		User:     viper.GetString("events_db_user"),
		Password: viper.GetString("events_db_pass"),
	}
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	em := eventmanager.New(cfg)
	events, err := em.GetEventsByMonth(uid, time.April)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Error("empty slice")
	} else {
		t.Log(len(events))
	}
	for _, e := range events {
		t.Log(*e)
	}
}
