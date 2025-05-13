package eventmanager

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lim-bo/calendar/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrLackOrWrongMaster = errors.New("event created by other user or event doesn't exist")
	ErrNoSuchEvent       = errors.New("there is no event with such ID")
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

func New(cfg DBConfig) *EventManager {
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
	dateStart := time.Date(time.Now().Year(), month, 1, 0, 0, 0, 0, time.Now().Location())
	dateEnd := time.Date(time.Now().Year(), month+time.Month(1), 1, 0, 0, 0, 0, time.Now().Location())
	cursor, err := em.cli.Database("calend_db").Collection("events").Find(context.Background(), bson.M{
		"parts.uid": master,
		"start": bson.M{
			"$gte": dateStart,
		},
		"end": bson.M{
			"$lt": dateEnd,
		},
	})
	if err != nil {
		return nil, errors.New("error getting events by month: " + err.Error())
	}
	defer cursor.Close(context.Background())
	result := make([]*models.Event, 0, cursor.RemainingBatchLength())
	for cursor.Next(context.Background()) {
		var event models.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, errors.New("decoding event error: " + err.Error())
		}
		result = append(result, &event)
	}
	return result, nil
}

func (em *EventManager) GetEventsByWeek(master uuid.UUID) ([]*models.Event, error) {
	daysToMonday := int(time.Now().Weekday()) - 1
	if daysToMonday < 0 {
		daysToMonday = 6
	}
	now := time.Now()
	weekStart := time.Date(now.Year(),
		now.Month(),
		now.Day()-daysToMonday, 0, 0, 0, 0,
		now.Location(),
	)
	weekEnd := weekStart.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
	filter := bson.M{
		"$and": []bson.M{
			{"start": bson.M{"$lte": weekEnd}},
			{"end": bson.M{"$gte": weekStart}},
		},
		"parts.uid": master,
	}
	cursor, err := em.cli.Database("calend_db").Collection("events").Find(context.Background(), filter)
	if err != nil {
		return nil, errors.New("error getting events by week: " + err.Error())
	}
	defer cursor.Close(context.Background())
	result := make([]*models.Event, 0, cursor.RemainingBatchLength())
	for cursor.Next(context.Background()) {
		var event models.Event
		cursor.Decode(&event)
		result = append(result, &event)
	}
	return result, nil
}

func (em *EventManager) GetEventsByDay(master uuid.UUID, day time.Time) ([]*models.Event, error) {
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	nextDay := day.Add(time.Hour * 24)
	cursor, err := em.cli.Database("calend_db").Collection("events").Find(context.Background(), bson.M{
		"$and": []bson.M{
			{"start": bson.M{"$lte": nextDay}},
			{"end": bson.M{"$gte": day}},
		},
		"parts.uid": master,
	})
	if err != nil {
		return nil, errors.New("errors getting events by day: " + err.Error())
	}
	defer cursor.Close(context.Background())
	result := make([]*models.Event, 0, cursor.RemainingBatchLength())
	for cursor.Next(context.Background()) {
		var event models.Event
		cursor.Decode(&event)
		result = append(result, &event)
	}
	return result, nil
}

func (em *EventManager) GetEvents(master uuid.UUID) ([]*models.Event, error) {
	cursor, err := em.cli.Database("calend_db").Collection("events").Find(context.Background(), bson.M{"parts.uid": master})
	if err != nil {
		return nil, errors.New("searching docs error: " + err.Error())
	}
	defer cursor.Close(context.Background())
	result := make([]*models.Event, 0, cursor.RemainingBatchLength())
	for cursor.Next(context.Background()) {
		var event models.Event
		if err = cursor.Decode(&event); err != nil {
			return nil, errors.New("decoding event error: " + err.Error())
		}
		result = append(result, &event)
	}
	return result, nil
}

func (em *EventManager) DeleteEvent(master uuid.UUID, id primitive.ObjectID) error {
	result, err := em.cli.Database("calend_db").Collection("events").DeleteOne(context.Background(), bson.M{"master": master, "_id": id})
	if err != nil {
		return errors.New("deleting event error: " + err.Error())
	}
	if result.DeletedCount == 0 {
		return ErrLackOrWrongMaster
	}
	return nil
}

func (em *EventManager) UpdateEvent(event *models.Event) error {
	res, err := em.cli.Database("calend_db").Collection("events").UpdateOne(context.Background(), bson.M{
		"_id": event.ID,
	}, bson.M{
		"$set": bson.M{
			"name":  event.Name,
			"desc":  event.Description,
			"type":  event.Type,
			"prior": event.Prior,
			"start": event.Start,
			"end":   event.End,
			"parts": event.Participants,
		},
	})
	if err != nil {
		return errors.New("updating event error: " + err.Error())
	}
	if res.MatchedCount == 0 {
		return ErrNoSuchEvent
	}
	return nil
}

func (em *EventManager) GetEventByID(id primitive.ObjectID) (*models.Event, error) {
	res := em.cli.Database("calend_db").Collection("events").FindOne(context.Background(), bson.M{
		"_id": id,
	})
	if res.Err() != nil {
		return nil, errors.New("error getting event by id: " + res.Err().Error())
	}
	var event models.Event
	err := res.Decode(&event)
	if err != nil {
		return nil, errors.New("error unmarshalling eventByID results: " + err.Error())
	}
	return &event, nil
}

func (em *EventManager) ChangeUserAcceptance(eventID primitive.ObjectID, uid uuid.UUID, accepted bool) error {
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{
				"elem.uid": uid,
			},
		},
	}
	res, err := em.cli.Database("calend_db").Collection("events").UpdateOne(context.Background(), bson.M{
		"_id":       eventID,
		"parts.uid": uid,
	}, bson.M{
		"$set": bson.M{
			"parts.$[elem].accepted": accepted,
		},
	}, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return errors.New("updating user state on event error: " + err.Error())
	}
	if res.MatchedCount == 0 {
		return ErrNoSuchEvent
	}
	return nil
}
