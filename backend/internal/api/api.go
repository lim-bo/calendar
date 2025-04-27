package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	eventmanager "github.com/lim-bo/calendar/backend/internal/event_manager"
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
)

type UserManagerI interface {
	Register(creds *models.UserCredentialsRegister) error
	Login(creds *models.UserCredentials) (uuid.UUID, error)
	UpdateUser(newCreds *models.UserCredentialsRegister, uid uuid.UUID) error
	ChangePassword(newPass string, uid uuid.UUID) error
	GetProfileInfo(uid uuid.UUID) (*models.UserCredentialsRegister, error)
	GetUUIDS(mails []string) ([]uuid.UUID, error)
}

type EventManagerI interface {
	AddEvent(event *models.Event) error
	GetEvents(master uuid.UUID) ([]*models.Event, error)
	GetEventsByMonth(master uuid.UUID, month time.Month) ([]*models.Event, error)
	DeleteEvent(master uuid.UUID, id primitive.ObjectID) error
	GetEventsByWeek(master uuid.UUID) ([]*models.Event, error)
	GetEventsByDay(master uuid.UUID, day time.Time) ([]*models.Event, error)
}

type API struct {
	r  *chi.Mux
	um UserManagerI
	em EventManagerI
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
	return &API{
		r:  chi.NewMux(),
		um: usermanager.New(usersDBcfg),
		em: eventmanager.New(eventDBcfg),
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
	})
}

func (api *API) Run() error {
	host, port := viper.GetString("api_host"), viper.GetString("api_port")
	fmt.Printf("server started at %s:%s\n", host, port)
	return http.ListenAndServe(host+":"+port, api.r)
}
