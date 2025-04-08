package api

import (
	"log/slog"
	"net/http"

	"github.com/bytedance/sonic"
	usermanager "github.com/lim-bo/calendar/backend/internal/user_manager"
	"github.com/lim-bo/calendar/backend/models"
)

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var creds models.UserCredentials
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&creds)
	defer r.Body.Close()
	if err != nil {
		slog.Error("unmarshalling body error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrProccessingBody)
		return
	}
	uid, err := api.um.Login(&creds)
	if err != nil {
		if err != usermanager.ErrUnregistered || err != usermanager.ErrWrongPass {
			slog.Error("login request internal error", slog.String("error_value", err.Error()), slog.String("from", r.RemoteAddr))
			w.WriteHeader(http.StatusInternalServerError)
			WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
			return
		}
		slog.Error("request with unregistered user or wrong credentials", slog.String("error_value", err.Error()), slog.String("from", r.RemoteAddr))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrLogin)
		return
	}
	slog.Info("successful login request", slog.String("from", r.RemoteAddr), slog.String("uid", uid.String()))
	WriteLoginResponse(w, uid)
}

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	var creds models.UserCredentialsRegister
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&creds)
	defer r.Body.Close()
	if err != nil {
		slog.Error("unmarshalling body error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrProccessingBody)
		return
	}
	slog.Debug("creds", slog.Any("value", creds))
	err = api.um.Register(&creds)
	if err != nil {
		if err == usermanager.ErrRegistered {
			slog.Error("incoming register request with already registered creds", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr))
			w.WriteHeader(http.StatusBadRequest)
			WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
			return
		}
		slog.Error("internal error while registration", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, usermanager.ErrInternal)
		return
	}
	slog.Info("successful registration", slog.String("from", r.RemoteAddr))
}

func (api *API) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		next.ServeHTTP(w, r)
	})
}
