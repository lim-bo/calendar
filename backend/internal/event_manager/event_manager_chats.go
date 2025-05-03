package eventmanager

import (
	"context"
	"errors"
	"time"

	"github.com/lim-bo/calendar/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (em *EventManager) SendMessage(eventID primitive.ObjectID, msg *models.MessageWithMail) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := em.cli.Database("calend_db").Collection("chats").UpdateOne(ctx, bson.M{
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

func (em *EventManager) GetMessages(eventID primitive.ObjectID) (*models.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res := em.cli.Database("calend_db").Collection("chats").FindOne(ctx, bson.M{
		"event_id": eventID,
	})
	if res.Err() != nil {
		return nil, errors.New("getting chat error: " + res.Err().Error())
	}
	var chat models.Chat
	err := res.Decode(&chat)
	if err != nil {
		return nil, errors.New("error decoding chat content: " + err.Error())
	}
	return &chat, nil
}

func (em *EventManager) DeleteChat(eventID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := em.cli.Database("calend_db").Collection("chats").DeleteOne(ctx, bson.M{
		"event_id": eventID,
	})
	if err != nil {
		return errors.New("error deleting chat: " + err.Error())
	}
	return nil
}
