package eventmanager

import (
	"context"
	"errors"

	"github.com/lim-bo/calendar/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (em *EventManager) SendMessage(eventID primitive.ObjectID, msg models.Message) error {
	_, err := em.cli.Database("calend_db").Collection("chats").UpdateOne(context.Background(), bson.M{
		"event_id": eventID,
	}, bson.M{
		"$push": bson.M{
			"messages": msg,
		},
		"$setOnInsert": bson.M{
			"event_id": eventID,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return errors.New("sending message error: " + err.Error())
	}
	return nil
}
