package eventmanager

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lim-bo/calendar/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventManager struct {
	cli *mongo.Client
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func New(cfg *DBConfig) *EventManager {
	connUrl := "mongodb://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connUrl).SetConnectTimeout(time.Second*5))
	if err != nil {
		log.Fatal(err)
	}
	return &EventManager{
		cli: client,
	}
}

func (em *EventManager) AddEvent(event *models.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_, err := em.cli.Database("calend_db").Collection("events").InsertOne(ctx, *event)
	if err != nil {
		return errors.New("event manager error: " + err.Error())
	}
	return nil
}

func (em *EventManager) GetEventsByMonth(master uuid.UUID, month time.Month) ([]*models.Event, error) {
	var result []*models.Event

	return result, nil
}

func (em *EventManager) GetEvents(master uuid.UUID) ([]*models.Event, error) {
	result := make([]*models.Event, 0, 2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	cursor, err := em.cli.Database("calend_db").Collection("events").Find(ctx, bson.M{"master": master.String()})
	if err != nil {
		return nil, errors.New("searching docs error: " + err.Error())
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var event models.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, errors.New("decoding event error: " + err.Error())
		}
		result = append(result, &event)
	}
	return result, nil
}

func (em *EventManager) DeleteEvent(master uuid.UUID, name string) error {

	return nil
}
