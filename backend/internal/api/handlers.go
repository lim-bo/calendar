package api

import (
	"log/slog"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	usermanager "github.com/lim-bo/calendar/backend/internal/user_manager"
	"github.com/lim-bo/calendar/backend/models"
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
	slog.Debug("creds", slog.Any("value", creds))
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
	var event models.Event
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&event)
	defer r.Body.Close()
	if err != nil {
		slog.Error("error unmarshalling json", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/add"))
		w.WriteHeader(http.StatusBadRequest)
		WriteErrorResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}
	err = api.em.AddEvent(&event)
	if err != nil {
		slog.Error("event insertion error", slog.String("error_desc", err.Error()), slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/add"))
		w.WriteHeader(http.StatusInternalServerError)
		WriteErrorResponse(w, http.StatusInternalServerError, ErrRepository)
	}
	slog.Info("event successfuly added", slog.String("from", r.RemoteAddr), slog.String("endpoint", "events/add"), slog.String("uid", event.Master.String()))
}

func (api *API) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
