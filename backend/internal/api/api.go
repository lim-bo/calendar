package api

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	attachmanager "github.com/lim-bo/calendar/backend/internal/attachments_manager"
	eventmanager "github.com/lim-bo/calendar/backend/internal/event_manager"
	"github.com/lim-bo/calendar/backend/internal/rabbit"
	usermanager "github.com/lim-bo/calendar/backend/internal/user_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrProccessingBody = errors.New("proccessing body error")
	ErrBadRequest      = errors.New("invalid request")
	ErrRepository      = errors.New("error while fetching db")
	ErrLogin           = errors.New("unregistered or wrong credentials")
	ErrResponse        = errors.New("error while responsing")
	ErrInvalidEmail    = errors.New("invalid email")
)

type UserManagerI interface {
	Register(creds *models.UserCredentialsRegister) error
	Login(creds *models.UserCredentials) (uuid.UUID, error)
	UpdateUser(newCreds *models.UserCredentialsRegister, uid uuid.UUID) error
	ChangePassword(newPass string, uid uuid.UUID) error
	GetProfileInfo(uid uuid.UUID) (*models.UserCredentialsRegister, error)
	GetUUIDS(mails []string) ([]uuid.UUID, error)
	GetEmails(uids []uuid.UUID) ([]string, error)
}

type EventManagerI interface {
	AddEvent(event *models.Event) error
	GetEvents(master uuid.UUID) ([]*models.Event, error)
	GetEventsByMonth(master uuid.UUID, month time.Month) ([]*models.Event, error)
	DeleteEvent(master uuid.UUID, id primitive.ObjectID) error
	UpdateEvent(event *models.Event) error
	GetEventsByWeek(master uuid.UUID) ([]*models.Event, error)
	GetEventsByDay(master uuid.UUID, day time.Time) ([]*models.Event, error)
	GetEventByID(id primitive.ObjectID) (*models.Event, error)
	ChangeUserAcceptance(eventID primitive.ObjectID, uid uuid.UUID, accepted bool) error
	GetPartsList(eventID primitive.ObjectID) ([]models.Participant, error)

	DeleteChat(eventID primitive.ObjectID) error
	GetMessages(eventID primitive.ObjectID) (*models.Chat, error)
	SendMessage(eventID primitive.ObjectID, msg *models.MessageWithMail) error
}

type AttachmentsManagerI interface {
	LoadAttachment(eventID primitive.ObjectID, file *models.FileLoad) error
	GetAttachments(eventID primitive.ObjectID) ([]*models.FileDownload, error)
}

type MQProducerI interface {
	ProduceWithJSON(jsonMsg []byte) error
}

type API struct {
	r              *chi.Mux
	um             UserManagerI
	em             EventManagerI
	am             AttachmentsManagerI
	p              MQProducerI
	producerCancel func()
}

func New() *API {
	err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	usersDBcfg := usermanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	eventDBcfg := eventmanager.DBConfig{
		Host:     viper.GetString("events_db_host"),
		Port:     viper.GetString("events_db_port"),
		User:     viper.GetString("events_db_user"),
		Password: viper.GetString("events_db_pass"),
	}
	s3cfg := attachmanager.MinioCfg{
		Address:    viper.GetString("minio_addr"),
		User:       viper.GetString("minio_user"),
		Pass:       viper.GetString("minio_pass"),
		BucketName: viper.GetString("minio_bucket"),
	}
	sqlcfg := attachmanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	rabbitCfg := rabbit.RabbitCfg{
		Host:     viper.GetString("rabbit_host"),
		Port:     viper.GetString("rabbit_port"),
		Username: viper.GetString("rabbit_user"),
		Password: viper.GetString("rabbit_pass"),
	}
	prod, cancel := rabbit.NewProducer(rabbitCfg, "notifications")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		for i := range sigCh {
			i.String()
			cancel()
			return
		}
	}()
	return &API{
		r:              chi.NewMux(),
		um:             usermanager.New(usersDBcfg),
		em:             eventmanager.New(eventDBcfg),
		am:             attachmanager.New(&s3cfg, &sqlcfg),
		p:              prod,
		producerCancel: cancel,
	}
}

func (api *API) MountEndpoint() {
	api.r.Route("/users", func(r chi.Router) {
		r.Use(api.CORSMiddleware)
		r.Post("/login", api.Login)
		r.Post("/register", api.Register)
		r.Post("/{uid}/update", api.UpdateUser)
		r.Post("/{uid}/changepass", api.ChangePassword)
		r.Get("/{uid}/profile", api.GetUserInfo)
	})
	api.r.Route("/events", func(r chi.Router) {
		r.Use(api.CORSMiddleware)
		r.Post("/add", api.AddEvent)
		r.Get("/{uid}/month", api.GetEventsByMonth)
		r.Get("/{uid}/week", api.GetEventsByWeek)
		r.Get("/{uid}/day", api.GetEventsByDay)
		r.Delete("/{uid}/delete", api.DeleteEvent)
		r.Post("/update", api.UpdateEvent)
		r.Post("/{eventID}/{uid}", api.ChangeParticipantState)
		r.Get("/{eventID}/parts", api.GetEventParticipants)
	})
	api.r.Route("/chats", func(r chi.Router) {
		r.Use(api.CORSMiddleware)
		r.Post("/{eventID}", api.SendMessage)
		r.Get("/{eventID}", api.GetMessages)
	})
	api.r.Route("/attachs", func(r chi.Router) {
		r.Use(api.CORSMiddleware)
		r.Post("/{eventID}", api.LoadAttachment)
		r.Get("/{eventID}", api.GetAttachments)
	})
}

func (api *API) Run() error {
	defer api.producerCancel()
	host, port := viper.GetString("api_host"), viper.GetString("api_port")
	fmt.Printf("server started at %s:%s\n", host, port)
	return http.ListenAndServe(host+":"+port, api.r)
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
	raw, err := sonic.Marshal(msg)
	if err != nil {
		slog.Error("error marshalling notification message", slog.String("error_desc", err.Error()))
		return
	}
	err = api.p.ProduceWithJSON(raw)
	if err != nil {
		slog.Error("error sending notification message", slog.String("error_desc", err.Error()))
		return
	}
	slog.Info("successfuly sended new message notification", slog.Any("to", mails), slog.String("eventID", string(eventID.Hex())))
}
