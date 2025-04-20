package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	usermanager "github.com/lim-bo/calendar/backend/internal/user_manager"
	"github.com/lim-bo/calendar/backend/internal/util"
	"github.com/lim-bo/calendar/backend/models"
	"github.com/spf13/viper"
)

var (
	ErrProccessingBody = errors.New("proccessing body error")
	ErrBadRequest      = errors.New("invalid request")
	ErrRepository      = errors.New("error while fetching db")
	ErrLogin           = errors.New("unregistered or wrong credentials")
)

type UserManagerI interface {
	Register(creds *models.UserCredentialsRegister) error
	Login(creds *models.UserCredentials) (uuid.UUID, error)
	UpdateUser(newCreds *models.UserCredentialsRegister, uid uuid.UUID) error
	ChangePassword(newPass string, uid uuid.UUID) error
	GetProfileInfo(uid uuid.UUID) (*models.UserCredentialsRegister, error)
}

type API struct {
	r  *chi.Mux
	um UserManagerI
}

func New() *API {
	err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg := usermanager.DBConfig{
		Host:     viper.GetString("users_db_host"),
		Port:     viper.GetString("users_db_port"),
		DBName:   viper.GetString("users_db_name"),
		User:     viper.GetString("users_db_user"),
		Password: viper.GetString("users_db_pass"),
	}
	return &API{
		r:  chi.NewMux(),
		um: usermanager.New(cfg),
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
}

func (api *API) Run() error {
	host, port := viper.GetString("api_host"), viper.GetString("api_port")
	fmt.Printf("server started at %s:%s\n", host, port)
	return http.ListenAndServe(host+":"+port, api.r)
}
