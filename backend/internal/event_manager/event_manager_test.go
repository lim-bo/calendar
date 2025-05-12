package eventmanager_test

import (
	"fmt"
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

var cfg eventmanager.DBConfig

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	err := util.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	cfg = eventmanager.DBConfig{
		Host:     "localhost",
		Port:     viper.GetString("events_db_port"),
		User:     viper.GetString("events_db_user"),
		Password: viper.GetString("events_db_pass"),
	}
	m.Run()
}

func TestAddEvent(t *testing.T) {
	em := eventmanager.New(cfg)
	timestamp := time.Now()
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	event := &models.Event{
		EventBase: models.EventBase{
			Master:      uid,
			Name:        "FOR NEW STRUCT",
			Description: "тусняк висняк",
			Type:        "TUSA",
			Start:       timestamp,
			End:         timestamp.Add(time.Hour * 24),
			Prior:       models.PriorityHigh,
		},
		Participants: []models.Participant{{UID: uid, Accepted: true}},
	}
	err := em.AddEvent(event)
	if err != nil {
		t.Error(err)
	}

}

func TestGetEvents(t *testing.T) {
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
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	em := eventmanager.New(cfg)
	events, err := em.GetEventsByMonth(uid, time.June)
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

func TestGetEventsByWeek(t *testing.T) {
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	em := eventmanager.New(cfg)
	events, err := em.GetEventsByWeek(uid)
	if err != nil {
		log.Fatal(err)
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

func TestGetEventsByDay(t *testing.T) {
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	em := eventmanager.New(cfg)
	day := time.Now()
	events, err := em.GetEventsByDay(uid, day)
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

func TestSendMessage(t *testing.T) {
	em := eventmanager.New(cfg)
	objID, err := primitive.ObjectIDFromHex("680e312715e5c401ed91290a")
	if err != nil {
		t.Fatal(err)
	}
	for i := range 10 {
		err = em.SendMessage(objID, &models.MessageWithMail{
			Sender:  "sender",
			Content: fmt.Sprintf("message No. %d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestGetMessages(t *testing.T) {
	em := eventmanager.New(cfg)
	objID, err := primitive.ObjectIDFromHex("680e312715e5c401ed91290a")
	if err != nil {
		t.Fatal(err)
	}
	chat, err := em.GetMessages(objID)
	if err != nil {
		t.Fatal(err)
	}
	for i, msg := range chat.Messages {
		t.Logf("message %d: %s", i+1, msg.Content)
	}
}

func TestDeleteChat(t *testing.T) {
	em := eventmanager.New(cfg)
	objID, err := primitive.ObjectIDFromHex("680e312715e5c401ed91290a")
	if err != nil {
		t.Fatal(err)
	}
	err = em.DeleteChat(objID)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateEvent(t *testing.T) {
	objID, err := primitive.ObjectIDFromHex("680fd43950d2b42f701061c0")
	if err != nil {
		t.Fatal(err)
	}
	em := eventmanager.New(cfg)
	uid := uuid.MustParse("c882bd5c-e2fb-4ca7-b291-6d751addf2d9")
	timestamp := time.Now().Add(time.Hour * 24 * 31)
	event := &models.Event{
		EventBase: models.EventBase{
			ID:          objID,
			Master:      uid,
			Name:        "FOR UPDATE",
			Description: "тусняк висняк",
			Type:        "Чаепитие",
			Start:       timestamp,
			End:         timestamp.Add(time.Hour * 24),
			Prior:       models.PriorityHigh,
		},
		Participants: []models.Participant{{UID: uid, Accepted: true}},
	}
	err = em.UpdateEvent(event)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEventByID(t *testing.T) {
	objID, err := primitive.ObjectIDFromHex("6814a8011117e998968fcc97")
	if err != nil {
		t.Fatal(err)
	}
	em := eventmanager.New(cfg)
	event, err := em.GetEventByID(objID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(event)
}
