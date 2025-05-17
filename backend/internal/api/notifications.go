package api

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/lim-bo/calendar/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (api *API) sendNotification(notification *models.Notification) {
	raw, err := sonic.Marshal(notification)
	if err != nil {
		slog.Error("error marshalling notification message", slog.String("error_desc", err.Error()))
		return
	}
	err = api.p.ProduceWithJSON(raw)
	if err != nil {
		slog.Error("error sending notification message", slog.String("error_desc", err.Error()))
		return
	}
	slog.Info("successfuly sended new message notification", slog.Any("to", notification.To))
}

func (api *API) SendChatMessageNotification(mails []string, eventID primitive.ObjectID) {
	event, err := api.em.GetEventByID(eventID)
	if err != nil {
		slog.Error("fetching event db error", slog.String("error_desc", err.Error()))
		return
	}
	var msg models.Notification
	msg.To = mails
	msg.Subject = "В чате события новое сообщение"
	msg.Content = fmt.Sprintf("Пользователь, в чате события \"%s\" новое сообщение.\nПроверьте, вдруг это важно))", event.Name)
	api.sendNotification(&msg)
}

func (api *API) SendUpdateNotification(mails []string, eventID primitive.ObjectID) {
	event, err := api.em.GetEventByID(eventID)
	if err != nil {
		slog.Error("fetching event db error", slog.String("error_desc", err.Error()))
		return
	}
	var msg models.Notification
	msg.To = mails
	msg.Subject = "Событие было обновлено"
	msg.Content = fmt.Sprintf("Пользователь, событие \"%s\" было обновлено.\nПроверьте систему.", event.Name)
	api.sendNotification(&msg)
}

func (api *API) SendDelayedEventStartNotification(eventID primitive.ObjectID, deadline time.Time) error {
	event, err := api.em.GetEventByID(eventID)
	if err != nil {
		return errors.New("error getting event: " + err.Error())
	}
	emails, err := api.um.GetEmails(event.ParticipantsUUIDS())
	if err != nil {
		return errors.New("error getting mails: " + err.Error())
	}
	var msg models.Notification
	msg.To = emails
	msg.Subject = "Уведомление о начале события"
	msg.Content = fmt.Sprintf("Пользователь, событие \"%s\" скоро начнётся.", event.Name)
	msg.Delayed = true
	msg.Deadline = deadline
	raw, err := sonic.Marshal(msg)
	if err != nil {
		return errors.New("marshalling notification error: " + err.Error())
	}
	err = api.p.ProduceWithJSON(raw)
	if err != nil {
		return errors.New("producing message error: " + err.Error())
	}
	return nil
}
