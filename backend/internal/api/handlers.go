package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	eventmanager "github.com/lim-bo/calendar/backend/internal/event_manager"
	usermanager "github.com/lim-bo/calendar/backend/internal/user_manager"
	"github.com/lim-bo/calendar/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.UserCredentials
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&creds)
	defer r.Body.Close()
	if err != nil {
		slog.Error("unmarshalling body error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/login"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrProccessingBody)
		return
	}
	uid, err := api.um.Login(&creds)
	if err != nil {
		if err != usermanager.ErrUnregistered || err != usermanager.ErrWrongPass {
			slog.Error("login request internal error", slog.String("error_value", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/login"))
			w.WriteHeader(http.StatusInternalServerError)
			WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
			return
		}
		slog.Error("request with unregistered user or wrong credentials", slog.String("error_value", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/login"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrLogin)
		return
	}
	slog.Info("successful login request", slog.String("from", r.RemoteAddr), slog.String("uid", uid.String()), slog.String("endpoint", "/users/login"))
	WriteLoginResponse(w, uid)
}

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	var creds models.UserCredentialsRegister
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&creds)
	defer r.Body.Close()
	if err != nil {
		slog.Error("unmarshalling body error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/register"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrProccessingBody)
		return
	}
	if !ValidateEmail(creds.Email) {
		slog.Error("register request with invalid email", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/register"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidEmail)
		return
	}
	err = api.um.Register(&creds)
	if err != nil {
		if err == usermanager.ErrRegistered {
			slog.Error("incoming register request with already registered creds", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/login"))
			w.WriteHeader(http.StatusBadRequest)
			WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
			return
		}
		slog.Error("internal error while registration", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/register"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, usermanager.ErrInternal)
		return
	}
	slog.Info("successful registration", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/register"))
}

func (api *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	uidStr := r.PathValue("uid")
	if uidStr == "" {
		slog.Error("lack of uid in pathvalues", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/update"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		slog.Error("wrong uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/update"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	var newCreds models.UserCredentialsRegister
	err = sonic.ConfigDefault.NewDecoder(r.Body).Decode(&newCreds)
	defer r.Body.Close()
	if err != nil {
		slog.Error("unmarshalling body error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr),
			slog.String("endpoint", "/users/{uid}/update"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrProccessingBody)
		return
	}
	err = api.um.UpdateUser(&newCreds, uid)
	if err != nil {
		slog.Error("repository error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr),
			slog.String("endpoint", "/users/{uid}/update"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("successful update", slog.String("uid", uidStr), slog.String("from", r.RemoteAddr),
		slog.String("endpoint", "/users/{uid}/update"))
}

func (api *API) ChangePassword(w http.ResponseWriter, r *http.Request) {
	uidStr := r.PathValue("uid")
	if uidStr == "" {
		slog.Error("lack of uid in pathvalues", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/changepass"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		slog.Error("wrong uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/changepass"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	newPass := make(map[string]interface{}, 1)
	err = sonic.ConfigDefault.NewDecoder(r.Body).Decode(&newPass)
	defer r.Body.Close()
	if err != nil {
		slog.Error("error unmarshalling json", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/changepass"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	pass, ok := newPass["pass"]
	if !ok {
		slog.Error("lack of new password in request body", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/changepass"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	passStr, ok := pass.(string)
	if !ok {
		slog.Error("incompatible data in request body", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/changepass"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	err = api.um.ChangePassword(passStr, uid)
	if err != nil {
		slog.Error("repository error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr),
			slog.String("endpoint", "/users/{uid}/changepass"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("successful changing password", slog.String("uid", uidStr), slog.String("from", r.RemoteAddr),
		slog.String("endpoint", "/users/{uid}/changepass"))
}

func (api *API) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	uidStr := r.PathValue("uid")
	if uidStr == "" {
		slog.Error("lack of uid in pathvalues", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/profile"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		slog.Error("wrong uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/profile"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	info, err := api.um.GetProfileInfo(uid)
	if err != nil {
		slog.Error("getting profile error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/profile"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	WriteGetProfileResponse(w, info)
	slog.Info("successful get profile request", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/users/{uid}/profile"), slog.String("uid", uid.String()))

}

func (api *API) AddEvent(w http.ResponseWriter, r *http.Request) {
	var eventRequest models.EventWithMails
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&eventRequest)
	defer r.Body.Close()
	if err != nil {
		slog.Error("error unmarshalling json", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/add"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	var event models.Event
	uids, err := api.um.GetUUIDS(eventRequest.Participants)
	if err != nil {
		slog.Error("getting uids error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/add"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	var parts []models.Participant
	for _, uid := range uids {
		event.Participants = append(event.Participants, models.Participant{UID: uid, Accepted: false})
	}
	event.EventBase = eventRequest.EventBase
	event.Participants = append(parts, models.Participant{UID: eventRequest.Master, Accepted: false})
	err = api.em.AddEvent(&event)
	if err != nil {
		slog.Error("event insertion error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/add"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("event successfuly added", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/add"), slog.String("uid", event.Master.String()))
}

func (api *API) GetEventsByMonth(w http.ResponseWriter, r *http.Request) {
	uidStr := r.PathValue("uid")
	if uidStr == "" {
		slog.Error("lack of uid in pathvalues", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/month"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		slog.Error("wrong uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/month"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil || month < 1 || month > 12 {
		slog.Error("invalid query month param", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/month"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	events, err := api.em.GetEventsByMonth(uid, time.Month(month))
	if err != nil {
		slog.Error("getting events by month error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/month"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	err = sonic.ConfigDefault.NewEncoder(w).Encode(map[string]interface{}{"events": events, "cod": 200})
	if err != nil {
		slog.Error("error marshalling events result", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/month"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrResponse)
		return
	}
	slog.Info("successfuly provided events data", slog.String("uid", uidStr), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/month"))
}

func (api *API) GetEventsByWeek(w http.ResponseWriter, r *http.Request) {
	uidStr := r.PathValue("uid")
	if uidStr == "" {
		slog.Error("lack of uid in pathvalues", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/week"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		slog.Error("wrong uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/week"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	events, err := api.em.GetEventsByWeek(uid)
	if err != nil {
		slog.Error("getting events by week error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/week"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	err = sonic.ConfigDefault.NewEncoder(w).Encode(map[string]interface{}{"events": events, "cod": 200})
	if err != nil {
		slog.Error("error marshalling events result", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/week"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrResponse)
		return
	}
	slog.Info("successfuly provided events data", slog.String("uid", uidStr), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/week"))

}

func (api *API) GetEventsByDay(w http.ResponseWriter, r *http.Request) {
	uidStr := r.PathValue("uid")
	if uidStr == "" {
		slog.Error("lack of uid in pathvalues", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/day"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		slog.Error("wrong uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/day"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	day, err := time.Parse("2006-01-02", r.URL.Query().Get("day"))
	if err != nil {
		slog.Error("invalid day query param", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/day"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	events, err := api.em.GetEventsByDay(uid, day)
	if err != nil {
		slog.Error("error getting events by day", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/day"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	err = sonic.ConfigDefault.NewEncoder(w).Encode(map[string]interface{}{"events": events, "cod": 200})
	if err != nil {
		slog.Error("error marshalling events result", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/day"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrResponse)
		return
	}
	slog.Info("successfuly provided events by day", slog.String("uid", uidStr), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/day"))
}

func (api *API) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	uidStr := r.PathValue("uid")
	if uidStr == "" {
		slog.Error("lack of uid in pathvalues", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/delete"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		slog.Error("wrong uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/delete"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	obj := r.URL.Query().Get("id")
	if obj == "" {
		slog.Error("request with lack of event id", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/delete"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(obj)
	if err != nil {
		slog.Error("request with invalid event id", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/delete"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
	}
	err = api.em.DeleteEvent(uid, id)
	if err != nil {
		if err == eventmanager.ErrLackOrWrongMaster {
			slog.Error("deletion request with unexisted event", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/delete"))
			w.WriteHeader(http.StatusBadRequest)
			WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
			return
		}
		slog.Error("deletion event error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/delete"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("successful event deletion", slog.String("from", r.RemoteAddr), slog.String("endpoint", "/events/{uid}/delete"), slog.String("uid", uidStr))
}

func (api *API) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var eventRequest models.EventWithMails
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&eventRequest)
	defer r.Body.Close()
	if err != nil {
		slog.Error("error unmarshalling json", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/update"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	var event models.Event
	uids, err := api.um.GetUUIDS(eventRequest.Participants)
	if err != nil {
		slog.Error("getting uids error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/update"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	var parts []models.Participant
	for i, uid := range uids {
		event.Participants = append(event.Participants, models.Participant{UID: uid, Accepted: event.Participants[i].Accepted})
	}
	event.EventBase = eventRequest.EventBase
	event.Participants = append(parts, models.Participant{UID: eventRequest.Master, Accepted: false})
	//TO-DO: add notification
	err = api.em.UpdateEvent(&event)
	if err != nil {
		slog.Error("updating event error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/update"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("successfully updated event", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/update"))
}

func (api *API) SendMessage(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	if eventID == "" {
		slog.Error("send message request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	eventIDPrim, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		slog.Error("send message request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	var msg models.Message
	err = sonic.ConfigDefault.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		slog.Error("error unmarshalling json", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	defer r.Body.Close()

	mails, err := api.um.GetEmails([]uuid.UUID{msg.Sender})
	if err != nil {
		slog.Error("error mapping uid to mail", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	err = api.em.SendMessage(eventIDPrim, &models.MessageWithMail{
		Sender:  mails[0],
		Content: msg.Content,
	})
	if err != nil {
		slog.Error("error sending message", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	go api.SendChatMessageNotification(mails, eventIDPrim)
	slog.Info("successfuly sended message to event chat", slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"), slog.String("event", eventID))
}

func (api *API) GetMessages(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	if eventID == "" {
		slog.Error("get messages request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	eventIDPrim, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		slog.Error("get messages request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	chat, err := api.em.GetMessages(eventIDPrim)
	if err != nil {
		slog.Error("getting messages error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	err = sonic.ConfigDefault.NewEncoder(w).Encode(chat)
	if err != nil {
		slog.Error("error marshalling content", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrResponse)
		return
	}
	slog.Info("successfuly provided messages", slog.String("from", r.RemoteAddr), slog.String("endpoint", "chats/{eventID}"), slog.String("event", eventID))
}

func (api *API) LoadAttachment(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	if eventID == "" {
		slog.Error("load attachment request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	eventIDPrim, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		slog.Error("load attachment request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	var file models.FileLoad
	err = sonic.ConfigDefault.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		slog.Error("error unmarshalling json", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	err = api.am.LoadAttachment(eventIDPrim, &file)
	if err != nil {
		slog.Error("loading attachment error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("successfuly loaded attachment", slog.String("event", eventID), slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
}

func (api *API) GetAttachments(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	if eventID == "" {
		slog.Error("download attachment request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	eventIDPrim, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		slog.Error("download attachment request with invalid eventid", slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	result, err := api.am.GetAttachments(eventIDPrim)
	if err != nil {
		slog.Error("getting attachments for event error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	err = sonic.ConfigDefault.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("marshalling result error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrResponse)
		return
	}
	slog.Info("successfuly provided attachments for event", slog.String("event", eventID), slog.String("from", r.RemoteAddr), slog.String("endpoint", "attachs/{eventID}"))
}

func (api *API) ChangeParticipantState(w http.ResponseWriter, r *http.Request) {
	eventID, err := primitive.ObjectIDFromHex(r.PathValue("eventID"))
	if err != nil {
		slog.Error("state changing request with invalid eventID in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/{uid}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	uid, err := uuid.Parse(r.PathValue("uid"))
	if err != nil {
		slog.Error("state changing request with invalid uid in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/{uid}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	state := r.URL.Query().Get("state")
	if state == "" {
		slog.Error("state changing request with lack of query params", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/{uid}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	var acceptance bool
	switch state {
	case "1":
		acceptance = true
	case "0":
		acceptance = false
	default:
		slog.Error("state changing request with invalid query param", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/{uid}"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	err = api.em.ChangeUserAcceptance(eventID, uid, acceptance)
	if err != nil {
		slog.Error("error changing user state", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/{uid}"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("successfully changed state", slog.String("from", r.RemoteAddr),
		slog.String("endpoint", "events/{eventID}/{uid}"),
		slog.String("uid", uid.String()),
		slog.String("eventID", eventID.Hex()),
	)
}

func (api *API) GetEventParticipants(w http.ResponseWriter, r *http.Request) {
	eventID, err := primitive.ObjectIDFromHex(r.PathValue("eventID"))
	if err != nil {
		slog.Error("state changing request with invalid eventID in path", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/parts"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	partsWithUUIDS, err := api.em.GetPartsList(eventID)
	if err != nil {
		slog.Error("getting participants error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/parts"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	var uuids []uuid.UUID
	for _, p := range partsWithUUIDS {
		uuids = append(uuids, p.UID)
	}
	var partsWithEmails []models.ParticipantWithEmail
	emails, err := api.um.GetEmails(uuids)
	if err != nil {
		slog.Error("mapping emails error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/parts"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	for i, p := range partsWithUUIDS {
		partsWithEmails = append(partsWithEmails, models.ParticipantWithEmail{
			Accepted: p.Accepted,
			Email:    emails[i],
		})
	}
	err = sonic.ConfigDefault.NewEncoder(w).Encode(map[string]interface{}{
		"cod":   200,
		"parts": partsWithEmails,
	})
	if err != nil {
		slog.Error("marshalling results error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/parts"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
		return
	}
	slog.Info("successfully provided participants list", slog.String("event", eventID.Hex()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/{eventID}/parts"))
}

func (api *API) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
